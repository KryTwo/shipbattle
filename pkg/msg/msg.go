package msg

import "fmt"

const (
	MsgWelcome = "Добро пожаловать в игру!\n" +
		"Выберите пункт меню:\n" +
		"1. Начать игру против компьютера.\n" +
		"2. Начать игру против человека.\n" +
		"0. Выйти из игры"
	MsgHowToSetShip = "Выберите вариант расстановки кораблей:\n" +
		"1. Автоматически\n" +
		"2. Вручную"
	SelectCellToShoot = "Выберите куда стрелять:"
)

func Print(s string) {
	fmt.Println(s)
}
