/**
 * @Author: cyj19
 * @Date: 2021/12/3 16:42
 */

package global

import (
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var (
	RootDir string
	once    = new(sync.Once)
)

func Init() {
	once.Do(func() {
		interRootDir()
	})
}

// 推断出项目根目录
func interRootDir() {
	// 返回一个对应当前工作目录的根路径
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var inter func(d string) string
	inter = func(d string) string {
		// 确保项目根目录下存在template目录
		if exists(d + "/template") {
			return d
		}

		// 递归查找 Dir返回路径除去最后一个路径元素的部分，即该路径最后一个元素所在的目录(上一级)
		return inter(filepath.Dir(d))
	}

	RootDir = inter(cwd)
}

func exists(filename string) bool {
	// Stat返回一个描述filename指定的文件对象的FileInfo
	_, err := os.Stat(filename)
	// 返回一个布尔值说明该错误是否表示一个文件或目录已经存在
	return err == nil || os.IsExist(err)
}
