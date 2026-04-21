package user

import (
	"errors"
)

type Profile struct {
	bio       Bio
	websites  []Website
	social    SocialLinks
	avatarUrl AvatarUrl
}

func (u *User) Bio() string { return u.profile.bio.Value() }
func (u *User) Websites() []string {
	if u.profile == nil {
		return nil
	}

	res := make([]string, len(u.profile.websites))
	for i, w := range u.profile.websites {
		res[i] = w.Value()
	}

	return res
}
func (u *User) Social() SocialLinks { return u.profile.social }
func (u *User) AvatarUrl() string   { return u.profile.avatarUrl.Value() }

func (u *User) ChangeBio(bio Bio) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	if u.profile.bio == bio {
		return
	}
	u.profile.bio = bio
	u.dirty[BioField] = struct{}{}
}

func (u *User) ChangeAvatarUrl(url AvatarUrl) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	if u.profile.avatarUrl.Value() == url.Value() {
		return
	}
	u.profile.avatarUrl = url
	u.dirty[AvatarUrlField] = struct{}{}
}

func (u *User) ChangeSocialLinks(newSocial SocialLinks) {
	if u.profile == nil {
		u.profile = &Profile{}
	}

	if u.profile.social == newSocial {
		return
	}

	u.profile.social = newSocial
	u.dirty[SocialField] = struct{}{}
}

func (u *User) ChangeWebsites(websites []Website) error {
	if len(websites) > 10 {
		return errors.New("limite de 10 websites atingido")
	}

	if u.profile == nil {
		u.profile = &Profile{}
	}

	u.profile.websites = websites
	u.dirty[WebsitesField] = struct{}{}
	return nil
}
