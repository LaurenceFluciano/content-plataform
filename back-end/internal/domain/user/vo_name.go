package user

import (
	"errors"
	"strings"
)

type Name struct {
	value string
}

func NewName(value string) (Name, error) {
	if len(value) < 2 {
		return Name{}, errors.New("Nome deve ter pelo menos 2 caracteres")
	}
	if len(value) > 50 {
		return Name{}, errors.New("Nome deve ter no máximo 50 caracteres")
	}
	return Name{value: strings.ToLower(value)}, nil
}

func (n *Name) Value() string {
	return n.value
}

func (n *Name) IsEmpty() bool {
	return n.value == ""
}
