# MIT-X11 Open Source License
# Copyright (c) 2024, Advanced Micro Devices, Inc.
build:
	make clean
	go mod init src
	go get golang.org/x/exp/slices
	go get github.com/prometheus/client_golang
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promhttp
	go get github.com/amd/go_amd_smi@master
	go build -o amd_smi_exporter main.go cpu_data.go

clean:
	rm -rf amd_smi_exporter go.mod go.sum

install:
	sudo cp amd_smi_exporter /usr/local/bin
	sudo cp amd-smi-exporter.service /etc/systemd/system
	sudo cp amd-smi-custom-rules.yml /etc/prometheus

container:
	mkdir -p ./install/library
	mkdir -p ./install/binconfig
	sudo cp -r /opt/goamdsmi ./install/library
	sudo cp -r /opt/e-sms ./install/library
	sudo cp -r /opt/rocm ./install/library
	sudo cp amd_smi_exporter ./install/binconfig
	sudo docker build -t amd_smi_exporter_container:0.1 .
	sudo docker save -o k8/amd_smi_exporter_container.tar amd_smi_exporter_container:0.1

container_clean:
	sudo rm -rf install
	sudo rm -rf k8/amd_smi_exporter_container.tar
	sudo docker rmi amd_smi_exporter_container:0.1
