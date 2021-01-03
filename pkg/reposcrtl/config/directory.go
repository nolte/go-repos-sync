package config

import "fmt"

type DirectoryType string

const (
	DirectoryTypeGroup   DirectoryType = "group"
	DirectoryTypeProject DirectoryType = "project"
	DirectoryTypeUser    DirectoryType = "user"
)

func (r DirectoryType) String() string {
	return string(r)
}

func DirectoryTypeFromString(directoryType string) (DirectoryType, error) {
	switch directoryType {
	case DirectoryTypeGroup.String():
		return DirectoryTypeGroup, nil
	case DirectoryTypeProject.String():
		return DirectoryTypeProject, nil
	case DirectoryTypeUser.String():
		return DirectoryTypeUser, nil
	default:
		return "", fmt.Errorf("not Supported Type String %s", directoryType)
	}
}

type Directory struct {
	Path     string              `json:"path"`
	Contents []DirectoryContents `json:"contents,omitempty"`
}

type Checkout struct {
	FsPrefix                 string `json:"fs_prefix"`
	Path                     string `json:"path"`
	KeepsGroupPath           bool   `json:"keepsGroupPath"`
	KeepsBaseGroupNameInPath bool   `json:"keepsBaseGroupNameInPath"`
}

type DirectoryContents struct {
	Ref      DirectoryContentsRef `json:"ref"`
	Checkout Checkout             `json:"checkout"`
}

type DirectoryContentsRef struct {
	Remote string                 `json:"remote"`
	IDS    []string               `json:"ids"`
	Kind   DirectoryType          `json:"kind"`
	Filter DirectoryContentFilter `json:"filter"`
}

type DirectoryContentFilter struct {
	Excludes              []string `json:"excludes,omitempty"`
	Includes              []string `json:"includes,omitempty"`
	IncludesArchivedRepos bool     `json:"includesArchivedRepos"`
	ArchivedReposOnly     bool     `json:"archivedReposOnly"`
}
