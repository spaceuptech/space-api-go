package realtime

import (
	"github.com/spaceuptech/space-api-go/api/model"
)

// LiveQuerySubscription represents the realtime subscription
type LiveQuerySubscription struct {
	unsubscribeFunc UnsubscribeFunction
	snapshot        *model.LiveData
}

func (l *LiveQuerySubscription) GetSnapshot() *model.LiveData {
	return l.snapshot
}

func (l *LiveQuerySubscription) Unsubscribe() () {
	l.unsubscribeFunc()
}
