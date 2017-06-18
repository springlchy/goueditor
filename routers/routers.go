package routers

import (
	"github.com/astaxie/beego"
	"../controllers"
)

func init() {
	beego.AutoRouter(&controllers.MainController{})

	beego.Router("/ueditor-upload/handle", &controllers.UEditorController{}, "*:Handle")
}