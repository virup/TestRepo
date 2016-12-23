# libera

- Get the repo and the binariers
cd /Users/YOURUSERNAME
mkdir src
cd src
git clone github.com/LiberaLabs/libera

- Server runs in a rocksdb container
docker pull az3r/golang-rocksdb
docker run -p 8080:8080 --privileged -it -v /Users/YOURUSERNAME/src/libera/libera:/libera c036c6e84f35  /bin/bash
- Inside the container
cd /libera/bin 
../buildit.sh
./server

- buildit.sh compiles the proto file and the server.
