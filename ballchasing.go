package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const tokenFile = ".bctoken"
const apiURL = "https://ballchasing.com/api/replays"

func uploadReplay(filePath string) error {
	token := getToken()
	for {
		if err := verifyToken(token); err != nil {
			fmt.Println("Invalid token, try again.")
			os.Remove(tokenFile)
			token = getToken()
			continue
		}
		break
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add file field (must be named "file")
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://ballchasing.com/api/v2/upload", &body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 201 {
		fmt.Println("Upload successful ✅")
	} else if resp.StatusCode == 409 {
		fmt.Println("Upload duplicate ✅")
	} else {
		return fmt.Errorf("Upload failed: %s\n%s", resp.Status, string(respBody))
	}

	fmt.Println(string(respBody))

	return nil
}

func getToken() string {
	// Try reading from file
	if data, err := os.ReadFile(tokenFile); err == nil {
		return strings.TrimSpace(string(data))
	}

	// Prompt user
	fmt.Print("Enter your ballchasing.com API token: ")
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')
	token = strings.TrimSpace(token)

	// Save to file
	err := os.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		fmt.Println("Warning: could not save token:", err)
	}

	return token
}

func verifyToken(token string) error {
	req, err := http.NewRequest("GET", "https://ballchasing.com/api/", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("token invalid: %s\n%s", resp.Status, string(body))
}
