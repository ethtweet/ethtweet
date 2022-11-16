package p2pNet

import (
	"bufio"
	"context"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"io"
)

func ReadData(online *OnlineNode, ctx context.Context) {
	logs.PrintlnInfo("listen's ", online.Pi.ID.String())
	online.IsLister = true
	defer func() {
		online.Close()
		logs.PrintlnInfo("closed listen's ", online.Pi.ID.String())
	}()
	for {
		select {
		case <-ctx.Done():
			logs.PrintErr("main close....", ctx.Err())
			return
		default:
		}
		msg, err := NewP2pNetMessageDecode(online.Rw)
		if err != nil {
			if err == io.EOF {
				continue
			}
			logs.PrintErr("read err ", err)
			return
		}
		logs.PrintlnInfo("read new msg........... from", online.Pi.ID.String(), msg)
		go func() {
			msg.ReceiveHandle(ctx, online)
			online.AddOnlineNodesTry()
		}()
	}
}

func WriteData(rw *bufio.ReadWriter, data P2pNetMessageReceiveInterface) error {
	pm, err := NewP2pNetMessageEncode(data)
	if err != nil {
		return err
	}
	err = pm.PushHandle(rw)
	if err != nil {
		return err
	}
	return rw.Flush()
}

func StreamToRw(s network.Stream) (*bufio.Reader, *bufio.Writer) {
	return bufio.NewReader(s), bufio.NewWriter(s)
}

func IdToOnline(usrNode *UserNode, id string, ctx context.Context) (*OnlineNode, error) {
	//尝试充在线列表里获取连接
	var oln *OnlineNode
	oln = usrNode.GetOnlineNode(id)
	if oln != nil {
		logs.PrintlnSuccess("get online nodes .............")
		return oln, nil
	}
	maddr, err := multiaddr.NewMultiaddr("/ipfs/" + id)
	if err != nil {
		return nil, err
	}
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}
	usrNode.Host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	s, err := usrNode.NewStreamCtx(ctx, *info)
	if err != nil {
		return nil, err
	}
	oln = NewOnlineNode(usrNode, s, false)
	oln.AddOnlineNodesTry()
	oln.ListenRead()
	return oln, nil
}
