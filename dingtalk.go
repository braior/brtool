package brtool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DingTalkMessage 钉钉消息结构体
type DingTalkMessage struct {
	Message string // 消息实体
	Title   string // markdown标题
	Type    string // 消息类型
}

// DingTalkClient 钉钉机器人client
type DingTalkClient struct {
	RobotURL string
	MsgInfo  *DingTalkMessage
}

// DingTalkResponse 。。。
type DingTalkResponse struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// SendMessage 通过钉钉机器人发送消息
func (d *DingTalkClient) SendMessage() (bool, error) {
	var message string
	dingTalkRresponse := new(DingTalkResponse)

	switch d.MsgInfo.Type {
	case "text":
		message = fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, d.MsgInfo.Message)
	case "markdown":
		message = fmt.Sprintf(`{"msgtype": "markdown", "markdown": {"title": "%s", "text": "%s"}}`, d.MsgInfo.Title, d.MsgInfo.Message)
	default:
		message = fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, d.MsgInfo.Message)
	}

	client := &http.Client{}
	request, _ := http.NewRequest("POST", d.RobotURL, bytes.NewBuffer([]byte(message)))
	request.Header.Set("Content-type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return false, fmt.Errorf("connect dingtalk url failed: %s", err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	err = json.Unmarshal(body, &dingTalkRresponse)
	if err != nil {
		return false, fmt.Errorf("parse json data error: %s", err)
	}

	if dingTalkRresponse.ErrCode != int64(0) {
		return false, fmt.Errorf("send message failed: %s", string(body))
	}

	if response.StatusCode != 200 {
		return false, fmt.Errorf("dingtalk response satus code is not 200 but is %d, respone body is: %s", response.StatusCode, string(body))
	}
	return true, nil
}
