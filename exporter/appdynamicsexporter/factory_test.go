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
	"encoding/binary"
	"errors"
	"fmt"
	av1 "github.com/Appdynamics/opentelemetry-ingest/gen/go/pb/appdynamics/v1"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer/pdata"
	semconventions "go.opentelemetry.io/collector/translator/conventions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
	"net"
	"testing"
	"time"
)

// NOTE : this trace data type and test cases are taken from :
//  internal/data/testdata/trace_test.go
//  may want to request all of these cases made exportable

type traceTestCase struct {
	name string
	td   pdata.Traces
}

func generateAllTraceTestCases() []traceTestCase {
	return []traceTestCase{
		{
			name: "empty-traces",
			td:   constructEmptyTrace(),
		},
		{
			name: "one-empty-trace",
			td :  constructOneEmptyTrace(),
		},
		{
			name: "one-empty-one-nil-trace",
			td :  constructOneEmptyOneNilTrace(),
		},
		{
			name: "one-normal-trace",
			td:   constructOneNormalTrace(),
		},
	}
}

func constructEmptyTrace() pdata.Traces {
	traces := pdata.NewTraces()
	return traces
}

func constructOneEmptyTrace() pdata.Traces {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	return traces
}

func constructOneEmptyOneNilTrace() pdata.Traces {
	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)

	resSpan := pdata.NewResourceSpans()
	traces.ResourceSpans().Append(&resSpan)

	return traces
}

func constructOneNormalTrace() pdata.Traces {
	resource := constructResource()

	traces := pdata.NewTraces()
	traces.ResourceSpans().Resize(1)
	rspans := traces.ResourceSpans().At(0)
	resource.CopyTo(rspans.Resource())
	rspans.InstrumentationLibrarySpans().Resize(1)
	ispans := rspans.InstrumentationLibrarySpans().At(0)
	ispans.Spans().Resize(2)
	constructHTTPClientSpan().CopyTo(ispans.Spans().At(0))
	constructHTTPServerSpan().CopyTo(ispans.Spans().At(1))
	return traces
}

func constructResource() pdata.Resource {
	resource := pdata.NewResource()
	resource.InitEmpty()
	attrs := pdata.NewAttributeMap()
	attrs.InsertString(semconventions.AttributeServiceName, "signup_aggregator")
	attrs.InsertString(semconventions.AttributeContainerName, "signup_aggregator")
	attrs.InsertString(semconventions.AttributeContainerImage, "otel/signupaggregator")
	attrs.CopyTo(resource.Attributes())
	return resource
}

func constructHTTPClientSpan() pdata.Span {
	attributes := make(map[string]interface{})
	attributes[semconventions.AttributeComponent] = semconventions.ComponentTypeHTTP
	attributes[semconventions.AttributeHTTPMethod] = "GET"
	attributes[semconventions.AttributeHTTPURL] = "https://api.example.com/users/junit"
	attributes[semconventions.AttributeHTTPStatusCode] = 200
	endTime := time.Now().Round(time.Second)
	startTime := endTime.Add(-90 * time.Second)
	spanAttributes := constructSpanAttributes(attributes)

	span := pdata.NewSpan()
	span.InitEmpty()
	span.SetTraceID(newTraceID())
	span.SetSpanID(newSegmentID())
	span.SetParentSpanID(newSegmentID())
	span.SetName("/users/junit")
	span.SetKind(pdata.SpanKindCLIENT)
	span.SetStartTime(pdata.TimestampUnixNano(startTime.UnixNano()))
	span.SetEndTime(pdata.TimestampUnixNano(endTime.UnixNano()))

	status := pdata.NewSpanStatus()
	status.InitEmpty()
	status.SetCode(0)
	status.SetMessage("OK")
	status.CopyTo(span.Status())

	spanAttributes.CopyTo(span.Attributes())
	return span
}

func constructHTTPServerSpan() pdata.Span {
	attributes := make(map[string]interface{})
	attributes[semconventions.AttributeComponent] = semconventions.ComponentTypeHTTP
	attributes[semconventions.AttributeHTTPMethod] = "GET"
	attributes[semconventions.AttributeHTTPURL] = "https://api.example.com/users/junit"
	attributes[semconventions.AttributeHTTPClientIP] = "192.168.15.32"
	attributes[semconventions.AttributeHTTPStatusCode] = 200
	endTime := time.Now().Round(time.Second)
	startTime := endTime.Add(-90 * time.Second)
	spanAttributes := constructSpanAttributes(attributes)

	span := pdata.NewSpan()
	span.InitEmpty()
	span.SetTraceID(newTraceID())
	span.SetSpanID(newSegmentID())
	span.SetParentSpanID(newSegmentID())
	span.SetName("/users/junit")
	span.SetKind(pdata.SpanKindSERVER)
	span.SetStartTime(pdata.TimestampUnixNano(startTime.UnixNano()))
	span.SetEndTime(pdata.TimestampUnixNano(endTime.UnixNano()))

	status := pdata.NewSpanStatus()
	status.InitEmpty()
	status.SetCode(0)
	status.SetMessage("OK")
	status.CopyTo(span.Status())

	spanAttributes.CopyTo(span.Attributes())
	return span
}

func constructSpanAttributes(attributes map[string]interface{}) pdata.AttributeMap {
	attrs := pdata.NewAttributeMap()
	for key, value := range attributes {
		if cast, ok := value.(int); ok {
			attrs.InsertInt(key, int64(cast))
		} else if cast, ok := value.(int64); ok {
			attrs.InsertInt(key, cast)
		} else {
			attrs.InsertString(key, fmt.Sprintf("%v", value))
		}
	}
	return attrs
}

func newTraceID() []byte {
	var r [16]byte
	epoch := time.Now().Unix()
	binary.BigEndian.PutUint32(r[0:4], uint32(epoch))
	_, err := rand.Read(r[4:])
	if err != nil {
		panic(err)
	}
	return r[:]
}

func newSegmentID() []byte {
	var r [8]byte
	_, err := rand.Read(r[:])
	if err != nil {
		panic(err)
	}
	return r[:]
}


func TestTransform(t *testing.T) {
	testsData := generateAllTraceTestCases()
	var buff bytes.Buffer

	// nothing should start in the buffer
	assert.True(t, buff.Len() == 0)

	// took tips from : https://developerpod.com/2019/05/27/adding-uber-go-zap-logger-to-golang-project/#3.0
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	writeSyncer := zapcore.AddSync(&buff)
	core := zapcore.NewCore(encoder, writeSyncer, zap.WarnLevel)
	logger := zap.New(core)

	// TODO : more test cases as in AWS stuff
	// TODO : need to write a function to verify the output of transform
	for _, elem := range testsData {
		Transform(elem.td, *logger)
	}

	// flush, then check there were no errors
	logger.Sync()
	fmt.Println(buff.String())
	assert.True(t, buff.Len() == 0)
}

// Helper function to create a GRPC config
func createGRPCClientConfig(endpoint string) configgrpc.GRPCClientSettings {
	return configgrpc.GRPCClientSettings{
		Headers: map[string]string{
			"serverType": "test",
		},
		Endpoint:    endpoint,
		Compression: "gzip",
		TLSSetting: configtls.TLSClientSetting{
			Insecure: true,
		},
		Keepalive: &configgrpc.KeepaliveClientConfig{
			Time:                time.Second,
			Timeout:             time.Second,
			PermitWithoutStream: true,
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		WaitForReady:    true,
	}
}

// Helper function to create a config for the exporter
func createExporterConfig(config *configgrpc.GRPCClientSettings) Config {
	return Config{
		ExporterSettings: configmodels.ExporterSettings{},
		AccessKey:        "AccKey",
		Account:          "Acct",
		Metrics:          MetricsConfig{},
		Traces:           TracesConfig{
			Grpc: config,
		},
	}
}

// A server base encapsulates the functionality that will be needed
//  for all the servers
	type myListener interface {
		net.Listener
		ToDial(context.Context, string) (net.Conn, error)
	}
	type server interface {
		av1.SpanHandlerServer
		GetConn() myListener
	}
	type serverBase struct {
		*av1.UnimplementedSpanHandlerServer
		lis myListener
	}
	func (s serverBase) GetConn() myListener {
		return s.lis
	}
	func initServer(server server) *grpc.Server{
		s := grpc.NewServer()
		av1.RegisterSpanHandlerServer(s, server)
		go func() {
			if err := s.Serve(server.GetConn()); err != nil {
				fmt.Println("server no longer serving : " +
				err.Error())
			}
		}()
		return s
	}

// These listeners will be used for testing
	// acts as expected
	type safeListener struct {
		*bufconn.Listener
	}
	func (myLis safeListener) ToDial(_ context.Context, _ string) (net.Conn, error) {
		return myLis.Dial()
	}
	// will not be able to dial the server
	type badDialListener struct {
		*bufconn.Listener
	}
	func (myLis badDialListener) ToDial(_ context.Context, _ string) (net.Conn, error) {
		return nil, errors.New("A proper Dial error")
	}
	// will not close the connection correctly
	type badCloseListener struct {
		*bufconn.Listener
	}
	func (myLis badCloseListener) ToDial(_ context.Context, _ string) (net.Conn, error) {
		return myLis.Dial()
	}
	func (myLis badCloseListener) Close() error {
		return errors.New("a proper Close error")
	}

// The normal server used for the typical case
	type normalServer struct {
		*serverBase
	}
	func (_ normalServer) HandleSpans(_ context.Context, _ *av1.SpansRequest) (*av1.SpansResponse, error) {
		resp := av1.SpansResponse{}
		return &resp, nil
	}

func TestDataPusherNormalServer(t *testing.T) {
	// First create the Exporter params (essentially a logger)
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
			func(e zapcore.Entry) error {
				t.Fatal(
					"there should be no logging in the normal case")
				return nil
			})))
	p := component.ExporterCreateParams{Logger: logger}

	// Create a configuration for the exporter
	GRPClientConfig := createGRPCClientConfig("")
	config := createExporterConfig(&GRPClientConfig)

	// Test with a normal server
	server := normalServer{
		&serverBase{
			&av1.UnimplementedSpanHandlerServer{},
			safeListener{
				bufconn.Listen(1024 * 1024),
			},
		},
	}
	initServer(server)

	// Create the test data and the test function
	testCases := generateAllTraceTestCases()
	testDataPusher := createDataPusher(
		&config, p, grpc.WithContextDialer(server.lis.ToDial))

	droppedSpans, err := testDataPusher(context.Background(), testCases[3].td)
	assert.True(t, droppedSpans == 0)
	assert.True(t, err == nil)

	// Flush, the hook should assert there were no errors
	logger.Sync()
}

// test the case the server throws an error when processing the span
	type badHandleServer struct {
		*serverBase
	}
	func (s badHandleServer) HandleSpans(_ context.Context, _ *av1.SpansRequest) (*av1.SpansResponse, error) {
		return nil, errors.New("cannot process span")
	}

func TestDataPusherBadHandleServer(t *testing.T) {
	// First create the Exporter params (essentially a logger)
	logAppeared := make(chan bool, 1)
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
		func(e zapcore.Entry) error {
			if e.Level != zap.InfoLevel {
				t.Fatal(
					"there should only be one INFO level log")
			} else {
				if e.Message != "rpc error: code = Unknown desc = cannot process span" {
					t.Fatal("incorrect log entry")
				} else {
					logAppeared <- true
					return nil
				}
			}
			return nil
		})))
	p := component.ExporterCreateParams{Logger: logger}

	// Create a configuration for the exporter
	GRPClientConfig := createGRPCClientConfig("")
	config := createExporterConfig(&GRPClientConfig)

	// Test with a normal server
	server := badHandleServer{
		&serverBase{
			&av1.UnimplementedSpanHandlerServer{},
			safeListener{
				bufconn.Listen(1024 * 1024),
			},
		},
	}
	initServer(server)

	// Create the test data and the test function
	testCases := generateAllTraceTestCases()
	testDataPusher := createDataPusher(
		&config, p, grpc.WithContextDialer(server.lis.ToDial))

	droppedSpans, err := testDataPusher(context.Background(), testCases[3].td)
	assert.True(t, droppedSpans == 1)
	assert.True(t, err != nil)

	// flush, then assert the correct log was written
	logger.Sync()
	select {
		case b :=<-logAppeared: {
			assert.True(t, b)
		}
		case <-time.After(time.Second):
			t.Fatal("No log appeared")
	}
}

// test the case where the cannot dial the server
	type badDialServer struct {
		*serverBase
	}
	func (s badDialServer) HandleSpans(_ context.Context, _ *av1.SpansRequest) (*av1.SpansResponse, error) {
		// this will not be used
		return nil, nil
	}

func TestDataPusherBadDialServer(t *testing.T) {
	// First create the Exporter params (essentially a logger)
	logAppeared := make(chan bool, 1)
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
		func(e zapcore.Entry) error {
			if e.Level != zap.ErrorLevel {
				t.Fatal(
					"there should only be one ERROR level log")
			} else {
				if e.Message != "fail to dial badEndpointcontext deadline exceeded" {
					t.Fatal("incorrect log entry")
				} else {
					logAppeared <- true
					return nil
				}
			}
			return nil
		})))
	p := component.ExporterCreateParams{Logger: logger}

	// Create a configuration for the exporter
	GRPClientConfig := createGRPCClientConfig("badEndpoint")
	config := createExporterConfig(&GRPClientConfig)

	// Test with a normal server
	server := badDialServer{
		&serverBase{
			&av1.UnimplementedSpanHandlerServer{},
			badDialListener{
				bufconn.Listen(1024 * 1024),
			},
		},
	}
	initServer(server)

	// Create the test data and the test function
	testCases := generateAllTraceTestCases()
	// TODO : 'WithTimeout' is depricated, will need to use
	//  DialWithContext within createDataPusher and set
	//  a deadline within the context
	testDataPusher := createDataPusher(&config, p,
		grpc.WithBlock(), grpc.WithTimeout(time.Second))

	droppedSpans, err := testDataPusher(context.Background(), testCases[3].td)
	assert.True(t, droppedSpans == 1)
	assert.True(t, err != nil)

	// flush, then check for the correct error
	logger.Sync()
	select {
	case b :=<-logAppeared: {
		assert.True(t, b)
	}
	case <-time.After(time.Second):
		t.Fatal("No log appeared")
	}
}

// test the case that we cannot close the connection properly
	type badCloseServer struct {
		*serverBase
	}
	func (s badCloseServer) HandleSpans(_ context.Context, _ *av1.SpansRequest) (*av1.SpansResponse, error) {
		resp := av1.SpansResponse{}
		return &resp, nil
	}

	// TODO : no log is appearing on CLOSE of the connection
func TestDataPusherCloseServer(t *testing.T) {
	// First create the Exporter params (essentially a logger)
	logAppeared := make(chan bool, 1)
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(
		func(e zapcore.Entry) error {
			if e.Level != zap.ErrorLevel {
				t.Fatal(
					"there should only be one ERROR level log")
			} else {
				// TODO : when can get the right error out of
				//  the test, replace "test"
				if e.Message != "test" {
					t.Fatal("incorrect log entry")
				} else {
					logAppeared <- true
					return nil
				}
			}
			return nil
		})))
	p := component.ExporterCreateParams{Logger: logger}

	// Create a configuration for the exporter
	GRPClientConfig := createGRPCClientConfig("")
	config := createExporterConfig(&GRPClientConfig)

	// Test with a normal server
	server := badCloseServer{
		&serverBase{
			&av1.UnimplementedSpanHandlerServer{},
			badCloseListener{
				bufconn.Listen(1024 * 1024),
			},
		},
	}
	initServer(server)

	// Create the test data and the test function
	testCases := generateAllTraceTestCases()
	testDataPusher := createDataPusher(
		&config, p, grpc.WithContextDialer(server.lis.ToDial),
		grpc.WithTimeout(time.Second))

	droppedSpans, err := testDataPusher(context.Background(), testCases[3].td)
	assert.True(t, droppedSpans == 0)
	assert.True(t, err == nil)

	// flush, then check for the correct error
	logger.Sync()
	select {
	case b :=<-logAppeared: {
		assert.True(t, b)
	}
	case <-time.After(time.Second):
		t.Fatal("No log appeared")
	}
}

