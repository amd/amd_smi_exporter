# MIT-X11 Open Source License
# Copyright (c) 2024, Advanced Micro Devices, Inc.

#Ubuntu based container
FROM ubuntu:20.04
# For rhel8 based container
# FROM registry.access.redhat.com/ubi8/ubi
MAINTAINER muthusamy.ramalingam@amd.com
ENV LD_LIBRARY_PATH=/opt/goamdsmi/lib:/opt/e-sms/e_smi/lib:/opt/rocm/lib:/opt/rocm/lib64
COPY ./install/library ./opt
COPY ./install/binconfig/amd_smi_exporter ./usr/local/bin
ENTRYPOINT ["/usr/local/bin/amd_smi_exporter"]
