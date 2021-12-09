package bitmap

import (
	"encoding/binary"
	"github.com/kisunSea/gopkg/logging"
	"io"
	"os"
	"sync"
)

var logger, _ = logging.GLogger()

type Bit struct {
	Offset int64
	Length int
	Hash   uint64
}

type BitmapFile struct {
	sync.RWMutex

	path     string
	fp       *os.File
	isOpened bool
	bits     chan *Bit
}

func (b *BitmapFile) GetBlockSize() int {
	return 8
}

func (b *BitmapFile) WriteByOffset(offset int64, hash uint64) (err error) {
	b.Lock()
	defer b.Unlock()

	if err = b.initHandle(); err != nil {
		return
	}

	tmpBuf := make([]byte, b.GetBlockSize())
	binary.BigEndian.PutUint64(tmpBuf, hash)

	if _, err = b.fp.WriteAt(tmpBuf, offset); err != nil {
		return
	}

	return nil
}

func (b *BitmapFile) WriteByIndex(index int64) (err error) {
	_ = index
	// TODO
	return nil
}

func (b *BitmapFile) ReadByOffset(offset int64) (hash uint64, err error) {
	b.Lock()
	defer b.Unlock()

	if err := b.initHandle(); err != nil {
		return 0, err
	}

	var buf = make([]byte, 8)
	if _, err = b.fp.ReadAt(buf, offset); err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(buf), nil
}

func (b *BitmapFile) Iter() {
	if err := b.initHandle(); err != nil {
		return
	}

	go b.__iter(0)
}

func (b *BitmapFile) __iter(__start int64) {
	b.Lock()

	defer b.Unlock()
	defer close(b.bits)

	var count int64

	buf := make([]byte, b.GetBlockSize())
	for true {
		n, err := b.fp.Read(buf)
		if err == io.EOF {
			logger.DebugF("Bitmap read finished")
			return
		} else if err != nil {
			logger.ErrorF("Bitmap read err: %v", err)
			return
		} else {
			var _bit = new(Bit)
			_bit.Offset = __start + count*int64(b.GetBlockSize())
			_bit.Length = b.GetBlockSize()
			_bit.Hash = binary.BigEndian.Uint64(buf[:n])
			b.bits <- _bit
			count++
		}
	}
}

func (b *BitmapFile) initHandle() (err error) {
	if !b.isOpened {

		b.fp, err = os.OpenFile(b.path, os.O_RDWR|os.O_CREATE, 0)
		if err != nil {
			logger.ErrorF("BitmapFile.initHandle OpenFile Failed: %v", err)
		}
		logger.DebugF("---------init handle successfully: %v--------", b.path)
		b.isOpened = true
	}

	return
}

func (b *BitmapFile) Close() (err error) {
	if err = b.fp.Close(); err != nil {
		logger.ErrorF("close `%v` Failed: %v", b.path, err)
	}

	b.isOpened = false
	logger.DebugF("---------close handle successfully: %v--------", b.path)
	return nil
}
