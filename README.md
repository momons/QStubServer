# Q Stub Server

いろいろな開発においてAPIへのアクセスやテスト等で自分自身がほしいな〜と思ったスタブサーバを作成しました。

## 機能

- 特定URLからローカルファイルを転送します。
- 特定URLから情報を別サーバへ転送します。
- リクエスト、レスポンスの内容をファイル出力します。

#### 特定のURLからローカルファイルを転送する。

`Setting.json`の内容を編集します。

アクセスを全て`/user/hoge/webroot/`配下のローカルファイルのアクセスに変換する。<br>
正規表現で対応します。

```JSON
{
  "convertList": [
    {
      "srcURL" : "/(.*)",
      "destPath" : "/user/hoge/webroot/$1"
    }
  ]
}
```

`/`へのアクセスに`index.html`を返却する。その他は上記と同じに。<br>
上から順に処理されていくため、優先するものは上位に設定します。

```JSON
{
  "convertList": [
    {
      "srcURL" : "/",
      "destPath" : "/user/hoge/webroot/index.html"
    },
    {
      "srcURL" : "/(.*)",
      "destPath" : "/user/hoge/webroot/$1"
    }
  ]
}
```

ファイルに別途ContentTypeを指定する。<br>
`index.html`の場合`Content-Type: application/octet-stream`が指定される。<br>
省略時は`ContentTypeList.json`から拡張子で自動で設定されます。

```JSON
{
  "convertList": [
    {
      "srcURL" : "/",
      "destPath" : "/user/hoge/webroot/index.html",
      "contentType" : "application/octet-stream"
    },
    {
      "srcURL" : "/(.*)",
      "destPath" : "/user/hoge/webroot/$1"
    }
  ]
}
```

#### 特定のURLを別サーバへ転送する。

アクセスを全て`http://192.168.1.10:9999/`に転送する。

```JSON
{
  "convertList": [
    {
      "srcURL" : "/(.*)",
      "destURL" : "http://192.168.1.10:9999/$1"
    }
  ]
}
```

`/hoge/`の場合ローカルファイルにアクセスする。その他は上記と同じに。

```JSON
{
  "convertList": [
    {
      "srcURL" : "/hoge/(.*)",
      "destPath" : "/user/hoge/webroot/$1"
    },
    {
      "srcURL" : "/(.*)",
      "destURL" : "http://192.168.1.10:9999/$1"
    }
  ]
}
```

## コマンドについて

ポートを9999に設定する。(デフォルトは8080)

```
> QStubServer -port 9999
```

`Setting.json`のパスを指定する。(デフォルトは`[カレントディレクトリ]+Setting.json`)

```
> QStubServer -setting /usr/hoge/QStubServer/Setting.json
```

`ContentTypeList.json`のパスを指定する。(デフォルトは`[カレントディレクトリ]+ContentTypeList.json`)

```
> QStubServer -contenttype /usr/hoge/QStubServer/ContentTypeList.json
```

リクエスト、レスポンス内容をファイルに出力する。<br>
`/usr/hoge/QStubServer/log`配下へ内容を出力します。

```
> QStubServer -outputlog /usr/hoge/QStubServer/log
```

## ビルド環境について

このライブラリのビルドは<br>
https://github.com/constabulary/gb<br>
で行っております。
