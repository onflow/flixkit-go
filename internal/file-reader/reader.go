package filereader

import "os"

type DefaultFileReader struct{}

func GetDefaultFileReader() DefaultFileReader {
	return DefaultFileReader{}
}

func (f DefaultFileReader) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
