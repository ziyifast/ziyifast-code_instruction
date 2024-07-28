package main

import "net/http"
import _ "net/http/pprof"

func main() {
	//引入pprof
	http.ListenAndServe(":8080", nil)
	select {}
}
