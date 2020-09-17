package mob

import (
	"archive/zip"
	"fmt"
	"github.com/guestin/mob/mio"
	"io"
	"os"
	"path/filepath"
)

type zipFileItem struct {
	rawFileName   string
	inZipFileName string
}

type FileGroup struct {
	rawDir          string
	dirName         string
	childFiles      []*zipFileItem
	childFileGroups []*FileGroup
}

func (this *FileGroup) GetDirName() string {
	return this.dirName
}

func (this *FileGroup) SetCustomDirName(dirName string) *FileGroup {
	this.dirName = dirName + separatorStr
	return this
}

const separatorStr = string(filepath.Separator)

func NewFileGroup(dir string, childFiles ...string) *FileGroup {
	out := &FileGroup{
		rawDir:     dir,
		dirName:    "",
		childFiles: nil,
	}
	if len(dir) != 0 {
		out.dirName = fmt.Sprintf("%s%c", filepath.Base(dir), filepath.Separator)
	}
	if len(childFiles) == 0 {
		return out
	}
	out.childFiles = make([]*zipFileItem, 0, len(childFiles))
	for idx := range childFiles {
		_ = out.AddFileItem(childFiles[idx])
	}
	return out
}

func (this *FileGroup) addChildFileGroup(fg FileGroup) {
	fg.rawDir = filepath.Join(this.rawDir, filepath.Base(fg.rawDir))
	fg.dirName = filepath.Join(this.dirName, fg.dirName) + separatorStr
	this.childFileGroups = append(this.childFileGroups, &fg)
}

func (this *FileGroup) AddChildFileGroup(fgs ...*FileGroup) *FileGroup {
	if this.childFileGroups == nil {
		this.childFileGroups = make([]*FileGroup, 0, len(fgs))
	}
	for _, fg := range fgs {
		this.addChildFileGroup(*fg)
	}
	return this
}

func (this *FileGroup) AddFileItem(file string) *FileGroup {
	baseName := filepath.Base(file)
	this.childFiles = append(this.childFiles, &zipFileItem{
		rawFileName:   baseName,
		inZipFileName: baseName,
	})
	return this
}

func (this *FileGroup) AddFileItemWithAlias(file, alias string) *FileGroup {
	baseName := filepath.Base(file)
	this.childFiles = append(this.childFiles, &zipFileItem{
		rawFileName:   baseName,
		inZipFileName: alias,
	})
	return this
}

func zipFilesDetails(zipWriter *zip.Writer, fgs []*FileGroup) error {
	curDir := "."
	for _, fg := range fgs {
		if len(fg.dirName) != 0 {
			_, err := zipWriter.Create(fg.dirName)
			if err != nil {
				return err
			}
			curDir = fg.dirName
		} else {
			curDir = "."
		}
		for _, it := range fg.childFiles {
			writer, err := zipWriter.Create(filepath.Join(curDir, it.inZipFileName))
			if err != nil {
				return err
			}
			file, err := os.Open(filepath.Join(fg.rawDir, it.rawFileName))
			if err != nil {
				return err
			}
			_, err = io.Copy(writer, file)
			if err != nil {
				_ = file.Close()
				return err
			}
			_ = file.Close()
		}
		if len(fg.childFileGroups) == 0 {
			continue
		}
		err := zipFilesDetails(zipWriter, fg.childFileGroups)
		if err != nil {
			return err
		}
	}
	return nil
}

func Zip(output io.Writer, fgs []*FileGroup) error {
	zipWriter := zip.NewWriter(output)
	defer mio.CloseIgnoreErr(zipWriter) // implicit flush to file
	return zipFilesDetails(zipWriter, fgs)
}

func UnZip(zipFilePath, output string) error {
	zrc, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer mio.CloseIgnoreErr(zrc)
	for _, f := range zrc.File {
		info := f.FileInfo()
		mode := info.Mode()
		outputTarget := filepath.Join(output, f.Name)
		if info.IsDir() {
			err := os.MkdirAll(outputTarget, mode|0755)
			if err != nil {
				return err
			}
			continue
		}
		file, err := os.OpenFile(
			outputTarget,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			mode|0766)
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(file, rc)
		_ = rc.Close()
		_ = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
