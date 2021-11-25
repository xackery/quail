package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	start := time.Now()
	err := run()
	if err != nil {
		fmt.Println("failed:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("finished in %0.2f seconds\n", time.Since(start).Seconds())
}

func run() error {
	return nil
}
