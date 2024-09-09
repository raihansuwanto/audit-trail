package activitylog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func TestNewActivityLogMiddleware(t *testing.T) {
	publisher := gochannel.NewGoChannel(gochannel.Config{}, nil)
	cfg := ActivityLogConfig{
		ServiceName:               "testService",
		ActorType:                 "testActor",
		IsRecordRequestBody:       true,
		IsRecordHeader:            true,
		IsRecordResponseBody:      true,
		IsRecordResponseCode:      true,
		IsPublishWhenNoActivities: true,
		TopicName:                 "testTopic",
	}

	subscriber, _ := publisher.Subscribe(context.Background(), cfg.TopicName)

	middleware := NewActivityLogMiddleware(publisher, cfg)

	r := chi.NewRouter()
	r.Use(middleware...)
	r.Post("/test", func(w http.ResponseWriter, r *http.Request) {

		body := make(map[string]interface{})
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			t.Errorf("Error decoding request body: %v", err)
		}
		fmt.Println(body)

		trx := FromContext(r.Context())

		trx.SetTransactionEventType("testEvent").Start()

		segment := trx.StartAction("testAction", "testMessage")
		segment.SetTargetUserID("testUserID")
		segment.SetTargetBusinessID(123)
		segment.SetDataBefore("testDataBefore")
		segment.SetDataAfter("testDataAfter")
		segment.SetRequestData("testRequestData")
		segment.SetResponseData("testResponseData")
		segment.Succeed()
		segment.End()

		render.Status(r, 200)
		render.JSON(w, r, map[string]string{"message": "success"})
	})

	reqBody := map[string]string{"message": "test"}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(reqBodyBytes))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	msgs := <-subscriber

	var result Transaction
	json.Unmarshal(msgs.Payload, &result)

	if result.Service != "testService" {
		t.Errorf("Expected Service to be %s, but got %s", "testService", result.Service)
	}

	if result.TimeStart.IsZero() {
		t.Errorf("Expected TimeStart to be not zero, but got %v", result.TimeStart)
	}

	if result.TimeEnd.IsZero() {
		t.Errorf("Expected TimeEnd to be not zero, but got %v", result.TimeEnd)
	}

	if result.ActorType != "testActor" {
		t.Errorf("Expected Type to be %s, but got %s", "testActor", result.ActorType)
	}

	if result.Target != "POST /test" {
		t.Errorf("Expected Target to be %s, but got %s", "GET /test", result.Target)
	}

	if result.RequestBody["message"] != "test" {
		t.Errorf("Expected RequestBody to be %s, but got %s", "test", result.RequestBody["message"])
	}

	if result.ResponseBody["message"] != "success" {
		t.Errorf("Expected ResponseBody to be %s, but got %s", "success", result.ResponseBody["message"])
	}

	if result.ResponseCode != 200 {
		t.Errorf("Expected ResponseCode to be %d, but got %d", 200, result.ResponseCode)
	}

	if result.TargetUserID != "testUserID" {
		t.Errorf("Expected TargetUserID to be %s, but got %s", "testUserID", result.TargetUserID)
	}

	if result.TargetBusinessID != "123" {
		t.Errorf("Expected TargetBusinessID to be %s, but got %s", "123", result.TargetBusinessID)
	}

	if len(result.Activities) != 1 {
		t.Errorf("Expected 1 activity, but got %d", len(result.Activities))
	} else {
		activity := result.Activities[0]

		if activity.Action != "testAction" {
			t.Errorf("Expected Action to be %s, but got %s", "testAction", activity.Action)
		}

		if activity.Message != "testMessage" {
			t.Errorf("Expected Message to be %s, but got %s", "testMessage", activity.Message)
		}

		if activity.DataBefore != "testDataBefore" {
			t.Errorf("Expected DataBefore to be %s, but got %v", "testDataBefore", activity.DataBefore)
		}

		if activity.DataAfter != "testDataAfter" {
			t.Errorf("Expected DataAfter to be %s, but got %v", "testDataAfter", activity.DataAfter)
		}

		if activity.RequestData != "testRequestData" {
			t.Errorf("Expected RequestData to be %s, but got %v", "testRequestData", activity.RequestData)
		}

		if activity.ResponseData != "testResponseData" {
			t.Errorf("Expected ResponseData to be %s, but got %v", "testResponseData", activity.ResponseData)
		}

		if activity.Timestamp.IsZero() {
			t.Errorf("Expected Timestamp to be set, but it was not")
		}

		if activity.Status != "success" {
			t.Errorf("Expected Status to be %s, but got %s", "success", activity.Status)
		}
	}

}
