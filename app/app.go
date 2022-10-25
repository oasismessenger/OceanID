package app

type Application interface {
	GetName() string
	Setup() error
	Start() error
	Shutdown() error
}
