package user

var editableFieldsByRole = map[RoleType][]string{
	RoleReader:   {"Name", "AvatarUrl"},
	RoleProducer: {"Bio", "Websites", "Social"},
}

func GetEditableFieldsMap(roles []RoleType) map[string]bool {
	permissions := make(map[string]bool)

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
