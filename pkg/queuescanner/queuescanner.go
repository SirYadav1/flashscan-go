package queuescanner

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/term"
)

// ANSI Color codes
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorBold    = "\033[1m"
)

type Ctx struct {
	ScanComplete int64
	SuccessCount int64
	startTime    int64
	lastStatTime int64
	statInterval int64 // in nanoseconds

	hostList     []string
	mu           sync.Mutex
	OutputFile   string
	lastResults  []string // Buffer for last N results
	resultsMutex sync.Mutex
	maxResults   int // Dynamic based on screen height
}

type QueueScanner struct {
	threads  int
	scanFunc func(c *Ctx, host string)
	queue    chan string
	wg       sync.WaitGroup
	ctx      *Ctx
}

func nowNano() int64 {
	return time.Now().UnixNano()
}

func formatETA(seconds float64) string {
	if seconds < 0 {
		return "--"
	}
	d := time.Duration(seconds * float64(time.Second))
	return d.Truncate(time.Second).String()
}

func formatDuration(seconds float64) string {
	d := time.Duration(seconds * float64(time.Second))
	return d.Truncate(time.Second).String()
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func printBanner() {
	banner := `  ______ _           _      _____                 
 |  ____| |         | |    / ____|                
 | |__  | | __ _ ___| |__ | (___   ___ __ _ _ __  
 |  __| | |/ _` + "`" + ` / __| '_ \ \___ \ / __/ _` + "`" + ` | '_ \ 
 | |    | | (_| \__ \ | | |____) | (_| (_| | | | |
 |_|    |_|\__,_|___/_| |_|_____/ \___\__,_|_| |_| v2.0`
	fmt.Printf("%s%s%s\n", ColorCyan+ColorBold, banner, ColorReset)
}

// Calculate dynamic max results based on terminal height
func getMaxResults() int {
	_, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 10 // Default fallback
	}

	// Banner: 4 lines
	// Progress box: 5 lines
	// Header: 2 lines
	// Table header: 3 lines
	// Footer: 2 lines
	// Total overhead: ~16 lines

	available := height - 16
	if available < 5 {
		return 5 // Minimum
	}
	if available > 50 {
		return 50 // Maximum
	}
	return available
}

// Add result to buffer
func (ctx *Ctx) Log(a ...any) {
	msg := fmt.Sprint(a...)

	ctx.resultsMutex.Lock()
	ctx.lastResults = append(ctx.lastResults, msg)
	// Keep only last N results
	if len(ctx.lastResults) > ctx.maxResults {
		ctx.lastResults = ctx.lastResults[1:]
	}
	ctx.resultsMutex.Unlock()
}

// Helper to strip ANSI codes for length calculation
func stripANSI(str string) string {
	// Simple ANSI code stripper
	result := ""
	inEscape := false
	for _, char := range str {
		if char == '\033' {
			inEscape = true
		} else if inEscape && char == 'm' {
			inEscape = false
		} else if !inEscape {
			result += string(char)
		}
	}
	return result
}

// Redraw entire screen with progress and results
func (ctx *Ctx) LogStat() {
	if ctx.statInterval > 0 {
		now := nowNano()
		if now-atomic.LoadInt64(&ctx.lastStatTime) < ctx.statInterval {
			return
		}
		atomic.StoreInt64(&ctx.lastStatTime, now)
	}

	scanSuccess := atomic.LoadInt64(&ctx.SuccessCount)
	scanComplete := atomic.LoadInt64(&ctx.ScanComplete)
	total := len(ctx.hostList)
	if total == 0 {
		return
	}
	failed := scanComplete - scanSuccess
	percentage := float64(scanComplete) / float64(total) * 100

	// Calculate stats
	elapsed := float64(nowNano()-ctx.startTime) / 1e9 // seconds
	var speed float64
	if elapsed > 0 {
		speed = float64(scanComplete) / elapsed
	}

	var etaSec float64
	if speed > 0 {
		remaining := float64(total - int(scanComplete))
		etaSec = remaining / speed
	}
	eta := formatETA(etaSec)

	// Progress bar - FIXED: prevent negative repeat counts
	barWidth := 40
	filled := int(percentage / 100 * float64(barWidth))
	if filled < 0 {
		filled = 0
	}
	if filled > barWidth {
		filled = barWidth
	}
	remaining := barWidth - filled
	if remaining < 0 {
		remaining = 0
	}
	bar := ColorGreen + strings.Repeat("━", filled) + ColorWhite + strings.Repeat("─", remaining) + ColorReset

	// Clear screen and redraw
	fmt.Print("\033[2J\033[H")

	// Banner
	printBanner()
	fmt.Println()

	// Progress box
	fmt.Printf("%s┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓%s\n", ColorBlue, ColorReset)
	
	// Progress bar line
	progressLine := fmt.Sprintf("┃ %s⚡ SCANNING... %s[%s] %s%.1f%%%s",
		ColorWhite+ColorBold, ColorReset, bar, ColorMagenta, percentage, ColorReset)
	progressLineStripped := stripANSI(progressLine)
	padding := 67 - len(progressLineStripped)
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s%s┃%s\n", ColorBlue, progressLine, strings.Repeat(" ", padding), ColorBlue, ColorReset)

	fmt.Printf("%s┠──────────────────────────────────────────────────────────────────┨%s\n", ColorBlue, ColorReset)

	// Stats line 1
	statsLine1 := fmt.Sprintf("┃ %s✔ Success: %s%-5d %s┃ %s✖ Failed: %s%-5d %s┃ %s🚀 Speed: %s%-6.0f %s┃",
		ColorGreen, ColorWhite, scanSuccess, ColorBlue,
		ColorRed, ColorWhite, failed, ColorBlue,
		ColorMagenta, ColorWhite, speed, ColorBlue)
	fmt.Printf("%s%s%s\n", ColorBlue, statsLine1, ColorReset)

	// Stats line 2
	statsLine2 := fmt.Sprintf("┃ %s⏱  ETA: %s%-12s %s┃ %s📂 Scanned: %s%d/%d %s┃                    ┃",
		ColorYellow, ColorWhite, eta, ColorBlue,
		ColorCyan, ColorWhite, scanComplete, total, ColorBlue)
	fmt.Printf("%s%s%s\n", ColorBlue, statsLine2, ColorReset)

	fmt.Printf("%s┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛%s\n", ColorBlue, ColorReset)
	fmt.Println()

	// Results table
	fmt.Printf("%s✅ LATEST %d RESULTS:%s\n", ColorGreen+ColorBold, ctx.maxResults, ColorReset)

	ctx.resultsMutex.Lock()
	resultsCount := len(ctx.lastResults)
	if resultsCount > 0 {
		for _, result := range ctx.lastResults {
			// Color code based on content
			if strings.Contains(result, "200") || strings.Contains(result, "✓") {
				fmt.Printf("%s%s%s\n", ColorGreen, result, ColorReset)
			} else if strings.Contains(result, "timeout") || strings.Contains(result, "failed") || strings.Contains(result, "✗") {
				fmt.Printf("%s%s%s\n", ColorRed, result, ColorReset)
			} else if strings.Contains(result, "301") || strings.Contains(result, "302") {
				fmt.Printf("%s%s%s\n", ColorYellow, result, ColorReset)
			} else {
				fmt.Printf("%s\n", result)
			}
		}
	} else {
		fmt.Printf("%sWaiting for results...%s\n", ColorCyan, ColorReset)
	}
	ctx.resultsMutex.Unlock()

	// Footer info
	if ctx.OutputFile != "" {
		fmt.Printf("\n%s💾 Results saved to:%s %s%s%s\n",
			ColorGreen, ColorReset, ColorCyan, ctx.OutputFile, ColorReset)
	}
}

// Print final summary
func (ctx *Ctx) PrintSummary() {
	total := len(ctx.hostList)
	if total == 0 {
		return
	}
	success := atomic.LoadInt64(&ctx.SuccessCount)
	failed := int64(total) - success
	elapsed := float64(nowNano()-ctx.startTime) / 1e9

	fmt.Print("\033[2J\033[H")

	fmt.Printf("\n%s╔══════════════════════════════════════════════════════════════════╗%s\n", ColorGreen+ColorBold, ColorReset)
	fmt.Printf("%s║                    📊 SCAN COMPLETED                             ║%s\n", ColorGreen+ColorBold, ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════╝%s\n\n", ColorGreen+ColorBold, ColorReset)

	fmt.Printf("%s📈 Statistics:%s\n", ColorBlue+ColorBold, ColorReset)
	fmt.Printf("   • Total Scanned: %s%d%s hosts\n", ColorMagenta, total, ColorReset)
	fmt.Printf("   • %sSuccessful:%s %s%d%s (%.1f%%)\n",
		ColorGreen, ColorReset, ColorGreen, success, ColorReset,
		float64(success)/float64(total)*100)
	fmt.Printf("   • %sFailed:%s %s%d%s (%.1f%%)\n",
		ColorRed, ColorReset, ColorRed, failed, ColorReset,
		float64(failed)/float64(total)*100)
	fmt.Printf("   • Time Elapsed: %s%s%s\n", ColorMagenta, formatDuration(elapsed), ColorReset)

	if elapsed > 0 {
		fmt.Printf("   • Average Speed: %s%.1f hosts/sec%s\n",
			ColorMagenta, float64(total)/elapsed, ColorReset)
	}

	if ctx.OutputFile != "" {
		fmt.Printf("\n%s💾 Results saved to:%s %s%s%s\n",
			ColorGreen, ColorReset, ColorCyan, ctx.OutputFile, ColorReset)
	}

	fmt.Printf("\n%s✨ Thank you for using FlashScan-Go! ✨%s\n\n",
		ColorGreen+ColorBold, ColorReset)
}

func (ctx *Ctx) ScanSuccess(result any) {
	if str, ok := result.(string); ok && ctx.OutputFile != "" {
		ctx.mu.Lock()
		file, err := os.OpenFile(ctx.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			file.WriteString(str + "\n")
			file.Close()
		}
		ctx.mu.Unlock()
	}

	atomic.AddInt64(&ctx.SuccessCount, 1)
}

func New(threads int, scanFunc func(c *Ctx, host string)) *QueueScanner {
	scanner := &QueueScanner{
		threads:  threads,
		scanFunc: scanFunc,
		queue:    make(chan string, threads*10), // Increased buffer
		ctx: &Ctx{
			maxResults:  getMaxResults(),
			lastResults: make([]string, 0),
		},
	}

	for i := 0; i < scanner.threads; i++ {
		scanner.wg.Add(1)
		go scanner.run()
	}

	return scanner
}

func (qs *QueueScanner) SetOptions(hostList []string, outputFile string, statInterval float64) {
	qs.ctx.hostList = hostList
	qs.ctx.OutputFile = outputFile
	qs.ctx.statInterval = int64(statInterval * 1e9)
	qs.ctx.maxResults = getMaxResults() // Update based on current terminal size
}

func (qs *QueueScanner) Start() {
	qs.ctx.startTime = nowNano()
	hideCursor()
	defer showCursor()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		showCursor()
		qs.ctx.PrintSummary()
		os.Exit(0)
	}()

	// Initial display
	qs.ctx.LogStat()

	for _, host := range qs.ctx.hostList {
		qs.queue <- host
	}
	close(qs.queue)

	qs.wg.Wait()

	// Final summary
	qs.ctx.PrintSummary()
}

func (qs *QueueScanner) run() {
	defer qs.wg.Done()

	for {
		host, ok := <-qs.queue
		if !ok {
			break
		}

		qs.scanFunc(qs.ctx, host)

		atomic.AddInt64(&qs.ctx.ScanComplete, 1)
		qs.ctx.LogStat()
	}
}
