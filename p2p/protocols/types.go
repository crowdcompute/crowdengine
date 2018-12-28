package protocols

type (
	Observer interface {
		onNotify()
	}

	Notifier interface {
		Register(Observer)
		Deregister(Observer)
		Notify()
	}
)
