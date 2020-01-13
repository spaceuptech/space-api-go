package model

import (
	"encoding/json"

	"github.com/spaceuptech/space-api-go/utils"
)

const RealtimeUnsubscribe string = "realtime-unsubscribe"

type DbStore map[string]ColStore
type ColStore map[string]IdStore
type IdStore map[string]*Store

// RealtimeRequest is the object sent for realtime requests
type RealtimeRequest struct {
	Token   string                 `json:"token"`
	DBType  string                 `json:"dbType"`
	Project string                 `json:"project"`
	Group   string                 `json:"group"` // Group is the collection name
	Type    string                 `json:"type"`  // Can either be subscribe or unsubscribe
	ID      string                 `json:"id"`    // id is the query id
	Where   map[string]interface{} `json:"where"`
	Options *LiveQueryOptions      `json:"options"`
}

// LiveQueryOptions is used to set the options for the live query
type LiveQueryOptions struct {
	ChangesOnly bool
	SkipInitial bool `json:"skipInitial"`
}

type Store struct {
	QueryOptions       LiveQueryOptions
	Snapshot           []*SnapshotData
	C                  chan *SubscriptionEvent
	SubscriptionObject StoreSubscriptionObject
	Find               interface{}
	Options            interface{}
}

func LiveQuerySubscriptionInit(unsubscribeFunc func(), snapshot []SnapshotData, c chan *SubscriptionEvent) StoreSubscriptionObject {
	return StoreSubscriptionObject{unsubscribeFunc: unsubscribeFunc, snapshot: snapshot, C: c}
}

func (l *StoreSubscriptionObject) GetSnapshot() []SnapshotData {
	return l.snapshot
}

func (l *StoreSubscriptionObject) Unsubscribe() {
	l.unsubscribeFunc()
}

type RealtimeParams struct {
	Find utils.M
}

type StoreSubscriptionObject struct {
	unsubscribeFunc func()
	snapshot        []SnapshotData
	C               chan *SubscriptionEvent
}

type SubscriptionEvent struct {
	err    error
	find   map[string]interface{}
	doc    interface{}
	evType string
}

func NewSubscriptionEvent(evType string, doc interface{}, find map[string]interface{}, err error) *SubscriptionEvent {
	return &SubscriptionEvent{evType: evType, doc: doc, find: find, err: err}
}

func (s *SubscriptionEvent) Unmarshal(vPtr interface{}) error {
	data, err := json.Marshal(s.doc)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, vPtr)
}

func (s *SubscriptionEvent) Type() string {
	return s.evType
}

func (s *SubscriptionEvent) Find() map[string]interface{} {
	return s.find
}

func (s *SubscriptionEvent) Err() error {
	return s.err
}

type SnapshotFunction func(interface{}, string, interface{})
type OnErrorFunction func(error)

type SnapshotData struct {
	Find      map[string]interface{}
	Time      int64
	Payload   interface{}
	IsDeleted bool
}
