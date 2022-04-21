package sourcecode

import "net/http"

func test() {
	http.ListenAndServe(":8080", nil)
}
