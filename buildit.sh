go_back() {
 cd - > /dev/null
}

build_protoc() {
 echo "Building proto files..."
 cd $SCRIPTDIR/bin
 protoc -I$SCRIPTDIR/src/server/rpcdef/ --go_out=plugins=grpc:$SCRIPTDIR/src/server/rpcdef/ $SCRIPTDIR/src/server/rpcdef/serverrpc.proto
 go_back
}

build_grpctest() {
 echo "Cleaning up old grpctest..."
 rm -f $OUTDIR/grpctest
 echo "Building grpctest..."
 cd src/test/grpctest
 go build -gcflags "-N -l" -o $OUTDIR/grpctest
 go_back
}

build_clienttest() {
 echo "Cleaning up old clientest..."
 rm -f $OUTDIR/clientest
 echo "Building clientest..."
 cd src/test/clientest
 go build -gcflags "-N -l" -o $OUTDIR/clientest
 go_back
}



build_server() {
 echo "Cleaning up old server..."
 rm -f $OUTDIR/server
 cd src/server/server

 echo "Building server..."
 go build -gcflags "-N -l" -o $OUTDIR/server
 go_back
}

# Build everything. At a later date we could give command line options to
# build only specific things
main() {
 build_protoc
 build_server
 build_grpctest
 build_clienttest
}

# Get the parent directory of this script and store all the binaries in the
# bin directory
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTDIR=$SCRIPTDIR/bin/

cd $SCRIPTDIR
export GOPATH=$SCRIPTDIR
export PATH=$PATH:$SCRIPTDIR/bin
main
