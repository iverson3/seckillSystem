package process

import (
	"seckillsystem/test/partUpDown/client/model"
)

type UploadMgr struct {

}

func (um *UploadMgr) NewLocalFileInfo(localPath string, bufSize int) *model.LocalFile {
	return &model.LocalFile{
		LocalFileMeta: model.LocalFileMeta{
			Path: localPath,
		},
		BufSize: bufSize,
	}
}
