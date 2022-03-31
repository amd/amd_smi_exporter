build:
	make clean
	go mod init src
	go get github.com/prometheus/client_golang
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promhttp
	go get github.com/amd/go_amd_smi
	go build -o amd_smi_exporter main.go cpu_data.go

clean:
	rm -rf amd_smi_exporter go.mod go.sum

install:
	sudo cp amd_smi_exporter /usr/local/bin
	sudo cp amd-smi-exporter.service /etc/systemd/system
	sudo cp amd-smi-custom-rules.yml /usr/local/bin/prometheus
