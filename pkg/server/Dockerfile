FROM centos:7
RUN yum -y install go

COPY run_server.sh /run_server.sh
RUN chmod 777 /run_server.sh
EXPOSE 8080
ENTRYPOINT ["/run_server.sh"]
