package models

import "errors"

// ErrNotFound возвращается когда не найден объект.
var ErrNotFound = errors.New("object not found")
