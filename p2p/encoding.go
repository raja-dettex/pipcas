package p2p

import (
	"fmt"
	"io"
	"strings"
)

type Decoder interface {
	Decode(io.Reader, *RPCMessage) error
}

type DefaultDecoder struct {
}

func (dec *DefaultDecoder) Decode(r io.Reader, msg *RPCMessage) error {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil {
			return err
		}
		mPStr := string(buf[:n])
		aStr := strings.Split(mPStr, " ")
		fmt.Println("astr ", aStr)
		if aStr[0] == "write" {
			if len(aStr) < 2 {
				return fmt.Errorf("invalid message")
			}
			msg.Header = RPCWriteHeader
			msg.Key = []byte(aStr[1])
			str := ""
			for i := 2; i < len(aStr); i++ {
				str += aStr[i]
				str += " "
			}
			msg.Paylaod = []byte(str)
			return nil
		} else if aStr[0] == "read" {
			if len(aStr) < 1 {
				return fmt.Errorf("invalid message")
			}
			msg.Header = RPCReadHeader
			msg.Key = []byte(aStr[1])
			return nil
		} else if aStr[0] == "delete" {
			if len(aStr) < 1 {
				return fmt.Errorf("invalid message")
			}
			msg.Header = RPCRemoveHeader
			msg.Key = []byte(aStr[1])
			return nil
		}

	}
	return nil
}
