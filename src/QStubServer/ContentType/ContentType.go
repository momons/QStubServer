// コンテントタイプ
package ContentType

import (
	"QStubServer/ConsoleLog"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// コンテントタイプ一覧
var ContentTypeList map[string]string

// コンテントタイプファイルパス
var filePath string

// 初期設定
func Setup(contentTypeFilePath string) bool {

	// 退避
	filePath = contentTypeFilePath

	// ファイルパスがなかったら
	if len(filePath) <= 0 {
		// カレントディレクトリ取得
		curDir, _ := os.Getwd()
		filePath = curDir + "/ContentTypeList.json"
	}

	// ファイルパスを出力
	ConsoleLog.Info(fmt.Sprintf("ContetTypeリストファイルパス: %s", filePath))

	// ファイル読み込み
	var isSuccess bool
	ContentTypeList, isSuccess = read()
	if !isSuccess {
		return false
	}

	return true
}

// ファイル読込
func read() (map[string]string, bool) {

	// ファイル読み込み
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("ContetTypeリストファイルの読み込みに失敗しました。%v", err))
		return nil, false
	}

	var list map[string]string

	// JSONに変換
	reader := strings.NewReader(string(data))
	dec := json.NewDecoder(reader)
	dec.Decode(&list)

	// ログ出力
	for key, obj := range list {
		ConsoleLog.Output(fmt.Sprintf("%s : %s", key, obj))
	}

	return list, true
}

// ファイル名からContentTypeを返却する
func GetContentType(filePath string) string {

	assigned := regexp.MustCompile("(.*)(?:\\.([^.]+$))")
	assignedGp := assigned.FindStringSubmatch(filePath)
	if assignedGp != nil && len(assignedGp) == 3 {
		if contentType, ok := ContentTypeList[assignedGp[2]]; ok {
			return contentType
		}
	}

	// 何も見つからない場合はデフォルトを返却
	return "application/octet-stream"
}
