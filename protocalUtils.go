package hsp_utils

import (
	"encoding/binary"
	"github.com/smallnest/ringbuffer"
	"log"
)

//消息头默认长度
const SELFMSGHEADLEN = 20

var (
	PROTOCAL_OK         = 1 //解析成功
	PROTOCAL_DIV        = 2 //数据未接收完整，暂不解析
	PROTOCAL_ERRMSG     = 3 //错误的消息
	PROTOCAL_ERRPROCESS = 4 //解析时发生错误
)

type SelfProtocal struct {
	IdentifyComm int    //通讯标志位
	CmdID        int    //命令码
	MsgLen       int    //消息体长度
	BodyDatas    []byte //消息体数据
	Status       int
}

func AnalysisSelfProtocalsMsg(rbuf *ringbuffer.RingBuffer) *SelfProtocal {
	res := &SelfProtocal{}

	headBuffer := make([]byte, SELFMSGHEADLEN)
	rmvBuffer := make([]byte, 1)
	for {
		//判断ringbuf内容是否超过或等于消息头长度
		if rbuf.Length() < SELFMSGHEADLEN {
			res.Status = PROTOCAL_DIV
			break
		}
		//拷贝前20个字节消息头
		n, err := rbuf.Copy(headBuffer)
		if err != nil || n != SELFMSGHEADLEN {
			log.Println("protocalUtils:read error from ringbuf:", err)
			res.Status = PROTOCAL_ERRPROCESS
			break
		}

		//判断消息头标志
		if headBuffer[0] == 0x05 && headBuffer[1] == 0x0A && headBuffer[2] == 0x05 && headBuffer[3] == 0x0A {
			tmpSlice := headBuffer[9:11]
			msgLen := (int)(binary.LittleEndian.Uint16(tmpSlice))
			msgLen = (msgLen/492)*512 + SELFMSGHEADLEN + msgLen%492
			if rbuf.Length() >= msgLen {
				//ringbuf中的数据多于消息体
				msgItem := make([]byte, msgLen)
				n, err = rbuf.Read(msgItem)
				if n != msgLen || err != nil {
					log.Println("protocalUtils:read error from ringbuf1:", err)
					res.Status = PROTOCAL_ERRPROCESS
					break
				}
				//处理单条指令
				if err = singleMsgProcess(msgItem, msgLen, res); err != nil {
					return res
				}
			} else {
				//ringbuf未收到足够的数据
				res.Status = PROTOCAL_DIV
				break
			}
		} else {
			rbuf.Read(rmvBuffer)
			log.Println("protocalUtils:ringbuf romove byte:", rmvBuffer)
		}
	}

	return res
}

func singleMsgProcess(msgItem []byte, msgLen int, res *SelfProtocal) error {

	return nil
}
