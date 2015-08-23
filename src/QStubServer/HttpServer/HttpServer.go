// HTTPサーバ
package HttpServer

import (
	"QStubServer/ConsoleLog"
	"QStubServer/ContentType"
	"QStubServer/ConvertInfo"
	"QStubServer/HttpClient"
	"QStubServer/OutputLog"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// セットアップ
func Setup(portNo int) bool {

	// ハンドラを登録してウェブページを表示させる
	http.HandleFunc("/", httpHandler)

	// ポート設定
	portNoStr := strconv.Itoa(portNo)

	// スタート
	ConsoleLog.Info(fmt.Sprintf("Qstub起動 ポート番号: %s", portNoStr))
	err := http.ListenAndServe(":"+portNoStr, nil)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("Qstub起動エラー: %v", err))
		return false
	}

	return true
}

// HTTPハンドラ
func httpHandler(w http.ResponseWriter, r *http.Request) {

	ConsoleLog.InfoStrong(fmt.Sprintf("リクエスト受信: %s　　　　　　　　　　　　　", r.URL.String()))

	// ログ出力
	outputReqest(r)

	// 該当のURLがあるかを検索
	convertEntity, replaseStr := ConvertInfo.SearchURL(r.URL.String())
	if convertEntity == nil {
		// 情報返却
		errMsg := fmt.Sprintf("該当のURLに対する設定が見つかりません。%s", r.URL.String())
		ConsoleLog.Warning(errMsg)
		w.WriteHeader(404)
		fmt.Fprintf(w, errMsg)
		return
	}

	if len(convertEntity.DestURL) > 0 {
		// 別のサイトへ通信
		connectOtherSite(replaseStr, convertEntity.ContentType, w, r)
	} else {
		// ローカルアクセス
		readFile(replaseStr, convertEntity.ContentType, w, r)
	}

}

// 別サイトへ通信
func connectOtherSite(url string, contentType string, w http.ResponseWriter, r *http.Request) {

	httpClient := HttpClient.Init()
	httpClient.RequestMethod = r.Method
	httpClient.RequestUrl = url

	// リクエストヘッダ設定
	ConsoleLog.Info("リクエストヘッダ開始")
	ConsoleLog.Output(fmt.Sprintf("%s %s", r.Proto, r.Method))
	for key, _ := range r.Header {
		httpClient.RequestHeader[key] = r.Header.Get(key)
		ConsoleLog.Output(fmt.Sprintf("%s : %s", key, r.Header.Get(key)))
	}
	ConsoleLog.Info("リクエストヘッダ終了")

	// コンテントタイプ指定がある場合は設定
	if len(contentType) > 0 {
		httpClient.RequestHeader["Content-Type"] = contentType
	}

	// リクエストボディ
	ConsoleLog.Info("リクエストボディ開始")
	bodyData := new(bytes.Buffer)
	bodyData.ReadFrom(r.Body)
	httpClient.RequestBody = bodyData.Bytes()
	ConsoleLog.Output(bodyData.String())
	ConsoleLog.Info("リクエストボディ終了")

	// リクエスト送信＆レスポンス受信
	httpClient.Send()

	// ヘッダ情報設定
	ConsoleLog.Info("レスポンスヘッダ開始")
	w.WriteHeader(httpClient.ResponseHttpStatus)
	for key, obj := range httpClient.ResponseHeader {
		w.Header().Set(key, obj)
		ConsoleLog.Output(fmt.Sprintf("%s : %s", key, obj))
	}
	ConsoleLog.Info("レスポンスヘッダ終了")

	ConsoleLog.Info("レスポンスボディ開始")
	ConsoleLog.Output(string(httpClient.ResponseBody))
	ConsoleLog.Info("レスポンスボディ終了")

	// ボディ情報設定
	w.Write(httpClient.ResponseBody)

	// レスポンスログ出力
	outputResponse(r.URL.String(), httpClient.ResponseHeaderBytes, httpClient.ResponseBody)
}

// ファイルを読み込んで返却
func readFile(filePath string, contentType string, w http.ResponseWriter, r *http.Request) {

	// ファイル読み込み
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 404返却
		errMsg := fmt.Sprintf("ファイルの読み込みに失敗しました。%v", err)
		ConsoleLog.Error(errMsg)
		w.WriteHeader(404)
		fmt.Fprintf(w, errMsg)
		return
	}

	if len(contentType) > 0 {
		// コンテントタイプ指定がある場合は設定
		w.Header().Set("Content-Type", contentType)
	} else {
		// ファイル名からタイプ設定
		w.Header().Set("Content-Type", ContentType.GetContentType(filePath))
	}

	// ボディ情報設定
	w.Write(data)
}

// リクエスト情報を出力
func outputReqest(r *http.Request) {

	// リクエストヘッダログ出力
	logMsg := ""
	for key, _ := range r.Header {
		logMsg += fmt.Sprintf("%s : %s\n", key, r.Header.Get(key))
	}
	OutputLog.Output(OutputLog.LogTypeReqestHeader, r.URL.String(), []byte(logMsg))

	// リクエストボディログ出力
	bodyData := new(bytes.Buffer)
	bodyData.ReadFrom(r.Body)
	OutputLog.Output(OutputLog.LogTypeReqestBody, r.URL.String(), bodyData.Bytes())
}

// レスポンス情報を出力
func outputResponse(url string, header []byte, body []byte) {

	// リクエストヘッダログ出力
	OutputLog.Output(OutputLog.LogTypeResponseHeader, url, header)

	// リクエストボディログ出力
	OutputLog.Output(OutputLog.LogTypeResponseBody, url, body)
}
