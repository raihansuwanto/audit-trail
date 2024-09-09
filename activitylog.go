package audittrail

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
)

type Transaction struct {

	// EventID is enum string.
	// It is used to store the event ID of the event log.
	EventID string `json:"eventID"`

	// EventType is enum string.
	// EventType is used to store the type of the event log.
	EventType string `json:"eventType"`

	// Service is used to store the specific service that called the event log.
	Service string `json:"service"`

	// Actor is used to store the keycloak id of the actor.
	Actor string `json:"actor"`

	// ActorEmail is used to store the email of the actor.
	ActorEmail string `json:"actorEmail"`

	// ActorType is used to store the type of the actor.
	ActorType string `json:"actorType"`

	// TargetUserID is used to store the keycloak id of the target user.
	TargetUserID string `json:"targetUserId"`

	// TargetBusinessID is used to store the business id of the target user.
	TargetBusinessID string `json:"targetBusinessId"`

	// Target is used to store the endpoint of the event log.
	Target string `json:"target"`

	// Header is used to store the header of the event log request.
	Header map[string]interface{} `json:"header"`

	// RequestBody is used to store the request body of the event log.
	RequestBody map[string]interface{} `json:"requestBody"`

	// ResponseBody is used to store the response body of the event log.
	ResponseBody map[string]interface{} `json:"responseBody"`

	// ResponseCode is used to store the response code of the event log.
	ResponseCode int `json:"responseCode"`

	// Activities is used to store the detail activities of the event log.
	Activities []Activity `json:"activities"`

	TimeStart time.Time `json:"timeStart"`
	TimeEnd   time.Time `json:"timeEnd"`

	// Resource
	Resource string `json:"resource"`

	// type
	Type string `json:"type"`

	// IsHtppMiddleware is used to determine whether the event log is created by the middleware.
	IsHtppMiddleware bool `json:"-"`

	// Publisher is used to store the publisher of the event log.
	Publisher message.Publisher `json:"-"`
}

type Activity struct {

	// Action is used to store the action of the action log.
	Action string `json:"action"`

	// Message is used to store the message of the action log.
	// The message will be displayed on CRM.
	Message string `json:"message"`

	// Status is used to store the status of the action log.
	// The status can be either "success" or "failed".
	Status string `json:"status"`

	// RequestData is used to store the request data of the action log.
	RequestData interface{} `json:"requestData"`

	// ResponseData is used to store the response data of the action log.
	ResponseData interface{} `json:"responseData"`

	// DataBefore is used to store the data before the action log.
	DataBefore interface{} `json:"dataBefore"`

	// DataAfter is used to store the data after the action log.
	DataAfter interface{} `json:"dataAfter"`

	// Timestamp is used to store the timestamp of the action log.
	Timestamp time.Time `json:"timestamp"`

	// IsVisible is used to determine whether the activity log is visible to the user.
	IsVisible bool `json:"isVisible"`
}

type ITransaction interface {

	// Start starts a new event log.
	Start() ITransaction

	// End ends the event log.
	End()

	// SetTransactionEventType sets the event type of the event log.
	SetTransactionEventType(eventType string) ITransaction

	// SetActor sets the actor of the event log.
	SetActor(actor string) ITransaction

	// SetActorEmail sets the actor email of the event log.
	SetActorEmail(actorEmail string) ITransaction

	// SetActorType sets the actor type of the event log.
	SetActorType(actorType string) ITransaction

	// GetPayloadTransaction returns the payload byte of the event log.
	GetPayloadTransaction() []byte

	// SetHeader
	SetHeader(header map[string]interface{}) ITransaction

	// Set Type
	SetType(typeString string) ITransaction

	// Set Resource
	SetResource(resource string) ITransaction

	//Publisher is used to send the event log to the message broker. (Google PubSub)
	Publish(topicName string)
}

// Start a new event log
func (c *Transaction) Start() ITransaction {
	c.TimeStart = time.Now()
	return c
}

// End the event log
func (c *Transaction) End() {
	c.TimeEnd = time.Now()
}

// To create a new event log
func (c *Transaction) SetTransactionEventType(eventType string) ITransaction {
	c.EventType = eventType
	return c
}

// To Set Actor of the event log
func (c *Transaction) SetActor(actor string) ITransaction {
	c.Actor = actor
	return c
}

// To Set Actor Email of the event log
func (c *Transaction) SetActorEmail(actorEmail string) ITransaction {
	c.ActorEmail = actorEmail
	return c
}

func (c *Transaction) SetActorType(actorType string) ITransaction {
	c.ActorType = actorType
	return c
}

func (c *Transaction) SetHeader(header map[string]interface{}) ITransaction {
	c.Header = header
	return c
}

func (c *Transaction) SetType(typeString string) ITransaction {
	c.Type = typeString
	return c
}

func (c *Transaction) SetResource(resource string) ITransaction {
	c.Resource = resource
	return c
}

// To get payload byte of the event log
func (c *Transaction) GetPayloadTransaction() []byte {
	payload, _ := json.Marshal(c)
	return payload
}

// To start a new action log before append it to the event log.
// Make sure to add the event log to the context after creating it.
//
// Example:
//
//	activitylog.NewContext(ctx, &src.Transaction{
//		EventID: "uuid",
//		Service: "POST /api/v1/user",
//	})
//
//	tx := actifitylog.FromContext(ctx).StartAction("create", "create user")
//	defer tx.End()
//
//	// ... function code here ...
//
//	tx.Succeed()
func (c *Transaction) StartAction(action string, message string) *Segment {
	return &Segment{
		root: c,
		Activity: Activity{
			Action:  action,
			Message: message,
			Status:  "failed",
		},
	}
}

// Publish
func (c *Transaction) Publish(topicName string) {
	PublishLog(context.Background(), c.Publisher, c, topicName)
}

type Segment struct {
	root     *Transaction
	Activity Activity
}

type ISegment interface {

	// End ends the action log.
	End()

	// SetRequestData sets the request data of the action log.
	SetRequestData(data interface{}) ISegment

	// SetResponseData sets the response data of the action log.
	SetResponseData(data interface{}) ISegment

	// SetDataAfter sets the data after change.
	SetDataAfter(data interface{}) ISegment

	// SetDataBefore sets the data before change.
	SetDataBefore(data interface{}) ISegment

	// SetTargetBusinessID sets the target business ID of the action log.
	SetTargetBusinessID(businessID int64) ISegment

	// SetTargetUserID sets the target user ID of the action log.
	SetTargetUserID(userID string) ISegment

	// SetVisibility sets the visibility of the action log.
	SetVisibility(isVisible bool) ISegment

	// Succeed marks the action log as successful.
	Succeed() ISegment
}

// To set the actor keycloak ID of the event log
func (c *Segment) SetTargetUserID(userID string) ISegment {
	c.root.TargetUserID = userID
	return c
}

// To set the actor business ID of the event log
func (c *Segment) SetTargetBusinessID(businessID int64) ISegment {
	c.root.TargetBusinessID = strconv.Itoa(int(businessID))
	return c
}

// To set the data before change
func (c *Segment) SetDataBefore(data interface{}) ISegment {
	c.Activity.DataBefore = data
	return c
}

// To set the data after change
func (c *Segment) SetDataAfter(data interface{}) ISegment {
	c.Activity.DataAfter = data
	return c
}

// To set the request data of the action log
func (c *Segment) SetRequestData(data interface{}) ISegment {
	c.Activity.RequestData = data
	return c
}

// To set the response data of the action log
func (c *Segment) SetResponseData(data interface{}) ISegment {
	c.Activity.ResponseData = data
	return c
}

// To set the request data of the action log
func (c *Segment) Succeed() ISegment {
	c.Activity.Status = "success"
	return c
}

// To set the visibility of the action log
func (c *Segment) SetVisibility(isVisible bool) ISegment {
	c.Activity.IsVisible = isVisible
	return c
}

// To Append the action log to the event log
func (c *Segment) End() {
	c.Activity.Timestamp = time.Now()
	c.root.Activities = append(c.root.Activities, c.Activity)
}
