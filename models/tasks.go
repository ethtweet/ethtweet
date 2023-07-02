package models

import (
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models/mField"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TasksTypeUpIpfsAndBroadcastTweet = "syncIpfsAndBroadcastTweet"

	TasksStatusWait     = 0
	TasksStatusIng      = 1
	TasksStatusComplete = 2
	TasksStatusFail     = 3
)

type Tasks struct {
	ID                           string `gorm:"primarykey"`
	Type                         string `gorm:"default:'';comment:执行类型，根据不同的type分配不同的方法"`
	Sort                         int64  `gorm:"default:100;comment:执行顺序优先级"`
	Status                       uint8  `gorm:"default:0;comment:执行状态:0待执行，1执行中，2执行完成, 3执行失败;index:status"`
	NextExecTime                 int64  `gorm:"default:0;comment:下一次执行时间,必须是小于当前时间的才能够执行(主要对于执行失败了但是不能被丢弃需要重新执行的任务);index:status"`
	MaxExecLockTime              int64  `gorm:"default:0;comment:最大锁定时间，也就是执行中的时候最大的一个时间 超过该时间可被获取 用于锁定一条执行中的记录"`
	mField.FieldsExtendsJsonType        //执行的上下文数据
	mField.FieldsTimeUnixModel
}

func (ts *Tasks) TableName() string {
	return "tasks"
}

func (ts *Tasks) BeforeCreate(tx *gorm.DB) error {
	ts.ID = uuid.New().String()
	return nil
}

func AddUpIpfsAndBroadcastTweetTask(tw *Tweets, sort int64) (*Tasks, error) {
	logs.PrintlnInfo("AddUpIpfsAndBroadcastTweetTask ", tw.Id, tw.UserId)
	if sort == 0 {
		sort = 100
	}
	task := &Tasks{
		Type:         TasksTypeUpIpfsAndBroadcastTweet,
		Sort:         sort,
		Status:       TasksStatusWait,
		NextExecTime: 0, //立即执行
	}
	task.SetExtendsJson("twId", tw.Id)
	err := global.GetDB().Create(task).Error
	if err != nil {

		return nil, err
	}
	logs.PrintlnSuccess("AddUpIpfsAndBroadcastTweetTask success ", tw.Id, tw.UserId)
	return task, nil
}
