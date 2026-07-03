// Package pkg implements a stand-in for estargz's seekable-tar-gz reader.
// It is deliberately tiny: this PoC only needs estargz to be an independent
// Go module with a real dependency on shared-lib.
package pkg

import "github.com/abhi4u1947/go-multimodule-poc/entities/shared-lib/logger"

// Reader is a minimal placeholder for an estargz archive reader.
type Reader struct {
	log *logger.Logger
}

// NewReader returns a Reader that logs via shared-lib's logger.
func NewReader() *Reader {
	return &Reader{log: logger.New("estargz")}
}

// Describe logs and returns a human-readable description of the reader.
func (r *Reader) Describe() string {
	r.log.Info("estargz reader initialized")
	return "estargz.Reader (stand-in implementation)"
}

// TOCDigest returns a placeholder digest for the archive's table of
// contents.
//
// Added in estargz v0.18.2.
func (r *Reader) TOCDigest() string {
	r.log.Info("computing TOC digest")
	return "sha256:stand-in-toc-digest"
}
