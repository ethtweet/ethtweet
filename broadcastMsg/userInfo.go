package broadcastMsg

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models"
	"github.com/ethtweet/ethtweet/p2pNet"
	"time"
)

const (
	UserInfoReceiveHandleTypeUpdate   = 1
	UserInfoReceiveHandleTypeAsk      = 2
	UserInfoReceiveHandleTypeAskReply = 3
	UserInfoReceiveHandleTypeGotoAsk  = 4
)

//该结构保存用户可广播的资料
type UserInfo struct {
	Id                        string
	Name                      string
	Desc                      string
	LatestCid                 string
	Avatar                    string
	Nonce                     uint64
	PeerId                    string
	HasPeerId                 uint8
	Sign                      string
	PublicKey                 string
	IpfsHash                  string
	UpdatedSignUnix           int64
	CreatedAt                 int64
	ReceiveHandleType         uint8
	IsRemoveNonceBeforeTweets bool

	user *models.User
}

func NewUserInfo(user *models.User) *UserInfo {
	return &UserInfo{
		Id:                        user.Id,
		Name:                      user.Name,
		Desc:                      user.Desc,
		LatestCid:                 user.LatestCid,
		Avatar:                    user.Avatar,
		Nonce:                     user.Nonce,
		PeerId:                    user.PeerId,
		HasPeerId:                 user.HasPeerId,
		UpdatedSignUnix:           user.UpdatedSignUnix,
		PublicKey:                 user.PubKey,
		Sign:                      user.Sign,
		IpfsHash:                  user.IpfsHash,
		CreatedAt:                 user.CreatedAt,
		ReceiveHandleType:         UserInfoReceiveHandleTypeUpdate,
		IsRemoveNonceBeforeTweets: false,
		user:                      user,
	}
}

func NewUserInfoAsk(user *models.User) *UserInfo {
	return &UserInfo{
		Id:                        user.Id,
		Name:                      user.Name,
		Desc:                      user.Desc,
		LatestCid:                 user.LatestCid,
		Avatar:                    user.Avatar,
		Nonce:                     user.Nonce,
		PeerId:                    user.PeerId,
		HasPeerId:                 user.HasPeerId,
		UpdatedSignUnix:           user.UpdatedSignUnix,
		PublicKey:                 user.PubKey,
		Sign:                      user.Sign,
		IpfsHash:                  user.IpfsHash,
		CreatedAt:                 user.CreatedAt,
		ReceiveHandleType:         UserInfoReceiveHandleTypeAsk,
		IsRemoveNonceBeforeTweets: false,
		user:                      user,
	}
}

func NewUserInfoGotoAsk(user *models.User, isRemoveTw bool) *UserInfo {
	return &UserInfo{
		Id:                        user.Id,
		Name:                      user.Name,
		Desc:                      user.Desc,
		LatestCid:                 user.LatestCid,
		Avatar:                    user.Avatar,
		Nonce:                     user.Nonce,
		PeerId:                    user.PeerId,
		HasPeerId:                 user.HasPeerId,
		UpdatedSignUnix:           user.UpdatedSignUnix,
		PublicKey:                 user.PubKey,
		Sign:                      user.Sign,
		IpfsHash:                  user.IpfsHash,
		CreatedAt:                 user.CreatedAt,
		ReceiveHandleType:         UserInfoReceiveHandleTypeGotoAsk,
		IsRemoveNonceBeforeTweets: isRemoveTw,
		user:                      user,
	}
}

func (usrInfo *UserInfo) ReceiveHandle(ctx context.Context, node *p2pNet.OnlineNode) {
	if usrInfo.ReceiveHandleType == UserInfoReceiveHandleTypeAsk {
		usrInfo.ReceiveHandleAsk(ctx, node)
		return
	} else if usrInfo.ReceiveHandleType == UserInfoReceiveHandleTypeGotoAsk { //需要用户去询问
		usrInfo.ReceiveHandleType = UserInfoReceiveHandleTypeAsk
		err := usrInfo.DoAsk()
		if err != nil {
			logs.PrintErr("UserInfoReceiveHandleTypeGotoAsk err ", err)
		}
		return
	}
	usrInfo.ReceiveHandleUpdate(ctx, node)
	return
}

//接收方法
func (usrInfo *UserInfo) ReceiveHandleUpdate(ctx context.Context, node *p2pNet.OnlineNode) {
	//这里只处理更新或者回复询问的数据
	if usrInfo.ReceiveHandleType != UserInfoReceiveHandleTypeAskReply && usrInfo.ReceiveHandleType != UserInfoReceiveHandleTypeUpdate {
		logs.PrintlnWarning("invalid ReceiveHandleType")
		return
	}
	logs.PrintlnInfo("receive update user info request....", usrInfo.Id, usrInfo)
	usr := &models.User{}
	db := global.GetDB()
	if db.Where("id", usrInfo.Id).Find(usr).RowsAffected == 0 {
		usr.Id = usrInfo.Id
		usr.LocalUser = global.IsNo
		usr.CreatedAt = usr.CreatedAt
		usr.PubKey = usrInfo.PublicKey
	}
	//更新用户 如果是其他节点回复询问的话 不要验证signUnix
	if usrInfo.ReceiveHandleType != UserInfoReceiveHandleTypeAskReply && usr.UpdatedSignUnix >= usrInfo.UpdatedSignUnix {
		logs.PrintlnWarning(fmt.Sprintf("local user UpdatedSignUnix %d >= receive user UpdatedSignUnix %d", usr.UpdatedSignUnix, usrInfo.UpdatedSignUnix))
		return
	}

	//回复的nonce小于本地的则不作处理
	if usr.Nonce > usrInfo.Nonce {
		logs.PrintlnWarning(fmt.Sprintf("refuse update......... local user nonce %d > receive user nonce %d", usr.Nonce, usrInfo.Nonce))
		return
	} else if usrInfo.LatestCid == "" && usr.LatestCid != "" {
		logs.PrintlnWarning("usrInfo.LatestCid is empty, but usr.LatestCid is not empty, so not update...")
	}
	//不是本地用户可更新节点id
	if usr.LocalUser == global.IsNo && usr.HasPeerId == global.IsNo {
		usr.PeerId = usrInfo.PeerId
		usr.HasPeerId = usrInfo.HasPeerId
	}

	usr.Name = usrInfo.Name
	usr.Desc = usrInfo.Desc
	usr.Avatar = usrInfo.Avatar
	usr.UpdatedSignUnix = usrInfo.UpdatedSignUnix
	usr.Sign = usrInfo.Sign
	usr.UpdatedAt = time.Now().Unix()

	usr.Nonce = usrInfo.Nonce
	usr.LatestCid = usrInfo.LatestCid
	usr.IpfsHash = usrInfo.IpfsHash
	//验证签名
	if !keys.VerifySignatureByAddress(usrInfo.Id, usrInfo.Sign, usr.GetSignMsg()) {
		logs.PrintErr(fmt.Sprintf("receive update user info but sign fail; signMsg:%s, userId:%s, sign:%s, peerId:%s, ReceiveHandleType:%d, address:%s",
			usr.GetSignMsg(),
			usrInfo.Id,
			usrInfo.Sign,
			node.GetIdPretty(),
			usrInfo.ReceiveHandleType,
			node.Stream.Conn().RemoteMultiaddr().String()))
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		//如果需要则删除nonce之前的所有推文 并且设置localNonce为0 等待重新同步
		if usrInfo.IsRemoveNonceBeforeTweets {
			tx.Where("user_id = ? and nonce <= ?", usr.Id, usr.Nonce).Delete(&models.Tweets{})
			usr.LocalNonce = 0
		}

		if err := tx.Session(&gorm.Session{SkipHooks: true}).Save(usr).Error; err != nil {
			return fmt.Errorf("update user fail ", err)
		} else {
			_ = usr.SetSyncStatusComplete(tx)
			logs.PrintlnSuccess("update user success")
		}
		return nil
	})
	if err != nil {
		logs.PrintErr(err)
		return
	}
	if usrInfo.IsRemoveNonceBeforeTweets {
		cu := models.GetCurrentUser()
		if cu == nil {
			return
		}
		go func() {
			AskUsersTweets([]*models.User{usr}, cu.UsrNode, global.GetGlobalCtx())
		}()
	}
}

func (usrInfo *UserInfo) DoAsk() (err error) {
	if usrInfo.ReceiveHandleType != UserInfoReceiveHandleTypeAsk {
		return fmt.Errorf("invalid ReceiveHandleType")
	}
	usr := models.GetCurrentUser()
	//在节点同步的时间+2s 也就是说再一次同步在线节点的时候任然没有找到其他的节点即超时
	re := 0
	waitOnlineNodeNum := 5
	defer func() {
		if err != nil {
			if errors.Is(err, global.ErrUserAskTimeout) {
				//主要是检查网络是否通顺
				_, err2 := global.GetIpfsInfo(usr.IpfsHash)
				if err2 != nil {
					logs.PrintErr("do ask network err ", err)
					return
				}
				if usrInfo.user == nil {
					usrInfo.user = &models.User{}
					if global.GetDB().Where("id", usrInfo.Id).First(usrInfo.user).RowsAffected == 0 {
						return
					}
				}
				_ = usrInfo.user.SetSyncStatusComplete(nil)
			}
		}
	}()
RE:
	t := time.NewTimer(p2pNet.OnlineNodesSyncDuration + 2*time.Second)
	select {
	case <-t.C:
		t.Stop()
		logs.PrintlnWarning("wait node timeout ............... ")
		if waitOnlineNodeNum != 0 {
			waitOnlineNodeNum = 0
			goto RE
		}
		err = global.ErrUserAskTimeout
		return err
	case <-usr.UsrNode.WaitOnlineNodeLimitNum(waitOnlineNodeNum):
		logs.PrintlnInfo("wait online node is ok..........")
		t.Stop()
	}
	hasWriteOk := false
	usr.UsrNode.EachOnlineNodes(func(node *p2pNet.OnlineNode) bool {
		logs.PrintlnInfo("broadcast update info ask req to ", node.Pi.ID)
		err := p2pNet.WriteData(node.Rw, usrInfo)
		if err != nil {
			logs.PrintErr("broadcast update info ask err ", node.Pi.ID, err)
		}
		hasWriteOk = true
		return true
	})
	if !hasWriteOk {
		if re == 0 {
			re++
			goto RE
		}
		err = global.ErrUserAskWriteAllFail
		return err
	}
	if usrInfo.user == nil {
		usrInfo.user = &models.User{}
		if global.GetDB().Where("id", usrInfo.Id).First(usrInfo.user).RowsAffected == 0 {
			return nil
		}
	}
	err = usrInfo.user.SetSyncStatusIng()
	if err != nil {
		logs.PrintlnWarning("SetSyncStatusIng err ", err)
	}
	return nil
}

func (usrInfo *UserInfo) ReceiveHandleAsk(ctx context.Context, node *p2pNet.OnlineNode) {
	if usrInfo.ReceiveHandleType != UserInfoReceiveHandleTypeAsk {
		logs.PrintlnWarning("invalid ReceiveHandleType")
		return
	}
	logs.PrintlnInfo("receive ask user info ", usrInfo.Id)
	localUser := &models.User{}
	if global.GetDB().Where("id", usrInfo.Id).Limit(1).Find(localUser).RowsAffected == 0 {
		logs.PrintlnWarning("receive ask user info, but not found user ", usrInfo.Id)
		return
	}
	if localUser.Sign == "" {
		logs.PrintlnWarning("sign is null ", localUser.Id)
		return
	}
	userInfo := NewUserInfo(localUser)
	userInfo.IsRemoveNonceBeforeTweets = usrInfo.IsRemoveNonceBeforeTweets
	userInfo.ReceiveHandleType = UserInfoReceiveHandleTypeAskReply
	err := node.WriteData(userInfo)
	if err != nil {
		logs.PrintErr("reply user ask fail ", err)
	} else {
		logs.PrintlnSuccess("reply user ask success")
	}
}

func SyncUserInfo(usr *models.User, force bool) error {
	if !force && !usr.SyncStatusWait() {
		return fmt.Errorf("syncing ....")
	}
	logs.PrintlnInfo("start sync user info ", usr.Id)
	usrInfo := NewUserInfoAsk(usr)
	defer logs.PrintlnSuccess("start sync user end ", usr.Id)
	return usrInfo.DoAsk()
}
