package main

import ("fmt"
				"github.com/fogleman/gg"
				"os"
				"errors"
				"image/color"
				"image"
				"image/png"
				"path/filepath"
)

var background_image_filename string = "bkg_img"
var output_file_name string = "image_to_post"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	dc := gg.NewContext(1200, 628)
	background_image, err := gg.LoadImage(background_image_filename)
	if err != nil {
		return errors.Unwrap(err)
	}
	dc.DrawImage(background_image, 0, 0)
	//Aggiungi overlay semi trasparente
	margin := 20.0
	x := margin
	y := margin
	w := float64(dc.Width()) - (2.0 * margin)
	h := float64(dc.Height()) - (2.0 * margin)
	dc.SetColor(color.RGBA{0, 0, 0, 204})
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()

	//Aggiungi nome
	fontPath := filepath.Join("./fonts", "OpenSans-Bold.ttf")
	if err := dc.LoadFontFace(fontPath, 40); err != nil {
		return errors.Unwrap(err)
	}
	dc.SetColor(color.White)
	s := "@Owen10825549"
	marginX := 30.0
	marginY := -5.0
	textWidth, textHeight := dc.MeasureString(s)
	x = float64(dc.Width()) - textWidth - marginX
	y = float64(dc.Height()) - textHeight - marginY
	dc.DrawString(s, x, y)

	//Aggiungi eventi


	//Salva immagine alla fine
	if err := dc.SavePNG(output_file_name); err != nil {
	return errors.Unwrap(err)
	}
	return nil
}

func SavePNG(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
