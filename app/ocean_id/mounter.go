package oceanID

type Mounter interface {
	SetOI(oi IDPool)
}

func Mount[T any](oi IDPool, service T) T {
	mounter := any(service).(Mounter)
	mounter.SetOI(oi)
	return service
}
