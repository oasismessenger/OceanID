package oceanID

type Mounter interface {
	SetOI(oi OI)
}

func Mount[T any](oi OI, service T) T {
	mounter := any(service).(Mounter)
	mounter.SetOI(oi)
	return service
}
