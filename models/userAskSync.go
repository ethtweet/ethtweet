package models

type UserAskSync struct {
	UserId          string `gorm:"type:varchar(100);primarykey"`
	SyncStatus      uint8  `gorm:"default:0;index:idx_sync_status;comment:同步状态 0未同步 1同步中 2已同步,该状态只在本地存在作为同步校验"`
	SyncTimeoutUnix int64  `gorm:"default:0;comment:同步超时时间，用于判断同步的超时 如果当前时间戳大于该时间并且状态在同步中则表示同步超时"`
}

func (uas *UserAskSync) TableName() string {
	return "user_ask_sync"
}
