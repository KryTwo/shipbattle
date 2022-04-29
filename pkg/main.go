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
	M "main/pkg/msg"
	"os"
	"regexp"
	"strconv"
	"time"
)

var lastHit fieldBuilder.Cell
var Direction string
var SelectLimit []int

func main() {
	MyField := fieldBuilder.NewField()
	EnemyField := fieldBuilder.NewField()

	MyField.SetShipRandom()
	EnemyField.SetShipRandom()

	MyField.ShowField()
	EnemyField.ShowField()

	chm := make(chan string)

	//MyField.ShowField()
	//MyField.SetShipRandom()
	//MyField.ShowField()
	//MyField := fieldBuilder.NewField()
	//EnemyField := fieldBuilder.NewField()
	//
	//EnemyField.SetShipRandom()
	//MyField.SetShipRandom()
	//
	//EnemyField.ShowField()
	//fmt.Println("---------------------")
	//MyField.ShowField()

	/*
		1. Начать игру с ботом
		 -> 1.1 рандомное расположение своих кораблей
		 -> 1.2 расстановка кораблей вручную
		2. Начать игру с соперником
		3. stop - остановить игру, вызывается из любого места программы
	*/
	//var in string
	//for {
	//	M.Print(M.MsgWelcome)
	//	fmt.Fscan(os.Stdin, &in)
	//	switch in {
	//	case "1": //начать игру с ботом
	//		selectStartOption(MyField, EnemyField)
	//	case "2": //начать игру онлайн
	//		fmt.Println("В разработке")
	//	case "0": //выйти
	//		fmt.Println("Выход в разработке")
	//	default:
	//		fmt.Println("игнор")
	//
	//		//неизвестная команда
	//	}
	//}

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
	M.Print(M.MsgHowToSetShip)
	fmt.Fscan(os.Stdin, &in)

	switch in {
	case "1":
		e.SetShipRandom()
		fmt.Println("Выбрано \"автоматически\"")
		m.SetShipRandom()
		fmt.Println("Ваше поле:")
		m.ShowField()
		fmt.Println("Поле врага:")
		//e.ShowField()
		e.ShowHiddenField()
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

func startGame(m, e *fieldBuilder.Field) {
	//coinFlipping()
	shipLeftMe := 20    //ячеек кораблей осталось
	shipLeftEnemy := 20 //ячеек кораблей осталось
	for shipLeftMe != 0 || shipLeftEnemy != 0 {
		//if shootMe(e, &shipLeftEnemy) {
		//	if shipLeftEnemy != 0 {
		//		shootMe(e, &shipLeftEnemy)
		//	} else {
		//		fmt.Println("вы выиграли")
		//		return
		//	}
		//}
		for {
			shootEnemy(m, &shipLeftMe)
			time.Sleep(500 * time.Millisecond)
		}
		//if shootEnemy(m, &shipLeftMe) {
		//	if shipLeftEnemy != 0 {
		//		shootEnemy(m, &shipLeftMe)
		//	} else {
		//		fmt.Println("вы проиграли")
		//		return
		//	}
		//}

	}

}

func CustomConvert(i fieldBuilder.Cell) string {
	c := strconv.Itoa(i.Column)
	r := strconv.Itoa(i.Row)
	buf := []byte(c)
	buf[0] += 49

	return string(buf[0]) + r
}

func shootEnemy(f *fieldBuilder.Field, shipLeftEnemy *int) bool {
	var leftEmptyCells []fieldBuilder.Cell
	for c := 0; c < 10; c++ {
		for r := 0; r < 10; r++ {
			if f[c][r].Hidden == fieldBuilder.Hidden {
				leftEmptyCells = append(leftEmptyCells, f[c][r])
			}
		}
	}

	var i int

	if lastHit.Hidden == "" {
		i = fieldBuilder.RandNum(len(leftEmptyCells))
	} else {
		i = selectNearCell(f, leftEmptyCells)
	}

	in := CustomConvert(leftEmptyCells[i])
	col, row := fieldBuilder.Convert(in)

	shoot(f, in, shipLeftEnemy)
	if f[col][row].Hidden == fieldBuilder.DestroyedShip {
		Direction = ""
		f.ShowField()
		f.ShowHiddenField()
		lastHit.Hidden = ""
		SelectLimit = nil
		Direction = ""
		return true
	}
	if f[col][row].Hidden == fieldBuilder.EmptyShip && lastHit.Hidden != "" {
		f.ShowField()
		f.ShowHiddenField()
		//time.Sleep(1000 * time.Millisecond)
		SelectLimit = nil
		return false
	}
	if f[col][row].Hidden == fieldBuilder.EmptyShip {
		lastHit.Hidden = ""
		f.ShowField()
		f.ShowHiddenField()
		//time.Sleep(1000 * time.Millisecond)
		return false
	}
	if f[col][row].Hidden == fieldBuilder.InjuredShip {
		//rotation := getRotationShip()
		SelectLimit = nil
		lastHit = leftEmptyCells[i]

		f.ShowField()
		f.ShowHiddenField()
		//time.Sleep(1000 * time.Millisecond)
		return true
	}

	return true
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

func shootMe(e *fieldBuilder.Field, shipLeftEnemy *int) bool {
	var in string
	M.Print(M.SelectCellToShoot)
	fmt.Fscan(os.Stdin, &in)
	if matched, _ := regexp.Match(`^[a-jA-J]\d$`, []byte(in)); matched == true {
		fmt.Println(*shipLeftEnemy)
		if !shoot(e, in, shipLeftEnemy) {
			e.ShowHiddenField()
			fmt.Println("missed")
			return false
		}
	}

	return true
}

func shoot(field *fieldBuilder.Field, in string, shipLeft *int) bool {
	col, row := fieldBuilder.Convert(in)

	isHere := isShipHere(field, col, row)
	allow := allowedToShoot(field, col, row)

	if !isHere && field[col][row].Hidden == fieldBuilder.Hidden {
		field[col][row].Hidden = fieldBuilder.EmptyShip
		return false
	}

	if !isHere && allow {
		fmt.Println("сюда стрелять нельзя")
		return true
	}

	if isHere && allow {
		size := fieldBuilder.ShipSize(field, col, row)
		leftDeck := getLeftDeck(field, col, row)
		if size != 1 { // однопалубник
			switch leftDeck {
			case 1: // будет убит
				destroy(field, col, row)
				*shipLeft--
			default: // будет ранен
				injure(field, col, row)
				*shipLeft--
			}
		} else {
			completelyDestroy(field, col, row)
			*shipLeft--
		}
		field.ShowHiddenField()
		return true
	}
	return false
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

func allowedToShoot(enemyField *fieldBuilder.Field, col, row int) bool {
	if enemyField[col][row].Hidden == fieldBuilder.InjuredShip {
		fmt.Println("вы уже стреляли туда")
		return false
	}
	if enemyField[col][row].Hidden == fieldBuilder.DestroyedShip {
		fmt.Println("вы уже стреляли туда")
		return false
	}
	if enemyField[col][row].Hidden == fieldBuilder.EmptyShip {
		return false
	}
	return true
}
