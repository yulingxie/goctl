package git

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
	"strings"

	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/console"
	"gitlab.kaiqitech.com/k7game/server/tools/goctl/util/execx"
	"gitlab.kaiqitech.com/nitro/nitro/v3/util/pathx"
)

type Git struct {
	console.Console
	gitClient *gitlab.Client
}

func NewGit(token string) *Git {
	gitClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(GITLAB_URL))
	if err != nil {
		panic(err)
	}
	return &Git{Console: console.NewConsole(false), gitClient: gitClient}
}

func (self *Git) CreateSngGateway(projectName, description string) (*gitlab.Project, error) {
	// 创建项目
	project, _, err := self.gitClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:                 gitlab.String(projectName),
		NamespaceID:          gitlab.Int(SNG_GATEWAY_GROUP),
		Visibility:           gitlab.Visibility(gitlab.PrivateVisibility),
		Description:          gitlab.String(description),
		InitializeWithReadme: gitlab.Bool(true),
	})

	// 创建dev、master分支
	self.gitClient.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
		Branch: gitlab.String("dev"),
		Ref:    gitlab.String("master"),
	})

	// 设置默认分支为dev
	self.gitClient.Projects.EditProject(project.ID, &gitlab.EditProjectOptions{
		DefaultBranch: gitlab.String("dev"),
	})

	// 设置各分支权限
	self.gitClient.Branches.ProtectBranch(project.ID, "dev", &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  gitlab.Bool(true),
		DevelopersCanMerge: gitlab.Bool(true),
	})

	self.gitClient.Branches.ProtectBranch(project.ID, "master", &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  gitlab.Bool(false),
		DevelopersCanMerge: gitlab.Bool(false),
	})

	// 设置推送dev自动构建的webhook
	self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
		PushEvents:             gitlab.Bool(true),
		PushEventsBranchFilter: gitlab.String("dev"),
		URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_GATEWAY_URL_DEV, projectName)),
		Token:                  gitlab.String(SNG_JENKINS_TOKEN_DEV),
	})
	// 设置推送master自动构建的webhook
	self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
		PushEvents:             gitlab.Bool(true),
		PushEventsBranchFilter: gitlab.String("master"),
		URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_GATEWAY_URL_MASTER, projectName)),
		Token:                  gitlab.String(SNG_JENKINS_TOKEN_DEV),
	})

	self.Info("创建项目: %d %s", project.ID, project.WebURL)
	return project, err
}

func (self *Git) UpdateAllSngGateway() error {
	// 获取项目
	projects, _, err := self.gitClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: gitlab.String("-gw"),
	})
	if err != nil {
		return err
	}
	for _, project := range projects {
		if !strings.Contains(project.WebURL, "https://gitlab.kaiqitech.com/k7game/server/supports/") {
			continue
		}
		projectName := project.Name
		// 删除老的所有hook
		hooks, _, _ := self.gitClient.Projects.ListProjectHooks(project.ID, nil)
		for _, hook := range hooks {
			self.gitClient.Projects.DeleteProjectHook(project.ID, hook.ID)
		}
		// 设置推送dev自动构建的webhook
		self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
			PushEvents:             gitlab.Bool(true),
			MergeRequestsEvents:    gitlab.Bool(true),
			PushEventsBranchFilter: gitlab.String("dev"),
			URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_GATEWAY_URL_DEV, projectName)),
			Token:                  gitlab.String(SNG_JENKINS_TOKEN_DEV),
		})
		// 设置推送master自动构建的webhook
		self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
			PushEvents:             gitlab.Bool(true),
			MergeRequestsEvents:    gitlab.Bool(true),
			PushEventsBranchFilter: gitlab.String("master"),
			URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_GATEWAY_URL_MASTER, projectName)),
			Token:                  gitlab.String(SNG_JENKINS_TOKEN_MASTER),
		})
	}
	return err
}

func (self *Git) CreateSngService(projectName, description string) (*gitlab.Project, error) {
	// 创建项目
	project, _, err := self.gitClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:                 gitlab.String(projectName),
		NamespaceID:          gitlab.Int(SNG_SERVICE_GROUP),
		Visibility:           gitlab.Visibility(gitlab.PrivateVisibility),
		Description:          gitlab.String(description),
		InitializeWithReadme: gitlab.Bool(true),
	})

	// 创建dev、master分支
	self.gitClient.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
		Branch: gitlab.String("dev"),
		Ref:    gitlab.String("master"),
	})

	// 设置默认分支为dev
	self.gitClient.Projects.EditProject(project.ID, &gitlab.EditProjectOptions{
		DefaultBranch: gitlab.String("dev"),
	})

	// 设置各分支权限
	self.gitClient.Branches.ProtectBranch(project.ID, "dev", &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  gitlab.Bool(true),
		DevelopersCanMerge: gitlab.Bool(true),
	})

	self.gitClient.Branches.ProtectBranch(project.ID, "master", &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  gitlab.Bool(false),
		DevelopersCanMerge: gitlab.Bool(false),
	})

	// 设置推送dev自动构建的webhook
	self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
		PushEvents:             gitlab.Bool(true),
		PushEventsBranchFilter: gitlab.String("dev"),
		URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_URL_DEV, projectName)),
		Token:                  gitlab.String(SNG_JENKINS_TOKEN_DEV),
	})
	// 设置推送master自动构建的webhook
	self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
		PushEvents:             gitlab.Bool(true),
		PushEventsBranchFilter: gitlab.String("master"),
		URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_URL_MASTER, projectName)),
		Token:                  gitlab.String(SNG_JENKINS_TOKEN_DEV),
	})

	self.Info("创建项目: %d %s", project.ID, project.WebURL)
	return project, err
}

func (self *Git) UpdateAllSngService() error {
	// 获取项目
	projects, _, err := self.gitClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: gitlab.String("service"),
	})
	if err != nil {
		return err
	}
	for _, project := range projects {
		if !strings.Contains(project.WebURL, "https://gitlab.kaiqitech.com/k7game/server/services/") {
			continue
		}
		projectName := project.Name
		// 删除老的所有hook
		hooks, _, _ := self.gitClient.Projects.ListProjectHooks(project.ID, nil)
		for _, hook := range hooks {
			self.gitClient.Projects.DeleteProjectHook(project.ID, hook.ID)
		}
		// 设置推送dev自动构建的webhook
		self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
			PushEvents:             gitlab.Bool(true),
			MergeRequestsEvents:    gitlab.Bool(true),
			PushEventsBranchFilter: gitlab.String("dev"),
			URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_URL_DEV, projectName)),
			Token:                  gitlab.String(SNG_JENKINS_TOKEN_DEV),
		})
		// 设置推送master自动构建的webhook
		self.gitClient.Projects.AddProjectHook(project.ID, &gitlab.AddProjectHookOptions{
			PushEvents:             gitlab.Bool(true),
			MergeRequestsEvents:    gitlab.Bool(true),
			PushEventsBranchFilter: gitlab.String("master"),
			URL:                    gitlab.String(fmt.Sprintf(SNG_JENKINS_URL_MASTER, projectName)),
			Token:                  gitlab.String(SNG_JENKINS_TOKEN_MASTER),
		})
	}
	return err
}

func (self *Git) CreateSngServiceTest(projectName, description string) (*gitlab.Project, error) {
	// 创建项目
	project, _, err := self.gitClient.Projects.CreateProject(&gitlab.CreateProjectOptions{
		Name:                 gitlab.String(projectName),
		NamespaceID:          gitlab.Int(SNG_SERVICE_TEST_GROUP),
		Visibility:           gitlab.Visibility(gitlab.PrivateVisibility),
		Description:          gitlab.String(description),
		InitializeWithReadme: gitlab.Bool(true),
	})

	// 创建test分支
	self.gitClient.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
		Branch: gitlab.String("test"),
		Ref:    gitlab.String("master"),
	})

	// 设置默认分支为test
	self.gitClient.Projects.EditProject(project.ID, &gitlab.EditProjectOptions{
		DefaultBranch: gitlab.String("test"),
	})

	// 设置各分支权限
	self.gitClient.Branches.ProtectBranch(project.ID, "test", &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  gitlab.Bool(true),
		DevelopersCanMerge: gitlab.Bool(true),
	})

	self.gitClient.Branches.ProtectBranch(project.ID, "master", &gitlab.ProtectBranchOptions{
		DevelopersCanPush:  gitlab.Bool(false),
		DevelopersCanMerge: gitlab.Bool(false),
	})
	self.Info("创建项目: %d %s", project.ID, project.WebURL)
	return project, err
}

func (self *Git) GetProject(name string) error {
	projects, _, err := self.gitClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: gitlab.String(name),
	})
	if err != nil {
		return err
	}
	for _, project := range projects {
		self.Info("%d %s", project.ID, project.WebURL)
	}
	return nil
}

func (self *Git) OpenProject(name string) error {
	projects, _, err := self.gitClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: gitlab.String(name),
	})
	if err != nil {
		return err
	}
	if len(projects) > 0 {
		if _, err := execx.Run(fmt.Sprintf("open %s", projects[0].WebURL), "./"); err != nil {
			return err
		}
	}

	return nil
}

func (self *Git) CloneProject(names []string, ids []int) {
	//如果指定了id，则使用id查找项目信息
	urls := []string{}

	for _, id := range ids {
		project, _, err := self.gitClient.Projects.GetProject(id, nil)
		if err != nil {
			self.Error("未找到服务: %v, %s", id, err)
			continue
		}
		urls = append(urls, project.WebURL)
	}

	for _, name := range names {
		if len(name) == 0 {
			continue
		}
		projects, _, err := self.gitClient.Projects.ListProjects(&gitlab.ListProjectsOptions{
			Search: gitlab.String(name),
		})
		if err != nil {
			self.Error("未匹配到服务: %v, %v", name, err)
			continue
		}
		if len(projects) == 0 {
			self.Error("未匹配到服务: %v", name)
			continue
		}
		if len(projects) > 1 {
			self.Error("%v 匹配到多个服务服务:", name)
			for _, project := range projects {
				self.Error("%d %s", project.ID, project.WebURL)
			}
			continue
		}
		urls = append(urls, projects[0].WebURL)
	}

	gitlabDir := os.Getenv("GITLAB")
	if len(gitlabDir) == 0 {
		gitlabDir = "./"
	}
	for _, url := range urls {
		projectDir := gitlabDir + strings.TrimPrefix(url, "https://gitlab.kaiqitech.com/")
		pathx.MkdirIfNotExist(projectDir)
		if info, err := execx.Run(fmt.Sprintf("git clone %s.git %s", url, projectDir), "./"); err != nil {
			self.Error(err.Error())
			continue
		} else {
			self.Info(info)
		}
		self.Info("拉取项目: %v", projectDir)
	}
}

func (self *Git) DeleteProject(projectIds ...int) error {
	for _, projectId := range projectIds {
		if _, err := self.gitClient.Projects.DeleteProject(projectId); err != nil {
			return err
		}
		self.Info("删除项目: %d", projectId)
	}
	return nil
}
