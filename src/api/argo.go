package api

import (
	"fmt"
	"gepaplexx/git-workflows/logger"
	"gepaplexx/git-workflows/model"
	"gepaplexx/git-workflows/utils"
	"github.com/go-git/go-git/v5"
	"github.com/otiai10/copy"
	"os"
	"strings"
)

const (
	ApplicationsetLocation = "%s/argocd/applicationset.yaml"
	ValuesLocation         = "%s/apps/env/%s/values.yaml"
	TemplateLocation       = "%s/apps/env/%s"
)

func UpdateArgoApplicationSet(c *model.Config, repo *git.Repository) {
	logger.Info("Updating ArgoCD Application")
	logger.Debug("Env: %s", c.Env)

	wt := checkout(repo, "main", false)

	filePath := fmt.Sprintf(ValuesLocation, wt.Filesystem.Root(), c.Env)
	logger.Debug("Updating file: %s", filePath)
	updateImageTag(c, filePath)

	commitAndPush(c, wt, repo, fmt.Sprintf("updated image tag to %s\nTriggered by ref: %s", c.ImageTag, c.CommitRef))
}

func ArgoCreateEnvironment(c *model.Config, repo *git.Repository) {
	logger.Info("Creating entry in ArgoCD ApplicationSet")
	logger.Debug("Env: %s", c.Env)

	wt := checkout(repo, "main", false)
	filePath := fmt.Sprintf(ApplicationsetLocation, wt.Filesystem.Root())
	addEnvironmentToApplicationSet(c, filePath)
	copyTemplateDir(wt, c, filePath)

	copyApplicationSet(c, filePath)
	commitAndPush(c, wt, repo, fmt.Sprintf("Added ApplicationSet entry and templates for %s", c.Env))
}

func DeleteArgoEnvironment(c *model.Config, repo *git.Repository) {
	logger.Info("Deleting entry in ArgoCD ApplicationSet")
	logger.Debug("Env: %s", c.Env)

	protectEnvironments(c)

	wt := checkout(repo, "main", false)
	filePath := fmt.Sprintf(ApplicationsetLocation, wt.Filesystem.Root())
	removeEnvironmentFromApplicationSet(c, filePath)
	deleteTemplateDir(wt, c)
	copyApplicationSet(c, filePath)
	commitAndPush(c, wt, repo, fmt.Sprintf("Removed ApplicationSet entry and templates for %s", c.Env))
}

func protectEnvironments(c *model.Config) {
	for _, stage := range c.Stages {
		if stage == c.Branch && c.Force == false {
			logger.Fatal("%s is a predefined stage. It cannot be deleted via automation. Override this check by setting --force flag", stage)
		}
	}
}

func UpdateAllStages(c *model.Config, wt *git.Worktree, repo *git.Repository) {
	logger.Info("Updating all stages to new image tag to prepare deployment")

	values := ParseYaml(fmt.Sprintf(ValuesLocation, wt.Filesystem.Root(), "main"))
	imagetagNode, err := FindNode(values.Content[0], c.TagLocation)
	utils.CheckIfError(err)
	c.ImageTag = imagetagNode.Value

	for _, stage := range c.Stages {
		filePath := fmt.Sprintf(ValuesLocation, wt.Filesystem.Root(), stage)
		logger.Debug("Updating file: %s", filePath)
		updateImageTag(c, filePath)
	}
	commitAndPush(c, wt, repo, fmt.Sprintf("updated image tag to %s", c.ImageTag))
}

func updateImageTag(c *model.Config, filepath string) {
	nodes := ParseYaml(filepath)
	tagNode, err := FindNode(nodes.Content[0], c.TagLocation)
	utils.CheckIfError(err)
	tagNode.Value = c.ImageTag
	WriteYaml(nodes, filepath)
}

func addEnvironmentToApplicationSet(c *model.Config, path string) {
	logger.Info("Adding (%s, %s) to ApplicationSet %s", c.Env, "main", path)
	nodes := ParseYaml(path)
	envNode, err := FindNode(nodes.Content[0], AppsetEnvPath)
	utils.CheckIfError(err)
	envNode.Content = append(envNode.Content, NewEnvNode(c.Env, "main"))
	WriteYaml(nodes, path)
}

func copyTemplateDir(wt *git.Worktree, c *model.Config, applicationset string) {
	nodes := ParseYaml(applicationset)
	fromBranch := strings.ReplaceAll(strings.ReplaceAll(c.FromBranch, "/", "-"), "_", "-")
	res, err := FindClusterWithBranch(nodes.Content[0], fromBranch)
	utils.CheckIfError(err)

	sourceDir := fmt.Sprintf(TemplateLocation, wt.Filesystem.Root(), res)
	targetTemplateDir := fmt.Sprintf(TemplateLocation, wt.Filesystem.Root(), c.Env)

	logger.Info("Copying %s to %s", sourceDir, targetTemplateDir)
	err = copy.Copy(sourceDir, targetTemplateDir)
	utils.CheckIfError(err)
}

func deleteTemplateDir(wt *git.Worktree, c *model.Config) {
	targetTemplateDir := fmt.Sprintf(TemplateLocation, wt.Filesystem.Root(), c.Env)

	logger.Info("Deleting %s", targetTemplateDir)
	err := os.RemoveAll(targetTemplateDir)
	utils.CheckIfError(err)
}

func removeEnvironmentFromApplicationSet(c *model.Config, path string) {
	logger.Info("Removing %s from ApplicationSet %s", c.Env, path)
	rootNode := ParseYaml(path)
	err := DeleteEnvFromApplicationset(rootNode.Content[0], c.Env)
	utils.CheckIfError(err)
	WriteYaml(rootNode, path)
}

func copyApplicationSet(c *model.Config, filePath string) {
	logger.Info("Copying ApplicationSet")
	err := copy.Copy(filePath, fmt.Sprintf("%s/application.yaml", c.BaseDir))
	utils.CheckIfError(err)
}
