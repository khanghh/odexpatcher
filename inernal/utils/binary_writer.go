package utils

import "io"

type BinaryWriter struct {
	buf []byte
	pos int64
}

func NewBinaryWriter(buf []byte) *BinaryWriter {
	return &BinaryWriter{buf: buf}
}

func (w *BinaryWriter) Len() int64 {
	return int64(len(w.buf))
}

func (w *BinaryWriter) Read(p []byte) (n int, err error) {
	toRead := w.pos + int64(len(p))
	switch {
	case w.pos == w.Len():
		p = []byte{}
	case toRead <= w.Len():
		p = w.buf[w.pos : int(w.pos)+len(p)]
	case toRead > w.Len():
		p = w.buf[w.pos:]
	}
	n = len(p)
	w.pos += int64(len(p))
	return
}

func (w *BinaryWriter) Write(p []byte) (n int, err error) {
	writeLen := int64(len(p))
	switch {
	case w.Len() == 0:
		w.buf = p
		w.pos = int64(len(p))
	case w.pos == w.Len():
		w.buf = append(w.buf, p...)
		w.pos += writeLen
	case w.pos < w.Len():
		switch {
		case w.pos+writeLen > w.Len():
			w.buf = append(w.buf[:w.pos], p...)
		case w.pos+writeLen <= w.Len():
			w.buf = append(w.buf[:w.pos], append(p, w.buf[w.pos+writeLen:]...)...)
		}
		w.pos += writeLen
	}
	return len(p), err
}

// Seek sets the offset for the next Read or Write to offset, interpreted according to whence
func (w *BinaryWriter) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		w.pos = 0 + offset
	case io.SeekCurrent:
		w.pos = w.pos + offset
	case io.SeekEnd:
		w.pos = w.Len() + offset
	}
	return w.pos, nil
}
