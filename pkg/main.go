/*
	TODO:...
	1. Сделать поле
	2. Инструмент ручной установки кораблей
	2.1 Ограничения расстановки кораблей
	3. Инструмент рандомной генерации кораблей на поле - для бота или для ленивого
	4. Разделение полей свой\чужой
	5. Инструмент реализующий выстрелы
	...
*/
package main

import (
	"fmt"
	"log"
	"main/pkg/fieldBuilder"
	"main/pkg/msg"
	s "main/pkg/status"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var lastHit fieldBuilder.Cell
var Direction string
var SelectLimit []int
var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Can't exec Run, %v", err)
		}
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatalf("Can't exec Run, %v", err)
		}
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
func main() {
	callClear()

	var in string
	for {
		MyField := fieldBuilder.NewField()
		EnemyField := fieldBuilder.NewField()
		fmt.Println(msg.MsgWelcome)
		if _, err := fmt.Fscan(os.Stdin, &in); err != nil {
			log.Fatalf("Input error: %v", err)
		}
		switch in {
		case "1": //начать игру с ботом
			callClear()
			selectStartOption(MyField, EnemyField)
		case "2": //начать игру онлайн
			callClear()
			fmt.Println("В разработке")
		//case "0": //выйти
		//	callClear()
		//	fmt.Println("Куда ты собрался выходить???")
		default:
			callClear()
			fmt.Println("Я тебя игнорирую")

		}
	}
}

// coinFlipping определяет право первого хода
//func coinFlipping() bool {
//	i := fieldBuilder.RandNum(2)
//	if i == 0 {
//		return true
//	}
//	return false
//}

func clearVariables(m, e *fieldBuilder.Field) {
	m = fieldBuilder.NewField()
	e = fieldBuilder.NewField()
	SelectLimit = nil
}

func selectStartOption(m, e *fieldBuilder.Field) {
	var in string
	fmt.Println(msg.MsgHowToSetShip)
	if _, err := fmt.Fscan(os.Stdin, &in); err != nil {
		log.Fatalf("Input errpr, %v", err)
	}

	switch in {
	case "1":
		e.SetShipRandom()
		m.SetShipRandom()
		fieldBuilder.ShowField(m, e)
		startGame(m, e)
		clearVariables(m, e)
		//автоматические корабли
	case "2":
		fieldBuilder.ShowField(m, e)
		e.SetShipRandom()
		manualSet(m, e)
		startGame(m, e)
		//Корабли вручную
	case "0":
		fmt.Println("Выход в предыдущее меню")
		break
		//выйти
	}

}

func manualSet(m, e *fieldBuilder.Field) {
	var in string
	var i int

	for i < 20 {
		fmt.Println(msg.MsgSelectCellToSetShip)
		if _, err := fmt.Fscan(os.Stdin, &in); err != nil {
			log.Fatalf("Input errpr, %v", err)
		}

		if matched, _ := regexp.Match(`^[a-jA-J]\d$`, []byte(in)); matched != true {

			fieldBuilder.ShowField(m, e)
			continue
		}
		if m.ManualSetShip(in) {
			i++
		}
		fieldBuilder.ShowField(m, e)

	}
	fieldBuilder.ClearShipsOnField()

}
func checkScore(me, enemy int) bool {
	if me == 0 {
		fmt.Println("Потрачено")
		return true
	}
	if enemy == 0 {
		fmt.Println("Сегодня ты выиграл, но везение не вечно")
		return true
	}
	return false
}
func startGame(m, e *fieldBuilder.Field) {
	//coinFlipping()
	shipLeftMe := 20    //ячеек кораблей осталось
	shipLeftEnemy := 20 //ячеек кораблей осталось

	for {
		shootMe(e, m, &shipLeftEnemy)
		if checkScore(shipLeftMe, shipLeftEnemy) {
			return
		}

		shootEnemy(m, e, &shipLeftMe)
		if checkScore(shipLeftMe, shipLeftEnemy) {
			return
		}
	}

}

func waitingCommand(m, e *fieldBuilder.Field) string {
	var in string
	fmt.Println(msg.MsgSelectCellToShoot)
	if _, err := fmt.Fscan(os.Stdin, &in); err != nil {
		log.Fatalf("Input errpr, %v", err)
	}

	if matched, _ := regexp.Match(`^[a-jA-J]\d$`, []byte(in)); matched == true {
		return in
	}

	fieldBuilder.ShowField(m, e)
	fmt.Println(msg.MsgWrongCommand)

	return waitingCommand(m, e)
}

func shootMe(e, m *fieldBuilder.Field, shipLeftEnemy *int) bool {
	in := waitingCommand(m, e)

	result := shoot(e, in, shipLeftEnemy)

	if result == s.Miss {
		fmt.Println(result)
		return false
	}

	fieldBuilder.ShowField(m, e)
	fmt.Println(result)

	return shootMe(e, m, shipLeftEnemy)
}

func shoot(field *fieldBuilder.Field, in string, shipLeft *int) string {
	col, row := fieldBuilder.Convert(in)

	isHere := isShipHere(field, col, row)
	allow := allowedToShoot(field, col, row)

	if !allow {
		return s.DoubleShot
	}

	if !isHere && field[col][row].Hidden == fieldBuilder.Hidden {
		field[col][row].Hidden = fieldBuilder.EmptyShip
		field[col][row].HiddenMe = fieldBuilder.EmptyShip
		return s.Miss
	}

	if isHere && allow {
		size := fieldBuilder.ShipSize(field, col, row)
		leftDeck := getLeftDeck(field, col, row)
		if size != 1 { // однопалубник
			switch leftDeck {
			case 1: // будет убит
				destroy(field, col, row)
				*shipLeft--
				return s.Destroy
			default: // будет ранен
				injure(field, col, row)
				*shipLeft--
				return s.Injured
			}
		} else {
			completelyDestroy(field, col, row)
			*shipLeft--
			return s.Destroy
		}
	}
	return ""
}

func shootEnemy(m, e *fieldBuilder.Field, shipLeftEnemy *int) bool {
	defer fieldBuilder.ShowField(m, e)
	var leftEmptyCells []fieldBuilder.Cell
	for c := 0; c < 10; c++ {
		for r := 0; r < 10; r++ {
			if m[c][r].Hidden == fieldBuilder.Hidden {
				leftEmptyCells = append(leftEmptyCells, m[c][r])
			}
		}
	}

	var i int

	if lastHit.Hidden == "" {
		i = fieldBuilder.RandNum(len(leftEmptyCells))
	} else {
		i = selectNearCell(m, leftEmptyCells)
	}

	in := convert(leftEmptyCells[i])
	col, row := fieldBuilder.Convert(in)

	shoot(m, in, shipLeftEnemy)
	if m[col][row].Hidden == fieldBuilder.DestroyedShip {
		Direction = ""
		lastHit.Hidden = ""
		SelectLimit = nil
		return shootEnemy(m, e, shipLeftEnemy)
	}
	if m[col][row].Hidden == fieldBuilder.EmptyShip && lastHit.Hidden != "" {
		//time.Sleep(1000 * time.Millisecond)
		SelectLimit = nil
		return false
	}
	if m[col][row].Hidden == fieldBuilder.EmptyShip {
		lastHit.Hidden = ""
		//time.Sleep(1000 * time.Millisecond)
		return false
	}
	if m[col][row].Hidden == fieldBuilder.InjuredShip {
		//rotation := getRotationShip()
		SelectLimit = nil
		lastHit = leftEmptyCells[i]

		//time.Sleep(1000 * time.Millisecond)
		return shootEnemy(m, e, shipLeftEnemy)
	}

	return shootEnemy(m, e, shipLeftEnemy)
}

func allowedToShoot(enemyField *fieldBuilder.Field, col, row int) bool {
	if enemyField[col][row].Hidden == fieldBuilder.InjuredShip {
		return false
	}
	if enemyField[col][row].Hidden == fieldBuilder.DestroyedShip {
		return false
	}
	if enemyField[col][row].Hidden == fieldBuilder.EmptyShip {
		return false
	}
	return true
}

func convert(i fieldBuilder.Cell) string {
	c := strconv.Itoa(i.Column)
	r := strconv.Itoa(i.Row)
	buf := []byte(c)
	buf[0] += 49

	return string(buf[0]) + r
}

func selectNearCell(f *fieldBuilder.Field, leftEmptyCell []fieldBuilder.Cell) int {
	/*
		1 - up
		2 - down
		3 - left
		4 - right
	*/

	//update lasthit, if 4 sides missed
	col := lastHit.Column
	row := lastHit.Row

	if len(SelectLimit) == 4 {
		SelectLimit = nil
		switch Direction {
		case "right":
			if col-1 >= 0 && f[col-1][row].Hidden == fieldBuilder.InjuredShip {
				lastHit.Column -= 1
				col -= 1
			} else {
				lastHit.Column += 1
				col += 1
			}
		case "bottom":
			if row-1 >= 0 && f[col][row-1].Hidden == fieldBuilder.InjuredShip {
				lastHit.Row -= 1
				row -= 1
			} else {
				lastHit.Row += 1
				row += 1
			}
		}
	}
	r := fieldBuilder.RandNum(4) + 1
	for ContainsInt(SelectLimit, r) {
		r = fieldBuilder.RandNum(4) + 1
	}
	//if !ContainsInt(SelectLimit, r) {
	SelectLimit = append(SelectLimit, r)
	//}

	switch r {
	case 1: // up
		row -= 1
	case 2: // down
		row += 1
	case 3: // left
		col -= 1
	case 4: // right
		col += 1
	}
	if col >= 0 && col < 10 && row >= 0 && row < 10 && f[col][row].Hidden == fieldBuilder.Hidden {
		if r == 1 || r == 2 {
			Direction = "bottom"
		} else {
			Direction = "right"
		}
		if ContainsCell(leftEmptyCell, f[col][row]) {
			q := find(leftEmptyCell, f[col][row])
			return q
		}
	}
	return selectNearCell(f, leftEmptyCell)
}

func ContainsCell(a []fieldBuilder.Cell, x fieldBuilder.Cell) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func ContainsInt(s []int, i int) bool {
	for _, n := range s {
		if i == n {
			return true
		}
	}
	return false
}

func find(a []fieldBuilder.Cell, x fieldBuilder.Cell) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

func completelyDestroy(e *fieldBuilder.Field, col, row int) {
	e[col][row].Hidden = fieldBuilder.DestroyedShip
	e[col][row].HiddenMe = fieldBuilder.DestroyedShip
	setEmptyShipDiag(e, col, row)
	setEmptyShipAround(e, col, row)
}

func destroy(e *fieldBuilder.Field, col, row int) {

	var dir string

	if col-1 >= 0 && e[col-1][row].StatusCode == fieldBuilder.Ship {
		dir = "right"
	}
	if col+1 < 10 && e[col+1][row].StatusCode == fieldBuilder.Ship {
		dir = "right"
	}
	if row-1 >= 0 && e[col][row-1].StatusCode == fieldBuilder.Ship {
		dir = "bottom"
	}
	if row+1 < 10 && e[col][row+1].StatusCode == fieldBuilder.Ship {
		dir = "bottom"
	}

	switch dir {
	case "right":
		for c := col; c < col+4; c++ {
			if c < 10 {
				if e[c][row].StatusCode == fieldBuilder.Ship {
					setEmptyShipDiag(e, c, row)
					e[c][row].Hidden = fieldBuilder.DestroyedShip
					e[c][row].HiddenMe = fieldBuilder.DestroyedShip
					setEmptyShipAround(e, c, row)
				} else {
					break
				}
			}
		}
		for c := col; c > col-4; c-- {
			if c >= 0 {
				if e[c][row].StatusCode == fieldBuilder.Ship {
					setEmptyShipDiag(e, c, row)
					e[c][row].Hidden = fieldBuilder.DestroyedShip
					e[c][row].HiddenMe = fieldBuilder.DestroyedShip
					setEmptyShipAround(e, c, row)
				} else {
					break
				}
			}
		}
	case "bottom":
		//c5 - 2|5
		for r := row; r < row+4; r++ {
			if r < 10 {
				if e[col][r].StatusCode == fieldBuilder.Ship {
					setEmptyShipDiag(e, col, r)
					e[col][r].Hidden = fieldBuilder.DestroyedShip
					e[col][r].HiddenMe = fieldBuilder.DestroyedShip
					setEmptyShipAround(e, col, r)
				} else {
					break
				}
			}
		}
		for r := row; r > row-4; r-- {
			if r >= 0 {
				if e[col][r].StatusCode == fieldBuilder.Ship {
					setEmptyShipDiag(e, col, r)
					e[col][r].Hidden = fieldBuilder.DestroyedShip
					e[col][r].HiddenMe = fieldBuilder.DestroyedShip
					setEmptyShipAround(e, col, r)
				} else {
					break
				}
			}
		}
	default:
		fmt.Println("cant get direction - destroy")
	}
}

func injure(e *fieldBuilder.Field, col, row int) {
	e[col][row].Hidden = fieldBuilder.InjuredShip
	e[col][row].HiddenMe = fieldBuilder.InjuredShip
	setEmptyShipDiag(e, col, row)

}

func setEmptyShipAround(e *fieldBuilder.Field, col, row int) {
	for c := col - 1; c <= col+1; c++ {
		for r := row - 1; r <= row+1; r++ {
			if c >= 0 && c <= 9 && r >= 0 && r <= 9 {
				if e[c][r].Hidden == fieldBuilder.Hidden {
					e[c][r].Hidden = fieldBuilder.EmptyShip
					e[c][r].HiddenMe = fieldBuilder.EmptyShip
				}
			}
		}
	}
}

func setEmptyShipDiag(e *fieldBuilder.Field, col, row int) {
	for r := row - 1; r <= row+1; r += 2 {
		for c := col - 1; c <= col+1; c += 2 {
			if c >= 0 && c <= 9 && r >= 0 && r <= 9 {
				e[c][r].Hidden = fieldBuilder.EmptyShip
				e[c][r].HiddenMe = fieldBuilder.EmptyShip
			}
		}
	}
}

func isShipHere(f *fieldBuilder.Field, col, row int) bool {
	if f[col][row].StatusCode == fieldBuilder.Ship {
		return true
	}

	return false
}

func getLeftDeck(f *fieldBuilder.Field, col, row int) int {
	leftDeck := 1

	for i := 1; i <= 4; i++ {
		if col-i >= 0 {
			if f[col-i][row].StatusCode == fieldBuilder.Ship {
				if f[col-i][row].Hidden == fieldBuilder.Hidden {
					leftDeck++
				}

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if col+i <= 9 {
			if f[col+i][row].StatusCode == fieldBuilder.Ship {
				if f[col+i][row].Hidden == fieldBuilder.Hidden {
					leftDeck++
				}

			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if row-i >= 0 {
			if f[col][row-i].StatusCode == fieldBuilder.Ship {
				if f[col][row-i].Hidden == fieldBuilder.Hidden {
					leftDeck++
				}
			} else {
				break
			}
		}
	}

	for i := 1; i <= 4; i++ {
		if row+i <= 9 {
			if f[col][row+i].StatusCode == fieldBuilder.Ship {
				if f[col][row+i].Hidden == fieldBuilder.Hidden {
					leftDeck++
				}

			} else {
				break
			}
		}
	}

	return leftDeck

}
