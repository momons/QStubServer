// コンバート情報
package ConvertInfo

import (
	"QStubServer/ConsoleLog"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// コンバート情報Entity
type ConvertInfoEntity struct {
	ConvertList []ConvertInfoDetailEntity `json:"convertList"`
}

// コンバート情報詳細Entity
type ConvertInfoDetailEntity struct {
	// 元URL
	SrcURL string `json:"srcURL"`
	// 接続先URL
	DestURL string `json:"destURL"`
	// 接続先パス
	DestPath string `json:"destPath"`
	// コンテントタイプ
	ContentType string `json:"contentType"`
}

// 設定ファイルパス
var filePath string

// 変換情報
var ConvertInfo *ConvertInfoEntity

// 初期設定
func Setup(settingFilePath string) bool {

	// 退避
	filePath = settingFilePath

	// ファイルパスがなかったら
	if len(filePath) <= 0 {
		// カレントディレクトリ取得
		curDir, _ := os.Getwd()
		filePath = curDir + "/Setting.json"
	}

	// ファイルパスを出力
	ConsoleLog.Info(fmt.Sprintf("設定ファイルパス: %s", filePath))

	// ファイル読み込み
	var isSuccess bool
	ConvertInfo, isSuccess = read()
	if !isSuccess {
		return false
	}

	return true
}

// ファイル読込
func read() (*ConvertInfoEntity, bool) {

	// ファイル読み込み
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("設定ファイルの読み込みに失敗しました。%v", err))
		return nil, false
	}

	var entity ConvertInfoEntity

	// JSONに変換
	reader := strings.NewReader(string(data))
	dec := json.NewDecoder(reader)
	dec.Decode(&entity)

	// ログ出力
	for _, obj := range entity.ConvertList {
		ConsoleLog.Output(fmt.Sprintf("%v", obj))
	}

	return &entity, true
}

// 検索
func SearchURL(url string) (*ConvertInfoDetailEntity, string) {

	// 上から順に該当するものを検索
	for _, entity := range ConvertInfo.ConvertList {
		assined := regexp.MustCompile(entity.SrcURL)
		assinedGp := assined.FindStringSubmatch(url)
		if assinedGp != nil {
			if len(entity.DestURL) > 0 {
				return &entity, assined.ReplaceAllString(url, entity.DestURL)
			} else {
				return &entity, assined.ReplaceAllString(url, entity.DestPath)
			}
		}
	}

	return nil, ""
}
