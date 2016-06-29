package compressor

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

type ZipCompressor struct{}

func NewZipCompressor() *ZipCompressor {
	return &ZipCompressor{}
}

func (c *ZipCompressor) Compress(compressedFile io.Writer, targetDir string, files []string) error {
	w := zip.NewWriter(compressedFile)

	for _, filename := range files {
		filepath := fmt.Sprintf("%s/%s", targetDir, filename)
		info, err := os.Stat(filepath)
		if err != nil {
			return err
		}

		if info.IsDir() {
			continue
		}

		file, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		hdr.Name = filename

		local := time.Now().Local()

		//現時刻のオフセットを取得
		_, offset := local.Zone()

		//差分を追加
		hdr.SetModTime(hdr.ModTime().Add(time.Duration(offset) * time.Second))

		f, err := w.CreateHeader(hdr)
		if err != nil {
			return err
		}

		contents, _ := ioutil.ReadFile(filepath)
		_, err = f.Write(contents)
		if err != nil {
			return err
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}
