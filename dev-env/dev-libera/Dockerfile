FROM dev-go

# Install compression libraries
RUN yum -y update && \
    yum install -y zlib-devel bzip2-devel snappy-devel &&\
    yum install -y make gcc-c++

# Get and build rocksdb 4.9
RUN cd / && \
    wget https://github.com/facebook/rocksdb/archive/v4.9.tar.gz && \
    tar -xzf v4.9.tar.gz && \
    rm -f v4.9.tar.gz && \
    cd rocksdb-4.9 && \
    make static_lib && \
    cd /

# Get golang wrapper of rocksdb and point it to the rocksdb include directory
# and static library
RUN CGO_CFLAGS="-I/rocksdb-4.9/include" \
    CGO_LDFLAGS="-L/rocksdb-4.9/ \
    -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy"   \
    go get github.com/tecbot/gorocksdb

RUN go get golang.org/x/net/context && \
    go get github.com/coreos/etcd/client
