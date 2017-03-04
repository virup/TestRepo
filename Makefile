DRUN = docker run --rm -i
DBUILD = docker build

DSTOP := docker stop
DRM := docker rm
DRMI := docker rmi
DSTART := docker start

BASEDIR := ${CURDIR}
SRCDIR := ${BASEDIR}/src
BIN_DIR := ${BASEDIR}/bin
INT_DIR := /libera
INT_BIN_DIR := ${INTDIR}/bin

COMPILER_IMAGE := soulfit/compiler
COMPILER_RESOURCE_DIR := ${BASEDIR}/pkg/compiler/
COMPILER_RESOURCE_DOCKERFILE := ${COMPILER_RESOURCE_DIR}/Dockerfile
COMPILER_CONTAINER := soulfit-compiler

DATABASE_RESOURCE_DIR := ${BASEDIR}/pkg/db
DATABASE_RESOURCE_DOCKERFILE := ${DATABASE_RESOURCE_DIR}/Dockerfile
DATABASE_IMAGE := soulfit/db
DATABASE_CONTAINER := soulfit-db
DATA_DIR := /opt/soulfit/db
INT_DATA_DIR := /var/lib/mysql/

SERVER_RESOURCE_DIR := ${BASEDIR}/pkg/server
SERVER_RESOURCE_DOCKERFILE := ${SERVER_RESOURCE_DIR}/Dockerfile
SERVER_IMAGE := soulfit/server
SERVER_CONTAINER := soulfit-server

setup:

	${DBUILD} -f ${COMPILER_RESOURCE_DOCKERFILE} -t ${COMPILER_IMAGE} ${COMPILER_RESOURCE_DIR}


clean: stop

	- ${DRMI} ${SERVER_IMAGE}
	- ${DRMI} ${DATABASE_IMAGE}

	rm -rf ${DATA_DIR}/*

	${DRUN} -v ${BASEDIR}:${INT_DIR}:z -e OP=CLEAN ${COMPILER_IMAGE}


compile:

	${DBUILD} -f ${DATABASE_RESOURCE_DOCKERFILE} -t ${DATABASE_IMAGE} ${DATABASE_RESOURCE_DIR}
	${DRUN} --name ${COMPILER_CONTAINER} -v ${BASEDIR}:${INT_DIR} -e OP=BUILD ${COMPILER_IMAGE}
	${DBUILD} -f ${SERVER_RESOURCE_DOCKERFILE} -t ${SERVER_IMAGE} ${SERVER_RESOURCE_DIR}


run:

	mkdir -p ${DATA_DIR}
	${MAKE} startcontainers

stop:

	- ${DSTOP} ${DATABASE_CONTAINER}
	- ${DRM} ${DATABASE_CONTAINER}

	- ${DSTOP} ${SERVER_CONTAINER}
	- ${DRM} ${SERVER_CONTAINER}

start: startcontainers

startcontainers:

	docker run -d --name ${DATABASE_CONTAINER} -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=soulfitdb -v ${DATA_DIR}:${INT_DATA_DIR} ${DATABASE_IMAGE}
	docker run -p 8099:8099 -d --name ${SERVER_CONTAINER} --link soulfit-db:sf-db -v ${BASEDIR}:${INT_DIR} ${SERVER_IMAGE}

