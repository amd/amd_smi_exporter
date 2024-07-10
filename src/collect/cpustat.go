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

// Package cpustat provides an example parser for Linux CPU utilization statistics.
package collect
import "github.com/amd/go_amd_smi"

var UINT16_MAX = uint16(0xFFFF)
var UINT32_MAX = uint32(0xFFFFFFFF)
var UINT64_MAX = uint64(0xFFFFFFFFFFFFFFFF)

type AMDParams struct {
	CoreEnergy [768]float64
	CoreBoost [768]float64
	SocketEnergy [8]float64
	SocketPower [8]float64
	PowerLimit [8]float64
	ProchotStatus [8]float64
	Sockets uint
	Threads uint
	ThreadsPerCore uint
	NumGPUs uint
	GPUDevId [24]float64
	GPUPowerCap [24]float64
	GPUPower [24]float64
	GPUTemperature [24]float64
	GPUSCLK [24]float64
	GPUMCLK [24]float64
	GPUUsage [24]float64
	GPUMemoryUsage [24]float64
}

func (amdParams *AMDParams) Init() {
	amdParams.Sockets = 0
	amdParams.Threads = 0
	amdParams.ThreadsPerCore = 0

	amdParams.NumGPUs = 0

	for socketLoopCounter := 0; socketLoopCounter < len(amdParams.SocketEnergy); socketLoopCounter++	{
		amdParams.SocketEnergy[socketLoopCounter] = -1
		amdParams.SocketPower[socketLoopCounter] = -1
		amdParams.PowerLimit[socketLoopCounter] = -1
		amdParams.ProchotStatus[socketLoopCounter] = -1
	}

	for logicalCoreLoopCounter := 0; logicalCoreLoopCounter < len(amdParams.CoreEnergy); logicalCoreLoopCounter++	{
		amdParams.CoreEnergy[logicalCoreLoopCounter] = -1
		amdParams.CoreBoost[logicalCoreLoopCounter] = -1
	}

	for gpuLoopCounter := 0; gpuLoopCounter < len(amdParams.GPUDevId); gpuLoopCounter++	{
		amdParams.GPUDevId[gpuLoopCounter] = -1
		amdParams.GPUPowerCap[gpuLoopCounter] = -1
		amdParams.GPUPower[gpuLoopCounter] = -1
		amdParams.GPUTemperature[gpuLoopCounter] = -1
		amdParams.GPUSCLK[gpuLoopCounter] = -1
		amdParams.GPUMCLK[gpuLoopCounter] = -1
		amdParams.GPUUsage[gpuLoopCounter] = -1
		amdParams.GPUMemoryUsage[gpuLoopCounter] = -1
	}
}

func Scan() (AMDParams) {

	var stat AMDParams
	stat.Init()

	value64 := uint64(0)
	value32 := uint32(0)
	value16 := uint16(0)

	if true == goamdsmi.GO_cpu_init() {

		num_sockets := int(goamdsmi.GO_cpu_number_of_sockets_get())
		num_threads := int(goamdsmi.GO_cpu_number_of_threads_get())
		num_threads_per_core := int(goamdsmi.GO_cpu_threads_per_core_get())

		stat.Sockets = uint(num_sockets)
		stat.Threads = uint(num_threads)
		stat.ThreadsPerCore = uint(num_threads_per_core)

		for i := 0; i < num_threads ; i++ {
			value64 = uint64(goamdsmi.GO_cpu_core_energy_get(i))
			if UINT64_MAX != value64 { stat.CoreEnergy[i] = float64(value64) }
			value64 = 0

			value32 = uint32(goamdsmi.GO_cpu_core_boostlimit_get(i))
			if UINT32_MAX != value32 { stat.CoreBoost[i] = float64(value32) }
			value32 = 0
		}

		for i := 0; i < num_sockets ; i++ {
			value64 = uint64(goamdsmi.GO_cpu_socket_energy_get(i))
			if UINT64_MAX != value64 { stat.SocketEnergy[i] = float64(value64) }
			value64 = 0

			value32 = uint32(goamdsmi.GO_cpu_socket_power_get(i))
			if UINT32_MAX != value32 { stat.SocketPower[i] = float64(value32) }
			value32 = 0

			value32 = uint32(goamdsmi.GO_cpu_socket_power_cap_get(i))
			if UINT32_MAX != value32 { stat.PowerLimit[i] = float64(value32) }
			value32 = 0

			value32 = uint32(goamdsmi.GO_cpu_prochot_status_get(i))
			if UINT32_MAX != value32 { stat.ProchotStatus[i] = float64(value32) }
			value32 = 0
		}
	}


	if true == goamdsmi.GO_gpu_init() {

		num_gpus := int(goamdsmi.GO_gpu_num_monitor_devices())
		stat.NumGPUs = uint(num_gpus)

		for i := 0; i < num_gpus ; i++ {
			value16 = uint16(goamdsmi.GO_gpu_dev_id_get(i))
			if UINT16_MAX != value16 { stat.GPUDevId[i] = float64(value16) }
			value16 = 0

			value64 = uint64(goamdsmi.GO_gpu_dev_power_cap_get(i))
			if UINT64_MAX != value64 { stat.GPUPowerCap[i] = float64(value64) }
			value64 = 0

			value64 = uint64(goamdsmi.GO_gpu_dev_power_get(i))
			if UINT64_MAX != value64 { stat.GPUPower[i] = float64(value64) }
			value64 = 0

			//Get the value for GPU current temperature. Sensor = 0(GPU), Metric = 0(current)
			value64 = uint64(goamdsmi.GO_gpu_dev_temp_metric_get(i, 0, 0))
			if UINT64_MAX == value64 {
				//Sensor = 1 (GPU Junction Temp)
				value64 = uint64(goamdsmi.GO_gpu_dev_temp_metric_get(i, 1, 0))
			}
			if UINT64_MAX != value64 { stat.GPUTemperature[i] = float64(value64) }
			value64 = 0

			value64 = uint64(goamdsmi.GO_gpu_dev_gpu_clk_freq_get_sclk(i))
			if UINT64_MAX != value64 { stat.GPUSCLK[i] = float64(value64) }
			value64 = 0

			value64 = uint64(goamdsmi.GO_gpu_dev_gpu_clk_freq_get_mclk(i))
			if UINT64_MAX != value64 { stat.GPUMCLK[i] = float64(value64) }
			value64 = 0

			value32 = uint32(goamdsmi.GO_gpu_dev_gpu_busy_percent_get(i))
			if UINT32_MAX != value32 { stat.GPUUsage[i] = float64(value32) }
			value32 = 0

			value64 = uint64(goamdsmi.GO_gpu_dev_gpu_memory_busy_percent_get(i))
			if UINT64_MAX != value64 { stat.GPUMemoryUsage[i] = float64(value64) }
			value64 = 0
		}
	}

	return stat
}

