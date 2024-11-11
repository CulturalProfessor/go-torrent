package peers

import (
	"bytes"
	// "context"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/CulturalProfessor/go-torrent/metainfo"
	"github.com/CulturalProfessor/go-torrent/parseTorrent"
	udptracker "github.com/CulturalProfessor/go-torrent/trackers"
	bencode "github.com/jackpal/bencode-go"
)

type Peer struct {
	IP   [4]byte
	Port uint16
}

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type AnnounceRequest struct {
	InfoHash metainfo.Hash // Required
	PeerID   metainfo.Hash // Required

	Uploaded   int64  // Required, but default: 0, which should be only used for test or first.
	Downloaded int64  // Required, but default: 0, which should be only used for test or first.
	Left       int64  // Required, but default: 0, which should be only used for test or last.
	Event      uint32 // Required, but default: 0

	IP      net.IP // Optional
	Key     int32  // Optional
	NumWant int32  // Optional
	Port    uint16 // Optional
}

type tclient struct {
	url  string
	udp  *udptracker.Client
	exts []udptracker.Extension
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

func (c *tclient) ExtractTrackerURL(announce string, peerID [20]byte, port uint16, torrentData parseTorrent.BencodeTorrent) (string, error) {
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
	if base.Scheme == "udp" {
		// udpAnnounceResponse := udptracker.Announce(context.Background(), udptracker.AnnounceRequest{
		// 	InfoHash:   metainfo.NewHash(infoHash[:]),
		// 	PeerID:     metainfo.NewHash(peerID[:]),
		// 	Downloaded: 0,
		// 	Left:       int64(torrentData.Info.Length),
		// 	Uploaded:   0,

		// 	Event: 0,

		// 	Key:     0,
		// 	NumWant: -1,
		// 	Port:    port,
		// })

		return "", fmt.Errorf("UDP trackers are not supported")
	} else {
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

		resp, err = client.Get(httpURL)
		if err != nil {
			return "", fmt.Errorf("both HTTPS and HTTP requests failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("tracker did not return 200 OK: %v", resp.Status)
		}

		return httpURL, nil
	}

}

func hashInfo(info parseTorrent.BencodeInfo) ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, info)
	if err != nil {
		return [20]byte{}, err
	}
	return sha1.Sum(buf.Bytes()), nil
}

func RequestPeers(torrentData parseTorrent.BencodeTorrent) ([]string, error) {
	peerID := generatePeerID()
	port := uint16(6881)
	trackerURLs := append([]string{torrentData.Announce}, flattenAnnounceList(torrentData.AnnounceList)...)

	var lastError error
	var workingPeers []string

	for _, trackerURL := range trackerURLs {
		url, err := (&tclient{}).ExtractTrackerURL(trackerURL, peerID, port, torrentData)
		fmt.Println("Trying tracker URL:", url)
		if err != nil {
			lastError = err
			fmt.Println("Error extracting tracker URL:", err)
			continue
		}

		client := &http.Client{
			Timeout: 10 * time.Second, // Set timeout to 10 seconds
		}

		resp, err := client.Get(url)
		if err != nil {
			lastError = fmt.Errorf("error making GET request to tracker: %v", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastError = fmt.Errorf("tracker did not return 200 OK: %v", resp.Status)
			continue
		}

		var trackerResponse TrackerResponse
		err = bencode.Unmarshal(resp.Body, &trackerResponse)
		if err != nil {
			lastError = fmt.Errorf("error decoding tracker response: %v", err)
			continue
		}

		peers := UnmarshalPeers([]byte(trackerResponse.Peers))
		for _, peer := range peers {
			ip := net.IP(peer.IP[:]).String()
			if net.ParseIP(ip) == nil {
				continue // Skip invalid IPs
			}
			workingPeers = append(workingPeers, fmt.Sprintf("%s:%d", ip, peer.Port))
		}
	}

	if len(workingPeers) == 0 {
		return nil, lastError
	}

	return workingPeers, nil
}

func flattenAnnounceList(announceList [][]string) []string {
	var result []string
	for _, sublist := range announceList {
		result = append(result, sublist...)
	}
	return result
}

func generatePeerID() [20]byte {
	return [20]byte{'-', 'P', 'C', '0', '0', '0', '1', '-', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '1', '2'}
}
