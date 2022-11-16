package broadcastMsg

import (
	"context"
	"errors"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/p2pNet"
	"gorm.io/gorm"
	"sync"
	"time"
)

var twReceiveMu sync.Mutex

/*var broadSendLogLen = 256
var broadSendLog = make(map[string]struct{}, broadSendLogLen)*/

type TweetInfo struct {
	Tw *models.Tweets
}

// 发布中心用户的文章发布
// 即 自有用户自己签名后 发布转发
func CenterUserRelease(tw *models.Tweets) error {
	var uNonce int64 = -1
	err := global.GetDB().Transaction(func(tx *gorm.DB) error {
		userLock := &models.User{}
		if global.LockForUpdate(tx.Model(userLock).Where("id = ?", tw.UserId)).Find(userLock).RowsAffected == 0 {
			return fmt.Errorf("not found user")
		}
		if !keys.VerifySignatureByAddress(tw.UserId, tw.Sign, tw.GetSignMsg()) {
			return fmt.Errorf("tw sign err %s %s %s", userLock.Id, tw.Sign, tw.GetSignMsg())
		}
		if !userLock.SyncStatusComplete() {
			return global.ErrWaitUserSync
		}
		if tw.Nonce == 0 {
			if userLock.LatestCid != "" {
				return fmt.Errorf("invalid nonce, need 1 --- 0, latestCid ", userLock.LatestCid)
			}
		} else if tw.Nonce != userLock.Nonce+1 {
			return fmt.Errorf("invalid nonce, need %d --- %s", userLock.Nonce+1, userLock.Id)
		}

		userLock.Nonce = tw.Nonce
		var err error
		if err = tx.Save(userLock).Error; err != nil {
			return err
		}
		if err = tx.Create(tw).Error; err != nil {
			var i int64
			//nonce已存在
			if tx.Model(&models.Tweets{}).Where("user_id = ? and nonce = ?", tw.UserId, tw.Nonce).Count(&i); i > 0 {
				uNonce = int64(tw.Nonce)
			}
			return err
		}
		return nil
	})
	if err != nil {
		if uNonce >= 0 {
			logs.PrintlnWarning("Tw nonce already exists ", tw.Nonce, " update user nonce")
			global.GetDB().Where("id", tw.UserId).Update("nonce", uNonce)
		}
		return err
	}
	go BroadcastTweet(tw)
	return nil
}

func ReleaseTweet(user *models.User, keyName, content, attachment, forwardId, topicTag string, createdAt int64) (*models.Tweets, error) {
	if !user.SyncStatusComplete() {
		return nil, global.ErrWaitUserSync
	}
	isOk := false
	tx := global.GetDB().Begin()
	var uNonce int64 = -1
	defer func() {
		if r := recover(); r != nil || !isOk {
			logs.PrintErr(r)
			tx.Rollback()
		} else {
			tx.Commit()
		}
		if uNonce >= 0 {
			global.GetDB().Where("user_id", user.Id).Update("nonce", uNonce)
		}
	}()
	if global.LockForUpdate(tx.Where("id", user.Id)).Find(user).RowsAffected == 0 {
		return nil, fmt.Errorf("user lock fail")
	}
	//todo 读取用户表，获得上一条推文的cid
	if user.Nonce != 0 {
		user.Nonce++
	} else {
		if user.LatestCid != "" {
			user.Nonce = 1
		}
	}
	Nonce := user.Nonce

	tw := &models.Tweets{
		UserId:     user.Id,
		Content:    content,
		Attachment: attachment,
		Nonce:      Nonce,
		TopicTag:   topicTag,
	}

	if forwardId != "" {
		originTw := &models.Tweets{}
		if tx.Model(originTw).Where("id = ?", forwardId).Find(originTw).RowsAffected == 0 {
			return nil, errors.New("invalid forward id " + forwardId)
		}
		tw.OriginTwId = originTw.Id
		tw.OriginTw = originTw
		tw.OriginUserId = originTw.UserId

		if originTw.OriginTwId != "" {
			return nil, errors.New("You can't forward Because it's not original")
		}
		originUser := &models.User{}
		if tx.Where("id = ?", originTw.UserId).Find(originUser).RowsAffected == 0 {
			return nil, errors.New("invalid origin user")
		}
		//因为是转发 所以这里也广播一下原文
		logs.PrintlnInfo("forward broadcast origin")
		go BroadcastTweet(originTw)
	}

	tw.Nonce = Nonce
	if createdAt > 0 {
		lastTw := &models.Tweets{}
		if tx.Where(
			"user_id = ? and nonce = ?", tw.UserId, tw.Nonce-1).Find(&lastTw).RowsAffected > 0 {
			if lastTw.CreatedAt >= createdAt {
				return nil, fmt.Errorf("createdAt must > %d", lastTw.CreatedAt)
			}
		}
		tw.CreatedAt = createdAt
	} else {
		tw.CreatedAt = time.Now().Unix()
	}
	msg := tw.GenerateSignMsg()
	var err error
	tw.Sign, err = models.GetCurrentUser().UsrNode.SignMsg(keyName, msg)
	if err != nil {
		return nil, err
	}
	if err = tw.Create(tx); err != nil {
		var i int64
		//nonce已存在
		if tx.Where("user_id = ? and nonce = ?", tw.UserId, tw.Nonce).Count(&i); i > 0 && tw.Nonce > user.Nonce {
			logs.PrintlnWarning("tw nonce already exists ", tw.Nonce, " update user nonce")
			uNonce = int64(tw.Nonce)
		}
		return nil, fmt.Errorf("create tweet err %s", err.Error())
	}
	if err = tx.Model(user).Save(user).Error; err != nil {
		return nil, fmt.Errorf("update user info err %s", err.Error())
	}
	isOk = true
	go func() {
		logs.PrintlnInfo("start broadcast tweet")
		BroadcastTweet(tw)
		logs.PrintlnSuccess("broadcast tweet success!", tw.Content)
	}()
	return tw, nil
}

func BroadcastTweetSync(tw *models.Tweets) {
	usr := models.GetCurrentUser()
	logs.PrintlnInfo("Wait nodes connect ..........")
	<-usr.UsrNode.WaitOnlineNode()
	logs.PrintlnInfo("Nodes connect ok........")
	if tw.UserInfo == nil {
		tw.UserInfo = &models.User{}
		global.GetDB().Where("id", tw.UserId).First(tw.UserInfo)
		tw.UserInfo = tw.UserInfo.GetUserInfoToPublic()
	}
	usr.UsrNode.EachOnlineNodes(func(nlo *p2pNet.OnlineNode) bool {
		if tw.UserInfo != nil && nlo.Pi.ID.String() == tw.UserInfo.PeerId {
			logs.PrintlnInfo("Skip author.................")
			return true
		}
		err := nlo.WriteData(&TweetInfo{Tw: tw})
		if err != nil {
			logs.PrintErr("WriteData err ", err, tw.Content, nlo.Pi.ID.String())
			logs.PrintlnWarning("Remove peer ", nlo.Pi.ID.String())
			usr.UsrNode.RemoveOnlineNodeLocked(nlo.Pi.ID.String())
			return true
		}
		logs.PrintlnSuccess("Send tw success ", tw.Content, nlo.Pi.ID.String())
		return true
	})
}

func BroadcastTweet(tw *models.Tweets) {
	_, err := models.AddUpIpfsAndBroadcastTweetTask(tw, 0)
	if err != nil {
		logs.PrintlnWarning("AddUpIpfsAndBroadcastTweetTask err ", err)
		logs.PrintlnInfo("Start Broadcast tweet sync... ", tw.Id)
		BroadcastTweetSync(tw)
	}
}

var Filter = bloom.NewWithEstimates(1000000, 0.01)

func (twInfo *TweetInfo) ReceiveHandle(ctx context.Context, node *p2pNet.OnlineNode) {
	twReceiveMu.Lock()
	defer twReceiveMu.Unlock()
	//存储
	if twInfo.Tw.UserInfo == nil {
		logs.PrintErr("not user info，refuse...", twInfo.Tw.Content)
		return
	}
	if twInfo.Tw.UserInfo.PubKey == "" {
		logs.PrintErr("public key is empty，refuse...", twInfo.Tw.Content)
		return
	}
	if twInfo.Tw.UserInfo.Id != twInfo.Tw.UserId {
		logs.PrintErr("userId atypism", twInfo.Tw.UserId, twInfo.Tw.UserInfo.Id)
		return
	}
	/*uj, _ := json.Marshal(twInfo.Tw.UserInfo)
	logs.PrintlnSuccess("============= 收到广播 ", string(uj), node.GetIdPretty())*/
	if !keys.VerifySignatureByAddress(twInfo.Tw.UserInfo.Id, twInfo.Tw.UserInfo.Sign, twInfo.Tw.UserInfo.GetSignMsg()) {
		logs.PrintErr("verify msg user info err, refuse...", twInfo.Tw.Content, "user id => ",
			twInfo.Tw.UserId,
			twInfo.Tw.UserInfo.Sign,
			twInfo.Tw.UserInfo.GetSignMsg(),
			node.GetIdPretty(),
		)
		return
	}
	if !keys.VerifySignatureByAddress(twInfo.Tw.UserId, twInfo.Tw.Sign, twInfo.Tw.GetSignMsg()) {
		logs.PrintErr("verify msg err, refuse...", twInfo.Tw.Content, "user id => ", twInfo.Tw.UserId, node.GetIdPretty())
		return
	}
	isOk := false
	sdb := global.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil || !isOk {
			sdb.Rollback()
			logs.PrintErr(r)
		} else {
			sdb.Commit()
		}
	}()

	usr := &models.User{}
	//查询对应的用户信息
	twUsrRowsAffected := global.LockForUpdate(sdb.Model(usr).Where("id = ?", twInfo.Tw.UserId)).Find(usr).RowsAffected
	if twUsrRowsAffected == 0 {
		usr = twInfo.Tw.UserInfo
		usr.LocalUser = global.IsNo
		usr.LocalNonce = 0
		usr.LastCheckTweetTime = 0
		logs.PrintlnInfo("new User info...")
		if err := sdb.Create(usr).Error; err != nil {
			logs.PrintErr(err)
			return
		}
	}
	var i int64
	if sdb.Table(twInfo.Tw.TableName()).Where("id = ?", twInfo.Tw.Id).Limit(1).Count(&i); i > 0 {
		logs.PrintlnInfo("data is saved，no more forwarding", twInfo.Tw.Content)
		return
	}

	//查询如果找到推文本地数据的推文nonce比接收推文的大但是创建时间却在接收推文之前 则要求广播推文方重置他的本地推文数据
	if sdb.Table(twInfo.Tw.TableName()).Where("user_id = ? and nonce > ? and created_at < ?", twInfo.Tw.UserId, twInfo.Tw.Nonce, twInfo.Tw.CreatedAt).Count(&i); i > 0 {
		logs.PrintlnWarning("need goto ask.........")
		go func() {
			//通知用户询问并且重置本地nonce和tweet数据
			uak := NewUserInfoGotoAsk(usr, true)
			_ = node.WriteData(uak)
		}()
		return
	}

	//更新用户信息
	if twUsrRowsAffected != 0 {
		usr.Name = twInfo.Tw.UserInfo.Name
		usr.Desc = twInfo.Tw.UserInfo.Desc
		usr.Avatar = twInfo.Tw.UserInfo.Avatar
		usr.Sign = twInfo.Tw.UserInfo.Sign
		usr.PeerId = twInfo.Tw.UserInfo.PeerId
		usr.PubKey = twInfo.Tw.UserInfo.PubKey
		usr.UpdatedSignUnix = twInfo.Tw.UserInfo.UpdatedSignUnix
		if twInfo.Tw.UserInfo.Nonce > usr.Nonce || (usr.Nonce == 0 && usr.LatestCid == "") {
			usr.Nonce = twInfo.Tw.UserInfo.Nonce
			usr.LatestCid = twInfo.Tw.UserInfo.LatestCid
		}
		logs.PrintlnInfo("update User info...")
		sdb.Save(usr)
	}

	if usr.Nonce != twInfo.Tw.Nonce {
		logs.PrintlnWarning("tweet nonce must ", usr.Nonce, " but not is", twInfo.Tw.Nonce, " refuse...", twInfo.Tw.Content)
		//更新计数器
		if twInfo.Tw.Nonce > usr.Nonce {
			usr.Nonce = twInfo.Tw.Nonce
			sdb.Select("Nonce").Save(usr)
		}
	}
	if err := sdb.Create(twInfo.Tw).Error; err != nil {
		logs.PrintErr("tweets save err ", twInfo.Tw.Content)
		twLocal := &models.Tweets{}
		if sdb.Where("user_id = ? and nonce = ?", twInfo.Tw.UserId, twInfo.Tw.Nonce).Limit(1).Find(twLocal).RowsAffected > 0 && twLocal.Id != twInfo.Tw.Id {
			go func() {
				//通知用户询问并且重置本地nonce和tweet数据
				uak := NewUserInfoGotoAsk(usr, true)
				_ = node.WriteData(uak)
			}()
		}
		return
	} else {
		logs.PrintlnSuccess("tweet save success：", twInfo.Tw.Content, twInfo.Tw.UserInfo.Id)
	}
	isOk = true
	//延迟广播
	go func() {
		//一条推文只广播一次
		if Filter.TestAndAddString(twInfo.Tw.Id) {
			return
		}
		time.Sleep(time.Second * 60)
		logs.PrintlnInfo("Start broadcast....", twInfo.Tw.Id)
		BroadcastTweet(twInfo.Tw)
		logs.PrintlnSuccess("Broadcast ok....", twInfo.Tw.Id)
	}()
}

func BroadcastNewestTweet(usr *models.User) {
	if usr.Nonce == 0 && usr.LatestCid == "" {
		return
	}
	sdb := global.GetDB()
	tweet := &models.Tweets{}
	rowsAffected := sdb.Model(models.Tweets{}).
		Where("user_id = ? ", usr.Id).
		Where("nonce = ? ", usr.Nonce).
		Limit(1).
		Find(tweet).RowsAffected
	if rowsAffected == 0 {
		logs.PrintErr("get gt nonce tweets fail, not found data")
		return
	}
	if tweet.Sign != "" {
		BroadcastTweet(tweet)
	}
}
