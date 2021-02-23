package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego/httplib"
	"github.com/pkg/errors"
)

const (
	// 提示语句
	eat   = "痴迷于工作的小伙伴们，吃饭啦！"
	drink = "喝喝水摸摸鱼！"
	on    = "上班打卡！"
	off   = "下班打卡！"

	// 企业微信机器人
	url = ``
)

func PostEat() {
	_, _ = post(url, body(eat))
}

func PostDrink() {
	_, _ = post(url, body(drink))
}

func PostOn() {
	_, _ = post(url, body(on))
}

func PostOff() {
	_, _ = post(url, body(off))
}

func body(text string) interface{} {
	return map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": text,
		},
	}
}

func post(url string, body interface{}) ([]byte, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "report body json marshal error.")
	}

	resp, err := httplib.Post(url).Body(bodyByte).DoRequest()
	if err != nil {
		return nil, errors.Wrap(err, "send http post request error.")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "get http response error.")
	}
	return bytes, nil
}
