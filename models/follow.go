package models

type Follow struct {
	UserId     string `gorm:"type:varchar(100);NOT NULL;comment:用户的id;uniqueIndex:idx_uid_follow_unique;index"`
	FollowedID string `gorm:"type:varchar(100);NOT NULL;comment:关注的用户id;uniqueIndex:idx_uid_follow_unique;index"`
	CreatedAt  int64  `gorm:"index;default:0;autoCreateTime"`
	UpdatedAt  int64  `gorm:"autoUpdateTime;default:0"`
}

func (flw *Follow) TableName() string {
	return "follow"
}
