package clide

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

func DefaultEditor() (string, error) {
	editorPath := os.Getenv("EDITOR")
	if editorPath == "" {
		switch runtime.GOOS {
		case "windows":
			return "C:\\WINDOWS\\system32\\notepad.exe", nil
		case "linux":
			return "/usr/bin/vi", nil
		case "darwin":
			return "/usr/bin/vi", nil
		}
	}

	return "", errors.New("unrecognized operating system")
}

type EditorOptions struct {
	Editor   string
	FilePath string
}

func NewEditor(opts EditorOptions) (*Editor, error) {
	var editor string
	var err error
	if opts.Editor == "" {
		editor, err = DefaultEditor()
		if err != nil {
			return nil, err
		}
	} else {
		editor = opts.Editor
	}

	fPath := opts.FilePath
	var fd *os.File
	if fPath == "" {
		fd, err = ioutil.TempFile("/tmp", "clide")
		if err != nil {
			return nil, fmt.Errorf("failed to create temporary file: %+w", err)
		}
	} else {
		fd, err = os.Open(fPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file \"%s\": %+w", fPath, err)
		}
	}

	return &Editor{
		fd:       fd,
		filePath: fPath,
		editor:   editor,
	}, nil
}

type Editor struct {
	filePath string
	fd       *os.File
	editor   string
}

func (e *Editor) ReadAll() ([]byte, error) {
	_, err := e.fd.Seek(0, 0)
	if err != nil {
		return []byte(""), fmt.Errorf("failed to seek to beginning of file: %+w", err)
	}

	return ioutil.ReadAll(e.fd)
}

func (e *Editor) Close() error {
	err := e.fd.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %+w", err)
	}

	return os.Remove(e.filePath)
}

func (e *Editor) Run() error {
	cmd := exec.Command(e.editor, e.filePath)
	err := cmd.Run()

	return err
}
