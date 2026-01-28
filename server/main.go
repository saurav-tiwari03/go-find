package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Install script that downloads the pre-built binary
const installScript = `#!/bin/bash

set -e

# go-find installer
# Usage: curl -skL https://go-find.sauravdev.in | bash

REPO_URL="https://go-find.sauravdev.in"
INSTALL_DIR="${HOME}/.local/bin"
BINARY_NAME="go-find"

echo ""
echo " ██████╗  ██████╗       ███████╗██╗███╗   ██╗██████╗ "
echo "██╔════╝ ██╔═══██╗      ██╔════╝██║████╗  ██║██╔══██╗"
echo "██║  ███╗██║   ██║█████╗█████╗  ██║██╔██╗ ██║██║  ██║"
echo "██║   ██║██║   ██║╚════╝██╔══╝  ██║██║╚██╗██║██║  ██║"
echo "╚██████╔╝╚██████╔╝      ██║     ██║██║ ╚████║██████╔╝"
echo " ╚═════╝  ╚═════╝       ╚═╝     ╚═╝╚═╝  ╚═══╝╚═════╝ "
echo ""
echo "Fast, minimal file explorer written in Go"
echo "─────────────────────────────────────────"
echo ""

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    mingw*|msys*|cygwin*) OS="windows" ;;
    *)
        echo "❌ Unsupported OS: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7l|armv6l) ARCH="arm" ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

DOWNLOAD_URL="${REPO_URL}/download/${OS}/${ARCH}/go-find"

echo "📋 Detected: ${OS}/${ARCH}"
echo "📥 Downloading go-find..."

# Create temp directory
TEMP_DIR=$(mktemp -d)
cleanup() {
    rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

# Download binary
if command -v curl &> /dev/null; then
    HTTP_CODE=$(curl -skL -w "%{http_code}" -o "${TEMP_DIR}/${BINARY_NAME}" "${DOWNLOAD_URL}")
    if [ "$HTTP_CODE" != "200" ]; then
        echo "❌ Failed to download (HTTP ${HTTP_CODE})"
        echo "   URL: ${DOWNLOAD_URL}"
        exit 1
    fi
elif command -v wget &> /dev/null; then
    wget -q -O "${TEMP_DIR}/${BINARY_NAME}" "${DOWNLOAD_URL}" || {
        echo "❌ Failed to download"
        exit 1
    }
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make executable
chmod +x "${TEMP_DIR}/${BINARY_NAME}"

# Check if user wants to install or just run
if [ "$1" = "install" ]; then
    echo "📦 Installing to ${INSTALL_DIR}..."
    mkdir -p "$INSTALL_DIR"
    mv "${TEMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Check if INSTALL_DIR is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo ""
        echo "⚠️  Add this to your shell profile (.bashrc, .zshrc, etc.):"
        echo "   export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
    fi
    
    echo "✅ Installed! Run 'go-find' to use."
else
    # Just run it
    echo "▶️  Running go-find..."
    echo ""
    TARGET_DIR="${1:-.}"
    "${TEMP_DIR}/${BINARY_NAME}" "$TARGET_DIR"
fi

echo ""
echo "─────────────────────────────────────────"
echo "💡 Tip: Run 'curl -skL ${REPO_URL} | bash -s install' to install permanently"
echo ""
`

// Landing page HTML
const landingPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="icon" type="favicon" href="/favicon.ico">
    <title>go-find | Fast File Explorer</title>
    <meta name="description" content="Fast, minimal file explorer written in Go. View directory trees with file sizes, colorized output, and instant results.">
    
    <!-- Open Graph / Facebook / WhatsApp -->
    <meta property="og:type" content="website">
    <meta property="og:url" content="https://go-find.sauravdev.in/">
    <meta property="og:title" content="go-find | Fast File Explorer">
    <meta property="og:description" content="Fast, minimal file explorer written in Go. View directory trees with file sizes, colorized output, and instant results. Install with: curl -skL go-find.sauravdev.in | bash">
    <meta property="og:image" content="https://go-find.sauravdev.in/og-image.svg">
    <meta property="og:site_name" content="go-find">
    
    <!-- Twitter -->
    <meta name="twitter:card" content="summary_large_image">
    <meta name="twitter:url" content="https://go-find.sauravdev.in/">
    <meta name="twitter:title" content="go-find | Fast File Explorer">
    <meta name="twitter:description" content="Fast, minimal file explorer written in Go. View directory trees with file sizes, colorized output, and instant results.">
    <meta name="twitter:image" content="https://go-find.sauravdev.in/og-image.svg">
    
    <!-- Additional meta -->
    <meta name="author" content="saurav-tiwari03">
    <meta name="keywords" content="go, golang, file explorer, tree, cli, terminal, command line">
    <meta name="theme-color" content="#00ADD8">
    
    <link rel="icon" type="image/x-icon" href="/favicon.ico">
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;600&family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: linear-gradient(135deg, #0d0d0d 0%, #1a1a2e 50%, #0d0d0d 100%);
            min-height: 100vh;
            display: flex;
            flex-direction: column;
            align-items: center;
            justify-content: center;
            color: #fff;
            padding: 20px;
        }

        .container {
            text-align: center;
            max-width: 600px;
        }

        .logo {
            margin-bottom: 20px;
        }

        .logo svg {
            width: 80px;
            height: 80px;
        }

        h1 {
            font-family: 'JetBrains Mono', monospace;
            font-size: 3rem;
            font-weight: 700;
            background: linear-gradient(135deg, #00ADD8 0%, #5DC9E2 50%, #00ADD8 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            margin-bottom: 12px;
            letter-spacing: -1px;
        }

        .subtitle {
            font-size: 1.1rem;
            color: #888;
            margin-bottom: 40px;
            font-weight: 400;
        }

        .features {
            display: flex;
            gap: 24px;
            justify-content: center;
            margin-bottom: 40px;
            flex-wrap: wrap;
        }

        .feature {
            display: flex;
            align-items: center;
            gap: 8px;
            color: #aaa;
            font-size: 0.9rem;
        }

        .feature-icon {
            font-size: 1.2rem;
        }

        .divider {
            color: #444;
            font-size: 0.85rem;
            margin: 30px 0;
            display: flex;
            align-items: center;
            gap: 16px;
        }

        .divider::before,
        .divider::after {
            content: '';
            flex: 1;
            height: 1px;
            background: linear-gradient(90deg, transparent, #333, transparent);
        }

        .code-block {
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 12px;
            padding: 16px 24px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 16px;
            margin-bottom: 16px;
            backdrop-filter: blur(10px);
        }

        .code-block code {
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.95rem;
            color: #e0e0e0;
            user-select: all;
        }

        .copy-btn {
            background: rgba(255, 255, 255, 0.1);
            border: 1px solid rgba(255, 255, 255, 0.15);
            color: #aaa;
            padding: 8px 16px;
            border-radius: 8px;
            cursor: pointer;
            font-family: 'Inter', sans-serif;
            font-size: 0.85rem;
            transition: all 0.2s ease;
        }

        .copy-btn:hover {
            background: rgba(255, 255, 255, 0.15);
            color: #fff;
        }

        .copy-btn.copied {
            background: rgba(0, 173, 216, 0.2);
            border-color: rgba(0, 173, 216, 0.3);
            color: #00ADD8;
        }

        .note {
            color: #555;
            font-size: 0.85rem;
            margin-top: 8px;
        }

        .footer {
            margin-top: 60px;
            display: flex;
            gap: 24px;
            align-items: center;
        }

        .footer a {
            color: #666;
            text-decoration: none;
            font-size: 0.9rem;
            transition: color 0.2s ease;
        }

        .footer a:hover {
            color: #00ADD8;
        }

        .footer-divider {
            color: #333;
        }

        .install-options {
            margin-top: 20px;
        }

        .install-label {
            color: #666;
            font-size: 0.8rem;
            margin-bottom: 8px;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .usage-section {
            margin-top: 40px;
        }

        .usage-grid {
            display: flex;
            gap: 16px;
            justify-content: center;
            flex-wrap: wrap;
        }

        .usage-item {
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            padding: 16px 20px;
            display: flex;
            flex-direction: column;
            gap: 8px;
            min-width: 200px;
        }

        .usage-item code {
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.9rem;
            color: #00ADD8;
        }

        .usage-desc {
            font-size: 0.8rem;
            color: #666;
        }

        @media (max-width: 600px) {
            h1 {
                font-size: 2.2rem;
            }

            .code-block {
                flex-direction: column;
                padding: 16px;
            }

            .code-block code {
                font-size: 0.8rem;
            }

            .features {
                flex-direction: column;
                gap: 12px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <svg viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
                <circle cx="50" cy="50" r="45" fill="#00ADD8" fill-opacity="0.1" stroke="#00ADD8" stroke-width="2"/>
                <path d="M30 50C30 38.954 38.954 30 50 30" stroke="#00ADD8" stroke-width="4" stroke-linecap="round"/>
                <path d="M50 30C61.046 30 70 38.954 70 50C70 61.046 61.046 70 50 70" stroke="#5DC9E2" stroke-width="4" stroke-linecap="round"/>
                <circle cx="50" cy="50" r="8" fill="#00ADD8"/>
            </svg>
        </div>

        <h1>go-find</h1>
        <p class="subtitle">fast, minimal file explorer written in Go</p>

        <div class="features">
            <div class="feature">
                <span class="feature-icon">📁</span>
                <span>Tree View</span>
            </div>
            <div class="feature">
                <span class="feature-icon">📊</span>
                <span>File Sizes</span>
            </div>
            <div class="feature">
                <span class="feature-icon">🎨</span>
                <span>Colorized</span>
            </div>
            <div class="feature">
                <span class="feature-icon">⚡</span>
                <span>Instant</span>
            </div>
        </div>

        <div class="divider">run in terminal</div>

        <div class="code-block">
            <code>curl -skL go-find.sauravdev.in | bash</code>
            <button class="copy-btn" onclick="copyCommand(this, 'curl -skL go-find.sauravdev.in | bash')">copy</button>
        </div>

        <div class="install-options">
            <div class="install-label">or install permanently</div>
            <div class="code-block">
                <code>curl -skL go-find.sauravdev.in | bash -s install</code>
                <button class="copy-btn" onclick="copyCommand(this, 'curl -skL go-find.sauravdev.in | bash -s install')">copy</button>
            </div>
        </div>

        <p class="note">works on macOS, Linux & Windows (WSL)</p>

        <div class="usage-section">
            <div class="divider">after installation</div>
            <div class="usage-grid">
                <div class="usage-item">
                    <code>go-find</code>
                    <span class="usage-desc">scan current directory</span>
                </div>
                <div class="usage-item">
                    <code>go-find ~/projects</code>
                    <span class="usage-desc">scan specific directory</span>
                </div>
            </div>
        </div>

        <div class="footer">
            <a href="https://github.com/saurav-tiwari03/go-find" target="_blank">github</a>
            <span class="footer-divider">·</span>
            <a href="https://sauravdev.in?utm_source=go-find" target="_blank">sauravdev.in</a>
        </div>
    </div>

    <script>
        function copyCommand(btn, text) {
            navigator.clipboard.writeText(text).then(() => {
                btn.textContent = 'copied!';
                btn.classList.add('copied');
                setTimeout(() => {
                    btn.textContent = 'copy';
                    btn.classList.remove('copied');
                }, 2000);
            });
        }
    </script>
</body>
</html>`

// Go gopher favicon as SVG
const faviconSVG = `<svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
  <circle cx="50" cy="50" r="45" fill="#00ADD8"/>
  <circle cx="50" cy="50" r="40" fill="#5DC9E2"/>
  <ellipse cx="35" cy="42" rx="8" ry="10" fill="white"/>
  <ellipse cx="65" cy="42" rx="8" ry="10" fill="white"/>
  <circle cx="37" cy="42" r="4" fill="#333"/>
  <circle cx="67" cy="42" r="4" fill="#333"/>
  <ellipse cx="50" cy="60" rx="12" ry="8" fill="#F6D2A2"/>
  <path d="M44 58 Q50 68 56 58" stroke="#333" stroke-width="2" fill="none" stroke-linecap="round"/>
</svg>`

// Open Graph image for social sharing (1200x630 recommended)
const ogImageSVG = `<svg width="1200" height="630" viewBox="0 0 1200 630" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <linearGradient id="bg" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" style="stop-color:#0d0d0d"/>
      <stop offset="50%" style="stop-color:#1a1a2e"/>
      <stop offset="100%" style="stop-color:#0d0d0d"/>
    </linearGradient>
    <linearGradient id="textGrad" x1="0%" y1="0%" x2="100%" y2="0%">
      <stop offset="0%" style="stop-color:#00ADD8"/>
      <stop offset="50%" style="stop-color:#5DC9E2"/>
      <stop offset="100%" style="stop-color:#00ADD8"/>
    </linearGradient>
  </defs>
  
  <!-- Background -->
  <rect width="1200" height="630" fill="url(#bg)"/>
  
  <!-- Logo circle -->
  <circle cx="600" cy="200" r="70" fill="#00ADD8" fill-opacity="0.1" stroke="#00ADD8" stroke-width="3"/>
  <path d="M550 200C550 172.386 572.386 150 600 150" stroke="#00ADD8" stroke-width="6" stroke-linecap="round" fill="none"/>
  <path d="M600 150C627.614 150 650 172.386 650 200C650 227.614 627.614 250 600 250" stroke="#5DC9E2" stroke-width="6" stroke-linecap="round" fill="none"/>
  <circle cx="600" cy="200" r="12" fill="#00ADD8"/>
  
  <!-- Title -->
  <text x="600" y="340" font-family="monospace" font-size="72" font-weight="bold" fill="url(#textGrad)" text-anchor="middle">go-find</text>
  
  <!-- Subtitle -->
  <text x="600" y="400" font-family="sans-serif" font-size="28" fill="#888888" text-anchor="middle">fast, minimal file explorer written in Go</text>
  
  <!-- Command box -->
  <rect x="300" y="450" width="600" height="60" rx="12" fill="rgba(255,255,255,0.05)" stroke="rgba(255,255,255,0.1)" stroke-width="1"/>
  <text x="600" y="490" font-family="monospace" font-size="22" fill="#e0e0e0" text-anchor="middle">curl -skL go-find.sauravdev.in | bash</text>
  
  <!-- Features -->
  <text x="350" y="570" font-family="sans-serif" font-size="18" fill="#666666" text-anchor="middle">📁 Tree View</text>
  <text x="500" y="570" font-family="sans-serif" font-size="18" fill="#666666" text-anchor="middle">📊 File Sizes</text>
  <text x="650" y="570" font-family="sans-serif" font-size="18" fill="#666666" text-anchor="middle">🎨 Colorized</text>
  <text x="800" y="570" font-family="sans-serif" font-size="18" fill="#666666" text-anchor="middle">⚡ Instant</text>
</svg>`

// Go favicon as ICO (base64 encoded simple 16x16 icon)
// This is a simple blue circle representing Go
var faviconICO []byte

func init() {
	// Simple 16x16 ICO file with Go blue color
	icoData := `AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAAMMOAADDDgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/AAAA/wAAAP8AAAD/AAAA/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/AAAAtQCtrf8ArKz/AK2t/wCsrP8AtbX/AAAA/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP8ArKz/AK2t/wCtrf8Ara3/AK2t/wCsrP8AAAD/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/AKys/wCtrf8Ara3/AK2t/wCtrf8Ara3/AKys/wAAAP8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA/wCsrP8Ara3/AK2t/wCtrf8Ara3/AK2t/wCsrP8AAAD/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP8ArKz/AK2t/wCtrf8Ara3/AK2t/wCtrf8ArKz/AAAA/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/AAAA/wCsrP8Ara3/AK2t/wCtrf8ArKz/AAAA/wAAAP8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/AAAA/wAAAP8AAAD/AAAA/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=`
	faviconICO, _ = base64.StdEncoding.DecodeString(icoData)
}

// Check if request is from a browser or curl/wget
func isBrowser(r *http.Request) bool {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))
	accept := r.Header.Get("Accept")

	// Check if it's curl, wget, or other CLI tools
	if strings.Contains(userAgent, "curl") ||
		strings.Contains(userAgent, "wget") ||
		strings.Contains(userAgent, "httpie") ||
		strings.Contains(userAgent, "fetch") {
		return false
	}

	// Check if Accept header indicates HTML preference
	if strings.Contains(accept, "text/html") {
		return true
	}

	// Check for common browser user agents
	if strings.Contains(userAgent, "mozilla") ||
		strings.Contains(userAgent, "chrome") ||
		strings.Contains(userAgent, "safari") ||
		strings.Contains(userAgent, "firefox") ||
		strings.Contains(userAgent, "edge") {
		return true
	}

	return false
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Serve favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Write(faviconICO)
	})

	http.HandleFunc("/favicon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fmt.Fprint(w, faviconSVG)
	})

	// Serve Open Graph image for social sharing
	http.HandleFunc("/og-image.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fmt.Fprint(w, ogImageSVG)
	})

	http.HandleFunc("/og-image.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fmt.Fprint(w, ogImageSVG)
	})

	// Serve landing page or install script at root based on client
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/install" && r.URL.Path != "/install.sh" {
			http.NotFound(w, r)
			return
		}

		// Force install script for specific paths
		if r.URL.Path == "/install" || r.URL.Path == "/install.sh" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("Content-Disposition", "inline")
			fmt.Fprint(w, installScript)
			log.Printf("Served install script to %s", r.RemoteAddr)
			return
		}

		// Serve based on client type
		if isBrowser(r) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprint(w, landingPage)
			log.Printf("Served landing page to %s", r.RemoteAddr)
		} else {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("Content-Disposition", "inline")
			fmt.Fprint(w, installScript)
			log.Printf("Served install script to %s", r.RemoteAddr)
		}
	})

	// Serve pre-built binaries
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/download/", http.FileServer(http.Dir("/app/binaries"))).ServeHTTP(w, r)
	})

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	log.Printf("🚀 go-find server starting on port %s", port)
	log.Printf("📋 Install: curl -sL http://localhost:%s | bash", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
