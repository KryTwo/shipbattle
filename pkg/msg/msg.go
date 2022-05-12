package msg

const (
	MsgWelcome = "Добро пожаловать в игру!\n" +
		"Выберите пункт меню:\n" +
		"1. Начать игру против компьютера.\n" +
		"2. Начать игру против человека.\n" +
		"0. Выйти из игры"

	MsgHowToSetShip = "Выберите вариант расстановки кораблей:\n" +
		"1. Автоматически\n" +
		"2. Вручную\n" +
		"0. Назад"

	MsgSelectCellToShoot = "Выберите куда стрелять:"

	MsgWrongCommand = "Неправильный ввод"

	MsgDoubleShot = "Вы не можете выстрелить сюда"

	MsgShipAlreadySetHere  = "Здесь уже есть корабль"
	MsgLimitOneDeckShips   = "Достигнут лимит однопалубных кораблей"
	MsgLimitTwoDeckShips   = "Достигнут лимит двупалубных кораблей"
	MsgLimitThreeDeckShips = "Достигнут лимит трехалубных кораблей"
	MsgLimitFourDeckShips  = "Достигнут лимит четырехпалубных кораблей"

	MsgShipTooBig          = "Корабль слишком большой"
	MsgInvalidShipPosition = "Недопустимое расположение корабля"
	MsgSelectCellToSetShip = "Выберите где поставить корабль:"
)
