#Ubuntu based container
FROM ubuntu:20.04
MAINTAINER krishnan.c@amd.com
ENV LD_LIBRARY_PATH=/opt/goamdsmi/lib:/opt/e-sms/e_smi/lib:/opt/rocm/rocm_smi/lib
COPY ./install/library ./opt
COPY ./install/binconfig/amd_smi_exporter ./usr/local/bin
ENTRYPOINT ["/usr/local/bin/amd_smi_exporter"]
