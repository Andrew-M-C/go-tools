package compress

// Algorithm indicates suported Algorithms in this package
type Algorithm string

const (
	// Gzip is gzip compression
	Gzip = Algorithm("gzip")
)

// Error is a inner error type
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNotSupported indicates that given algorithm is not supported.
	ErrNotSupported = Error("not supported")
	// ErrNilInput indicates that given data is nil.
	ErrNilInput = Error("nil input")
)

// Options indicates compress/uncompress options
type Options struct {
	Name    string
	Comment string
	MaxRead int
}

// Compress compresses data with specified algorithm
func Compress(algorithm Algorithm, in []byte, opts ...Options) (out []byte, err error) {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	switch algorithm {
	default:
		return nil, ErrNotSupported
	case Gzip:
		return gzipCompress(in, opt)
	}
}

// Uncompress uncompress data with specified algorithm
func Uncompress(algorithm Algorithm, in []byte, opts ...Options) (out []byte, err error) {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	switch algorithm {
	default:
		return nil, ErrNotSupported
	case Gzip:
		return gzipUncompress(in, opt)
	}
}
