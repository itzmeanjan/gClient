SHELL := /bin/bash

build_pub:
	pushd publisher; go build -o ../gclient-pub; popd

build_sub:
	pushd subscriber; go build -o ../gclient-sub; popd

docker_pub:
	docker build -t pub -f ./publisher/Dockerfile .

docker_sub:
	docker build -t sub -f ./subscriber/Dockerfile .

run_pub:
	docker run --name pub -d pub

run_sub:
	docker run --name sub -d sub
