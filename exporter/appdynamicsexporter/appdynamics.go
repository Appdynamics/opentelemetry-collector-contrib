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
	pb "github.com/pavankrish123/ot-sample-lib/gen/ingest/v1"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumerdata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AppDynamicsExporter struct {
	endpoint  string
	accessKey string
}

func (c *AppDynamicsExporter) ConsumeMetricsData(ctx context.Context, md consumerdata.MetricsData) error {
	return nil
}

func (c *AppDynamicsExporter) ConsumeTraceData(ctx context.Context, td consumerdata.TraceData) error {
	conn, err := grpc.Dial(c.endpoint)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx = metadata.AppendToOutgoingContext(ctx, "appd-access-key", c.accessKey)
	_, err = pb.NewIngestClient(conn).IngestTraces(ctx, &pb.IngestTraceServiceRequest{
		Spans: td.Spans,
	})
	return err
}

func (c *AppDynamicsExporter) Start(ctx context.Context, host component.Host) error {
	panic("implement me")
}

func (c *AppDynamicsExporter) Shutdown(ctx context.Context) error {
	return nil
}
