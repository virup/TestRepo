#!/bin/bash

# Build base image dev-go, dev-libera
docker build -t dev-go dev-go/.
docker build -t dev-libera dev-libera/.
