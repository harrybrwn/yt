package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var (
	// Logging is a variable serves as a toggle for builtin logging
	Logging = false
	client  = &http.Client{}
)

const (
	badchars = `\/:*?"<>|.`
	// agent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36"
	agent = "video download cli tool"
	// agent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36"
)

// ToJSON converts any object to a json string
func ToJSON(obj interface{}) string {
	jsondata, err := json.MarshalIndent(obj, "  ", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsondata)
}

// Pprint prints an object as a json object
func Pprint(obj interface{}) {
	fmt.Println(ToJSON(obj))
}

func get(url string) ([]byte, error) {
	_log(callGraph())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", fmt.Sprintf("%s%d", agent, time.Now().Nanosecond()))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_log("get req recieved")
	data, err := ioutil.ReadAll(resp.Body)
	return data, err
}

func safeFileName(name string) string {
	for i := range badchars {
		if strings.Contains(name, string(badchars[i])) {
			name = strings.Replace(name, string(badchars[i]), "", -1)
		}
	}
	return name
}

func functionCall(level int) string {
	fpcs := make([]uintptr, 1)
	if runtime.Callers(level, fpcs) == 0 {
		return "n/a"
	}
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "n/a"
	}
	return fun.Name()
}

func _log(msg ...interface{}) {
	if !Logging {
		return
	}
	fname := functionCall(3) // 3 is one level above _log()
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("[youtube log] %s:%d %s()\n    ", file, line, fname)
	fmt.Println(msg...)
}

func callGraph() string {
	var funcs [3]string
	for i := 0; i < 3; i++ {
		funcs[i] = functionCall(i + 3)
	}
	return fmt.Sprintf("%s()\n      %s()\n        %s()", funcs[2], funcs[1], funcs[0])
}
