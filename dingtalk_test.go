package brtool

import "testing"

// SendByDingTalkRobot 通过钉钉发送消息通知
func TestSendMessage(t *testing.T) {

	robotURL := "https://oapi.dingtalk.com/robot/send?access_token=xxxxxx"

	dingtalk := &DingTalkClient{
		RobotURL: robotURL,
		MsgInfo: &DingTalkMessage{
			Type:    "text",
			Message: "TEST123456",
			Title:   "",
		},
	}

	ok, err := dingtalk.SendMessage()
	if err != nil {
		dingFields := map[string]interface{}{
			"entryType":      "DingTalkRobot",
			"dingTalkRobot ": robotURL,
		}
		t.Errorf("send %s failed: %s", dingFields, err)
	}
	if !ok{
		t.Errorf("send failed: %s",err)
	}
}
