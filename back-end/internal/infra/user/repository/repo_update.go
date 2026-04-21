package repository

import (
	"log"

	"github.com/laurencefluciano/content-api/internal/domain/user"
	"github.com/laurencefluciano/content-api/internal/infra/database"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *gormUserRepository) Update(entity *user.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		userChanges := make(map[string]any)
		profileChanges := make(map[string]any)
		producerProfileChanges := make(map[string]any)

		dirtyFields := entity.GetDirtyFields()

		for field := range dirtyFields {
			switch field {
			case user.NameField:
				userChanges["name"] = entity.Name()
			case user.StatusField:
				userChanges["status"] = entity.Status()
			case user.AvatarUrlField:
				profileChanges["avatar_url"] = entity.AvatarUrl()
			case user.BioField:
				producerProfileChanges["bio"] = entity.Bio()
			case user.SocialField:
				producerProfileChanges["social"] = entity.Social()
			case user.WebsitesField:
				producerProfileChanges["websites"] = pq.StringArray(entity.Websites())
			case user.RolesField:
				roles := toRoleTable(entity.Roles())
				err := tx.Model(&database.UserTable{ID: entity.ID()}).
					Association("Roles").
					Replace(roles)

				if err != nil {
					return err
				}

			}
		}

		log.Printf("Entiy: %v\n", entity.Websites())
		log.Printf("Changes: %v\n", producerProfileChanges["websites"])

		if len(userChanges) > 0 {
			err := tx.
				Model(&database.UserTable{}).
				Where("id = ?", entity.ID()).
				Updates(userChanges).
				Error

			if err != nil {
				return err
			}
		}

		if len(profileChanges) > 0 {

			profileChanges["user_id"] = entity.ID()

			err := tx.
				Model(&database.UserProfileTable{}).
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_id"}},
					DoUpdates: clause.Assignments(profileChanges),
				}).
				Create(profileChanges).
				Error

			if err != nil {
				return err
			}
		}

		if len(producerProfileChanges) > 0 {

			producerProfileChanges["user_id"] = entity.ID()

			err := tx.
				Model(&database.ProducerProfileTable{}).
				Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "user_id"}},
					DoUpdates: clause.Assignments(producerProfileChanges),
				}).
				Create(producerProfileChanges).
				Error

			if err != nil {
				return err
			}
		}

		return nil
	})
}
