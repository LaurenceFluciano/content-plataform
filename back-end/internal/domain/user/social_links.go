package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
)

type SocialLinks struct {
	Facebook  string `json:"facebook"`
	Youtube   string `json:"youtube"`
	Instagram string `json:"instagram"`
	Tiktok    string `json:"tiktok"`
	X         string `json:"x"`
}

func (s SocialLinks) Merge(new SocialLinks) SocialLinks {
	choose := func(old, incoming string) string {
		incoming = strings.TrimSpace(incoming)

		if incoming == "-" {
			return ""
		}

		if incoming == "" {
			return old
		}
		return incoming
	}

	return SocialLinks{
		Facebook:  choose(s.Facebook, new.Facebook),
		Youtube:   choose(s.Youtube, new.Youtube),
		Instagram: choose(s.Instagram, new.Instagram),
		Tiktok:    choose(s.Tiktok, new.Tiktok),
		X:         choose(s.X, new.X),
	}
}

func (s SocialLinks) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *SocialLinks) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, s)
}
