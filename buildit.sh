go_back() {
 cd - > /dev/null
}

build_protocsql() {
 echo "Building proto sql files..."
 cd $SCRIPTDIR/bin
 protoc -I$SCRIPTDIR/src/server/rpcdefsql/ --go_out=plugins=grpc:$SCRIPTDIR/src/server/rpcdefsql/ $SCRIPTDIR/src/server/rpcdefsql/serverrpc.proto
 go_back
}

build_grpcsqltest() {
 echo "Cleaning up old grpcsqltest..."
 rm -f $OUTDIR/grpcsqltest
 echo "Building grpcsqltest..."
 cd src/test/grpcsqltest
 go build -gcflags "-N -l" -o $OUTDIR/grpcsqltest
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


build_serversql() {
 echo "Cleaning up old sql server..."
 rm -f $OUTDIR/serversql
 cd src/server/serversql

 echo "Building server sql..."
 go build -gcflags "-N -l" -o $OUTDIR/serversql
 go_back
}

# Build everything. At a later date we could give command line options to
# build only specific things
main() {
 build_protocsql
 build_serversql
 build_grpcsqltest
 #build_clienttest
}

# Get the parent directory of this script and store all the binaries in the
# bin directory
SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
OUTDIR=$SCRIPTDIR/bin/

cd $SCRIPTDIR
export GOPATH=$SCRIPTDIR
export PATH=$PATH:$SCRIPTDIR/bin
main
