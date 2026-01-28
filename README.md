# go-find

A command-line utility written in Go that displays directory structures in a tree format with file sizes and colorized output.

## Features

- 📁 **Tree View Display**: Shows directory structure in a hierarchical tree format
- 📊 **File Size Information**: Displays human-readable file sizes (B, KB, MB, GB, TB, PB, EB)
- 🎨 **Colorized Output**: Color-coded output for better readability
  - Directories in blue
  - Files in white
  - Size information in light black
- 📈 **Statistics**: Tracks total files, folders, and combined size
- 🎯 **Clean Interface**: ASCII tree connectors for easy visualization

## Installation

### Quick Install (Recommended)

```bash
curl -sL https://go-find.sauravdev.in | bash
```

### Manual Installation

```bash
git clone https://github.com/saurav-tiwari03/go-find.git
cd go-find
go build -o go-find
./go-find
```

### From Source

```bash
git clone https://github.com/saurav-tiwari03/go-find.git
cd go-find
go install
```

## Usage

After installation, simply run:

```bash
go-find
```

Or to scan a specific directory:

```bash
go-find /path/to/directory
```

**Example:**
```bash
go-find .
```

This will display the directory structure with file sizes.

## Dependencies

- [fatih/color](https://github.com/fatih/color) - For colored terminal output

Install dependencies:
```bash
go mod download
```

## Docker Deployment

### Using EasyPanel

1. Connect your GitHub repository to EasyPanel
2. EasyPanel will automatically detect the Dockerfile
3. Deploy with a single click

### Building Docker Image Locally

```bash
docker build -t go-find:latest .
```

### Running with Docker

```bash
# Scan current directory
docker run --rm -v $(pwd):/scan go-find:latest /scan

# Scan specific directory
docker run --rm -v /path/to/scan:/scan go-find:latest /scan
```

### Docker Compose

Create a `docker-compose.yml`:
```yaml
version: '3.8'
services:
  go-find:
    build: .
    volumes:
      - /path/to/scan:/scan
    command: /scan
```

Run with:
```bash
docker-compose up
```

## How It Works

1. **Banner Display**: Shows the project banner on startup
2. **Directory Traversal**: Recursively walks through directories
3. **Organization**: Displays directories first, then files
4. **Size Calculation**: Computes and displays human-readable file sizes
5. **Statistics**: Accumulates total files, folders, and size information

## Project Structure

```
go-find/
├── main.go          # Main application file
├── go.mod           # Go module definition
└── README.md        # Project documentation
```

## Functions

- `banner()` - Displays the ASCII art banner
- `iconDecide(isDir bool)` - Returns appropriate icon (📁 for directory, 📄 for file)
- `humanSize(bytes int64)` - Converts byte size to human-readable format
- `sizeCalc(path string)` - Calculates file size and updates total
- `tree(path string, prefix string)` - Recursively displays directory tree

## Go Version

Requires Go 1.24.1 or later

## License

[Add your license here]

## Author

[saurav-tiwari03](https://github.com/saurav-tiwari03)

## Contributing

Contributions are welcome! Feel free to open issues and pull requests.
