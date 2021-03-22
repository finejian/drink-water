package news

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
	Url         string `json:"url"`
	PicURL      string `json:"picurl"`
}

func Once(offset int) (msg *MsgNews, err error) {
	var zhihu, golang, vue, test, other *Article
	if zhihu, err = OnceZhuhu(offset); err != nil {
		return
	}
	if golang, err = OnceGolang(offset); err != nil {
		return
	}
	if vue, err = OnceVue(offset); err != nil {
		return
	}
	if test, err = OnceTest(offset); err != nil {
		return
	}
	if other, err = OnceOther(offset); err != nil {
		return
	}
	msg = &MsgNews{
		MsgType: "news",
		News: News{
			Articles: []Article{
				*zhihu,
				*golang,
				*vue,
				*test,
				*other,
			},
		},
	}
	return
}
