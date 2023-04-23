#!/bin/bash

testedDirs=( parser/ )
dir=$(dirname ${BASH_SOURCE})

for d in ${testedDirs[@]} 
do
    echo "Testing: $d"
    go test -v "$dir/$d"
done
