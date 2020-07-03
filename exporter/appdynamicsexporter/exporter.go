package appdynamicsexporter

import (
	"context"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	apiclient "github.com/pavankrish123/cms-openapi/v1/go/client"
	"github.com/pavankrish123/cms-openapi/v1/go/client/data"
	cms "github.com/pavankrish123/cms-openapi/v1/go/models"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/consumer/pdatautil"
	"strconv"
)

const (
	OTResourceType = "ot_resource"

	// keywords used in AppD Payload
	AppDKey         = "key"
	AppDType        = "type"
	AppDKind        = "kind"
	AppDDescription = "description"
	AppDUnit        = "unit"
)

type AppDMetricsTransformer func(metrics pdata.Metrics) *cms.Payload

type CMSExporter struct {
	params component.ExporterCreateParams
}


func (e *CMSExporter) pDataToCMS(md pdata.Metrics) (*cms.Payload, int) {
	entities := make([]*cms.Entity, 0)
	totalDropped := 0
	idm := pdatautil.MetricsToInternalMetrics(md)

	for ix :=0; ix < idm.ResourceMetrics().Len(); ix++ {
		e.params.Logger.Info("received" +  strconv.Itoa(idm.ResourceMetrics().Len()) + "resource metrics")
		resourceMetric := idm.ResourceMetrics().At(ix)
		resource := resourceMetric.Resource()
		if resourceMetric.IsNil() || resource.IsNil() {
			e.params.Logger.Warn(fmt.Sprintf("Dropping metric %v %v", resourceMetric.IsNil(), resource.IsNil()))
			totalDropped += 1
			return nil, totalDropped
		}
		props := e.cmsPropsfromResourceMetricAttrs(resource.Attributes())
		metrics, dropped := e.cmsMetricsFromIMetrics(resourceMetric.InstrumentationLibraryMetrics())
		if dropped > 0 {
			totalDropped += 1
		}
		otType := OTResourceType
		entities = append(entities, &cms.Entity{
			Attributes: nil,
			Metrics:    metrics,
			Properties: props,
			Type:       &otType,
		})
	}

	return &cms.Payload{
		Entities:      entities,
		Labels:        nil,
		Relationships: nil,
	}, totalDropped
}


func (e *CMSExporter) cmsMetricsFromIMetrics(imslice pdata.InstrumentationLibraryMetricsSlice) ([]*cms.Metric, int){
	var droppedSeries int
	var meters []*cms.Metric
	for ix := 0; ix < imslice.Len(); ix ++ {
		e.params.Logger.Debug(fmt.Sprintf("received %v I metrics", imslice.Len()))
		if im := imslice.At(ix); !im.IsNil() {
			if !im.InstrumentationLibrary().IsNil(){
				e.params.Logger.Info("INAME:" + im.InstrumentationLibrary().Name())
			}
			for jx := 0; jx < im.Metrics().Len(); jx++ {
				if mtr := im.Metrics().At(jx); !mtr.IsNil() && !mtr.MetricDescriptor().IsNil(){
					dp, ok := e.cmsMetricFromPdataMetric(&mtr)
					if !ok {
						droppedSeries += 1
					}
					meters = append(meters, dp)
				}
			}
		}
	}
	return meters, droppedSeries
}

func (e *CMSExporter) cmsMetricFromPdataMetric(mtr *pdata.Metric)(*cms.Metric, bool){
	points, ok := e.cmsPointsFromPdataMetric(mtr)
	if !ok{
		return nil, false
	}
	labels := e.cmsMetricLabelsFromMetricDescriptor(mtr.MetricDescriptor())
	e.params.Logger.Debug(fmt.Sprintf("received %v labels", labels))

	return &cms.Metric{Labels: labels, Points: points}, true
}

func (e *CMSExporter)  cmsPointsFromPdataMetric(mtr *pdata.Metric) ([]*cms.DataPoint, bool) {
	var mDataPoints []*cms.DataPoint
	switch mtr.MetricDescriptor().Type() {
	case pdata.MetricTypeMonotonicInt64, pdata.MetricTypeInt64:
		dpSeries :=   mtr.Int64DataPoints()
		for i := 0; i < dpSeries.Len(); i++ {
			timestamp := fmt.Sprintf("%v", dpSeries.At(i).Timestamp())
			val := fmt.Sprintf("%v", dpSeries.At(i).Value())
			mDataPoints = append(mDataPoints, &cms.DataPoint{
				Timestamp: &timestamp,
				Value:     &val,
			})
		}
		return mDataPoints, true
	case pdata.MetricTypeMonotonicDouble, pdata.MetricTypeDouble:
		dpSeries :=   mtr.DoubleDataPoints()
		for i := 0; i < dpSeries.Len(); i++ {
			timestamp := fmt.Sprintf("%v", dpSeries.At(i).Timestamp())
			val := fmt.Sprintf("%v", dpSeries.At(i).Value())
			mDataPoints = append(mDataPoints, &cms.DataPoint{
				Timestamp: &timestamp,
				Value:     &val,
			})
		}
		return mDataPoints, true
	case pdata.MetricTypeHistogram, pdata.MetricTypeSummary:
		return nil, false
	default:
		return nil, false
	}
}

func (e *CMSExporter)  cmsMetricLabelsFromMetricDescriptor(mtd pdata.MetricDescriptor) cms.Map {
	return map[string]string{
		AppDKey:         mtd.Name(),
		AppDType:        mtd.Type().String(),
		AppDKind:        mtd.Type().String(),
		AppDDescription: mtd.Description(),
		AppDUnit:        mtd.Unit(),
	}
}

func (e *CMSExporter)  cmsPropsfromResourceMetricAttrs(attrs pdata.AttributeMap) cms.Map {
	props := map[string]string{}
	attrs.ForEach(func(k string, v pdata.AttributeValue) {
		// checking v.IsNil() not supported - this may lead to panic
		switch v.Type() {
		case pdata.AttributeValueBOOL:
			props[k] = fmt.Sprintf("%v", v.BoolVal())
		case pdata.AttributeValueDOUBLE:
			props[k] = fmt.Sprintf("%v", v.DoubleVal())
		case pdata.AttributeValueSTRING:
			props[k] = v.StringVal()
		case pdata.AttributeValueINT:
			props[k] = fmt.Sprintf("%v", v.IntVal())
		default:
			return
		}
	})
	return props
}

func (e *CMSExporter)  pushPayloadToCMS(ctx context.Context, config *Config, params component.ExporterCreateParams, payload *cms.Payload) error {
	cl, err := config.Metrics.HTTP.ToClient()
	if err != nil {
		params.Logger.Info(err.Error())
		return err
	}

	dp := &data.IngestDataParams{
		AppdynamicsAPIKey:  &config.AccessKey,
		AppdynamicsAccount: &config.Account,
		Body:               payload,
		Context:            ctx,
		HTTPClient:         cl,
	}

	transport := httptransport.New(config.Metrics.HTTP.Endpoint, "", nil)
	apiclient.Default.SetTransport(transport)
	if resp, err := apiclient.Default.Data.IngestData(dp); err != nil {
		params.Logger.Info(err.Error())
		if resp != nil {
			params.Logger.Info(resp.Error())
		}
		return err
	}
	return nil
}
