package models

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/ethtweet/ethtweet/models/mField"
	"github.com/ethtweet/ethtweet/p2pNet"
	"time"

	cryptoEth "github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	keystore "github.com/ipfs/go-ipfs-keystore"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/mr-tron/base58"
	"gorm.io/gorm"
)

const (
	UserSyncStatusWait     = 0
	UserSyncStatusIng      = 1
	UserSyncStatusComplete = 2

	UserSyncTimeoutDuration = p2pNet.OnlineNodesSyncDuration * 4
)

var currentUsr *User

type User struct {
	Id                         string `gorm:"primarykey;type:varchar(100);NOT NULL;comment:用户id 保存为以太坊的地址"`
	PeerId                     string `gorm:"type:varchar(100);NOT NULL;index:idx_peer_id,unique;comment:节点id 可为空"`
	Name                       string `gorm:"type:varchar(100);NOT NULL;comment:用户昵称;index:idx_user_name"`
	Desc                       string `gorm:"type:varchar(200);comment:用户简介;default:''"`
	LatestCid                  string `gorm:"type:varchar(200);comment:最新一条cid;default:''"`
	Avatar                     string `gorm:"type:text;comment:用户简介;"`
	Nonce                      uint64 `gorm:"default:0;comment:最新nonce"`
	LocalNonce                 uint64 `gorm:"default:0;comment:本地nonce"`
	LocalUser                  uint64 `gorm:"default:0;comment:是否是本地用户"`
	Sign                       string `gorm:"type:varchar(200);NOT NULL;comment:用户签名;"`
	PubKey                     string `gorm:"type:text;comment:用户公钥"`
	IpfsHash                   string `gorm:"type:varchar(200);comment:ipfs的hash地址;default:''"`
	HasPeerId                  uint8  `gorm:"default:0;index:idx_has_peer;comment:节点id是否是有效的id 通过中心化接口创建的用户 因为没有运行节点 所以不会存在peerId 程序将为其生成一个唯一的但不是有效的节点id作为替代"`
	LastCheckTweetTime         int64  `gorm:"default:0;index:last_check_tweet_time"`
	mField.FieldsTimeUnixModel `json:"-"`
	mField.FieldsExtendsJsonType
	UpdatedSignUnix int64            `gorm:"default:0;comment:更新时间戳，用于签名使用"`
	UsrNode         *p2pNet.UserNode `gorm:"-" json:"-"`
}

func (usr *User) TableName() string {
	return "user"
}

func (usr *User) AfterCreate(tx *gorm.DB) (err error) {
	logs.PrintlnInfo("AfterCreate...................... ", usr.Id)
	err = usr.UploadIpfs(tx)
	if err != nil {
		logs.PrintlnWarning("upload ipfs err ", err)
	}
	return nil
}

func (usr *User) BeforeUpdate(tx *gorm.DB) (err error) {
	logs.PrintlnInfo("BeforeUpdate...................... ", usr.Id)
	err = usr.UploadIpfs(tx)
	if err != nil {
		logs.PrintlnWarning("upload ipfs err ", err)
	}
	return nil
}

func (usr *User) SyncStatusWait() bool {
	uas := &UserAskSync{}
	if global.GetDB().Where("user_id", usr.Id).Limit(1).Find(uas).RowsAffected == 0 {
		return true
	}
	return uas.SyncStatus == UserSyncStatusWait
}

func (usr *User) SyncStatusIng() bool {
	uas := &UserAskSync{}
	if global.GetDB().Where("user_id", usr.Id).Limit(1).Find(uas).RowsAffected == 0 {
		return false
	}
	if uas.SyncStatus == UserSyncStatusIng {
		if uas.SyncTimeoutUnix <= time.Now().Unix() {
			uas.SyncStatus = UserSyncStatusComplete
			uas.SyncTimeoutUnix = 0
			global.GetDB().Model(uas).Where("user_id", usr.Id).Update("SyncStatus", UserSyncStatusComplete)
			return false
		}
		return true
	}
	return false
}

func (usr *User) SyncStatusComplete() bool {
	uas := &UserAskSync{}
	if global.GetDB().Where("user_id", usr.Id).Limit(1).Find(uas).RowsAffected == 0 {
		return false
	}
	if uas.SyncStatus == UserSyncStatusIng && uas.SyncTimeoutUnix < time.Now().Unix() {
		uas.SyncStatus = UserSyncStatusComplete
		uas.SyncTimeoutUnix = 0
		global.GetDB().Save(uas)
		return true
	}
	return uas.SyncStatus == UserSyncStatusComplete
}

func (usr *User) SetSyncStatusIng() error {
	if usr.SyncStatusIng() {
		return fmt.Errorf("sync status is ing")
	}
	if usr.SyncStatusComplete() {
		return fmt.Errorf("sync status is complete")
	}
	uas := &UserAskSync{}
	if global.GetDB().Where("user_id", usr.Id).Limit(1).Find(uas).RowsAffected == 0 {
		uas.SyncStatus = UserSyncStatusIng
		uas.SyncTimeoutUnix = time.Now().Add(UserSyncTimeoutDuration).Unix()
		uas.UserId = usr.Id
		return global.GetDB().Create(uas).Error
	}
	uas.SyncStatus = UserSyncStatusIng
	uas.SyncTimeoutUnix = time.Now().Add(UserSyncTimeoutDuration).Unix()
	return global.GetDB().Where("sync_status", UserSyncStatusWait).Select("SyncStatus", "SyncTimeoutUnix").Save(uas).Error
}

// SetSyncStatusComplete 设置成功同步状态
func (usr *User) SetSyncStatusComplete(tx *gorm.DB) error {
	if usr.SyncStatusComplete() {
		return fmt.Errorf("sync status is ing")
	}
	if tx == nil {
		tx = global.GetDB()
	}
	uas := &UserAskSync{}
	if tx.Where("user_id", usr.Id).Limit(1).Find(uas).RowsAffected == 0 {
		uas.SyncStatus = UserSyncStatusComplete
		uas.SyncTimeoutUnix = 0
		uas.UserId = usr.Id
		return tx.Create(uas).Error
	}
	uas.SyncStatus = UserSyncStatusComplete
	uas.SyncTimeoutUnix = 0
	return tx.Select("SyncStatus", "SyncTimeoutUnix").Save(uas).Error
}

func SetCurrentUser(usrNode *p2pNet.UserNode) error {
	usr := &User{}
	sdb := global.GetDB()
	priKey, err := usrNode.GetPriKey("")
	if err != nil {
		return err
	}
	if err := sdb.Model(usr).Where("id = ?", priKey.GetEthAddress()).Find(usr).Error; err != nil {
		return err
	}
	defer func() {
		if currentUsr != nil {
			currentUsr.UsrNode = usrNode
			if usr.Sign == "" {
				_ = usr.ReloadSign(true)
			}
		}
	}()
	if usr.Id != "" {
		if usr.PubKey == "" {
			logs.PrintlnInfo("empty pubKey, creating.....")
			usr.PubKey = priKey.Encode58Public()
			if err := sdb.Save(usr).Error; err != nil {
				return err
			}
		}
		usr.PeerId = usrNode.IDPretty()
		usr.HasPeerId = global.IsYes
		if err := sdb.Save(usr).Error; err != nil {
			return err
		}
		currentUsr = usr
		return nil
	}

	//新建用户
	usr, err = GetOrCreateUserByPri(priKey)
	if err != nil {
		return err
	}
	currentUsr = usr
	return nil
}

func ClearCurrentUser() {
	currentUsr = nil
}

func GetCurrentUser() *User {
	return currentUsr
}

func GetUserByIpfs(h string) (*User, error) {
	r, err := global.GetIpfsInfo(h)
	if err != nil {
		return nil, err
	}
	user := &User{}
	err = json.Unmarshal(r, user)
	if err != nil {
		return nil, err
	}
	user.IpfsHash = h
	user.LocalNonce = 0
	user.LocalUser = global.IsNo
	return user, nil
}

// 验证签名时调用
func (usr *User) GetSignMsg() string {
	return fmt.Sprintf("%s|%s|%s|%s|%d", usr.Id, usr.Name, usr.Desc, usr.Avatar, usr.UpdatedSignUnix)
}

// 主要用户生成签名
func (usr *User) GenerateSignMsg() string {
	usr.UpdatedSignUnix = time.Now().Unix()
	return usr.GetSignMsg()
}

func (usr *User) ReloadSign(isSave bool) error {
	msg := usr.GenerateSignMsg()
	if usr.UsrNode == nil {
		return fmt.Errorf("userNode is invalid")
	}
	var err error
	usr.Sign, err = usr.UsrNode.SignMsg("", msg)
	if err != nil {
		return err
	}
	if isSave {
		return global.GetDB().Save(usr).Error
	}
	return nil
}

// 返回用户可以公开的字段信息
func (usr *User) GetUserInfoToPublic() *User {
	u := &User{
		Id:              usr.Id,
		Name:            usr.Name,
		Desc:            usr.Desc,
		Nonce:           usr.Nonce,
		LatestCid:       usr.LatestCid,
		Avatar:          usr.Avatar,
		Sign:            usr.Sign,
		PubKey:          usr.PubKey,
		PeerId:          usr.PeerId,
		HasPeerId:       usr.HasPeerId,
		IpfsHash:        usr.IpfsHash,
		UpdatedSignUnix: usr.UpdatedSignUnix,
	}
	u.CreatedAt = usr.CreatedAt
	u.UpdatedAt = usr.UpdatedAt
	return u
}

func (usr *User) UploadIpfs(tx *gorm.DB) error {
	h, err := global.UploadIpfs(usr.GetUserInfoToPublic())
	if err != nil {
		return fmt.Errorf("ipfs upload err %s", err.Error())
	}
	usr.IpfsHash = h
	if tx == nil {
		return nil
	}
	return tx.Table(usr.TableName()).Where("id = ?", usr.Id).Updates(map[string]interface{}{
		"ipfs_hash":  usr.IpfsHash,
		"updated_at": time.Now().Unix(),
	}).Error
}

// 加载一个临时的节点id
func (usr *User) ReloadTmpPeerId() {
	_uuid := uuid.New()
	usr.HasPeerId = global.IsNo
	usr.PeerId = _uuid.String()
}

func GetOrCreateUserByPri(pri *keys.PrivateKey) (*User, error) {
	peerId, err := peer.IDFromPublicKey(pri.LibP2pPrivate.GetPublic())
	if err != nil {
		return nil, err
	}
	address := pri.GetEthAddress()
	usr := &User{}
	if global.GetDB().Where("id = ?", address.String()).Find(usr).RowsAffected > 0 {
		return usr, nil
	}
	usr.Id = address.String()
	usr.PeerId = peerId.String()
	usr.HasPeerId = global.IsYes
	usr.Name = global.RandStringRunes(6)
	usr.PubKey = pri.Encode58Public()
	usr.LocalUser = global.IsYes
	usr.CreatedAt = time.Now().Unix()
	_ = usr.ReloadSign(false)
	if err := global.GetDB().Create(usr).Error; err != nil {
		return nil, err
	}
	return usr, nil
}

// 通过公钥创建一个用户 生成一个零食的peerId
func GetOrCreateUserByPub(pub *ecdsa.PublicKey) (*User, error) {
	address := cryptoEth.PubkeyToAddress(*pub)
	usr := &User{}
	if global.GetDB().Where("id = ?", address.String()).Find(usr).RowsAffected > 0 {
		return usr, nil
	}
	usr.Id = address.String()
	usr.Name = global.RandStringRunes(6)
	usr.PubKey = base58.Encode(cryptoEth.FromECDSAPub(pub))
	usr.LocalUser = 0
	usr.ReloadTmpPeerId()
	//_ = usr.ReloadSign(false)
	usr.Sign = ""
	if err := global.GetDB().Create(usr).Error; err != nil {
		return nil, err
	}
	return usr, nil
}

func GetUserByKeyName(dir, key string, autoCreate bool) (*User, error) {
	ks, err := keystore.NewFSKeystore(dir)
	if err != nil {
		return nil, err
	}
	if ok, _ := ks.Has(key); !ok {
		if autoCreate {
			priKey, err := keys.NewPrivateKey()
			if err != nil {
				return nil, err
			}
			err = ks.Put(key, priKey.LibP2pPrivate)
			if err != nil {
				return nil, err
			}
			return GetOrCreateUserByPri(priKey)
		}
		return nil, fmt.Errorf("invalid key")
	}
	pri, err := ks.Get(key)
	if err != nil {
		return nil, err
	}
	peerId, err := peer.IDFromPublicKey(pri.GetPublic())
	if err != nil {
		return nil, err
	}
	usr := &User{}
	if global.GetDB().Where("peer_id = ?", peerId.String()).Find(usr).RowsAffected > 0 {
		return usr, nil
	}
	if autoCreate {
		priKey, err := keys.NewPrivateKeyByLibP2pPri(pri)
		if err != nil {
			return nil, err
		}
		return GetOrCreateUserByPri(priKey)
	}
	return nil, fmt.Errorf("not found user")
}
