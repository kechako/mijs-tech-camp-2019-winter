(module
  (import "js" "hello" (func $hello (param i32 i32)))
  (memory (import "js" "mem") 1)
  (data (i32.const 0) "Hello, WebAssembly")
  (func (export "hello")
	i32.const 0
	i32.const 18
	call $hello))
