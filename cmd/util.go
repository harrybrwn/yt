package cmd

// func functionCall(level int) string {
// 	fpcs := make([]uintptr, 1)
// 	if runtime.Callers(level, fpcs) == 0 {
// 		fmt.Println("logging broke! -", "'runtime.Callers(2, fpcs) == 0'")
// 	}
// 	fun := runtime.FuncForPC(fpcs[0] - 1)
// 	if fun == nil {
// 		return "n/a"
// 	}
// 	return fun.Name()
// }

// func log(msg ...interface{}) {
// 	if logging {
// 		fname := functionCall(3) // 3 is one level above _log()
// 		_, file, line, _ := runtime.Caller(1)
// 		fmt.Printf("[yt log] %s:%d %s()\n    ", file, line, fname)
// 		fmt.Println(msg...)
// 	}
// }
