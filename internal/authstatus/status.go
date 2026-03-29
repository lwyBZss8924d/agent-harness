package authstatus

import (
	"os"
	"os/exec"

	"github.com/lwyBZss8924d/agent-harness/aih-toolkit/internal/config"
)

type ToolStatus struct {
	Name            string `json:"name"`
	Required        bool   `json:"required"`
	ConfigDir       string `json:"config_dir,omitempty"`
	ConfigDirExists bool   `json:"config_dir_exists"`
	AuthFile        string `json:"auth_file,omitempty"`
	AuthFileExists  bool   `json:"auth_file_exists"`
	Binary          string `json:"binary,omitempty"`
	Available       bool   `json:"available"`
}

type Status struct {
	Tools []ToolStatus `json:"tools"`
}

func Collect(cfg config.Config) Status {
	tools := make([]ToolStatus, 0, len(cfg.Auth.Tools))
	for _, item := range cfg.Auth.Tools {
		binaryPath := commandPath(item.Binary)
		tools = append(tools, ToolStatus{
			Name:            item.Name,
			Required:        item.Required,
			ConfigDir:       item.ConfigDir,
			ConfigDirExists: isDir(item.ConfigDir),
			AuthFile:        item.AuthFile,
			AuthFileExists:  isFile(item.AuthFile),
			Binary:          binaryPath,
			Available:       binaryPath != "",
		})
	}
	return Status{Tools: tools}
}

func commandPath(name string) string {
	if name == "" {
		return ""
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return ""
	}
	return path
}

func isDir(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isFile(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
