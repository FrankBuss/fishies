// main.go
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"fishies/model"
	"fishies/render"

	"github.com/rwxrob/bonzai/anim"
)

type config struct {
	numFish int
	ground  bool
}

func parseFlags() config {
	cfg := config{}

	flag.IntVar(&cfg.numFish, "fish", 3, "number of fish to simulate")
	flag.IntVar(&cfg.numFish, "f", 3, "shorthand for --fish")
	flag.BoolVar(&cfg.ground, "ground", false, "enable ground plane")
	flag.BoolVar(&cfg.ground, "g", false, "shorthand for --ground")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if cfg.numFish <= 0 {
		fmt.Fprintf(os.Stderr, "Error: number of fish must be positive\n")
		flag.Usage()
		os.Exit(1)
	}

	return cfg
}

// getTerminalSize retrieves current terminal dimensions
func getTerminalSize() (width int, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// fallback to some reasonable defaults if we can't get the size
		return 80, 24
	}
	return width, height
}

// clamp ensures a value stays within given range
func clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// colorToANSI converts render.Color to ANSI terminal color string
func colorToANSI(c render.Color) string {
	r := int(c.R * 255)
	g := int(c.G * 255)
	b := int(c.B * 255)
	r = clamp(r, 0, 255)
	g = clamp(g, 0, 255)
	b = clamp(b, 0, 255)

	if r == 0 && g == 0 && b == 0 {
		return " " // Just space for black (eyes)
	}

	// Use X character with appropriate color
	return fmt.Sprintf("\033[38;2;%d;%d;%dm%s", r, g, b, "X")
}

func main() {
	cfg := parseFlags()

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	err := anim.SimpleAnimationScreen()
	if err != nil {
		fmt.Printf("Error initializing animation screen: %v\n", err)
		return
	}

	// Initialize string builder for frame buffer
	var frameBuilder strings.Builder
	frameBuilder.Grow(16384)

	scene, light := model.CreateScene(cfg.numFish, cfg.ground)

	lastTime := time.Now()
	for {
		currentTime := time.Now()
		deltaTime := float64(currentTime.Sub(lastTime).Seconds())
		lastTime = currentTime

		width, height := getTerminalSize()

		// Calculate dimensions maintaining aspect ratio
		smallestDimension := width
		if height*2 < smallestDimension {
			smallestDimension = height * 2
		}

		aspectRatio := float64(width) / float64(height*2)

		frameBuilder.Reset()
		frameBuilder.WriteString("\033[H")
		frameBuilder.WriteString("\033[0m")

		model.UpdateScene(scene, deltaTime)

		// Buffer for building output lines
		line := make([]string, width+1)
		line[width] = "\n"

		lineIndex := 0
		render.RenderScene(scene, light, width, height, aspectRatio, func(x, y int, color render.Color) {
			// Start new line when reaching width
			if x == 0 && y > 0 {
				for _, s := range line {
					frameBuilder.WriteString(s)
				}
				lineIndex = 0
			}
			line[lineIndex] = colorToANSI(color)
			lineIndex++
		})

		// Write final line
		for _, s := range line {
			frameBuilder.WriteString(s)
		}

		// Output the frame
		fmt.Print(frameBuilder.String())

		// Cap frame rate at ~30fps
		time.Sleep(33 * time.Millisecond)
	}
}
