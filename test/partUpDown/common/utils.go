package common

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
)

// 网络传输者
type Transfer struct {
	Conn net.Conn
	// 数据传输时的缓冲空间
	Buf [70*1024*1024]byte
}

// 从连接中读取客户端发送过来的数据
func (this *Transfer) ReadPkg() (mess Message, err error) {
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		return
	}

	// 获取到数据的长度信息 (借助binary包方法将byte数据转为int数字)
	pkgLen := binary.BigEndian.Uint32(this.Buf[:4])

	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		return
	}

	err = json.Unmarshal(this.Buf[:pkgLen], &mess)
	if err != nil {
		// 自定义error信息
		err = errors.New("反序列化失败")
		return
	}
	return
}

// 向该连接的客户端发送数据
func (this *Transfer) WritePkg(data []byte) (err error)  {
	// 先向客户端发送数据的长度信息
	pkgLen := uint32(len(data))

	var pkgLenByte [4]byte // 4 * 8 = 32 (uint32)
	// 将一个int类型的数字，转成byte切片
	binary.BigEndian.PutUint32(pkgLenByte[0:4], pkgLen)
	n, err := this.Conn.Write(pkgLenByte[0:4])
	if n != len(pkgLenByte) || err != nil {
		return
	}

	// 接着向客户端发送真正的数据
	n, err = this.Conn.Write(data)

	if n != int(pkgLen) {
		return errors.New("发送给客户端的数据的长度与数据实际长度不匹配")
	}
	return
}

func Md5String(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}