package mod

import (
	"archive/zip"
	"fmt"
	"github.com/guestin/mob/merrors"
	"github.com/guestin/mob/mio"
	"os"
	"testing"
)

func TestUnZip(t *testing.T) {
	err := os.Mkdir("zip_output", 0755)
	merrors.AssertError(err, "mkdir zip_output")
	defer func() {
		_ = os.RemoveAll("zip_output")
	}()
	err = UnZip("sample.zip", "zip_output")
	merrors.AssertError(err, "unzip")
}

func TestZip(t *testing.T) {
	wd, err := os.Getwd()
	merrors.AssertError(err, "getwd")
	fmt.Println("wd=", wd)
	file, err := os.OpenFile("output.zip", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	merrors.AssertError(err, "create output")
	defer mio.CloseIgnoreErr(file)
	err = Zip(file, []*FileGroup{
		NewFileGroup(".", "go.mod").SetCustomDirName("c").AddFileItemWithAlias("go.sum", "goSum"),
		NewFileGroup(".idea", "misc.xml", "modules.xml", "workspace.xml").SetCustomDirName(".").
			AddChildFileGroup(NewFileGroup("inspectionProfiles", "Project_Default.xml").SetCustomDirName("a")),
	})
	merrors.AssertError(err, "zip failed")
}

func TestZip2(t *testing.T) {
	file, err := os.OpenFile("output.zip", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	merrors.AssertError(err, "create output")

	writer := zip.NewWriter(file)
	defer func() {
		mio.CloseIgnoreErr(writer)
		mio.CloseIgnoreErr(file)
	}()
	_, err = writer.Create("a/")
	merrors.AssertError(err, "create a/")
	_, err = writer.Create("a/b/")
	merrors.AssertError(err, "create a/b/")
	_, err = writer.Create("a/")
	merrors.AssertError(err, "create a/ again")
}
