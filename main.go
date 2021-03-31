package main

import (
	"time"

	"github.com/jinzhu/now"
	"github.com/robfig/cron"
)

const (
	perHour    = "0 29 9-19 * * *"
	perDay0855 = "0 55 8 * * *"
	perDay1200 = "0 59 11 * * *"
	perDay1901 = "0 1 19 * * *"
)

var c *cron.Cron

func main() {
	c = cron.New()
	if err := Run(); err != nil {
		panic("run err: " + err.Error())
	}
	c.Start()
	select {}
}

/*
大小周工作日：
9:00 上班打卡
12:00 午饭
19:00 下班打卡
每1小时提醒喝水1次
每逢单周时周六上班
*/
func Run() error {
	if err := AddCron(perDay1200, WorkDay, PostEat); err != nil {
		return err
	}
	if err := AddCron(perDay0855, WorkDay, PostOn); err != nil {
		return err
	}
	if err := AddCron(perDay1901, WorkDay, PostOff); err != nil {
		return err
	}
	if err := AddCron(perHour, WorkDay, PostDrink); err != nil {
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
	todayEnd := now.BeginningOfDay().Add(time.Hour * 19)
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

	PostDrink()
}
