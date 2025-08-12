package uio

import (
	"bufio"
	"fmt"
	uos "github.com/hxxshidage/myutils/os"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func ReadAll(path string) ([]byte, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

type LineParser[R any] func(string) (R, error)

func ReadLinesAndParse[R any](path string, lp LineParser[R]) ([]R, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var lines []R
	var line string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line = scanner.Text()
		if len(line) == 0 {
			continue
		}

		if parsed, err := lp(line); err != nil {
			return nil, err
		} else {
			lines = append(lines, parsed)
		}
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(name string, contents []byte) error {
	return os.WriteFile(name, contents, 0644)
}

func WritePath(path string, contents []byte) error {
	if err := mkdirp(path); err != nil {
		return err
	}

	return WriteFile(path, contents)
}

func AppendFile(name string, contents []byte) error {
	f, err := os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = f.Write(contents)

	if cErr := f.Close(); cErr != nil && err == nil {
		err = cErr
	}
	return err
}

func AppendPath(path string, contents []byte) error {
	if err := mkdirp(path); err != nil {
		return err
	}

	return AppendFile(path, contents)
}

func WriteFileWithStream(name string, rd io.Reader) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer f.Close()

	// 缓冲区
	writer := bufio.NewWriterSize(f, 64*1024)
	_, err = io.Copy(writer, rd)
	if err != nil {
		return err
	}

	// 确保缓冲区内容写入磁盘
	if err = writer.Flush(); err != nil {
		return err
	}

	return err
}

func WritePathWithStream(path string, rd io.Reader) error {
	if err := mkdirp(path); err != nil {
		return err
	}

	return WriteFileWithStream(path, rd)
}

// 分块写入: 适用于contents比较大, 但不能太大, 太大考虑stream写入
func WriteFileChunk(name string, contents []byte, chunkSize int) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	defer f.Close()

	if chunkSize == 0 {
		// 默认4MB
		chunkSize = 4 * 1024 * 1024
	}

	// 缓冲区
	writer := bufio.NewWriterSize(f, 64*1024)

	cl := len(contents)

	for i := 0; i < cl; i += chunkSize {
		end := i + chunkSize
		if end > cl {
			end = cl
		}
		_, err = writer.Write(contents[i:end])
		if err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	return nil
}

func WritePathChunk(path string, contents []byte, chunkSize int) error {
	if err := mkdirp(path); err != nil {
		return err
	}

	return WriteFileChunk(path, contents, chunkSize)
}

func mkdirp(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0755)
}

// windows路径风格转linux风格, 支持网络路径, 盘符路径
func ToUnixPath(path string) string {
	if path == "" {
		return path
	}

	if uos.PlatformWin() {
		if strings.HasPrefix(path, `\\`) {
			return strings.ReplaceAll(path, `\`, `/`)
		}

		if len(path) >= 2 && path[1] == ':' {
			drive := unicode.ToLower(rune(path[0]))

			rest := path[2:]
			rest = strings.ReplaceAll(rest, `\`, `/`)

			rest = strings.Trim(rest, `/`)

			if rest == "" {
				return fmt.Sprintf("/%c", drive)
			}
			return fmt.Sprintf("/%c/%s", drive, rest)
		}

		return strings.ReplaceAll(path, `\`, `/`)
	}

	return path
}
