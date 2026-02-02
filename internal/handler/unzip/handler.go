package unzip

import (
	"converter/internal/service/unzip"
	"errors"
	"fmt"
)

type UnzipHandler struct {
	s *unzip.UnzipService
}

func NewUnzipHandler(unzipSerivce *unzip.UnzipService) *UnzipHandler {
	return &UnzipHandler{
		s: unzipSerivce,
	}
}

func (uh *UnzipHandler) Unzip(root string) (string, error) {
	if root == "" {
		return "", errors.New("cant unpack archive with empty filepath")
	}

	result, err := uh.s.UnzipArchive(root)
	if err != nil {
		return "", fmt.Errorf("unzip extract files err: %s", err)
	}

	return result.DestDir, nil
}
