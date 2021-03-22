package main

import (
	"encoding/json"
	"finejian/drink-water/news"
	"io/ioutil"
	"time"

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
	wxURL = ``
)

func PostEat() {
	_, _ = postWx(bodyText(eat))
}

func PostDrink() {
	_ = postNews()
}

func PostOn() {
	_, _ = postWx(bodyText(on))
}

func PostOff() {
	_, _ = postWx(bodyText(off))
}

func bodyText(text string) interface{} {
	return map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content": text,
		},
	}
}

func postWx(body interface{}) ([]byte, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(err, "report body json marshal error.")
	}

	resp, err := httplib.Post(wxURL).Body(bodyByte).DoRequest()
	if err != nil {
		return nil, errors.Wrap(err, "send http post request error.")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "get http response error.")
	}
	return bytes, nil
}

func postNews() (err error) {
	defer func() {
		if err != nil {
			_, _ = postWx(drink)
		}
	}()

	offset := time.Now().Hour() - 9
	msg, err := news.Once(offset)
	if err != nil {
		err = errors.New("get news error")
		return
	}
	_, err = postWx(msg)
	return
}
