package main

import (
	"time"

	"github.com/jinzhu/now"
	"github.com/robfig/cron"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	Mode995  = "995"
	ModeWeek = "week" // 大小周
)

var Mode = kingpin.Flag("mode", "run mode of dw. week(default)/995").Default(ModeWeek).Short('m').String()
var FirstBigWeekDay = kingpin.Flag("bigWeekDay", "first big week, before today, eg: 2021-03-01").Default("2021-03-01").Short('b').String()

const (
	perHour    = "0 29 9-19 * * *"
	perDay0855 = "0 55 8 * * *"
	perDay1200 = "0 59 11 * * *"
	perDay1800 = "0 0 18 * * *"
	perDay1901 = "0 1 19 * * *"
	perDay2100 = "0 0 21 * * *"
)

var c *cron.Cron

func main() {
	kingpin.Parse()
	c = cron.New()
	switch *Mode {
	case Mode995:
		if err := Run995(); err != nil {
			panic("run 995 err: " + err.Error())
		}
	case ModeWeek:
		if err := RunWeek(); err != nil {
			panic("run week err: " + err.Error())
		}
	default:
		panic("unknown mode")
	}
	c.Start()
	select {}
}

/*
995工作日：
8:40 上班打卡
12:00 午饭
18:00 晚饭
21:00 下班打卡
每1小时提醒喝水1次
*/
func Run995() error {
	if err := AddCron(perDay1200, WorkDay995, PostEat); err != nil {
		return err
	}
	if err := AddCron(perDay1800, WorkDay995, PostEat); err != nil {
		return err
	}
	if err := AddCron(perDay0855, WorkDay995, PostOn); err != nil {
		return err
	}
	if err := AddCron(perDay2100, WorkDay995, PostOff); err != nil {
		return err
	}
	if err := AddCron(perHour, WorkDay995, PostDrink); err != nil {
		return err
	}
	return nil
}

/*
大小周工作日：
9:00 上班打卡
12:00 午饭
19:00 下班打卡
每1小时提醒喝水1次
每逢单周时周六上班
*/
func RunWeek() error {
	if err := AddCron(perDay1200, WorkDayWeek, PostEat); err != nil {
		return err
	}
	if err := AddCron(perDay0855, WorkDayWeek, PostOn); err != nil {
		return err
	}
	if err := AddCron(perDay1901, WorkDayWeek, PostOff); err != nil {
		return err
	}
	if err := AddCron(perHour, WorkDayWeek, PostDrink); err != nil {
		return err
	}
	return nil
}

func AddCron(spec string, workDay func() bool, post func()) error {
	return c.AddFunc(spec, func() {
		if !workDay() {
			return
		}
		if spec == perHour {
			PostPerHour()
		} else {
			post()
		}
	})
}

func PostPerHour() {
	todayNow := time.Now()
	// 上班
	todayStart := now.BeginningOfDay().Add(time.Hour * 9)

	// 下班
	var todayEnd time.Time
	switch *Mode {
	case Mode995:
		todayEnd = now.BeginningOfDay().Add(time.Hour * 21)
	case ModeWeek:
		todayEnd = now.BeginningOfDay().Add(time.Hour * 19)
	}
	// 上班前，下班后提醒
	if todayNow.Before(todayStart) || todayNow.After(todayEnd) {
		return
	}

	// 午餐不提醒
	todayLunchStart := now.BeginningOfDay().Add(time.Hour * 12)
	todayLunchEnd := now.BeginningOfDay().Add(time.Hour * 14)
	if todayNow.Before(todayLunchEnd) && todayNow.After(todayLunchStart) {
		return
	}

	// 晚餐不提醒
	todayDinnerStart := now.BeginningOfDay().Add(time.Hour * 18)
	todayDinnerEnd := now.BeginningOfDay().Add(time.Hour * 19)
	if *Mode == Mode995 && todayNow.Before(todayDinnerEnd) && todayNow.After(todayDinnerStart) {
		return
	}

	PostDrink()
}
