package broadcastMsg

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"

	"gorm.io/gorm"
)

func IpfsSync() {
	for {
		logs.Println("IpfsSync 1")
		check()
		time.Sleep(time.Second * 60)
	}
}

func check() {
	var users []*models.User
	sdb := global.GetDB()
	sdb.Limit(10).Where("local_nonce < nonce or (local_nonce = 0 and latest_cid != '')").Order("id desc").Find(&users)
	for _, user := range users {
		go func(user *models.User) {
			if !user.SyncStatusComplete() {
				uak := NewUserInfoAsk(user)
				_ = uak.DoAsk()
				return
			}
			SyncByUser(user, nil)
		}(user)
	}
}

func SyncByUser(usr *models.User, sdb *gorm.DB) {
	cid := usr.LatestCid
	if sdb == nil {
		sdb = global.GetDB()
	}
	logs.PrintlnInfo("---------------------------- sync user info ", usr.Id)
	nonces := make(map[uint64]struct{}, 1024)
	isFinish := false
	var lastNonce uint64 = 0
	for {
		if cid == "" {
			return
		}
		r := SyncByCid(cid)
		if r == nil {
			return
		}

		if len(r.Sign) < 10 {
			return
		}
		tw := &models.Tweets{
			UserId:     r.UserId,
			Content:    r.Content,
			Attachment: r.Attachment,
			Nonce:      r.Nonce,
			Sign:       r.Sign,
			CreatedAt:  r.CreatedAt,
		}
		cid = r.PreviousCid
		if !keys.VerifySignature(usr.PubKey, r.Sign, tw.GetSignMsg()) {
			logs.Println("usr.PubKey " + usr.PubKey)
			logs.PrintErr("===========sync verify msg err, refuse...", tw.Content, "user id => ", tw.UserId)
			return
		}
		logs.PrintlnInfo(fmt.Sprintf("======================get user %s, lastCid %s, tweetNonce %d", usr.Id, cid, tw.Nonce))
		if usr.LocalNonce == tw.Nonce {
			logs.PrintlnInfo(fmt.Sprintf("======================= wait finish handle........... userId:%s, localNonce:%d, nonce:%d", usr.Id, usr.LocalNonce, tw.Nonce))
			//获取同步进度
			tmp_nonce := tw.Nonce
			for tmp_nonce <= usr.Nonce {
				tmp_nonce++
				if _, ok := nonces[tmp_nonce]; ok {
					continue
				} else {
					tmp_nonce--
					break
				}
			}
			usr.LocalNonce = tmp_nonce
			sdb.Select("LocalNonce").Save(usr)

			//标记可以结束
			isFinish = true
		}

		tw.GenerateId()

		var c int64
		if sdb.Table(tw.TableName()).Where("id = ?", tw.Id).Count(&c); c != 0 {
			nonces[tw.Nonce] = struct{}{}
			logs.PrintlnInfo("data is saved，no more forwarding", tw.Content)
		} else {
			if err := sdb.Create(tw).Error; err != nil {
				sdb.Where("user_id = ? and nonce=?", tw.UserId, tw.Nonce).Delete(&models.Tweets{})
				logs.PrintErr("tweets save err ", tw.Content)
			} else {
				nonces[tw.Nonce] = struct{}{}
				logs.PrintlnSuccess("ipfs tweet save success：", tw.Content, tw.Id)
			}
		}

		//这里表示ipfs上的推文发生了断层
		if r.PreviousCid == "" && tw.Nonce < lastNonce-1 {
			logs.PrintlnWarning(fmt.Sprintf("user %s has tweet fault, set localNonce is %d, origin localNonce is %d", usr.Id, usr.Nonce, usr.LocalNonce))
			for k, _ := range nonces {
				if k > usr.LocalNonce {
					usr.LocalNonce = k
				}
			}
			sdb.Select("LocalNonce").Save(usr)
		}
		lastNonce = tw.Nonce
		if isFinish {
			logs.PrintlnSuccess("ipfs sync tweet nonce success：", usr.LocalNonce)
			logs.PrintlnSuccess("ipfs sync finish")
			return
		}
	}
}

func SyncByCid(cid string) *models.TweetJson {
	body, err := global.GetIpfsInfo(cid)
	if err != nil {
		logs.PrintErr(err)
		return nil
	}
	r := &models.TweetJson{}
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil
	}
	return r
}
