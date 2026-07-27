package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	xminFlag   = flag.Float64("xmin", -2.0, "minimum x coordinate")
	xmaxFlag   = flag.Float64("xmax", 2.0, "maximum x coordinate")
	yminFlag   = flag.Float64("ymin", -2.0, "minimum y coordinate")
	ymaxFlag   = flag.Float64("ymax", 2.0, "maximum y coordinate")
	widthFlag  = flag.Int("width", 8000, "image width")
	heightFlag = flag.Int("height", 8000, "image height")
	iterFlag   = flag.Int("iter", 500, "maximum iterations")
	outputFlag = flag.String("output", "mandelbrot.bmp", "output filename (supports .bmp, .tga, .png)")
	formatFlag = flag.String("format", "bmp", "output format: bmp, tga, png, raw")
)

func main() {
	flag.Parse()
	
	width := *widthFlag
	height := *heightFlag
	iterations := *iterFlag
	
	// Calculate expected file size
	uncompressedSize := width * height * 4 // RGBA
	fmt.Printf("Generating %dx%d image\n", width, height)
	fmt.Printf("Uncompressed size: %.2f MB\n", float64(uncompressedSize)/(1024*1024))
	
	// Different formats have different sizes
	var estimatedSize float64
	switch *formatFlag {
	case "bmp", "tga", "raw":
		estimatedSize = float64(uncompressedSize) / (1024 * 1024)
	case "png":
		estimatedSize = float64(uncompressedSize) / (1024 * 1024 * 2) // PNG ~50% compression
	default:
		estimatedSize = float64(uncompressedSize) / (1024 * 1024)
	}
	fmt.Printf("Estimated %s file size: %.2f MB\n", *formatFlag, estimatedSize)
	fmt.Println("Starting generation...")
	
	start := time.Now()
	
	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// Use all CPU cores
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	fmt.Printf("Using %d CPU cores\n", numCPU)
	
	// Create a wait group and a channel for rows
	var wg sync.WaitGroup
	rowChan := make(chan int, height)
	
	// Progress tracking
	progress := make(chan int, height)
	
	// Start workers
	for w := 0; w < numCPU; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for py := range rowChan {
				// Process a single row
				for px := 0; px < width; px++ {
					x := float64(px)/float64(width)*(*xmaxFlag-*xminFlag) + *xminFlag
					y := float64(py)/float64(height)*(*ymaxFlag-*yminFlag) + *yminFlag
					z := complex(x, y)
					
					// Calculate color for this pixel
					c := mandelbrot(z, iterations)
					img.Set(px, py, c)
				}
				progress <- py
			}
		}()
	}
	
	// Send rows to workers
	go func() {
		for py := 0; py < height; py++ {
			rowChan <- py
		}
		close(rowChan)
	}()
	
	// Monitor progress
	go func() {
		rowsDone := 0
		lastPercent := 0
		for range progress {
			rowsDone++
			percent := rowsDone * 100 / height
			if percent > lastPercent {
				lastPercent = percent
				elapsed := time.Since(start)
				estimatedTotal := time.Duration(float64(elapsed) / float64(percent) * 100)
				remaining := estimatedTotal - elapsed
				fmt.Printf("\rProgress: %d%% | Elapsed: %v | Remaining: %v", 
					percent, elapsed.Round(time.Second), remaining.Round(time.Second))
			}
		}
		fmt.Println()
	}()
	
	// Wait for all workers to finish
	wg.Wait()
	close(progress)
	
	fmt.Printf("\nImage generation complete in %v\n", time.Since(start))
	fmt.Println("Saving image...")
	
	// Save the image in the requested format
	encodeStart := time.Now()
	err := saveImage(img, *outputFlag, *formatFlag)
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("Image save complete in %v\n", time.Since(encodeStart))
	
	// Get final file size
	fileInfo, _ := os.Stat(*outputFlag)
	fmt.Printf("Output file: %s (%.2f MB)\n", *outputFlag, float64(fileInfo.Size())/(1024*1024))
	fmt.Println("Done!")
}

func saveImage(img *image.RGBA, filename, format string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	
	switch format {
	case "bmp":
		return saveBMP(f, img)
	case "tga":
		return saveTGA(f, img)
	case "raw":
		return saveRaw(f, img)
	case "png":
		// PNG with no compression (fastest, largest file)
		encoder := &png.Encoder{
			CompressionLevel: png.NoCompression,
		}
		return encoder.Encode(f, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

// Save as BMP (uncompressed, ~4 bytes per pixel)
func saveBMP(f *os.File, img *image.RGBA) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// BMP header
	// File header (14 bytes)
	fileSize := 14 + 40 + width*height*4
	binary.Write(f, binary.LittleEndian, uint16(0x4D42)) // BM
	binary.Write(f, binary.LittleEndian, uint32(fileSize))
	binary.Write(f, binary.LittleEndian, uint32(0)) // Reserved
	binary.Write(f, binary.LittleEndian, uint32(14+40)) // Data offset
	
	// Info header (40 bytes)
	binary.Write(f, binary.LittleEndian, uint32(40)) // Header size
	binary.Write(f, binary.LittleEndian, uint32(width))
	binary.Write(f, binary.LittleEndian, uint32(height))
	binary.Write(f, binary.LittleEndian, uint16(1)) // Planes
	binary.Write(f, binary.LittleEndian, uint16(32)) // Bits per pixel
	binary.Write(f, binary.LittleEndian, uint32(0)) // Compression (0 = uncompressed)
	binary.Write(f, binary.LittleEndian, uint32(width*height*4)) // Image size
	binary.Write(f, binary.LittleEndian, uint32(2835)) // X pixels per meter
	binary.Write(f, binary.LittleEndian, uint32(2835)) // Y pixels per meter
	binary.Write(f, binary.LittleEndian, uint32(0)) // Colors used
	binary.Write(f, binary.LittleEndian, uint32(0)) // Important colors
	
	// Pixel data (BMP is stored bottom-up)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			c := img.At(x, y).(color.RGBA)
			// BMP stores BGR, not RGB
			f.Write([]byte{c.B, c.G, c.R, c.A})
		}
	}
	
	return nil
}

// Save as TGA (Truevision TGA, uncompressed)
func saveTGA(f *os.File, img *image.RGBA) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// TGA header (18 bytes)
	header := []byte{
		0,          // ID length
		0,          // Color map type (0 = none)
		2,          // Image type (2 = uncompressed true-color)
		0, 0, 0, 0, // Color map specification
		0, 0,       // X origin
		0, 0,       // Y origin
		byte(width & 0xFF), byte((width >> 8) & 0xFF),
		byte(height & 0xFF), byte((height >> 8) & 0xFF),
		32,         // Bits per pixel (32-bit with alpha)
		0b00101000, // Image descriptor (top-left, 8 bits alpha)
	}
	f.Write(header)
	
	// Pixel data (TGA stores BGR, top-left origin)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y).(color.RGBA)
			f.Write([]byte{c.B, c.G, c.R, c.A})
		}
	}
	
	return nil
}

// Save as raw (just the pixel data, no header)
func saveRaw(f *os.File, img *image.RGBA) error {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	
	// Just write the raw RGBA data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.At(x, y).(color.RGBA)
			f.Write([]byte{c.R, c.G, c.B, c.A})
		}
	}
	
	return nil
}

func mandelbrot(z complex128, iterations int) color.Color {
	var v complex128
	
	for n := 0; n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			// Smooth coloring
			value := float64(n) / float64(iterations)
			return getColor(value)
		}
	}
	return color.Black
}

func getColor(t float64) color.Color {
	// Create a gradient from blue to red through cyan, green, yellow
	var r, g, b float64
	
	if t < 0.16 {
		r = 0
		g = 0
		b = t / 0.16
	} else if t < 0.33 {
		r = 0
		g = (t - 0.16) / 0.17
		b = 1
	} else if t < 0.5 {
		r = 0
		g = 1
		b = 1 - (t-0.33)/0.17
	} else if t < 0.66 {
		r = (t - 0.5) / 0.16
		g = 1
		b = 0
	} else if t < 0.83 {
		r = 1
		g = 1 - (t-0.66)/0.17
		b = 0
	} else {
		r = 1
		g = 0
		b = (t - 0.83) / 0.17
	}
	
	r = r * r * 1.2
	g = g * g * 1.2
	b = b * b * 1.2
	
	if r > 1 { r = 1 }
	if g > 1 { g = 1 }
	if b > 1 { b = 1 }
	
	return color.RGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: 255,
	}
}
