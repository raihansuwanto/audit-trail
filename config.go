package audittrail

type ActivityLogConfig struct {
	ServiceName               string
	ActorType                 string
	ActorEmail                string
	TopicName                 string
	IsRecordRequestBody       bool
	IsRecordResponseBody      bool
	IsRecordHeader            bool
	IsRecordResponseCode      bool
	IsPublishWhenNoActivities bool
}
