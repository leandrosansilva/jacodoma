#!/bin/bash

set -e

JACODOMA_SOURCE="$PWD"

docker run --rm -it -e DISPLAY=$DISPLAY \
  -v /tmp/.X11-unix:/tmp/.X11-unix \
  -v "$JACODOMA_SOURCE":/gopath/src/jacodoma leandrosansilva/jacodoma \
  "$@"
