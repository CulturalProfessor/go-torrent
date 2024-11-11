package main

import (
	"fmt"
	"log"
	"os"
	"github.com/CulturalProfessor/go-torrent/download"
	parseTorrent "github.com/CulturalProfessor/go-torrent/parseTorrent"
	// "github.com/CulturalProfessor/go-torrent/peers"
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

	download.DownloadFromPeer(torrentData)
}
