package main

import (
	"strconv"
	"src/collect" // This has the implementation of the Scan() function
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = &amd_data{}

type amd_data struct {
	DataDesc *prometheus.Desc
	CoreEnergy *prometheus.Desc
	SocketEnergy *prometheus.Desc
	BoostLimit *prometheus.Desc
	SocketPower *prometheus.Desc
	PowerLimit *prometheus.Desc
	ProchotStatus *prometheus.Desc
	Sockets *prometheus.Desc
	Threads *prometheus.Desc
	ThreadsPerCore *prometheus.Desc
	Data func() (collect.AMDParams)
}

func NewCollector(handle func() (collect.AMDParams)) prometheus.Collector {
	return &amd_data{
		DataDesc: prometheus.NewDesc(
			"amd_data",// Name of the metric.
			"AMD Params",// The metric's help text.
			[]string{"socket"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		CoreEnergy: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "core_energy"),
			"AMD Params",// The metric's help text.
			[]string{"thread"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		SocketEnergy: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "socket_energy"),
			"AMD Params",// The metric's help text.
			[]string{"socket_energy"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		BoostLimit: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "boost_limit"),
			"AMD Params",// The metric's help text.
			[]string{"thread"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		SocketPower: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "socket_power"),
			"AMD Params",// The metric's help text.
			[]string{"socket_power"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		PowerLimit: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "power_limit"),
			"AMD Params",// The metric's help text.
			[]string{"power_limit"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		ProchotStatus: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "prochot_status"),
			"AMD Params",// The metric's help text.
			[]string{"prochot_status"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		Sockets: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "num_sockets"),
			"AMD Params",// The metric's help text.
			[]string{"num_sockets"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		Threads: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "num_threads"),
			"AMD Params",// The metric's help text.
			[]string{"num_threads"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		ThreadsPerCore: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "num_threads_per_core"),
			"AMD Params",// The metric's help text.
			[]string{"num_threads_per_core"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),

		Data: handle, //This is the Scan() function handle
	}
}

func (c *amd_data) Describe(ch chan<- *prometheus.Desc) {

	ds := []*prometheus.Desc{
		c.DataDesc,
	}

	for _, d := range ds {
		ch <- d
	}
}

func (c *amd_data) Collect(ch chan<- prometheus.Metric) {

	data := c.Data() //Call the Scan() function here and get AMDParams

	for i,s := range data.CoreEnergy{
		if uint(i) > (data.Threads - 1) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.CoreEnergy,
			prometheus.CounterValue, float64(s), strconv.Itoa(i))
	}

	for i,s := range data.CoreBoost{
		if uint(i) > (data.Threads - 1) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.BoostLimit,
			prometheus.GaugeValue, float64(s), strconv.Itoa(i))
	}

	for i,s := range data.SocketEnergy{
		if uint(i) > (data.Sockets - 1) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.SocketEnergy,
			prometheus.CounterValue, float64(s), strconv.Itoa(i))
	}

	for i,s := range data.SocketPower{
		if uint(i) > (data.Sockets - 1) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.SocketPower,
			prometheus.GaugeValue, float64(s), strconv.Itoa(i))
	}

	for i,s := range data.PowerLimit{
		if uint(i) > (data.Sockets - 1) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.PowerLimit,
			prometheus.GaugeValue, float64(s), strconv.Itoa(i))
	}

	for i,s := range data.ProchotStatus{
		if uint(i) > (data.Sockets - 1) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.ProchotStatus,
			prometheus.GaugeValue, float64(s), strconv.Itoa(i))
	}

	ch <- prometheus.MustNewConstMetric(c.Sockets,
		prometheus.GaugeValue, float64(data.Sockets), "")
	ch <- prometheus.MustNewConstMetric(c.Threads,
		prometheus.GaugeValue, float64(data.Threads), "")
	ch <- prometheus.MustNewConstMetric(c.ThreadsPerCore,
		prometheus.GaugeValue, float64(data.ThreadsPerCore), "")
}
