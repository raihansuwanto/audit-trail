package activitylog

import "context"

const (
	ActivityLogCtx = "activity_log"
)

// Set the event log to the context
func NewContext(ctx context.Context, log *Transaction) context.Context {
	return context.WithValue(ctx, ActivityLogCtx, log)
}

// Get the event log from the context
func FromContext(ctx context.Context) *Transaction {
	if nil == ctx {
		return nil
	}
	h, _ := ctx.Value(ActivityLogCtx).(*Transaction)
	return h
}
