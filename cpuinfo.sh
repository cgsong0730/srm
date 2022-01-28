#! /bin/bash

cat /proc/cpuinfo | grep "processor" | awk '{print $3}' | wc -l
