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
# Usage: curl -sL https://go-find.sauravdev.in | bash

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

DOWNLOAD_URL="${REPO_URL}/download/${OS}/${ARCH}"

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
    HTTP_CODE=$(curl -sL -w "%{http_code}" -o "${TEMP_DIR}/${BINARY_NAME}" "${DOWNLOAD_URL}")
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
echo "💡 Tip: Run 'curl -sL ${REPO_URL} | bash -s install' to install permanently"
echo ""
`

// Landing page HTML
const landingPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>go-find | Fast File Explorer</title>
    <link rel="icon" type="image/svg+xml" href="/favicon.svg">
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
            <code>curl -sL go-find.sauravdev.in | bash</code>
            <button class="copy-btn" onclick="copyCommand(this, 'curl -sL go-find.sauravdev.in | bash')">copy</button>
        </div>

        <div class="install-options">
            <div class="install-label">or install permanently</div>
            <div class="code-block">
                <code>curl -sL go-find.sauravdev.in | bash -s install</code>
                <button class="copy-btn" onclick="copyCommand(this, 'curl -sL go-find.sauravdev.in | bash -s install')">copy</button>
            </div>
        </div>

        <p class="note">works on macOS, Linux & Windows (WSL)</p>

        <div class="footer">
            <a href="https://github.com/saurav-tiwari03/go-find" target="_blank">github</a>
            <span class="footer-divider">·</span>
            <a href="https://sauravdev.in" target="_blank">sauravdev.in</a>
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
