package user

import (
	"errors"
)

type Profile struct {
	bio       string
	websites  []string
	social    SocialLinks
	avatarUrl string
}

func (p *Profile) Bio() string         { return p.bio }
func (p *Profile) Websites() []string  { return p.websites }
func (p *Profile) Social() SocialLinks { return p.social }
func (p *Profile) AvatarUrl() string   { return p.avatarUrl }

func (u *User) SetBio(bio string) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	if u.profile.bio == bio {
		return
	}
	u.profile.bio = bio
	u.dirty[BioField] = struct{}{}
}

func (u *User) SetAvatarUrl(url string) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	if u.profile.avatarUrl == url {
		return
	}
	u.profile.avatarUrl = url
	u.dirty[AvatarUrlField] = struct{}{}
}

func (u *User) ChangeSocialLinks(newSocial SocialLinks) {
	if u.profile == nil {
		u.profile = &Profile{}
	}

	merged := u.profile.social.Merge(newSocial)

	if u.profile.social == merged {
		return
	}

	u.profile.social = merged
	u.dirty[SocialField] = struct{}{}
}

func (u *User) AddWebsite(url string) error {
	if u.profile == nil {
		u.profile = &Profile{}
	}

	if len(u.profile.websites) >= 10 {
		return errors.New("limite de 10 websites atingido")
	}

	u.profile.websites = append(u.profile.websites, url)
	u.dirty[WebsitesField] = struct{}{}
	return nil
}

func (u *User) RemoveWebsite(url string) {
	if u.profile == nil || len(u.profile.websites) == 0 {
		return
	}

	for i, w := range u.profile.websites {
		if w == url {
			u.profile.websites = append(u.profile.websites[:i], u.profile.websites[i+1:]...)
			u.dirty[WebsitesField] = struct{}{}
			return
		}
	}
}
