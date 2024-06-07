package peers

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"github.com/CulturalProfessor/go-torrent/parseTorrent" 
	bencode "github.com/jackpal/bencode-go"
)

type Peer struct {
	IP   [4]byte
	Port uint16
}

func UnmarshalPeers(response []byte) []Peer {
	const peerSize = 6 // 4 bytes for IP and 2 bytes for port
	numPeers := len(response) / peerSize
	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		copy(peers[i].IP[:], response[offset:offset+4])
		peers[i].Port = binary.BigEndian.Uint16(response[offset+4 : offset+6])
	}

	return peers
}

func ExtractTrackerURL(peerID [20]byte, port uint16, torrentData parseTorrent.BencodeTorrent) (string, error) {
	announce := torrentData.Announce
	length := torrentData.Info.Length
	infoHash, err := hashInfo(torrentData.Info)

	if err != nil {
		return "", fmt.Errorf("error hashing torrent info: %v", err)
	}

	base, err := url.Parse(announce)
	if err != nil {
		return "", fmt.Errorf("error parsing announce URL: %v", err)
	}

	// Remove the port from the base URL if it is present
	base.Host = base.Hostname()

	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(length)},
	}

	// Attempt HTTPS connection first
	base.Scheme = "https"
	base.RawQuery = params.Encode()
	httpsURL := base.String()

	client := &http.Client{
		Timeout: 10 * time.Second, // Set timeout to 10 seconds
	}
	resp, err := client.Get(httpsURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		return httpsURL, nil
	}

	// Fall back to HTTP if HTTPS failed
	base.Scheme = "http"
	base.RawQuery = params.Encode()
	httpURL := base.String()

	resp, err = http.Get(httpURL)
	if err != nil {
		return "", fmt.Errorf("both HTTPS and HTTP requests failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tracker did not return 200 OK: %v", resp.Status)
	}

	return httpURL, nil
}

func hashInfo(info parseTorrent.BencodeInfo) ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func RequestPeers(torrentData parseTorrent.BencodeTorrent) ([]Peer, error) {

	peerID:=generatePeerID()
	port := uint16(6881)

	trackerURL,err:= ExtractTrackerURL(peerID, port, torrentData)
	fmt.Println(trackerURL)
	if err != nil {
		return nil, fmt.Errorf("error extracting tracker URL: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second, // Set timeout to 10 seconds
	}

	resp, err := client.Get(trackerURL)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to tracker: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tracker did not return 200 OK: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading tracker response body: %v", err)
	}

	return UnmarshalPeers(body), err
}

func generatePeerID() [20]byte {
	return [20]byte{'-', 'P', 'C', '0', '0', '0', '1', '-', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2'}
}
