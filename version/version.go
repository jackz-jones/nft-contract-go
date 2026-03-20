// Package version 版本信息
package version

import "fmt"

var (
	// CommitID git commit hash
	CommitID = ""

	// BuildTime when to build
	BuildTime = ""

	// Version what version
	Version = ""
)

// VerInfo version info
type VerInfo struct {
	Version   string `json:"version"`
	CommitID  string `json:"commitID"`
	BuildTime string `json:"buildTime"`
}

// Print print current cmd version
func (v VerInfo) Print() string {
	return fmt.Sprintf("\nCurrent version: %s\nCommit hash: %s\nBuild time: %s\n\n", v.Version, v.CommitID, v.BuildTime)
}
