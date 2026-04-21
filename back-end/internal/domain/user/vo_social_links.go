package user

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type SocialLinks struct {
	Facebook  string `json:"facebook"`
	Youtube   string `json:"youtube"`
	Instagram string `json:"instagram"`
	Tiktok    string `json:"tiktok"`
	X         string `json:"x"`
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
