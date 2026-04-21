package user

import "github.com/google/uuid"

func RestoreUser(id uuid.UUID, authId string, name string, status Status, roles []RoleType) *User {
	return &User{
		id:     id,
		authId: authId,
		name:   Name{value: name},
		status: status,
		roles:  roles,
		dirty:  make(map[Field]struct{}),
	}
}

func (u *User) RestoreProfile(avatarUrl string) {
	if u.profile == nil {
		u.profile = &Profile{}
	}
	u.profile.avatarUrl = AvatarUrl{value: avatarUrl}
}

func (u *User) RestoreProducer(bio string, websites []string, social SocialLinks) {
	if u.profile == nil {
		u.profile = &Profile{}
	}

	var voWebsites []Website
	for _, w := range websites {
		voWebsites = append(voWebsites, Website{value: w})
	}

	u.profile.bio = Bio{value: bio}
	u.profile.websites = voWebsites
	u.profile.social = social
}
