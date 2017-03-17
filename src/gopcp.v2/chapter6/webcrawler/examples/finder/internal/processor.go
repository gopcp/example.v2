package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopcp.v2/chapter6/webcrawler/module"
)

// genItemProcessors 用于生成条目处理器。
func genItemProcessors(dirPath string) []module.ProcessItem {
	savePicture := func(item module.Item) (result module.Item, err error) {
		if item == nil {
			return nil, errors.New("invalid item!")
		}
		// 检查和准备数据。
		var absDirPath string
		if absDirPath, err = checkDirPath(dirPath); err != nil {
			return
		}
		v := item["reader"]
		reader, ok := v.(io.Reader)
		if !ok {
			return nil, fmt.Errorf("incorrect reader type: %T", v)
		}
		readCloser, ok := reader.(io.ReadCloser)
		if ok {
			defer readCloser.Close()
		}
		v = item["name"]
		name, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("incorrect name type: %T", v)
		}
		// 创建图片文件。
		fileName := name
		filePath := filepath.Join(absDirPath, fileName)
		file, err := os.Create(filePath)
		if err != nil {
			return nil, fmt.Errorf("couldn't create file: %s (path: %s)",
				err, filePath)
		}
		defer file.Close()
		// 写图片文件。
		_, err = io.Copy(file, reader)
		if err != nil {
			return nil, err
		}
		// 生成新的条目。
		result = make(map[string]interface{})
		for k, v := range item {
			result[k] = v
		}
		result["file_path"] = filePath
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}
		result["file_size"] = fileInfo.Size()
		return result, nil
	}
	recordPicture := func(item module.Item) (result module.Item, err error) {
		v := item["file_path"]
		path, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("incorrect file path type: %T", v)
		}
		v = item["file_size"]
		size, ok := v.(int64)
		if !ok {
			return nil, fmt.Errorf("incorrect file name type: %T", v)
		}
		logger.Infof("Saved file: %s, size: %d byte(s).", path, size)
		return nil, nil
	}
	return []module.ProcessItem{savePicture, recordPicture}
}
