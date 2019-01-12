// +build js,wasm

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/nfnt/resize"
)

type Resizer struct {
	convertImageFunc js.Func
	shutdownFunc     js.Func
	imageBytes       []byte
	convertedBytes   []byte

	done context.CancelFunc
}

func NewResizer() *Resizer {
	r := &Resizer{}
	r.init()

	return r
}

func (r *Resizer) init() {
	r.convertImageFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) < 2 {
			return nil
		}

		return r.convertImage(args[0].Int(), ConvertMode(args[1].Int()))
	})
	js.Global().Set("convertImage", r.convertImageFunc)

	r.shutdownFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if r.done != nil {
			r.done()
		}

		return nil
	})
	js.Global().Call("addEventListener", "beforeunload", r.shutdownFunc)
}

func (r *Resizer) Close() error {
	r.convertImageFunc.Release()
	r.shutdownFunc.Release()

	fmt.Println("Closed")

	return nil
}

func (r *Resizer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	r.done = cancel

	<-ctx.Done()

	return nil
}

type ConvertMode int

const (
	ConvertModeReize ConvertMode = iota
	ConvertModeGrayscale
)

func (r *Resizer) convertImage(length int, mode ConvertMode) error {
	r.imageBytes = make([]byte, length)
	offset := slicePointer(r.imageBytes)
	js.Global().Call("setFileBytesToMem", offset)

	fmt.Printf("source image offset : %d\n", offset)
	fmt.Printf("source image length : %d\n", length)

	img, name, err := image.Decode(bytes.NewReader(r.imageBytes))
	if err != nil {
		return err
	}
	fmt.Printf("source image type : %s\n", name)

	var convImage image.Image

	switch mode {
	case ConvertModeReize:
		convImage = r.resizeImage(img)
	case ConvertModeGrayscale:
		convImage = r.grayscaleImage(img)
	}

	var convBuf bytes.Buffer
	switch name {
	case "png":
		err = png.Encode(&convBuf, convImage)
	case "jpeg":
		err = jpeg.Encode(&convBuf, convImage, &jpeg.Options{Quality: 80})
	case "gif":
		err = gif.Encode(&convBuf, convImage, &gif.Options{NumColors: 256})
	default:
		return errors.New("unsupported format")
	}
	if err != nil {
		return err
	}

	r.convertedBytes = convBuf.Bytes()
	convLen := len(r.convertedBytes)

	convOffset := slicePointer(r.convertedBytes)
	fmt.Printf("converted image offset : %d\n", convOffset)
	fmt.Printf("converted image length : %d\n", convLen)

	js.Global().Call("setResult", convOffset, convLen, name)

	return nil
}

func (r *Resizer) resizeImage(img image.Image) image.Image {
	size := img.Bounds().Size()
	return resize.Resize(uint(size.X/10), uint(size.Y/10), img, resize.Lanczos3)
}

func (r *Resizer) grayscaleImage(img image.Image) image.Image {
	bounds := img.Bounds()
	dest := image.NewGray16(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.Gray16Model.Convert(img.At(x, y))
			gray, _ := c.(color.Gray16)
			dest.Set(x, y, gray)
		}
	}
	return dest
}

func slicePointer(b []byte) uintptr {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	return uintptr(unsafe.Pointer(header.Data))
}

func main() {
	r := NewResizer()
	defer r.Close()

	r.Run(context.Background())
}
