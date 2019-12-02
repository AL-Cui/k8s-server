package routers

import (
	"k8s-server/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	APIs := beego.NewNamespace("/api",
		beego.NSNamespace("/pods",
			beego.NSInclude(
				&controllers.Pod{},	
			),
		),
	)
	beego.AddNamespaces(APIs)
}
