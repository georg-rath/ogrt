#!/usr/bin/env bash

export LD_PRELOAD=$(pwd)/ogrt/lib/libogrt.so
export OGRT_ACTIVE=1
export OGRT_SILENT=1
export PATH=$(pwd)/ogrt/bin/:$PATH
export LD_LIBRARY_PATH=$(pwd)/sw/liba/:$(pwd)/sw/libb/:$LD_LIBRARY_PATH
