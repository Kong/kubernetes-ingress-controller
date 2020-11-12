#!/bin/bash

#    Copyright 2015 Grafana Labs
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
#   Kong Inc. has modified the original Grafana source of this script.
#   The modifications add additional instructions to its output text.

set -eo pipefail

# Verify that Go is properly installed and available
command -v go >/dev/null 2>&1 || { echo 'please install Go or use an image that has it'; exit 1; }

backup_go_mod_files()
{
    mod=$(mktemp)
    cp go.mod "$mod"

    sum=$(mktemp)
    cp go.sum "$sum"
}

restore_go_mod_files()
{
    cp "$mod" go.mod
    rm "$mod"

    cp "$sum" go.sum
    rm "$sum"
}

# Backup current go.mod and go.sum files
backup_go_mod_files

# Defer the go.mod and go.sum files backup recovery
trap restore_go_mod_files EXIT

# Tidy go.mod and go.sum files
go mod tidy

diff "$mod" go.mod || { echo "your go.mod is inconsistent. please run \"go mod tidy\""; exit 1; }
diff "$sum" go.sum || { echo "your go.sum is inconsistent. please run \"go mod tidy\""; exit 1; }
