package backend

import (
	"bytes"
	"errors"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cblomart/vsphere-graphite/backend/ThinInfluxClient"
	influxclient "github.com/influxdata/influxdb/client/v2"
	"github.com/marpaia/graphite-golang"
)

// Point : Information collected for a point
type Point struct {
	VCenter      string   `influx:"tag,name"`
	ObjectType   string   `influx:"tag,type"`
	ObjectName   string   `influx:"tag,name"`
	Group        string   `influx:"key,1"`
	Counter      string   `influx:"key,2"`
	Instance     string   `influx:"tag,instance"`
	Rollup       string   `influx:"key,3"`
	Value        int64    `influx:"value"`
	Datastore    []string `influx:"tag,datastore"`
	ESXi         string   `influx:"tag,host"`
	Cluster      string   `influx:"tag,cluster"`
	Network      []string `influx:"tag,network"`
	ResourcePool string   `influx:"tag,resourcepool"`
	Folder       string   `influx:"tag,folder"`
	ViTags       []string `influx:"tag,vitags"`
	NumCPU       int32    `influx:"tag,numcpu"`
	MemorySizeMB int32    `influx:"tag,memorysizemb"`
	Timestamp    int64    `influx:"time"`
}

// InfluxPoint is the representation of the parts of a point for influx
type InfluxPoint struct {
	Key       string
	Fields    map[string]string
	Tags      map[string]string
	Timestamp int64
}

// Backend : storage backend
type Backend struct {
	Hostname     string
	ValueField   string
	Database     string
	Username     string
	Password     string
	Type         string
	Port         int
	NoArray      bool
	Encrypted    bool
	carbon       *graphite.Graphite
	influx       *influxclient.Client
	thininfluxdb *ThinInfluxClient.ThinInfluxClient
}

const (
	// Graphite name of the graphite backend
	Graphite = "graphite"
	// InfluxDB name of the influx db backend
	InfluxDB = "influxdb"
	// ThinInfluxDB name of the thin influx db backend
	ThinInfluxDB = "thininfluxdb"
	// InfluxTag is the tag for influxdb
	InfluxTag = "influx"
)

var stdlog, errlog *log.Logger

// StringMaptoString converts a string map to csv or get the first value
func StringMaptoString(value []string, separator string, noarray bool) string {
	if len(value) == 0 {
		return ""
	}
	if noarray {
		return value[0]
	}
	return strings.Join(value, separator)
}

// IntMaptoString converts a int map to csv or get the first value
func IntMaptoString(value []int, separator string, noarray bool) string {
	if len(value) == 0 {
		return ""
	}
	if noarray {
		return strconv.Itoa(value[0])
	}
	var strval []string
	for _, i := range value {
		strval = append(strval, strconv.Itoa(i))
	}
	return strings.Join(strval, separator)
}

// Int32MaptoString converts a int32 map to csv or get the first value
func Int32MaptoString(value []int32, separator string, noarray bool) string {
	if len(value) == 0 {
		return ""
	}
	if noarray {
		return strconv.FormatInt(int64(value[0]), 10)
	}
	var strval []string
	for _, i := range value {
		strval = append(strval, strconv.FormatInt(int64(i), 10))
	}
	return strings.Join(strval, separator)
}

// Int64MaptoString converts a int64 map to csv or get the first value
func Int64MaptoString(value []int64, separator string, noarray bool) string {
	if len(value) == 0 {
		return ""
	}
	if noarray {
		return strconv.FormatInt(value[0], 10)
	}
	var strval []string
	for _, i := range value {
		strval = append(strval, strconv.FormatInt(i, 10))
	}
	return strings.Join(strval, separator)
}

// ValToString : try to convert interface to string. Separated by separator if slice
func ValToString(value interface{}, separator string, noarray bool) string {
	switch v := value.(type) {
	case string:
		return v
	case []string:
		return StringMaptoString(v, separator, noarray)
	case int:
		return strconv.Itoa(v)
	case []int:
		return IntMaptoString(v, separator, noarray)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case []int32:
		return Int32MaptoString(v, separator, noarray)
	case int64:
		return strconv.FormatInt(v, 10)
	case []int64:
		return Int64MaptoString(v, separator, noarray)
	default:
		return ""
	}
}

// Join map[int]string into a string
func Join(values map[int]string, separator string) string {
	var keys []int
	for k := range values {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// create a map with the key parts in order
	var tmp []string
	for _, k := range keys {
		tmp = append(tmp, values[k])
	}
	return strings.Join(tmp, separator)
}

// MustAtoi converts a string to integer and return 0 i case of error
func MustAtoi(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		i = 0
	}
	return i
}

// GetInfluxPoint : convert a point to an influxpoint
func (p *Point) GetInfluxPoint(noarray bool, valuefield string) *InfluxPoint {
	keyParts := make(map[int]string)
	ip := InfluxPoint{
		Fields: make(map[string]string),
		Tags:   make(map[string]string),
	}
	v := reflect.ValueOf(p).Elem()
	for i := 0; i < v.NumField(); i++ {
		vfield := v.Field(i)
		tfield := v.Type().Field(i)
		tag := tfield.Tag.Get(InfluxTag)
		tagfields := strings.Split(tag, ",")
		if len(tagfields) == 0 || len(tagfields) > 2 {
			stdlog.Println("tag field ignored: " + tag)
			continue
		}
		tagtype := tagfields[0]
		tagname := strings.ToLower(tfield.Name)
		if len(tagfields) == 2 {
			tagname = tagfields[1]
		}
		switch tagtype {
		case "key":
			keyParts[MustAtoi(tagname)] = ValToString(vfield.Interface(), "_", false)
		case "tag":
			ip.Tags[tagname] = ValToString(vfield.Interface(), "\\,", noarray)
		case "value":
			ip.Fields[valuefield] = ValToString(vfield.Interface(), ",", true) + "i"
		case "time":
			ip.Timestamp = vfield.Int()
		default:
		}
	}
	// sort key part keys and join them
	ip.Key = Join(keyParts, "_")
	return &ip
}

// ConvertToKV converts a map[string]string to a csv with k=v pairs
func ConvertToKV(values map[string]string) string {
	var tmp []string
	for key, val := range values {
		if len(val) == 0 {
			continue
		}
		tbuf := bytes.NewBuffer(nil)
		tbuf.WriteString(key)
		tbuf.WriteRune('=')
		tbuf.WriteString(val)
		tmp = append(tmp, tbuf.String())
	}
	return strings.Join(tmp, ",")
}

// ToInflux converts the influx point to influx string format
func (ip *InfluxPoint) ToInflux(noarray bool, valuefield string) string {
	// buffer containing the resulting line
	buff := bytes.NewBuffer(nil)
	// key of the mesurement
	buff.WriteString(ip.Key)
	buff.WriteRune(',')
	// Tags
	buff.WriteString(ConvertToKV(ip.Tags))
	// separator
	buff.WriteRune(' ')
	// fields
	buff.WriteString(ConvertToKV(ip.Fields))
	// separator
	buff.WriteRune(' ')
	// timestamp
	buff.WriteString(strconv.FormatInt(ip.Timestamp, 10))
	return buff.String()
}

// ToInflux serialises the data to be consumed by influx line protocol
// see https://docs.influxdata.com/influxdb/v1.2/write_protocols/line_protocol_tutorial/
func (p *Point) ToInflux(noarray bool, valuefield string) string {
	return p.GetInfluxPoint(noarray, valuefield).ToInflux(noarray, valuefield)
}

// Init : initialize a backend
func (backend *Backend) Init(standardLogs *log.Logger, errorLogs *log.Logger) error {
	stdlog = standardLogs
	errlog = errorLogs
	if len(backend.ValueField) == 0 {
		// for compatibility reason with previous version
		// can now be changed in the config file.
		// the default can later be changed to another value.
		// most probably "value" (lower case)
		backend.ValueField = "Value"
	}
	switch backendType := strings.ToLower(backend.Type); backendType {
	case Graphite:
		// Initialize Graphite
		stdlog.Println("Intializing " + backendType + " backend")
		carbon, err := graphite.NewGraphite(backend.Hostname, backend.Port)
		if err != nil {
			errlog.Println("Error connecting to graphite")
			return err
		}
		backend.carbon = carbon
		return nil
	case InfluxDB:
		//Initialize Influx DB
		stdlog.Println("Intializing " + backendType + " backend")
		influxclt, err := influxclient.NewHTTPClient(influxclient.HTTPConfig{
			Addr:     "http://" + backend.Hostname + ":" + strconv.Itoa(backend.Port),
			Username: backend.Username,
			Password: backend.Password,
		})
		if err != nil {
			errlog.Println("Error connecting to InfluxDB")
			return err
		}
		backend.influx = &influxclt
		return nil
	case ThinInfluxDB:
		//Initialize thin Influx DB client
		stdlog.Println("Initializing " + backendType + " backend")
		thininfluxclt, err := ThinInfluxClient.NewThinInlfuxClient(backend.Hostname, backend.Port, backend.Database, backend.Username, backend.Password, "s", backend.Encrypted)
		if err != nil {
			errlog.Println("Error creating thin InfluxDB client")
			return err
		}
		backend.thininfluxdb = &thininfluxclt
		return nil
	default:
		errlog.Println("Backend " + backendType + " unknown.")
		return errors.New("Backend " + backendType + " unknown.")
	}
}

// Disconnect : disconnect from backend
func (backend *Backend) Disconnect() {
	switch backendType := strings.ToLower(backend.Type); backendType {
	case Graphite:
		// Disconnect from graphite
		stdlog.Println("Disconnecting from graphite")
		err := backend.carbon.Disconnect()
		if err != nil {
			errlog.Println("Error disconnecting from graphite: ", err)
		}
	case InfluxDB:
		// Disconnect from influxdb
		stdlog.Println("Disconnecting from influxdb")
	case ThinInfluxDB:
		// Disconnect from thin influx db
		errlog.Println("Disconnecting from thininfluxdb")
	default:
		errlog.Println("Backend " + backendType + " unknown.")
	}
}

// SendMetrics : send metrics to backend
func (backend *Backend) SendMetrics(metrics []*Point) {
	switch backendType := strings.ToLower(backend.Type); backendType {
	case Graphite:
		var graphiteMetrics []graphite.Metric
		for _, point := range metrics {
			if point == nil {
				continue
			}
			//key := "vsphere." + vcName + "." + entityName + "." + name + "." + metricName
			key := "vsphere." + point.VCenter + "." + point.ObjectType + "." + point.ObjectName + "." + point.Group + "." + point.Counter + "." + point.Rollup
			if len(point.Instance) > 0 {
				key += "." + strings.ToLower(strings.Replace(point.Instance, ".", "_", -1))
			}
			graphiteMetrics = append(graphiteMetrics, graphite.Metric{Name: key, Value: strconv.FormatInt(point.Value, 10), Timestamp: point.Timestamp})
		}
		err := backend.carbon.SendMetrics(graphiteMetrics)
		if err != nil {
			errlog.Println("Error sending metrics (trying to reconnect): ", err)
			err := backend.carbon.Connect()
			if err != nil {
				errlog.Println("could not connect to graphite: ", err)
			}
		}
	case InfluxDB:
		//Influx batch points
		bp, err := influxclient.NewBatchPoints(influxclient.BatchPointsConfig{
			Database:  backend.Database,
			Precision: "s",
		})
		if err != nil {
			errlog.Println("Error creating influx batchpoint")
			errlog.Println(err)
			return
		}
		for _, point := range metrics {
			if point == nil {
				continue
			}
			key := point.Group + "_" + point.Counter + "_" + point.Rollup
			tags := map[string]string{}
			tags["vcenter"] = point.VCenter
			tags["type"] = point.ObjectType
			tags["name"] = point.ObjectName
			if backend.NoArray {
				if len(point.Datastore) > 0 {
					tags["datastore"] = point.Datastore[0]
				}
			} else {
				if len(point.Datastore) > 0 {
					tags["datastore"] = strings.Join(point.Datastore, "\\,")
				}
			}
			if backend.NoArray {
				if len(point.Network) > 0 {
					tags["network"] = point.Network[0]
				}
			} else {
				if len(point.Network) > 0 {
					tags["network"] = strings.Join(point.Network, "\\,")
				}
			}
			if len(point.ESXi) > 0 {
				tags["host"] = point.ESXi
			}
			if len(point.Cluster) > 0 {
				tags["cluster"] = point.Cluster
			}
			if len(point.Instance) > 0 {
				tags["instance"] = point.Instance
			}
			if len(point.ResourcePool) > 0 {
				tags["resourcepool"] = point.ResourcePool
			}
			if len(point.Folder) > 0 {
				tags["folder"] = point.Folder
			}
			if backend.NoArray {
				if len(point.ViTags) > 0 {
					tags["vitags"] = point.ViTags[0]
				}
			} else {
				if len(point.ViTags) > 0 {
					tags["vitags"] = strings.Join(point.ViTags, "\\,")
				}
			}
			if point.NumCPU != 0 {
				tags["numcpu"] = strconv.FormatInt(int64(point.NumCPU), 10)
			}
			if point.MemorySizeMB != 0 {
				tags["memorysizemb"] = strconv.FormatInt(int64(point.MemorySizeMB), 10)
			}
			fields := make(map[string]interface{})
			fields[backend.ValueField] = point.Value
			pt, err := influxclient.NewPoint(key, tags, fields, time.Unix(point.Timestamp, 0)) // nolint: vetshadow
			if err != nil {
				errlog.Println("Could not create influxdb point")
				errlog.Println(err)
				continue
			}
			bp.AddPoint(pt)
		}
		err = (*backend.influx).Write(bp)
		if err != nil {
			errlog.Println("Error sending metrics: ", err)
		}
	case ThinInfluxDB:
		lines := []string{}
		for _, point := range metrics {
			if point == nil {
				continue
			}
			lines = append(lines, point.ToInflux(backend.NoArray, backend.ValueField))
		}
		count := 3
		for count > 0 {
			err := backend.thininfluxdb.Send(lines)
			if err != nil {
				errlog.Println("Error sending metrics: ", err)
				if err.Error() == "Server Busy: timeout" {
					errlog.Println("waiting .5 second to continue")
					time.Sleep(500 * time.Millisecond)
					count--
				} else {
					break
				}
			} else {
				break
			}
		}
		err := backend.thininfluxdb.Send(lines)
		if err != nil {
			errlog.Println("Error sendg metrics: ", err)
		}
	default:
		errlog.Println("Backend " + backendType + " unknown.")
	}
}
