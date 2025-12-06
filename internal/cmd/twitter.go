package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/OrkWard/wormface/internal/utils"
	"github.com/OrkWard/wormface/pkg/twitter"
)

func CmdTwitter(args []string) {
	// Parse args
	twitterFlags := flag.NewFlagSet("wormface twitter", flag.ExitOnError)
	noVideo := twitterFlags.Bool("v", false, "Only download images")
	noImage := twitterFlags.Bool("i", false, "Only download videos")
	maxCount := twitterFlags.Int("l", 0, "Limit the number of media to download")
	twitterFlags.Parse(args)

	positionArgs := twitterFlags.Args()
	if len(positionArgs) < 1 {
		fmt.Println("Usage: wormface-cli twitter [options] <username>")
		twitterFlags.PrintDefaults()
		os.Exit(1)
	}
	userName := positionArgs[0]

	fmt.Printf("[INPUT] Scraping user: %s\n", userName)
	fmt.Printf("[OPTION] No video: %v, No image: %v, Limit: %d\n", *noVideo, *noImage, *maxCount)

	// Init client
	headers := http.Header{}
	if utils.Config.Twitter.CsrfToken == "" || utils.Config.Twitter.Authorization == "" || utils.Config.Twitter.Cookie == "" {
		panic("twitter config not set")
	}
	headers.Set("Cookie", utils.Config.Twitter.Cookie)
	headers.Set("X-CSRF-Token", utils.Config.Twitter.CsrfToken)
	headers.Set("Authorization", utils.Config.Twitter.Authorization)
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	headers.Set("Referer", fmt.Sprintf("https://x.com/%s/media", userName))

	client := twitter.NewTwitterClient(context.Background(), headers)
	defer client.Close()

	// Get user ID
	userId, err := client.GetUserId(userName)
	if err != nil {
		fmt.Printf("Error getting user ID: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[INFO] User ID: %s\n", userId)

	// Create output directory
	outputDir := fmt.Sprintf("output/%s", userName)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Printf("[ERROR] Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Implement media fetching loop
	var images, videos []string
	var cursor *string
	for {
		result, err := client.GetUserMedia(userId, cursor)
		if err != nil {
			fmt.Printf("\n[ERROR] Error getting user media: %v", err)
			os.Exit(1)
		}

		if len(result.Images) == 0 && len(result.Videos) == 0 {
			fmt.Print("\n[INFO] No more media found.")
			break
		}

		images = append(images, result.Images...)
		videos = append(videos, result.Videos...)

		fmt.Printf("\r[PROGRESS] %c Fetched %d images and %d videos...", utils.GetChar(), len(images), len(videos))

		if result.Cursor == "" {
			break
		}
		cursor = &result.Cursor

		if *maxCount > 0 && (len(images)+len(videos)) >= *maxCount {
			break
		}
	}
	fmt.Println()

	fmt.Printf("[INFO] Total images: %d, Total videos: %d\n", len(images), len(videos))

	// Save metadata
	saveMetadata(outputDir, "all_image.json", images)
	saveMetadata(outputDir, "all_video.json", videos)

	// Download files
	if !*noImage {
		fmt.Println("[INFO] Downloading images...")
		utils.DownloadAll(images, outputDir)
	}

	if !*noVideo {
		fmt.Println("[INFO] Downloading videos...")
		utils.DownloadAll(videos, outputDir)
	}
}

func saveMetadata(outputDir, fileName string, data []string) {
	filePath := fmt.Sprintf("%s/%s", outputDir, fileName)
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling metadata to JSON: %v\n", err)
		return
	}
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		fmt.Printf("Error writing metadata file %s: %v\n", filePath, err)
	}
}
