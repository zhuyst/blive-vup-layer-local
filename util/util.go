package util

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v3/pkg/application"
	"math/rand"
	"runtime/debug"
)

func MapToStruct(m map[string]interface{}, s interface{}) error {
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, s)
}

func IsRepeatedChar(s string) bool {
	if s == "" {
		return false
	}
	cs := []rune(s)
	char := cs[0]

	for i := 1; i < len(cs); i++ {
		if cs[i] != char {
			return false
		}
	}
	return true
}

func GetRandomStr(randomArr []string) string {
	i := rand.Intn(len(randomArr))
	return randomArr[i]
}

func RunGr(fn func()) {
	go func() {
		defer Recover()
		fn()
	}()
}

func Recover() {
	if err := recover(); err != nil {
		log.Errorf("panic: %s, stack: %s", err, string(debug.Stack()))

		dialog := application.ErrorDialog()
		dialog.SetTitle("程序发生崩溃，已恢复")
		dialog.SetMessage(fmt.Sprintf("程序发生崩溃，已恢复\npanic: %s\nstack: %s", err, string(debug.Stack())))
		dialog.Show()
	}
}

func ShowErrorDialog(errMsg string) {
	dialog := application.ErrorDialog()
	dialog.SetTitle("程序发生错误")
	dialog.SetMessage(fmt.Sprintf("程序发生错误\npanic: %s\nstack: %s", errMsg, string(debug.Stack())))
	dialog.Show()
}
