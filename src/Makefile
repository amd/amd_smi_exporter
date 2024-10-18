# MIT-X11 Open Source License
# Copyright (c) 2024, Advanced Micro Devices, Inc.
build:
	make clean
	go mod init src
	go get golang.org/x/exp/slices
	go get github.com/prometheus/client_golang
	go get github.com/prometheus/client_golang/prometheus
	go get github.com/prometheus/client_golang/prometheus/promhttp
	go get github.com/ROCm/amdsmi@amd-staging
	go build -o amd_smi_exporter main.go cpu_data.go

clean:
	rm -rf amd_smi_exporter go.mod go.sum

install: copy_exporter_rules
	
	sudo cp amd_smi_exporter /usr/local/bin
	sudo cp amd-smi-exporter.service /etc/systemd/system
	sudo mkdir -p /usr/share/doc/amd-smi-exporter
	sudo cp ../README.md /usr/share/doc/amd-smi-exporter

copy_exporter_rules:
ifeq ($(shell test -d "/etc/prometheus" && echo -n yes),yes)
	sudo cp amd-smi-custom-rules.yml /etc/prometheus
else
	$(info @@@Follow README.md to install and configure prometheus@@@)
	sudo mkdir -p /etc/prometheus
	sudo cp amd-smi-custom-rules.yml /etc/prometheus
endif

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

# Generate packages using fpm tool
# Basic package information
PKG_NAME = amd_smi_exporter
PKG_DESCRIPTION = "Package AMD_SMI_EXPORTER"
PKG_VERSION = 1.0.0
PKG_RELEASE = 1
PKG_MAINTAINER = "Muralidhara M K \<muralimk@amd.com\>"
PKG_ARCH = all
PKG_ARCH_RPM = noarch

# These vars probably need no change
PKG_DEB = ${PKG_NAME}_${PKG_VERSION}-${PKG_RELEASE}_${PKG_ARCH}.deb
PKG_RPM = ${PKG_NAME}-${PKG_VERSION}-${PKG_RELEASE}.${PKG_ARCH_RPM}.rpm
FPM_OPTS = -s dir -n $(PKG_NAME) -v $(PKG_VERSION) --iteration $(PKG_RELEASE) --maintainer ${PKG_MAINTAINER} --description $(PKG_DESCRIPTION) -a $(PKG_ARCH)

BINS = /usr/local/bin/amd_smi_exporter /etc/systemd/system/amd-smi-exporter.service \
	/etc/prometheus/amd-smi-custom-rules.yml /usr/share/doc/amd-smi-exporter

FPM_BIN := $(shell command -v fpm 2> /dev/null)
check:
	$(if $(FPM_BIN),, \
		$(error Please install "ruby-dev ruby-dotenv" and then "gem install fpm"))


# Generate a deb and rpm packages using fpm
package:build \
	install \
	check \
	
	fpm -t deb -p $(PKG_DEB) $(FPM_OPTS) $(BINS)
	
	fpm -t rpm -p $(PKG_RPM) $(FPM_OPTS) $(BINS)
