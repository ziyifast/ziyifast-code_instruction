package main

import "fmt"

const Separator = "--"

// FileSystemNode 文件系统接口：文件和目录都要实现该接口
type FileSystemNode interface {
	Display(separator string)
}

// FileCommonFunc 文件通用功能
type FileCommonFunc struct {
	fileName string
}

func (f *FileCommonFunc) SetFileName(fileName string) {
	f.fileName = fileName
}

// FileNode 文件类
type FileNode struct {
	FileCommonFunc
}

// Display 文件类显示文件内容
func (f *FileNode) Display(separator string) {
	fmt.Println(separator + f.fileName + "   文件内容为：Hello，world")
}

// DirectoryNode 目录类
type DirectoryNode struct {
	FileCommonFunc
	nodes []FileSystemNode
}

// Display 目录类展示文件名
func (d *DirectoryNode) Display(separator string) {
	fmt.Println(separator + d.fileName)
	for _, node := range d.nodes {
		node.Display(separator + Separator)
	}
}

func (d *DirectoryNode) Add(f FileSystemNode) {
	d.nodes = append(d.nodes, f)
}

func main() {
	//初始化
	biji := DirectoryNode{}
	biji.SetFileName("笔记")

	huiyi := DirectoryNode{}
	huiyi.SetFileName("会议")

	chenhui := FileNode{}
	chenhui.SetFileName("晨会.md")

	zhouhui := FileNode{}
	zhouhui.SetFileName("周会.md")
	//组装
	biji.Add(&huiyi)
	huiyi.Add(&chenhui)
	huiyi.Add(&zhouhui)
	//显示
	biji.Display(Separator)
}
