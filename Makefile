SHELL := /bin/bash

build_pub:
	pushd publisher; go build -o ../gclient-pub; popd
