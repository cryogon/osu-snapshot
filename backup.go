package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

var osuBeatmapSetIDs []int

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	return err
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := path.Join(src, entry.Name())
		dstPath := path.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0o755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error retrieving user home directory:", err)
		return
	}

	osuPath := path.Join(homeDir, ".local/share/osu-wine/osu!")
	fmt.Println("Osu! installation path:", osuPath)
	songFolderPath := path.Join(osuPath, "Songs")
	fmt.Println("Osu! Songs folder path:", songFolderPath)

	if _, err := os.Stat(songFolderPath); os.IsNotExist(err) {
		fmt.Println("Songs folder does not exist.")
		return
	}

	files, err := os.ReadDir(songFolderPath)

	if err != nil {
		fmt.Println("Error reading Songs directory:", err)
		return
	}

	for _, file := range files {
		// filter only directories
		if !file.IsDir() {
			continue
		}

		chunks := strings.Split(file.Name(), " ")

		if len(chunks) < 2 {
			continue
		}

		beatmapsetID := chunks[0]

		id, err1 := strconv.Atoi(beatmapsetID) // Atoi converts string to int
		if err1 != nil {
			// read the .osu file inside the folder to get the beatmapset ID from the [Metadata] section (if available)
			osuFiles, err2 := os.ReadDir(path.Join(songFolderPath, file.Name()))
			if err2 != nil {
				fmt.Println("Error reading beatmap directory:", err2)
				continue
			}

			foundBeatmapId := false
			// Look for .osu files
			for _, osuFile := range osuFiles {
				if !strings.HasSuffix(osuFile.Name(), ".osu") {
					continue
				}

				// Read and parse the .osu file
				content, err := os.ReadFile(path.Join(songFolderPath, file.Name(), osuFile.Name()))
				if err != nil {
					fmt.Println("Error reading .osu file:", err)
					continue
				}

				// Look for BeatmapSetID in the content
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					if strings.HasPrefix(line, "BeatmapSetID:") {
						idStr := strings.TrimSpace(strings.TrimPrefix(line, "BeatmapSetID:"))
						if id, err := strconv.Atoi(idStr); err == nil {
							osuBeatmapSetIDs = append(osuBeatmapSetIDs, id)
							foundBeatmapId = true
							break
						}
					}
				}

				if foundBeatmapId {
					break
				}
			}
			continue
		}

		osuBeatmapSetIDs = append(osuBeatmapSetIDs, id)
	}

	// Create a map to store unique IDs
	unique := make(map[int]bool)
	// Add all IDs to the map
	for _, id := range osuBeatmapSetIDs {
		unique[id] = true
	}
	// Clear the slice and add only unique IDs back
	osuBeatmapSetIDs = osuBeatmapSetIDs[:0]
	for id := range unique {
		osuBeatmapSetIDs = append(osuBeatmapSetIDs, id)
	}
	fmt.Printf("Total beatmaps found: %d\n", len(osuBeatmapSetIDs))

	backupDir := path.Join(osuPath, "backup")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		fmt.Println("Error creating backup directory:", err)
		return
	}

	if len(osuBeatmapSetIDs) == 0 {
		fmt.Println("No beatmap IDs to backup.")
		return
	}

	// Build a JSON array of numbers without adding a new import
	parts := make([]string, len(osuBeatmapSetIDs))
	for i, id := range osuBeatmapSetIDs {
		parts[i] = strconv.Itoa(id)
	}
	jsonData := "[" + strings.Join(parts, ",") + "]"

	filePath := "./backup/beatmap_ids.json"
	if err := os.WriteFile(filePath, []byte(jsonData), 0o644); err != nil {
		fmt.Println("Error writing JSON backup:", err)
		return
	}

	// Files to copy from osu! directory to backup
	filesToCopy := []string{"Replays", "scores.db", "osu!.db", "collection.db", "presence.db", "Data"}

	for _, file := range filesToCopy {
		srcPath := path.Join(osuPath, file)
		dstPath := path.Join("backup", file)

		// Skip if source doesn't exist
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			fmt.Printf("Skipping %s: file not found\n", file)
			continue
		}

		// Copy directories recursively
		if info, err := os.Stat(srcPath); err == nil && info.IsDir() {
			if err := os.RemoveAll(dstPath); err != nil {
				fmt.Printf("Error removing existing %s directory: %v\n", file, err)
				continue
			}
			if err := os.MkdirAll(dstPath, 0o755); err != nil {
				fmt.Printf("Error creating %s directory: %v\n", file, err)
				continue
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				fmt.Printf("Error copying %s directory: %v\n", file, err)
			}
			continue
		}

		// Copy regular files
		if err := copyFile(srcPath, dstPath); err != nil {
			fmt.Printf("Error copying %s: %v\n", file, err)
		}
	}

	fmt.Println("Backup saved to", filePath)
}
