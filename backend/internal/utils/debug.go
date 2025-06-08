package utils

import (
	"encoding/json"
	"log"
	"runtime"
)

// DebugObject prints a formatted JSON representation of any object
func DebugObject(prefix string, obj interface{}) {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Printf("Error marshaling object for debug: %v", err)
		return
	}

	// Get caller information
	_, file, line, ok := runtime.Caller(1)
	if ok {
		log.Printf("[DEBUG] %s:%d - %s: %s", file, line, prefix, string(data))
	} else {
		log.Printf("[DEBUG] %s: %s", prefix, string(data))
	}
}

// TraceFunction logs entry and exit of functions
func TraceFunction() func() {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return func() {}
	}

	fn := runtime.FuncForPC(pc)
	funcName := fn.Name()
	log.Printf("[TRACE] Entering %s (%s:%d)", funcName, file, line)

	return func() {
		log.Printf("[TRACE] Exiting %s", funcName)
	}
}
