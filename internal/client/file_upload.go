package client

import "path/filepath"

func analyzeLocalFile(filename string) (*InternalClientMessage, error) {

	// TODO: Clean this properly
	_, fileOnly := filepath.Split(filename)

	return &InternalClientMessage{
		LocalFilename: filename,
		CleanFilename: fileOnly,
	}, nil
}
