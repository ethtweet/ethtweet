package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethtweet/ethtweet/broadcastMsg"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"

	"gorm.io/gorm"
)

func RunTasks(ctx context.Context) {
	t := time.NewTicker(2 * time.Hour)
	defer t.Stop()
	global.GetDB().Where("status = ? or status = ?", models.TasksStatusFail, models.TasksStatusComplete).Delete(&models.Tasks{})
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			global.GetDB().Where("status = ? or status = ?", models.TasksStatusFail, models.TasksStatusComplete).Delete(&models.Tasks{})
		default:
		}
		task := &models.Tasks{}
		now := time.Now().Unix()
		if global.GetDB().
			Where("status = ? or (status = ? and max_exec_lock_time <= ?)", models.TasksStatusWait, models.TasksStatusIng, now).
			Where("next_exec_time <= ?", now).
			Order("sort asc").Limit(1).Find(task).RowsAffected == 0 {
			logs.PrintlnWarning("not tasks... sleep 3 sec")
			time.Sleep(3 * time.Second)
			continue
		}
		//防止并发的操作
		if global.GetDB().Table(task.TableName()).Where("id", task.ID).
			Where("next_exec_time <= ?", now).
			Where("status = ? or (status = ? and max_exec_lock_time <= ?)", models.TasksStatusWait, models.TasksStatusIng, now).
			Updates(map[string]interface{}{
				"max_exec_lock_time": time.Now().Add(5 * time.Minute).Unix(),
				"status":             models.TasksStatusIng,
			}).RowsAffected == 0 {
			logs.PrintlnWarning("task is locked ", task.ID)
			continue
		}
		task.Status = models.TasksStatusIng
		logs.PrintlnInfo(fmt.Sprintf("get task ... type %s id %s", task.Type, task.ID))
		_ = global.GetDB().Transaction(func(tx *gorm.DB) error {
			var err error
			switch task.Type {
			case models.TasksTypeUpIpfsAndBroadcastTweet:
				err = execUpIpfsAndBroadcastTweet(task, tx)
			default:
				task.Status = models.TasksStatusFail
				task.SetExtendsJson("err", "invalid type..")
				tx.Save(task)
				return nil
			}

			ExecAfter(task, tx, err)
			return nil
		})
	}
}

// 在事务中使用
func execUpIpfsAndBroadcastTweet(ts *models.Tasks, tx *gorm.DB) (err error) {
	if ts.Type != models.TasksTypeUpIpfsAndBroadcastTweet {
		return fmt.Errorf("invalid type ", ts.Type)
	}
	if tx == nil {
		return fmt.Errorf("invalid tx...")
	}
	twId := ts.GetExtendsJson("twId").String()
	//查找推文
	tw := &models.Tweets{}
	if tx.Where("id", twId).Find(tw).RowsAffected == 0 {
		ts.Status = models.TasksStatusFail
		tx.Select("status").Save(ts)
		err = fmt.Errorf("not found tweet by tweet id %s", twId)
		return
	}

	userId := tw.UserId
	//查找用户
	user := &models.User{}
	if tx.Where("id", userId).Find(user).RowsAffected == 0 {
		ts.Status = models.TasksStatusFail
		tx.Select("status").Save(ts)
		err = fmt.Errorf("not found user by user id %s", userId)
		return
	}
	defer func() {
		if ts.Status == models.TasksStatusComplete {
			logs.PrintlnSuccess("execBroadcastTweet is ok start BroadcastTweet..........")
			tw.UserInfo = user.GetUserInfoToPublic()
			broadcastMsg.BroadcastTweetSync(tw)
		}
	}()
	r := models.TweetJson{}

	//是否需要包含开始的nonce
	isEqStartNonce := false
	if user.LatestCid == "" {
		r.Nonce = 0
		isEqStartNonce = true
	} else {
		b, err2 := global.GetIpfsInfo(user.LatestCid)
		if err2 != nil {
			err = err2
			return
		}
		err2 = json.Unmarshal(b, &r)
		if err2 != nil {
			r.Nonce = 0
			user.LatestCid = ""
			isEqStartNonce = true
		}
	}

	if isEqStartNonce {
		if r.Nonce > tw.Nonce {
			ts.Status = models.TasksStatusComplete
			err = tx.Select("status").Save(ts).Error
			return
		}
		if r.Nonce == tw.Nonce {
			_, err = tw.UpIpfs(user)
			if err != nil {
				return
			}
			//标记完成
			ts.Status = models.TasksStatusComplete
			err = tx.Select("status").Save(ts).Error
			if err != nil {
				return
			}
			err = tx.Select("LatestCid").Save(user).Error
			return
		}
	} else {
		if r.Nonce >= tw.Nonce {
			ts.Status = models.TasksStatusComplete
			err = tx.Select("status").Save(ts).Error
			return
		}
		//这里表示是连续的
		if r.Nonce+1 == tw.Nonce {
			_, err = tw.UpIpfs(user)
			if err != nil {
				return
			}
			//标记完成
			ts.Status = models.TasksStatusComplete
			tx.Select("status").Save(ts)
			err = tx.Select("LatestCid").Save(user).Error
			return
		}
	}

	query := tx.Where("user_id", userId)
	if isEqStartNonce {
		query.Where("nonce >= ? and nonce <= ?", r.Nonce, tw.Nonce)
	} else {
		query.Where("nonce > ? and nonce <= ?", r.Nonce, tw.Nonce)
	}
	limit := 10
	tws := make([]*models.Tweets, 0, limit)
	if query.Limit(10).Order("nonce asc").Find(&tws).RowsAffected == 0 {
		ts.Status = models.TasksStatusFail
		err = tx.Select("status").Save(ts).Error
		if err != nil {
			return
		}
		err = fmt.Errorf("not found user %s tweets range nonce %d~%d", userId, r.Nonce, tw.Nonce)
		return
	}
	ok := 0
	for k, rTw := range tws {
		//这里要确保连续性
		if k > 0 && rTw.Nonce != tws[k-1].Nonce+1 {
			break
		}
		_, err = rTw.UpIpfs(user)
		if err != nil {
			break
		}
		ok++
		time.Sleep(500 * time.Millisecond)
	}
	l := len(tws)
	if ok == l && tws[l-1].Nonce == tw.Nonce {
		ts.Status = models.TasksStatusComplete
		err = tx.Select("status").Save(ts).Error
		if err != nil {
			return
		}
	}
	err = tx.Select("LatestCid").Save(user).Error
	return
}

func ExecAfter(ts *models.Tasks, tx *gorm.DB, err error) {
	if err != nil {
		ts.SetExtendsJson("err", err.Error())
		logs.PrintErr(fmt.Sprintf("err task ... type %s id %s, err %s", ts.Type, ts.ID, err.Error()))
	} else {
		logs.PrintlnSuccess(fmt.Sprintf("exec task success ... type %s id %s", ts.Type, ts.ID))
	}
	//这里需要修改状态
	if ts.Status == models.TasksStatusIng || ts.Status == models.TasksStatusWait {
		ts.Status = models.TasksStatusWait
		ts.MaxExecLockTime = 0
		ts.NextExecTime = time.Now().Add(3 * time.Minute).Unix()
		tx.Save(ts)
	} else if ts.Status == models.TasksStatusFail {
		tx.Select("ExtraData").Save(ts)
	}
}
