# goueditor
It shows the usage of [ueditorgoback](https://github.com/springlchy/ueditorgoback) within [beego](https://beego.me)

# Usage

Take beego framework for example.

## Step 1: Download UEditor

[http://ueditor.baidu.com/website/download.html](http://ueditor.baidu.com/website/download.html)

Choose the Latest(1.4.3.3) PHP version, UTF-8

## Step 2: Unpack
Unpack it to `plugins` directory

```
| + conf
| + controllers
| + models
| - plugins
    | - ueditor
   	   | + dialogs
       | + lang
       | - php
       	   config.json
       	   action_upload.php
       | + themes
       | + third-party
       index.html
       ueditor.all.js
       ueditor.all.min.js
       ueditor.config.job.js
       ueditor.config.js
       ueditor.parse.js
       ueditor.parse.min.js
| + routes
| + static
| + upload
| + views
```

Remove all the .php files under the plugins/ueditor/php for safety.

## Step 3: Configure ueditor

ueditor.config.js, line 33

```
, serverUrl: "php/controller.php"
```

change it to

```
, serverUrl: "/ueditor-upload/handle"
```

## Step 4: Download `ueditorgoback`

```
go get github.com/springlchy/ueditorgoback
```

## Step 5: wrap `ueditorgoback` within a controller

uedigor.go

``` go
package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/context" // DO NOT OMIT THIS
	"github.com/springlchy/ueditorgoback"
)

type UEditorController struct{
	beego.Controller
}

func (this *UEditorController) Handle() {
	ueditorgoback.HandleUpload(this.Ctx.ResponseWriter, this.Ctx.Request)
}

// DO NOT OMIT THIS
func (this *UEditorController) Render() error{
	return nil
}
```

## Step 6: Congifure router

Add this line to the `init()` function of routers.go
``` go
beego.Router("/ueditor-upload/handle", &controllers.UEditorController{}, "*:Handle")
```

## Step 7: Configure static file
Add the following lines to the `main()` function of main.go
``` go
beego.SetStaticPath("/plugins", "plugins")
beego.SetStaticPath("/upload", "upload")
```
If you use other server (nginx/apache) to serve the static file, this step can be skipped.
Configure the respective server (nginx/apache) instead.

## Last Step: Test

default.go (Irrelative code has been removed)
``` go
package controllers

import (
	"github.com/astaxie/beego"
)


func init() {

}

type MainController struct {
	beego.Controller
}

func (this *MainController) Article() {
	this.TplName = "article.tpl"
}
```

router.go
``` go
package routers

import (
	"github.com/astaxie/beego"
	"../controllers"
)

func init() {
	beego.AutoRouter(&controllers.MainController{})

	beego.Router("/ueditor-upload/handle", &controllers.UEditorController{}, "*:Handle")
}
```

article.tpl is under the `views` directory.

```
go build main.go
```
```
main
```

visit http://localhost:9090/main/article to make a test. Enjoy!
