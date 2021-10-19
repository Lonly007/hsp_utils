package hsp_utils

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/smallnest/ringbuffer"
	"log"
)

//消息头默认长度
const SELFMSGHEADLEN = 20
const PACKAGELEN = 512

var (
	PROTOCAL_OK         = 1 //解析成功
	PROTOCAL_DIV        = 2 //数据未接收完整，暂不解析
	PROTOCAL_ERRMSG     = 3 //错误的消息
	PROTOCAL_ERRPROCESS = 4 //解析时发生错误
)
var (
	ERR_BODYXOR    = errors.New("ERR_BODY_XORSUM")
	ERR_HEADXOR    = errors.New("ERR_HEAD_XORSUM")
	ERR_HEADLENGTH = errors.New("ERR_HEAD_LENGTH")
	ERR_CMDID      = errors.New("ERR_MSG_CMDID")
	ERR_COMMID     = errors.New("ERR_MSG_COMMID")
	ERR_PACKERR    = errors.New("ERR_PACKAGE_BODY_TOOLANG")
)

type SelfProtocal struct {
	IdentifyComm int    //通讯标志位
	CmdID        int    //命令码
	MsgLen       int    //消息体长度
	BodyDatas    []byte //消息体数据
	Status       int
}

//协议解析
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
			msgLen = (msgLen/492)*PACKAGELEN + SELFMSGHEADLEN + msgLen%492
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
				if err = singleMsgProcess(msgItem, msgLen, res); err == nil {
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

	packageCnt := msgLen/PACKAGELEN + 1
	for i := 0; i < packageCnt; i++ {
		packageLen := PACKAGELEN
		if i == packageCnt-1 {
			//最后一包时packagelen不一定是512
			packageLen = msgLen % PACKAGELEN
			if packageLen == 0 {
				//最后一包正好是512字节的处理
				packageLen = PACKAGELEN
			}
		}

		//消息体异或校验
		var bodyXor byte = 0
		for j := SELFMSGHEADLEN; j < packageLen; j++ {
			bodyXor = bodyXor ^ msgItem[i*PACKAGELEN+j]
		}
		if bodyXor != msgItem[PACKAGELEN*i+14] {
			// logMsg := "protocalUtils:Error when check xor of body:" + string(msgItem)
			logMsg := "protocalUtils:Error when check xor of body:" + fmt.Sprintf("%d", bodyXor) + fmt.Sprintf("--%d", msgItem[PACKAGELEN*i+14])
			log.Println(logMsg)
			res.Status = PROTOCAL_ERRMSG
			return ERR_BODYXOR
		}

		//消息头校验
		var headXor byte = 0
		for j := 0; j < SELFMSGHEADLEN-1; j++ {
			headXor = headXor ^ msgItem[PACKAGELEN*i+j]
		}
		if headXor != msgItem[PACKAGELEN*i+19] {
			logMsg := "protocalUtils:Error when check xor of head:" + string(msgItem)
			log.Println(logMsg)
			res.Status = PROTOCAL_ERRMSG
			return ERR_HEADXOR
		}

		//cmdID获取和校验
		if i == 0 {
			res.CmdID = int(binary.LittleEndian.Uint32(msgItem[5:9]))
		} else {
			if res.CmdID != int(binary.LittleEndian.Uint32(msgItem[PACKAGELEN*i+5:PACKAGELEN*i+9])) {
				//cmdID不匹配
				logMsg := "protocalUtils:Error when check cmdID of multi-package:" + string(msgItem)
				log.Println(logMsg)
				res.Status = PROTOCAL_ERRMSG
				return ERR_CMDID
			}
		}

		//通讯标志位获取和校验
		if i == 0 {
			res.IdentifyComm = int(msgItem[4:5][0])
		} else {
			if res.IdentifyComm != int(msgItem[PACKAGELEN*i+4 : PACKAGELEN*i+5][0]) {
				//通讯标志位不匹配
				logMsg := "protocalUtils:Error when check comm-identify of multi-package:" + string(msgItem)
				log.Println(logMsg)
				res.Status = PROTOCAL_ERRMSG
				return ERR_COMMID
			}
		}

		//body获取
		if i == packageCnt-1 {
			//最后一包的处理
			body := msgItem[PACKAGELEN*i+SELFMSGHEADLEN:]
			res.BodyDatas = append(res.BodyDatas, body...)
		} else {
			body := msgItem[PACKAGELEN*i+SELFMSGHEADLEN : PACKAGELEN*(i+1)]
			res.BodyDatas = append(res.BodyDatas, body...)
		}
	}

	res.Status = PROTOCAL_OK
	return nil
}

//协议打包
func PackageSelfProtocalsMsg(cmdInfo [5]byte, msgInfo []byte) [][]byte {
	msgLen := len(msgInfo)
	packageNum := msgLen/(PACKAGELEN-SELFMSGHEADLEN) + 1
	sendPackage := make([][]byte, packageNum)

	var frameLen, frameID int
	var bLastFrame bool = false
	if packageNum == 1 {
		//单包协议
		frameLen = msgLen + SELFMSGHEADLEN
		frameID = 1
		bLastFrame = true
	} else {
		frameLen = PACKAGELEN
		frameID = 1
		bLastFrame = false
	}

	firstFrame, err := packageSingleFrame(cmdInfo, msgInfo[:frameLen-SELFMSGHEADLEN], frameID, bLastFrame, msgLen)
	if err != nil {
		log.Println("Error when packet msg cmdInfo: ", cmdInfo)
		log.Println("Error when packet msg msgInfo: ", msgInfo)
		return nil
	} else {
		sendPackage[0] = firstFrame
	}

	//middle frames package (multi-package case, not used here)
	frameLen = PACKAGELEN
	bLastFrame = false
	for i := 1; i < packageNum-1; i++ {
		frameID = i + 1
		middleFrame, err := packageSingleFrame(cmdInfo, msgInfo[i*(PACKAGELEN-SELFMSGHEADLEN):(i+1)*(PACKAGELEN-SELFMSGHEADLEN)], frameID, bLastFrame, msgLen)
		if err != nil {
			log.Println("Error when packet msg cmdInfo: ", cmdInfo)
			log.Println("Error when packet msg msgInfo: ", msgInfo)
			return nil
		} else {
			sendPackage[i] = middleFrame
		}
	}

	//last frame package (multi-package case, not used here)
	frameLen = msgLen%(PACKAGELEN-SELFMSGHEADLEN) + SELFMSGHEADLEN
	bLastFrame = true
	lastFrame, err := packageSingleFrame(cmdInfo, msgInfo[(packageNum-1)*(PACKAGELEN-SELFMSGHEADLEN):], frameID+1, bLastFrame, msgLen)
	sendPackage[packageNum-1] = lastFrame

	return sendPackage

}

func packageSingleFrame(cmdInfo [5]byte, msgInfo []byte, frameID int, bLastFrame bool, msgLen int) ([]byte, error) {
	frameLen := len(msgInfo) + SELFMSGHEADLEN
	if frameLen > PACKAGELEN {
		log.Println("packageSingleFrame receives too lang msg-body")
		return nil, ERR_PACKERR
	}
	msgFrame := make([]byte, frameLen)
	//msg head packet
	//identification
	msgFrame[0] = 0x05
	msgFrame[1] = 0x0A
	msgFrame[2] = 0x05
	msgFrame[3] = 0x0A
	//通信标志位和cmdid
	copy(msgFrame[4:9], cmdInfo[:])
	//length of body
	binary.LittleEndian.PutUint16(msgFrame[9:11], uint16(msgLen))
	//frameID
	binary.LittleEndian.PutUint16(msgFrame[11:13], uint16(frameID))
	//last frame
	if bLastFrame {
		msgFrame[13] = 1
	} else {
		msgFrame[13] = 0
	}
	//body xor
	var bodyXor byte = 0
	for i := 0; i < frameLen-SELFMSGHEADLEN; i++ {
		bodyXor = bodyXor ^ msgInfo[i]
	}
	msgFrame[14] = bodyXor
	//head Xor
	var headXor byte = 0
	for i := 0; i < SELFMSGHEADLEN-1; i++ {
		headXor = headXor ^ msgFrame[i]
	}
	msgFrame[19] = headXor
	//body info
	copy(msgFrame[20:], msgInfo[:])

	return msgFrame, nil
}
