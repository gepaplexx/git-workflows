package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
	"strings"
)

var deleteCmd = &cobra.Command{
	Use:   "argo-delete",
	Short: "Deletes an argocd application on branch deletion and updates git repository accordingly.",
	Run: func(cmd *cobra.Command, args []string) {
		deleteArgo(&Config)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteArgo(c *model.Config) {
	c.Env = strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	if c.Development {
		logger.Debug("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}
	argoDeletePrerequisites(c)
	repo := api.CloneRepo(c, "main", false)
	api.ArgoCreateEnvironment(c, repo)
}

func argoDeletePrerequisites(c *model.Config) {
	prerequisites(c)
}
