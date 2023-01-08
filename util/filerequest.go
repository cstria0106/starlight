package util

import (
	"encoding/binary"
	"io"
)

type FileRequest struct {
	Source, SourceOffset, CompressedSize int64
}

func (f *FileRequest) Write(w io.Writer) error {
	var err error
	if err = binary.Write(w, binary.BigEndian, f.Source); err != nil {
		return err
	}
	if err = binary.Write(w, binary.BigEndian, f.SourceOffset); err != nil {
		return err
	}
	if err = binary.Write(w, binary.BigEndian, f.CompressedSize); err != nil {
		return err
	}
	return nil
}

func ReadFileRequest(r io.Reader, fr *FileRequest) error {
	var err error
	if err = binary.Read(r, binary.BigEndian, &fr.Source); err != nil {
		return err
	}
	if err = binary.Read(r, binary.BigEndian, &fr.SourceOffset); err != nil {
		return err
	}
	if err = binary.Read(r, binary.BigEndian, &fr.CompressedSize); err != nil {
		return err
	}
	return nil
}
