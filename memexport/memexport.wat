(module
  (import "js" "hello" (func $hello (param i32 i32)))
  (memory 1)
  (data (i32.const 0) "Hello, WebAssembly")
  (func (export "hello")
	i32.const 0
	i32.const 18
	call $hello)
  (export "mem" (memory 0)))
