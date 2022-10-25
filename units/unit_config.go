package units

import (
	"OceanID/config"
	pkgCtl "github.com/RealFax/pkg-ctl"
)

func init() {
	pkgCtl.RegisterHandler(1, "config-loader", config.NewLoader)
}
