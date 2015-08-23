// ログ出力
package OutputLog

import (
	"QStubServer/ConsoleLog"
	"QStubServer/FileManager"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
)

// ログ種別
// リクエストヘッダ
var LogTypeReqestHeader int = 1

// リクエストボディ
var LogTypeReqestBody int = 2

// レスポンスヘッダ
var LogTypeResponseHeader int = 3

// レスポンスボディ
var LogTypeResponseBody int = 4

// ログ出力先パス
var filePath string

// 初期設定
func Setup(outputLogPath string) bool {

	// 退避
	filePath = outputLogPath
	// 指定がない場合は何もしない
	if len(filePath) <= 0 {
		return true
	}

	// ディレクトリ作成
	err := FileManager.CreateDir(filePath)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("ディレクトリの作成に失敗しました。%v", err))
		return false
	}

	ConsoleLog.Info(fmt.Sprintf("ログ出力先パス: %s", filePath))

	return true
}

// 出力
func Output(logType int, url string, logMsg []byte) {

	// 指定がない場合は何もしない
	if len(filePath) <= 0 {
		return
	}

	// 出力ファイル名設定
	outputFile := time.Now().Format("2006.01.02_15.04.05")
	// URLより一番後ろのパスを設定する
	assined := regexp.MustCompile("\\/([\\w.]*)([#?]+.*)*$")
	assinedGp := assined.FindStringSubmatch(url)
	if assinedGp != nil && len(assinedGp) >= 2 {
		outputFile += "_" + assinedGp[1]
	}
	//
	switch logType {
	case LogTypeReqestHeader:
		outputFile += "_reqhed"
	case LogTypeReqestBody:
		outputFile += "_reqbody"
	case LogTypeResponseHeader:
		outputFile += "_reshed"
	case LogTypeResponseBody:
		outputFile += "_resbody"
	}
	outputFile += ".log"

	// ファイル出力
	err := ioutil.WriteFile(path.Join(filePath, outputFile), logMsg, os.ModePerm)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("ログファイルの出力に失敗しました。%v", err))
	}

}
