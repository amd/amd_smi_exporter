1. [Introduction to the AMD SMI Exporter](#desc)
2. [Building the GO Exporter](#build)
3. [Library Dependencies](#lib)
4. [Kernel Dependencies](#kernel)
5. [Running the GO Exporter](#running)
6. [Supported Hardware](#hw)
7. [Software needed for the build](#sw)
8. [Objects exported by the AMD SMI Exporter](#objects)
a. [At the Core level](#core)
b. [At the Socket level](#socket)
c. [At the System level](#system)
d. [GPU Metrics](#gpu)
e. [Custom rules](#custom)

<a name="desc"></a>
# AMD SMI Prometheus Exporter

The AMD SMI Exporter is a standalone app that can be run as a daemon, written in GO Language,
that exports AMD CPU  & GPU metrics to the Prometheus server. The AMD SMI Prometheus Exporter
employs the [E-SMI In-Band C library](https://github.com/amd/esmi_ib_library.git) &
[ROCm SMI Library](https://github.com/RadeonOpenCompute/rocm_smi_lib.git) for its data
acquisition. The exporter and the E-SMI/ROCm-SMI library have a
[GO binding](https://github.com/amd/go_amd_smi.git) that provides an interface between the
e-smi,rocm-smi C,C++ library and the GO exporter code.

## Important note about Versioning and Backward Compatibility

The AMD SMI Exporter follows the E-SMI In-band library and the ROCm library in its releases,
as it is dependent on the underlying libraries for its data. The Exporter is currently under
development, and therefore subject to change in the features it offers and at the interface
with the GO binding.

While every effort will be made to ensure stable backward compatibility in software releases
with a major version greater than 0, any code/interface may be subject to rivision/change while
the major version remains 0.

<a name="build"></a>
# Building the GO Exporter

## Dowloading the source

The source code for the GO Exporter is available at [AMD SMI Exporter](https://github.com/amd/amd_smi_exporter.git).

## Directory stucture of the source

Once the exporter source has been cloned to a local Linux machine, the directory structure of
source is as below:
* `$ src/` Contains exporter source for package main
* `$ src/collect` Contains the implementation of the Scan function of the collector.

## Building

The GO Exporter may be built from the src directory as follows:

* Change the directory to amd_smi_exporter/src

	```$ cd amd_smi_exporter/src```

* Execute "make clean" to clean pre-existing binaries and GO module files

	```amd_smi_exporter/src$ make clean```

* Execute "make" to perform a "go get" of dependent modules such as
	* github.com/prometheus/client_golang
	* github.com/prometheus/client_golang/prometheus
	* github.com/prometheus/client_golang/prometheus/promhttp
	* github.com/amd/go_amd_smi

NOTE: Before executing the GO exporter as a standalone executable or as a service, one needs to ensure that the e-smi , goamdsmi_shim, and rocm-smi library dependencies are met. Please refer to the steps to build and install the library dependencies in the respective README of these repositories. The environment variable for the LD_LIBRARY_PATH is to be set to "/opt/e-sms/e_smi/lib:/opt/rocm/rocm_smi/lib:/opt/goamdsmi/lib". The user may edit this environment variable to reflect the installation path where the dependent libraries are installed.

	```amd_smi_exporter/src$ make```

The aforementioned steps will create the "amd_smi_exporter" GO binary file. To install the binary in
/usr/local/bin, and install the service file in /etc/systemd/system directory, one may execute:

	```$ sudo make install```

<a name="lib"></a>
# Library dependencies

Before executing the GO exporter as a standalone executable or as a service, one needs to ensure
that the e-smi , goamdsmi_shim, and rocm-smi library dependencies are met by ensuring that they are
installed in the "/opt/e-sms/e_smi/lib", "/opt/goamdsmi/lib" and "/opt/rocm/rocm_smi/lib" directories
respectively. Please refer to <https://github.com/amd/esmi_ib_library/docs/README.md>,
<https://github.com/amd/go_amd_smi/README.md>, and
<https://github.com/RadeonOpenCompute/rocm_smi_lib/README.md> for the build and installation
instructions.

<a name="kernel"></a>
# Kernel dependencies

The E-SMI Library, and inturn the GO exporter, depends on the following device drivers from Linux
to manage the system management features. The amd_hsmp driver is available in upstream kernel version 5.18. The driver git repo is as follows:

	* amd_hsmp <https://github.com/amd/amd_hsmp.git>
	* amd_energy <https://github.com/amd/amd_energy.git>

<a name="running"></a>
#Running the GO Exporter

1. The GO exporter may be run manually by executing the "amd_smi_exporter" GO binary

Prerequisite: Please ensure that the prometheus systemd service is installed in /etc/systemd/system/prometheus.service and that it is running with the configs specified in /usr/local/bin/prometheus/prometheus.yml.

	```$ ./amd_smi_exporter```

** OR **

2. The GO exporter may be started as a systemd daemon as follows:

Prerequisite: Edit the /usr/local/bin/prometheus/prometheus.yml and add the following rule_files and scrape_configs:
rule_files:
  - "amd-smi-custom-rules.yml"


scrape_configs:
  - job_name: "prometheus"
  - job_name: "amd-smi-exporter"
    static_configs:
      - targets: ["localhost:2021"]

	```$ sudo systemctl daemon-reload```
	```$ sudo service prometheus restart```
	```$ sudo service amd-smi-exporter start```

NOTE: The environment variable for the LD_LIBRARY_PATH is set to /opt/e-sms/e_smi/lib:/opt/rocm/rocm_smi/lib:/opt/goamdsmi/lib

<a name="hw"></a>
# Supported hardware

AMD Zen3 based CPU Family `19h` Models `0h-Fh` and `30h-3Fh`, and `17h` Model `30h`.

<a name="sw"></a>
# Additional required software for building

In order to build the GO Exporter, the following components are required. Note that the software versions
listed are what is being used in development. Earlier versions are not guaranteed to work:

* go1.17.3 linux/amd64


<a name="core"></a>
## At the Core level

### 1. amd_core_energy
	### Description: Displays the per-core energy consumption of the processor so far.
This object may be queried at the core level or the thread level. The values reported by
the threads in a hyperthreaded core will be the same. This object query will report the
energy counter values for all threads. To query a single thread (lets say the thread number
is 101), the user may use the following query:

	amd_core_energy{thread="101"}

	### Type: Counter
	### Property: Read-only

### 2. amd_boost_limit
	### Description: Displays the per-core boost limit that the core is operating at.
	### Type: Gauge
	### Property: Read-only

<a name="socket"></a>
## Socket

### 3. amd_socket_energy
	### Description: Displays the per-socket cumulative energy consumed by all cores
so far. This value excludes the energy consumed by the AID (Active Interposer Die).To query
a single socket (lets say socket 2), the user may use the following query:

	amd_socket_energy{socket="2"}

	### Type: Counter
	### Property: Read-only

### 4. amd_socket_power
	### Description: Displays the per-socket power consumed. This is a real time gauge
value that is queried at a time interval set by the scrape interval.
	### Type: Gauge
	### Property: Read-only

### 5. amd_power_limit
	### Description: Displays the power limit at which the processor is operating at.
	### Type: Gauge
	### Property: Read-only

### 6. amd_prochot_status
	### Description: Displays a binary value of "0" or "1", where "1" implies that the
PROC_HOT status of the processor has been triggered.
	### Type: Gauge
	### Property: Read-only

<a name="system"></a>
## System

### 7. amd_num_sockets
	### Description: Displays the number of sockets which the processor is seated in.
	### Type: Gauge
	### Property: Read-only

### 8. amd_num_threads
	### Description: Displays the total number of threads (logical CPUs) in all.
	### Type: Gauge
	### Property: Read-only

## 9. amd_num_threads_per_core
	### Description: Displays the number of threads (logical CPUs) per core.
	### Type: Gauge
	### Property: Read-only

<a name="gpu"></a>
## GPU Metrics

## 10. amd_num_gpus
	### Description: Displays the number of gpus
	### Type: Gauge
	### Property: Read-only

## 11. amd_gpu_dev_id
	### Description: Displays the dev id of the gpu
	### Type: Gauge
	### Property: Read-only

## 12. amd_gpu_power_cap
	### Description: Displays the gpu power cap
	### Type: Gauge
	### Property: Read-only

## 13. amd_gpu_power_avg
	### Description: Displays the gpu average power consumed
	### Type: Counter
	### Property: Read-only

## 14. amd_gpu_current_temperature
	### Description: Displays the current temperature of the gpu
	### Type: Gauge
	### Property: Read-only

## 15. amd_gpu_SCLK
	### Description: Displays the GPU SCLK frequency
	### Type: Gauge
	### Property: Read-only

## 16. amd_gpu_MCLK
	### Description: Displays the GPU MCLK frequency
	### Type: Gauge
	### Property: Read-only

<a name="custom"></a>
## Custom rules

The prometheus query language allows the user to customize his queries based on user requirements. The customizations may be added to the /usr/local/bin/prometheus/amd-smi-custom-rules.yml file".
Here are a few sample queries that may be built over the aforementioned objects:

* > ### amd_core_energy{thread="101"}/1000000
	Displays the core energy of core 101 shifted by six decimal points.

* > ### amd_socket_power/100 > 650.00
	Rule to check if socket power consumption has gone over 650.00

* > ### amd_prochot_status != 0
	Alert to check if PROC_HOT status has been triggered
