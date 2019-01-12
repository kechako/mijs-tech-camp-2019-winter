// +build js,wasm

package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"reflect"
	"syscall/js"
	"unsafe"

	"github.com/nfnt/resize"
)

type Resizer struct {
	resizeImageFunc js.Func
	shutdownFunc    js.Func
	imageBytes      []byte
	resizedBytes    []byte

	done context.CancelFunc
}

func NewResizer() *Resizer {
	r := &Resizer{}
	r.init()

	return r
}

func (r *Resizer) init() {
	r.resizeImageFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) == 0 {
			return nil
		}

		return r.resizeImage(args[0].Int())
	})
	js.Global().Set("resizeImage", r.resizeImageFunc)

	r.shutdownFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if r.done != nil {
			r.done()
		}

		return nil
	})
	js.Global().Call("addEventListener", "beforeunload", r.shutdownFunc)
}

func (r *Resizer) Close() error {
	r.resizeImageFunc.Release()
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

func (r *Resizer) resizeImage(length int) error {
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

	size := img.Bounds().Size()
	resizedImg := resize.Resize(uint(size.X/10), uint(size.Y/10), img, resize.Lanczos3)

	var rbuf bytes.Buffer

	switch name {
	case "png":
		err = png.Encode(&rbuf, resizedImg)
	case "jpeg":
		err = jpeg.Encode(&rbuf, resizedImg, &jpeg.Options{Quality: 80})
	case "gif":
		err = gif.Encode(&rbuf, resizedImg, &gif.Options{NumColors: 256})
	default:
		return errors.New("unsupported format")
	}
	if err != nil {
		return err
	}

	r.resizedBytes = rbuf.Bytes()
	rlen := len(r.resizedBytes)

	roffset := slicePointer(r.resizedBytes)
	fmt.Printf("resized image offset : %d\n", roffset)
	fmt.Printf("resized image length : %d\n", rlen)

	js.Global().Call("setResult", roffset, rlen, name)

	return nil
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
