package fieldBuilder

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	Empty         = "~" // пустая ячейка
	Ship          = "S" // ячейка с кораблем
	InjuredShip   = "I" // поврежденный корабль
	DestroyedShip = "D" // уничтоженный корабль
	MissShot      = "m" // промах
	EmptyShip     = "×" // пустая ячейка гарантированно без корабля
	Hidden        = "·" // неизвестное поле для отображения сопернику

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

type Cell struct {
	Column     int
	Row        int
	StatusCode string
	Hidden     string
}

type Field [10][10]Cell

func NewField() *Field {
	var field Field
	for c := 0; c < 10; c++ {
		for r := 0; r < 10; r++ {
			field[c][r].Column = c
			field[c][r].Row = r
			field[c][r].StatusCode = Empty
			field[c][r].Hidden = Hidden
		}
	}
	return &field
}

func ShowField(MyField, EnemyField *Field) {
	chm := make(chan string)
	che := make(chan string)

	go ShowHiddenField(MyField, chm)
	go ShowHiddenField(EnemyField, che)

	for i := 0; i < 12; i++ {
		fmt.Print(<-che)
		fmt.Print("      ")
		fmt.Print(<-chm)
		fmt.Println()
	}
}

func Show(f *Field, ch chan string) {
	var t string
	ch <- "  0 1 2 3 4 5 6 7 8 9 "
	ch <- "  A B C D E F G H I J "
	for i := 0; i < 10; i++ {
		t = t + strconv.Itoa(i) + " "
		for r := 0; r < 10; r++ {
			t = t + f[r][i].StatusCode
			t = t + " "
		}
		ch <- t
		t = ""
	}
	defer close(ch)
}
func ShowHiddenField(f *Field, ch chan string) {
	var t string
	ch <- "  0 1 2 3 4 5 6 7 8 9 "
	ch <- "  A B C D E F G H I J "
	for i := 0; i < 10; i++ {
		t = t + strconv.Itoa(i) + " "
		for r := 0; r < 10; r++ {
			t = t + f[r][i].Hidden
			t = t + " "
		}
		ch <- t
		t = ""
	}
	defer close(ch)
}

func Convert(s string) (int, int) {
	r := []byte(s)
	r[0] -= 49

	col := int(r[0] - 48)
	row := int(r[1] - 48)

	return col, row
}

func (f *Field) SetShip(s string) {
	col, row := Convert(strings.ToLower(s))
	if IsCan(f, col, row) {
		f[col][row].StatusCode = Ship
	}
}

func CheckSizeLim(s int) bool {
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

// IsShipHere false если попытка установить на уже существующий корабль
func IsShipHere(f *Field, col, row int) bool {
	if f[col][row].StatusCode == Ship {
		fmt.Println("корабль уже стоит")
		return false
	}

	return true
}

func IsCan(f *Field, col, row int) bool {
	// false при попытке установит корабль на место существующего
	if !IsShipHere(f, col, row) {
		return false
	}

	// ShipSize узнает текущий размер корабля, к которому мы добавляем палубу
	size := ShipSize(f, col, row)

	// CheckSizeLim возвращает false, если достигнут лимит в количестве кораблей
	if !CheckSizeLim(size) {
		return false
	}

	// максимальный размер корабля не выше 4
	if size > 4 {
		fmt.Println("корабль слишком большой")
		return false
	}

	// проверка диагональных ячеек на наличие корабля
	if !CheckDiag(f, col, row) {
		return false
	}

	return true
}

// CheckDiag false if wrong place, true if right place
func CheckDiag(f *Field, col, row int) bool {
	// проверка диагональных ячеек на наличие корабля
	for r := row - 1; r <= row+1; r += 2 {
		for c := col - 1; c <= col+1; c += 2 {
			if c >= 0 && c <= 9 && r >= 0 && r <= 9 {
				if f[c][r].StatusCode == Ship {
					fmt.Println("Недопустимое расположение корабля")
					return false
				}
			}
		}
	}
	return true
}

func ShipSize(f *Field, col, row int) int {
	size := 1

	for i := 1; i <= 4; i++ {
		if col-i >= 0 {
			if f[col-i][row].StatusCode == Ship {
				size++

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if col+i <= 9 {
			if f[col+i][row].StatusCode == Ship {
				size++

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if row-i >= 0 {
			if f[col][row-i].StatusCode == Ship {
				size++

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if row+i <= 9 {
			if f[col][row+i].StatusCode == Ship {
				size++

			} else {
				break
			}
		}
	}

	return size

}

func RandNum(l int) int {
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r1.Intn(l)
}

// RandDir выбирает случайное направление размещения корабля
func RandDir() string {
	s1 := RandNum(2)
	switch s1 {
	case 1:
		return "right"
	default:
		return "bottom"
	}
}

func (f *Field) SetShipRandom() bool {
	setCuatroShip(f)
	i := 0
	for i < 2 {
		i = setTresShip(f, i)
	}
	i2 := 0
	for i2 < 3 {
		i2 = setDoesShip(f, i2)
	}

	i3 := 0
	for i3 < 4 {
		i3 = setUnoShip(f, i3)
	}
	return true // true if all ok
}

func setUnoShip(f *Field, in int) int {
	//fmt.Println("1 pal try")
	col := RandNum(10)
	row := RandNum(10)

	//fmt.Printf("Ship:%v, col:%v, row:%v\n", 1, col, row)
	if !IsShipHere(f, col, row) || !CheckDiag(f, col, row) || ShipSize(f, col, row) > 1 {
		return in
	} else {

		f[col][row].StatusCode = Ship
		fmt.Println("Ship set")

		return in + 1
	}
}

func setDoesShip(f *Field, in int) int {
	//fmt.Println("2 pal try")
	col := RandNum(9)
	row := RandNum(9)

	direction := RandDir()
	//fmt.Printf("Ship: %v, col:%v, row:%v, dir:%v\n", 2, col, row, direction)
	switch direction {
	case "right":
		for c := col; c < col+2; c++ {
			if !IsShipHere(f, c, row) || !CheckDiag(f, c, row) || ShipSize(f, c, row) > 2 {
				return in
			}
		}
		for c := col; c < col+2; c++ {
			f[c][row].StatusCode = Ship
			fmt.Println("Ship set")
		}
		return in + 1
	default:
		for r := row; r < row+2; r++ {
			if !IsShipHere(f, col, r) || !CheckDiag(f, col, r) || ShipSize(f, col, r) > 2 {
				return in
			}
		}
		for r := row; r < row+2; r++ {
			f[col][r].StatusCode = Ship
			fmt.Println("Ship set")
		}
		return in + 1

	}

}

func setTresShip(f *Field, in int) int {
	//fmt.Println("3 pal try")
	col := RandNum(8)
	row := RandNum(8)

	direction := RandDir()
	//fmt.Printf("Ship: %v, col:%v, row:%v, dir:%v\n", 3, col, row, direction)
	switch direction {
	case "right":
		for c := col; c < col+3; c++ {
			if !IsShipHere(f, c, row) || !CheckDiag(f, c, row) || ShipSize(f, c, row) > 3 {
				return in
			}
		}
		for c := col; c < col+3; c++ {
			f[c][row].StatusCode = Ship
			fmt.Println("Ship set")
		}
		return in + 1
	default:
		for r := row; r < row+3; r++ {
			if !IsShipHere(f, col, r) || !CheckDiag(f, col, r) || ShipSize(f, col, r) > 3 {
				return in
			}
		}
		for r := row; r < row+3; r++ {
			f[col][r].StatusCode = Ship
			fmt.Println("Ship set")
		}
		return in + 1

	}

}

func setCuatroShip(f *Field) {
	//fmt.Println("4 pal try")
	col := 0
	row := 1
	//col := RandNum(7)
	//row := RandNum(7)

	direction := "right" //RandDir()
	switch direction {
	case "right":
		for c := col; c < col+4; c++ {
			f[c][row].StatusCode = Ship
		}
	default: //"bottom"
		for r := row; r < row+4; r++ {
			f[col][r].StatusCode = Ship
		}
	}
	defer func() {
		shipsOnField["cuatro"] += 1

	}()
}
