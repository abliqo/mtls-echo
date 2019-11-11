#!/bin/bash

docker run \
	-e CONFIG_PATH=/opt/mtls-echo/config/dev.json \
	-v $(pwd)/test/certs:/opt/mtls-echo/test/certs \
	-d \
	--rm \
	--net=host \
	--name="mtls-echo" mtls-echo