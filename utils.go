package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/eiannone/keyboard"
)

var itoa = strconv.Itoa

func (o statusInf) drop() {
	o.dt = time.Now()
	o.outgoDay = 0
	o.outgoMonth = 0
	o.outgoPrevMonth = 0
	o.outgoYear = 0
}

func askAndGetYes() bool {
	ch, _, _ := keyboard.GetSingleKey()
	return ch == 'y' || ch == 'Y' || ch == '1'
}

func sayAndWait(txt string) {
	fmt.Println(txt)
	keyboard.GetSingleKey()
}

func (o outgoItem) copy() outgoItem {
	return outgoItem{sum: o.sum, code: o.code, dt: o.dt}
}

func (o outgoItem) toString() string {
	s := fmt.Sprintf("%7d %-22s %s", o.sum, getOutgo(o.code), o.dt.Format("02.01.2006"))
	return s
}

func (o outgoStr) toString() string {
	s := fmt.Sprintf("%7d %-22s", o.sumToday, o.name)
	return s
}
func (o outgoStr) toStringM() string {
	s := fmt.Sprintf("%7d %-22s %10d", o.sumToday, o.name, o.sumMonth)
	return s
}

func getOutgo(code rune) string {
	name, _ := ogClassifier[code]
	return name
}

func compareDates(a, b time.Time) int {
	aa := a.Year() + a.YearDay()
	bb := b.Year() + b.YearDay()
	return aa - bb
}

func isCurDate(dt time.Time) bool {
	return compareDates(dt, stat.dt) == 0
}

func getNMonth(dt time.Time) int {
	return (dt.Year() << 4) + int(dt.Month())
}

func isCurMonth(dt time.Time) bool {
	return getNMonth(dt) == getNMonth(stat.dt)
}

func isPrevMonth(dt time.Time) bool {
	return getNMonth(dt) == getNMonth(stat.dt.AddDate(0, -1, 0))
}

func isCurYear(dt time.Time) bool {
	return dt.Year() == stat.dt.Year()
}

func clearScreen() {
	//fmt.Printf("\n\n\n\n\n================================================")
	//fmt.Println("================================================")
	//fmt.Printf("================================================\n\n\n\n\n")

	switch runtime.GOOS {
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func weekDatStr(dt time.Time) string {
	wd := dt.Weekday()
	switch wd {
	case time.Monday:
		return "Пн"
	case time.Tuesday:
		return "Вт"
	case time.Wednesday:
		return "Ср"
	case time.Thursday:
		return "Чт"
	case time.Friday:
		return "Пт"
	case time.Saturday:
		return "Сб"
	default:
		return "Вс"
	}
}

func dateFmtExt(dt time.Time) string {
	return dt.Format("02.01.2006") + " " + weekDatStr(dt)
}

func beginOfDay(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0, time.UTC)
}

func beginOfMonth(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func beginOfYear(dt time.Time) time.Time {
	return time.Date(dt.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func endOfDay(dt time.Time) time.Time {
	return beginOfDay(dt.AddDate(0, 0, 1)).Add(-1)
}

func endOfMonth(dt time.Time) time.Time {
	return beginOfMonth(dt.AddDate(0, 1, 0)).Add(-1)
}

func endOfYear(dt time.Time) time.Time {
	return beginOfYear(dt.AddDate(1, 0, 0)).Add(-1)
}

func getTimer() func() time.Duration {
	start := time.Now()
	return func() time.Duration {
		return time.Since(start)
	}
}
