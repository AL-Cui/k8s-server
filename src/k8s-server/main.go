package main

import (
	_ "k8s-server/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

