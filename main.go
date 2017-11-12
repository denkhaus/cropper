package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/muesli/smartcrop"
	"github.com/nfnt/resize"
	"github.com/urfave/cli"
)

var (
	ErrParseWidthHeight  = errors.New("malformed width/height info")
	ErrNoSubImageSupport = errors.New("no subimage support")
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "cropper"
	app.EnableBashCompletion = true
	app.Usage = "A simple image croper"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "wh",
			Value: "580x434",
			Usage: "width and hight of new image in format <with>x<height>",
		},
	}

	app.Action = func(ctx *cli.Context) {
		widthHeight := ctx.GlobalString("wh")

		a := ctx.Args()
		if !a.Present() {
			cli.ShowAppHelp(ctx)
			return
		}

		files, err := getFiles(a)
		if err != nil {
			log.Fatalf("error getting input file(s): %v", err)
		}

		width, height, err := parseWidthHeight(widthHeight)
		if err != nil {
			log.Fatal(err)
		}

		if err := processFiles(files, widthHeight, width, height); err != nil {
			log.Fatal(err)
		}

	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func processFiles(files []string, prefix string, width, height int) error {
	analyzer := smartcrop.NewAnalyzer()

	for _, fPath := range files {
		ext := strings.ToLower(filepath.Ext(fPath))

		if ext != ".jpg" && ext != ".png" {
			continue
		}

		fi, err := os.Open(fPath)
		if err != nil {
			return errors.Annotate(err, "Open")
		}
		defer fi.Close()

		img, _, err := image.Decode(fi)
		if err != nil {
			return errors.Annotate(err, "Decode")
		}

		topCrop, err := analyzer.FindBestCrop(img, width, height)
		if err != nil {
			return errors.Annotate(err, "SmartCrop")
		}

		sub, ok := img.(SubImager)
		if !ok {
			return ErrNoSubImageSupport
		}

		cropImage := sub.SubImage(topCrop)
		newImage := resize.Resize(uint(width), uint(height), cropImage, resize.Lanczos3)

		fileName := filepath.Base(fPath)
		fileName = fileName[0 : len(fileName)-len(ext)]
		newFilePath := filepath.Join(filepath.Dir(fPath), fmt.Sprintf("%s_%s%s", fileName, prefix, ext))

		if ext == ".jpg" {
			if err := outputJpeg(&newImage, newFilePath); err != nil {
				return errors.Annotate(err, "outputJpeg")
			}
		}
		if ext == ".png" {
			if err := outputPng(&newImage, newFilePath); err != nil {
				return errors.Annotate(err, "outputPng")
			}
		}
	}

	return nil
}

func parseWidthHeight(in string) (int, int, error) {
	parts := strings.Split(in, "x")
	if len(parts) != 2 {
		return 0, 0, ErrParseWidthHeight
	}

	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, ErrParseWidthHeight
	}

	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, ErrParseWidthHeight
	}

	return width, height, nil
}

func getFiles(args cli.Args) ([]string, error) {
	files := make([]string, 0)
	for _, a := range args {
		fi, err := os.Stat(a)
		if err != nil {
			return nil, errors.Annotate(err, "Stat")
		}

		dir, err := filepath.Abs(a)
		if err != nil {
			return nil, errors.Annotate(err, "Abs")
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			fs, err := ioutil.ReadDir(dir)
			if err != nil {
				return nil, errors.Annotate(err, "ReadDir")
			}

			for _, f := range fs {
				files = append(files, filepath.Join(dir, f.Name()))
			}

		case mode.IsRegular():
			files = append(files, dir)
		}
	}

	return files, nil
}

func outputJpeg(img *image.Image, name string) error {
	fso, err := os.Create(name)
	if err != nil {
		return errors.Annotate(err, "Create")
	}
	defer fso.Close()

	jpeg.Encode(fso, (*img), &jpeg.Options{Quality: 100})
	return nil
}

func outputPng(img *image.Image, name string) error {
	fso, err := os.Create(name)
	if err != nil {
		return errors.Annotate(err, "Create")
	}
	defer fso.Close()

	png.Encode(fso, (*img))
	return nil
}
