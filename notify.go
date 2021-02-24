package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/tidwall/gjson"

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
	wxURL      = ``
	englishURL = ``
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

type MsgNews struct {
	MsgType string `json:"msgtype"`
	News    News   `json:"news"`
}

type News struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"wxURL"`
	PicURL      string `json:"picurl"`
}

func bodyNews(title, description, url, picurl string) interface{} {
	return MsgNews{
		MsgType: "news",
		News: News{
			Articles: []Article{
				{
					Title:       title,
					Description: description,
					Url:         url,
					PicURL:      picurl,
				},
			},
		},
	}
}

/*
每日一句
https://www.tianapi.com/apiview/174

{
  "code": 200,
  "msg": "success",
  "newslist": [
    {
      "id": 4048,
      "content": "You’re never really done for as long as you’ve got a good story and someone to tell it to.",
      "source": "新版每日一句",
      "note": "只要你还有个好故事，还有一个能够倾诉的人，你就永远不会完蛋。",
      "tts": "https://staticedu-wps.cache.iciba.com/audio/25645f512d763b20b976abd0673c6948.mp3",
      "imgurl": "https://staticedu-wps.cache.iciba.com/image/c182da719ae82ca62c6ef76be4e7776e.png",
      "date": "2021-02-24"
    }
  ]
}
*/
func postNews() (err error) {
	defer func() {
		if err != nil {
			_, _ = postWx(drink)
		}
	}()

	resp, err := httplib.Get(englishURL).DoRequest()
	if err != nil {
		return errors.Wrap(err, "send http get request error.")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "get http response error.")
	}
	if value := gjson.Get(string(bytes), `code`).Array(); len(value) < 0 || value[0].Int() != 200 {
		return errors.Wrap(err, "request error")
	}

	var title, description, audio, picurl string
	if value := gjson.Get(string(bytes), `newslist.note`).Array(); len(value) > 0 {
		title = value[0].String()
	}
	if value := gjson.Get(string(bytes), `newslist.content`).Array(); len(value) > 0 {
		description = value[0].String()
	}
	if value := gjson.Get(string(bytes), `newslist.tts`).Array(); len(value) > 0 {
		audio = value[0].String()
	}
	if value := gjson.Get(string(bytes), `newslist.imgurl`).Array(); len(value) > 0 {
		picurl = value[0].String()
	}
	content := bodyNews(title, description, audio, picurl)
	_, err = postWx(content)
	return
}
