package migrates

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/models"
)

type Migrate_202111191445 struct{}

func (m *Migrate_202111191445) GetId() string {
	return "Migrate_202111191445"
}

func (m *Migrate_202111191445) Migrate() *gormigrate.Gormigrate {
	db := global.GetDB()
	mArr := []ModelTableNameInterface{
		&models.User{},
		&models.Follow{},
		&models.Tweets{},
		&models.UserAskSync{},
		&models.Tasks{},
	}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: m.GetId(),
			Migrate: func(tx *gorm.DB) error {
				for _, ma := range mArr {
					if tx.Migrator().HasTable(ma) {
						err := tx.Migrator().DropTable(ma)
						if err != nil {
							return err
						}
					}
				}
				for _, ma := range mArr {
					err := tx.Migrator().AutoMigrate(ma)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
	})
}
