#! /bin/bash

cat /proc/cpuinfo | grep "core id" | awk '{print $4}'
