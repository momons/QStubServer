// ログ管理
package ConsoleLog

import (
	"fmt"
)

// ログを出力するか
var IsLogging bool = true

func Output(msg string) {
	if IsLogging {
		fmt.Println(msg)
	}
}

// 情報ログ出力
func Info(msg string) {
	Output("\x1b[32m" + msg + "\x1b[m")
}

// 情報ログ出力(強)
func InfoStrong(msg string) {
	Output("\x1b[46m" + msg + "\x1b[m")
}

// ワーニングログ出力
func Warning(msg string) {
	Output("\x1b[33m" + msg + "\x1b[m")
}

// エラーログ出力
func Error(msg string) {
	Output("\x1b[31m" + msg + "\x1b[m")
}
