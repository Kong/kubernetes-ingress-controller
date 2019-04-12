/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"testing"
)

// resetForTesting clears all flag state and sets the usage function as directed.
// After calling resetForTesting, parse errors in flag handling will not
// exit the program.
// Extracted from https://github.com/golang/go/blob/master/src/flag/export_test.go
func resetForTesting(usage func()) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.Usage = usage
}

func TestDefaults(t *testing.T) {
	resetForTesting(func() { t.Fatal("bad parse") })

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "--publish-service", "namespace/test"}

	showVersion, conf, err := parseFlags()
	if err != nil {
		t.Fatalf("unexpected error parsing default flags: %v", err)
	}

	if showVersion {
		t.Fatal("expected false but true was returned for flag show-version")
	}

	if conf == nil {
		t.Fatal("expected a configuration but nil returned")
	}
}
