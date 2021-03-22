package news

import (
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego/httplib"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const (
	// https://juejin.cn/tag/%E6%B5%8B%E8%AF%95
	testURL = `https://api.juejin.cn/recommend_api/v1/article/recommend_tag_feed`
)

/*
{
    "err_no": 0,
    "err_msg": "success",
    "data": [
        {
            "article_id": "6932857399506436104",
            "article_info": {
                "article_id": "6932857399506436104",
                "user_id": "2313028193225319",
                "category_id": "6809637769959178254",
                "tag_ids": [
                    6809640364677267469
                ],
                "visible_level": 0,
                "link_url": "",
                "cover_image": "",
                "is_gfw": 0,
                "title": "聊聊gost的DeltaCompare",
                "brief_content": "gost的提供了DeltaCompareFloat32、DeltaCompareFloat64方法用于对比两个float类型的差值是否小于等于指定的delta。",
                "is_english": 0,
                "is_original": 1,
                "user_index": 4.050490241716322,
                "original_type": 0,
                "original_author": "",
                "content": "",
                "ctime": "1614181737",
                "mtime": "1614181739",
                "rtime": "1614181739",
                "draft_id": "6932855870569054216",
                "view_count": 306,
                "collect_count": 0,
                "digg_count": 0,
                "comment_count": 0,
                "hot_index": 15,
                "is_hot": 0,
                "rank_index": 0.5369891,
                "status": 2,
                "verify_status": 1,
                "audit_status": 2,
                "mark_content": ""
            },
            "author_user_info": {
                "user_id": "2313028193225319",
                "user_name": "go4it",
                "company": "",
                "job_title": "",
                "avatar_large": "https://sf1-ttcdn-tos.pstatp.com/img/user-avatar/6619c4d0978378785a89c945d78266a4~300x300.image",
                "level": 4,
                "description": "",
                "followee_count": 5,
                "follower_count": 493,
                "post_article_count": 1233,
                "digg_article_count": 0,
                "got_digg_count": 837,
                "got_view_count": 552982,
                "post_shortmsg_count": 0,
                "digg_shortmsg_count": 0,
                "isfollowed": false,
                "favorable_author": 1,
                "power": 6373,
                "study_point": 0,
                "university": {
                    "university_id": "0",
                    "name": "",
                    "logo": ""
                },
                "major": {
                    "major_id": "0",
                    "parent_id": "0",
                    "name": ""
                },
                "student_status": 0,
                "select_event_count": 0,
                "select_online_course_count": 0,
                "identity": 0,
                "is_select_annual": false,
                "select_annual_rank": 0,
                "annual_list_type": 0
            },
            "category": {
                "category_id": "6809637769959178254",
                "category_name": "后端",
                "category_url": "backend",
                "rank": 1,
                "ctime": 1457483880,
                "mtime": 1432503193,
                "show_type": 3
            },
            "tags": [
                {
                    "id": 2546494,
                    "tag_id": "6809640364677267469",
                    "tag_name": "Go",
                    "color": "#64D7E3",
                    "icon": "https://lc-gold-cdn.xitu.io/1aae38ab22d12b654cfa.png",
                    "back_ground": "",
                    "show_navi": 0,
                    "ctime": 1432234497,
                    "mtime": 1614214734,
                    "id_type": 9,
                    "tag_alias": "",
                    "post_article_count": 6768,
                    "concern_user_count": 80159
                }
            ],
            "user_interact": {
                "id": 6932857399506436104,
                "omitempty": 2,
                "user_id": 0,
                "is_digg": false,
                "is_follow": false,
                "is_collect": false
            },
            "org": {
                "org_info": null,
                "org_user": null,
                "is_followed": false
            }
        }
    ],
    "cursor": "eyJ2IjoiNjkzMjg1NzM5OTUwNjQzNjEwNCIsImkiOjIwfQ==",
    "count": 2247,
    "has_more": true
}
*/

func OnceTest(offset int) (article *Article, err error) {
	param, _ := json.Marshal(map[string]interface{}{
		"id_type":   2,
		"sort_type": 300,
		"cursor":    "0",
		"tag_ids":   []string{"6809640427465998350"},
	})
	resp, err := httplib.Post(golangURL).
		Header("Content-Type", "application/json; charset=utf-8").
		Body(param).DoRequest()
	if err != nil {
		err = errors.Wrap(err, "send http get request error.")
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "get http response error.")
		return
	}

	body := string(bytes)
	idValue := gjson.Get(body, `data.#.article_info.article_id`).Array()
	titleValue := gjson.Get(body, `data.#.article_info.title`).Array()
	if len(idValue) == 0 || len(titleValue) == 0 {
		err = errors.Wrap(err, "empty error")
		return
	}

	article = new(Article)
	article.Title = titleValue[offset].Array()[0].String()
	article.Url = "https://juejin.cn/post/" + idValue[offset].Array()[0].String()

	picurlValue := gjson.Get(body, `data.#.tags.icon`).Array()
	if len(picurlValue) > 0 {
		article.PicURL = picurlValue[offset].Array()[0].String()
	}
	return
}
