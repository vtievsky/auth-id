// Модели для обмена данными между слоями
package models

type Privilege struct {
	ID          int    // Ключевое поле для БД
	Code        string // Ключевое поле для интерфейса
	Name        string
	Description string
}
