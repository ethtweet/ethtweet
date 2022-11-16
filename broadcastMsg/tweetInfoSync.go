package broadcastMsg

import (
	"context"
	"fmt"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/p2pNet"
	"time"
)

const (
	TweetInfoSyncTypeAsk   = 1
	TweetInfoSyncTypeReply = 2
)

type TweetInfoSync struct {
	UserAddress string
	TwNonce     uint64
	Size        int
	Type        uint8
	ReplyTweets []*models.Tweets //对于询问回复的推文信息

	//是否还有其他的新推文没有回复 该字段不代表已回复的推文列表是否满足size大小 而是是否还有比已回复的推文更新的内容
	//对于size是否完成 通过比较size和len(ReplyTweets)
	ReplyTwSurplusNewNum uint64
}

func NewTweetInfoSyncAsk(userAddress string, twNonce uint64) *TweetInfoSync {
	return &TweetInfoSync{
		UserAddress:          userAddress,
		TwNonce:              twNonce,
		Type:                 TweetInfoSyncTypeAsk,
		Size:                 50,
		ReplyTweets:          nil,
		ReplyTwSurplusNewNum: 0,
	}
}

func (tis *TweetInfoSync) ReceiveHandle(ctx context.Context, node *p2pNet.OnlineNode) {
	if tis.Type == TweetInfoSyncTypeAsk {
		tis.ReceiveHandleAsk(ctx, node)
		return
	}
	tis.ReceiveHandleReply(ctx, node)
	return
}

func (tis *TweetInfoSync) ReceiveHandleAsk(ctx context.Context, node *p2pNet.OnlineNode) {
	if tis.Type != TweetInfoSyncTypeAsk {
		return
	}
	tws := make([]*models.Tweets, 0, tis.Size)
	//没有更多的推文
	query := global.GetDB()
	query = query.Where("nonce >= ?", tis.TwNonce)
	if query.
		Where("user_id", tis.UserAddress).
		Order("nonce asc").
		Limit(tis.Size).Find(&tws).RowsAffected == 0 {
		return
	}

	if len(tws) == 0 {
		logs.PrintlnInfo("ReceiveHandleAsk no new tweet")
		return
	}

	//这里处理一下必须是nonce连续的推文才放置到回复列表里去
	tis.ReplyTweets = make([]*models.Tweets, 0, len(tws))
	var maxNonce uint64 = 0
	for k, rTw := range tws {
		if k > 0 && rTw.Nonce != tws[k-1].Nonce+1 {
			break
		}
		tis.ReplyTweets = append(tis.ReplyTweets, rTw)
		maxNonce = rTw.Nonce
	}
	if len(tis.ReplyTweets) == 0 {
		return
	}
	//查询一下是否有更新的连续数据
	var i int64 = 0
	global.GetDB().Model(&models.Tweets{}).Where("user_id", tis.UserAddress).
		Where("nonce >= ?", maxNonce+1).Count(&i)
	tis.ReplyTwSurplusNewNum = uint64(i)

	//修改一下类型回复出去
	// todo 如果没有更新的，那就不回复
	tis.Type = TweetInfoSyncTypeReply
	logs.PrintlnInfo(fmt.Sprintf("ReceiveHandleAsk userId: %s, twNonce: %d, size: %d, ReplyTwSurplusNewNum: %d", tis.UserAddress, tis.TwNonce, tis.Size, tis.ReplyTwSurplusNewNum))
	err := node.WriteData(tis)
	if err != nil {
		logs.PrintlnWarning(fmt.Sprintf("reply user:%s, tweet err %s", tis.UserAddress, err.Error()))
	}
}

func (tis *TweetInfoSync) ReceiveHandleReply(ctx context.Context, node *p2pNet.OnlineNode) {
	if tis.Type != TweetInfoSyncTypeReply || len(tis.ReplyTweets) == 0 {
		return
	}
	tx := global.GetDB()
	user := &models.User{}
	if tx.Where("id", tis.UserAddress).Find(user).RowsAffected == 0 {
		return
	}
	var i int64 = 0
	var maxNonce uint64 = 0
	var isContinueAsk = true
	var isUpdateLocalNonce = true
	for k, rtw := range tis.ReplyTweets {
		if tx.Table(rtw.TableName()).Where("id", rtw.Id).Count(&i); i > 0 {
			logs.PrintlnInfo("tweet id exist...")
			if rtw.Nonce > maxNonce {
				maxNonce = rtw.Nonce
			}
			continue
		}
		//验证签名
		if !keys.VerifySignatureByAddress(user.Id, rtw.Sign, rtw.GetSignMsg()) {
			isContinueAsk = false
			logs.PrintErr("verify msg err, refuse...", rtw.Content, "user id => ", rtw.UserId)
			break
		}
		//连续性验证
		if k == 0 {
			if rtw.Nonce != user.LocalNonce+1 && rtw.Nonce != user.LocalNonce {
				isUpdateLocalNonce = false
			}
		} else {
			if rtw.Nonce != tis.ReplyTweets[k-1].Nonce+1 {
				logs.PrintlnWarning("Not continuous data ", tis.ReplyTweets[k-1].Nonce, "=====", rtw.Nonce)
				break
			}
		}
		if err := tx.Create(rtw).Error; err != nil {
			isContinueAsk = false
			logs.PrintErr(fmt.Sprintf("create user %s tweet err %s", rtw.UserId, err.Error()))
			break
		}
		if rtw.Nonce > maxNonce {
			maxNonce = rtw.Nonce
		}
	}
	if isContinueAsk && tis.ReplyTwSurplusNewNum > 0 {
		logs.PrintlnInfo("start two ask err")
		err := node.WriteData(NewTweetInfoSyncAsk(tis.UserAddress, maxNonce+1))
		if err != nil {
			logs.PrintErr("two ask err ", err)
		} else {
			logs.PrintlnSuccess("two ask success")
		}
	}
	if isUpdateLocalNonce && maxNonce > user.LocalNonce {
		tx.Model(user).Where("local_nonce < ?", maxNonce).Update("LocalNonce", maxNonce)
	}
	return
}

// SyncUserTweets 循环拉去网络最新推文
func SyncUserTweets(ctx context.Context) error {
	cUser := models.GetCurrentUser()
	if cUser == nil || cUser.UsrNode == nil {
		return fmt.Errorf("user node is invalid")
	}
	t := time.NewTicker(45 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
		}
		logs.PrintlnInfo("Start sync tweets.................")
		users := make([]*models.User, 0, 20)
		if global.GetDB().Limit(20).
			Where("last_check_tweet_time < ?", time.Now().Add(-(time.Minute*60)).Unix()). //更新频率控制
			Where("local_nonce < nonce or (local_nonce = 0 and latest_cid != '')").Order("last_check_tweet_time asc").
			Find(&users).RowsAffected > 0 {
			c, cc := context.WithTimeout(ctx, 10*time.Second)
			AskUsersTweets(users, cUser.UsrNode, c)
			cc()
		} else {
			logs.PrintlnInfo("Not need sync user tweets...")
		}
	}
}

// 询问单个用的推文
func AskUsersTweets(users []*models.User, node *p2pNet.UserNode, ctx context.Context) {
	logs.PrintlnInfo("AskUsersTweets wait online node............")
	select {
	case <-ctx.Done():
		return
	case <-node.WaitOnlineNode():
	}
	logs.PrintlnSuccess("AskUsersTweets wait online node ok............")
	tx := global.GetDB()
	for _, user := range users {
		logs.PrintlnInfo(fmt.Sprintf("prepare sync user's %s tweet.....", user.Id))
		if !user.SyncStatusComplete() {
			logs.PrintlnWarning(fmt.Sprintf("user %s is not SyncStatusComplete", user.Id))
			uak := NewUserInfoAsk(user)
			_ = uak.DoAsk()
			tx.Model(user).Update("LastCheckTweetTime", time.Now().Unix())
			continue
		}
		logs.PrintlnSuccess(fmt.Sprintf("Start sync user's %s tweet.....", user.Id))
		tis := NewTweetInfoSyncAsk(user.Id, user.LocalNonce)
		node.EachOnlineNodes(func(node *p2pNet.OnlineNode) bool {
			err := node.WriteData(tis)
			if err != nil {
				logs.PrintErr("Tis write err ", err)
			} else {
				logs.PrintlnSuccess("Tis success user id ", tis.UserAddress)
			}
			return true
		})
		user.LastCheckTweetTime = time.Now().Unix()
		tx.Save(user)
		logs.PrintlnSuccess(fmt.Sprintf("Success sync user's %s tweet.....", user.Id))
	}
}
