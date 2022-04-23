/*
	1. исключить повторную установку кораблей
*/
package main

import (
	"fmt"
	m "main/pkg"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	empty       = "~" // пустая ячейка
	ship        = "S" // ячейка с кораблем
	injuredShip = "I" // поврежденный корабль
	brokenShip  = "D" // уничтоженный корабль
	missShot    = "m" // промах
	emptyShip   = "X" // пустая ячейка без корабля

	// ships size quantity
	unoSize    = 4
	doesSize   = 3
	tresSize   = 2
	cuatroSize = 1
)

var shipsOnField = map[string]int{
	"uno":    0,
	"does":   0,
	"tres":   0,
	"cuatro": 0,
}

type cell struct {
	digit      int
	statusCode string
}

type field [10][10]cell

func newField() *field {
	var field field
	for c := 0; c < 10; c++ {
		for r := 0; r < 10; r++ {
			field[c][r].digit = c
			field[c][r].statusCode = empty
		}
	}
	return &field
}

func (f *field) showField() {
	fmt.Println("  0 1 2 3 4 5 6 7 8 9")
	fmt.Println("  A B C D E F G H I J")
	for i := 0; i < 10; i++ {
		fmt.Print(strconv.Itoa(i) + " ")
		for r := 0; r < 10; r++ {
			fmt.Print(f[r][i].statusCode)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func convert(s string) (int, int) {
	r := []byte(s)
	r[0] -= 49

	col := int(r[0] - 48)
	row := int(r[1] - 48)

	return col, row
}

func (f *field) setShip(s string) {
	col, row := convert(strings.ToLower(s))
	if isCan(f, col, row) {
		f[col][row].statusCode = ship
	}
}

func checkSizeLim(s int) bool {
	switch s {
	case 1:
		if shipsOnField["uno"] < unoSize {
			shipsOnField["uno"] += 1
		} else {
			fmt.Println("лимит однопалубных кораблей")
			return false
		}
	case 2:
		if shipsOnField["does"] < doesSize {
			shipsOnField["uno"] -= 1
			shipsOnField["does"] += 1
		} else {
			fmt.Println("лимит двупалубных кораблей")
			return false
		}
	case 3:
		if shipsOnField["tres"] < tresSize {
			shipsOnField["does"] -= 1
			shipsOnField["tres"] += 1
		} else {
			fmt.Println("лимит трехпалубных кораблей")
			return false
		}
	case 4:
		if shipsOnField["cuatro"] < cuatroSize {
			shipsOnField["tres"] -= 1
			shipsOnField["cuatro"] += 1
		} else {
			fmt.Println("лимит четырехпалубных кораблей")
			return false
		}

	}
	return true

}

func isCan(f *field, col, row int) bool {
	// false при попытке установит корабль на место существующего
	if f[col][row].statusCode == ship {
		fmt.Println("корабль уже стоит")
		return false
	}

	// shipSize узнает текущий размер корабля, к которому мы добавляем палубу
	size := shipSize(f, col, row)

	// checkSizeLim возвращает false, если достигнут лимит в количестве кораблей
	if !checkSizeLim(size) {
		return false
	}

	// максимальный размер корабля не выше 4
	if size > 4 {
		fmt.Println("корабль слишком большой")
		return false
	}

	// проверка диагональных ячеек на наличие корабля
	for r := row - 1; r <= row+1; r += 2 {
		for c := col - 1; c <= col+1; c += 2 {
			if c >= 0 && c <= 9 && r >= 0 && r <= 9 {
				if f[c][r].statusCode != empty && f[c][r].statusCode != emptyShip {
					fmt.Println(m.WrongPlace)
					return false
				}
			}
		}
	}

	return true
}

func shipSize(f *field, col, row int) int {
	size := 1

	for i := 1; i <= 4; i++ {
		if col-i >= 0 {
			if f[col-i][row].statusCode == ship {
				size++

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if col+i <= 9 {
			if f[col+i][row].statusCode == ship {
				size++

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if row-i >= 0 {
			if f[col][row-i].statusCode == ship {
				size++

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if row+i <= 9 {
			if f[col][row+i].statusCode == ship {
				size++

			} else {
				break
			}
		}
	}

	return size

}

func main() {
	f := newField()

	f.showField()

	var i string
	for {
		fmt.Fscan(os.Stdin, &i)
		if matched, _ := regexp.Match(`^[a-jA-J]\d$`, []byte(i)); matched != true {
			fmt.Println("Неверный ввод. Пример: h0 или D4")
			f.showField()

			continue
		}

		f.setShip(i)
		f.showField()

	}
}
