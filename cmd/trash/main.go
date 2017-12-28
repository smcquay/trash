package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"mcquay.me/trash"
)

var algo = flag.String("a", "caca", "algorithm to use")

func main() {
	flag.Parse()
	var r io.Reader
	switch *algo {
	case "lo", "low", "00", "0", "nil", "null", "zeros":
		r = trash.Zeros
	case "ones", "ff", "hi", "high":
		r = trash.Fs
	case "trash", "caca":
		r = trash.Reader
	case "hilo", "aa", "a":
		r = trash.HiLo
	case "lohi", "55", "5":
		r = trash.LoHi
	case "rand", "random":
		r = trash.Random
	default:
		fmt.Fprintf(os.Stderr, "unsupported algorithm: %v\ntry one of 'lo(w)', 'hi(gh)', 'hilo', 'lohi', 'trash'\n", *algo)
		os.Exit(1)
	}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		fmt.Fprintf(os.Stderr, "problem copying to stdout: %v", err)
		os.Exit(1)
	}
}
