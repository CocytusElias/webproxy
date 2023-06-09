package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"webProxy/extern/logger"
)

func init() {
	logger.Init("client-module-gen.go")
}

func main() {
	var modulesPath string
	var err error

	var moduleNum = 0
	var templateImportStr = ""
	var templateInitStr = ""

	// 获取模块目录的绝对路径
	if modulesPath, err = filepath.Abs("client/module"); err != nil {
		logger.Panic(err.Error())
	}

	// 遍历模块目录并提取模块
	_ = filepath.Walk(modulesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error(err.Error())
			return nil
		}

		// 检查当前路径是否为目录
		if info.IsDir() && path != modulesPath {
			var needLoad bool
			if needLoad, err = hasGoFiles(path); err != nil {
				logger.Error(err.Error())
				return nil
			}

			// 仅带有 go 文件的目录才需要加载
			if !needLoad {
				return nil
			}

			moduleNum++
			templateImportStr += fmt.Sprintf("\t\"webProxy/client/module/%s\"\n", info.Name())
			templateInitStr += fmt.Sprintf("\n\tif modules[\"%s\"], err = %s.Init(); err != nil {\n\t\tlogger.Panic(err.Error())\n\t}\n", info.Name(), info.Name())
		}

		return nil
	})

	if moduleNum == 0 {
		logger.Warn("no modules to load")
	} else {
		gen(templateImportStr, templateInitStr)
		logger.Info(fmt.Sprintf("modules loaded success. module num: %v", moduleNum))
	}

}

// 判断是否存在 go 文件
func hasGoFiles(dirPath string) (bool, error) {
	var dir *os.File
	var err error

	// 读取目录内容
	if dir, err = os.Open(dirPath); err != nil {
		logger.Error(err.Error())
		return false, err
	}
	defer func(dir *os.File) {
		if err = dir.Close(); err != nil {
			logger.Error(err.Error())
		}
	}(dir)

	// 遍历目录中的文件
	var files []os.FileInfo
	if files, err = dir.Readdir(0); err != nil {
		logger.Error(err.Error())
		return false, err
	}

	// 只要一个 go 文件就返回
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") {
			// 目录中存在 Go 文件
			return true, nil
		}
	}

	// 一个 go 文件都没有
	return false, nil
}

// 生成
func gen(templateImportStr, templateInitStr string) {
	var template []byte
	var err error

	if //goland:noinspection ALL
	template, err = ioutil.ReadFile("client/module/module.template"); err != nil {
		logger.Panic(err.Error())
	}

	processedContent := string(template)

	processedContent = strings.ReplaceAll(processedContent, "%{{import}}%", templateImportStr)

	processedContent = strings.ReplaceAll(processedContent, "%{{init}}%", fmt.Sprintf("\n\tvar err error\n%s", templateInitStr))

	template = []byte(processedContent)
	if template, err = format.Source(template); err != nil {
		logger.Panic(err.Error())
	}

	if //goland:noinspection ALL
	err = ioutil.WriteFile("client/module/module.go", template, 0644); err != nil {
		logger.Panic(err.Error())
	}
}
