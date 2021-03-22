package news

import (
	"io/ioutil"
	"strings"

	"github.com/astaxie/beego/httplib"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const (
	// https://www.zhihu.com/hot
	zhihuURL = `https://www.zhihu.com/api/v3/feed/topstory/hot-lists/total?limit=50&desktop=true`
)

/*
{
    "data": [
        {
            "type": "hot_list_feed",
            "style_type": "1",
            "id": "0_1614219077.2096949",
            "card_id": "Q_444936831",
            "target": {
                "id": 444936831,
                "title": "看《你好，李焕英》时孩子哭得稀里哗啦，但是回家后也没多刷一个碗，说两句，还一样顶嘴，为什么呢？",
                "url": "https://api.zhihu.com/questions/444936831",
                "type": "question",
                "created": 1613572645,
                "answer_count": 1936,
                "follower_count": 4716,
                "author": {
                    "type": "people",
                    "user_type": "people",
                    "id": "0",
                    "url_token": "",
                    "url": "",
                    "name": "用户",
                    "headline": "",
                    "avatar_url": "https://pica.zhimg.com/aadd7b895_s.jpg"
                },
                "bound_topic_ids": [
                    988,
                    5268,
                    65149,
                    99533,
                    572235
                ],
                "comment_count": 142,
                "is_following": false,
                "excerpt": ""
            },
            "attached_info": "CkAIi5HDzJDW8tcaEAMaCDYxMTUzNDQ4IKXUtIEGMI4BOOwkQAByCTQ0NDkzNjgzMXgAqgEJYmlsbGJvYXJk0gEA",
            "detail_text": "1233 万热度",
            "trend": 0,
            "debut": false,
            "children": [
                {
                    "type": "answer",
                    "thumbnail": "https://pic1.zhimg.com/80/v2-ba441e5346c5ecb832740380912339aa_720w.jpg?source=1940ef5c"
                }
            ]
        },
    ],
    "paging": {
        "is_end": true,
        "next": "",
        "previous": ""
    },
    "fresh_text": "热榜已更新"
}
*/

func OnceZhuhu(offset int) (article *Article, err error) {
	resp, err := httplib.Get(zhihuURL).DoRequest()
	if err != nil {
		err = errors.Wrap(err, "send http get request error.")
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "get http response error.")
		return
	}

	titleValue := gjson.Get(string(bytes), `data.#.target.title`).Array()
	urlValue := gjson.Get(string(bytes), `data.#.target.url`).Array()
	if len(titleValue) == 0 || len(urlValue) == 0 {
		err = errors.Wrap(err, "empty error")
		return
	}
	article = new(Article)
	article.Title = titleValue[offset].Array()[0].String()
	article.Url = urlValue[offset].Array()[0].String()
	article.Url = strings.ReplaceAll(article.Url, "/api.", "/www.")
	article.Url = strings.ReplaceAll(article.Url, "questions", "question")

	picurlValue := gjson.Get(string(bytes), `data.#.children.#.thumbnail`).Array()
	if len(picurlValue) > 0 && len(picurlValue[offset].Array()) > 0 {
		url := picurlValue[offset].Array()[0].String()
		arr := strings.Split(url, "?")
		if len(arr) > 0 {
			article.PicURL = arr[0]
		}
	}
	return
}
