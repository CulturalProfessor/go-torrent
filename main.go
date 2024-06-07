package main

import (
	"fmt"
	"log"
	"os"

	parseTorrent "github.com/CulturalProfessor/go-torrent/parseTorrent"
	"github.com/CulturalProfessor/go-torrent/peers"
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

	torrentData, err := parseTorrent.Open(ioReader)
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("torrentData.txt", []byte(fmt.Sprintf("%+v", torrentData.Info)), 0644)

	peers, err := peers.RequestPeers(torrentData)
	if err != nil {
		fmt.Println("Error requesting peers")
	}

	for _, peer := range peers {
		fmt.Println(peer)
	}
}
