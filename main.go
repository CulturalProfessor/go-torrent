package main

import (
	"fmt"
	"log"
	"os"

	parsetorrent "github.com/CulturalProfessor/go-torrent/parseTorrent"
	peers "github.com/CulturalProfessor/go-torrent/peers"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./program_name path/to/torrent_file")
		os.Exit(1)
	}
	filePath := os.Args[1]
	_, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	ioReader, _ := os.Open(filePath)

	torrentData, err := parsetorrent.Open(ioReader)
	if err != nil {
		log.Fatal(err)
	}

	peerID := generatePeerID()
	port := uint16(6881)

	// List of public trackers
	trackers := []string{
		"https://tracker.gbitt.info/announce",
	}

	for _, tracker := range trackers {
		torrentData.Announce = tracker
		trackerURL, err := peers.ExtractTrackerURL(peerID, port, torrentData)
		fmt.Println(trackerURL)
		if err != nil {
			log.Printf("Failed to extract tracker URL for %s: %v\n", tracker, err)
			continue
		}

		peerList, err := peers.RequestPeers(torrentData)
		if err != nil {
			log.Printf("Failed to request peers from %s: %v\n", tracker, err)
			continue
		}

		for _, peer := range peerList {
			fmt.Printf("Peer: %v:%d\n", peer.IP, peer.Port)
		}
		break
	}
}

func generatePeerID() [20]byte {
	return [20]byte{'-', 'P', 'C', '0', '0', '0', '1', '-', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2'}
}


// TODO : not parsing tracker url correctly