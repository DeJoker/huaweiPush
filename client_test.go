package huaweiPush

import (
	"testing"
	"time"
)


var appId string = ""
var appSecret string = ""
var packageName string = ""

var regID string = ""



func TestViVoPush_Send(t *testing.T) {
	client := NewClient(appId, appSecret)

	message := NewMessage()

	message.SetTitle("发送标题")
	time.Now().String()
	message.SetContent("发送内容"+time.Now().Format("2006-01-02T15:04:04.333"))
	message.SetAppPkgName(packageName)
	message.SetBigTag(time.Now().String() + "-")



	//msg1 := NewVivoMessage("标题", "内容", "")
	result := client.PushMsgToList([]string{regID}, message.Json(), 1800)
	t.Logf("result=%+v\n", result)
	time.Sleep(10*time.Second)
	result = client.PushMsgToList([]string{regID}, message.Json(), 1800)
	t.Logf("result=%+v\n", result)
	time.Sleep(3*time.Second)
	result = client.PushMsgToList([]string{regID}, message.Json(), 1800)
	t.Logf("result=%+v\n", result)
}