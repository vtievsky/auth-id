// Модели для обмена данными между слоями
package models

type Privilege struct {
	ID          uint64 // Ключевое поле для БД
	Code        string // Ключевое поле для интерфейса
	Name        string
	Description string
}
