package migrates

import (
	"errors"
	"github.com/go-gormigrate/gormigrate/v2"
)

type MigrateInterface interface {
	GetId() string
	Migrate() *gormigrate.Gormigrate
}

type ModelTableNameInterface interface {
	TableName() string
}

var _migrates = []MigrateInterface{
	&MigrateInit{},
}

func Migrate() error {
	_migrates = append(_migrates, mm...)
	for _, m := range _migrates {
		gm := m.Migrate()
		if err := gm.Migrate(); err != nil {
			_ = gm.RollbackLast()
			return errors.New("迁移" + m.GetId() + "执行失败 " + err.Error())
		}
	}
	return nil
}

func Rollback(id string) error {
	if id == "" {
		return nil
	}
	_migrates = append(_migrates, mm...)
	for _, m := range _migrates {
		if id == m.GetId() {
			gm := m.Migrate()
			if err := gm.RollbackLast(); err != nil {
				return errors.New("迁移回滚" + m.GetId() + "执行失败" + err.Error())
			}
			break
		}
	}
	return nil
}

//运行迁移
func MigrateFunc(migrateCmd, mRollbackId string) error {
	if migrateCmd == "" {
		return nil
	}
	if migrateCmd == "run" {
		return Migrate()
	} else if migrateCmd == "rollback" {
		if mRollbackId == "" {
			return errors.New("无效的回退版本号【请填写参数 mRollbackId")
		}
		return Rollback(mRollbackId)
	}
	return nil
}
