## Go bindings for libjpeg

These bindings expose themselves through the normal image interface -- simply import them and they'll automatically handle calls to `image.Decode` when they appear to be jpeg-encoded.

    import (
        "image"
        _ "github.com/lye/libjpeg"
    )

The advantage of using `libjpeg` over the normal `image/jpeg` package is simply that `libjpeg` supports progressive jpegs; `image/jpeg` currently does not. If you don't need support for progressive jpegs (e.g., you're not decoding arbitrary data) then `image/jpeg` is probably a safer choice.

### Dependencies

* `libjpeg`; if it complains about not being able to find the files, add in the proper search paths and send me a pull request.

### Known issues

* They're using the default libjpeg `jpeg_error_mgr`, which on a catastrophic failure calls `exit(3)`, thus killing your application. A "catastrophic failure", I believe, is only when the libjpeg API is used incorrectly, which should never happen. There's no clean way around it -- the docs say "just set `setjmp`", but that isn't really a tenable option in Go. Swapping it out for a no-op function will probably cause memory corruption.
* It doesn't support progressive loading in progressive mode (e.g., it just loads the final version of the file).
* The test suite and documentation are non-existent.
* It conflicts with `image/jpeg` for the entire binary. Not just your application source, but every package it depends on -- if something somewhere requires `image/jpeg` it's a crapshoot which one will be handling your files. This is an unfortunate byproduct of the `image` decoder registration process. The only real workaround, if for some reason you need access to both, is to call `image/jpeg.Decode` and `libjpeg.Decode` directly.
* `Decode` does not attempt to compute the size of the image before handing it to libjpeg: it reads the input stream until `EOF`. This may cause problems if you're, e.g., attempting to read an image off a network connection that isn't closed by the remote host (or if you have excessively large jpeg files). The reason for this deficiency is _callbacks_.
