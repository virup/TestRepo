# libera

- Get the repo and the binaries

`cd ~`

`mkdir soulfit`

`cd soulfit`

`https://gitlab.com/LiberaLabs/libera.git`


-- Server runs in a centos golang container, lets build it.

# Build the developer image
`docker build -f ~/soulfit/libera/pkg/developer/Dockerfile -t soulfit/developer ~/soulfit/libera/pkg/developer`

`make setup`

`make clean`

`make compile`

# Run the developer image
`sudo mkdir -p /opt/soulfit/db`

`docker run -d --name soulfit-db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=soulfitdb -v /opt/soulfit/db:/var/lib/mysql/ soulfit/db`

`docker run --name soulfit-developer --link soulfit-db:sf-db -it -v ~/soulfit/libera:/libera soulfit/developer  /bin/bash`

# Inside the developer container, run the server
`cd /libera/bin`

- Compile the proto file and the server.

`../buildit.sh`

- Run the server

`./server`



# For compiling and running the production app

make setup

make clean

make compile

sudo make run


# To restart

make stop

make start
