package user

import (
	"errors"
	"path/filepath"
	"strings"
)

type AvatarUrl struct {
	value string
}

func NewAvatarUrl(path string) (AvatarUrl, error) {
	if path == "" {
		return AvatarUrl{}, errors.New("O caminho do avatar não pode ser vazio")
	}

	if !strings.HasPrefix(path, "avatars/") {
		return AvatarUrl{}, errors.New("O caminho do avatar deve começar com 'avatars/'")
	}

	ext := strings.ToLower(filepath.Ext(path))
	validExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}

	if !validExtensions[ext] {
		return AvatarUrl{}, errors.New("Formato de imagem inválido para o avatar")
	}

	if len(path) > 2048 {
		return AvatarUrl{}, errors.New("Máximo de 2048 caracteres atingido")
	}

	return AvatarUrl{value: path}, nil
}

func (a AvatarUrl) Value() string {
	return a.value
}
