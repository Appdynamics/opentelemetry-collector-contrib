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
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.uber.org/zap"
)

const (
	typeStr = "appdynamics"
)

// Factory is the factory for the LightStep exporter.
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

// CreateTraceExporter creates a LightStep trace exporter for this configuration.
func (f *Factory) CreateTraceExporter(_ *zap.Logger, cfg configmodels.Exporter) (component.TraceExporterOld, error) {
	config := cfg.(*Config)
	return &AppDynamicsExporter{
		endpoint:  config.EndPoint,
		accessKey: config.AccessToken,
	}, nil
}

// CreateMetricsExporter returns nil.
func (f *Factory) CreateMetricsExporter(logger *zap.Logger, cfg configmodels.Exporter) (component.MetricsExporterOld, error) {
	config := cfg.(*Config)
	return &AppDynamicsExporter{
		endpoint:  config.EndPoint,
		accessKey: config.AccessToken,
	}, nil
}
