package huaweiPush

import (
	"testing"
	"time"
)


var appId string = ""
var appSecret string = ""
var packageName string = ""

var regID string = ""


func TestHuaweiPush_Send(t *testing.T) {
	client := NewClient(appId, appSecret)

	message := NewMessage()

	message.SetTitle("发送标题")

	//message.SetContent("发送内容"+time.Now().Format("2006-01-02T15:04:04.333"))
	message.SetAppPkgName(packageName)
	message.SetBigTag(time.Now().String() + "-")



	//msg1 := NewVivoMessage("标题", "内容", "")
	for i:=0; i<3; i++ {
		message.SetContent("发送内容"+time.Now().Format("2006-01-02T15:04:04.333"))
		result,err := client.PushMsgToList([]string{regID}, message.Json(), 1800)
		t.Logf("%+v\n", result)
		t.Log(err)
		time.Sleep(2*time.Second)
	}
}