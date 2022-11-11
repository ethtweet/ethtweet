package migrates

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"github.com/ethtweet/ethtweet/global"
)

type MigrateInit struct{}

func (m *MigrateInit) GetId() string {
	return "init"
}

func (m *MigrateInit) Migrate() *gormigrate.Gormigrate {
	db := global.GetDB()
	mArr := []ModelTableNameInterface{}
	return gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: m.GetId(),
			Migrate: func(tx *gorm.DB) error {
				for _, ma := range mArr {
					err := tx.Migrator().AutoMigrate(ma)
					if err != nil {
						return err
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				for _, ma := range mArr {
					err := tx.Migrator().DropTable(ma.TableName())
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
	})
}
