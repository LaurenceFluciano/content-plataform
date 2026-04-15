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
		return Name{}, errors.New("Name must have at lesat 2 chars")
	}
	return Name{value: strings.ToLower(value)}, nil
}

func (n *Name) Value() string {
	return n.value
}

func (n *Name) IsEmpty() bool {
	return n.value == ""
}
