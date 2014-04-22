package libjpeg

// #include <stddef.h>
// #include <stdio.h>
// #include <stdlib.h>
// #include <jpeglib.h>
// typedef unsigned char *PUCHAR;
import "C"
import (
	"io"
	"unsafe"
	"image"
)

type Options struct {
	Quality int
}

func Encode(w io.Writer, m image.Image, o *Options) error {
	quality := 75
	if o != nil {
		quality = o.Quality
	}

	var cinfo C.struct_jpeg_compress_struct
	var jerr C.struct_jpeg_error_mgr
	var workBuf *C.uchar
	var workBufLen C.ulong

	cinfo.err = C.jpeg_std_error(&jerr)
	C.jpeg_CreateCompress(&cinfo, C.JPEG_LIB_VERSION, C.size_t(unsafe.Sizeof(cinfo)))
	C.jpeg_mem_dest(&cinfo, &workBuf, &workBufLen)

	bounds := m.Bounds()
	cinfo.image_width = C.JDIMENSION(bounds.Dx())
	cinfo.image_height = C.JDIMENSION(bounds.Dy())
	cinfo.input_components = 3
	cinfo.in_color_space = C.JCS_RGB

	C.jpeg_set_defaults(&cinfo);
	C.jpeg_set_quality(&cinfo, C.int(quality), C.boolean(1))
	C.jpeg_start_compress(&cinfo, C.boolean(1))

	rowBuf := make([]byte, cinfo.image_width * 3)

	for cinfo.next_scanline < cinfo.image_height {
		for x := 0; x < int(cinfo.image_width); x += 1 {
			r, g, b, _ := m.At(x, int(cinfo.next_scanline)).RGBA()
			rowBuf[x*3] = byte(r >> 8)
			rowBuf[x*3+1] = byte(g >> 8)
			rowBuf[x*3+2] = byte(b >> 8)
		}

		rowPointer := C.JSAMPROW(unsafe.Pointer(&rowBuf[0]))
		C.jpeg_write_scanlines(&cinfo, &rowPointer, 1);
	}

	C.jpeg_finish_compress(&cinfo);
	C.jpeg_destroy_compress(&cinfo);

	outBs := C.GoBytes(unsafe.Pointer(workBuf), C.int(workBufLen))
	w.Write(outBs)
	C.free(unsafe.Pointer(workBuf))

	return nil
}
