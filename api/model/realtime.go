package model

import "github.com/spaceuptech/space-api-go/api/utils"

const RealtimeUnsubscribe string = "realtime-unsubscribe"


type TypeStore map[string]map[string]map[string]Store

// RealtimeRequest is the object sent for realtime requests
type RealtimeRequest struct {
	Token   string                 `json:"token"`
	DBType  string                 `json:"dbType"`
	Project string                 `json:"project"`
	Group   string                 `json:"group"` // Group is the collection name
	Type    string                 `json:"type"`  // Can either be subscribe or unsubscribe
	ID      string                 `json:"id"`    // id is the query id
	Where   map[string]interface{} `json:"where"`
	Options *LiveQueryOptions       `json:"options"`
}

// LiveQueryOptions is used to set the options for the live query
type LiveQueryOptions struct {
	ChangesOnly bool
	SkipInitial bool `json:"skipInitial"`
}

type Store struct {
	QueryOptions       LiveQueryOptions
	Snapshot           []SnapshotData
	Subscription       StoreSubscription
	SubscriptionObject StoreSubscriptionObject
	Find               interface{}
	Options            interface{}
}

type RealtimeParams struct {
	Find utils.M
}

type StoreSubscription struct {
	OnSnapShot SnapshotFunction
	OnError    OnErrorFunction
}

type StoreSubscriptionObject struct {
	SnapShot []SnapshotData
}

type SnapshotFunction func(interface{}, string, interface{})
type OnErrorFunction func(error)

type SnapshotData struct {
	Id        string
	Time      int64
	Payload   interface{}
	IsDeleted bool
}
