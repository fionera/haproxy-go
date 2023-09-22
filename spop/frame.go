package spop

import (
	"encoding/binary"
	"fmt"
	"github.com/fionera/haproxy-go/pkg/buffer"
	"io"
	"sync"

	"github.com/fionera/haproxy-go/pkg/encoding"
)

const uint32Len = 4

var framePool = sync.Pool{
	New: func() any {
		return &frame{
			length: make([]byte, uint32Len),
			buf:    buffer.NewSliceBuffer(maxFrameSize),
		}
	},
}

func acquireFrame() *frame {
	return framePool.Get().(*frame)
}

func releaseFrame(f *frame) {
	f.buf.Reset()

	framePool.Put(f)
}

type frameMetadata struct {
	Flags    frameFlag
	StreamID int64
	FrameID  int64
}

type frame struct {
	length []byte
	buf    *buffer.SliceBuffer

	frameType frameType
	meta      frameMetadata
}

func (f *frame) ReadFrom(r io.Reader) (int64, error) {
	if _, err := r.Read(f.length); err != nil {
		return 0, fmt.Errorf("reading frame length: %w", err)
	}
	frameLen := binary.BigEndian.Uint32(f.length)

	f.buf.Reset()
	dataBuf := f.buf.WriteNBytes(int(frameLen))

	// read full frame into buffer
	n, err := r.Read(dataBuf)
	if err != nil {
		return int64(n + len(f.length)), fmt.Errorf("reading frame payload: %w", err)
	}

	if n != int(frameLen) {
		return int64(n + len(f.length)), io.ErrUnexpectedEOF
	}

	return int64(n + len(f.length)), f.decodeHeader()
}

func (f *frame) WriteTo(w io.Writer) (int64, error) {
	binary.BigEndian.PutUint32(f.length, uint32(f.buf.Len()))

	if n, err := w.Write(f.length); err != nil {
		return int64(n), err
	}

	n, err := w.Write(f.buf.ReadBytes())
	return int64(n + len(f.length)), err
}

func (f *frame) encodeHeader() error {
	f.buf.WriteNBytes(1)[0] = byte(f.frameType)

	binary.BigEndian.PutUint32(f.buf.WriteNBytes(uint32Len), uint32(f.meta.Flags))

	n, err := encoding.PutVarint(f.buf.WriteBytes(), f.meta.StreamID)
	if err != nil {
		return err
	}
	f.buf.AdvanceW(n)

	n, err = encoding.PutVarint(f.buf.WriteBytes(), f.meta.FrameID)
	if err != nil {
		return err
	}
	f.buf.AdvanceW(n)

	return nil
}

func (f *frame) decodeHeader() error {
	// We don't need to validate here,
	// there is validation further down the chain
	f.frameType = frameType(f.buf.ReadNBytes(1)[0])

	f.meta.Flags = frameFlag(binary.BigEndian.Uint32(f.buf.ReadNBytes(uint32Len)))

	streamID, n, err := encoding.Varint(f.buf.ReadBytes())
	if err != nil {
		return err
	}

	f.meta.StreamID = streamID
	f.buf.AdvanceR(n)

	frameID, n, err := encoding.Varint(f.buf.ReadBytes())
	if err != nil {
		return err
	}
	f.meta.FrameID = frameID
	f.buf.AdvanceR(n)

	return nil
}
