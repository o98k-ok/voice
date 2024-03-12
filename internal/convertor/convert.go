package convertor

import (
	"fmt"
	"io"
	"os"

	"github.com/duke-git/lancet/v2/system"
)

type Convertor interface {
	ConvertM4AToWav(reader io.Reader, writer io.Writer) error
}

type AfconvertConvertor struct {
	WorkDir string
}

func NewAfconvertConvertor(workDir string) *AfconvertConvertor {
	return &AfconvertConvertor{
		WorkDir: workDir,
	}
}

func (ac *AfconvertConvertor) ConvertM4AToWav(reader io.Reader, writer io.Writer) error {
	src, err := os.CreateTemp(ac.WorkDir, "afconvert_input_*.mp4")
	if err != nil {
		return err
	}
	defer src.Close()
	defer os.Remove(src.Name())
	if _, err = io.Copy(src, reader); err != nil {
		return err
	}

	dstname := func() string {
		dst, err := os.CreateTemp(ac.WorkDir, "afconvert_output_*.wav")
		if err != nil {
			return ""
		}
		dst.Close()
		return dst.Name()
	}()
	defer os.Remove(dstname)

	// UI8 目前用的数据格式，有些数据格式没办法快进
	c := fmt.Sprintf("afconvert %s -f WAVE  -d UI8 %s", src.Name(), dstname)
	if _, _, err := system.ExecCommand(c); err != nil {
		return err
	}

	dst, err := os.Open(dstname)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(writer, dst); err != nil {
		return err
	}
	return nil
}
