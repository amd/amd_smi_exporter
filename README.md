AMD SMI Prometheus Exporter
----------------------------

The AMD SMI Exporter is a standalone app that can be run as a daemon, written in GO Language,
that exports AMD CPU & GPU metrics to the Prometheus server. The AMD SMI Prometheus Exporter
employs
* [AMDSMI Library](https://github.com/ROCm/amdsmi.git) for its data acquisition and 
* [GO binding](https://github.com/ROCm/amdsmi/blob/amd-staging/goamdsmi.go) that provides an interface between the amdsmi and the GO exporter code.

Note:
----
This [AMD SMI Exporter](https://github.com/amd/amd_smi_exporter) repository will no longer receive further updates. 
Moving forward, the [AMD Device Metrics Exporter](https://github.com/rocm/device-metrics-exporter) repository will serve as the official replacement. 
Please transition to [AMD Device Metrics Exporter](https://github.com/rocm/device-metrics-exporter) for continued improvements, updates, and support.

### Important note about Versioning and Backward Compatibility

The AMD SMI Exporter follows the AMDSMI library in its releases, as it is dependent on the underlying libraries for its data. The Exporter is currently under development, and therefore subject to change in the features it offers and at the interface with the GO binding.

While every effort will be made to ensure stable backward compatibility in software releases
with a major version greater than 0, any code/interface may be subject to revision/change while the major version remains 0.

## Building the exporter

The standalone GO Exporter may be built from the src directory as follows:

### Downloading the source

The source code for the GO Exporter is [AMD SMI Exporter](https://github.com/amd/amd_smi_exporter.git).

### Directory stucture of the source

Once the exporter source has been cloned to a local Linux machine, the directory structure of
source is as below:
* `$ src/` Contains exporter source for package main
* `$ src/collect/` Contains the implementation of the Scan function of the collector.
* `$ grafana/` Contains the JSON files for Grafana dashboard.

* Change the directory to amd_smi_exporter/src

	```$ cd amd_smi_exporter/src```

* Execute "make clean" to clean pre-existing binaries and GO module files

	```amd_smi_exporter/src$ make clean```

* Execute "make" to perform a "go get" of dependent modules such as
	* github.com/prometheus/client_golang
	* github.com/prometheus/client_golang/prometheus
	* github.com/prometheus/client_golang/prometheus/promhttp
	* github.com/ROCm/amdsmi

	```amd_smi_exporter/src$ make```

The aforementioned steps will create the "amd_smi_exporter" GO binary file.
To install the binary in /usr/local/bin, and install the service file in
/etc/systemd/system directory, one may execute:

	```$ sudo make install```

## Building the container for the GO Exporter

Once the GO Exporter is built, one may proceed to create a containerized micro service of the go executable by executing the following commands:

Prerequisite: docker version 20.10.12 or later must be installed on the build server for the
container build to succeed.

* Execute "make container_clean" to clean pre-existing images and configuration of the container image.

	```amd_smi_exporter/src$ make container_clean```

* Build the fresh container image with the following command:

	```amd_smi_exporter/src$ make container```

  This command will build the container image and will be listed when the user issues the
```sudo docker images``` command.
A tarball of the container image file "k8/amd_smi_exporter_container.tar" is also saved in
the "k8" directory, and this may be used to deploy the container manually on respective
nodes of the kubernetes cluster using the "k8/daemonset.yaml" file.

## Grafana Dashboard:

JSON files for Grafana dashboard are available under grafana/ of this repo
* AMDSmiExporter_CPU_GrafanaDashboard.json
* AMDSmiExporter_GPU_GrafanaDashboard.json

## Dependencies

Please ensure the following are in place
1. amdsmi library with  goamdsmi_shim bindings installed under "/opt/rocm"
2. GO v1.20
3. Docker (tested on v20.10.12 or later)

### GO Installation:

To run on AMD rocm dockers, GO installation through apt install on Linux is only supported till 1.18.
Manual installation can be done from here: <https://go.dev/dl/>
Below is an example of installing 1.20.12 of go.

	$ wget -L "https://golang.org/dl/go1.20.12.linux-amd64.tar.gz"
	$ tar -xf "go1.20.12.linux-amd64.tar.gz"
	$ cd go/
	$ ls -l
	$ cd ..
	$ sudo chown -R root:root ./go
	$ sudo mv -v go /usr/local
	$ export GOPATH=$HOME/go
	$ export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
	

Add amdsmi library path to LD_LIBRARY_PATH environment variable and export.

	$ export LD_LIBRARY_PATH=<path_to_amdsmi_library>

## Running the GO Exporter

NOTE: Only one instance of the GO Exporter may be run on the server, either as a 
standalone service, or as a containerized micro service (started with "docker run"
or as a daemonSet of a kubernetes deployment).

Prerequisite: To ensure that AMD custom parameters defined in the 
amd-smi-custom-rules.yml file are found in the promql queries, add 
the following rule_files and scrape_configs to the 
/etc/prometheus/prometheus.yml file:

rule_files:
  - "amd-smi-custom-rules.yml"

scrape_configs:
  - job_name: "prometheus"
  - job_name: "amd-smi-exporter"
    static_configs:
      - targets: ["localhost:2021"]

### Custom rules

The prometheus query language allows the user to customize his queries based on user requirements.
The customizations may be added to the /usr/local/bin/prometheus/amd-smi-custom-rules.yml file".
Here are a few sample queries that may be built over the aforementioned objects:

* > ### amd_core_energy{thread="101"}/1000000
	Displays the core energy of core 101 shifted by six decimal points.

* > ### amd_socket_power/100 > 650.00
	Rule to check if socket power consumption has gone over 650.00

* > ### amd_prochot_status != 0
	Alert to check if PROC_HOT status has been triggered

## Executing the Go Exporter

The GO exporter may be run manually in the following ways

### 1. Executing the "amd_smi_exporter" GO binary:

	```amd_smi_exporter/src$ ./amd_smi_exporter```

### 2. As a systemd daemon:

	```$ sudo systemctl daemon-reload```

	```$ sudo service prometheus restart```

	```$ sudo service amd-smi-exporter start```

### 3. As a containerized micro service that may be started manually or as a kubernetes daemonSet:

Assuming user has a running docker daemon and a kubernetes cluster.

#### On a server node that is not a part of a kubernetes cluster, one may execute the following command:

	```$ sudo docker run -d --name amd-exporter --device=/dev/cpu --device=/dev/kfd
           --device=/dev/dri --privileged -p 2021:2021 amd_smi_exporter_container:0.1```

   Alternatively, the docker image tarball of the container may be copied to individual
   kubernetes cluster node and loaded on the worker node. The daemonSet may then be applied
   from the master node as follows:

#### On the worker node, copy the amd_smi_exporter_container.tar image file and execute:

	```$ sudo docker load -i amd_smi_exporter_container.tar```

	On the master node, copy the daemonset.yaml file and execute:

	```$ kubectl apply -f daemonset.yaml```

	This will deploy a single running instance of the AMD SMI Exporter container micro
	service on the worker nodes of the kubernetes cluster. The daemonset.yaml file may
	be edited to apply taints for nodes where the exporter is not expected to run in
	the cluster.

## Supported hardware

AMD EPYC TM line of server CPU Families:

* AMD CPU Family `19h` Models `0h-Fh` (Milan), `10h-1Fh` (Genoa), `A0h-AFh`. 
* AMD CPU Family `1Ah` Models `0h-Fh` (Turin), `10h-1Fh`.
* AMD APU Family `19h` Models `90h-9fh` and 
* AMD GPUs MI200 and MI300.

## Examples

### CPU core metrics

#### 1. amd_core_energy
	### Description: Displays the per-core energy consumption of the processor so far.
This object may be queried at the core level or the thread level. The values reported by
the threads in a hyperthreaded core will be the same. This object query will report the
energy counter values for all threads. To query a single thread (lets say the thread number
is 101), the user may use the following query:

	amd_core_energy{thread="101"}

	### Type: Counter
	### Property: Read-only

#### 2. amd_boost_limit
	### Description: Displays the per-core boost limit that the core is operating at.
	### Type: Gauge
	### Property: Read-only

### CPU Socket metrics

#### 3. amd_socket_energy
	### Description: Displays the per-socket cumulative energy consumed by all cores
so far. This value excludes the energy consumed by the AID (Active Interposer Die).To query
a single socket (lets say socket 2), the user may use the following query:

	amd_socket_energy{socket="2"}

	### Type: Counter
	### Property: Read-only

#### 4. amd_socket_power
	### Description: Displays the per-socket power consumed. This is a real time gauge
value that is queried at a time interval set by the scrape interval.
	### Type: Gauge
	### Property: Read-only

#### 5. amd_power_limit
	### Description: Displays the power limit at which the processor is operating at.
	### Type: Gauge
	### Property: Read-only

#### 6. amd_prochot_status
	### Description: Displays a binary value of "0" or "1", where "1" implies that the
PROC_HOT status of the processor has been triggered.
	### Type: Gauge
	### Property: Read-only

### System

#### 7. amd_num_sockets
	### Description: Displays the number of sockets which the processor is seated in.
	### Type: Gauge
	### Property: Read-only

#### 8. amd_num_threads
	### Description: Displays the total number of threads (logical CPUs) in all.
	### Type: Gauge
	### Property: Read-only

### 9. amd_num_threads_per_core
	### Description: Displays the number of threads (logical CPUs) per core.
	### Type: Gauge
	### Property: Read-only

### GPU Metrics

#### 10. amd_num_gpus
	### Description: Displays the number of gpus
	### Type: Gauge
	### Property: Read-only

#### 11. amd_gpu_dev_id
	### Description: Displays the dev id of the gpu
	### Type: Gauge
	### Property: Read-only

#### 12. amd_gpu_power_cap
	### Description: Displays the gpu power cap
	### Type: Gauge
	### Property: Read-only

#### 13. amd_gpu_power_avg
	### Description: Displays the gpu average power consumed
	### Type: Counter
	### Property: Read-only

#### 14. amd_gpu_current_temperature
	### Description: Displays the current temperature of the gpu
	### Type: Gauge
	### Property: Read-only

#### 15. amd_gpu_SCLK
	### Description: Displays the GPU SCLK frequency
	### Type: Gauge
	### Property: Read-only

#### 16. amd_gpu_MCLK
	### Description: Displays the GPU MCLK frequency
	### Type: Gauge
	### Property: Read-only

#### 17. amd_gpu_Usage
        ### Description: Displays the GPU Use percent
        ### Type: Gauge
        ### Property: Read-only

#### 18. amd_gpu_memory_busy percent
        ### Description: Displays the GPU Memory busy percent
        ### Type: Gauge
        ### Property: Read-only

## FAQs:

* If the prometheus service fails to start properly,
   run the command ```journalctl -u prometheus -f --no-pager``` 

* If an issue is related to "Web lister busy" or "Port is already in use",
  Please change Port from 9090 to 9091 in the following files

	* /etc/systemd/system/prometheus.service file
		- under line "--web.listen-address=0.0.0.0:9090"
	* /etc/prometheus/prometheus.yml file
		- under line "targets: ["localhost:9090"]

	and restart the systemd service using command "service prometheus restart".


