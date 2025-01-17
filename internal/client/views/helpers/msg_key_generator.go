package helpers

import "fmt"

func GenMsgKey(msg any) string {
	return fmt.Sprintf("%T-Handler", msg)
}
