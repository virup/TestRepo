FROM centos:7
RUN yum -y install go

COPY compile.sh /compile.sh
RUN chmod 777 /compile.sh

ENTRYPOINT ["/compile.sh"]
