package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func Restore() {
	backupFolder := "./backup"
	if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
		fmt.Println("No backup found to restore.")
		return
	}
	beatmapIds, err := os.ReadFile(path.Join(backupFolder, "beatmap_ids.json"))
	if err != nil {
		fmt.Println("Error reading beatmap IDs:", err)
		return
	}
	var ids []int
	if err := json.Unmarshal(beatmapIds, &ids); err != nil {
		fmt.Println("Error parsing beatmap IDs:", err)
		return
	}
	if len(ids) == 0 {
		fmt.Println("No beatmaps to restore")
		return
	}
	// only downloading the first map for testing
	beatmapId := ids[0]

	// Create a client that doesn't follow redirects automatically
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			fmt.Printf("Redirect to: %s\n", req.URL.String())
			return nil
		},
	}

	url := fmt.Sprintf("https://osu.ppy.sh/osu/%d", beatmapId)
	fmt.Printf("Downloading from: %s\n", url)

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error downloading beatmap %d: %v\n", beatmapId, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Content-Length: %d\n", resp.ContentLength)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))

	// Read the body to see what we're actually getting
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Create the file and write the body
	file, err := os.Create(fmt.Sprintf("%d.osu", beatmapId))
	if err != nil {
		fmt.Printf("Error creating file for beatmap %d: %v\n", beatmapId, err)
		return
	}
	defer file.Close()

	written, err := file.Write(body)
	if err != nil {
		fmt.Printf("Error saving beatmap %d: %v\n", beatmapId, err)
		return
	}

	fmt.Printf("Successfully restored beatmap %d (%d bytes written)\n", beatmapId, written)
}
