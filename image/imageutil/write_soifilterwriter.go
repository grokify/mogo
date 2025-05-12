package imageutil

import (
	"io"
)

// SOIFilterWriter is a writer that filters out the SOI marker (0xFF 0xD8) from JPEG data.
type SOIFilterWriter struct {
	w     io.Writer
	state int // 0: initial, 1: saw 0xFF, 2: saw 0xD8
}

func NewSOIFilterWriter(w io.Writer) *SOIFilterWriter {
	return &SOIFilterWriter{w: w}
}

func (s *SOIFilterWriter) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	// If we're in state 2, we've already seen the SOI marker, write everything
	if s.state == 2 {
		return s.w.Write(p)
	}

	// Process the data byte by byte
	for i := 0; i < len(p); i++ {
		switch s.state {
		case 0: // initial state
			if p[i] == JPEGMarkerPrefix {
				s.state = 1
			} else {
				if _, err := s.w.Write(p[i : i+1]); err != nil {
					return n, err
				}
				n++
			}
		case 1: // saw 0xFF
			if p[i] == JPEGMarkerSOI {
				s.state = 2
			} else {
				// Write the 0xFF we saw earlier
				if _, err := s.w.Write([]byte{JPEGMarkerPrefix}); err != nil {
					return n, err
				}
				n++
				// Write current byte if it's not 0xFF
				if p[i] != JPEGMarkerPrefix {
					if _, err := s.w.Write(p[i : i+1]); err != nil {
						return n, err
					}
					n++
				}
				s.state = 0
			}
		}
	}

	// If we're still in state 1 at the end of the buffer, write the 0xFF
	if s.state == 1 {
		if _, err := s.w.Write([]byte{JPEGMarkerPrefix}); err != nil {
			return n, err
		}
		n++
		s.state = 0
	}

	return n, nil
}

// Close implements io.Closer
func (s *SOIFilterWriter) Close() error {
	// If we're in state 1, write the final 0xFF
	if s.state == 1 {
		if _, err := s.w.Write([]byte{JPEGMarkerPrefix}); err != nil {
			return err
		}
	}
	return nil
}
