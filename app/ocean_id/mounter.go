package oceanId

type Mounter interface {
	SetOI(oi IdPool)
}

func Mount[T any](oi IdPool, service T) T {
	mounter := any(service).(Mounter)
	mounter.SetOI(oi)
	return service
}
