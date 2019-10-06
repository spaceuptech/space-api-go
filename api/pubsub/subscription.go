package pubsub

type pubsubSubscription struct {
	subject             string
	unSubscribeFunction func()
}

func pubsubSubscriptionInit(subject string, unSubscribeFunction func()) *pubsubSubscription {
	return &pubsubSubscription{subject: subject, unSubscribeFunction: unSubscribeFunction}
}

func (p *pubsubSubscription) getSubject() string {
	return p.subject
}

func (p *pubsubSubscription) unSubscribe() {
	p.unSubscribeFunction()
}
