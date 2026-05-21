package rule

import (
	"runtime"
	"strings"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
)

var _ RuleItem = (*ProcessPathItem)(nil)

type ProcessPathItem struct {
	processes  []string
	processMap map[string]bool
}

func NewProcessPathItem(processNameList []string) *ProcessPathItem {
	rule := &ProcessPathItem{
		processes:  processNameList,
		processMap: make(map[string]bool),
	}
	isWindows := runtime.GOOS == "windows"
	for _, processName := range processNameList {
		if isWindows {
			rule.processMap[strings.ToLower(processName)] = true
		} else {
			rule.processMap[processName] = true
		}
	}
	return rule
}

func (r *ProcessPathItem) Match(metadata *adapter.InboundContext) bool {
	if metadata.ProcessInfo == nil {
		return false
	}
	processPath := metadata.ProcessInfo.ProcessPath
	if processPath != "" {
		if runtime.GOOS == "windows" {
			processPath = strings.ToLower(processPath)
		}
		if r.processMap[processPath] {
			return true
		}
	}
	if C.IsAndroid {
		for _, packageName := range metadata.ProcessInfo.AndroidPackageNames {
			if r.processMap[packageName] {
				return true
			}
		}
	}
	return false
}

func (r *ProcessPathItem) String() string {
	var description string
	pLen := len(r.processes)
	if pLen == 1 {
		description = "process_path=" + r.processes[0]
	} else {
		description = "process_path=[" + strings.Join(r.processes, " ") + "]"
	}
	return description
}
