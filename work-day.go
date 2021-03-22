package main

import (
	"io/ioutil"
	"time"

	"github.com/tidwall/gjson"

	"github.com/astaxie/beego/httplib"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

var today time.Time
var workDay bool

const holidayURL = `http://timor.tech/api/holiday/info`

func WorkDay995() bool {
	refreshToday()
	return workDay
}

func WorkDayWeek() bool {
	refreshToday()

	now.WeekStartDay = time.Monday
	// 第一个大周的第一天
	firstBigWeekDay, _ := time.ParseInLocation("2006-01-02", *FirstBigWeekDay, time.Now().Location())
	firstBigWeekDay = now.New(firstBigWeekDay).BeginningOfWeek()

	today = time.Now()
	// 距离第一个大周的第一天的天数
	days := int(today.Sub(firstBigWeekDay).Hours() / 24)
	isSmall := days/7%2 == 1
	// 当前是小周，且今天是周六
	if isSmall && today.Weekday() == time.Saturday {
		workDay = true
	}
	return workDay
}

func refreshToday() {
	if now.BeginningOfDay().Equal(today) {
		return
	}

	temp, err := isWorkDay(holidayURL)
	if err == nil {
		today = now.BeginningOfDay()
		workDay = temp
		return
	}
	// 查询失败时，周一~周五工作日，周六周日休息日
	switch time.Now().Weekday() {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		workDay = true
	default:
		workDay = false
	}
	today = now.BeginningOfDay()
}

// 在线查询今天是否是工作日
// api地址： http://timor.tech/api/holiday/
func isWorkDay(url string) (bool, error) {
	resp, err := httplib.Get(url).DoRequest()
	if err != nil {
		return false, errors.Wrap(err, "send http get request error.")
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "get http response error.")
	}

	// "type": enum(0, 1, 2, 3), // 节假日类型，分别表示 工作日、周末、节日、调休。
	var workDay bool
	if value := gjson.Get(string(bytes), `type.type`).Array(); len(value) > 0 {
		workDay = value[0].Int() == 0 || value[0].Int() == 3
	}
	return workDay, nil
}
