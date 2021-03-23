#!/bin/bash

set -euox pipefail

cd railgun/
make test.integration
