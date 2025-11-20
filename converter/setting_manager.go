package converter

/*
 * @Author: SimingLiu siming.liu@dbjtech.com
 * @Date: 2025-01-30 12:00:00
 * @LastEditors: SimingLiu siming.liu@dbjtech.com
 * @LastEditTime: 2025-01-30 12:00:00
 * @FilePath: \go_809_converter\converter\setting_manager.go
 * @Description: 配置管理API处理
 *
 */

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/dbjtech/go_809_converter/libs"
	"github.com/gin-gonic/gin"
	"github.com/gookit/config/v2"
	"github.com/linketech/microg/v4"
)

type NodeFor809 struct {
	IA1                 int64    `json:"IA1"`
	IC1                 int64    `json:"IC1"`
	M1                  int64    `json:"M1"`
	CryptoPacket        []string `json:"cryptoPacket"`
	Enabled             bool     `json:"enabled"`
	EncryptKey          int64    `json:"encryptKey"`
	ExtendVersion       bool     `json:"extendVersion"`
	GovServerIP         string   `json:"govServerIP"`
	GovServerPort       int64    `json:"govServerPort"`
	LocalServerIP       string   `json:"localServerIP"`
	LocalServerPort     int64    `json:"localServerPort"`
	Name                string   `json:"name"`
	OpenCrypto          bool     `json:"openCrypto"`
	PlatformId          int64    `json:"platformId"`
	PlatformPassword    string   `json:"platformPassword"`
	PlatformUserId      int64    `json:"platformUserId"`
	ProtocolVersion     string   `json:"protocolVersion"`
	ThirdpartPort       int64    `json:"thirdpartPort"`
	UseLocationInterval bool     `json:"useLocationInterval"`
}

// SettingRequest 配置请求结构
type SettingRequest struct {
	Config    map[string]any `json:"config,omitempty"`
	Key       string         `json:"key,omitempty"`
	Value     any            `json:"value,omitempty"`
	Operation string         `json:"operation"`
	Timestamp string         `json:"timestamp,omitempty"`
}

// SettingResponse 配置响应结构
type SettingResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message,omitempty"`
	Config  map[string]any `json:"config,omitempty"`
	History []HistoryItem  `json:"history,omitempty"`
}
type EnvEntity struct {
	ConsolePort int64                 `json:"consolePort"`
	Converter   map[string]NodeFor809 `json:"converter"`
}

type MysqlDbEntity struct {
	Database      string `json:"database"`
	Host          string `json:"host"`
	Password      string `json:"password"`
	PoolIdleConns int64  `json:"pool_idle_conns"`
	PoolSize      int64  `json:"pool_size"`
	Port          int64  `json:"port"`
	ShowSQL       bool   `json:"showSQL"`
	User          string `json:"user"`
}

// HistoryItem 历史记录项
type HistoryItem struct {
	Timestamp  string         `json:"timestamp"`
	Operation  map[string]any `json:"operation"`
	FullConfig map[string]any `json:"full_config"`
}

func currentConfigPath() string {
	return libs.GetConfigPath(libs.ConfigType, libs.ConfigFile)
}

func historyFilePath() string {
	return filepath.Join(filepath.Dir(currentConfigPath()), "config.history")
}

func backupFilePath() string {
	return currentConfigPath() + ".backup"
}

func templateFilePath() string {
	return filepath.Join(filepath.Dir(currentConfigPath()), "configuration.toml.template")
}

// getCurrentConfig 获取当前配置
func getCurrentConfig(c *gin.Context) {
	configData := config.Data()

	response := SettingResponse{
		Success: true,
		Config:  configData,
	}

	c.JSON(http.StatusOK, response)
}

// saveConfig 保存配置
func saveConfig(c *gin.Context) {
	var req SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SettingResponse{
			Success: false,
			Message: "请求格式错误: " + err.Error(),
		})
		return
	}

	// 备份当前配置
	if err := backupCurrentConfig(); err != nil {
		microg.W("Failed to backup config: %v", err)
	}

	// 记录历史
	historyData := map[string]any{
		"operation": req.Operation,
		"config":    req.Config,
	}
	if err := recordHistory(historyData, config.Data()); err != nil {
		microg.W("Failed to record history: %v", err)
	}

	if strings.ToLower(strings.TrimSpace(req.Operation)) == "add_subproject" {
		req.Config = ensureAddSubProjectDefaults(req.Config)
	}
	flatReq := flattenConfig(req.Config, "")
	for key, value := range flatReq {
		if err := config.Set(key, value, true); err != nil {
			c.JSON(http.StatusInternalServerError, SettingResponse{
				Success: false,
				Message: fmt.Sprintf("设置配置项 %s 失败: %v", key, err),
			})
			return
		}
	}

	// 保存配置到文件
	if err := saveConfigToFile(); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "保存配置文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		Message: "配置保存成功",
	})
}

// addConfig 添加配置项
func addConfig(c *gin.Context) {
	var req SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SettingResponse{
			Success: false,
			Message: "请求格式错误: " + err.Error(),
		})
		return
	}

	// 检查配置项是否已存在
	if config.Exists(req.Key) {
		c.JSON(http.StatusConflict, SettingResponse{
			Success: false,
			Message: fmt.Sprintf("配置项 %s 已存在", req.Key),
		})
		return
	}

	// 备份当前配置
	if err := backupCurrentConfig(); err != nil {
		microg.W("Failed to backup config: %v", err)
	}

	// 记录历史
	historyData := map[string]any{
		"operation": req.Operation,
		"config": map[string]any{
			req.Key: req.Value,
		},
	}
	if err := recordHistory(historyData, config.Data()); err != nil {
		microg.W("Failed to record history: %v", err)
	}

	// 添加配置项
	if err := config.Set(req.Key, req.Value, true); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: fmt.Sprintf("添加配置项失败: %v", err),
		})
		return
	}

	// 保存配置到文件
	if err := saveConfigToFile(); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "保存配置文件失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		Message: fmt.Sprintf("配置项 %s 添加成功", req.Key),
	})
}

// deleteConfig 删除配置项
func deleteConfig(c *gin.Context) {
	var req SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SettingResponse{
			Success: false,
			Message: "请求格式错误: " + err.Error(),
		})
		return
	}

	// 检查配置项是否存在
	if !config.Exists(req.Key) {
		c.JSON(http.StatusNotFound, SettingResponse{
			Success: false,
			Message: fmt.Sprintf("配置项 %s 不存在", req.Key),
		})
		return
	}

	// 备份当前配置
	if err := backupCurrentConfig(); err != nil {
		microg.W("Failed to backup config: %v", err)
	}

	// 记录历史
	historyData := map[string]any{
		"operation": req.Operation,
		"key":       req.Key,
		"old_value": config.Get(req.Key),
	}
	if err := recordHistory(historyData, config.Data()); err != nil {
		microg.W("Failed to record history: %v", err)
	}

	// 删除配置项（通过重新构建配置来实现）
	if err := removeConfigKey(req.Key); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: fmt.Sprintf("删除配置项失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		Message: fmt.Sprintf("配置项 %s 删除成功", req.Key),
	})
}

// resetConfig 重置配置
func resetConfig(c *gin.Context) {
	// 备份当前配置
	if err := backupCurrentConfig(); err != nil {
		microg.W("Failed to backup config: %v", err)
	}

	// 记录历史
	historyData := map[string]any{
		"operation": "reset",
		"message":   "重置为默认配置",
	}
	if err := recordHistory(historyData, config.Data()); err != nil {
		microg.W("Failed to record history: %v", err)
	}

	// 重新加载默认配置
	if err := resetToDefaultConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "重置配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		Message: "配置已重置为默认值",
	})
}

// getHistory 获取历史记录
func getHistory(c *gin.Context) {
	history, err := loadHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "加载历史记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		History: history,
	})
}

// rollbackConfig 回滚配置
func rollbackConfig(c *gin.Context) {
	var req SettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SettingResponse{
			Success: false,
			Message: "请求格式错误: " + err.Error(),
		})
		return
	}

	// 查找指定时间戳的历史记录
	history, err := loadHistory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "加载历史记录失败: " + err.Error(),
		})
		return
	}

	var targetHistory *HistoryItem
	for _, item := range history {
		if item.Timestamp == req.Timestamp {
			targetHistory = &item
			break
		}
	}

	if targetHistory == nil {
		c.JSON(http.StatusNotFound, SettingResponse{
			Success: false,
			Message: "未找到指定时间戳的历史记录",
		})
		return
	}

	// 备份当前配置
	if err := backupCurrentConfig(); err != nil {
		microg.W("Failed to backup config: %v", err)
	}

	// 执行回滚到目标快照
	if targetHistory.FullConfig == nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "历史记录缺少完整配置快照",
		})
		return
	}

	// 清空并写入目标配置
	config.ClearAll()
	flat := flattenConfig(targetHistory.FullConfig, "")
	for k, v := range flat {
		if err := config.Set(k, v, true); err != nil {
			c.JSON(http.StatusInternalServerError, SettingResponse{
				Success: false,
				Message: fmt.Sprintf("回滚设置项 %s 失败: %v", k, err),
			})
			return
		}
	}

	if err := saveConfigToFile(); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "保存回滚后的配置失败: " + err.Error(),
		})
		return
	}

	// 记录回滚后的完整配置快照
	op := map[string]any{
		"operation":        "rollback",
		"target_timestamp": req.Timestamp,
	}
	if err := recordHistory(op, config.Data()); err != nil {
		microg.W("Failed to record rollback snapshot: %v", err)
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		Message: fmt.Sprintf("已回滚到 %s 的配置", req.Timestamp),
	})
}

// clearHistory 清空历史记录
func clearHistory(c *gin.Context) {
	if err := os.Remove(historyFilePath()); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, SettingResponse{
			Success: false,
			Message: "清空历史记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SettingResponse{
		Success: true,
		Message: "历史记录已清空",
	})
}

// 辅助函数

// flattenConfig 将嵌套配置扁平化
func flattenConfig(data any, prefix string) map[string]any {
	result := make(map[string]any)

	switch v := data.(type) {
	case map[string]any:
		for key, value := range v {
			newKey := key
			if prefix != "" {
				newKey = prefix + "." + key
			}

			if nested, ok := value.(map[string]any); ok {
				maps.Copy(result, flattenConfig(nested, newKey))
			} else {
				result[newKey] = value
			}
		}
	default:
		if prefix != "" {
			result[prefix] = v
		}
	}

	return result
}

// backupCurrentConfig 备份当前配置
func backupCurrentConfig() error {
	configFile := currentConfigPath()
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return nil // 配置文件不存在，无需备份
	}

	input, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return os.WriteFile(backupFilePath(), input, 0644)
}

// saveConfigToFile 保存配置到文件
func saveConfigToFile() error {
	configFile := currentConfigPath()
	configData := config.Data()
	// 对 configData 进行深度复制，避免修改原始数据
	configDataBackup := make(map[string]any)
	maps.Copy(configDataBackup, configData)
	// 删除 env jtwTcp normalTcp字段
	delete(configDataBackup, "env")
	delete(configDataBackup, "jtwTcp")
	delete(configDataBackup, "normalTcp")
	// 构建TOML格式的配置内容
	content := buildTOMLContent(configDataBackup)

	return os.WriteFile(configFile, []byte(content), 0644)
}

// buildTOMLContent 构建TOML格式内容
func buildTOMLContent(data map[string]any) string {
	var builder strings.Builder

	// 递归输出：先输出当前段的标量键，再深入子段
	var writeSection func(prefix string, obj map[string]any)
	writeSection = func(prefix string, obj map[string]any) {
		scalars := make(map[string]any)
		nested := make(map[string]map[string]any)

		for k, v := range obj {
			if m, ok := v.(map[string]any); ok {
				nested[k] = m
			} else {
				scalars[k] = v
			}
		}

		if len(scalars) > 0 {
			builder.WriteString(fmt.Sprintf("[%s]\n", prefix))
			keys := make([]string, 0, len(scalars))
			for k := range scalars {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				builder.WriteString(fmt.Sprintf("%s = %s\n", k, formatTOMLValue(scalars[k])))
			}
			builder.WriteString("\n")
		}

		nkeys := make([]string, 0, len(nested))
		for k := range nested {
			nkeys = append(nkeys, k)
		}
		sort.Strings(nkeys)
		for _, k := range nkeys {
			writeSection(prefix+"."+k, nested[k])
		}
	}

	// 顶级遍历：对map类型写为段，顶级标量直接写键值
	topKeys := make([]string, 0, len(data))
	for k := range data {
		topKeys = append(topKeys, k)
	}
	sort.Strings(topKeys)

	for _, k := range topKeys {
		if m, ok := data[k].(map[string]any); ok {
			writeSection(k, m)
		} else {
			builder.WriteString(fmt.Sprintf("%s = %s\n", k, formatTOMLValue(data[k])))
		}
	}

	return builder.String()
}

// formatTOMLValue 格式化TOML值
func formatTOMLValue(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case bool:
		return strconv.FormatBool(v)
	case int, int32, int64:
		return fmt.Sprintf("%d", v)
	case float32:
		return formatFloatNoSci(float64(v))
	case float64:
		return formatFloatNoSci(v)
	case json.Number:
		s := v.String()
		if strings.ContainsAny(s, "eE") {
			if f, err := v.Float64(); err == nil {
				return formatFloatNoSci(f)
			}
		}
		return s
	case []any:
		var items []string
		for _, item := range v {
			items = append(items, formatTOMLValue(item))
		}
		return fmt.Sprintf("[%s]", strings.Join(items, ", "))
	case map[string]any:
		// 嵌套的map不应该在这里处理，它们应该作为独立的段
		// 如果到了这里，说明有配置结构问题，我们跳过它
		return `""`
	default:
		return fmt.Sprintf(`"%v"`, v)
	}
}

func formatFloatNoSci(f float64) string {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return "0"
	}
	if math.Abs(f-math.Round(f)) < 1e-9 {
		return strconv.FormatInt(int64(math.Round(f)), 10)
	}
	s := strconv.FormatFloat(f, 'f', 6, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	if s == "" {
		return "0"
	}
	return s
}

// removeConfigKey 删除配置键
func removeConfigKey(key string) error {
	configData := config.Data()

	// 从配置数据中删除指定键
	parts := strings.Split(key, ".")
	if len(parts) == 1 {
		delete(configData, key)
	} else {
		// 处理嵌套键
		current := configData
		for i, part := range parts[:len(parts)-1] {
			if next, ok := current[part].(map[string]any); ok {
				current = next
			} else {
				return fmt.Errorf("配置路径 %s 不存在", strings.Join(parts[:i+1], "."))
			}
		}
		delete(current, parts[len(parts)-1])
	}

	// 重新构建配置
	config.ClearAll()
	for k, v := range configData {
		config.Set(k, v, true)
	}

	return saveConfigToFile()
}

// resetToDefaultConfig 重置为默认配置
func resetToDefaultConfig() error {
	templateFile := templateFilePath()
	configFile := currentConfigPath()

	if _, err := os.Stat(templateFile); os.IsNotExist(err) {
		return fmt.Errorf("默认配置模板文件不存在: %s", templateFile)
	}

	input, err := os.ReadFile(templateFile)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, input, 0644); err != nil {
		return err
	}

	// 重新加载配置
	config.ClearAll()
	return config.LoadFiles(configFile)
}

// recordHistory 记录历史
func recordHistory(operation map[string]interface{}, fullConfig map[string]interface{}) error {
	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(historyFilePath()), 0755); err != nil {
		return err
	}

	// 构建历史记录项
	historyItem := HistoryItem{
		Timestamp:  time.Now().Format("2006-01-02 15:04:05"),
		Operation:  operation,
		FullConfig: fullConfig,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(historyItem)
	if err != nil {
		return err
	}

	// 追加到历史文件
	file, err := os.OpenFile(historyFilePath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(string(jsonData) + "\n")
	return err
}

// loadHistory 加载历史记录
func loadHistory() ([]HistoryItem, error) {
	var history []HistoryItem

	if _, err := os.Stat(historyFilePath()); os.IsNotExist(err) {
		return history, nil // 文件不存在，返回空历史
	}

	file, err := os.Open(historyFilePath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	validLines := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var item HistoryItem
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			microg.W("Failed to parse history line: %s, error: %v", line, err)
			continue
		}

		history = append(history, item)
		validLines++
		if validLines >= 1000 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// 按时间戳倒序排列（最新的在前面）
	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp > history[j].Timestamp
	})

	return history, nil
}

// deleteHistoryItem 删除某条历史记录
func deleteHistoryItem(c *gin.Context) {
	var req struct {
		Timestamp string `json:"timestamp"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Timestamp) == "" {
		c.JSON(http.StatusBadRequest, SettingResponse{Success: false, Message: "请求格式错误或缺少时间戳"})
		return
	}

	// 读取历史文件并过滤
	var lines []string
	if _, err := os.Stat(historyFilePath()); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, SettingResponse{Success: false, Message: "历史记录文件不存在"})
		return
	}
	data, err := os.ReadFile(historyFilePath())
	if err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{Success: false, Message: "读取历史记录失败"})
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		var item HistoryItem
		if err := json.Unmarshal([]byte(l), &item); err != nil {
			continue
		}
		if item.Timestamp != req.Timestamp {
			lines = append(lines, l)
		}
	}
	// 重写文件
	if err := os.WriteFile(historyFilePath(), []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		c.JSON(http.StatusInternalServerError, SettingResponse{Success: false, Message: "删除历史记录失败"})
		return
	}
	c.JSON(http.StatusOK, SettingResponse{Success: true, Message: "历史记录已删除"})
}
func ensureAddSubProjectDefaults(cfg map[string]any) map[string]any {
	m := cfg
	for envKey, envVal := range m {
		envMap, ok := envVal.(map[string]any)
		if !ok {
			continue
		}
		convVal, ok := envMap["converter"]
		if !ok {
			continue
		}
		convMap, ok := convVal.(map[string]any)
		if !ok {
			continue
		}
		for nodeKey, nodeVal := range convMap {
			nodeMap, ok := nodeVal.(map[string]any)
			if !ok {
				continue
			}
			defaults := map[string]any{
				"cryptoPacket":                []string{},
				"extendVersion":               true,
				"jtw809ConverterDownLinkIp":   "127.0.0.1",
				"jtw809ConverterDownLinkPort": int64(1302),
				"jtw809ConverterIp":           "127.0.0.1",
				"jtw809ConverterPort":         int64(1311),
				"thirdpartPort":               int64(11223),
				"useLocationInterval":         false,
				"IA1":                         "20000000",
				"IC1":                         "30000000",
				"M1":                          "10000000",
			}
			for k, v := range defaults {
				if _, exists := nodeMap[k]; !exists {
					nodeMap[k] = v
				}
			}
			convMap[nodeKey] = nodeMap
		}
		envMap["converter"] = convMap
		m[envKey] = envMap
	}
	return m
}
