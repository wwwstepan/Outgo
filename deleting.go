package main

import (
	"fmt"

	clrtxt "github.com/daviddengcn/go-colortext"
	"github.com/eiannone/keyboard"
)

func deleteOutgo() {
	if len(outgo) == 0 {
		sayAndWait("В базе нет записей")
		return
	}
	maxPage := (len(outgo) - 1) / 9
	page := maxPage
	for {
		showRecords(page)

		fmt.Println()
		if page > 0 {
			fmt.Println("  <-, v Показать более ранние записи")
		}
		if page < maxPage {
			fmt.Println("  ->, ^ Показать более поздние записи")
		}
		fmt.Println("  1-9 - Удалить запись с соответствующим номером")
		fmt.Println("  Любая другая клавиша (ESC, Пробел, 0) - выход")

		ch, keyCode, err := keyboard.GetSingleKey()
		if err != nil {
			return
		}

		if keyCode == keyboard.KeyArrowLeft ||
			keyCode == keyboard.KeyArrowDown {
			if page > 0 {
				page--
			}
		} else if keyCode == keyboard.KeyArrowRight ||
			keyCode == keyboard.KeyArrowUp {
			if page < maxPage {
				page++
			}
		} else {
			switch ch {
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				deleteRecord(page*9 + int(ch) - int('1'))
			default:
				return
			}
		}
	}
}

func showRecords(page int) {
	showStatusLine()

	fmt.Println()
	clrtxt.Foreground(clrtxt.Green, true)
	fmt.Println("\n===== Сумма == Статья расходов === Дата ====")
	clrtxt.ResetColor()

	nFirstR := page * 9
	if nFirstR < 0 {
		nFirstR = 0
	}

	n := 1
	for i := nFirstR; i < len(outgo) && n <= 9; i++ {
		if outgo[i].sum != 0 {
			s := "'" + itoa(n) + "' " + outgo[i].toString()
			fmt.Println(s)
			n++
		}
	}
}

func deleteRecord(nrec int) {
	if nrec < 0 || nrec >= len(outgo) {
		return
	}

	fmt.Println()
	fmt.Println("Удаление записи:")
	fmt.Println(outgo[nrec].toString())
	fmt.Println("Удалить? (Y/N) / (1/0) ?")
	if !askAndGetYes() {
		return
	}
	stat.addSum(stat.dt, -outgo[nrec].sum)
	outgo[nrec].sum = 0
	wasChanges = true
}
