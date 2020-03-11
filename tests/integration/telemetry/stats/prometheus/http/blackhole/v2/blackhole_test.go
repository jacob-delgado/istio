//  Copyright 2020 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package mixer

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/components/prometheus"
	"istio.io/istio/pkg/test/framework/label"
	"istio.io/istio/pkg/test/framework/resource"
	telemetry "istio.io/istio/tests/integration/telemetry/stats/prometheus"
)

var (
	sleepNs namespace.Instance
	ist     istio.Instance
	prom    prometheus.Instance
)

func TestBlackHoleCluster_HttpMetric(t *testing.T) {
	// no matching virtual outbound listener will be found, so it'll use the default
	// tcp matching virtual outbound listener
	// Turn it back on once issue is fixed.
	t.Skip("https://github.com/istio/istio/issues/21385")
	telemetry.RunExternalRequestTest(&telemetry.Config{
		SleepNs:    sleepNs,
		Prometheus: prom,
		Metric:     "istio_requests_total",
		Labels:     `destination_service="prow.istio.io",destination_service_name="BlackHoleCluster",response_code="502"`,
		URL:        "http://prow.istio.io",
	}, false, t)
}

func TestMain(m *testing.M) {
	framework.
		NewSuite("telemetry_v2_blackhole_metrics", m).
		RequireEnvironment(environment.Kube).
		Label(label.CustomSetup).
		SetupOnEnv(environment.Kube, istio.Setup(&ist, func(cfg *istio.Config) {
			cfg.ControlPlaneValues = `
values:
  prometheus:
    enabled: true
  global:
    disablePolicyChecks: false
    outboundTrafficPolicy:
      mode: REGISTRY_ONLY
  telemetry:
    v1:
      enabled: false
    v2:
      enabled: true
components:
  policy:
    enabled: true
  telemetry:
    enabled: false`
		})).
		Setup(testsetup).
		Run()
}

func testsetup(ctx resource.Context) (err error) {
	sleepNs, err = namespace.New(ctx, namespace.Config{
		Prefix: "istio-sleep",
		Inject: true,
	})
	if err != nil {
		return
	}
	prom, err = prometheus.New(ctx)
	if err != nil {
		return err
	}

	return nil
}
