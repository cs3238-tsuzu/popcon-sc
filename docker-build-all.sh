#! /bin/bash

set -eu

./ppweb/build/build.sh &&\
./ppjc/build/build.sh &&\
./ppms/build/build.sh &&\
./ppjudge/build/build.sh &&\
./modules/build.sh
