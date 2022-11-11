package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ethtweet/ethtweet/global"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/mr-tron/base58"
)

type TweetJson struct {
	PreviousCid string
	UserId      string
	Nonce       uint64
	Content     string
	Attachment  string
	Sign        string
	TopicTag    string
	OriginTwId  string
	OriginTw    *Tweets
	CreatedAt   int64
}

func NewTweetJson(tw *Tweets, latestCid string) *TweetJson {
	return &TweetJson{
		PreviousCid: latestCid,
		UserId:      tw.UserId,
		Content:     tw.Content,
		Attachment:  tw.Attachment,
		Nonce:       tw.Nonce,
		Sign:        tw.Sign,
		OriginTwId:  tw.OriginTwId,
		CreatedAt:   tw.CreatedAt,
		OriginTw:    tw.OriginTw,
		TopicTag:    tw.TopicTag,
	}
}

type Tweets struct {
	Id         string `gorm:"primarykey;type:varchar(100);NOT NULL;comment:sha2(用户地址+nonce) 然后进行base58编码;"`
	UserId     string `gorm:"type:varchar(100);NOT NULL;comment:用户的id，即用户地址;index:idx_uid_nonce,unique"`
	Content    string `gorm:"NOT NULL;comment:推文内容"`
	Attachment string `gorm:"comment:附件"`
	Nonce      uint64 `gorm:"default:0;comment:根据用户的推文自增;index:idx_uid_nonce,unique"`
	Sign       string `gorm:"type:varchar(200);NOT NULL;comment:签名(用户地址+nonce+推文)，然后base58编码;index:tweet_sign,unique"`
	TopicTag   string `gorm:"type:varchar(200);default:'';comment:话题标签;index:topic_tag"`

	UserInfo     *User  `gorm:"foreignKey:UserId;references:Id"` //用于发送消息时嵌入的用户信息 这个数据不保存在数据库中
	CreatedAt    int64  `gorm:"index:idx_created_at;autoCreateTime;comment:推文的发布时间"`
	UpdatedAt    int64  `gorm:"autoUpdateTime;default:0;comment:更新时间 也就是同步到推文的时间"`
	OriginUserId string `gorm:"type:varchar(100);default:'';comment:作者id，如果没有被转发 则为空字符串;index:idx_origin_user_id"`
	OriginTwId   string `gorm:"varchar(100);comment:原始文章Id;index:idx_origin_tw_id"`

	OriginTw *Tweets `gorm:"foreignKey:OriginTwId;references:Id"`
}

func (tw *Tweets) BeforeCreate(tx *gorm.DB) (err error) {
	tw.UpdatedAt = time.Now().Unix()
	return nil
}

func (tw *Tweets) BeforeUpdate(tx *gorm.DB) (err error) {
	tw.UpdatedAt = time.Now().Unix()
	return nil
}

func (tw *Tweets) TableName() string {
	return "tweets"
}

func (tw *Tweets) UpIpfs(usr *User) (string, error) {
	if usr.Id != tw.UserId {
		return "", errors.New("invalid user")
	}
	cid, err := global.UploadIpfs(NewTweetJson(tw, usr.LatestCid))
	if err != nil {
		return "", fmt.Errorf("tweet upload ipfs err %s", err.Error())
	}
	fmt.Printf("publish to https://ipfs.io/ipfs/%s\n", cid)
	usr.LatestCid = cid
	return cid, nil
}

func (tw *Tweets) GenerateId() {
	tw.Id = fmt.Sprintf("%s|%d|%s", tw.UserId, tw.Nonce, tw.Content)
	tw.Id = fmt.Sprintf("%x", sha256.Sum256([]byte(tw.Id)))
	tw.Id = base58.Encode([]byte(tw.Id))
}

//生成twId sign参数
func (tw *Tweets) GenerateSysParams() {
	tw.GenerateId()
}

func (tw *Tweets) Create(tx *gorm.DB) error {
	if tx == nil {
		tx = global.GetDB()
	}
	if tw.UserId == "" {
		return errors.New("invalid uid")
	}
	if tw.Content == "" {
		return errors.New("content is not null")
	}
	tw.GenerateSysParams()
	if err := tx.Create(tw).Error; err != nil {
		return nil
	}
	return nil
}

func (tw *Tweets) GetSignMsg() string {
	return fmt.Sprintf("%s|%d|%s|%s|%s|%s|%d", strings.ToLower(tw.UserId), tw.Nonce, tw.Content, tw.Attachment, strings.ToLower(tw.OriginTwId), strings.ToLower(tw.OriginUserId), tw.CreatedAt)
}

func (tw *Tweets) GenerateSignMsg() string {
	if tw.CreatedAt == 0 {
		tw.CreatedAt = time.Now().Unix()
	}
	return tw.GetSignMsg()
}
