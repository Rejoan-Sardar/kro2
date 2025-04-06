package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	// Parse command line flags
	var (
		lspAddr    = flag.String("lsp-addr", "127.0.0.1:5001", "LSP server listen address")
		healthAddr = flag.String("addr", "127.0.0.1:5000", "Health server listen address")
		port       = flag.Int("port", 0, "Alternative way to specify port for the LSP server")
		healthPort = flag.Int("health-port", 0, "Alternative way to specify port for the health server")
		verbose    = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Configure logging
	log.SetPrefix("[kro-lsp] ")
	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(log.LstdFlags)
	}

	// Override lspAddr if port is specified
	if *port > 0 {
		*lspAddr = "127.0.0.1:" + itoa(*port)
	}

	// Override healthAddr if healthPort is specified
	if *healthPort > 0 {
		*healthAddr = "127.0.0.1:" + itoa(*healthPort)
	}

	// Start the server
	log.Printf("Starting Kro LSP Server (version %s)", ServerVersion)
	server := NewLSPServer(*lspAddr, *healthAddr)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
		os.Exit(1)
	}
}

// itoa converts an integer to a string without importing strconv
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	
	// Max int64 is 9223372036854775807, which has 19 digits
	digits := make([]byte, 0, 20)
	
	// Handle negative numbers
	negative := false
	if n < 0 {
		negative = true
		n = -n
	}
	
	// Generate digits in reverse order
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	
	// Add negative sign if needed
	if negative {
		digits = append(digits, '-')
	}
	
	// Reverse the slice
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	
	return string(digits)
}
