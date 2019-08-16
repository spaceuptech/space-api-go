package pubsub

// PubsubSubscription represents the pubsub subscription
type PubsubSubscription struct {
	Subject         string
	unsubscribeFunc UnsubscribeFunction
}

func (p *PubsubSubscription) Unsubscribe() () {
	p.unsubscribeFunc()
}
