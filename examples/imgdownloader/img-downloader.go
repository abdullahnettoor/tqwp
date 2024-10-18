package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/abdullahnettoor/tqwp"
)

// FileDownloadTask represents a task to download a file from a URL.
type FileDownloadTask struct {
	tqwp.TaskModel
	Id       uint
	URL      string
	FileName string
}

// Process downloads the file from URL and saves it inside the "downloads" folder in the "imgdownloader" directory.
func (t *FileDownloadTask) Process() error {
	// Ensure the "downloads" folder exists within "imgdownloader".
	downloadPath := filepath.Join("examples", "imgdownloader", "downloads")
	err := os.MkdirAll(downloadPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create downloads folder: %v", err)
	}

	// Download the file.
	resp, err := http.Get(t.URL)
	if err != nil {
		return fmt.Errorf("failed to download file from %s: %v", t.URL, err)
	}
	defer resp.Body.Close()

	// Return error if we get status code other than 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file from %s --- %v: File Not Found", t.URL, resp.StatusCode)
	}

	// Create the file inside the "downloads" folder.
	filePath := filepath.Join(downloadPath, t.FileName)
	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer out.Close()

	// Write the response body to file.
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %v", filePath, err)
	}

	fmt.Printf("File downloaded and saved: %s\n", filePath)
	return nil
}
func main() {
	urlList := []string{
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexoto-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/downloeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo57.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg&fm=jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
		"https://wallpapers.com/images/featured/4k-gaming-33vov45f7zqi6t75.jpg",
		"https://asset.gecdesigns.com/img/wallpapers/beautiful-fantasy-wallpaper-ultra-hd-wallpaper-4k-sr10012418-1706506236698-cover.webp",
		"https://cdn.hero.page/wallpapers/6d7193c5-9ef1-45dc-a154-7d82180887f3-cosmic-galleries-interstellar-vistas-wallpaper-1.png",
		"https://images.wallpapersden.com/image/download/hellsweeper-vr-4k-gaming_bmZnaWqUmZqaraWkpJRqZmdlrWdtbWU.jpg",
		"https://images.pexels.com/photos/1366957/pexels-photo-1366957.jpeg?cs=srgb&dl=pexels-iriser-1366957.jpg",
	}

	var numOfWorkers, maxRetries uint = 3, 2

	// Create and start the worker pool.
	wp := tqwp.New(&tqwp.WorkerPoolConfig{
		NumOfWorkers: numOfWorkers,
		MaxRetries:   maxRetries,
		QueueSize:    10,
	})
	defer wp.Summary()
	defer wp.Stop()

	wp.Start()

	// Enqueue a task for each file download.
	for i, url := range urlList {
		fileName := fmt.Sprintf("file%d.jpg", i+1)

		task := &FileDownloadTask{
			Id:       uint(i + 1),
			URL:      url,
			FileName: fileName,
		}
		wp.EnqueueTask(task)
	}

}
