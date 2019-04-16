package huaweiPush

import (
	"encoding/json"
	"fmt"
	"os"
)

// url
const (
	TOKEN_URL = "https://login.cloud.huawei.com/oauth2/v2/token"
	PUSH_URL  = "https://api.push.hicloud.com/pushsend.do"
)

// config
const (
	GRANTTYPE = "client_credentials"
	NSP_SVC   = "openpush.message.api.send"
)

/**
 **************************************** 结构体
 */

type HuaweiPushClient struct {
	ClientId     string
	ClientSecret string
	NspCtx       string
}

type Vers struct {
	Ver   string `json:"ver"`
	AppID string `json:"appId"`
}

type TokenResStruct struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Token_type   string `json:"token_type"`
}

var (
	PostRetryTimes = 3
)

/**
 **************************************** 消息体
 */

type Message struct {
	Hps Hps `json:"hps"`
}

type Hps struct {
	Msg Msg `json:"msg"`
	Ext Ext `json:"ext"`
}
type Msg struct {
	Type   int    `json:"type"`
	Body   Body   `json:"body"`
	Action Action `json:"action"`
}
type Body struct {
	Content string `json:"content"`
	Title   string `json:"title"`
}
type Action struct {
	Type  int   `json:"type"`
	Param Param `json:"param"`
}
type Param struct {
	Intent string `json:"intent"`
	AppPkgName string `json:"appPkgName"`
	Url string `json:"url"`
}

type ExtObj struct {
	Name string
}
type Ext struct {
	BadgeAddNum string `json:"badgeAddNum"`	//设置应用角标数值，取值范围1-99
	BadgeClass string `json:"badgeClass"`		//桌面图标对应的应用入口Activity类  例如“com.test.badge.MainActivity”
	BigTag string `json:"biTag"`				//设置消息标签，如果带了这个标签，会在回执中推送给CP用于检测某种类型消息的到达率和状态。
	Customize []map[string]string `json:"customize"`		//"hps"->"ext"->"customize"下的内容为用户自定义扩展信息，开发者可以通过该消息实现onEvent点击事件的触发
}

/**
 **************************************** 封装
 */

func (this *Message) SetContent(content string) *Message {
	this.Hps.Msg.Body.Content = content
	return this
}

func (this *Message) SetTitle(title string) *Message {
	this.Hps.Msg.Body.Title = title
	return this
}

func (this *Message) SetIntent(intent string) *Message {
	this.Hps.Msg.Action.Param.Intent = intent
	return this
}

func (this *Message) SetAppPkgName(appPkgName string) *Message {
	this.Hps.Msg.Action.Param.AppPkgName = appPkgName
	return this
}

func (this *Message) SetActionType(actionType int) *Message {
	this.Hps.Msg.Action.Type = actionType
	return this
}

func (this *Message) SetExtCustmoize(exts map[string]string) *Message {
	huaweiCustomize := []map[string]string{}

	for key, value := range exts {
		m := map[string]string{}
		m[key] = value
		huaweiCustomize = append(huaweiCustomize, m)
	}

	this.Hps.Ext.Customize = huaweiCustomize
	return this
}

func (this *Message) Json() string {
	bytes, err := json.Marshal(this)
	if err != nil {
		fmt.Println(os.Stderr, err.Error())
		return ""
	}
	return string(bytes)
}
