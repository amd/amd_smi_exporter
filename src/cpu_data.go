/*
 * MIT-X11 Open Source License
 *
 * Copyright (c) 2022, Advanced Micro Devices, Inc.
 * All rights reserved.
 *
 * Developed by:
 *
 *                 AMD Research and AMD Software Development
 *
 *                 Advanced Micro Devices, Inc.
 *
 *                 www.amd.com
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or
 * sellcopies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 *  - The above copyright notice and this permission notice shall be included in
 *    all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 *
 * Except as contained in this notice, the name of the Advanced Micro Devices,
 * Inc. shall not be used in advertising or otherwise to promote the sale, use
 * or other dealings in this Software without prior written authorization from
 * the Advanced Micro Devices, Inc.
 *
 */

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
	NumGPUs *prometheus.Desc
	GPUDevId *prometheus.Desc
	GPUPowerCap *prometheus.Desc
	GPUPower *prometheus.Desc
	GPUTemperature *prometheus.Desc
	GPUSCLK *prometheus.Desc
	GPUMCLK *prometheus.Desc
	GPUUsage *prometheus.Desc
	GPUMemoryUsage *prometheus.Desc
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
			[]string{"socket"},// The metric's variable label dimensions.
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
			[]string{"socket"},// The metric's variable label dimensions.
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
		NumGPUs: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "num_gpus"),
			"AMD Params",// The metric's help text.
			[]string{"num_gpus"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUDevId: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "gpu_dev_id"),
			"AMD Params",// The metric's help text.
			[]string{"gpu_dev_id", "productname"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUPowerCap: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "gpu_power_cap"),
			"AMD Params",// The metric's help text.
			[]string{"gpu_power_cap", "productname"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUPower: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "gpu_power"),
			"AMD Params",// The metric's help text.
			[]string{"gpu_power", "productname"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUTemperature: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "gpu_current_temperature"),
			"AMD Params",// The metric's help text.
			[]string{"gpu_current_temperature", "productname"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUSCLK: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "gpu_SCLK"),
			"AMD Params",// The metric's help text.
			[]string{"gpu_SCLK", "productname"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUMCLK: prometheus.NewDesc(
			prometheus.BuildFQName("amd", "", "gpu_MCLK"),
			"AMD Params",// The metric's help text.
			[]string{"gpu_MCLK", "productname"},// The metric's variable label dimensions.
			nil,// The metric's constant label dimensions.
		),
		GPUUsage: prometheus.NewDesc(
                        prometheus.BuildFQName("amd", "", "gpu_use_percent"),
                        "AMD Params",// The metric's help text.
                        []string{"gpu_use_percent", "productname"},// The metric's variable label dimensions.
                        nil,// The metric's constant label dimensions.
                ),
                GPUMemoryUsage: prometheus.NewDesc(
                        prometheus.BuildFQName("amd", "", "gpu_memory_use_percent"),
                        "AMD Params",// The metric's help text.
                        []string{"gpu_memory_use_percent", "productname"},// The metric's variable label dimensions.
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

	if 0 != data.Threads {
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
	}

	if 0 != data.Sockets {
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
	}

	if 0 != data.NumGPUs {
		for i,s := range data.GPUDevId{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUDevId,
				prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUPowerCap{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUPowerCap,
				prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUPower{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUPower,
				prometheus.CounterValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUTemperature{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUTemperature,
				prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUSCLK{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUSCLK,
				prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUMCLK{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUMCLK,
				prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUUsage{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUUsage,
			prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}

		for i,s := range data.GPUMemoryUsage{
			if uint(i) > (data.NumGPUs - 1) {
				continue
			}
			ch <- prometheus.MustNewConstMetric(c.GPUMemoryUsage,
			prometheus.GaugeValue, float64(s), strconv.Itoa(i), gGPUProductNames[i])
		}
	}

	ch <- prometheus.MustNewConstMetric(c.Sockets,
		prometheus.GaugeValue, float64(data.Sockets), "")
	ch <- prometheus.MustNewConstMetric(c.Threads,
		prometheus.GaugeValue, float64(data.Threads), "")
	ch <- prometheus.MustNewConstMetric(c.ThreadsPerCore,
		prometheus.GaugeValue, float64(data.ThreadsPerCore), "")
	ch <- prometheus.MustNewConstMetric(c.NumGPUs,
		prometheus.GaugeValue, float64(data.NumGPUs), "")
}
