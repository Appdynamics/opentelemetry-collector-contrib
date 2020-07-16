// Copyright 2020 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package appdynamicsexporter

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.uber.org/zap"
	"net/http"
)

const (
	typeStr = "appdynamics"
)

// Factory is the factory for the AppDynamics exporter.
type Factory struct {
}

// Type gets the type of the exporter configuration created by this factory.
func (f *Factory) Type() configmodels.Type {
	return configmodels.Type(typeStr)
}

// CreateDefaultConfig creates a default configuration for this exporter.
func (f *Factory) CreateDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: configmodels.Type(typeStr),
			NameVal: typeStr,
		},
		AccessToken: "",
		EndPoint:    "https://saas.appdynamics.com:443",
	}
}

// CreateTraceExporter creates a AppDynamics trace exporter for this configuration.
func (f *Factory) CreateTraceExporter(
	ctx context.Context,
	_ component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.TraceExporter, error) {
	config := cfg.(*Config)
	//_ := http.DefaultClient
	return exporterhelper.NewTraceExporter(config,
		func(ctx context.Context, td pdata.Traces) (droppedSpans int, err error) {
			otldpData := pdata.TracesToOtlp(td)
			var errs []error
			m := jsonpb.Marshaler{OrigName: true}
			for _, otlp := range otldpData {
				var b bytes.Buffer
				if err := m.Marshal(&b, otlp); err != nil {
					droppedSpans += 1
					errs = append(errs, err)
				} else {
					_, err := http.DefaultClient.Post(config.EndPoint, "application/json", &b)
					if err != nil {
						errs = append(errs)
						droppedSpans += 1
					}
				}
			}
			return droppedSpans, componenterror.CombineErrors(errs)
		})
}


// CreateMetricsExporter always returns nil.
func (f *Factory) CreateMetricsExporter(
	_ context.Context,
	_ component.ExporterCreateParams,
	_ configmodels.Exporter,
) (component.MetricsExporter, error) {
	return nil, configerror.ErrDataTypeIsNotSupported
}
