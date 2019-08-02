package build

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	portainer "github.com/hazik1024/portainer/api"
	"github.com/hazik1024/portainer/api/bolt/registry"
	"github.com/hazik1024/portainer/api/custom/mysqldb"
	httperror "github.com/portainer/libhttp/error"
)

const (
	sourcePath    = "/data/build/sourceCode"
	imagePath     = "/data/build/buildImage"
	sourcePathWin = "C:\\data\\build\\sourceCode"
	imagePathWin  = "C:\\data\\build\\buildImage"
)

// Service 定义BuildService
type Service struct {
	db              *mysqldb.MySQLDb
	registryService *registry.Service
}

// NewService 返回BuildService指针
func NewService(db *mysqldb.MySQLDb, registryService *registry.Service) *Service {
	initDir()
	return &Service{
		db:              db,
		registryService: registryService,
	}
}

func initDir() {
	if !pathExists(sourcePath) {
		os.MkdirAll(sourcePath, os.ModePerm)
	}
	if !pathExists(imagePath) {
		os.MkdirAll(imagePath, os.ModePerm)
	}
}

// BuildAndPushImage 编译并推送
func (s *Service) buildAndPushImage(req reqPayload) ([6]string, *httperror.HandlerError) {
	gitPaths := strings.Split(req.GitPath, "/")
	projectDir := strings.Replace(gitPaths[len(gitPaths)-1], ".git", "", 1)
	cloneOutStr, err1 := s.cloneSourceCode(req, projectDir)
	var outStrs [6]string
	if err1 != nil {
		return outStrs, err1
	}
	outStrs[0] = "1.拉取分支代码到本地:\n"
	outStrs[1] = cloneOutStr
	tag := req.ImageTag
	if len(tag) < 1 {
		str := cloneOutStr[:len(cloneOutStr)-2]
		arr := strings.Split(str, "\n")
		tag = arr[len(arr)-1]
	}
	log.Println("tag: ", tag)

	registry, err := s.registryService.Registry(portainer.RegistryID(req.RegistryID))
	if err != nil {
		return outStrs, &httperror.HandlerError{
			StatusCode: http.StatusBadRequest,
			Message:    "找不到镜像仓库",
			Err:        err,
		}
	}

	var buffer bytes.Buffer
	buffer.WriteString(req.ImageName)
	buffer.WriteString(":")
	buffer.WriteString(req.GitBranch)
	if len(tag) > 0 {
		buffer.WriteString("-")
		buffer.WriteString(tag)
	}
	imageName := buffer.String()

	buildOutStr, err2 := s.buildImage(req, imageName, projectDir, registry)
	if err2 != nil {
		return outStrs, err2
	}
	outStrs[2] = "2.构建docker镜像:\n"
	outStrs[3] = buildOutStr

	pushOutStr, err3 := s.pushImage(req, imageName, projectDir, registry)
	if err3 != nil {
		return outStrs, err3
	}
	outStrs[2] = "3.推送镜像到镜像仓库:\n"
	outStrs[3] = pushOutStr
	return outStrs, nil
}

// cloneSourceCode 拷贝源码到本地
func (s *Service) cloneSourceCode(req reqPayload, projectDir string) (string, *httperror.HandlerError) {
	var gitPath string
	var template string
	if strings.Contains(req.GitPath, "http://") {
		gitPath = strings.Replace(req.GitPath, "http://", "", 1)
		template = "http://%s:%s@%s"
	} else if strings.Contains(req.GitPath, "https://") {
		gitPath = strings.Replace(req.GitPath, "https://", "", 1)
		template = "https://%s:%s@%s"
	} else {
		return "", &httperror.HandlerError{
			StatusCode: http.StatusBadRequest,
			Message:    "GitPath应为http/https开头",
			Err:        nil,
		}
	}
	gitPath = fmt.Sprintf(template, req.GitUser, req.GitPwd, gitPath)
	var buffer bytes.Buffer
	buffer.WriteString("cd ")
	buffer.WriteString(sourcePath)
	buffer.WriteString(" && rm -rf ")
	buffer.WriteString(projectDir)
	buffer.WriteString(" && git clone '")
	buffer.WriteString(gitPath)
	buffer.WriteString("' 2>&1 && cd ")
	buffer.WriteString(projectDir)
	buffer.WriteString(" && git checkout ")
	buffer.WriteString(req.GitBranch)
	buffer.WriteString(" && git symbolic-ref --short -q HEAD 2>&1 && git rev-parse --short HEAD 2>&1 ")
	// cmd := fmt.Sprint("cd %s && rm -rf %s && git clone '%s' 2>&1 && cd %s && git checkout %s && git symbolic-ref --short -q HEAD 2>&1 && git rev-parse --short HEAD 2>&1", sourcePath, projectDir, gitPath, projectDir, req.GitBranch)
	cmd := buffer.String()
	str, err := execShell(cmd)
	log.Println("执行命令:\n", cmd, "\n返回内容:\n", str)
	if err != nil {
		log.Println("执行出错:", err.Error())
		return "", &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "复制源码出错",
			Err:        err,
		}
	}
	return str, nil
}

// buildImage 构建Docker镜像
func (s *Service) buildImage(req reqPayload, imageName string, projectDir string, registry *portainer.Registry) (string, *httperror.HandlerError) {
	var buffer bytes.Buffer
	buffer.WriteString("cd ")
	buffer.WriteString(sourcePath)
	buffer.WriteString("/")
	buffer.WriteString(projectDir)
	buffer.WriteString(" && docker build -t '")
	buffer.WriteString(registry.URL)
	buffer.WriteString("/")
	buffer.WriteString(imageName)
	buffer.WriteString("' . 2>&1")

	cmd := buffer.String()
	str, err := execShell(cmd)
	log.Println("执行命令:\n", cmd, "\n返回内容:\n", str)
	if err != nil {
		log.Println("执行出错:", err.Error())
		return "", &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "构建Docker镜像出错",
			Err:        err,
		}
	}

	return str, nil
}

// pushImage 推送镜像
func (s *Service) pushImage(req reqPayload, imageName string, projectDir string, registry *portainer.Registry) (string, *httperror.HandlerError) {
	var buffer bytes.Buffer
	buffer.WriteString("cd ")
	buffer.WriteString(sourcePath)
	buffer.WriteString("/")
	buffer.WriteString(projectDir)
	buffer.WriteString(" && docker login -u ")
	buffer.WriteString(registry.Username)
	buffer.WriteString(" -p ")
	buffer.WriteString(registry.Password)
	buffer.WriteString(" ")
	buffer.WriteString(registry.URL)
	buffer.WriteString(" 2>&1 && docker push '")
	buffer.WriteString(imageName)
	buffer.WriteString("' 2>&1 && docker tag '")
	buffer.WriteString(imageName)
	buffer.WriteString("' '")
	buffer.WriteString(req.ImageName)
	buffer.WriteString(":latest' 2>&1 && docker push '")
	buffer.WriteString(req.ImageName)
	buffer.WriteString("' 2>&1 && docker rmi '")
	buffer.WriteString(imageName)
	buffer.WriteString("' && docker rmi '")
	buffer.WriteString(req.ImageName)
	buffer.WriteString(":latest'")

	cmd := buffer.String()
	str, err := execShell(cmd)
	log.Println("执行命令:\n", cmd, "\n返回内容:\n", str)
	if err != nil {
		log.Println("执行出错:", err.Error())
		return "", &httperror.HandlerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "推送镜像出错",
			Err:        err,
		}
	}

	return str, nil
}

// PathExists 判断路径是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// exec_shell 阻塞式调用shell脚本
func execShell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out
	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	return out.String(), err
}
