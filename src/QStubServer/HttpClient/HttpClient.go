// HTTPクライアント
package HttpClient

import (
	"QStubServer/ConsoleLog"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strconv"
	"strings"
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
		ConsoleLog.Error(fmt.Sprintln("データの送信に失敗しました: %v", err))
		return false
	}

	// ヘッダ情報を設定
	for key, obj := range cl.RequestHeader {
		req.Header.Set(key, obj)
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintln("データの送信に失敗しました: %v", err))
		return false
	}
	defer resp.Body.Close()

	// ヘッダ情報取得
	cl.ResponseHeaderBytes, err = httputil.DumpResponse(resp, false)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintln("データの受信に失敗しました: %v", err))
		return false
	}

	// ヘッダ情報解析
	cl.getResponseHeaders(cl.ResponseHeaderBytes)

	// ボディ部取得
	cl.ResponseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		ConsoleLog.Error(fmt.Sprintln("データの受信に失敗しました:", err))
		return false
	}

	return true
}

// ヘッダ情報解析
func (cl *HTTPClient) getResponseHeaders(dumpResp []byte) {

	isExistHttpStatus := false

	// ヘッダ情報解析
	headers := strings.Split(string(dumpResp), "\r\n")
	for _, header := range headers {
		if !isExistHttpStatus {
			// HTTPステータスを取得
			assined := regexp.MustCompile("HTTP\\/[\\s]*([0-9.]+)[\\s]*([0-9]+)[\\s]*([a-zA-Z ]+)$")
			assinedGp := assined.FindStringSubmatch(header)
			if assinedGp != nil {
				ConsoleLog.Output(fmt.Sprintf("%v", assinedGp))
				cl.ResponseHttpStatus, _ = strconv.Atoi(assinedGp[2])
				// 取得フラグON
				isExistHttpStatus = true
				continue
			}
		}
		if isExistHttpStatus {
			assined := regexp.MustCompile("[\\s]*([\\S]+)[\\s]*:(.+)$")
			assinedGp := assined.FindStringSubmatch(header)
			if len(assinedGp) == 3 {
				cl.ResponseHeader[assinedGp[1]] = strings.Trim(assinedGp[2], " ")
			}
		}
	}
}

