package realtime

import (
	"github.com/spaceuptech/space-api-go/model"
)

// LiveQuerySubscription represents the realtime subscription
type LiveQuerySubscription struct {
	unsubscribeFunc func()
	snapshot        model.SnapshotFunction
}

func LiveQuerySubscriptionInit(snapshot model.SnapshotFunction, unsubscribeFunc func()) *LiveQuerySubscription {
	return &LiveQuerySubscription{unsubscribeFunc: unsubscribeFunc, snapshot: snapshot}
}

func (l *LiveQuerySubscription) GetSnapshot() model.SnapshotFunction {
	return l.snapshot
}

func (l *LiveQuerySubscription) Unsubscribe() () {
	l.unsubscribeFunc()
}
