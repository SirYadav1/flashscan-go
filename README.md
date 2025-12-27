# FlashScan-Go

**High Performance SNI Bug Host Scanner** - A next-generation network scanner by **SirYadav1**.

[![Telegram](https://img.shields.io/badge/Telegram-Join%20Group-0088cc?style=flat-square&logo=telegram)](https://t.me/AdwanceSNI)
[![Go Version](https://img.shields.io/github/go-mod/go-version/SirYadav1/flashscan-go?style=flat-square)](https://github.com/SirYadav1/flashscan-go)
[![License](https://img.shields.io/github/license/SirYadav1/flashscan-go?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/SirYadav1/flashscan-go?style=flat-square)](https://github.com/SirYadav1/flashscan-go/releases)

## Installation

### Build from Source

```bash
go install github.com/SirYadav1/flashscan-go@latest
```

## Usage

### Quick Start
```bash
# Show all available commands
flashscan-go --help

# Direct scan
flashscan-go direct -f domains.txt -o output.txt

# CDN SSL scan
flashscan-go cdn-ssl -f proxies.txt --target example.com

# SNI scan with custom parameters
flashscan-go sni -f subdomains.txt --threads 128 --timeout 5
```

### Available Commands
- `sni`     - Scan Server Name Indication list from file
- `cdn-ssl` - Scan using CDN SSL proxy with payload injection
- `direct`  - Scan using direct connection to targets (Most accurate)
- `proxy`   - Scan using a proxy with payload
- `ping`    - Scan hosts using TCP ping

## Features
- **High Performance**: Optimized with DNS Caching & Buffer Pooling
- **Beautiful UI**: Modern, colorful, and adaptive terminal inter
### UI and UX Improvements
- [MODIFY] [pkg/queuescanner/queuescanner.go](file:///x:/project/AdwanceSNI-2.0/flashscan-go/pkg/queuescanner/queuescanner.go):
    - Enhance the banner with a more premium/modern ASCII art or border.
    - Improve the progress bar design (e.g., using more detailed characters or different colors).
    - Refine the results table layout for better readability on different terminal sizes.
    - Add micro-animations or smoother transitions for the stat updates.

### Documentation and CLI Consistency
- [MODIFY] [README.md](file:///x:/project/AdwanceSNI-2.0/flashscan-go/README.md):
    - Fix incorrect flag examples (e.g., `proxy-filename` should be `filename` or `-f`).
    - Update installation instructions if needed.
    - Ensure all commands described match the implementation in `cmd/`.
face
- **Dynamic Sizing**: Automatically adjusts to your screen size
- **Concurrent**: Scans thousands of hosts in seconds
- **Cross-platform**: Works on Windows, Linux, macOS

## Contributing
Contributions are welcome! Created by **SirYadav1**.

## Support
Join our [Telegram group](https://t.me/AdwanceSNI) for support.

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
