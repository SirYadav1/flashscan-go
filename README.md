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
- **Beautiful UI**: Modern, colorful, and adaptive terminal interface
- **Dynamic Sizing**: Automatically adjusts to your screen size
- **Concurrent**: Scans thousands of hosts in seconds
- **Cross-platform**: Works on Windows, Linux, macOS
- **DNS Caching**: Highly optimized IP resolution

## Contributing
Contributions are welcome! Created by **SirYadav1**.

## Support
Join our [Telegram group](https://t.me/AdwanceSNI) for support.

## License
This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
