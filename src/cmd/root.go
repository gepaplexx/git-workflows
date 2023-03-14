/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/spf13/cobra"
	"os"
)

var Version string
var Config model.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-workflows",
	Short: "handle git operations for argo workflows",
	Long: `cli application to handle git operations in argo workflows. 
	
Subcommands reflect the stages in the workflow. `,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	utils.CheckIfError(err)
}

func init() {

	// global flags
	rootCmd.PersistentFlags().BoolVar(&Config.Development, "dev", false, "enable development mode")
	rootCmd.PersistentFlags().StringVarP(&Config.BaseDir, "path", "p", "/mnt/out/", "base directory for all operations")
	rootCmd.PersistentFlags().StringVar(&Config.Username, "commit-user", "argo-ci", "username for git operations")
	rootCmd.PersistentFlags().StringVar(&Config.Email, "commit-email", "argo-ci@gepardec.com", "email for git operations")

	rootCmd.PersistentFlags().StringVarP(&Config.GitUrl, "url", "u", "", "git url for the repository")
	rootCmd.PersistentFlags().StringVar(&Config.Reponame, "name", "", "name of the repository")
	rootCmd.PersistentFlags().StringVarP(&Config.Branch, "branch", "b", "main", "branch to checkout")
	rootCmd.PersistentFlags().StringVar(&Config.SshConfigDir, "ssh-config-dir", "/root/.ssh/", "directory for ssh known_hosts and private key")
	rootCmd.PersistentFlags().StringVar(&Config.RepoToken, "token", "", "token to access the repository")

	rootCmd.PersistentFlags().StringVar(&Config.InfraRepoSuffix, "infra-repo-suffix", "-ci", "Suffix for infrastructure git repository")
	rootCmd.PersistentFlags().StringVar(&Config.ImageTag, "tag", "", "Commit-Hash/Image-Tag for the deployment change")
	rootCmd.PersistentFlags().StringVar(&Config.TagLocation, "image-tag-location", "image.tag", "Location of image-tag in the infrastructure repository")
	rootCmd.PersistentFlags().StringSliceVar(&Config.Stages, "stages", []string{"main", "dev", "qa", "prod"}, "deployment stages")
	rootCmd.PersistentFlags().StringVar(&Config.FromBranch, "from-branch", "main", "Base branch for argo-create | Branch from which to deploy from")
	rootCmd.PersistentFlags().StringVar(&Config.ToBranch, "to-branch", "", "Target brancht for deployments")

	updateCmd.PersistentFlags().StringVar(&Config.CommitRef, "commit-ref", "no reference supplied", "Reference to original commit that triggered the workflow")
	deleteCmd.PersistentFlags().BoolVarP(&Config.Force, "force", "f", false, "allows deletion of protected environments. Remember: with great power comes great responsibility!")

	deployCmd.PersistentFlags().BoolVar(&Config.ResourcesOnly, "resources-only", false, "only deploy resources, no application")
	descriptorCmd.PersistentFlags().StringVar(&Config.Descriptor, "descriptor", "workflow-descriptor.json", "full path and name to workflow descriptor")
}

func prerequisites(c *model.Config) {
	if c.GitUrl == "" {
		logger.Fatal("Git URL must be set")
	}

	if c.Reponame == "" {
		logger.Fatal("Reponame must be set")
	}

	if c.Branch == "" {
		logger.Fatal("Branch must be set")
	}
}

func developmentMode(c *model.Config) {
	logger.Info("Development mode enabled. Using local configuration.")
	c.BaseDir = "../../tmp/"
	c.SshConfigDir = os.Getenv("HOME") + "/.ssh/"
	c.ImageTag = "abcdefg"
	err := os.RemoveAll(c.BaseDir)
	logger.EnableDebug()
	utils.CheckIfError(err)
	c.GitUrl = "git@github.com:gepaplexx-demos/demo-microservice.git"
	c.Reponame = "demo-microservice"
}
