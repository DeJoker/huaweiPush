package huaweiPush

import (
	"encoding/json"
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
}

type ExtObj struct {
	Name string
}
type Ext struct {
	Action  string `json:"action"`
	Func    string `json:"func"`
	Collect string `json:"collect"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Url     string `json:"url"`
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

func (this *Message) SetExtAction(Action string) *Message {
	this.Hps.Ext.Action = Action
	return this
}
func (this *Message) SetExtFunc(Func string) *Message {
	this.Hps.Ext.Func = Func
	return this
}
func (this *Message) SetExtCollect(collect string) *Message {
	this.Hps.Ext.Collect = collect
	return this
}
func (this *Message) SetExtTitle(title string) *Message {
	this.Hps.Ext.Title = title
	return this
}
func (this *Message) SetExtContent(content string) *Message {
	this.Hps.Ext.Collect = content
	return this
}

func (this *Message) SetExtUrl(url string) *Message {
	this.Hps.Ext.Url = url
	return this
}

func (this *Message) Json() string {
	bytes, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
