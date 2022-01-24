#! /bin/bash

lsns | grep Singularity | awk '{print $4}'

