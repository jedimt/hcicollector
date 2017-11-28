package config

import (
	"github.com/cblomart/vsphere-graphite/backend"
	"github.com/cblomart/vsphere-graphite/vsphere"
)

// Configuration : configurarion base
type Configuration struct {
	VCenters     []*vsphere.VCenter
	Metrics      []*vsphere.Metric
	Interval     int
	Domain       string
	Backend      *backend.Backend
	CPUProfiling bool
	MEMProfiling bool
	FlushSize    int
}
