package consts

const (
	// JPEGMagicNumber is the magic number that identifies a JPEG file. This can be used to identify the file type
	// returned from HTTPBin backends when using `/image/*` endpoints.
	JPEGMagicNumber = "\xff\xd8\xff\xe0\x00\x10JFIF"

	// PNGMagicNumber is the magic number that identifies a PNG file. This can be used to identify the file type
	// returned from HTTPBin backends when using `/image/*` endpoints.
	PNGMagicNumber = "\x89PNG\r\n\x1a\n"
)
