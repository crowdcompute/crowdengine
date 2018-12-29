package protocols

type (
	// Observer
	Observer interface {
		onNotify()
	}
	// Notifier
	Notifier interface {
		Register(Observer)
		Deregister(Observer)
		Notify()
	}
)
