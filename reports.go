package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	clrtxt "github.com/daviddengcn/go-colortext"
	"github.com/eiannone/keyboard"
)

func reports() {
	fmt.Println("1 - В разрезе месяцев и статей расходов")
	fmt.Println("2 - Диаграмма")
	fmt.Println("3 - По дням")

	ch, _, err := keyboard.GetSingleKey()
	if err != nil {
		return
	}

	switch ch {
	case '1':
		editReportOptions(1, "расходы по месяцам и статьям расходов", reportMonthAndType)
	case '2':
		editReportOptions(2, "диаграмма по статьям расходов", reportDiagram)
	case '3':
		editReportOptions(3, "расходы по дням", reportDaily)
	}
}

func showHeader(descr string, dt1, dt2 time.Time, withKeys bool) {
	clearScreen()
	clrtxt.ChangeColor(clrtxt.White, true, clrtxt.Blue, false)
	txt := fmt.Sprintf("======== Отчет: %s ", descr)
	nrep := 57 - len([]rune(txt))
	if nrep <= 0 {
		nrep = 1
	}
	txt += strings.Repeat("=", nrep)
	fmt.Printf("%s\n", txt)

	if withKeys {
		fmt.Println("   Начальная дата: ", dateFmtExt(dt1),
			"   '1' - изменить      ")
		fmt.Println("    Конечная дата: ", dateFmtExt(dt2),
			"   '2' - изменить      ")
	} else {
		fmt.Println("За период", dateFmtExt(dt1), " - ", dateFmtExt(dt2), strings.Repeat(" ", 15))
	}
	clrtxt.ResetColor()
}

func ifEmptyReport() {
	fmt.Printf("Нет расходов за указанный период\n")
	fmt.Println("Нажмите любую клавишу для выхода")
	keyboard.GetSingleKey()
}

func waitExitReport() {
	fmt.Println()
	fmt.Println()
	fmt.Println("Нажмите любую клавишу для выхода")
	keyboard.GetSingleKey()
}

func editReportOptions(nRep int, descr string, showReport func(string, time.Time, time.Time)) {
	dt1 := beginOfMonth(stat.dt.AddDate(0, -6, 0))
	if nRep == 3 {
		dt1 = stat.dt.AddDate(0, 0, -10)
	}
	dt2 := endOfDay(stat.dt)
	for {
		showHeader(descr, dt1, dt2, true)
		fmt.Println("   Enter, '3' - Сформировать отчет")
		fmt.Println(" любая другая клавиша (ESC, Пробел, 0) - возврат в главное меню")

		ch, keyCode, err := keyboard.GetSingleKey()
		if err != nil {
			return
		}

		if ch == '3' || keyCode == keyboard.KeyEnter {
			showReport(descr, dt1, dt2)
			waitExitReport()
			return
		} else if ch == '1' {
			dt1 = inputDate(dt1)
		} else if ch == '2' {
			dt2 = inputDate(dt2)
		} else {
			return
		}
	}
}

func reportMonthAndType(descr string, dt1, dt2 time.Time) {
	showHeader(descr, dt1, dt2, false)

	dt1 = endOfDay(dt1.AddDate(0, 0, -1))
	dt2 = beginOfDay(dt2.AddDate(0, 0, 1))

	dat := make(map[rune]map[time.Time]int32)
	months := make(map[time.Time]int32)

	for _, v := range outgo {
		if v.dt.After(dt1) && v.dt.Before(dt2) {
			dt := beginOfMonth(v.dt)
			_, ok := dat[v.code]
			if !ok {
				mp := make(map[time.Time]int32)
				dat[v.code] = mp
			}
			dat[v.code][dt] += v.sum
			months[dt] += v.sum
		}
	}

	if len(months) == 0 {
		ifEmptyReport()
		return
	}

	aMoths := []time.Time{}
	for k := range months {
		aMoths = append(aMoths, k)
	}
	sort.Slice(aMoths, func(i, j int) bool {
		return aMoths[i].Before(aMoths[j])
	})

	fmt.Println()
	fmt.Println()

	s := "  Статьи расходов   "
	for _, mon := range aMoths {
		s += mon.Format("    2006-01")
	}
	clrtxt.Foreground(clrtxt.Green, true)
	fmt.Println(s)
	clrtxt.ResetColor()
	for _, v := range ogTypeClassifier {
		if mon, ok := dat[v.code]; ok {
			s := fmt.Sprintf("%-20s", v.name)
			for _, m := range aMoths {
				sum := mon[m]
				s += fmt.Sprintf("%11d", sum)
			}
			fmt.Println(s)
		}
	}
	fmt.Println()
	s = strings.Repeat(" ", 20)
	for _, mon := range aMoths {
		sum := months[mon]
		s += fmt.Sprintf("%11d", sum)
	}
	fmt.Println(s)
}

func reportDiagram(descr string, dt1, dt2 time.Time) {
	showHeader(descr, dt1, dt2, false)
	dt1 = endOfDay(dt1.AddDate(0, 0, -1))
	dt2 = beginOfDay(dt2.AddDate(0, 0, 1))

	dat := make(map[rune]int32)

	for _, v := range outgo {
		if v.dt.After(dt1) && v.dt.Before(dt2) {
			dat[v.code] += v.sum
		}
	}

	var max int32 = 0
	for _, v := range dat {
		if v > max {
			max = v
		}
	}

	if max == 0 {
		ifEmptyReport()
		return
	}

	fmt.Println()
	fmt.Println()

	sumTotal := 0

	for _, v := range ogTypeClassifier {
		if sum, ok := dat[v.code]; ok {
			sz := int(45.0 * float32(sum) / float32(max))
			clrtxt.ChangeColor(clrtxt.Green, true, clrtxt.Green, false)
			fmt.Print(strings.Repeat("|", sz))
			clrtxt.ResetColor()
			s := "    " + getOutgo(v.code) + " " + itoa(int(sum))
			fmt.Println(s)
			sumTotal += int(sum)
		}
	}
	fmt.Println("Итого " + fmt.Sprintf("%11d", sumTotal))
}

func reportDaily(descr string, dt1, dt2 time.Time) {

	if compareDates(dt2, dt1) > 40 {
		sayAndWait("Выберите период не дольше 40 дней")
		return
	}

	showHeader(descr, dt1, dt2, false)
	dt1 = endOfDay(dt1.AddDate(0, 0, -1))
	dt2 = beginOfDay(dt2.AddDate(0, 0, 1))

	dailySum := make(map[time.Time]int32)
	dailyDet := make(map[time.Time]string)

	for _, v := range outgo {
		if v.dt.After(dt1) && v.dt.Before(dt2) {
			dt := beginOfDay(v.dt)
			dailySum[dt] += v.sum
			dailyDet[dt] += "    " + getOutgo(v.code) + " " + itoa(int(v.sum))
		}
	}

	sumTotal := 0

	dt := beginOfDay(dt1)
	for compareDates(dt, dt2) < 0 {
		sum, ok1 := dailySum[dt]
		str, ok2 := dailyDet[dt]
		if ok1 && ok2 {
			fmt.Println(dateFmtExt(dt) +
				fmt.Sprintf("%11d", sum) + "   " +
				str)
			sumTotal += int(sum)
		}
		dt = dt.AddDate(0, 0, 1)
	}
	fmt.Println("Итого        " + fmt.Sprintf("%11d", sumTotal))
}
