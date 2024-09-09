# Audit Trail Logging

The Audit Trail Logging is a comprehensive solution for tracking and managing events and activities within your application. It provides a flexible and extensible framework for capturing detailed logs of user actions, system operations, and data changes, enabling robust auditing, monitoring, and troubleshooting capabilities.

## Features

- **Structured Event Logging**: The system allows for the creation of structured event logs, capturing a wide range of metadata such as event type, actor information, target resources, request/response data, and more.
- **Multi-level Activity Logging**: Events can be broken down into granular activities, each with its own set of detailed information, enabling a hierarchical view of the logged data.
- **Persistence and Retrieval**: The logged data can be persisted to various storage systems (e.g., databases, message brokers) for long-term retention and efficient retrieval.
- **Customizable Event Types**: Users can define custom event types and associated metadata to match the specific requirements of their application.
- **Flexible Logging Strategies**: The system supports different logging strategies, such as synchronous, asynchronous, or batch processing, to accommodate various performance and reliability needs.
- **Monitoring and Alerting**: The activity logs can be integrated with monitoring and alerting systems to detect and notify about critical events or anomalies.
- **Extensibility**: The system is designed to be easily extended, allowing for the integration of additional features, data sources, or custom data processing pipelines.

## Installation
To use the Activity Logging Middleware, add the following import to your Go project:

```bash
$ go get github.com/raihansuwanto/audit-trail/v1
```

### Usage

1. **Configure the Middleware**:
```go
import activitylog "github.com/raihansuwanto/audit-trail/v1"
   
   cfg := activitylog.ActivityLogConfig{
       ServiceName:               "your-service-name", // The name of the services
       ActorType:                 "user", // Set the default user that triggered the event
       ActorEmail:                "user@example.com", // Set the default Actor email
       TopicName:                 "activity-log-topic", 
       IsRecordRequestBody:       true,
       IsRecordResponseBody:      true,
       IsRecordHeader:            true,
       IsRecordResponseCode:      true,
       IsPublishWhenNoActivities: true,
   }
```

2. **Register the Middleware**:
```go
import activitylog "github.com/raihansuwanto/audit-trail/v1"

    r := chi.NewRouter()
    r.Use(activitylog.NewActivityLogMiddleware(publisher, cfg))
```
- Replace publisher with your message broker's publisher instance.

3. **Initialize Log**:
```go
import activitylog "github.com/raihansuwanto/audit-trail/v1"

    activitylog.FromContext(ctx).
        SetTransactionEventType("Update Data Project").
        SetActor("test123").
        SetActorEmail("example@mail.co")
```

3. **Start The Activities**:
```go
import activitylog "github.com/raihansuwanto/audit-trail/v1"

    activity := tx.StartAction("update", "update project A")
    
    activity.SetRequestData(userRequest)
    activity.SetDataBefore(map[string]interface{}{"date":"2024-08-08"})
    
    // Perform business logic

    activity.SetDataAfter(map[string]interface{}{"date":"2024-08-10"})
    activity.Succeed().End()
```

## Documentation
For more detailed usage examples, advanced configurations, and information on extending the Audit Trail Logging system, please refer to the documentation.

## Contributing
We welcome contributions to the Audit Trail Logging project! If you encounter any issues, have suggestions for improvements, or would like to contribute new features, please feel free to open an issue or submit a pull request on the GitHub repository.