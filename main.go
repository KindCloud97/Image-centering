package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func generatePath(file string, canvasMult *int) string {
	//newPath := strings.Split(path.Base(oldPath), ".")
	//return newPath[0] + strconv.Itoa(canvasMult) + "." + newPath[1]

	inputDir := filepath.Dir(file)
	inputBase := filepath.Base(file)
	inputExt := filepath.Ext(file)

	newName := fmt.Sprintf("%s-canvas-resized-%dx%s",
		strings.TrimSuffix(inputBase, inputExt), canvasMult, inputExt)
	return filepath.Join(inputDir, newName)

}

func run(canvasMult int, fpath string) error {
	file, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("Couldn't open file! %w", err)
	}
	defer file.Close()

	src, err := png.Decode(file)
	if err != nil {
		return err
	}

	rect := src.Bounds()
	rect.Min, rect.Max = rect.Min.Mul(canvasMult), rect.Max.Mul(canvasMult)
	newImg := image.NewRGBA(rect)

	centerRect := src.Bounds().Add(rect.Size().Div(2).Sub(src.Bounds().Size().Div(2)))
	draw.Draw(newImg, centerRect, src, image.Point{}, draw.Src)
	result, err := os.Create(generatePath(fpath, &canvasMult))
	if err != nil {
		return err
	}

	if err := png.Encode(result, newImg); err != nil {
		return err
	}

	return nil
}

func main() {
	canvasMult := flag.Int("cmult", 2, "Multiplies image's canvas in size.")
	fpath := flag.String("fpath", "", "Path to a file.")
	flag.Parse()
	if *fpath == "" {
		fmt.Println("Empty path to file!")
		return
	}
	if err := run(*canvasMult, *fpath); err != nil {
		fmt.Println(err)
	}

}
