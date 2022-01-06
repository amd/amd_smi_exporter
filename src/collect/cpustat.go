// Package cpustat provides an example parser for Linux CPU utilization statistics.
package collect

import "https://github.com/amd/goesmi"

type AMDParams struct {
	CoreEnergy [256]float64
	SocketEnergy [8]float64
	CoreBoost [256]float64
	SocketPower [8]float64
	PowerLimit [8]float64
	ProchotStatus [8]float64
	Sockets uint
	Threads uint
	ThreadsPerCore uint
}

func Scan() (AMDParams) {

	var stat AMDParams
	value64 := uint64(0)
	value32 := uint32(0)

	if 1 == goesmi.GO_esmi_init() {
		num_sockets := int(goesmi.GO_esmi_number_of_sockets_get())
		num_threads := int(goesmi.GO_esmi_number_of_threads_get())
		num_threads_per_core := int(goesmi.GO_esmi_threads_per_core_get())

		stat.Sockets = uint(num_sockets)
		stat.Threads = uint(num_threads)
		stat.ThreadsPerCore = uint(num_threads_per_core)

		for i := 0; i < num_threads ; i++ {
			value64 = uint64(goesmi.GO_esmi_core_energy_get(i))
			stat.CoreEnergy[i] = float64(value64)
			value64 = 0

			value32 = uint32(goesmi.GO_esmi_core_boostlimit_get(i))
			stat.CoreBoost[i] = float64(value32)
			value32 = 0
		}

		for i := 0; i < num_sockets ; i++ {
			value64 = uint64(goesmi.GO_esmi_socket_energy_get(i))
			stat.SocketEnergy[i] = float64(value64)
			value64 = 0

			value32 = uint32(goesmi.GO_esmi_socket_power_get(i))
			stat.SocketPower[i] = float64(value32)
			value32 = 0

			value32 = uint32(goesmi.GO_esmi_socket_power_cap_get(i))
			stat.PowerLimit[i] = float64(value32)
			value32 = 0

			value32 = uint32(goesmi.GO_esmi_prochot_status_get(i))
			stat.ProchotStatus[i] = float64(value32)
			value32 = 0
		}
	}

	return stat
}
