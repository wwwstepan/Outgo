package main

import (
	"fmt"
	"time"

	clrtxt "github.com/daviddengcn/go-colortext"
	"github.com/eiannone/keyboard"
)

func inputDate(dtStart time.Time) time.Time {
	dt := dtStart

	for {
		clearScreen()
		fmt.Printf("======== Ввод даты ========\n")
		fmt.Println("   <- (-) (+) ->   +- день")
		fmt.Println("   v  (-) (+) ^    +- неделя")
		fmt.Println(" PgDn (-) (+) PgUp +- месяц")
		fmt.Println("  End (-) (+) Home +- год")
		fmt.Println("     Enter         Завершить ввод")
		fmt.Println(" любая другая клавиша (ESC, Пробел, 0) - отмена редактирования даты")
		fmt.Println()
		clrtxt.Foreground(clrtxt.Cyan, true)
		fmt.Println("  ", dateFmtExt(dt))
		clrtxt.ResetColor()

		_, keyCode, err := keyboard.GetSingleKey()

		if err != nil {
			return dtStart
		}

		switch keyCode {
		case keyboard.KeyArrowLeft:
			dt = dt.AddDate(0, 0, -1)
		case keyboard.KeyArrowRight:
			dt = dt.AddDate(0, 0, 1)
		case keyboard.KeyArrowDown:
			dt = dt.AddDate(0, 0, -7)
		case keyboard.KeyArrowUp:
			dt = dt.AddDate(0, 0, 7)
		case keyboard.KeyPgdn:
			dt = dt.AddDate(0, -1, 0)
		case keyboard.KeyPgup:
			dt = dt.AddDate(0, 1, 0)
		case keyboard.KeyEnd:
			dt = dt.AddDate(-1, 0, 0)
		case keyboard.KeyHome:
			dt = dt.AddDate(1, 0, 0)
		case keyboard.KeyEnter:
			return dt
		default:
			return dtStart
		}

	}

}
