package def

import (
	"bytes"
	"crypto/md5"
	"crypto/rc4"
	"fmt"
	"io"
	"log"
	"overlook/cfg"
)

type (
	DWORD uint32
	WORD  uint16
)

//消息类型（tagTKHEADER::dwType） 定义
//高8位为类型定义区间
const (
	TK_REQ DWORD = 0x00000000 //请求消息类型，第32位为0
	TK_ACK DWORD = 0x80000000 //应答消息类型，第32位为1
	TK_ENC DWORD = 0x40000000 //加密消息类型，第31位为1
	TK_ZIP DWORD = 0x20000000 //压缩消息类型，第30位为1
	TK_SHR DWORD = 0x10000000 //共享消息类型，第29位为1
	TK_PRI DWORD = 0x08000000 //优先消息类型，第28位为1
	TK_DLY DWORD = 0x04000000 //延迟消息类型，第27位为1,只对TK_SHR类型的消息有效
	TK_SEC DWORD = 0x01000000 //安全消息类型，第26位为1
)

const (
	TK_ACKRESULT_SUCCESS    DWORD = 0
	TK_ACKRESULT_FAILED           = 1
	TK_ACKRESULT_SVRBUSY          = 2
	TK_ACKRESULT_LOWVERSION       = 3 //版本太低
	TK_ACKRESULT_NOTFINDOBJ       = 4 //对象没有找到
	TK_ACKRESULT_OBJEXIST         = 5 //对象已经存在
	TK_ACKRESULT_WAITASYNC        = 6 //等待异步返回结果
)

const (
	WATCH_TABLE_OK            int32 = 0 //观看ok
	WATCH_TABLE_MISSING             = 1 //没这个桌
	SOCKET_ASSISTUSER_ORIGINE WORD  = (258)
)

func GetMsgType(id DWORD) DWORD {
	return ^TK_ACK & id
}

func Encrypt(psw string) ([]byte, error) {
	md5Hash := md5.New()
	io.WriteString(md5Hash, cfg.GConfig.Judge.Password)

	md5psw := fmt.Sprintf("%x", md5Hash.Sum(nil))
	//log.Printf("psw md5 : %v", md5psw)
	rc4cipher, err := rc4.NewCipher(bytes.NewBufferString(md5psw).Bytes())
	if err != nil {
		log.Fatalln("rc4 en-crypt error, login failed")
		return nil, err
	}
	rc4Psw := make([]byte, len(md5psw))
	rc4cipher.XORKeyStream(rc4Psw, bytes.NewBufferString(psw).Bytes())
	return rc4Psw, nil
}
