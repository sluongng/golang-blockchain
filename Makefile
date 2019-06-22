# Go parameters
GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=${GOCMD} clean
GOTEST=${GOCMD} test
GOGET=${GOCMD} get
BUILD_PATH=./bin/
BINARY_NAME=bc_cli
BINARY_PATH=${BUILD_PATH}${BINARY_NAME}

docker:
	docker-compose up -d --force-recreate --build

database: docker

build:
	${GOBUILD} -o ${BINARY_PATH} -v

run: build
	${BINARY_PATH}

clean:
	${GOCLEAN}
	rm -f *.db
	rm -rf ./bin
	rm -rf ./.boltDb

all: build 
