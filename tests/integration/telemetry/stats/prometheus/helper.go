// Copyright 2019 Istio Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package promtheus

import (
	"fmt"
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/components/prometheus"
	"istio.io/istio/pkg/test/framework/components/sleep"
	util "istio.io/istio/tests/integration/mixer"
)

type Config struct {
	SleepNs    namespace.Instance
	Prometheus prometheus.Instance
	Metric     string
	Labels     string
	URL        string
}

// RunExternalRequestTest performs a query given the passed in config
// and validates it against data stored in prometheus
func RunExternalRequestTest(cfg *Config, expectErr bool, t *testing.T) {
	framework.
		NewTest(t).
		RequiresEnvironment(environment.Kube).
		Run(func(ctx framework.TestContext) {
			sleepInst := sleep.DeployOrFail(t, ctx, sleep.Config{Namespace: cfg.SleepNs, Cfg: sleep.Sleep})

			_, err := sleepInst.Curl(cfg.URL)
			errOut := err != nil
			if errOut != expectErr {
				t.Fatalf("did not receive expected condition (%s) from exec curl %s from sleep pod: %v", expectErr, cfg.URL, errOut)
			}
			query := fmt.Sprintf(`sum(%s{%s})`, cfg.Metric, cfg.Labels)
			util.ValidateMetric(t, cfg.Prometheus, query, cfg.Metric, 1)
		})
}
