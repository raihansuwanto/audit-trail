package activitylog

import (
	"encoding/json"
	"testing"
)

func TestTransactionAndSegmentMethods(t *testing.T) {
	t.Run("test transaction and segment methods", func(t *testing.T) {
		transaction := &Transaction{}

		// Test ITransaction methods
		transaction.SetTransactionEventType("testEvent")
		transaction.SetActor("testActor")
		payload := transaction.GetPayloadTransaction()

		if transaction.EventType != "testEvent" {
			t.Errorf("Expected EventType to be %s, but got %s", "testEvent", transaction.EventType)
		}

		if transaction.Actor != "testActor" {
			t.Errorf("Expected Actor to be %s, but got %s", "testActor", transaction.Actor)
		}

		expectedPayload, _ := json.Marshal(transaction)
		if string(payload) != string(expectedPayload) {
			t.Errorf("Expected payload to be %s, but got %s", string(expectedPayload), string(payload))
		}

		// Test ISegment methods
		segment := transaction.StartAction("testAction", "testMessage")
		segment.SetTargetUserID("testUserID")
		segment.SetTargetBusinessID(123)
		segment.SetDataBefore("testDataBefore")
		segment.SetDataAfter("testDataAfter")
		segment.Succeed()
		segment.End()

		if transaction.TargetUserID != "testUserID" {
			t.Errorf("Expected TargetUserID to be %s, but got %s", "testUserID", transaction.TargetUserID)
		}

		if transaction.TargetBusinessID != "123" {
			t.Errorf("Expected TargetBusinessID to be %s, but got %s", "123", transaction.TargetBusinessID)
		}

		if len(transaction.Activities) != 1 {
			t.Errorf("Expected 1 activity, but got %d", len(transaction.Activities))
		} else {
			activity := transaction.Activities[0]

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

			if activity.Status != "success" {
				t.Errorf("Expected Status to be %s, but got %s", "success", activity.Status)
			}

			if activity.Timestamp.IsZero() {
				t.Errorf("Expected Timestamp to be set, but it was not")
			}
		}
	})
}
