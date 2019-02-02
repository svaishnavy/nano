package node

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const packetSize = 512
const numberOfPeersToShare = 8

var conn *net.UDPConn

var DefaultPeer = Peer{
	//net.ParseIP("::ffff:192.168.0.70"),
	net.ParseIP("::ffff:94.130.105.241"),
	7075,
	nil,
}

var PeerList = []Peer{DefaultPeer}
var PeerSet = map[string]bool{DefaultPeer.String(): true}

func (p *Peer) SendMessage(m Message) error {
	now := time.Now()
	p.LastReachout = &now

	buf := bytes.NewBuffer(nil)
	err := m.Write(buf)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(buf.Bytes(), &net.UDPAddr{Port: int(p.Port), IP: p.IP})
	if err != nil {
		return err
	}

	return nil
}

func ListenForUdp() error {
	log.Printf("Listening for udp packets on 7075")
	var err error
	conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: 7075})
	if err != nil {
		return err
	}

	buf := make([]byte, packetSize)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error: UDP read error: %v", err)
			continue
		}
		if n > 0 {
			log.Println("Received message")
			handleMessage(bytes.NewBuffer(buf[:n]))
		}
	}
}

func SendKeepAlive(peer Peer) error {
	randomPeers := make([]Peer, 0)
	randIndices := rand.Perm(len(PeerList))
	for n, i := range randIndices {
		if n == numberOfPeersToShare {
			break
		}
		randomPeers = append(randomPeers, PeerList[i])
		fmt.Println(PeerList[i].IP, []byte(PeerList[i].IP))
	}

	m := CreateKeepAlive(randomPeers)
	log.Println("Sending keepalive")
	return peer.SendMessage(m)
}

func SendKeepAlives(params []interface{}) {
	peers := params[0].([]Peer)
	timeCutoff := time.Now().Add(-5 * time.Minute)

	for _, peer := range peers {
		if peer.LastReachout == nil || peer.LastReachout.Before(timeCutoff) {
			err := SendKeepAlive(peer)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
