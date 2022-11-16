package p2pNet

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ethtweet/ethtweet/global"
	"github.com/ethtweet/ethtweet/keys"
	"github.com/ethtweet/ethtweet/logs"
	"github.com/libp2p/go-libp2p/p2p/host/autorelay"
	connmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2pquic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"runtime"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/p2p/protocol/ping"

	keystore "github.com/ipfs/go-ipfs-keystore"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	routing2 "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	libp2ptcp "github.com/libp2p/go-libp2p/p2p/transport/tcp"
	websocket "github.com/libp2p/go-libp2p/p2p/transport/websocket"
)

const MaxOnlineNodesNum = 128
const OnlineNodesSyncDuration = time.Second * 15
const OnlineNodesChkDuration = time.Minute * 3
const NodePingTimeoutDuration = time.Second * 6

type OnlineNode struct {
	Pi              peer.AddrInfo
	Stream          network.Stream
	Rw              *bufio.ReadWriter
	LastChkTime     time.Time
	IsInOnlineNodes bool
	usrNode         *UserNode
	IsLister        bool
}

func NewOnlineNode(usrNode *UserNode, stream network.Stream, isInOnlineNodes bool) *OnlineNode {
	return &OnlineNode{
		Pi: peer.AddrInfo{
			ID: stream.Conn().RemotePeer(),
		},
		Stream:          stream,
		Rw:              bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)),
		LastChkTime:     time.Now(),
		IsInOnlineNodes: isInOnlineNodes,
		usrNode:         usrNode,
		IsLister:        false,
	}
}

func (oln *OnlineNode) ListenRead() {
	if oln.IsLister {
		return
	}
	go func() {
		c, cc := context.WithCancel(oln.usrNode.Ctx)
		defer cc()
		ReadData(oln, c)
	}()
}

func (oln *OnlineNode) AddOnlineNodesTry() {
	if oln.usrNode.onlineNodes == nil {
		return
	}
	if oln.IsInOnlineNodes {
		return
	}
	logs.PrintlnInfo("AddOnlineNodesTry..........")
	oln.usrNode.LockedOnlineNode()
	defer oln.usrNode.UnLockedOnlineNode()
	if len(oln.usrNode.onlineNodes) >= MaxOnlineNodesNum {
		return
	}
	id := oln.GetIdPretty()
	_, ok := oln.usrNode.onlineNodes[id]
	if !ok {
		oln.usrNode.onlineNodes[id] = oln
	}
	oln.IsInOnlineNodes = true
}

func (oln *OnlineNode) GetIdPretty() string {
	return oln.Pi.ID.String()
}

func (oln *OnlineNode) Close() {
	oln.IsLister = false
	peerId := oln.Stream.Conn().RemotePeer().String()
	_, ok := oln.usrNode.onlineNodes[peerId]
	if ok {
		oln.IsInOnlineNodes = true
		oln.usrNode.RemoveOnlineNode(oln.Pi.ID.String())
		logs.PrintlnInfo("close stream from online nodes ....", peerId)
		return
	}
	oln.IsInOnlineNodes = false
	_ = oln.Stream.Close()
	_ = oln.Stream.Conn().Close()
	logs.PrintlnInfo("direct close stream ....", peerId)
}

func (oln *OnlineNode) WriteData(data P2pNetMessageReceiveInterface) error {
	err := WriteData(oln.Rw, data)
	if err != nil {
		oln.Close()
	}
	return err
}

type UserNode struct {
	Ctx               context.Context    //p2p上下文
	Cancel            context.CancelFunc //用户取消当前上下文 平滑退出
	priKey            *keys.PrivateKey   //用户私钥
	Protocol          string             //当前节点协议
	Host              host.Host
	UserKey           string //存储私钥的key
	UserData          string //存储私钥的目录
	dht               *dht.IpfsDHT
	Port              int
	routingDiscovery  *routing2.RoutingDiscovery
	ps                *ping.PingService
	onlineNodes       map[string]*OnlineNode
	onlineMux         sync.RWMutex
	isOnlineMuxLocked bool
}

// 新建一个用户节点
func NewUserNode(port int, userKey, userData string) *UserNode {
	usr := &UserNode{
		Protocol: global.TwitterProtocol,
		Port:     port,
	}
	usr.Ctx, usr.Cancel = context.WithCancel(global.GetGlobalCtx())
	usr.UserKey = userKey
	if usr.UserKey == "" {
		usr.UserKey = "userKey"
	}
	if runtime.GOOS == "android" {
		userData = "/sdcard/" + userData
	}

	usr.UserData = userData
	if usr.UserData == "" {
		usr.UserData = "./"
	}
	return usr
}

func (usr *UserNode) LockedOnlineNode() {
	usr.onlineMux.Lock()
	usr.isOnlineMuxLocked = true
}

func (usr *UserNode) UnLockedOnlineNode() {
	usr.isOnlineMuxLocked = false
	usr.onlineMux.Unlock()
}

func (usr *UserNode) IsLockedOnlineNodes() bool {
	return usr.isOnlineMuxLocked
}

// 在这里处理退出流程
func (usr *UserNode) Exit() {
	err := usr.Host.Close()
	if err != nil {
		logs.PrintErr("close host err ", err)
	}
	err = usr.dht.Close()
	if err != nil {
		logs.PrintErr("close dht err ", err)
	}
	cc := usr.Cancel
	cc()
	logs.PrintlnInfo("exit....")
	return
}

func (usr *UserNode) Ping(pi peer.ID, ctx context.Context) <-chan ping.Result {
	return usr.ps.Ping(ctx, pi)
}

func (usr *UserNode) Connect(pi peer.AddrInfo) error {
	return usr.Host.Connect(usr.Ctx, pi)
}

func (usr *UserNode) IDPretty() string {
	return usr.Host.ID().String()
}

func (usr *UserNode) NewStream(pi peer.AddrInfo) (network.Stream, error) {
	return usr.Host.NewStream(usr.Ctx, pi.ID, protocol.ID(usr.Protocol))
}

func (usr *UserNode) NewStreamCtx(ctx context.Context, pi peer.AddrInfo) (network.Stream, error) {
	return usr.Host.NewStream(ctx, pi.ID, protocol.ID(usr.Protocol))
}
func (usr *UserNode) ConnectP2p() error {
	var err error
	usr.priKey, err = usr.GetPriKey(usr.UserKey)
	if err != nil {
		return err
	}

	connmgr_, _ := connmgr.NewConnManager(
		50,  // Lowwater
		200, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	usr.Host, err = libp2p.New(
		libp2p.Identity(usr.priKey.LibP2pPrivate),
		//尝试开启upnp协议
		libp2p.NATPortMap(),
		libp2p.EnableNATService(),

		libp2p.DefaultPeerstore,
		//注册使用路由
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			usr.dht, err = dht.New(usr.Ctx, h, dht.BootstrapPeers(dht.GetDefaultBootstrapPeerAddrInfos()...))
			return usr.dht, err
		}),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.Transport(libp2ptcp.NewTCPTransport),
		libp2p.Transport(websocket.New),

		libp2p.EnableRelay(),
		libp2p.EnableNATService(),
		libp2p.EnableRelayService(),
		libp2p.ForceReachabilityPublic(),
		libp2p.EnableAutoRelay(autorelay.WithDefaultStaticRelays(), autorelay.WithCircuitV1Support(), autorelay.WithNumRelays(20)),

		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", usr.Port),
			fmt.Sprintf("/ip6/::/tcp/%d", usr.Port),

			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", usr.Port),
			fmt.Sprintf("/ip6/::/udp/%d/quic", usr.Port),

			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/ws", usr.Port),
			fmt.Sprintf("/ip6/::/tcp/%d/ws", usr.Port),

			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic/webtransport", usr.Port),
			fmt.Sprintf("/ip6/::/udp/%d/quic/webtransport", usr.Port),
		),
		libp2p.ConnectionManager(connmgr_),
	)
	if err != nil {
		return err
	}

	usr.Host.SetStreamHandler(protocol.ID(usr.Protocol), usr.handleStream)
	//启动dht
	if err = usr.dht.Bootstrap(usr.Ctx); err != nil {
		return err
	}

	//每分钟广播一次
	go func() {
		time.Sleep(time.Minute * 1)
		usr.routingDiscovery.Advertise(usr.Ctx, usr.Protocol)

	}()
	//广播自己的位置
	usr.routingDiscovery = routing2.NewRoutingDiscovery(usr.dht)
	usr.onlineNodes = make(map[string]*OnlineNode, MaxOnlineNodesNum)
	usr.ps = ping.NewPingService(usr.Host)

	// 开启一个线程去同步节点
	go func() {
		for {
			logs.PrintDebug("开始同步节点状态.....")
			select {
			case <-usr.Ctx.Done():
				return
			default:
				logs.Println("SyncOnlineNodes 1")
				usr.SyncOnlineNodes()
				logs.Println("SyncOnlineNodes end")
			}
			time.Sleep(OnlineNodesSyncDuration)
		}
	}()

	return nil
}

func (usr *UserNode) handleStream(s network.Stream) {
	logs.PrintlnSuccess("Got a new stream!", s.Conn().RemotePeer().String(), s.ID(), s.Conn().ID(), len(s.Conn().GetStreams()))
	node := NewOnlineNode(usr, s, false)
	node.AddOnlineNodesTry()
	node.ListenRead()
}

// 获取用户私钥 没有则创建
func (usr *UserNode) GetPriKey(key string) (*keys.PrivateKey, error) {
	if (key == usr.UserKey || key == "") && usr.priKey != nil {
		return usr.priKey, nil
	}
	ks, err := keystore.NewFSKeystore(usr.UserData)
	if err != nil {
		return nil, err
	}
	ok, err := ks.Has(key)
	if err != nil {
		logs.PrintErr(err)
	}
	var priKey *keys.PrivateKey
	if ok {
		priKey2, err := ks.Get(key)
		if err != nil {
			return nil, err
		}
		return keys.NewPrivateKeyByLibP2pPri(priKey2)
	} else {
		priKey, err = keys.NewPrivateKey()
		if err != nil {
			return nil, err
		}
		err = ks.Put(key, priKey.LibP2pPrivate)
		if err != nil {
			return nil, err
		}
	}
	return priKey, nil
}

func (usr *UserNode) SignMsg(key string, msg string) (string, error) {
	if key == "" {
		key = usr.UserKey
	}
	priv, err := usr.GetPriKey(key)
	if err != nil {
		return "", err
	}
	return priv.Sign(msg)
}

func (usr *UserNode) FindPeers() (<-chan peer.AddrInfo, error) {
	//刷新路由表 防止节点下线后仍然被发现
	return usr.routingDiscovery.FindPeers(usr.Ctx, usr.Protocol)
}

func (usr *UserNode) WaitOnlineNodeLimitNum(num int) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		for {
			select {
			case <-usr.Ctx.Done():
				return
			default:
				n := len(usr.onlineNodes)
				if num == 0 {
					if n > num {
						c <- struct{}{}
					}
				} else if n >= num {
					c <- struct{}{}
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return c
}

// 检查并等待在线列表
func (usr *UserNode) WaitOnlineNode() <-chan struct{} {
	return usr.WaitOnlineNodeLimitNum(0)
}

func (usr *UserNode) SyncOnlineNodes() {
	//维护已有节点数
	var wg sync.WaitGroup
	usr.CheckOnlineNodes()
	//缓存节点数小于最大时 继续维护
	if len(usr.onlineNodes) < MaxOnlineNodesNum {
		usr.routingDiscovery.Advertise(usr.Ctx, usr.Protocol)
		logs.PrintlnInfo("find online new nodes..........")
		peers, err := usr.FindPeers()
		if err != nil {
			logs.PrintDebugErr(err)
			return
		}
		logs.PrintlnInfo("FindPeers: len", len(peers))
		for pi := range peers {
			if pi.ID.String() == usr.IDPretty() {
				continue
			}
			wg.Add(1)
			go func(pi peer.AddrInfo) {
				defer wg.Done()
				err := usr.setOnlineNodes(pi)
				if err != nil {
					logs.PrintDebugErr(err, pi.ID.String())
				}
			}(pi)
		}
		wg.Wait()
	}
	logs.PrintlnInfo("online nodes ok;  len ", len(usr.onlineNodes))
}

func (usr *UserNode) RemoveOnlineNode(id string) {
	if usr.IsLockedOnlineNodes() {
		usr.RemoveOnlineNodeLocked(id)
		return
	}
	usr.LockedOnlineNode()
	defer usr.UnLockedOnlineNode()
	usr.RemoveOnlineNodeLocked(id)
	return
}

func (usr *UserNode) RemoveOnlineNodeLocked(id string) {
	if oln, ok := usr.onlineNodes[id]; ok {
		_ = oln.Stream.Close()
		_ = oln.Stream.Conn().Close()
		oln.IsInOnlineNodes = false
		delete(usr.onlineNodes, id)
	}
}

func (usr *UserNode) GetOnlineNodesCount() int {
	return len(usr.onlineNodes)
}

func (usr *UserNode) EachOnlineNodes(fn func(node *OnlineNode) bool) {
	usr.LockedOnlineNode()
	defer usr.UnLockedOnlineNode()
	for _, node := range usr.onlineNodes {
		if !fn(node) {
			return
		}
	}
}

func (usr *UserNode) GetOnlineNode(id string) *OnlineNode {
	oln, ok := usr.onlineNodes[id]
	if !ok {
		return nil
	}
	return oln
}

func (usr *UserNode) setOnlineNodes(pi peer.AddrInfo) error {
	defer func() {
		//这里可能会产生并发操作map的错误
		//为了不影响速度 故不使用锁 所以这里捕获一下不让程序终止
		if r := recover(); r != nil {
			logs.Println("setOnlineNodes err end")
			logs.PrintDebugErr(r)
		}
	}()
	peerId := pi.ID.String()
	logs.PrintlnInfo("fund node ", peerId)

	_, ok := usr.onlineNodes[peerId]
	if len(usr.onlineNodes) < MaxOnlineNodesNum && !ok {
		c, cc := context.WithTimeout(usr.Ctx, NodePingTimeoutDuration)
		defer cc()
		err := usr.Host.Connect(c, pi)
		if err != nil {
			logs.PrintlnWarning("node Connect fail ", peerId, err)
			return err
		}
		s, err := usr.NewStreamCtx(c, pi)
		logs.PrintlnInfo("setOnlineNodes NewStream ", peerId)
		if err != nil {
			logs.PrintlnWarning("node NewStream fail ", peerId, err)
			return err
		}
		usr.LockedOnlineNode()
		usr.onlineNodes[peerId] = NewOnlineNode(usr, s, true)
		usr.onlineNodes[peerId].ListenRead()
		usr.UnLockedOnlineNode()
	} else {
		logs.PrintlnInfo("setOnlineNodes is  ", peerId)
	}

	logs.PrintlnSuccess("setOnlineNodes OK ", peerId)
	return nil
}

func (usr *UserNode) CheckOnlineNodes() {
	logs.PrintlnInfo("check online nodes.....")
	for _, node := range usr.onlineNodes {
		now := time.Now()
		id := node.GetIdPretty()
		if node.LastChkTime.After(now.Add(-OnlineNodesChkDuration)) {
			logs.Println(id, " online node sync duration")
			continue
		}
		err := node.WriteData(NewHearBeat())
		if err != nil {
			logs.PrintErr("online node ", id, " is closed... ", err)
			usr.RemoveOnlineNode(id)
			continue
		}
		node.LastChkTime = now
		logs.PrintlnSuccess("online node ", id, "  is ok")
	}
}
