package units

import (
	"OceanID/app"
	pkgCtl "github.com/RealFax/pkg-ctl"
)

func init() {
	pkgCtl.RegisterHandler(100, "app-service", app.NewService)
}
