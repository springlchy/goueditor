package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/context"
	"github.com/springlchy/ueditorgoback"
)

type UEditorController struct{
	beego.Controller
}

func (this *UEditorController) Handle() {
	ueditorgoback.HandleUpload(this.Ctx.ResponseWriter, this.Ctx.Request)
}

func (this *UEditorController) Render() error{
	return nil
}