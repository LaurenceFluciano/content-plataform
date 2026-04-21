package repository

import (
	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra/database"
)

func (r *gormUserRepository) toDomain(table *database.UserTable) *user.User {
	if table == nil {
		return nil
	}

	u := user.RestoreUser(
		table.ID,
		table.AuthId,
		table.Name,
		table.Status,
		toRoleType(table.Roles),
	)

	if table.UserProfile != nil {
		u.RestoreProfile(table.UserProfile.AvatarUrl)
	}

	websitesDomain := make([]string, len(table.ProducerProfile.Websites))
	copy(websitesDomain, table.ProducerProfile.Websites)

	if table.ProducerProfile != nil {
		u.RestoreProducer(
			table.ProducerProfile.Bio,
			websitesDomain,
			table.ProducerProfile.Social,
		)
	}

	return u
}

func toRoleTable(roles []user.RoleType) []database.RoleTable {
	tableRoles := make([]database.RoleTable, len(roles))
	for i, r := range roles {
		tableRoles[i] = database.RoleTable{
			RoleId: r,
		}
	}
	return tableRoles
}

func toRoleType(roles []database.RoleTable) []user.RoleType {
	rolesType := make([]user.RoleType, len(roles))
	for i, r := range roles {
		rolesType[i] = r.RoleId
	}
	return rolesType
}
