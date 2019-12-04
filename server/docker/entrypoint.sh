#!/bin/bash
cd /usr/local/ogrt
if [[ -e /usr/local/ogrt/ogrt.conf ]]; then
    ./ogrt-server
else
    echo "OGRT configuration file is missing"
fi
