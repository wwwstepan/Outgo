package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"mylogs"

	clrtxt "github.com/daviddengcn/go-colortext"
	"github.com/eiannone/keyboard"
)

var mainMenu = []string{
	" 1  Установить дату",
	" 2  Добавить расход",
	" 3  Удалить расход",
	" 4  Отчеты",
	" 0  Выход",
}

var stat statusInf
var wasChanges bool

var outgo []outgoItem
var outgoToday []outgoItem
var ogToday []outgoStr
var ogClassifier map[rune]string
var ogClassifierInd map[rune]int32

var ogTypeClassifier []outgoType

func main() {

	mylogs.Init("outgo.log")

	loadClassifier()

	stat = statusInf{dt: time.Now()}
	wasChanges = false

	ogClassifier = make(map[rune]string)
	ogClassifierInd = make(map[rune]int32)

	for _, o := range ogTypeClassifier {
		ogClassifier[o.code] = o.name
		ogClassifierInd[o.code] = o.ind
	}

	outgo = make([]outgoItem, 0)
	outgoToday = make([]outgoItem, 0)

	if err := loadOutgo(); err != nil {
		fmt.Println("Ошибка загрузки базы данных: ", err)
		fmt.Println("Выходим (Y/N) / (1/0) ?")
		if askAndGetYes() {
			return
		}
	}

	ogToday = getSortTodayOG()

	for {
		showMainMenu()

		ch, keyCode, err := keyboard.GetSingleKey()
		if err != nil {
			continue
		}

		if ch == '0' || keyCode == keyboard.KeyEsc {
			if quitProgram() {
				return
			}
			continue
		}

		switch ch {
		case '1':
			setCurrentDate()
		case '2':
			addOutgo()
		case '3':
			deleteOutgo()
		case '4':
			reports()
		default:

		}
	}
}

func setCurrentDate() {
	oldDate := stat.dt
	dt := inputDate(stat.dt)
	if compareDates(dt, oldDate) != 0 {
		stat.clear()
		stat.dt = dt
		stat.reCalculate()
		ogToday = getSortTodayOG()
	}
}

func quitProgram() bool {
	fmt.Println("Выход (Y/N) / (1/0) ?")
	if !askAndGetYes() {
		return false
	}

	if wasChanges {
		fmt.Println("Сохранить изменения? (Y/N) / (1/0) ?")
		if askAndGetYes() {
			if err := saveOutgo(); err != nil {
				fmt.Println("Ошибка сохранения: ", err)
				fmt.Println("Все равно выходим (Y/N) / (1/0) ?")
				if !askAndGetYes() {
					return false
				}
			}
		}
	}

	clearScreen()
	return true

}

func showStatusLine() {
	clearScreen()
	clrtxt.ChangeColor(clrtxt.White, true, clrtxt.Blue, false)
	fmt.Printf("Дата %s     Траты за день %-11d Траты за месяц %-11d\n", dateFmtExt(stat.dt), stat.outgoDay, stat.outgoMonth)
	fmt.Printf("           Траты за предыдущий месяц %-11d Траты за год   %-11d", stat.outgoPrevMonth, stat.outgoYear)
	clrtxt.ResetColor()
}

func showMainMenu() {
	showStatusLine()
	clrtxt.Foreground(clrtxt.Green, true)
	fmt.Println("\n\n" + strings.Repeat(" ", 30) +
		"=== Расходы за " +
		stat.dt.Format("02.01.06") + " ======= С " + stat.dt.AddDate(0, -1, 0).Format("02.01.06"))
	clrtxt.ResetColor()
	n := 0
	for _, strMenu := range mainMenu {
		var s string

		clrtxt.Foreground(clrtxt.Yellow, false)
		fmt.Print(strMenu)
		clrtxt.ResetColor()

		if n < len(ogToday) {
			lenStrMenu := len([]rune(strMenu))
			s = strings.Repeat(" ", 30-lenStrMenu) +
				ogToday[n].toStringM()
			fmt.Println(s)
		} else {
			fmt.Println()
		}
		n++
	}
	for n < len(ogToday) {
		s := strings.Repeat(" ", 30) + ogToday[n].toStringM()
		fmt.Println(s)
		n++
	}
}

func showMenuOutgo() {
	for _, o := range ogTypeClassifier {
		fmt.Println(string(o.code), o.name)
	}
}

func addOutgo() {
	showStatusLine()
	fmt.Printf("\n\nУкажите сумму расхода (1..10млн). Допускается ввод нескольких\nсумм через '+', они будут суммированы.\nПробелы при вводе не допускаются\n\n")
	var userInput string

	clrtxt.Foreground(clrtxt.Green, true)
	fmt.Scanln(&userInput)
	clrtxt.ResetColor()

	if strings.IndexByte(userInput, ' ') > 0 {
		return
	}

	var sum int
	if strings.IndexByte(userInput, '+') > 0 {
		sums := strings.Split(userInput, "+")
		for _, s := range sums {
			ss, err := strconv.Atoi(s)
			if !(err != nil) {
				sum += ss
			}
		}
	} else {
		var err error
		sum, err = strconv.Atoi(userInput)
		if err != nil {
			return
		}
	}

	if sum < 1 || sum > 9_999_999 {
		sayAndWait("Некорректная сумма")
		return
	}
	fmt.Println()

	showMenuOutgo()

	clrtxt.Foreground(clrtxt.White, true)
	fmt.Printf("\nСумма %d  Укажите код расхода\n", sum)
	clrtxt.ResetColor()

	ch, keyCode, _ := keyboard.GetSingleKey()
	if keyCode == keyboard.KeyEsc {
		return
	}
	if _, ok := ogClassifier[ch]; !ok && ch != '0' {
		sayAndWait("Неизвестный код расхода")
		return
	}

	it := outgoItem{int32(sum), ch, stat.dt}
	stat.addSum(stat.dt, int32(sum))

	outgo = append(outgo, it)
	if isCurDate(stat.dt) {
		outgoToday = append(outgoToday, it)
		ogToday = getSortTodayOG()
	}
	wasChanges = true
}

func getSortTodayOG() []outgoStr {
	dtMinusMonth := stat.dt.AddDate(0, -1, 1).Add(-1)
	dtEndCurDay := endOfDay(stat.dt).Add(1)

	a := make([]outgoStr, 0)

	for _, o := range outgo {
		if o.dt.After(dtMinusMonth) && o.dt.Before(dtEndCurDay) {
			find := false
			for i, v := range a {
				if v.code == o.code {
					a[i].sumMonth += o.sum
					if isCurDate(o.dt) {
						a[i].sumToday += o.sum
					}
					find = true
				}
			}
			if !find {
				v := outgoStr{code: o.code, name: getOutgo(o.code), sumMonth: o.sum}
				if isCurDate(o.dt) {
					v.sumToday = o.sum
				}
				a = append(a, v)
			}
		}
	}

	sort.Sort(byCodeOutgoStr(a))
	return a
}

func (o *statusInf) addSum(dt time.Time, sum int32) {
	if isCurDate(dt) {
		stat.outgoDay += sum
	}
	if isCurMonth(dt) {
		stat.outgoMonth += sum
	}
	if isPrevMonth(dt) {
		stat.outgoPrevMonth += sum
	}
	if isCurYear(dt) {
		stat.outgoYear += sum
	}
}

func (o *statusInf) clear() {
	o.dt = time.Now()
	o.outgoDay = 0
	o.outgoMonth = 0
	o.outgoPrevMonth = 0
	o.outgoYear = 0
}

func (o *statusInf) reCalculate() {
	for _, og := range outgo {
		o.addSum(og.dt, og.sum)
	}
}
