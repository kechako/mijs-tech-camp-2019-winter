# MIJS Tech Camp 2019 Winter - WebAssembly

MIJS ニューテクノロジー委員会 2019 年冬合宿 WebAssembly 分科会用のメモ。

Go 言語による WebAssembly についてまとめ。

## WebAssembly とは

Comming soon...?

## サンプルプログラムについて

- [hello](/hello)
  - Hello, WebAssembly
  - ブラウザーのデバッグコンソールに出力される
- [memexport](/memexport)
  - WebAssembly 内のメモリーを JavaScript 側にエクスポートするサンプル。
- [memimport](/memimport)
  - WebAssembly 内に JavaScript からメモリーをインポートするサンプル。
- [calc](/calc)
  - 入力された 2 つの値を足し合わせて結果を表示するだけのサンプル。
- [server](/server)
  - 簡単な Web サーバー

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

## DOM へのインタラクト

`syscall/js` パッケージを使用すると、DOM にインタラクトできる。

``` go
type Value struct {
    // contains filtered or unexported fields
}
```

`Value` は、JavaScript の値を表す。

``` go
func Global() Value
```

`Global` は、JavaScript のグローバルオブジェクト（通常は `window` または `global`）を返す。

``` go
func Null() Value
```

`Null` は、JavaScript の `null` 値を返す。

``` go
func Undefined() Value
```

`Undefined` は、JavaScript の `undefined` 値を返す。

``` go
func ValueOf(x interface{}) Value
```

`ValueOf` は、`x` を JavaScript の値として返す。

| Go                     | JavaScript             |
| ---------------------- | ---------------------- |
| js.Value               | [its value]            |
| js.TypedArray          | typed array            |
| js.Callback            | function               |
| nil                    | null                   |
| bool                   | boolean                |
| integers and floats    | number                 |
| string                 | string                 |
| []interface{}          | new array              |
| map[string]interface{} | new object             |

``` go
func (v Value) Get(p string) Value
```

`Get` は、JavaScript の値 `v` のプロパティ `p` を返す。

``` go
func (v Value) Set(p string, x interface{})
```

`Set` は、JavaScript の値 `v` のプロパティ `p` に `ValueOf(x)` を設定する。

``` go
type Callback struct {
    Value // the JavaScript function that queues the callback for execution
    // contains filtered or unexported fields
}
```

`Callback` は、JavaScript のコールバックとして使用するためにラップされた Go の関数。

``` go
func NewCallback(fn func(args []Value)) Callback
```

`NewCallback` は、ラップされたコールバック関数を返す。

``` go
package main

import (
    "fmt"
    "syscall/js"
)

func main() {
    var cb js.Callback
    cb = js.NewCallback(func(args []js.Value) {
        fmt.Println("button clicked")
        cb.Release() // release the callback if the button will not be clicked again
    })
    js.Global().Get("document").Call("getElementById", "myButton").Call("addEventListener", "click", cb)
}
```

Go 1.12 より、`Callback` 構造体と `NewCallback` 関数は、`Func` 構造体と `FuncOf` 関数にリネームされている。

## JavaScript API

### WebAssembly JavaScript オブジェクト

`WebAssembly` JavaScript オブジェクトは、すべての WebAssembly に関する機能の名前空間として振る舞う。

### WebAssembly.instantiateStreaming() 関数

``` js
Promise<ResultObject> WebAssembly.instantiateStreaming(source, importObject);
```

`WebAssembly.instantiateStreaming()` 関数はソースのストリームから直接 WebAssembly モジュールをコンパイルしてインスタンス化する。

`ResultObject` の2つのフィールド
- `module`
  - コンパイルされた `WebAssembly.Module` オブジェクト。
  - この `Module` は、再度インスタンス化、`postMessage()` 経由での共有、IndexDB へのキャッシュが可能。
- `instance`
  - すべてのエクスポートされた WebAssembly 関数を含む `WebAssembly.Instance` オブジェクト。

## WebAssembly のメモリー空間

WebAssembly からアクセスされる線形メモリーは、WebAssembly.Memory のインスタンスである。
これは、生のバイト列を持つリサイズ可能な ArrayBuffer である。

JavaScript から、WebAssembly の線形メモリーにアクセスする方法として、以下の2つの方法がある:

1. JavaScript で生成した WebAssembly.Memory のインスタンスを、importObject として WebAssemly に渡す。
2. WebAssembly のメモリーをエクスポートする。
  - `mem` という名称でエクスポートした場合、`WebAssembly.Instance.exports.mem` として JavaScript から参照できる。
  - JavaScript からは WebAssembly.Memory のインスタンスとして参照できる。

どちらの方法も、JavaScript からは WebAssembly.Memory のインスタンスとして操作可能なので、このインスタンスの `buffer` プロパティ（メモリーに関連付けられているバッファーを返す）を、`Uint8Array` などの `TypedArray` でオフセット、長さウィ指定してラップすることで、バッファー上の任意のオフセットにアクセスすることができる。

Go 言語では、`mem` という名称で、WebAssembly の線形メモリーがエクスポートされている。

Go 言語の main 関数を呼び出す場合、JavaScript 側からコマンドライン引数や環境変数を指定することができるが、そのデータの受け渡しは、`Instance.exports.mem.buffer` をラップした `Uint8Array` を使用して、WebAssembly のメモリーバッファーに直接書き込むことで実現している。

## wasm_exec.js

`wasm_exec.js` では、WebAssembly と JavaScript の連携に必要な機能が実装されている。

Go から呼び出される一部のシステムコールは、このスクリプト内で定義されている。

`Go` JavaScript オブジェクトでは、WebAssembly と JavaScript API を連携するためのメソッド等が実装されているので、Go でビルドされた WebAssembly アプリケーションを実行するのにこのオブジェクトは必須である。

## 参考リンク

- [Go Wiki](https://github.com/golang/go/wiki/WebAssembly)
- [syscall/js](https://golang.org/pkg/syscall/js/)
- [WebAssembly](https://developer.mozilla.org/ja/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly)
- [WebAssembly.instantiateStreaming()](https://developer.mozilla.org/ja/docs/Web/JavaScript/Reference/Global_Objects/WebAssembly/instantiateStreaming)
- [WebAssembly Examples](https://github.com/mdn/webassembly-examples)
- [WebAssembly テキストフォーマットを理解する](https://developer.mozilla.org/ja/docs/WebAssembly/Understanding_the_text_format)

