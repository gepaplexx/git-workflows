package cmd

import (
	"gepaplexx/git-workflows/api"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"github.com/spf13/cobra"
	"strings"
)

var updateCmd = &cobra.Command{
	Use:   "argo-update",
	Short: "Updates argocd application in infrastructure repository to handle deployments",
	Run: func(cmd *cobra.Command, args []string) {
		updateArgo(&Config)
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func updateArgo(c *model.Config) {
	c.Env = strings.ReplaceAll(strings.ReplaceAll(c.Branch, "/", "-"), "_", "-")
	if c.Development {
		logger.Debug("Development mode enabled. Using local configuration.")
		developmentMode(c)
	}
	argoUpdatePrerequisites(c)
	repo := api.CloneRepo(c, "main", false)
	api.UpdateArgoApplicationSet(c, repo)
}

func argoUpdatePrerequisites(c *model.Config) {
	if c.ImageTag == "" {
		logger.Fatal("Image tag must be set")
	}

	if c.Branch == "" {
		logger.Fatal("Branch must be set")
	}

}
