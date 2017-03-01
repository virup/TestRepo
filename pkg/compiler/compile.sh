#!/usr/bin/env bash

if [ ${OP} == "CLEAN" ]; then
    echo "Cleaning bin directory ..."
    rm -rf /libera/bin/server
    rm -rf /libera/bin/grpcsqltest
    echo "Done !"
    exit 0
fi

# Default option is to build
echo "Building ..."
cd /libera/bin
sh ../buildit.sh > /libera/build.log 2>/libera/build.err.log
echo "Done! "
