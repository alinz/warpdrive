#!/bin/bash

# Make sure we're in our git repository
cd $(dirname "$0")

if [ "$1" == "--long" ]; then
	git describe --tags --long --dirty
else
	git describe --tags --abbrev=0
fi
