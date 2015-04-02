package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type LinearGradient struct {
	// start
	x0, y0 float64
	// end
	x1, y1 float64
}

// Parse the line and create a new LinearGradient
func NewLinearGradient(line string) *LinearGradient {
	parts := strings.Split(line, " ")
	x0, _ := strconv.ParseFloat(parts[1], 64)
	y0, _ := strconv.ParseFloat(parts[2], 64)
	x1, _ := strconv.ParseFloat(parts[3], 64)
	y1, _ := strconv.ParseFloat(parts[4], 64)
	return &LinearGradient{x0: x0, y0: y0, x1: x1, y1: y1}
}

// Norm computes the norm of the gradient vector
func (lg *LinearGradient) norm() float64 {
	return math.Sqrt(math.Pow(lg.x1-lg.x0, 2) + math.Pow(lg.y1-lg.y0, 2))
}

// Dot computes the dot product of this gradient vector with the given one
func (lg *LinearGradient) dot(a, b float64) float64 {
	return a*(lg.x1-lg.x0) + b*(lg.y1-lg.y0)
}

// Color contains data regarding the colors characters
type Color struct {
	colors  string
	binSize float64
}

// Parse the line and create a new Color
func NewColor(line string) *Color {
	return &Color{colors: line}
}

// Color returns a color character depending on the given scalar value
func (c *Color) color(y float64) string {
	if y < 0 {
		bin := len(c.colors) + int(y/c.binSize) - 1
		if bin < 0 {
			return string(c.colors[0])
		}
		return string(c.colors[bin])
	}
	bin := int(y / c.binSize)
	if bin < len(c.colors) {
		return string(c.colors[bin])
	}
	return string(c.colors[len(c.colors)-1])
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Failed to open file", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// Get configuration
	width, height := float64(0), float64(0)
	var colors *Color
	var gradient *LinearGradient
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		switch i {
		case 0:
			parts := strings.Split(line, " ")
			w, _ := strconv.ParseInt(parts[0], 10, 32)
			h, _ := strconv.ParseInt(parts[1], 10, 32)
			width, height = float64(w), float64(h)
		case 1:
			colors = NewColor(line)
		case 2:
			switch {
			case strings.HasPrefix(line, "linear"):
				gradient = NewLinearGradient(scanner.Text())
			default:
				log.Fatalf("Unknown gradient: %s", line)
			}
		default:
			log.Fatal("Too many lines")
		}
	}
	if scanner.Err() != nil {
		log.Fatalf("Failed to parse config file", err)
	}

	norm := gradient.norm()
	colors.binSize = norm / float64(len(colors.colors))

	// Display the gradient
	for h := float64(0); h < height; h++ {
		for w := float64(0); w < width; w++ {
			projection := gradient.dot(w, h) / norm
			fmt.Print(colors.color(projection))
		}
		fmt.Println()
	}
}
