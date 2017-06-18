package controllers

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"encoding/json"
	"os"
	"path"
	"fmt"
	"mime/multipart"
	"time"
	"math/rand"
	"io"
	"net/http"
	"strconv"
)


type UEditorConfig map[string]interface{}
type UploadConfig struct {
	pathFormat string
	maxSize int64
	allowFiles []string
	origName string
}

var ueditorConfig UEditorConfig
var serverPath string

var stateMap = map[string]string{
	"0": "SUCCESS",
	"1" : "文件大小超出 upload_max_filesize 限制",
	"2" : "文件大小超出 MAX_FILE_SIZE 限制",
	"3" : "文件未被完整上传",
	"4" : "没有文件被上传",
	"5" : "上传文件为空",
    "ERROR_TMP_FILE" : "临时文件错误",
    "ERROR_TMP_FILE_NOT_FOUND" : "找不到临时文件",
    "ERROR_SIZE_EXCEED" : "文件大小超出网站限制",
    "ERROR_TYPE_NOT_ALLOWED" : "文件类型不允许",
    "ERROR_CREATE_DIR" : "目录创建失败",
    "ERROR_DIR_NOT_WRITEABLE" : "目录没有写权限",
    "ERROR_FILE_MOVE" : "文件保存时出错",
    "ERROR_FILE_NOT_FOUND" : "找不到上传文件",
    "ERROR_WRITE_CONTENT" : "写入文件内容错误",
    "ERROR_UNKNOWN" : "未知错误",
    "ERROR_DEAD_LINK" : "链接不可用",
    "ERROR_HTTP_LINK" : "链接不是http链接",
    "ERROR_HTTP_CONTENTTYPE" : "链接contentType不正确",
}

type Size interface {
	Size() int64
}

type Stat interface {
	Stat() (os.FileInfo, error)
}

func init() {
	serverPath = getCurrentPath()
	ueditorConfig = make(UEditorConfig)
	readConfig("/plugins/ueditor/php/config.json")
}

func main() {
	fmt.Printf(ueditorConfig["imageActionName"].(string))
}

type UEditorController struct{
	beego.Controller
}

type Uploader struct {
	request *http.Request
	fileField string
	file multipart.File
	base64 string
	config UEditorConfig
	oriName string
	fileName string
	fullName string
	filePath string
	fileSize int64
	fileType string
	stateInfo string
	optype string
}

func (this *UEditorController) Handle() {
	var fieldName string = "upload"
	var config UEditorConfig
	var fbase64 string
	var result map[string]string

	action := this.GetString("action")
	switch (action) {
	case "config":
		this.Data["json"] = &ueditorConfig
		this.ServeJSON()
		return
	case "uploadimage":
		config = UEditorConfig{
			"pathFormat" : ueditorConfig["imagePathFormat"],
			"maxSize" : ueditorConfig["imageMaxSize"],
			"allowFiles" : ueditorConfig["imageAllowFiles"],
		}
		fieldName = ueditorConfig["imageFieldName"].(string)
	case "uploadscrawl":
		config = UEditorConfig{
			"pathFormat" : ueditorConfig["scrawPathFormat"],
			"maxSize" : ueditorConfig["scrawMaxSize"],
			"allowFiles" : ueditorConfig["scrawlAllowFiles"],
			"oriname" : "scrawl.png",
		}
		fieldName = ueditorConfig["scrawFieldName"].(string)
		fbase64 = "base64"
	case "uploadvideo":
		config = UEditorConfig{
			"pathFormat" : ueditorConfig["videoPathFormat"],
			"maxSize": ueditorConfig["videoMaxSize"],
			"allowFiles" : ueditorConfig["videoAllowFiles"],
		}
		fieldName = ueditorConfig["videoFieldName"].(string)
	case "uploadfile":
		config = UEditorConfig{
			"pathFormat" : ueditorConfig["filePathFormat"],
			"maxSize" : ueditorConfig["fileMaxSize"],
			"allowFiles" : ueditorConfig["fileAllowFiles"],
		}
		fieldName = ueditorConfig["fileFieldName"].(string)
	default:
		this.Data["json"] = &map[string]string{
			"state" : "请求地址出错",
		}
	}
	config["maxSize"] = int64(config["maxSize"].(float64))
	uploader := NewUploader(this.Ctx.Request, fieldName, config, fbase64)

	uploader.upFile()
	result = uploader.getFileInfo()

	this.Data["json"] = &result
	this.ServeJSON()
}

func NewUploader(request *http.Request, fileField string, config UEditorConfig, optype string) (uploader *Uploader) {
	uploader = new(Uploader)
	uploader.request = request
	uploader.fileField = fileField
	uploader.config = config
	uploader.optype = optype

	return
}

func (this *Uploader) upFile() {
	this.request.ParseMultipartForm(this.config["maxSize"].(int64))
	file, fheader, err := this.request.FormFile(this.fileField)
	defer file.Close()
	if err != nil {
		this.stateInfo = err.Error()
		fmt.Printf("upload file error: %s", err)
	} else {
		this.oriName = fheader.Filename
		if stateInterface, ok := file.(Stat); ok {
			fileInfo, _ := stateInterface.Stat()
			this.fileSize = fileInfo.Size()
		} else {
			this.stateInfo = this.getStateInfo("ERROR_UNKNOWN")
		}

		this.fileType = this.getFileExt()
		this.fullName = this.getFullName()
		this.filePath = this.getFilePath()
		this.fileName = this.getFileName()

		dirname := path.Dir(this.filePath)

		if ! this.checkSize() {
			this.stateInfo = this.getStateInfo("ERROR_SIZE_EXCEED")
			return
		}

		if ! this.checkType() {
			this.stateInfo = this.getStateInfo("ERROR_TYPE_NOT_ALLOWED")
			return
		}

		dirInfo, err := os.Stat(dirname)
		if err != nil {
			err = os.MkdirAll(dirname, 0666)
			if err != nil {
				this.stateInfo = this.getStateInfo("ERROR_CREATE_DIR")
				fmt.Printf("Error create dir: %s", err)
				return
			}
		} else if dirInfo.Mode() & 0222 == 0 {
			this.stateInfo = this.getStateInfo("ERROR_DIR_NOT_WRITEABLE")
			return
		}

		fout, err := os.OpenFile(this.filePath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			this.stateInfo = this.getStateInfo("ERROR_FILE_MOVE");
			return
		}
		defer fout.Close()

		io.Copy(fout, file)
		// if err != nil {
		// 	this.stateInfo = this.getStateInfo("ERROR_FILE_MOVE");
		// 	return
		// }

		this.stateInfo = stateMap["0"]
	}
}

func (this *Uploader) getStateInfo(errCode string) string {
	if errInfo, ok := stateMap[errCode]; ok {
		return errInfo
	} else {
		return stateMap["ERROR_UNKNOWN"]
	}
}

func (this *Uploader) getFileExt() string {
	pos := strings.LastIndex(this.oriName,  ".")
	return strings.ToLower(this.oriName[pos:])
}

func (this *Uploader) getFullName() string {
	t := time.Now()
	format := this.config["pathFormat"].(string)
	format = strings.Replace(format, "{yyyy}", strconv.Itoa(t.Year()), 1)
	format = strings.Replace(format, "{mm}", strconv.Itoa(int(t.Month())), 1)
	format = strings.Replace(format, "{dd}", strconv.Itoa(t.Day()), 1)
	format = strings.Replace(format, "{hh}", strconv.Itoa(t.Hour()), 1)
	format = strings.Replace(format, "{ii}", strconv.Itoa(t.Minute()), 1)
	format = strings.Replace(format, "{ss}", strconv.Itoa(t.Second()), 1)
	format = strings.Replace(format, "{time}", strconv.FormatInt(t.Unix(), 10), 1)

	randNum := strconv.Itoa(rand.Intn(10000000)+10000000)
	format = format + randNum

	return format + this.getFileExt()
}

func (this *Uploader) getFileName() string {
	pos := strings.LastIndex(this.filePath, "/");
	return this.filePath[pos+1:]
}

func (this *Uploader) getFilePath() string {
	fullname := this.fullName

	if strings.LastIndex(fullname, "/") == (len(fullname)-1) {
		fullname = "/" + fullname
	}

	return serverPath + fullname
}

func (this *Uploader) checkType() bool {
	found := false;
	for _, ext := range this.config["allowFiles"].([]interface{}) {
		if ext.(string) == this.fileType {
			found = true
			break
		}
	}

	return found
}

func (this *Uploader) checkSize() bool {
	return this.fileSize <= (this.config["maxSize"].(int64))
}

func (this *Uploader) getFileInfo() map[string]string {
	return map[string]string{
		"state" : this.stateInfo,
		"url" : this.fullName,
		"title" : this.fileName,
		"original" : this.oriName,
		"type" : this.fileType,
		"size" : string(this.fileSize),
	}
}

func readConfig(configPath string) {
	fd, err := os.Open(serverPath + "/" + configPath)
	checkError(err)
	defer fd.Close()

	data, err := ioutil.ReadAll(fd)
	checkError(err)

	pattern := regexp.MustCompile(`/\*.+?\*/`)
	data = pattern.ReplaceAll(data, []byte(""))

	json.Unmarshal(data, &ueditorConfig)
	fmt.Println(string(data))
}

func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	checkError(err)
	
	return strings.Replace(dir, "\\", "/", -1)
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
