package activitylog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"bitbucket.org/tunaiku/amargo-core/pkg/logger"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	rw.body.Write(body)
	return rw.ResponseWriter.Write(body)
}

func NewActivityLogMiddleware(publisher message.Publisher, cfg ActivityLogConfig) chi.Middlewares {
	return chi.Middlewares{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				ctx := r.Context()

				// Read the body
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Failed to read request body", http.StatusInternalServerError)
					return
				}

				// Restore the body so it can be read again
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				log := createTransactionLog(cfg, r)

				log.Start()

				// Restore the body again for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				ctx = NewContext(ctx, log)

				rw := &responseWriter{ResponseWriter: w}
				next.ServeHTTP(rw, r.WithContext(ctx))

				updateLogWithResponse(cfg, rw, log)

				log.End()

				if len(log.Activities) != 0 || cfg.IsPublishWhenNoActivities {
					PublishLog(ctx, publisher, log, cfg.TopicName)
				}

			})
		},
	}
}

func createTransactionLog(cfg ActivityLogConfig, r *http.Request) *Transaction {
	log := &Transaction{
		Service:   cfg.ServiceName,
		ActorType: cfg.ActorType,
		Target:    fmt.Sprintf("%s %s", r.Method, getRoutePattern(r)),
	}

	if cfg.IsRecordRequestBody {
		log.RequestBody = parseBody(r)
	}

	if cfg.IsRecordHeader {
		log.Header = parseHeader(r)
	}

	return log
}

func updateLogWithResponse(cfg ActivityLogConfig, r *responseWriter, log *Transaction) {
	if cfg.IsRecordResponseBody {
		body := make(map[string]interface{})
		json.Unmarshal(r.body.Bytes(), &body)
		log.ResponseBody = body
	}

	if cfg.IsRecordResponseCode {
		log.ResponseCode = r.statusCode
	}
}

func PublishLog(ctx context.Context, publisher message.Publisher, log *Transaction, topicName string) {
	log.EventID = uuid.New().String()

	logger.IWithTraceId(ctx).Debug("publishing activity log ", logrus.Fields{
		"activityLog": fmt.Sprintf("%+v", *log),
		"logID":       log.EventID,
	})

	msg := message.NewMessage(log.EventID, log.GetPayloadTransaction())
	err := publisher.Publish(topicName, msg)
	if err != nil {
		logger.IWithTraceId(ctx).Error("error publishing activity log ", logrus.Fields{
			"logID": log.EventID,
			"err":   err})
	}
}

func parseHeader(r *http.Request) map[string]interface{} {
	header := make(map[string]interface{})
	for k, v := range r.Header {
		header[k] = v
	}
	return header
}

func parseBody(r *http.Request) map[string]interface{} {
	body := make(map[string]interface{})
	if err := render.DecodeJSON(r.Body, &body); err != nil {
		return nil
	}
	return body
}

func parseResponse(r *httptest.ResponseRecorder) map[string]interface{} {
	body := make(map[string]interface{})
	json.Unmarshal(r.Body.Bytes(), &body)
	return body
}

func getRoutePattern(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())

	if pattern := rctx.RoutePattern(); pattern != "" {
		// Pattern is already available
		return pattern
	}
	routePath := r.URL.Path

	if r.URL.RawPath != "" {
		routePath = r.URL.RawPath
	}

	tctx := chi.NewRouteContext()

	if !rctx.Routes.Match(tctx, r.Method, routePath) {
		// No matching pattern, so just return the request path.
		// Depending on your use case, it might make sense to
		// return an empty string or error here instead
		return routePath
	}

	// tctx has the updated pattern, since Match mutates it
	return tctx.RoutePattern()

}

func ProcessVendorActivityLog(ctx context.Context, log *Transaction, publisher message.Publisher, cfg ActivityLogConfig) {
	ctx = NewContext(ctx, log)

	log.End()

	PublishLog(ctx, publisher, log, cfg.TopicName)
}
