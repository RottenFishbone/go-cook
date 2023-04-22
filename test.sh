#!/bin/bash

testedDirs=( pkg/recipe/ )
dir=$(dirname ${BASH_SOURCE})

for d in ${testedDirs[@]} 
do
    echo "Testing: $d"
    go test -v "$dir/$d"
done
