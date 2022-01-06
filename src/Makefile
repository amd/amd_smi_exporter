build:
	make clean
	go mod init src
	go get github.com/prometheus/client_golang
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promhttp
	go get https://github.com/amd/goesmi
	go build -o amd_exporter main.go cpu_data.go

clean:
	rm -rf amd_exporter go.mod go.sum

install:
	sudo cp amd_exporter /usr/local/bin
	sudo cp amd-exporter.service /etc/systemd/system
