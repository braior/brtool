package brtool

import (
	"bytes"
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

// SendMessage 通过钉钉机器人发送消息
func (d *DingTalkClient) SendMessage() (bool, error) {
	var message string
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
		return false, fmt.Errorf("connect dingtalk url(%s) failed: %s", d.RobotURL, err)
	}

	if response.Status != "200" {
		body, _ := ioutil.ReadAll(response.Body)
		return false, fmt.Errorf("connect dingtalk url(%s) failed: %s", d.RobotURL, string(body))
	}

	ioutil.ReadAll(response.Body)
	return true, nil
}
