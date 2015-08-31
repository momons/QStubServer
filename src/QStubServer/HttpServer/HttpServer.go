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
	"net/http/httputil"
	"strconv"
)

// リクエストヘッダ退避用
var reqestHeaderBytes []byte
// リクエストボディ退避用
var reqestBodyBytes []byte

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

	// ヘッダ部取得
	reqestHeaderBytes, _ = httputil.DumpRequest(r, false)

	// ボディ部読み込み
	bodyData := new(bytes.Buffer)
	bodyData.ReadFrom(r.Body)
	reqestBodyBytes = bodyData.Bytes()
	defer r.Body.Close()

	// ログ出力
	outputReqest(r.URL.String(), reqestHeaderBytes, reqestBodyBytes)
	ConsoleLog.Info("リクエスト受信開始")
	ConsoleLog.Output(string(reqestHeaderBytes) + string(reqestBodyBytes))
	ConsoleLog.Info("リクエスト受信終了")

	// 該当のURLがあるかを検索
	convertEntity, replaseStr := ConvertInfo.SearchURL(r.URL.String(), string(reqestHeaderBytes), string(reqestBodyBytes))
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

	ConsoleLog.Info("転送リクエスト開始")
	// リクエストヘッダ設定
	ConsoleLog.Output(fmt.Sprintf("%s %s", r.Proto, r.Method))
	for key, _ := range r.Header {
		httpClient.RequestHeader[key] = r.Header.Get(key)
		ConsoleLog.Output(fmt.Sprintf("%s : %s", key, r.Header.Get(key)))
	}

	// リクエストボディ
	httpClient.RequestBody = reqestBodyBytes
	
	ConsoleLog.Output("\n" + string(reqestBodyBytes))

	// コンテントタイプ指定がある場合は設定
	if len(contentType) > 0 {
		contentTypeBk := httpClient.RequestHeader["Content-Type"]
		httpClient.RequestHeader["Content-Type"] = contentType
		ConsoleLog.Info(fmt.Sprintf("Content-Type変更 %s → %s", contentTypeBk, contentType))
	}

	ConsoleLog.Info("転送リクエスト終了")

	// リクエスト送信＆レスポンス受信
	httpClient.Send()

	// ヘッダ情報設定
	w.WriteHeader(httpClient.ResponseHttpStatus)
	for key, obj := range httpClient.ResponseHeader {
		w.Header().Set(key, obj)
	}

	// ボディ情報設定
	w.Write(httpClient.ResponseBody)

	// レスポンスログ出力
	outputResponse(r.URL.String(), httpClient.ResponseHeaderBytes, httpClient.ResponseBody)
	ConsoleLog.Info("レスポンス開始")
	ConsoleLog.Output(string(httpClient.ResponseHeaderBytes) + string(httpClient.ResponseBody))
	ConsoleLog.Info("レスポンス終了")

}

// ファイルを読み込んで返却
func readFile(filePath string, contentType string, w http.ResponseWriter, r *http.Request) {

	ConsoleLog.Info(fmt.Sprintf("ファイル転送 %s", filePath))

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
		ConsoleLog.Output(fmt.Sprintf("Content-Type %s", contentType))
	} else {
		// ファイル名からタイプ設定
		fileContentType := ContentType.GetContentType(filePath)
		w.Header().Set("Content-Type", fileContentType)
		ConsoleLog.Output(fmt.Sprintf("Content-Type %s", fileContentType))
	}

	// ボディ情報設定
	w.Write(data)

	ConsoleLog.Output(string(data))
}

// リクエスト情報を出力
func outputReqest(url string, header []byte, body []byte) {

	// リクエストヘッダログ出力
	OutputLog.Output(OutputLog.LogTypeReqestHeader, url, header)

	// リクエストボディログ出力
	OutputLog.Output(OutputLog.LogTypeReqestBody, url, body)
}

// レスポンス情報を出力
func outputResponse(url string, header []byte, body []byte) {

	// リクエストヘッダログ出力
	OutputLog.Output(OutputLog.LogTypeResponseHeader, url, header)

	// リクエストボディログ出力
	OutputLog.Output(OutputLog.LogTypeResponseBody, url, body)
}
