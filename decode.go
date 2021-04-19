package libjpeg

// #cgo LDFLAGS: -L/usr/local/lib -ljpeg
// #cgo CFLAGS: -I/usr/local/include
// #include <sys/types.h>
// #include <stdlib.h>
// #include <stddef.h>
// #include <stdio.h>
// #include <string.h>
// #include <jpeglib.h>
import "C"

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"unsafe"
)

func DecodeConfig(r io.Reader) (cfg image.Config, er error) {
	var soi, width, height int16
	var components int8

	if er = binary.Read(r, binary.LittleEndian, &soi); er != nil {
		return
	}

	if er = binary.Read(r, binary.LittleEndian, &width); er != nil {
		return
	}

	if er = binary.Read(r, binary.LittleEndian, &height); er != nil {
		return
	}

	if er = binary.Read(r, binary.LittleEndian, &components); er != nil {
		return
	}

	if components == 1 || components == 4 {
		cfg.ColorModel = color.RGBAModel

	} else if components == 3 {
		cfg.ColorModel = color.GrayModel

	} else {
		er = fmt.Errorf("Invalid number of components (%d)", components)
		return
	}

	cfg.Width = int(width)
	cfg.Height = int(height)
	return
}

func decodeGrayscale(cinfo *C.struct_jpeg_decompress_struct) image.Image {
	cbuflen := cinfo.output_width
	cbuf := C.malloc(C.size_t(cbuflen))

	scanLine := cbuf
	scanLines := C.JSAMPARRAY(unsafe.Pointer(&scanLine))

	img := image.NewGray(image.Rect(0, 0, int(cinfo.output_width), int(cinfo.output_height)))

	for cinfo.output_scanline < cinfo.output_height {
		y := int(cinfo.output_scanline)
		off := y * img.Stride

		C.jpeg_read_scanlines(cinfo, scanLines, 1)

		for x := 0; x < int(cinfo.output_width); x += 1 {
			base := uintptr(cbuf) + uintptr(x*1)
			img.Pix[off+x] = *(*uint8)(unsafe.Pointer(base + 0))
		}
	}

	C.free(unsafe.Pointer(cbuf))

	return img
}

func decodeRGB(cinfo *C.struct_jpeg_decompress_struct) image.Image {
	cbuflen := cinfo.output_width * 3
	cbuf := C.malloc(C.size_t(cbuflen))

	scanLine := cbuf
	scanLines := C.JSAMPARRAY(unsafe.Pointer(&scanLine))

	img := image.NewRGBA(image.Rect(0, 0, int(cinfo.output_width), int(cinfo.output_height)))

	for cinfo.output_scanline < cinfo.output_height {
		y := int(cinfo.output_scanline)
		off := y * img.Stride

		C.jpeg_read_scanlines(cinfo, scanLines, 1)

		for x := 0; x < int(cinfo.output_width); x += 1 {
			base := uintptr(cbuf) + uintptr(x*3)
			img.Pix[off+4*x+0] = *(*uint8)(unsafe.Pointer(base + 0))
			img.Pix[off+4*x+1] = *(*uint8)(unsafe.Pointer(base + 1))
			img.Pix[off+4*x+2] = *(*uint8)(unsafe.Pointer(base + 2))
			img.Pix[off+4*x+3] = 255
		}
	}

	C.free(unsafe.Pointer(cbuf))

	return img
}

func decodeCMYK(cinfo *C.struct_jpeg_decompress_struct) image.Image {
	cbuflen := cinfo.output_width * 4
	cbuf := C.malloc(C.size_t(cbuflen))

	scanLine := cbuf
	scanLines := C.JSAMPARRAY(unsafe.Pointer(&scanLine))

	img := image.NewRGBA(image.Rect(0, 0, int(cinfo.output_width), int(cinfo.output_height)))

	for cinfo.output_scanline < cinfo.output_height {
		y := int(cinfo.output_scanline)
		off := y * img.Stride

		C.jpeg_read_scanlines(cinfo, scanLines, 1)

		for x := 0; x < int(cinfo.output_width); x += 1 {
			base := uintptr(cbuf) + uintptr(x*4)

			c := *(*uint8)(unsafe.Pointer(base + 0))
			m := *(*uint8)(unsafe.Pointer(base + 1))
			y := *(*uint8)(unsafe.Pointer(base + 2))
			k := *(*uint8)(unsafe.Pointer(base + 3))

			r := uint8(255. * (1. - float64(c)) * (1. - float64(k)))
			g := uint8(255. * (1. - float64(m)) * (1. - float64(k)))
			b := uint8(255. * (1. - float64(y)) * (1. - float64(k)))

			img.Pix[off+4*x+0] = r
			img.Pix[off+4*x+1] = g
			img.Pix[off+4*x+2] = b
			img.Pix[off+4*x+3] = 255
		}
	}

	return img
}

func Decode(r io.Reader) (img image.Image, er error) {
	/* Reading the whole file in may be inefficient, but libjpeg wants callbacks
	 * to functions to read in more data, and that is a nightmare to implement. We
	 * don't want to read the entire stream, however, which means pulling the header.
	 * We may be able to read enough to call jpeg_read_header with a [10]byte, but
	 * I'll change it later if need be, since this probably doesn't play nicely
	 * with a non-closing io.Reader */

	var wholeFile []byte
	if wholeFile, er = ioutil.ReadAll(r); er != nil {
		return
	}

	fileBytes := C.CBytes(wholeFile)

	var cinfo *C.struct_jpeg_decompress_struct
	var jerr *C.struct_jpeg_error_mgr

	cinfolen := C.size_t(unsafe.Sizeof(C.struct_jpeg_decompress_struct{}))
	jerrlen := C.size_t(unsafe.Sizeof(C.struct_jpeg_error_mgr{}))

	cinfo = (*C.struct_jpeg_decompress_struct)(C.malloc(cinfolen))
	jerr = (*C.struct_jpeg_error_mgr)(C.malloc(jerrlen))
	C.memset(unsafe.Pointer(cinfo), 0, cinfolen)
	C.memset(unsafe.Pointer(jerr), 0, jerrlen)

	cinfo.err = C.jpeg_std_error(jerr)

	C.jpeg_CreateDecompress(cinfo, C.JPEG_LIB_VERSION, cinfolen)
	C.jpeg_mem_src(cinfo, (*C.uchar)(fileBytes), C.ulong(len(wholeFile)))

	if C.jpeg_read_header(cinfo, C.TRUE) == C.JPEG_HEADER_OK {
		C.jpeg_start_decompress(cinfo)

		if cinfo.num_components == 1 {
			img = decodeGrayscale(cinfo)

		} else if cinfo.num_components == 3 {
			img = decodeRGB(cinfo)

		} else if cinfo.num_components == 4 {
			img = decodeCMYK(cinfo)

		} else {
			er = fmt.Errorf("Invalid number of components (%d)", cinfo.num_components)
		}

		if er == nil {
			C.jpeg_finish_decompress(cinfo)
		}
	}

	C.jpeg_destroy_decompress(cinfo)
	C.free(unsafe.Pointer(jerr))
	C.free(unsafe.Pointer(cinfo))
	C.free(fileBytes)
	return
}

func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", Decode, DecodeConfig)
}
