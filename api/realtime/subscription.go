package realtime

import (
	"github.com/spaceuptech/space-api-go/api/model"
)

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
