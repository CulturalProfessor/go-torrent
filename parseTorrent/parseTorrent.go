package parseTorrent

import (
	"fmt"
	"io"
	"github.com/jackpal/bencode-go"
)

// BencodeInfo represents the info dictionary in a torrent file.
type BencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

// BencodeTorrent represents the structure of a torrent file.
type BencodeTorrent struct {
	Announce     string          `bencode:"announce"`
	AnnounceList [][]string      `bencode:"announce-list"`
	Info         BencodeInfo     `bencode:"info"`
}

// Open reads and parses a torrent file from the provided io.Reader.
func Open(r io.Reader) (BencodeTorrent, error) {
	bto := BencodeTorrent{}
	err := bencode.Unmarshal(r, &bto)
	if err != nil {
		fmt.Println("Error unmarshalling torrent data")
		return BencodeTorrent{}, err
	}
	fmt.Println("Primary announce URL:", bto.Announce)
	fmt.Println("Additional announce URLs:", bto.AnnounceList)
	return bto, nil
}
