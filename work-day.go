package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/jinzhu/now"
)

type Day struct {
	DayStr  string
	Weekday time.Weekday
	Workday bool
}

var DayFilename string
var DayMap = make(map[string]bool)

var today time.Time
var workDay bool

func InitMap() {
	filename := fmt.Sprintf("/opt/ban/data/%d.special", time.Now().Year())
	if len(DayMap) > 0 && DayFilename == filename {
		return
	}
	if _, err := os.Stat(filename); err != nil {
		log.Printf("year special file no exist err: %v\n", err)
		return
	}
	DayFilename = filename

	content, err := ioutil.ReadFile(DayFilename)
	if err != nil {
		log.Printf("read year special file err: %v\n", err)
		return
	}
	var days []Day
	if err := json.Unmarshal(content, &days); err != nil {
		log.Printf("json unmarshal year special file err: %v\n", err)
		return
	}
	for _, d := range days {
		DayMap[d.DayStr] = d.Workday
	}
	return
}

func refreshToday() {
	if now.BeginningOfDay().Equal(today) {
		return
	}
	today = now.BeginningOfDay()
	// 周一~周五工作日，周六周日休息日
	switch time.Now().Weekday() {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		workDay = true
	default:
		workDay = false
	}
	tempWorkday, exists := DayMap[today.Format("2006-01-02")]
	if exists {
		workDay = tempWorkday
	}
}

func WorkDay() bool {
	InitMap()
	refreshToday()
	return workDay
}
