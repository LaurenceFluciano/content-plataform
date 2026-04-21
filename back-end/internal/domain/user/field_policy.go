package user

var editableFieldsByRole = map[RoleType][]Field{
	RoleReader:   {NameField, AvatarUrlField},
	RoleProducer: {BioField, WebsitesField, SocialField},
}

func CanChangeField(roles []RoleType, field Field) bool {
	for _, role := range roles {
		fields, ok := editableFieldsByRole[role]
		if !ok {
			continue
		}
		for _, f := range fields {
			if f == field {
				return true
			}
		}
	}
	return false
}
