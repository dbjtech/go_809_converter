package converter

import (
	"os/exec"
	"runtime"
)

// ToolChecker 工具检测器
type ToolChecker struct {
	HasCat  bool `json:"has_cat"`
	HasGzip bool `json:"has_gzip"`
	HasGrep bool `json:"has_grep"`
}

// NewToolChecker 创建新的工具检测器
func NewToolChecker() *ToolChecker {
	return &ToolChecker{}
}

// CheckSystemTools 检测系统工具可用性
func (tc *ToolChecker) CheckSystemTools() {
	tc.HasCat = tc.checkTool("cat")
	tc.HasGzip = tc.checkTool("gzip")
	tc.HasGrep = tc.checkTool("grep")

	// Windows系统特殊处理
	if runtime.GOOS == "windows" {
		// Windows下可能使用type代替cat
		if !tc.HasCat {
			tc.HasCat = tc.checkTool("type")
		}
		// Windows下可能使用findstr代替grep
		if !tc.HasGrep {
			tc.HasGrep = tc.checkTool("findstr")
		}
	}
}

// checkTool 检测单个工具是否可用
func (tc *ToolChecker) checkTool(toolName string) bool {
	_, err := exec.LookPath(toolName)
	return err == nil
}

// GetAvailableTools 获取可用工具列表
func (tc *ToolChecker) GetAvailableTools() []string {
	var tools []string

	if tc.HasCat {
		tools = append(tools, "cat")
	}
	if tc.HasGzip {
		tools = append(tools, "gzip")
	}
	if tc.HasGrep {
		tools = append(tools, "grep")
	}

	return tools
}

// HasAllTools 检查是否所有工具都可用
func (tc *ToolChecker) HasAllTools() bool {
	return tc.HasCat && tc.HasGzip && tc.HasGrep
}

// GetToolStatus 获取工具状态描述
func (tc *ToolChecker) GetToolStatus() map[string]string {
	status := make(map[string]string)

	if tc.HasCat {
		status["cat"] = "可用"
	} else {
		status["cat"] = "不可用"
	}

	if tc.HasGzip {
		status["gzip"] = "可用"
	} else {
		status["gzip"] = "不可用"
	}

	if tc.HasGrep {
		status["grep"] = "可用"
	} else {
		status["grep"] = "不可用"
	}

	return status
}
