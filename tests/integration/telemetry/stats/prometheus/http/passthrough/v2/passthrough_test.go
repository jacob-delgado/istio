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

package v2

import (
	"testing"

	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/environment"
	"istio.io/istio/pkg/test/framework/components/galley"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/mixer"
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

func TestPassthroughCluster_TcpMetric(t *testing.T) {
	// Turn it on once issue is fixed.
	t.Skip("https://github.com/istio/istio/issues/21385")
	telemetry.RunExternalRequestTest(&telemetry.Config{
		SleepNs:    sleepNs,
		Prometheus: prom,
		Metric:     "istio_tcp_connections_closed_total",
		Labels:     `destination_service="PassthroughCluster",destination_service_name="PassthroughCluster"`,
		URL:        "https://prow.istio.io",
	}, true, t)
}

func TestMain(m *testing.M) {
	framework.
		NewSuite("mixer_telemetry_passthrough_metrics", m).
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
      mode: ALLOW_ANY
  telemetry:
    v1:
      enabled: false
    v2:
      enabled: true
components:
  policy:
    enabled: false
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
	g, err := galley.New(ctx, galley.Config{})
	if err != nil {
		return err
	}
	if _, err = mixer.New(ctx, mixer.Config{Galley: g}); err != nil {
		return err
	}
	prom, err = prometheus.New(ctx)
	if err != nil {
		return err
	}

	return nil
}
