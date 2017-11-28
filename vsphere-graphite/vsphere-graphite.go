package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path"
	"reflect"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cblomart/vsphere-graphite/backend"
	"github.com/cblomart/vsphere-graphite/config"
	"github.com/cblomart/vsphere-graphite/vsphere"

	"github.com/takama/daemon"

	"code.cloudfoundry.org/bytefmt"
)

const (
	// name of the service
	name        = "vsphere-graphite"
	description = "send vsphere stats to graphite"
)

var dependencies = []string{}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

func queryVCenter(vcenter vsphere.VCenter, conf config.Configuration, channel *chan backend.Point) {
	vcenter.Query(conf.Interval, conf.Domain, channel)
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: vsphere-graphite install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	stdlog.Println("Starting daemon:", path.Base(os.Args[0]))

	// read the configuration
	file, err := os.Open("/etc/" + path.Base(os.Args[0]) + ".json")
	if err != nil {
		return "Could not open configuration file", err
	}
	jsondec := json.NewDecoder(file)
	conf := config.Configuration{}
	err = jsondec.Decode(&conf)
	if err != nil {
		return "Could not decode configuration file", err
	}

	if conf.FlushSize == 0 {
		//conf.FlushSize = 1000 
		conf.FlushSize = 200
	}

	if conf.CPUProfiling {
		f, err := ioutil.TempFile("/tmp", "vsphere-graphite-cpu.profile") // nolint: vetshadow
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		stdlog.Println("Will write cpu profiling to: ", f.Name())
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	//force backend values to environement varialbles if present
	s := reflect.ValueOf(conf.Backend).Elem()
	numfields := s.NumField()
	for i := 0; i < numfields; i++ {
		f := s.Field(i)
		if f.CanSet() {
			//exported field
			envname := strings.ToUpper(s.Type().Name() + "_" + s.Type().Field(i).Name)
			envval := os.Getenv(envname)
			if len(envval) > 0 {
				//environment variable set with name
				switch ftype := f.Type().Name(); ftype {
				case "string":
					f.SetString(envval)
				case "int":
					val, err := strconv.ParseInt(envval, 10, 64) // nolint: vetshadow
					if err == nil {
						f.SetInt(val)
					}
				}
			}
		}
	}

	for _, vcenter := range conf.VCenters {
		vcenter.Init(conf.Metrics, stdlog, errlog)
	}

	err = conf.Backend.Init(stdlog, errlog)
	if err != nil {
		return "Could not initialize backend", err
	}
	defer conf.Backend.Disconnect()

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM) // nolint: megacheck

	// Set up a channel to receive the metrics
	metrics := make(chan backend.Point, conf.FlushSize)

	// Set up a ticker to collect metrics at givent interval
	ticker := time.NewTicker(time.Second * time.Duration(conf.Interval))
	defer ticker.Stop()

	// Start retriveing and sending metrics
	stdlog.Println("Retrieving metrics")
	for _, vcenter := range conf.VCenters {
		go queryVCenter(*vcenter, conf, &metrics)
	}

	// Memory statisctics
	var memstats runtime.MemStats
	// timer to execute memory collection
	memtimer := time.NewTimer(time.Second * time.Duration(10))
	// Memory profiling
	var mf *os.File
	if conf.MEMProfiling {
		mf, err = ioutil.TempFile("/tmp", "vsphere-graphite-mem.profile")
		if err != nil {
			log.Fatal("could not create MEM profile: ", err)
		}
		defer mf.Close() // nolint: errcheck
	}
	// buffer for points to send
	pointbuffer := make([]*backend.Point, conf.FlushSize)
	bufferindex := 0

	for {
		select {
		case value := <-metrics:
			// reset timer as a point has been revieved
			if !memtimer.Stop() {
				select {
				case <-memtimer.C:
				default:
				}
			}
			memtimer.Reset(time.Second * time.Duration(5))
			pointbuffer[bufferindex] = &value
			bufferindex++
			if bufferindex == len(pointbuffer) {
				conf.Backend.SendMetrics(pointbuffer)
				stdlog.Printf("Sent %d logs to backend", bufferindex)
				ClearBuffer(pointbuffer)
				bufferindex = 0
			}
		case <-ticker.C:
			stdlog.Println("Retrieving metrics")
			for _, vcenter := range conf.VCenters {
				go queryVCenter(*vcenter, conf, &metrics)
			}
		case <-memtimer.C:
			// sent remaining values
			conf.Backend.SendMetrics(pointbuffer)
			stdlog.Printf("Sent %d logs to backend", bufferindex)
			bufferindex = 0
			ClearBuffer(pointbuffer)
			runtime.GC()
			debug.FreeOSMemory()
			runtime.ReadMemStats(&memstats)
			stdlog.Printf("Memory usage : sys=%s alloc=%s\n", bytefmt.ByteSize(memstats.Sys), bytefmt.ByteSize(memstats.Alloc))
			if conf.MEMProfiling {
				stdlog.Println("Writing mem profiling to: ", mf.Name())
				debug.WriteHeapDump(mf.Fd())
			}
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if bufferindex > 0 {
				conf.Backend.SendMetrics(pointbuffer[:bufferindex])
				stdlog.Printf("Sent %d logs to backend", bufferindex)
			}
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

// ClearBuffer : set all values in pointer array to nil
func ClearBuffer(buffer []*backend.Point) {
	for i := 0; i < len(buffer); i++ {
		buffer[i] = nil
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	srv, err := daemon.New(name, description, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "Error: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
