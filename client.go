package huaweiPush

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/**
 * init
 */
func NewClient(clientID string, clientSecret string) *HuaweiPushClient {

	vers := &Vers{
		Ver:   "1",
		AppID: clientID,
	}
	nspCtx, _ := json.Marshal(vers)
	return &HuaweiPushClient{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		NspCtx:       string(nspCtx),
	}
}

/**
 * message init
 */
func NewMessage() *Message {
	return &Message{
		Hps: Hps{
			Msg: Msg{
				Type: 3, //1, 透传异步消息; 3, 系统通知栏异步消息;
				Body: Body{
					Content: "",
					Title:   "",
				},
				Action: Action{
					Type: 1, //1, 自定义行为; 2, 打开URL; 3, 打开App;
					Param: Param{
						Intent: "#Intent;compo=com.rvr/.Activity;S.W=U;end",
						AppPkgName:"",
					},
				},
			},
			Ext: Ext{ // 扩展信息, 含 BI 消息统计, 特定展示风格, 消息折叠;
				Action:  "",
				Func:    "",
				Collect: "",
				Title:   "",
				Content: "",
				Url:     "",
			},
		},
	}
}

/**
 * form post
 */
func FormPost(url string, data url.Values) ([]byte, error) {
	u := ioutil.NopCloser(strings.NewReader(data.Encode()))
	r, err := http.Post(url, "application/x-www-form-urlencoded", u)
	if err != nil {

		return []byte(""), err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {

		return []byte(""), err
	}
	return b, err
}

func doPost(ctx context.Context, url string, form url.Values) ([]byte, error) {
	var result []byte
	var req *http.Request
	var res *http.Response
	var err error
	req, err = http.NewRequest("POST", url, strings.NewReader(form.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	client := &http.Client{}
	tryTime := 0

	for tryTime < PostRetryTimes {
		res, err = ctxhttp.Do(ctx, client, req)
		if err != nil {
			fmt.Println("huawei push post err:", err, tryTime)
			select {
			case <-ctx.Done():
				return nil, err
			default:
			}
			tryTime += 1
			if tryTime < PostRetryTimes {
				continue
			}
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, errors.New("network error")
		}
		result, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	return []byte("unknow error"), nil
}

/**
 * get token
 */
func (this HuaweiPushClient) GetToken() string {
	reqUrl := TOKEN_URL
	param := make(url.Values)
	param["grant_type"] = []string{GRANTTYPE}
	param["client_id"] = []string{this.ClientId}
	param["client_secret"] = []string{this.ClientSecret}
	res, err := FormPost(reqUrl, param)

	if nil != err {
		return ""
	}
	var tokenRes = &TokenResStruct{}
	err = json.Unmarshal(res, tokenRes)
	if err != nil {
		return ""
	}
	return tokenRes.Access_token
}

/**
 * push msg
 */
func (this HuaweiPushClient) PushMsg(deviceToken, payload string) string {

	accessToken := this.GetToken()
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(this.NspCtx)


	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "[\"" + deviceToken + "\"]",
		"payload":           payload,
		"expire_time":       time.Now().Format("2006-01-02T15:04"),
	}


	param := make(url.Values)
	param["access_token"] = []string{originParam["access_token"]}
	param["nsp_svc"] = []string{originParam["nsp_svc"]}
	param["nsp_ts"] = []string{originParam["nsp_ts"]}
	param["device_token_list"] = []string{originParam["device_token_list"]}
	param["payload"] = []string{originParam["payload"]}

	// push
	res, _ := FormPost(reqUrl, param)

	return string(res)
}

func (this HuaweiPushClient) PushMsgToList(deviceTokens []string, payload string) string {

	accessToken := this.GetToken()
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(this.NspCtx)


	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "",
		"payload":           payload,
		"expire_time":       time.Now().Format("2006-01-02T15:04"),
	}

	jdeviceTokenArray, jsonErr := json.Marshal(deviceTokens)
	if jsonErr != nil {
		jsonErr.Error()
	}
	originParam["device_token_list"] = string(jdeviceTokenArray)


	param := make(url.Values)
	param["access_token"] = []string{originParam["access_token"]}
	param["nsp_svc"] = []string{originParam["nsp_svc"]}
	param["nsp_ts"] = []string{originParam["nsp_ts"]}
	param["device_token_list"] = []string{originParam["device_token_list"]}
	param["payload"] = []string{originParam["payload"]}

	// push
	res, _ := doPost(context.Background(), reqUrl, param)

	return string(res)
}
