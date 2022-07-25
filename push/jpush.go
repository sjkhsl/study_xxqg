package push

import (
	"encoding/base64"
	"errors"

	"github.com/imroc/req/v3"

	"github.com/huoxue1/study_xxqg/lib"
)

func PushMessage(title, content, message, pushID string) error {
	if !lib.GetConfig().JiGuangPush.Enable {
		return nil
	}

	c := req.C()
	response, err := c.R().SetBodyJsonMarshal(map[string]interface{}{
		"platform": "all",
		"audience": map[string][]string{
			"registration_id": {pushID},
		},
		"notification": map[string]interface{}{
			"alert": content,
		},
		"message": map[string]string{
			"msg_content": message,
		},
	}).SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(lib.GetConfig().JiGuangPush.AppKey+":"+lib.GetConfig().JiGuangPush.Secret))).Post("https://api.jpush.cn/v3/push")
	if err != nil {
		return err
	}
	if response.IsSuccess() {
		return nil
	}
	return errors.New("消息推送失败" + response.Response.Status)
}
