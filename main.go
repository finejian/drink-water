package main

import (
	"time"

	"github.com/jinzhu/now"
	"github.com/robfig/cron"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	Mode995      = "995"
	ModeOddWeek  = "odd"  // 单周大
	ModeEvenWeek = "even" // 双周大

)

// 不支持热切换，重启后生效
var Mode = kingpin.Flag("mode", "run mode of dw. 995(default)/odd/even").Default(Mode995).Short('m').String()

const (
	perHour    = "0 30 9-21 * * *"
	perDay0840 = "0 40 8 * * *"
	perDay0900 = "0 0 9 * * *"
	perDay1200 = "0 59 11 * * *"
	perDay1800 = "0 59 17 * * *"
	perDay1900 = "0 59 18 * * *"
	perDay2100 = "0 59 20 * * *"
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
	case ModeOddWeek:
		if err := RunOddWeek(); err != nil {
			panic("run odd week err: " + err.Error())
		}
	case ModeEvenWeek:
		if err := RunEvenWeek(); err != nil {
			panic("run even week err: " + err.Error())
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
	if err := AddCron(perDay0840, WorkDay995, PostOn); err != nil {
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
单周工作日：
9:00 上班打卡
12:00 午饭
19:00 下班打卡
每1小时提醒喝水1次
每逢单周时周六上班
*/
func RunOddWeek() error {
	if err := AddCron(perDay1200, WorkDayOddWeek, PostEat); err != nil {
		return err
	}
	if err := AddCron(perDay0900, WorkDayOddWeek, PostOn); err != nil {
		return err
	}
	if err := AddCron(perDay1900, WorkDayOddWeek, PostOff); err != nil {
		return err
	}
	if err := AddCron(perHour, WorkDayOddWeek, PostDrink); err != nil {
		return err
	}
	return nil
}

/*
双周工作日：
9:00 上班打卡
12:00 午饭
19:00 下班打卡
每1小时提醒喝水1次
每逢双周时周六上班
*/
func RunEvenWeek() error {
	if err := AddCron(perDay1200, WorkDayEvenWeek, PostEat); err != nil {
		return err
	}
	if err := AddCron(perDay0900, WorkDayEvenWeek, PostOn); err != nil {
		return err
	}
	if err := AddCron(perDay1900, WorkDayEvenWeek, PostOff); err != nil {
		return err
	}
	if err := AddCron(perHour, WorkDayEvenWeek, PostDrink); err != nil {
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
	case ModeOddWeek, ModeEvenWeek:
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
