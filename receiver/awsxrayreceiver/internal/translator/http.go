// Copyright The OpenTelemetry Authors
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

package translator

import (
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/translator/conventions"
	tracetranslator "go.opentelemetry.io/collector/translator/trace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/common/awsxray"
)

func addHTTP(seg *awsxray.Segment, span *pdata.Span) {
	if seg.HTTP == nil {
		return
	}

	attrs := span.Attributes()
	if req := seg.HTTP.Request; req != nil {
		// https://docs.aws.amazon.com/xray/latest/devguide/xray-api-segmentdocuments.html#api-segmentdocuments-http
		addString(req.Method, conventions.AttributeHTTPMethod, &attrs)

		if req.ClientIP != nil {
			// since the ClientIP is not nil, this means that this segment is generated
			// by a server serving an incoming request
			attrs.UpsertString(conventions.AttributeHTTPClientIP, *req.ClientIP)
		}

		addString(req.UserAgent, conventions.AttributeHTTPUserAgent, &attrs)
		addString(req.URL, conventions.AttributeHTTPURL, &attrs)
		addBool(req.XForwardedFor, awsxray.AWSXRayXForwardedForAttribute, &attrs)
	}

	if resp := seg.HTTP.Response; resp != nil {
		if resp.Status != nil {
			otStatus := tracetranslator.OCStatusCodeFromHTTP(int32(*resp.Status))
			// in X-Ray exporter, the segment status is set:
			// first via the span attribute, conventions.AttributeHTTPStatusCode
			// then the span status. Since we are also setting the span attribute
			// below, the span status code here will not be actually used
			span.Status().SetCode(pdata.StatusCode(otStatus))
			attrs.UpsertInt(conventions.AttributeHTTPStatusCode, *resp.Status)
		}

		addInt64(resp.ContentLength, conventions.AttributeHTTPResponseContentLength, &attrs)
	}

}
