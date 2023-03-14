package model

import (
	"fmt"
)

type Config struct {
	Development     bool     `json:"development"`
	BaseDir         string   `json:"base_dir"`
	GitUrl          string   `json:"git_url"`
	Username        string   `json:"username"`
	Email           string   `json:"email"`
	Reponame        string   `json:"reponame"`
	Branch          string   `json:"branch"`
	SshConfigDir    string   `json:"ssh_dir" `
	RepoToken       string   `json:"repo_token"`
	InfraRepoSuffix string   `json:"infra_repo_suffix"`
	ImageTag        string   `json:"image_tag"`
	TagLocation     string   `json:"tag_location"`
	Stages          []string `json:"stages"`
	Env             string   `json:"env"`
	FromBranch      string   `json:"from_branch"`
	ToBranch        string   `json:"to_branch"`
	Force           bool     `json:"force"`
	ResourcesOnly   bool     `json:"resources_only"`
	Descriptor      string   `json:"descriptor"`
	CommitRef       string   `json:"commit_ref"`
}

func (c *Config) ApplicationClonePath() string {

	return fmt.Sprintf("%s%s", c.BaseDir, c.Reponame)
}

func (c *Config) InfrastructureClonePath() string {
	return fmt.Sprintf("%s%s%s", c.BaseDir, c.Reponame, c.InfraRepoSuffix)
}

func (c *Config) IsPushEnabled() bool {
	return !c.Development
}
