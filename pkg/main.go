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
	"main/pkg/fieldBuilder"
	msg "main/pkg/msg"
	s "main/pkg/status"
	"os"
	"regexp"
	"strconv"
)

var lastHit fieldBuilder.Cell
var Direction string
var SelectLimit []int
var Wrong string

//func selectOne(chm, che chan string) {
//	for {
//		select {}
//	}
//}

func main() {
	MyField := fieldBuilder.NewField()
	EnemyField := fieldBuilder.NewField()
	/*
		1. Начать игру с ботом
		 -> 1.1 рандомное расположение своих кораблей
		 -> 1.2 расстановка кораблей вручную
		2. Начать игру с соперником
		3. stop - остановить игру, вызывается из любого места программы
	*/
	var in string
	for {
		fmt.Println(msg.MsgWelcome)
		fmt.Fscan(os.Stdin, &in)
		switch in {
		case "1": //начать игру с ботом
			selectStartOption(MyField, EnemyField)
		case "2": //начать игру онлайн
			fmt.Println("В разработке")
		case "0": //выйти
			fmt.Println("Куда ты собрался выходить???")
		default:
			fmt.Println("Я тебя игнорирую")

		}
	}

	//MyField.SetShip(in)
	//MyField.ShowField()

	//if matched, _ := regexp.Match(`^[a-jA-J]\d$`, []byte(in)); matched != true {
	//	fmt.Println("Неверный ввод. Пример: h0 или D4")
	//	MyField.ShowField()

	//if matched, _ := regexp.Match(`^\d$`, []byte(in)); !matched {
	//	fmt.Println("Неверный ввод. Пример: 1 или 3")
	//	continue
	//}

}

// coinFlipping определяет право первого хода
func coinFlipping() bool {
	i := fieldBuilder.RandNum(2)
	if i == 0 {
		return true
	}
	return false
}

func selectStartOption(m, e *fieldBuilder.Field) {
	var in string
	fmt.Println(msg.MsgHowToSetShip)
	fmt.Fscan(os.Stdin, &in)

	switch in {
	case "1":
		e.SetShipRandom()
		fmt.Println("Выбрано \"автоматически\"")
		m.SetShipRandom()
		fieldBuilder.ShowField(m, e)
		startGame(m, e)

		//автоматические корабли
	case "2":
		fmt.Println("ручной режим в разработке")
		break
		//Корабли вручную
	case "0":
		fmt.Println("Выход в предыдущее меню")
		break
		//выйти
	}

}
func checkScore(me, enemy int) bool {
	if me == 0 {
		fmt.Println("Лузер, проебал")
		return true
	}
	if enemy == 0 {
		fmt.Println("Блэд, ты выиграл...")
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

func waitingCommand() string {
	var in string
	fmt.Println(msg.MsgSelectCellToShoot)
	fmt.Fscan(os.Stdin, &in)
	if matched, _ := regexp.Match(`^[a-jA-J]\d$`, []byte(in)); matched == true {
		return in
	}
	fmt.Println(msg.MsgWrongCommand)
	return waitingCommand()
}

func shootMe(e, m *fieldBuilder.Field, shipLeftEnemy *int) bool {
	in := waitingCommand()

	result := shoot(e, in, shipLeftEnemy)
	if result == s.Miss {
		return false
	}
	fieldBuilder.ShowField(m, e)
	return shootMe(e, m, shipLeftEnemy)
}

func shoot(field *fieldBuilder.Field, in string, shipLeft *int) string {
	col, row := fieldBuilder.Convert(in)

	isHere := isShipHere(field, col, row)
	allow := allowedToShoot(field, col, row)

	if !allow {
		Wrong = msg.MsgDoubleShot
		return s.DoubleShot
	}

	if !isHere && field[col][row].Hidden == fieldBuilder.Hidden {
		field[col][row].Hidden = fieldBuilder.EmptyShip
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

	in := CustomConvert(leftEmptyCells[i])
	col, row := fieldBuilder.Convert(in)

	shoot(m, in, shipLeftEnemy)
	if m[col][row].Hidden == fieldBuilder.DestroyedShip {
		Direction = ""
		lastHit.Hidden = ""
		SelectLimit = nil
		Direction = ""
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

func CustomConvert(i fieldBuilder.Cell) string {
	c := strconv.Itoa(i.Column)
	r := strconv.Itoa(i.Row)
	buf := []byte(c)
	buf[0] += 49

	return string(buf[0]) + r
}

//func getRotationShip(f *fieldBuilder.Field, col, row int) string {
//	if col-1 >= 0 && col+1 <10 {
//		if f[col-1][row].StatusCode
//	}
//	return ""
//}

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
			q := Find(leftEmptyCell, f[col][row])
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

func Find(a []fieldBuilder.Cell, x fieldBuilder.Cell) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

func completelyDestroy(e *fieldBuilder.Field, col, row int) {
	e[col][row].Hidden = fieldBuilder.DestroyedShip
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
	setEmptyShipDiag(e, col, row)

}

func setEmptyShipAround(e *fieldBuilder.Field, col, row int) {
	for c := col - 1; c <= col+1; c++ {
		for r := row - 1; r <= row+1; r++ {
			if c >= 0 && c <= 9 && r >= 0 && r <= 9 {
				if e[c][r].Hidden == fieldBuilder.Hidden {
					e[c][r].Hidden = fieldBuilder.EmptyShip
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
