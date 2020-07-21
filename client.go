package huaweiPush

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
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
					Type: 3, //1, 自定义行为; 2, 打开URL; 3, 打开App;
					Param: Param{
						Intent:     "",
						AppPkgName: "",
						Url:        "",
					},
				},
			},
			Ext: Ext{ // 扩展信息, 含 BI 消息统计, 特定展示风格, 消息折叠;
				BigTag:    "",
				Customize: []map[string]string{},
			},
		},
	}
}


// http base

var (
	httpclient = &http.Client{
		Timeout : time.Second * 60,
	}
)


func FormPost(url string, data url.Values) ([]byte, error) {
	u := ioutil.NopCloser(strings.NewReader(data.Encode()))
	r, err := httpclient.Post(url, "application/x-www-form-urlencoded", u)
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
		return nil, errors.New(("create post request error"))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	tryTime := 0

	for tryTime < PostRetryTimes {
		res, err = ctxhttp.Do(ctx, httpclient, req)
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

		result, err = ioutil.ReadAll(res.Body)
		if res.StatusCode != http.StatusOK {
			return result, errors.New("network error, http code:"+strconv.Itoa(res.StatusCode))
		}
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

const (
	OneHour = 3600 * 1000
)

func (hc *HuaweiPushClient) GetToken() (string, error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if hc.AuthInfo != nil && hc.AuthInfo.AccessToken != "" && hc.AuthInfo.validTime>now {
		return hc.AuthInfo.AccessToken, nil
	}

	reqUrl := TOKEN_URL
	param := make(url.Values)
	param["grant_type"] = []string{GRANTTYPE}
	param["client_id"] = []string{hc.ClientId}
	param["client_secret"] = []string{hc.ClientSecret}
	res, err := FormPost(reqUrl, param)

	if nil != err {
		return "", err
	}
	var tokenRes = &TokenResult{}
	err = json.Unmarshal(res, tokenRes)
	if err != nil {
		return "", err
	}
	hc.AuthInfo = tokenRes
	hc.AuthInfo.validTime = now + int64(hc.AuthInfo.ExpiresIn*2/3) // 提前过期token时间到2/3
	return tokenRes.AccessToken, nil
}

/**
 * push msg with token push into
 */
func (hc *HuaweiPushClient) PushMsg(accessToken, deviceToken, payload string, timeToLive int) (string, error) {
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(hc.NspCtx)

	now := time.Now()
	expireSecond := time.Duration(timeToLive * 1e9)
	expireTime := now.Add(expireSecond)

	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "[\"" + deviceToken + "\"]",
		"payload":           payload,
		"expire_time":       expireTime.Format("2006-01-02T15:04"),
	}

	param := make(url.Values)
	param["access_token"] = []string{originParam["access_token"]}
	param["nsp_svc"] = []string{originParam["nsp_svc"]}
	param["nsp_ts"] = []string{originParam["nsp_ts"]}
	param["device_token_list"] = []string{originParam["device_token_list"]}
	param["payload"] = []string{originParam["payload"]}
	param["expire_time"] = []string{originParam["expire_time"]}

	// push
	res, err := FormPost(reqUrl, param)

	return string(res), err
}
func (hc *HuaweiPushClient) PushMsgToList(deviceTokens []string, payload string, timeToLive int) (string, error) {
	accessToken, err := hc.GetToken()
	if err != nil {
		return "", err
	}
	return hc.PushToListWithToken(accessToken, deviceTokens, payload, timeToLive)
}

func (hc *HuaweiPushClient) PushToListWithToken(accessToken string, deviceTokens []string, payload string, timeToLive int) (string, error) {
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(hc.NspCtx)

	now := time.Now()
	expireSecond := time.Duration(timeToLive * 1e9)
	expireTime := now.Add(expireSecond)

	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "",
		"payload":           payload,
		"expire_time":       expireTime.Format("2006-01-02T15:04"),
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
	param["expire_time"] = []string{originParam["expire_time"]}

	// push
	res, err := doPost(context.Background(), reqUrl, param)

	return string(res), err
}

func (hc *HuaweiPushClient) PushMsgToArrayNoExpire(deviceTokens []string, payload string) (string, error) {
	accessToken, err := hc.GetToken()
	if err != nil {
		return "", err
	}
	reqUrl := PUSH_URL + "?nsp_ctx=" + url.QueryEscape(hc.NspCtx)

	var originParam = map[string]string{
		"access_token":      accessToken,
		"nsp_svc":           NSP_SVC,
		"nsp_ts":            strconv.Itoa(int(time.Now().Unix())),
		"device_token_list": "",
		"payload":           payload,
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
	res, err := doPost(context.Background(), reqUrl, param)

	return string(res), err
}
