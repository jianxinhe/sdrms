package models

import (
	"container/list"
	"io/ioutil"
)

// GetTreeGrid 获取treegrid顺序的列表
func GetTreeGrid() []*Resource {
	list, err := TraverseFolder()
	if err != nil {
		return []*Resource{}
	}
	return resourceList2TreeGrid(list)
}

// TraverseFolder 获取文件列表
func TraverseFolder() ([]*Resource, error) {
	//storagePath := beego.AppConfig.String("course_ware_storage_path")
	storagePath := "E:/学院相关/师母项目/coursesuite/"

	return getFileList(storagePath)
}

// getFileList 递归遍历文件夹和其中的文件
func getFileList(path string) ([]*Resource, error) {
	// 最终存储数据的list
	result := make([]*Resource, 0)
	fs, _ := ioutil.ReadDir(path)
	queue := list.New()
	id := 1
	parent := &Resource{
		Id:     id,
		Name:   "",
		Parent: nil,
		Rtype:  -1,
		Icon:   "",
		UrlFor: "",
	}
	// 先进行根目录的遍历
	for _, file := range fs {
		if file.IsDir() {
			id += 1
			rs := &Resource{
				Id:     id,
				Name:   file.Name(),
				Parent: nil,
				Rtype:  parent.Rtype + 1,
				Icon:   "",
				UrlFor: path + file.Name() + "/",
			}
			// 把路径加入到队列当中
			queue.PushBack(rs)
			// 把路径加入到result当中
			result = append(result, rs)
		} else {
			// 如果本身就是文件不是路径的话
			id += 1
			rs := &Resource{
				Id:     id,
				Name:   file.Name(),
				Parent: nil,
				Rtype:  parent.Rtype + 1,
				Icon:   "",
				UrlFor: path + file.Name(),
			}
			// 直接把文件加入到result当中
			result = append(result, rs)
		}
	}

	// 对队列进行遍历
	for queue.Len() > 0 {
		// 出队一个元素
		i1 := queue.Front()
		queue.Remove(i1)
		//newPath := path + "/" + file.Name() + "/"
		fs, _ := ioutil.ReadDir(i1.Value.(*Resource).UrlFor)
		parent := i1.Value.(*Resource)
		for _, file := range fs {
			if file.IsDir() {
				id += 1
				rs := &Resource{
					Id:     id,
					Name:   file.Name(),
					Parent: parent,
					Rtype:  parent.Rtype + 1,
					Icon:   "",
					UrlFor: parent.UrlFor + file.Name() + "/",
				}
				// 把路径加入到队列当中
				queue.PushBack(rs)
				// 把路径加入到result当中
				result = append(result, rs)
			} else {
				// 如果本身就是文件不是路径的话
				id += 1
				rs := &Resource{
					Id:     id,
					Name:   file.Name(),
					Parent: parent,
					Rtype:  parent.Rtype + 1,
					Icon:   "",
					UrlFor: parent.UrlFor + file.Name(),
				}
				// 直接把文件加入到result当中
				result = append(result, rs)
			}
		}
	}

	return result, nil
}
