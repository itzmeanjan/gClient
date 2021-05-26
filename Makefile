SHELL := /bin/bash

build_pub:
	pushd publisher; go build -o ../gclient-pub; popd

build_sub:
	pushd subscriber; go build -o ../gclient-sub; popd
