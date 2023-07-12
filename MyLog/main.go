package main

import (
	"fmt"
	"io"
	"os"
)

type Reader interface {
	Read(os.File) ([]byte, error)
	ReadString(os.File) (string, error)
}

type Writer interface {
	Write(os.File, []byte) error
	WriteString(os.File, string) error
}

type FileLog struct {
}

func (f *FileLog) NewLog() Reader {
	panic("implement me")
}

func (f *FileLog) Read(file os.File) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileLog) ReadString(file os.File) (string, error) {
	//TODO implement me
	panic("implement me")
}

type TerminalLog struct {
}

func (t *TerminalLog) Write(file os.File, bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (t *TerminalLog) WriteString(file os.File, s string) error {
	//TODO implement me
	panic("implement me")
}

// 日志分级
// userScan
// logger.Trace()
// logger.Debug()
// logger.Warning()
// logger.Info()
// logger.Error(日志内容)

func main() {
	file, err := os.Open(".\\main.go")
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println("close file error:", err)
			return
		}
	}()

	tmp := make([]byte, 1024)
	_, err = file.Read(tmp)
	if err != nil && err != io.EOF {
		fmt.Println("read file error:", err)
		return
	}
}
