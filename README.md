# MIJS Tech Camp 2019 Winter - WebAssembly

## WebAssembly とは

## Golang で WebAssembly

Go 1.11 以降で WebAssembly がサポートされた。

``` go
package main

import "fmt"

func main() {
    fmt.Println("Hello, WebAssembly!")
}
```

WebAssembly としてコンパイルするには、`GOOS=js` 及び `GOARCH=wasm` を環境変数として設定する:

``` console
$ GOOS=js GOARCH=wasm go build -o main.wasm
```

これを実行すると、main.wasm という名前の実行可能な WebAssemlby モジュールファイルとしてビルドされる。
ファイルの拡張子を .wasm にしておくと、より簡単に正確な Content-Type ヘッダー付きで HTTP でサーブされる。

Go で生成された main.wasm をブラウザー上で実行するには、JavaScript サポートファイルと、これらを全部いっしょにつなげる HTML ページが必要。

JavaScript サポートファイルをコピー:

``` console
$ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

index.html を作成:

``` html
<html>
  <head>
    <meta charset="utf-8">
    <script src="wasm_exec.js"></script>
    <script>
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
      go.run(result.instance);
    });
    </script>
  </head>
  <body></body>
</html>
```

ブラウザーが `WebAssembly.instantiateStreaming` をまだサポートしていない場合は、以下の polyfill を使える:

``` js
if (!WebAssembly.instantiateStreaming) { // polyfill
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}
```

あとは、3つのファイル（index.html、wasm_exec.js、main.wasm）を Web サーバーでサーブする。

`goexec` を使う場合:

``` console
$ goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'
```

または、自分で基本的な HTTP サーバーを作る:

``` go
package main

import (
    "flag"
    "log"
    "net/http"
)

var (
    listen = flag.String("listen", ":8080", "listen address")
    dir    = flag.String("dir", ".", "directory to serve")
)

func main() {
    flag.Parse()
    log.Printf("listening on %q...", *listen)
    err := http.ListenAndServe(*listen, http.FileServer(http.Dir(*dir)))
    log.Fatalln(err)
}
```

http://localhost:8080/index.html を開き、JavaScript のデバッグコンソールを開くと、出力を見れる。

## 参考リンク

- [Go Wiki](https://github.com/golang/go/wiki/WebAssembly)

