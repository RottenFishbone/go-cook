#!/bin/bash

find . -name "*.go" -type f -exec gofmt -w $(dirname {}) \;
