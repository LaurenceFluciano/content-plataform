package user

import (
	"errors"
	"strings"
)

type Bio struct {
	value string
}

func NewBio(value string) (Bio, error) {
	if len(value) > 500 {
		return Bio{}, errors.New("Bio deve ter no máximo 500 caracteres.")
	}
	if len(value) < 10 {
		return Bio{}, errors.New("Bio deve ter no minimo 10 caracteres.")
	}
	return Bio{value: strings.ToLower(value)}, nil
}

func (vo *Bio) Value() string {
	return vo.value
}

func (vo *Bio) IsEmpty() bool {
	return vo.value == ""
}
