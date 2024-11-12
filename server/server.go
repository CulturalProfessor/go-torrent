package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

var (
	mainClient  *torrent.Client
	downloadDir = "./downloadedFiles"
)

func init() {
	cfg := torrent.NewDefaultClientConfig()
	cfg.Seed = false
	cfg.DataDir = downloadDir
	var err error
	mainClient, err = torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create torrent client: %v", err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

func uploadTorrent(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("torrent")
	if err != nil {
		http.Error(w, "Failed to read torrent file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mi, err := metainfo.Load(file)
	if err != nil {
		http.Error(w, "Failed to parse torrent file", http.StatusBadRequest)
		return
	}

	t, err := mainClient.AddTorrent(mi)
	if err != nil {
		http.Error(w, "Failed to add torrent", http.StatusInternalServerError)
		return
	}
	<-t.GotInfo()
	t.DownloadAll()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"torrentID": t.InfoHash().HexString()})
}

func getDownloadProgress(w http.ResponseWriter, r *http.Request) {
	torrentID := r.URL.Query().Get("id")
	infoHash := metainfo.NewHashFromHex(torrentID)
	if infoHash.IsZero() {
		http.Error(w, "Invalid torrent ID", http.StatusBadRequest)
		return
	}

	t, found := mainClient.Torrent(infoHash)
	if !found {
		http.Error(w, "Torrent not found", http.StatusNotFound)
		return
	}

	progress := float64(t.BytesCompleted()) / float64(t.Info().TotalLength()) * 100
	json.NewEncoder(w).Encode(map[string]float64{"progress": progress})
}

func serveDownloadedFile(w http.ResponseWriter, r *http.Request) {
	torrentID := r.URL.Query().Get("id")
	infoHash := metainfo.NewHashFromHex(torrentID)
	if infoHash.IsZero() {
		http.Error(w, "Invalid torrent ID", http.StatusBadRequest)
		return
	}

	t, found := mainClient.Torrent(infoHash)
	if !found {
		http.Error(w, "Torrent not found or incomplete", http.StatusNotFound)
		return
	}

	filePath := filepath.Join(downloadDir, t.Name())
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(filePath))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}

func main() {
	http.HandleFunc("/ping", ping)              
	http.HandleFunc("/upload", uploadTorrent)
	http.HandleFunc("/progress", getDownloadProgress)
	http.HandleFunc("/download", serveDownloadedFile)
	log.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
