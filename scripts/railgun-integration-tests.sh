#!/bin/bash

cd railgun/

if [ -z "$1" ]
then
    make test.integration
elif [ "$1" == "LEGACY" ]
then
    make test.integration.legacy
else
    echo "$1 isn't a valid option (valid options: LEGACY)"
    exit 1
fi
