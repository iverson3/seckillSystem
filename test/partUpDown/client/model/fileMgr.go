package model

import (
	"fmt"
	"os"
)

const defaultBufSize = 64 * 1024 * 1024   // 默认buffer大小 64MB

type LocalFileMeta struct {
	Path     string `json:"path"`     // 本地路径
	Length   int64  `json:"length"`   // 文件大小
	SliceMD5 []byte `json:"slicemd5"` // 文件前 requiredSliceLen (256KB) 切片的 md5 值
	MD5      []byte `json:"md5"`      // 文件的 md5
	CRC32    uint32 `json:"crc32"`    // 文件的 crc32
	ModTime  int64  `json:"modtime"`  // 修改日期
}
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

// https://chunlife.top/2019/04/09/%E6%9C%8D%E5%8A%A1%E5%99%A8%E4%B8%8A%E4%BC%A0%E4%B8%8B%E8%BD%BD%E9%97%AE%E9%A2%98%E4%B9%8B%E5%88%86%E5%9D%97%E4%B8%8A%E4%BC%A0%EF%BC%88%E6%96%AD%E7%82%B9%E7%BB%AD%E4%BC%A0%EF%BC%89/