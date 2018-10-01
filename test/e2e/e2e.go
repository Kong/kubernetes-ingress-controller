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

package e2e

import (
	"testing"

	"github.com/golang/glog"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega"
	"k8s.io/apiserver/pkg/util/logs"

	// required
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/kong/kubernetes-ingress-controller/test/e2e/framework"

	// tests to run
	_ "github.com/kong/kubernetes-ingress-controller/test/e2e/annotations"
	_ "github.com/kong/kubernetes-ingress-controller/test/e2e/defaultbackend"
	_ "github.com/kong/kubernetes-ingress-controller/test/e2e/settings"
	_ "github.com/kong/kubernetes-ingress-controller/test/e2e/ssl"
)

// RunE2ETests checks configuration parameters (specified through flags) and then runs
// E2E tests using the Ginkgo runner.
func RunE2ETests(t *testing.T) {
	logs.InitLogs()
	defer logs.FlushLogs()

	gomega.RegisterFailHandler(ginkgo.Fail)
	// Disable skipped tests unless they are explicitly requested.
	if config.GinkgoConfig.FocusString == "" && config.GinkgoConfig.SkipString == "" {
		config.GinkgoConfig.SkipString = `\[Flaky\]|\[Feature:.+\]`
	}

	glog.Infof("Starting e2e run %q on Ginkgo node %d", framework.RunID, config.GinkgoConfig.ParallelNode)
	ginkgo.RunSpecs(t, "nginx-ingress-controller e2e suite")
}
