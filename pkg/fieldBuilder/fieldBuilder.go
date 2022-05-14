package fieldBuilder

import (
	"fmt"
	"main/pkg/msg"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	Empty         = "." // пустая ячейка
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
var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func callClear() {
	//value, ok := clear[runtime.GOOS]
	//if ok {
	//	value()
	//} else {
	//	panic("Your platform is unsupported! I can't clear terminal screen :(")
	//}
}

type Cell struct {
	Column     int
	Row        int
	StatusCode string
	Hidden     string
	HiddenMe   string
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
			field[c][r].HiddenMe = Hidden
		}
	}
	return &field
}

func (f *Field) OldShowField() {
	fmt.Println("  0 1 2 3 4 5 6 7 8 9")
	fmt.Println("  A B C D E F G H I J")
	for i := 0; i < 10; i++ {
		fmt.Print(strconv.Itoa(i) + " ")
		for r := 0; r < 10; r++ {
			fmt.Print(f[r][i].StatusCode)
			fmt.Print(" ")
		}
		fmt.Println()
	}
}

func ShowField(MyField, EnemyField *Field) {
	callClear()
	chm := make(chan string)
	che := make(chan string)

	go Show(MyField, chm)
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
			t = t + f[r][i].HiddenMe
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

// SetShip - ручная расстановка кораблей
func (f *Field) SetShip(s string) bool {
	col, row := Convert(strings.ToLower(s))
	if ok, mErr := IsCan(f, col, row); ok {
		f[col][row].StatusCode = Ship
		f[col][row].HiddenMe = Ship
	} else {
		fmt.Println(mErr)
		return false
	}

	return true
}

func IsCan(f *Field, col, row int) (bool, string) {
	// false при попытке установит корабль на место существующего
	if ok, mErr := IsShipHere(f, col, row); !ok {
		return false, mErr
	}

	// проверка диагональных ячеек на наличие корабля
	if ok, mErr := CheckDiag(f, col, row); !ok {
		return false, mErr
	}

	// ShipSize узнает текущий размер корабля, к которому мы добавляем палубу
	size := ShipSize(f, col, row)

	// максимальный размер корабля не выше 4
	if size > 4 {
		return false, msg.MsgShipTooBig
	}

	// CheckSizeLim возвращает false, если достигнут лимит в количестве кораблей
	if ok, mErr := CheckSizeLim(size); !ok {
		return false, mErr
	}

	return true, ""
}

func CheckSizeLim(s int) (bool, string) {
	switch s {
	case 1:
		if shipsOnField["uno"] < unoSize {
			shipsOnField["uno"] += 1
		} else {
			return false, msg.MsgLimitOneDeckShips
		}
	case 2:
		if shipsOnField["does"] < doesSize {
			shipsOnField["uno"] -= 1
			shipsOnField["does"] += 1
		} else {
			return false, msg.MsgLimitTwoDeckShips
		}
	case 3:
		if shipsOnField["tres"] < tresSize {
			shipsOnField["does"] -= 1
			shipsOnField["tres"] += 1
		} else {
			return false, msg.MsgLimitThreeDeckShips
		}
	case 4:
		if shipsOnField["cuatro"] < cuatroSize {
			shipsOnField["tres"] -= 1
			shipsOnField["cuatro"] += 1
		} else {
			return false, msg.MsgLimitFourDeckShips
		}

	}
	return true, ""

}

// IsShipHere false если попытка установить на уже существующий корабль
func IsShipHere(f *Field, col, row int) (bool, string) {
	if f[col][row].StatusCode == Ship {
		return false, msg.MsgShipAlreadySetHere
	}

	return true, ""
}

// CheckDiag false if wrong place, true if right place
func CheckDiag(f *Field, col, row int) (bool, string) {
	// проверка диагональных ячеек на наличие корабля
	for r := row - 1; r <= row+1; r += 2 {
		for c := col - 1; c <= col+1; c += 2 {
			if c >= 0 && c <= 9 && r >= 0 && r <= 9 {
				if f[c][r].StatusCode == Ship {
					return false, msg.MsgInvalidShipPosition
				}
			}
		}
	}
	return true, ""
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

func (f *Field) ManualSetShip(in string) bool {
	if f.SetShip(in) {
		return true
	}
	return false
}
func ClearShipsOnField() {
	shipsOnField["uno"] = 0
	shipsOnField["does"] = 0
	shipsOnField["tres"] = 0
	shipsOnField["cuatro"] = 0
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

	ClearShipsOnField()
	return true // true if all ok
}

func setUnoShip(f *Field, in int) int {
	//fmt.Println("1 pal try")
	col := RandNum(10)
	row := RandNum(10)

	ok, _ := IsShipHere(f, col, row)
	ok2, _ := CheckDiag(f, col, row)

	//fmt.Printf("Ship:%v, col:%v, row:%v\n", 1, col, row)
	if !ok || !ok2 || ShipSize(f, col, row) > 1 {
		return in
	} else {

		f[col][row].StatusCode = Ship
		f[col][row].HiddenMe = Ship

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
			ok, _ := IsShipHere(f, c, row)
			ok2, _ := CheckDiag(f, c, row)
			if !ok || !ok2 || ShipSize(f, c, row) > 2 {
				return in
			}
		}
		for c := col; c < col+2; c++ {
			f[c][row].StatusCode = Ship
			f[c][row].HiddenMe = Ship
		}
		return in + 1
	default:
		for r := row; r < row+2; r++ {
			ok, _ := IsShipHere(f, col, r)
			ok2, _ := CheckDiag(f, col, r)
			if !ok || !ok2 || ShipSize(f, col, r) > 2 {
				return in
			}
		}
		for r := row; r < row+2; r++ {
			f[col][r].StatusCode = Ship
			f[col][r].HiddenMe = Ship
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
			ok, _ := IsShipHere(f, c, row)
			ok2, _ := CheckDiag(f, c, row)
			if !ok || !ok2 || ShipSize(f, c, row) > 3 {
				return in
			}
		}
		for c := col; c < col+3; c++ {
			f[c][row].StatusCode = Ship
			f[c][row].HiddenMe = Ship
		}
		return in + 1
	default:
		for r := row; r < row+3; r++ {
			ok, _ := IsShipHere(f, col, r)
			ok2, _ := CheckDiag(f, col, r)
			if !ok || !ok2 || ShipSize(f, col, r) > 3 {
				return in
			}
		}
		for r := row; r < row+3; r++ {
			f[col][r].StatusCode = Ship
			f[col][r].HiddenMe = Ship
		}
		return in + 1

	}

}

func setCuatroShip(f *Field) {
	col := RandNum(7)
	row := RandNum(7)

	direction := RandDir()
	switch direction {
	case "right":
		for c := col; c < col+4; c++ {
			f[c][row].StatusCode = Ship
			f[c][row].HiddenMe = Ship
		}
	default: //"bottom"
		for r := row; r < row+4; r++ {
			f[col][r].StatusCode = Ship
			f[col][r].HiddenMe = Ship

		}
	}
	defer func() {
		shipsOnField["cuatro"] += 1

	}()
}
