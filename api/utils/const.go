package utils

const (
	// All is used when all the records needs to be worked on
	All string = "all"

	// One is used when oly a single record needs to be worked on
	One string = "one"

	// Count is used to count the number of documents returned
	Count string = "count"

	// Distinct is used to get the distinct values
	Distinct string = "distinct"

	// Upsert is used to upsert documents
	Upsert string = "upsert"

	// Delete is used to delete documents
	Delete string = "delete"

	// Update is used to update documents
	Update string = "update"

	// Create is used to create documents
	Create string = "create"
)

const (
	// Mongo is the constant for selecting MongoDB
	Mongo string = "mongo"

	// MySQL is the constant for selected MySQL
	MySQL string = "sql-mysql"

	// Postgres is the constant for selected Postgres
	Postgres string = "sql-postgres"
)

const (
	// TypeRealtimeSubscribe is the request type for live query subscription
	TypeRealtimeSubscribe string = "realtime-subscribe"

	// TypeRealtimeUnsubscribe is the request type for live query subscription
	TypeRealtimeUnsubscribe string = "realtime-unsubscribe"

	// TypeRealtimeFeed is the response type for realtime feed
	TypeRealtimeFeed string = "realtime-feed"

	// TypeServiceRegister is the request type for service registration
	TypeServiceRegister string = "service-register"

	// TypeServiceUnregister is the request type for service removal
	TypeServiceUnregister string = "service-unregister"

	// TypeServiceRequest is type triggering a service's function
	TypeServiceRequest string = "service-request"

	// TypePubsubSubscribe is type triggering a pubsub subscribe
	TypePubsubSubscribe string = "pubsub-subscribe"

	// TypePubsubSubscribeFeed is type having a pubsub subscribe feed
	TypePubsubSubscribeFeed string = "pubsub-subscribe-feed"

	// TypePubsubUnsubscribe is type triggering a pubsub unsubscribe
	TypePubsubUnsubscribe string = "pubsub-unsubscribe"

	// TypePubsubUnsubscribeAll is type triggering a pubsub unsubscribe all
	TypePubsubUnsubscribeAll string = "pubsub-unsubscribe-all"
)

const (
	// RealtimeInsert is for create operations
	RealtimeInsert string = "insert"

	// RealtimeUpdate is for update operations
	RealtimeUpdate string = "update"

	// RealtimeDelete is for delete operations
	RealtimeDelete string = "delete"

	// RealtimeInitial is for initial operations
	RealtimeInitial string = "initial"
)
const (
	// PayloadSize is the size of the payload(in bytes) in file upload and download
	PayloadSize int = 256*1024 // 256 kB
)