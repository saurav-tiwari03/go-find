package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

var (
	totalSize    int64
	totalFiles   int
	totalFolders int
)

func banner() {
	fmt.Println(`
 ██████╗  ██████╗       ███████╗██╗███╗   ██╗██████╗ 
██╔════╝ ██╔═══██╗      ██╔════╝██║████╗  ██║██╔══██╗
██║  ███╗██║   ██║█████╗█████╗  ██║██╔██╗ ██║██║  ██║
██║   ██║██║   ██║╚════╝██╔══╝  ██║██║╚██╗██║██║  ██║
╚██████╔╝╚██████╔╝      ██║     ██║██║ ╚████║██████╔╝
 ╚═════╝  ╚═════╝       ╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝`)
}

/* -------------------- helpers -------------------- */

func iconDecide(isDir bool) string {
	if isDir {
		return "📁"
	}
	return "📄"
}

func humanSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func sizeCalc(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	if !info.IsDir() {
		totalSize += info.Size()
		return info.Size()
	}
	return 0
}

/* -------------------- tree logic -------------------- */

func tree(path string, prefix string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	// directories first, files later
	var dirs, files []os.DirEntry
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			files = append(files, e)
		}
	}
	entries = append(dirs, files...)

	for i, entry := range entries {
		isLast := i == len(entries)-1

		connector := "├── "
		nextPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			nextPrefix = prefix + "    "
		}

		fullPath := filepath.Join(path, entry.Name())
		icon := iconDecide(entry.IsDir())

		if entry.IsDir() {
			color.Blue("%s%s%s %s/", prefix, connector, icon, entry.Name())
			totalFolders++
			tree(fullPath, nextPrefix)
		} else {
			size := sizeCalc(fullPath)
			color.White("%s%s%s %s", prefix, connector, icon, entry.Name())
			color.HiBlack(" (%s)", humanSize(size))
			totalFiles++
		}
	}
}

/* -------------------- UI -------------------- */

func header() {
	banner()
	color.Cyan("go-find v0.1.0")
	color.HiBlack("Fast, minimal file explorer written in Go")
	color.HiBlack("────────────────────────────────────────")
}

/* -------------------- main -------------------- */

func main() {
	header()

	color.HiBlack("Scanning current directory...\n")
	fmt.Println(".")

	tree(".", "")

	color.HiBlack("\n────────────────────────────────────────")
	color.Cyan("Summary")
	fmt.Printf("  Size     : %s\n", humanSize(totalSize))
	fmt.Printf("  Files    : %d\n", totalFiles)
	fmt.Printf("  Folders  : %d\n", totalFolders)

	color.HiBlack("\nDone ✔")
}
