package server

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Note: Changing this might need tweeking of the format string
const maxFilenameLength = 30

func (s *Server) fileListFormatter(hfm HashFileMap, showTime bool) string {
	baseURI := s.config.URL
	res := make([]string, 0, len(hfm))

	for hash, file := range hfm {
		name := file.Filename
		if len(name) > maxFilenameLength {
			name = name[:maxFilenameLength-2] + ".."
		}

		if showTime {
			humanTime := time.Unix(file.CreatedAt, 0).Format("2006-01-02 15:04")
			res = append(res, fmt.Sprintf("%-30s \t%10s  %18s  %s%s", name, file.HumanSize(), humanTime, baseURI, hash))
		} else {
			res = append(res, fmt.Sprintf("%-30s \t%10s  %s%s", name, file.HumanSize(), baseURI, hash))
		}

	}

	// sort result by original filename, alphabetically
	sort.Strings(res)

	return strings.Join(res, "\n")
}
