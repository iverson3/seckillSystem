package common

import (
	"crypto/md5"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
)

const defaultBufSize = 64 * 1024 * 1024   // 默认buffer大小 64MB

type LocalFile struct {
	LocalFileMeta

	BufSize int
	Buf []byte
	File *os.File
}


// OpenPath 检查文件状态并获取文件的大小 (Length)
func (lf *LocalFile) OpenPath() (err error) {
	if lf.File != nil {
		lf.File.Close()
	}

	lf.File, err = os.Open(lf.Path)
	if err != nil {
		return
	}

	info, err := lf.File.Stat()
	if err != nil {
		return
	}

	lf.Length  = info.Size()
	lf.ModTime = info.ModTime().Unix()
	return
}

// Close 关闭文件
func (lf *LocalFile) Close() error {
	if lf.File == nil {
		return fmt.Errorf("file is nil")
	}
	return lf.File.Close()
}

func (lf *LocalFile) initBuf() {
	if lf.Buf == nil {
		if lf.BufSize != 0 {
			lf.Buf = make([]byte, lf.BufSize)
		} else {
			lf.Buf = make([]byte, defaultBufSize)
		}
	}
}

func (lf *LocalFile) repeatRead(ws ...io.Writer)  {
	if lf.File == nil {
		return
	}

	lf.initBuf()

	var (
		begin int64
		n int
		err error
	)

	handle := func() {
		begin += int64(n)
		for k := range ws {
			_, err = ws[k].Write(lf.Buf[:n])
			if err != nil {
				fmt.Println("write data failed! error: ", err)
			}
		}
	}

	// 读文件
	for {
		n, err = lf.File.ReadAt(lf.Buf, begin)
		if err != nil {
			if err == io.EOF {
				handle()
			} else {
				fmt.Println("read file failed! error: ", err)
			}
			break
		}
		handle()
	}
}

// Sum 计算文件摘要值
func (lf *LocalFile) Sum(cfg SumConfig) {
	var (
		md5w   hash.Hash
		crc32w hash.Hash32
	)

	ws := make([]io.Writer, 0, 2)
	if cfg.IsMD5Sum {
		md5w = md5.New()
		ws = append(ws, md5w)
	}
	if cfg.IsCRC32Sum {
		crc32w = crc32.NewIEEE()
		ws = append(ws, crc32w)
	}
	if cfg.IsSliceMD5Sum {
		lf.SliceMD5Sum()
	}

	lf.repeatRead(ws...)

	if cfg.IsMD5Sum {
		lf.MD5 = md5w.Sum(nil)
	}
	if cfg.IsCRC32Sum {
		lf.CRC32 = crc32w.Sum32()
	}
}

// Md5Sum 获取文件的 md5 值
func (lf *LocalFile) Md5Sum() {
	lf.Sum(SumConfig{
		IsMD5Sum: true,
	})
}

// SliceMD5Sum 获取文件前 requiredSliceLen (256KB) 切片的 md5 值
func (lf *LocalFile) SliceMD5Sum() {
	if lf.File == nil {
		return
	}

	// 获取前 256KB 文件切片的 md5
	lf.initBuf()

	m := md5.New()
	n, err := lf.File.ReadAt(lf.Buf, 0)
	if err != nil {
		if err == io.EOF {
			goto md5sum
		} else {
			fmt.Printf("SliceMD5Sum: %s\n", err)
			return
		}
	}

md5sum:
	m.Write(lf.Buf[:n])
	lf.SliceMD5 = m.Sum(nil)
}

// Crc32Sum 获取文件的 crc32 值
func (lf *LocalFile) Crc32Sum() {
	lf.Sum(SumConfig{
		IsCRC32Sum: true,
	})
}