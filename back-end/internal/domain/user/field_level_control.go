package user

var editableFieldsByRole = map[RoleType][]Field{
	RoleReader:   {NameField, AvatarUrlField},
	RoleProducer: {BioField, WebsitesField, SocialField},
}

func GetEditableFieldsMap(roles []RoleType) map[Field]bool {
	permissions := make(map[Field]bool)

	for _, role := range roles {
		editableFields, ok := editableFieldsByRole[role]
		if !ok {
			continue
		}

		for _, field := range editableFields {
			permissions[field] = true
		}
	}

	return permissions
}
