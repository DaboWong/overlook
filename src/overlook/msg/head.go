package msg

import (
	"overlook/def"
	"overlook/ds"
)

/*****
//TK消息包头定义
typedef struct tagTKHEADER
{
DWORD	dwMagic;			//消息魔数
DWORD	dwSerial;			//序列号
WORD	wOrigine;			//消息来源
WORD	wReserve;			//保留
DWORD	dwType;				//消息类型
DWORD	dwParam;			//消息参数（消息版本，返回值，标志位等）
DWORD	dwLength;			//实际数据长度，不包括消息头
}TKHEADER,*PTKHEADER;
*/

type TKHeader struct {
	Magic   def.DWORD //消息魔数
	Serial  def.DWORD //序列号
	Origin  def.WORD  //消息来源
	Reserve def.WORD  //保留
	Type    def.DWORD //消息类型
	Param   def.DWORD //消息参数（消息版本，返回值，标志位等）
	Length  def.DWORD //实际数据长度，不包括消息头
}

const (
	//Origin 预先定义好
	Origin def.WORD = 258

	//WS转发到match的消息
	WebToMatchQueryID string = "WebToMatchQuery"
	MatchToWebNtfID string = "MatchToWebNtf"
)

type JsonHeader struct {
	ID string `json:"id"`
}

type WebToMatchQuery struct {
	REQ      int32
	ACK      int32
	ClientID int32
	Transfer interface{}
	Return   ds.Notifier
}
