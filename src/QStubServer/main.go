// メイン
package main

import (
	"QStubServer/ConvertInfo"
	"QStubServer/ContentType"
	"QStubServer/HttpServer"
	"flag"
	"os"
)

// コマンド情報
// ポート番号
var commandPort int
// 設定ファイルパス
var commandSettingPath string
// コンテントタイプリストファイルパス
var commandContetTypeListPath string

// 終了コード
var exitCode = 0

// メイン
func main() {

	// セットアップ
	isSuccess := setup()
	if !isSuccess {
		exitCode = 1
	}

	os.Exit(exitCode)
}

// セットアップ
func setup() bool {

	// コマンドライン取得
	setupCommand()

	// コンバート情報設定
	isSuccess := ConvertInfo.Setup(commandSettingPath)
	if !isSuccess {
		return false
	}

	// コンテントタイプリスト設定
	isSuccess = ContentType.Setup(commandContetTypeListPath)
	if !isSuccess {
		return false
	}

	// HTTPサーバセットアップ
	isSuccess = HttpServer.Setup(commandPort)
	if !isSuccess {
		return false
	}

	return true
}

// コマンドライン設定
func setupCommand() {
	// ポート
	flag.IntVar(&commandPort, "port", 8080, "ポートを指定して下さい。")
	// 設定ファイルパス
	flag.StringVar(&commandSettingPath, "setting", "設定ファイルパス", "設定ファイルパスを指定して下さい。")
	// コンテントタイプリストファイルパス
	flag.StringVar(&commandContetTypeListPath, "contenttype", "ContetTypeリストファイルパス", "ContetTypeリストファイルパスを指定して下さい。")

	flag.Parse()

	// 設定ファイルパス
	if commandSettingPath == "設定ファイルパス" {
		commandSettingPath = ""
	}
	// コンテントタイプリストファイルパス
	if commandContetTypeListPath == "ContetTypeリストファイルパス" {
		commandContetTypeListPath = ""
	}

}
