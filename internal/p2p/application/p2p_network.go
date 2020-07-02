package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	multiplex "github.com/libp2p/go-libp2p-mplex"
	secio "github.com/libp2p/go-libp2p-secio"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-tcp-transport"
	ws "github.com/libp2p/go-ws-transport"
	"time"

	//"github.com/libp2p/go-ud"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"perch/internal/p2p"
	_ "perch/pkg/general/log"
	"syscall"
)

func main() {
	var p2pnetwork p2p.NetworkP2P
	var p2pOptions []libp2p.Option
	var err error

	ctx, cacel := context.WithCancel(context.Background())
	defer cacel()
	security := libp2p.Security(secio.ID, secio.New)
	muxers := libp2p.ChainOptions(
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
		libp2p.Muxer("/mplex/6.7.0", multiplex.DefaultTransport),
	)
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(ws.New),
	)
	listenAddr := libp2p.ListenAddrStrings(
		"/ip4/0.0.0.0/tcp/0", "/ip6/::/tcp/0/ws")

	p2pOptions = append(p2pOptions, muxers, security, listenAddr, transports)

	options := p2p.NetworkRuntimeOptions{
		Ctx:            ctx,
		NetworkOptions: p2pOptions,
	}

	p2pnetwork.StartNetworkP2P(options)

	for _, addr := range p2pnetwork.NetworkPeer.Addrs() {
		fmt.Printf("Listening P2P on %s/p2p/%s\n", addr.String(), p2pnetwork.NetworkPeer.ID().String())
	}

	//	pubs ,err := pubsub.NewGossipSub(ctx,p2pnetwork.NetworkPeer)
	//pubs, err := pubsub.NewGossipSub(ctx, p2pnetwork.NetworkPeer)
	pubs, err := p2p.PubsubgossipGen(ctx, p2pnetwork.NetworkPeer)
	if err != nil {
		log.Error(err)
	}
	sub, topsub, err := p2p.PubsubtopicsJoin(pubs, p2p.Pubsub_Default_Topic)
	if err != nil {
		log.Error(err)
	}

	err = p2p.MDNSDiscoverySetup(ctx, p2pnetwork.NetworkPeer, p2p.DiscoveryInterval, p2p.DiscoveryServiceTag)
	if err != nil {
		log.Error(err)
	}
	go func() {
		for {
			msg := new(p2p.PubsubMessage)
			msg.SenderPeer = p2pnetwork.NetworkPeer.ID().Pretty()
			msg.PMessageStr = "hello world"
			msg.SenderFrom = "from localhost"
			err = p2p.PubsubTopicPubish(ctx, *msg, topsub, nil)
			if err != nil {
				log.Error(err)
			}
			time.Sleep(3 * time.Second)
		}

	}()
	msgChan := make(chan interface{})
	go p2p.PubsubMsgHandler(sub, ctx, p2pnetwork.NetworkPeer, msgChan)
	/*	for {
		select {
		case msg := <- msgChan:
			fmt.Println("msg from msg chan is:",msg)
		case <-time.Tick(5*time.Second):
			fmt.Println("time.is over")
			//return
		}
	}*/
	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case msg := <-msgChan:
			fmt.Println("msg from msg chan is:", msg)

		case err := <-errChan:
			if err != nil {
				log.Println(err)
			}
			// 执行额外的清理操作
			for _, clean := range p2pnetwork.NetworkCleanFunc {
				clean()
			}
			p2pnetwork.NetworkPeer.Close()
			return
		case s := <-signalChan:
			log.Printf("捕获到信号%v，准备停止服务\n", s)
			p2pnetwork.NetworkPeer.Close()
			return
		}
	}
}
