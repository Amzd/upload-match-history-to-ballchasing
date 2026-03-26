package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dank/rlapi"
)

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	fmt.Println("Storing cache in ", getCacheDir())

	rpc, _ := RPC()
	defer rpc.Close()

	uploaded := loadUploadedCache()
	replays := getReplays(rpc)
	rpc.Close()

	for _, replay := range replays {
		if !uploaded[replay.Match.MatchGUID] {
			filePath, err := downloadFile(replay.ReplayUrl)
			if err != nil {
				fmt.Println("Download error:", err)
				continue
			}

			fmt.Println("Uploading:", filePath)

			err = uploadReplay(filePath)
			if err != nil {
				fmt.Println("Upload error:", err)
			} else {
				appendToUploadedCache(replay.Match.MatchGUID)
			}

			os.Remove(filePath)
		} else {
			fmt.Println("Skipping replay as it was already uploaded")
		}
	}
}

func getReplays(rpc *rlapi.PsyNetRPC) []rlapi.MatchEntry {
	apiCtx, apiCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer apiCancel()

	historyResp, err := rpc.GetMatchHistory(apiCtx)
	if err != nil {
		slog.Error("Failed to get match history", slog.Any("error", err))
		panic(err)
	}

	return historyResp
}

func downloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to download: %s", resp.Status)
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "*.replay")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func loadUploadedCache() map[string]bool {
	cache := make(map[string]bool)

	file, err := os.Open(uploadedCacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return cache // no cache yet
		}
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url != "" {
			cache[url] = true
		}
	}

	return cache
}

func appendToUploadedCache(url string) {
	f, err := os.OpenFile(uploadedCacheFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(url + "\n")
	if err != nil {
		panic(err)
	}
}


var rlTokenFile = filepath.Join(getCacheDir(), ".rltoken")
var bcTokenFile = filepath.Join(getCacheDir(), ".bctoken")
var uploadedCacheFile = filepath.Join(getCacheDir(), ".uploaded")

func getCacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	dir = filepath.Join(dir, "upload-match-history-to-ballchasing")
	os.MkdirAll(dir, 0700)
	return dir
}
