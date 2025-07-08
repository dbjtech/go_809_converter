package converter

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/linketech/microg/v4"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// LogProcessor 日志处理器接口
type LogProcessor interface {
	SearchInFile(filename, pattern string) ([]SearchResult, error)
	SearchInFileStream(ctx context.Context, filename, pattern string, callback func(SearchResult)) error
}

// SearchResult 搜索结果
type SearchResult struct {
	LineNumber int    `json:"line_number"`
	Content    string `json:"content"`
	Matched    string `json:"matched"`
}

func hasAroundParam(patterns []string) bool {
	for _, pattern := range patterns {
		// 将pattern按空格分割
		parts := strings.FieldsFunc(pattern, func(r rune) bool {
			return r == ' ' || r == '\t'
		})
		for _, p := range parts {
			if p == "-A" || p == "-B" {
				return true
			}
		}
	}
	return false
}

// filterGrepParams 过滤pattern中的grep参数，只保留实际的搜索模式
func filterGrepParams(pattern string) string {
	// 将pattern按空格分割
	parts := strings.Fields(pattern)
	if len(parts) == 0 {
		return pattern
	}

	var filteredParts []string
	i := 0

	for i < len(parts) {
		part := parts[i]

		// 跳过单字符选项参数
		if part == "-i" || part == "-w" || part == "-x" || part == "-F" || part == "-E" || part == "-n" {
			i++
			continue
		}

		// 跳过带数值的参数（-m N, -A N, -B N, -C N）
		if (part == "-m" || part == "-A" || part == "-B" || part == "-C") && i+1 < len(parts) {
			i += 2 // 跳过参数和其值
			continue
		}

		// 跳过-e参数及其模式
		if part == "-e" && i+1 < len(parts) {
			i += 2 // 跳过-e和其模式
			continue
		}

		// 跳过以-开头的其他参数
		if strings.HasPrefix(part, "-") {
			i++
			continue
		}

		// 保留非参数部分（实际的搜索模式）
		filteredParts = append(filteredParts, part)
		i++
	}

	// 如果没有找到任何模式，返回原始pattern
	if len(filteredParts) == 0 {
		return pattern
	}

	// 返回过滤后的模式，用空格连接
	return strings.Join(filteredParts, " ")
}

func (sr SearchResult) ToBytes() ([]byte, error) {
	result := fmt.Sprintf("line_number:%v&nl matched:%v&nl content:%v}", sr.LineNumber, sr.Matched, sr.Content)
	return []byte(result), nil
}

// MainLogProcessor 主日志处理器
type MainLogProcessor struct {
	toolChecker     *ToolChecker
	systemProcessor *SystemToolProcessor
	nativeProcessor *NativeGoProcessor
}

// NewLogProcessor 创建新的日志处理器
func NewLogProcessor(toolChecker *ToolChecker) *MainLogProcessor {
	return &MainLogProcessor{
		toolChecker:     toolChecker,
		systemProcessor: NewSystemToolProcessor(),
		nativeProcessor: NewNativeGoProcessor(),
	}
}

// SearchInFile 在文件中搜索
func (mlp *MainLogProcessor) SearchInFile(filename, pattern string) ([]SearchResult, error) {
	// 验证文件路径安全性
	if err := validateFilePath(filename); err != nil {
		return nil, err
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件不存在: %s", filename)
	}

	// 判断是否为压缩文件
	isCompressed := strings.HasSuffix(strings.ToLower(filename), ".gz")

	// 根据工具可用性选择处理策略
	if mlp.canUseSystemTools(isCompressed) {
		return mlp.systemProcessor.SearchInFile(filename, pattern)
	}

	return mlp.nativeProcessor.SearchInFile(filename, pattern)
}

// SearchInFileStream 流式搜索文件
func (mlp *MainLogProcessor) SearchInFileStream(ctx context.Context, filename, pattern string, callback func(SearchResult)) error {
	// 验证文件路径安全性
	if err := validateFilePath(filename); err != nil {
		return err
	}

	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filename)
	}

	// 判断是否为压缩文件
	isCompressed := strings.HasSuffix(strings.ToLower(filename), ".gz")

	// 根据工具可用性选择处理策略
	if mlp.canUseSystemTools(isCompressed) {
		return mlp.systemProcessor.SearchInFileStream(ctx, filename, pattern, callback)
	}

	return mlp.nativeProcessor.SearchInFileStream(ctx, filename, pattern, callback)
}

// canUseSystemTools 判断是否可以使用系统工具
func (mlp *MainLogProcessor) canUseSystemTools(isCompressed bool) bool {
	if isCompressed {
		return mlp.toolChecker.HasGzip && mlp.toolChecker.HasGrep
	}
	return mlp.toolChecker.HasCat && mlp.toolChecker.HasGrep
}

// SystemToolProcessor 系统工具处理器
type SystemToolProcessor struct{}

// NewSystemToolProcessor 创建系统工具处理器
func NewSystemToolProcessor() *SystemToolProcessor {
	return &SystemToolProcessor{}
}

// SearchInFile 使用系统工具搜索文件
func (stp *SystemToolProcessor) SearchInFile(filename, pattern string) ([]SearchResult, error) {
	isCompressed := strings.HasSuffix(strings.ToLower(filename), ".gz")

	var cmd *exec.Cmd

	// 检测是否有Unix工具可用（GitBash环境）
	hasUnixTools := stp.hasUnixTools()

	// 处理多层查询：检查pattern是否包含" | "分隔符
	patterns := strings.Split(pattern, " | ")
	var grepChain string
	hasAround := hasAroundParam(patterns)
	if len(patterns) > 1 {
		// 多层查询：构建多个grep命令的管道
		if runtime.GOOS == "windows" && !hasUnixTools {
			// Windows环境使用 findstr
			for i, p := range patterns {
				if i == 0 {
					grepChain = fmt.Sprintf(`findstr /n "%s"`, p)
				} else {
					grepChain += fmt.Sprintf(` | findstr "%s"`, p)
				}
			}
		} else {
			// Unix环境使用grep
			for i, p := range patterns {
				if i == 0 {
					grepChain = fmt.Sprintf(`grep -n %s`, p)
				} else {
					grepChain += fmt.Sprintf(` | grep %s`, p)
				}
			}
		}
	} else {
		// 单层查询：使用原有逻辑
		if runtime.GOOS == "windows" && !hasUnixTools {
			grepChain = fmt.Sprintf(`findstr /n "%s"`, pattern)
		} else {
			grepChain = fmt.Sprintf(`grep -n %s`, pattern)
		}
	}

	if isCompressed {
		// 使用 gunzip -d -c | grep链
		if runtime.GOOS == "windows" && !hasUnixTools {
			// 纯Windows环境，使用Windows命令
			cmd = exec.Command("cmd", "/c", fmt.Sprintf(`gzip -dc "%s" | %s`, filename, grepChain))
		} else {
			// Unix环境或GitBash环境
			cmd = exec.Command("sh", "-c", fmt.Sprintf(`gzip -dc "%s" | %s`, filename, grepChain))
		}
	} else {
		// 使用 cat | grep链
		if runtime.GOOS == "windows" && !hasUnixTools {
			// 纯Windows环境，使用Windows命令
			cmd = exec.Command("cmd", "/c", fmt.Sprintf(`type "%s" | %s`, filename, grepChain))
		} else {
			// Unix环境或GitBash环境
			cmd = exec.Command("sh", "-c", fmt.Sprintf(`cat "%s" | %s`, filename, grepChain))
		}
	}

	output, err := cmd.Output()
	if err != nil {
		// 如果命令退出状态是 1 说明是没有数据不用返回数据，也不用报错
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("执行系统命令失败: %v", err)
	}

	return parseGrepOutput(string(output), pattern, hasAround), nil
}

// SearchInFileStream 使用系统工具流式搜索文件
func (stp *SystemToolProcessor) SearchInFileStream(ctx context.Context, filename, pattern string, callback func(SearchResult)) error {
	isCompressed := strings.HasSuffix(strings.ToLower(filename), ".gz")

	var cmd *exec.Cmd

	// 检测是否有Unix工具可用（GitBash环境）
	hasUnixTools := stp.hasUnixTools()

	// 处理多层查询：检查pattern是否包含" | "分隔符
	patterns := strings.Split(pattern, " | ")
	var grepChain string
	hasAround := hasAroundParam(patterns)
	if len(patterns) > 1 {
		// 多层查询：构建多个grep命令的管道
		if runtime.GOOS == "windows" && !hasUnixTools {
			// Windows环境使用findstr
			for i, p := range patterns {
				if i == 0 {
					grepChain = fmt.Sprintf(`findstr /n "%s"`, p)
				} else {
					grepChain += fmt.Sprintf(` | findstr "%s"`, p)
				}
			}
		} else {
			// Unix环境使用grep
			for i, p := range patterns {
				if i == 0 {
					grepChain = fmt.Sprintf(`grep -n %s`, p)
				} else {
					grepChain += fmt.Sprintf(` | grep %s`, p)
				}
			}
		}
		pattern = patterns[len(patterns)-1]
	} else {
		// 单层查询：使用原有逻辑
		if runtime.GOOS == "windows" && !hasUnixTools {
			grepChain = fmt.Sprintf(`findstr /n "%s"`, pattern)
		} else {
			grepChain = fmt.Sprintf(`grep -n %s`, pattern)
		}
	}

	if isCompressed {
		// 使用 gzip -dc | grep链
		if runtime.GOOS == "windows" && !hasUnixTools {
			// 纯Windows环境，使用Windows命令
			cmd = exec.Command("cmd", "/c", fmt.Sprintf(`gzip -dc "%s" | %s`, filename, grepChain))
		} else {
			// Unix环境或GitBash环境
			cmd = exec.Command("sh", "-c", fmt.Sprintf(`gzip -dc "%s" | %s`, filename, grepChain))
		}
	} else {
		// 使用 cat | grep链
		if runtime.GOOS == "windows" && !hasUnixTools {
			// 纯Windows环境，使用Windows命令
			cmd = exec.Command("cmd", "/c", fmt.Sprintf(`type "%s" | %s`, filename, grepChain))
		} else {
			// Unix环境或GitBash环境
			cmd = exec.Command("sh", "-c", fmt.Sprintf(`cat "%s" | %s`, filename, grepChain))
		}
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("创建管道失败: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动命令失败: %v", err)
	}

	// 创建一个goroutine来处理上下文取消
	go func() {
		<-ctx.Done()
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	// pattern 可能包含 -i -w -x -F -E -m N -A N -B N -C N -e PTRN 等 grep 参数，需要把这些参数过滤掉
	filteredPattern := filterGrepParams(pattern)
	scanner := bufio.NewScanner(stdout)
	var result *SearchResult
	for scanner.Scan() {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()
		if line == "" {
			continue
		}
		// 解析grep -n的输出格式: 行号:内容 or 行号-内容
		re := regexp.MustCompile(`^(\d+)([-:])(.*)`) // 匹配行号、分隔符和内容
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 {
			lineNum, _ := strconv.Atoi(matches[1])
			delimiter := matches[2]
			content := matches[3]
			if !hasAround {
				if result == nil {
					result = new(SearchResult)
				}
				result.LineNumber = lineNum
				result.Content = content
				result.Matched = filteredPattern
				callback(*result)
				result = nil
			} else {
				switch delimiter {
				case ":":
					if result != nil {
						callback(*result)
					}
					// 匹配行号:内容
					result = &SearchResult{
						LineNumber: lineNum,
						Content:    content,
						Matched:    filteredPattern,
					}
				case "-":
					// 匹配行号-内容 追加到上一条搜索结果
					result.Content += "\n" + content
				}
			}
		}
	}
	if result != nil {
		callback(*result)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取输出失败: %v", err)
	}

	return cmd.Wait()
}

// hasUnixTools 检测是否有Unix工具可用（GitBash环境检测）
func (stp *SystemToolProcessor) hasUnixTools() bool {
	// 尝试执行sh命令来检测是否在Unix-like环境中
	cmd := exec.Command("sh", "-c", "echo test")
	err := cmd.Run()
	return err == nil
}

// NativeGoProcessor Go内置处理器
type NativeGoProcessor struct{}

// NewNativeGoProcessor 创建Go内置处理器
func NewNativeGoProcessor() *NativeGoProcessor {
	return &NativeGoProcessor{}
}

// SearchInFile 使用Go内置方法搜索文件
func (ngp *NativeGoProcessor) SearchInFile(filename, pattern string) ([]SearchResult, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			microg.E("Close file error: %s\n", err)
		}
	}(file)

	var scanner *bufio.Scanner

	// 判断是否为压缩文件
	if strings.HasSuffix(strings.ToLower(filename), ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("解压文件失败: %v", err)
		}
		defer func(gzReader *gzip.Reader) {
			err := gzReader.Close()
			if err != nil {
				microg.E("Close gzip reader error: %s\n", err)
			}
		}(gzReader)
		scanner = bufio.NewScanner(gzReader)
	} else {
		scanner = bufio.NewScanner(file)
	}

	// 编译正则表达式（如果模式看起来像正则表达式）
	var regex *regexp.Regexp
	if isRegexPattern(pattern) {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			// 如果正则表达式编译失败，回退到字符串匹配
			regex = nil
		}
	}

	var results []SearchResult
	lineNumber := 0
	// maxResults := 1000 // 限制结果数量防止内存溢出

	for scanner.Scan() && len(results) < setupConfig.Logging.MaxResults {
		lineNumber++
		line := scanner.Text()

		var matched bool
		var matchedText string

		if regex != nil {
			// 使用正则表达式匹配
			if regex.MatchString(line) {
				matched = true
				matchedText = regex.FindString(line)
			}
		} else {
			// 使用字符串包含匹配
			if strings.Contains(line, pattern) {
				matched = true
				// pattern 可能包含 grep 参数，需要过滤掉
				matchedText = filterGrepParams(pattern)
			}
		}

		if matched {
			results = append(results, SearchResult{
				LineNumber: lineNumber,
				Content:    line,
				Matched:    matchedText,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return results, nil
}

// SearchInFileStream 使用Go内置方法流式搜索文件
func (ngp *NativeGoProcessor) SearchInFileStream(ctx context.Context, filename, pattern string, callback func(SearchResult)) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			microg.E("Close file error: %s\n", err)
		}
	}(file)

	var scanner *bufio.Scanner

	// 判断是否为压缩文件
	if strings.HasSuffix(strings.ToLower(filename), ".gz") {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("解压文件失败: %v", err)
		}
		defer func(gzReader *gzip.Reader) {
			err := gzReader.Close()
			if err != nil {
				microg.E("Close gzip reader error: %s\n", err)
			}
		}(gzReader)
		scanner = bufio.NewScanner(gzReader)
	} else {
		scanner = bufio.NewScanner(file)
	}

	// 编译正则表达式（如果模式看起来像正则表达式）
	var regex *regexp.Regexp
	if isRegexPattern(pattern) {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			// 如果正则表达式编译失败，回退到字符串匹配
			regex = nil
		}
	}

	lineNumber := 0
	// maxResults := 1000 // 限制结果数量防止内存溢出
	resultCount := 0

	for scanner.Scan() {
		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		lineNumber++
		line := scanner.Text()

		var matched bool
		var matchedText string

		if regex != nil {
			// 使用正则表达式匹配
			if regex.MatchString(line) {
				matched = true
				matchedText = regex.FindString(line)
			}
		} else {
			// 使用字符串包含匹配
			if strings.Contains(line, pattern) {
				matched = true
				// pattern 可能包含 grep 参数，需要过滤掉
				matchedText = filterGrepParams(pattern)
			}
		}

		if matched {
			result := SearchResult{
				LineNumber: lineNumber,
				Content:    line,
				Matched:    matchedText,
			}
			callback(result)
			resultCount++
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	return nil
}

// validateFilePath 验证文件路径安全性
func validateFilePath(filename string) error {
	// 防止目录遍历攻击
	cleanPath := filepath.Clean(filename)
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("不安全的文件路径: %s", filename)
	}

	// 只允许访问当前目录及子目录的日志文件
	if !strings.HasSuffix(cleanPath, ".log") && !strings.HasSuffix(cleanPath, ".gz") {
		return fmt.Errorf("只允许访问日志文件: %s", filename)
	}

	return nil
}

// parseGrepOutput 解析grep命令输出
func parseGrepOutput(output, pattern string, hasAround bool) []SearchResult {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var results []SearchResult

	// pattern 可能包含 grep 参数，需要过滤掉
	filteredPattern := filterGrepParams(pattern)
	for _, line := range lines {
		if line == "" {
			continue
		}

		// 解析grep -n的输出格式: 行号:内容 or 行号-内容
		re := regexp.MustCompile(`^(\d+)([-:])(.*)`) // 匹配行号、分隔符和内容
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 {
			lineNum, _ := strconv.Atoi(matches[1])
			delimiter := matches[2]
			content := matches[3]
			if !hasAround {
				results = append(results, SearchResult{
					LineNumber: lineNum,
					Content:    content,
					Matched:    filteredPattern,
				})
			} else {
				switch delimiter {
				case ":":
					// 匹配行号:内容
					results = append(results, SearchResult{
						LineNumber: lineNum,
						Content:    content,
						Matched:    filteredPattern,
					})
				case "-":
					// 匹配行号-内容 追加到上一条搜索结果
					if len(results) > 0 {
						lastResult := &results[len(results)-1]
						lastResult.Content += "\n" + content
					}
				}
			}

			if len(results) >= setupConfig.Logging.MaxResults {
				break
			}
		}
	}

	return results
}

// isRegexPattern 判断是否为正则表达式模式
func isRegexPattern(pattern string) bool {
	// 简单判断：包含正则表达式特殊字符
	regexChars := []string{"[", "]", "(", ")", "{", "}", "^", "$", ".", "*", "+", "?", "|", "\\"}
	for _, char := range regexChars {
		if strings.Contains(pattern, char) {
			return true
		}
	}
	return false
}
