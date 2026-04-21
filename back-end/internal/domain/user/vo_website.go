package user

import "errors"

type Website struct {
	value string
}

func NewWebsite(url string) (Website, error) {
	if url == "" {
		return Website{}, errors.New("Um website não pode ser vazio. Para apagar considere enviar um novo estado")
	}
	if len(url) > 2048 {
		return Website{}, errors.New("A url fornecida é muito grande, o máximo de caracteres suportado é de 2048.")
	}

	return Website{value: url}, nil
}

func (w Website) Value() string { return w.value }
