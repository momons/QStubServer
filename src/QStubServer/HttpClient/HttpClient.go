// HTTPクライアント
package HttpClient

import (
	"QStubServer/ConsoleLog"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
)

// HTTPクライアント
type HTTPClient struct {
	// リクエストメソッド
	RequestMethod string
	// リクエストURL
	RequestUrl string
	// リクエストヘッダ情報
	RequestHeader map[string]string
	// リクエストボディ情報
	RequestBody []byte
	// レスポンスHTTPステータス
	ResponseHttpStatus int
	// レスポンスヘッダ
	ResponseHeader map[string]string
	// レスポンスヘッダ解析前
	ResponseHeaderBytes []byte
	// レスポンスボディ
	ResponseBody []byte
}

// 初期化
func Init() *HTTPClient {
	// 値を初期化
	httpClient := new(HTTPClient)
	httpClient.RequestMethod = ""
	httpClient.RequestUrl = ""
	httpClient.RequestHeader = map[string]string{}
	httpClient.RequestBody = []byte{}
	httpClient.ResponseHttpStatus = 0
	httpClient.ResponseHeader = map[string]string{}
	httpClient.ResponseBody = []byte{}
	return httpClient
}

// 送信
func (cl *HTTPClient) Send() bool {

	// メソッド、URL、ボディ部を設定
	req, err := http.NewRequest(cl.RequestMethod, cl.RequestUrl, bytes.NewReader(cl.RequestBody))
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("データの送信に失敗しました: %v", err))
		return false
	}

	// ヘッダ情報を設定
	for key, obj := range cl.RequestHeader {
		req.Header.Set(key, obj)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintf("データの送信に失敗しました: %v", err))
		return false
	}

	// ヘッダ情報取得
	cl.ResponseHeaderBytes, _ = httputil.DumpResponse(resp, false)
	cl.ResponseHttpStatus = resp.StatusCode
	for key, _ := range resp.Header {
		cl.ResponseHeader[key] = resp.Header.Get(key)
	}

	// ボディ部取得
	bodyData := new(bytes.Buffer)
	bodyData.ReadFrom(resp.Body)
	cl.ResponseBody = bodyData.Bytes()
	defer resp.Body.Close()

	return true
}
