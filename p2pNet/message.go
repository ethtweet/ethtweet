package p2pNet

import (
	"bufio"
	"context"
	"encoding/gob"
	"github.com/ethtweet/ethtweet/logs"
	"time"
)

type P2pNetMessageReceiveInterface interface {
	ReceiveHandle(ctx context.Context, node *OnlineNode)
}

type P2pNetMessagePushInterface interface {
	PushHandle(writer *bufio.ReadWriter) error
}

type P2pNetMessage struct {
	Msg       P2pNetMessageReceiveInterface
	CreatedAt time.Time
}

func NewP2pNetMessageEncode(data P2pNetMessageReceiveInterface) (P2pNetMessagePushInterface, error) {
	return &P2pNetMessage{
		Msg:       data,
		CreatedAt: time.Now(),
	}, nil
}

func NewP2pNetMessageDecode(rw *bufio.ReadWriter) (P2pNetMessageReceiveInterface, error) {
	d := gob.NewDecoder(rw)
	pm := &P2pNetMessage{}
	if err := d.Decode(pm); err != nil {
		return nil, err
	}
	return pm, nil
}

func (p *P2pNetMessage) PushHandle(rw *bufio.ReadWriter) error {
	e := gob.NewEncoder(rw)
	if err := e.Encode(p); err != nil {
		return err
	}
	return nil
}

func (p *P2pNetMessage) ReceiveHandle(ctx context.Context, node *OnlineNode) {
	p.Msg.ReceiveHandle(ctx, node)
}

type HearBeat struct {
	Time time.Time
}

func NewHearBeat() *HearBeat {
	return &HearBeat{Time: time.Now()}
}

func (h *HearBeat) ReceiveHandle(ctx context.Context, node *OnlineNode) {
	logs.PrintlnInfo("receive from node ", node.GetIdPretty(), "'s hear beat check....")
}
