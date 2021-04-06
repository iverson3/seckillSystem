package common

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type LocalFileMeta struct {
	Path     string `json:"path"`     // 本地路径
	Length   int64  `json:"length"`   // 文件大小
	SliceMD5 []byte `json:"slicemd5"` // 文件前 requiredSliceLen (256KB) 切片的 md5 值
	MD5      []byte `json:"md5"`      // 文件的 md5
	CRC32    uint32 `json:"crc32"`    // 文件的 crc32
	ModTime  int64  `json:"modtime"`  // 修改日期
}

// EqualLengthMD5 检测md5和大小是否相同
func (lfm *LocalFileMeta) EqualLengthMD5(m *LocalFileMeta) bool {
	if lfm.Length != m.Length {
		return false
	}
	if bytes.Compare(lfm.MD5, m.MD5) != 0 {
		return false
	}
	return true
}

// CompleteAbsPath 补齐绝对路径
func (lfm *LocalFileMeta) CompleteAbsPath() {
	if filepath.IsAbs(lfm.Path) {
		return
	}

	absPath, err := filepath.Abs(lfm.Path)
	if err != nil {
		return
	}

	lfm.Path = absPath
}

// GetFileSum 获取文件的大小, md5, 前256KB切片的 md5, crc32
func GetFileSum(localPath string, cfg *SumConfig) (lf *LocalFile, err error) {
	file, err := os.Open(localPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fileStat.IsDir() {
		return nil, fmt.Errorf("sum %s: is a directory", localPath)
	}

	lf = &LocalFile{
		File: file,
		LocalFileMeta: LocalFileMeta{
			Path:   localPath,
			Length: fileStat.Size(),
		},
	}

	lf.Sum(*cfg)

	return lf, nil
}