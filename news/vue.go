package news

import (
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego/httplib"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

const (
	// https://juejin.cn/tag/Vue.js
	vueURL = `https://api.juejin.cn/recommend_api/v1/article/recommend_tag_feed`
)

/*
{
    "err_no": 0,
    "err_msg": "success",
    "data": [
        {
            "article_id": "6933005298198118413",
            "article_info": {
                "article_id": "6933005298198118413",
                "user_id": "2348212566363197",
                "category_id": "6809637767543259144",
                "tag_ids": [
                    6809640369764958215
                ],
                "visible_level": 0,
                "link_url": "",
                "cover_image": "",
                "is_gfw": 0,
                "title": "TS 加持的 Vue 3，如何帮你轻松构建企业级前端应用",
                "brief_content": "在如今被三大框架支配的前端领域，已经很少有人不知道 Vue 了。2014 年，前 Google 工程师尤雨溪发布了所谓的渐进式（Progressive）前端应用框架 Vue，其简化的模版绑定和组件化思想给当时还是 jQuery 时代的前端领域产生了积极而深远的影响。Vue 的诞…",
                "is_english": 0,
                "is_original": 1,
                "user_index": 0,
                "original_type": 0,
                "original_author": "",
                "content": "",
                "ctime": "1614216261",
                "mtime": "1614219643",
                "rtime": "1614219643",
                "draft_id": "6933005171744047117",
                "view_count": 209,
                "collect_count": 0,
                "digg_count": 2,
                "comment_count": 0,
                "hot_index": 10,
                "is_hot": 0,
                "rank_index": 4.2619649,
                "status": 2,
                "verify_status": 1,
                "audit_status": 2,
                "mark_content": ""
            },
            "author_user_info": {
                "user_id": "2348212566363197",
                "user_name": "MarvinZhang",
                "company": "",
                "job_title": "搬砖工程师",
                "avatar_large": "https://sf1-ttcdn-tos.pstatp.com/img/user-avatar/2f44a2b12aa41d57813f3e5687fd484c~300x300.image",
                "level": 3,
                "description": "前端+爬虫+数据分析",
                "followee_count": 355,
                "follower_count": 1809,
                "post_article_count": 60,
                "digg_article_count": 317,
                "got_digg_count": 463,
                "got_view_count": 50289,
                "post_shortmsg_count": 25,
                "digg_shortmsg_count": 50,
                "isfollowed": false,
                "favorable_author": 0,
                "power": 1581,
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
                "category_id": "6809637767543259144",
                "category_name": "前端",
                "category_url": "frontend",
                "rank": 2,
                "ctime": 1457483942,
                "mtime": 1432503190,
                "show_type": 3
            },
            "tags": [
                {
                    "id": 2546498,
                    "tag_id": "6809640369764958215",
                    "tag_name": "Vue.js",
                    "color": "#41B883",
                    "icon": "https://lc-gold-cdn.xitu.io/7b5c3eb591b671749fee.png",
                    "back_ground": "",
                    "show_navi": 0,
                    "ctime": 1432234520,
                    "mtime": 1614222251,
                    "id_type": 9,
                    "tag_alias": "",
                    "post_article_count": 23483,
                    "concern_user_count": 291060
                }
            ],
            "user_interact": {
                "id": 6933005298198118413,
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
        },
    ],
    "cursor": "eyJ2IjoiNjkzMzAwNTI5ODE5ODExODQxMyIsImkiOjIwfQ==",
    "count": 6482,
    "has_more": true
}
*/

func OnceVue(offset int) (article *Article, err error) {
	param, _ := json.Marshal(map[string]interface{}{
		"id_type":   2,
		"sort_type": 300,
		"cursor":    "0",
		"tag_ids":   []string{"6809640369764958215"},
	})
	resp, err := httplib.Post(vueURL).
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
