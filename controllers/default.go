package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/astaxie/beego/utils/captcha"
)
var cpt * captcha.Captcha

func init() {
	store := cache.NewMemoryCache()
	cpt = captcha.NewWithFilter("/captcha/", store)
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["WebSite"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"

	v :=  this.GetSession("userId")
	if v == nil {
		this.SetSession("userId", int(5))
		this.Data["userId"] = 5
	} else {
		this.Data["userId"] = v.(int) + 1
		
	}
	this.TplName = "index.tpl"
}

func (this *MainController) Article() {
	this.TplName = "article.tpl"
}