package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	totalSize    int64
	totalFiles   int
	totalFolders int
)

func banner() {
	fmt.Println("\033[38;2;139;92;246m" + `
 ██████╗  ██████╗       ███████╗██╗███╗   ██╗██████╗ 
██╔════╝ ██╔═══██╗      ██╔════╝██║████╗  ██║██╔══██╗
██║  ███╗██║   ██║█████╗█████╗  ██║██╔██╗ ██║██║  ██║
██║   ██║██║   ██║╚════╝██╔══╝  ██║██║╚██╗██║██║  ██║
╚██████╔╝╚██████╔╝      ██║     ██║██║ ╚████║██████╔╝
 ╚═════╝  ╚═════╝       ╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝` + "\033[0m")
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

// Get the last git commit date/time for the given directory, if it's a git repo.
func getLastCommitDate(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "log", "-1", "--format=%cd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func confirmStart(scanRoot string) bool {
	color.Yellow("This can take a while depending on project size.")
	color.Yellow("Type START to begin scanning node_modules under: %s", scanRoot)
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	line = strings.TrimSpace(line)
	return line == "START"
}

/* -------------------- tree logic -------------------- */

func tree(path string, prefix string) {
	// If the current folder is node_modules: calculate its size, print it in red, print last-commit, and return
	if filepath.Base(path) == "node_modules" {
		// compute total size of node_modules recursively
		var nmSize int64 = 0
		filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				nmSize += info.Size()
			}
			return nil
		})
		color.Red("%s[skipping node_modules/] (total size: %s)", prefix, humanSize(nmSize))
		lastCommit, err := getLastCommitDate(path)
		if err == nil && len(lastCommit) > 0 {
			color.Red("%sLast commit: %s", prefix, lastCommit)
		} else {
			color.Red("%sLast commit: (not a git repo or error)", prefix)
		}
		totalSize += nmSize
		return
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	// directories first, files later
	var dirs, files []os.DirEntry
	for _, e := range entries {
		// Skip .git directory
		if e.IsDir() && e.Name() == ".git" {
			continue
		}
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
			// If the entry is node_modules, handle as above in this spot
			if entry.Name() == "node_modules" {
				var nmSize int64 = 0
				filepath.Walk(fullPath, func(_ string, info os.FileInfo, err error) error {
					if err != nil {
						return nil
					}
					if !info.IsDir() {
						nmSize += info.Size()
					}
					return nil
				})
				color.Red("%s%s[skipping %s/] (total size: %s)", prefix, connector, entry.Name(), humanSize(nmSize))
				lastCommit, err := getLastCommitDate(fullPath)
				if err == nil && len(lastCommit) > 0 {
					color.Red("%s%sLast commit: %s", prefix, connector, lastCommit)
				} else {
					color.Red("%s%sLast commit: (not a git repo or error)", prefix, connector)
				}
				totalFolders++
				totalSize += nmSize
				continue
			}
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
	color.Cyan("go-find v1.0.0")
	color.HiBlack("Fast, minimal file explorer written in Go")
	color.HiBlack("────────────────────────────────────────")
	color.HiBlack(`Commands: 
	  1. Start : start the scan process
	  2. Exit  : exit the program`)
}

/* -------------------- node_modules mode -------------------- */

type nodeModulesEntry struct {
	ID        int    `json:"id"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
}

type nodeModulesCache struct {
	Version   int                `json:"version"`
	Root      string             `json:"root"`
	CreatedAt time.Time          `json:"created_at"`
	Entries   []nodeModulesEntry `json:"entries"`
}

func cacheFilePath() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil || base == "" {
		return "", fmt.Errorf("cannot determine user cache dir")
	}
	dir := filepath.Join(base, "go-find")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "node_modules_cache.json"), nil
}

func loadNodeModulesCache(path string) (*nodeModulesCache, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c nodeModulesCache
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	if c.Version == 0 {
		c.Version = 1
	}
	return &c, nil
}

func saveNodeModulesCache(path string, c *nodeModulesCache) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	return os.WriteFile(path, b, 0o644)
}

func isLikelyIDList(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if strings.HasPrefix(s, "[") && strings.Contains(s, "]") {
		return true
	}
	if strings.Contains(s, ",") {
		return true
	}
	// allow single id like "7" or "id7"
	if strings.HasPrefix(strings.ToLower(s), "id") {
		_, err := strconv.Atoi(strings.TrimPrefix(strings.ToLower(s), "id"))
		return err == nil
	}
	_, err := strconv.Atoi(s)
	return err == nil
}

func parseIDs(s string) ([]int, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("empty id list")
	}
	parts := strings.Split(s, ",")
	out := make([]int, 0, len(parts))
	seen := make(map[int]bool, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		pLower := strings.ToLower(p)
		if strings.HasPrefix(pLower, "id") {
			p = strings.TrimPrefix(pLower, "id")
		}
		n, err := strconv.Atoi(p)
		if err != nil || n <= 0 {
			return nil, fmt.Errorf("invalid id: %q", p)
		}
		if !seen[n] {
			seen[n] = true
			out = append(out, n)
		}
	}
	sort.Ints(out)
	return out, nil
}

func dirSize(root string) int64 {
	var total int64
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// avoid symlink loops
		if d.Type()&os.ModeSymlink != 0 {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total
}

func scanNodeModules(root string) ([]nodeModulesEntry, error) {
	var entries []nodeModulesEntry
	id := 1

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == ".git" {
			return filepath.SkipDir
		}
		if d.Type()&os.ModeSymlink != 0 {
			return filepath.SkipDir
		}
		if name == "node_modules" {
			sz := dirSize(path)
			entries = append(entries, nodeModulesEntry{
				ID:        id,
				Path:      path,
				SizeBytes: sz,
			})
			id++
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// big ones first is usually most useful
	sort.Slice(entries, func(i, j int) bool { return entries[i].SizeBytes > entries[j].SizeBytes })
	for i := range entries {
		entries[i].ID = i + 1
	}
	return entries, nil
}

type nodeModulesOpts struct {
	root       string
	read       bool
	del        bool
	selectMode bool
	ids        []int
}

func parseNodeModulesArgs(args []string) (nodeModulesOpts, error) {
	opts := nodeModulesOpts{root: "."}

	for _, a := range args {
		switch a {
		case "--read":
			opts.read = true
		case "--delete":
			opts.del = true
		case "--select":
			opts.selectMode = true
		default:
			if strings.HasPrefix(a, "--") {
				return opts, fmt.Errorf("unknown flag: %s", a)
			}
			// non-flag: either root dir or id list
			if isLikelyIDList(a) {
				ids, err := parseIDs(a)
				if err != nil {
					return opts, err
				}
				opts.ids = ids
				continue
			}
			if opts.root == "." {
				opts.root = a
				continue
			}
			return opts, fmt.Errorf("unexpected argument: %s", a)
		}
	}

	if opts.selectMode && !opts.del {
		return opts, fmt.Errorf("--select is only supported with --delete")
	}
	if opts.selectMode && opts.del && len(opts.ids) == 0 {
		return opts, fmt.Errorf("missing ids list (example: [1,2,3] or id1,id2)")
	}
	if !opts.read && !opts.del {
		// default to read if user only typed: node_modules
		opts.read = true
	}
	return opts, nil
}

func printNodeModules(entries []nodeModulesEntry) {
	if len(entries) == 0 {
		color.Yellow("No node_modules directories found.")
		return
	}
	color.Cyan("Found %d node_modules directories:", len(entries))
	for _, e := range entries {
		fmt.Printf("  [%d] %s (%s)\n", e.ID, e.Path, humanSize(e.SizeBytes))
	}
}

func deleteNodeModules(entries []nodeModulesEntry) (deleted []nodeModulesEntry) {
	for _, e := range entries {
		err := os.RemoveAll(e.Path)
		if err != nil {
			color.Red("❌ Failed: [%d] %s (%v)", e.ID, e.Path, err)
			continue
		}
		color.Green("✅ Deleted: [%d] %s", e.ID, e.Path)
		deleted = append(deleted, e)
	}
	return deleted
}

func runNodeModulesMode(args []string) {
	opts, err := parseNodeModulesArgs(args)
	if err != nil {
		color.Red("❌ %v", err)
		color.HiBlack("Usage examples:")
		color.HiBlack("  go-find node_modules --read")
		color.HiBlack("  go-find node_modules --delete")
		color.HiBlack("  go-find node_modules --select --delete id1,id2,id3")
		color.HiBlack("  go-find node_modules --select --delete '[" + `1,2,3` + "]'  (quote brackets in zsh)")
		os.Exit(2)
	}

	// Validate root exists for scanning paths
	rootInfo, err := os.Stat(opts.root)
	if err != nil || !rootInfo.IsDir() {
		color.Red("❌ Error: %s is not a directory", opts.root)
		os.Exit(2)
	}

	cachePath, err := cacheFilePath()
	if err != nil {
		color.Red("❌ Cache error: %v", err)
		os.Exit(1)
	}

	absRoot, _ := filepath.Abs(opts.root)

	// Select+delete uses cache to avoid rescanning.
	if opts.selectMode && opts.del {
		c, err := loadNodeModulesCache(cachePath)
		if err != nil {
			color.Red("❌ No cached scan found. Run: go-find node_modules --read")
			os.Exit(2)
		}
		if c.Root != "" && absRoot != "" && c.Root != absRoot {
			color.Yellow("Cached scan root differs from current root.")
			color.Yellow("Cached: %s", c.Root)
			color.Yellow("Current: %s", absRoot)
		}

		byID := make(map[int]nodeModulesEntry, len(c.Entries))
		for _, e := range c.Entries {
			byID[e.ID] = e
		}
		var selected []nodeModulesEntry
		for _, id := range opts.ids {
			e, ok := byID[id]
			if !ok {
				color.Red("❌ Unknown id: %d (run --read again to refresh IDs)", id)
				os.Exit(2)
			}
			selected = append(selected, e)
		}

		color.HiBlack("Deleting %d selected node_modules...", len(selected))
		_ = deleteNodeModules(selected)

		// remove deleted entries from cache
		deletedSet := make(map[int]bool, len(selected))
		for _, e := range selected {
			deletedSet[e.ID] = true
		}
		var remaining []nodeModulesEntry
		for _, e := range c.Entries {
			if !deletedSet[e.ID] {
				remaining = append(remaining, e)
			}
		}
		c.Entries = remaining
		_ = saveNodeModulesCache(cachePath, c)
		return
	}

	// Anything else needs scanning (read-only or delete-all).
	if !confirmStart(absRoot) {
		color.HiBlack("Cancelled.")
		return
	}

	color.HiBlack("Scanning for node_modules under: %s", absRoot)
	entries, err := scanNodeModules(opts.root)
	if err != nil {
		color.Red("❌ Scan error: %v", err)
		os.Exit(1)
	}

	printNodeModules(entries)

	c := &nodeModulesCache{
		Version:   1,
		Root:      absRoot,
		CreatedAt: time.Now(),
		Entries:   entries,
	}
	_ = saveNodeModulesCache(cachePath, c)

	if opts.del {
		color.HiBlack("\nDeleting all found node_modules...")
		_ = deleteNodeModules(entries)
		// after delete-all, clear cache (stale IDs)
		c.Entries = nil
		_ = saveNodeModulesCache(cachePath, c)
	}
}

/* -------------------- main -------------------- */

func main() {
	header()

	var command string
	fmt.Scanf("%s", &command)

	// node_modules mode: go-find node_modules [root] --read|--delete|--select --delete [ids]
	if len(os.Args) > 1 && os.Args[1] == "node_modules" {
		runNodeModulesMode(os.Args[2:])
		return
	}

	// Default tree mode: Get target directory from command-line argument or use current directory
	targetDir := "."
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	// Validate directory exists
	info, err := os.Stat(targetDir)
	if err != nil {
		color.Red("❌ Error: %v", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		color.Red("❌ Error: %s is not a directory", targetDir)
		os.Exit(1)
	}

	// If the starting path is node_modules, handle that scenario directly before proceeding
	if filepath.Base(targetDir) == "node_modules" {
		var nmSize int64 = 0
		filepath.Walk(targetDir, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				nmSize += info.Size()
			}
			return nil
		})
		color.Red("Scanning directory: %s (total size: %s)", targetDir, humanSize(nmSize))
		lastCommit, err := getLastCommitDate(targetDir)
		if err == nil && len(lastCommit) > 0 {
			color.Red("Last commit: %s", lastCommit)
		} else {
			color.Red("Last commit: (not a git repo or error)")
		}
		totalSize = nmSize
		// Don't tree-print anything else. Show summary directly:
		color.HiBlack("\n────────────────────────────────────────")
		color.Cyan("Summary")
		fmt.Printf("  Size     : %s\n", humanSize(totalSize))
		fmt.Printf("  Files    : %d\n", totalFiles)
		fmt.Printf("  Folders  : %d\n", totalFolders)
		color.HiBlack("\nDone ✔")
		return
	}

	color.HiBlack("Scanning directory: %s\n", targetDir)
	fmt.Println(targetDir)

	tree(targetDir, "")

	color.HiBlack("\n────────────────────────────────────────")
	color.Cyan("Summary")
	fmt.Printf("  Size     : %s\n", humanSize(totalSize))
	fmt.Printf("  Files    : %d\n", totalFiles)
	fmt.Printf("  Folders  : %d\n", totalFolders)

	color.HiBlack("\nDone ✔")
}
