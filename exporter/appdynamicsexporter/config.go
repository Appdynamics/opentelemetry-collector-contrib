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
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
)


type TracesConfig struct {
	Grpc *configgrpc.GRPCClientSettings `mapstructure:"grpc"`
}

type MetricsConfig struct {
	HTTP *confighttp.HTTPClientSettings `mapstructure:"http"`
}

// Config defines configuration options for the LightStep exporter.
type Config struct {
	configmodels.ExporterSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	// AccessKey is your project access token in AppDynamics.
	AccessKey string `mapstructure:"access_key"`
	// Account
	Account string `mapstructure:"account"`
	// Metrics Ingestion config
	Metrics MetricsConfig `mapstructure:"metrics"`
	// Traces Ingestion config
	Traces TracesConfig `mapstructure:"traces"`

}
