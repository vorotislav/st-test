package models

import (
	"net/http"
	"time"
)

// Item описывает объект, включает в себя id, тело объекта как массив байт и дату, через которую надо удалить объект.
type Item struct {
	ID      int
	Body    []byte
	Expires time.Duration
}

// ToJSON возвращает объект как массив байт. Но т.к. у нас объект и так любое валидный json-объект, то данный метод является реализацией интерфейса.
func (i Item) ToJSON() ([]byte, error) {
	return i.Body, nil
}

// StatusCode возвращает статус код.
func (i Item) StatusCode() int {
	return http.StatusOK
}
