package converter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/linketech/microg/v4"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
	// "strings"
)

// Config 系统配置
type Config struct {
	// 服务器配置
	Server ServerConfig `json:"server"`

	// 日志配置
	Logging LoggingConfig `json:"logging"`

	// 搜索配置
	Search SearchConfig `json:"search"`

	// 安全配置
	Security SecurityConfig `json:"security"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `json:"port"`
	Host         string `json:"host"`
	StaticDir    string `json:"static_dir"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level       string   `json:"level"`
	LogDir      string   `json:"log_dir"`
	AllowedExts []string `json:"allowed_extensions"`
	MaxFileSize int64    `json:"max_file_size"`
	MaxResults  int      `json:"max_results"`
}

// SearchConfig 搜索配置
type SearchConfig struct {
	Timeout           int  `json:"timeout_seconds"`
	MaxConcurrent     int  `json:"max_concurrent"`
	PreferSystemTools bool `json:"prefer_system_tools"`
	CacheResults      bool `json:"cache_results"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	AllowedPaths    []string `json:"allowed_paths"`
	BlockedPatterns []string `json:"blocked_patterns"`
	EnableAuth      bool     `json:"enable_auth"`
	APIKey          string   `json:"api_key,omitempty"`
}

// IsAllowedFile 检查文件是否被允许访问
func (c *Config) IsAllowedFile(filename string) bool {
	// 检查文件扩展名
	ext := filepath.Ext(filename)
	allowed := false
	for _, allowedExt := range c.Logging.AllowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return false
	}

	// 检查路径安全性
	for _, blocked := range c.Security.BlockedPatterns {
		if strings.Contains(filename, blocked) {
			return false
		}
	}

	return true
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         8080,
			Host:         "localhost",
			StaticDir:    "static",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Logging: LoggingConfig{
			Level:       "info",
			LogDir:      ".",
			AllowedExts: []string{".log", ".gz", ".txt"},
			MaxFileSize: 200 * 1024 * 1024, // 200MB
			MaxResults:  10000,
		},
		Search: SearchConfig{
			Timeout:           30,
			MaxConcurrent:     5,
			PreferSystemTools: true,
			CacheResults:      false,
		},
		Security: SecurityConfig{
			AllowedPaths:    []string{"."},
			BlockedPatterns: []string{"../", "..\\"}, // 防止路径遍历
			EnableAuth:      false,
		},
	}
}

var setupConfig = DefaultConfig()

// 全局文件映射表，存储文件ID到完整路径的映射
var (
	filePathMap  = make(map[string]string)
	fileMapMutex sync.RWMutex
)

// 处理工具检测API
func handleToolsCheck(toolChecker *ToolChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		toolChecker.CheckSystemTools()
		c.JSON(http.StatusOK, toolChecker)
	}
}

// 处理日志文件列表API
func handleLogsList(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 microg.LogFilePath 所在目录
		if microg.LogFilePath == "" {
			microg.E("LogFilePath is empty")
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"msg": "LogFilePath is empty"})
			return
		}
		logDir := filepath.Dir(microg.LogFilePath)

		files, err := getLogFiles(logDir, config)
		if err != nil {
			microg.E("getLogFiles %v", err)
			c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{"msg": "Failed to get log files"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"files": files})
	}
}

// 处理日志搜索API
func handleLogsSearch(processor *MainLogProcessor, config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SearchRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			microg.E("DecodeRequest %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request body"})
			return
		}

		// 根据文件ID获取完整路径
		filePath, exists := getFilePathByID(req.FileID)
		if !exists {
			microg.W("File ID not found: %s", req.FileID)
			c.JSON(http.StatusNotFound, gin.H{"msg": "File not found"})
			return
		}

		// 验证文件路径安全性
		if !config.IsAllowedFile(filepath.Base(filePath)) {
			microg.W("Blocked file access attempt: %s", req.FileID)
			c.JSON(http.StatusForbidden, gin.H{"msg": "File access denied"})
			return
		}

		results, err := processor.SearchInFile(filePath, req.Pattern)
		if err != nil {
			microg.W("SearchInFile %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "Internal server error"})
			return
		}

		// 限制结果数量
		if len(results) > config.Logging.MaxResults {
			results = results[:config.Logging.MaxResults]
			microg.W("Search results truncated to %d items", config.Logging.MaxResults)
		}
		c.JSON(http.StatusOK, gin.H{
			"results":   results,
			"count":     len(results),
			"truncated": len(results) == config.Logging.MaxResults,
		})
	}
}

// 处理SSE搜索API
func handleLogsSearchStream(processor *MainLogProcessor, config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		// 设置SSE头部
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")
		w.Header().Set("X-Accel-Buffering", "no") // Disable nginx buffering
		w.Header().Set("Transfer-Encoding", "chunked")
		flusher, ok := w.(http.Flusher)
		if !ok {
			microg.E("ResponseWriter does not support flushing")
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "ResponseWriter does not support flushing"})
			return
		}
		var req SearchRequest
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			microg.E("DecodeRequest %v", err)
			_, err := fmt.Fprintf(w, "event: error\ndata: {\"error\": \"Invalid request body\"}\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
			return
		}

		// 根据文件ID获取完整路径
		filePath, exists := getFilePathByID(req.FileID)
		if !exists {
			microg.W("File ID not found: %s", req.FileID)
			_, err := fmt.Fprintf(w, "event: error\ndata: {\"error\": \"File not found\"}\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
			return
		}

		// 验证文件路径安全性
		if !config.IsAllowedFile(filepath.Base(filePath)) {
			microg.W("Blocked file access attempt: %s", req.FileID)
			_, err := fmt.Fprintf(w, "event: error\ndata: {\"error\": \"File access denied\"}\n\n")
			if err != nil {
				return
			}
			flusher.Flush()
			return
		}

		// 创建可取消的上下文
		ctx, cancel := context.WithCancel(c.Request.Context())
		defer cancel()

		// 监听客户端断开连接
		go func() {
			<-c.Request.Context().Done()
			cancel()
		}()

		// 发送开始事件
		_, err := fmt.Fprintf(w, "event: start\ndata: {\"message\": \"开始搜索...\"}\n\n")
		if err != nil {
			return
		}
		flusher.Flush()

		// 执行流式搜索
		err = processor.SearchInFileStream(ctx, filePath, req.Pattern, func(result SearchResult) {
			// 检查上下文是否已取消
			select {
			case <-ctx.Done():
				return
			default:
			}

			data, _ := result.ToBytes()
			fmt.Printf("event: result\ndata: %s\n\n", data)
			_, err := fmt.Fprintf(w, "event: result\ndata: %s\n\n", data)
			if err != nil {
				return
			}
			flusher.Flush()
		})

		if err != nil {
			// 检查是否是上下文取消导致的错误
			if errors.Is(err, context.Canceled) {
				microg.I("Search cancelled by client disconnect")
				return
			} // 如果命令退出状态是 1 说明是没有数据不用返回数据，也不用报错
			var exitError *exec.ExitError
			if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
				return
			}
			microg.E("SearchInFileStream %v", err)
			_, err := fmt.Fprintf(w, "event: error\ndata: {\"error\": \"%s\"}\n\n", err.Error())
			if err != nil {
				return
			}
		} else {
			_, err := fmt.Fprintf(w, "event: finished\ndata: {\"message\": \"搜索完成\"}\n\n")
			if err != nil {
				return
			}
		}
		flusher.Flush()
	}
}

// 获取日志文件列表
func getLogFiles(dir string, config *Config) ([]LogFile, error) {
	var files []LogFile

	// 清空并重新构建文件映射表
	fileMapMutex.Lock()
	filePathMap = make(map[string]string)
	fileMapMutex.Unlock()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 检查文件扩展名
			ext := filepath.Ext(path)
			allowed := false
			allowed = slices.Contains(config.Logging.AllowedExts, ext)

			// 检查文件大小
			if allowed && info.Size() <= config.Logging.MaxFileSize {
				// 生成相对于日志目录的相对路径作为ID
				relPath, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}

				// 存储ID到完整路径的映射
				fileMapMutex.Lock()
				filePathMap[relPath] = path
				fileMapMutex.Unlock()

				files = append(files, LogFile{
					Name:       info.Name(),
					ID:         relPath, // 使用相对路径作为ID
					Size:       info.Size(),
					Modified:   info.ModTime(),
					Compressed: filepath.Ext(path) == ".gz",
					path:       path, // 完整路径仅供后端使用
				})
			}
		}

		return nil
	})

	return files, err
}

// SearchRequest 请求结构体
type SearchRequest struct {
	FileID  string `json:"file_id"` // 改为使用文件ID
	Pattern string `json:"pattern"`
}

// getFilePathByID 根据文件ID获取完整路径
func getFilePathByID(fileID string) (string, bool) {
	fileMapMutex.RLock()
	defer fileMapMutex.RUnlock()
	path, exists := filePathMap[fileID]
	return path, exists
}

// LogFile 日志文件结构体
type LogFile struct {
	Name       string    `json:"name"`
	ID         string    `json:"id"` // 文件ID，用于前端识别
	Size       int64     `json:"size"`
	Modified   time.Time `json:"modified"`
	Compressed bool      `json:"compressed"`
	path       string    // 完整路径，不暴露给前端
}

// GetPath 获取文件的完整路径（仅供后端使用）
func (lf *LogFile) GetPath() string {
	return lf.path
}
