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
	"context"
	"encoding/json"
	"errors"
	av1 "github.com/Appdynamics/opentelemetry-ingest/gen/go/pb/appdynamics/v1"
	tv1 "github.com/Appdynamics/opentelemetry-ingest/gen/go/pb/trace/v1"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	typeStr = "appdynamics"
)

// Factory is the factory for the AppDynamics exporter.
type Factory struct {
}

// Type gets the type of the exporter configuration created by this factory.
func (f *Factory) Type() configmodels.Type {
	return typeStr
}

func (f *Factory) CreateDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{TypeVal: configmodels.Type(typeStr), NameVal: typeStr},
		AccessKey:        "",
		Account:          "",
		Metrics:          MetricsConfig{HTTP: &confighttp.HTTPClientSettings{}},
		Traces:           TracesConfig{Grpc: &configgrpc.GRPCClientSettings{}},
	}
}

// CreateTraceExporter creates a AppDynamics trace exporter for this configuration.
func (f *Factory) CreateTraceExporter(_ context.Context, p component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.TraceExporter, error) {
	config := cfg.(*Config)
	return exporterhelper.NewTraceExporter(config, createDataPusher(config, p))
}

func createDataPusher(config *Config, p component.ExporterCreateParams, additionalOps... grpc.DialOption) func(context.Context, pdata.Traces) (droppedSpans int, err error) {
	return func (ctx context.Context, td pdata.Traces) (droppedSpans int, err error) {
		opts, _ := config.Traces.Grpc.ToDialOptions()
		opts = append(opts, additionalOps...)

		var conn *grpc.ClientConn
		conn, err = grpc.Dial(config.Traces.Grpc.Endpoint, opts...)

		if err != nil {
			p.Logger.Error("fail to dial " + config.Traces.Grpc.Endpoint + err.Error())
			droppedSpans = td.ResourceSpans().Len()
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				p.Logger.Error(err.Error())
			}
		}()

		client := av1.NewSpanHandlerClient(conn)
		payload := Transform(td, *p.Logger)

		if _, err = client.HandleSpans(ctx, payload); err != nil {
			p.Logger.Info(err.Error())
			droppedSpans = td.ResourceSpans().Len()
		} else {
			droppedSpans = len(payload.ResourceSpans) - td.ResourceSpans().Len()
		}
		return
	}
}

func Transform(td pdata.Traces, logger zap.Logger) *av1.SpansRequest {
	otlpSpans := pdata.TracesToOtlp(td)
	var ret []*tv1.ResourceSpans
	for _, rs := range otlpSpans {
		// marshal to json from otlp.ResourceSpan
		marshal, err := json.Marshal(&rs)
		if err != nil {
			logger.Warn(err.Error())
			continue
		}
		// unmarshal it to IngestService.ResourceSpan
		var irs tv1.ResourceSpans
		if err := json.Unmarshal(marshal, &irs); err != nil {
			logger.Warn(err.Error())
			continue
		}
		ret = append(ret, &irs)
	}
	return &av1.SpansRequest{ResourceSpans: ret}
}

func (f *Factory) CreateMetricsExporter(_ context.Context, params component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.MetricsExporter, error) {
	// TODO : will need to create a metrics exporter
	return nil, errors.New("NOT YET IMPLEMENTED")
}
