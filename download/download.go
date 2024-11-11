package download

import (
	"fmt"
	"net"
	"time"

	"github.com/CulturalProfessor/go-torrent/parseTorrent"
	"github.com/CulturalProfessor/go-torrent/peers"
)

func DownloadFromPeer(torrentData parseTorrent.BencodeTorrent) {
	peerIPs, err := peers.RequestPeers(torrentData)
	if err != nil {
		fmt.Println("Error requesting peers:", err)
		return
	}

	for _, peer := range peerIPs {
		fmt.Printf("Attempting to connect to peer: %s\n", peer)
		conn, err := connectToPeer(peer, 3*time.Second)
		if err != nil {
			fmt.Printf("Error connecting to peer %s: %v\n", peer, err)
			continue
		}
		defer conn.Close()

		fmt.Println("Connected to peer", peer)
		// Handle communication with the peer
		break // Exit after successfully connecting to a peer
	}
}

func connectToPeer(peer string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout("tcp", peer, timeout)
}
