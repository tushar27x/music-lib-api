package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// DebugInfo holds debug information
type DebugInfo struct {
	Timestamp    time.Time        `json:"timestamp"`
	GoVersion    string           `json:"go_version"`
	Architecture string           `json:"architecture"`
	NumCPU       int              `json:"num_cpu"`
	NumGoroutine int              `json:"num_goroutine"`
	MemoryStats  runtime.MemStats `json:"memory_stats"`
}

// GetDebugInfo returns current debug information
func GetDebugInfo() DebugInfo {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return DebugInfo{
		Timestamp:    time.Now(),
		GoVersion:    runtime.Version(),
		Architecture: fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemoryStats:  memStats,
	}
}

// LogDebugInfo logs debug information
func LogDebugInfo() {
	info := GetDebugInfo()
	log.Printf("🐛 Debug Info - Go: %s, Arch: %s, CPU: %d, Goroutines: %d",
		info.GoVersion, info.Architecture, info.NumCPU, info.NumGoroutine)
}

// DebugMiddleware adds debug information to Gin context
func DebugMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		log.Printf("🔍 %s %s - %s - %v - %d",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			duration,
			c.Writer.Status())
	}
}

// PrettyPrintJSON prints JSON in a readable format
func PrettyPrintJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("❌ Error marshaling JSON: %v", err)
		return
	}
	log.Printf("📄 JSON Output:\n%s", string(jsonData))
}

// PanicRecovery recovers from panics and logs them
func PanicRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("💥 Panic recovered: %s", err)
		}
		c.AbortWithStatusJSON(500, gin.H{
			"error":   "Internal server error",
			"message": "A panic occurred and was recovered",
		})
	})
}
