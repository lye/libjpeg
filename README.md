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

* They're read-only right now, only because I had no need to write jpegs.
* They're using the default libjpeg `jpeg_error_mgr`, which on a catastrophic failure calls `exit(2)`, thus killing your application. A "catastrophic failure", I believe, is only when the libjpeg API is used incorrectly, which should never happen. There's no clean way around it -- the docs say "just set `setjmp`", but that isn't really a tenable option in Go. Swapping it out for a no-op function will probably cause memory corruption.
* It doesn't support progressive loading in progressive mode (e.g., it just loads the final version of the file).
* The test suite and documentation are non-existent.
