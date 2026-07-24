// Package util provides utility functions used across ProxyBridge.
package util

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// GetHomeDir returns the user's home directory
func GetHomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// GetOS returns the current operating system
func GetOS() string {
	return runtime.GOOS
}

// GetArch returns the current architecture
func GetArch() string {
	return runtime.GOARCH
}

// GetGoVersion returns the Go version
func GetGoVersion() string {
	return runtime.Version()
}

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return strings.EqualFold(runtime.GOOS, "windows")
}

// IsMacOS returns true if running on macOS
func IsMacOS() bool {
	return strings.EqualFold(runtime.GOOS, "darwin")
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return strings.EqualFold(runtime.GOOS, "linux")
}

// IsWSL returns true if running on WSL
func IsWSL() bool {
	if !IsLinux() {
		return false
	}

	// Check for WSL by looking at /proc/version
	if data, err := os.ReadFile("/proc/version"); err == nil {
		return strings.Contains(strings.ToLower(string(data)), "microsoft")
	}

	// Check for WSL by looking at the hostname
	if hostname, err := os.Hostname(); err == nil {
		return strings.Contains(strings.ToLower(hostname), "wsl")
	}

	return false
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirectoryExists checks if a directory exists
func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(path string) error {
	if !DirectoryExists(path) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

// FileSize returns the size of a file in bytes
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// MoveFile moves a file from src to dst
func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

// DeleteFile deletes a file
func DeleteFile(path string) error {
	return os.Remove(path)
}

// DeleteDirectory deletes a directory and all contents
func DeleteDirectory(path string) error {
	return os.RemoveAll(path)
}

// FileContents returns the contents of a file as a string
func FileContents(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFile writes content to a file
func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// FileLineCount returns the number of lines in a file
func FileLineCount(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	count := 0
	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		for i := 0; i < n; i++ {
			if buf[i] == '\n' {
				count++
			}
		}
		if n < len(buf) {
			break
		}
	}

	return count, nil
}

// FindExecutable looks for an executable in PATH
func FindExecutable(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, nil
	}
	return absPath, nil
}

// IsExecutable checks if a file is executable
func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := info.Mode()
	return mode&0111 != 0
}

// SplitLines splits a string into lines
func SplitLines(text string) []string {
	lines := strings.Split(text, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// JoinLines joins lines with newlines
func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// TruncateString truncates a string to max length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// EllipsisForPath creates an ellipsis representation of a long path
func EllipsisForPath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}

	home := GetHomeDir()
	relPath, err := filepath.Rel(home, path)
	if err == nil {
		path = "~/" + relPath
	}

	if len(path) <= maxLen {
		return path
	}

	// Build path with ellipsis
	parts := strings.Split(path, string(os.PathSeparator))
	if len(parts) <= 2 {
		return path
	}

	result := parts[0] + "/"
	for i := 1; i < len(parts)-1; i++ {
		if i == len(parts)/2 {
			result += "..."
		} else {
			result += parts[i] + "/"
		}
	}
	result += parts[len(parts)-1]

	if len(result) <= maxLen {
		return result
	}

	// Final fallback
	return path[len(path)-maxLen+3:]
}

// HumanReadableSize formats bytes as human readable size
func HumanReadableSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// HumanReadableDuration formats duration as human readable string
func HumanReadableDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%v", d)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

// SafePathJoin joins paths safely
func SafePathJoin(elem ...string) string {
	return filepath.Join(elem...)
}

// PathExists returns true if path exists
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsEmptyDirectory checks if a directory is empty
func IsEmptyDirectory(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

// CreateTempFile creates a temporary file
func CreateTempFile(dir, prefix string) (string, error) {
	file, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return file.Name(), nil
}

// CreateTempDirectory creates a temporary directory
func CreateTempDirectory(dir, prefix string) (string, error) {
	dirPath, err := os.MkdirTemp(dir, prefix)
	if err != nil {
		return "", err
	}
	return dirPath, nil
}

// ChmodFile changes file permissions
func ChmodFile(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// FilePermissions returns file permissions
func FilePermissions(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Mode().Perm(), nil
}

// IsProcessRunning checks if a process with the given PID is running
func IsProcessRunning(pid int) bool {
	// Try to kill process 0 to check if it exists
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = proc.Signal(0)
	return err == nil
}

// GetProcessCPUUsage returns CPU usage of a process
func GetProcessCPUUsage(pid int) (float64, error) {
	// Platform-specific implementation
	return 0, fmt.Errorf("not implemented")
}

// GetProcessMemoryUsage returns memory usage of a process
func GetProcessMemoryUsage(pid int) (int64, error) {
	// Platform-specific implementation
	return 0, fmt.Errorf("not implemented")
}
