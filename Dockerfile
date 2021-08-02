##
## To run TSAK-2 in container
##
FROM ubuntu:20.04
LABEL description="TSAK-2 services"
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-update
ARG MASTER_SCRIPT
ENV SNMP_LISTEN=127.0.0.1
ENV SNMP_AGENT_PORT=6161
ENV TRAP_RECEIVER_PORT=6163
ENV GEN=15
ENV SCRIPT=${MASTER_SCRIPT}
VOLUME /home/ubuntu/mibs
COPY /home/ubuntu/tsak2 /

CMD [ "/tsak2" "--mibs=/mibs/nri-snmp.db" "--snmplisten=${SNMP_LISTEN}" "--snmpport=${SNMP_AGENT_PORT}" "--trapport=${TRAP_RECEIVER_PORT}" "--gen=${GEN}" "run" "${SCRIPT}" "--erloop" ]
