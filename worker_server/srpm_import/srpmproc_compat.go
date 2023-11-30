package srpm_import

import (
	"github.com/openela/mothership/base/storage"
	"github.com/rocky-linux/srpmproc/pkg/data"
	"github.com/rocky-linux/srpmproc/pkg/misc"
	"strings"
)

type srpmprocBlobCompat struct {
	storage.Storage
}

func (s *srpmprocBlobCompat) Write(path string, content []byte) error {
	_, err := s.PutBytes(path, content)
	return err
}

func (s *srpmprocBlobCompat) Read(path string) ([]byte, error) {
	return s.Get(path)
}

type srpmprocImportModeCompat struct {
	data.ImportMode
}

func (s *srpmprocImportModeCompat) ImportName(pd *data.ProcessData, md *data.ModeData) string {
	if misc.GetTagImportRegex(pd).MatchString(md.TagBranch) {
		match := misc.GetTagImportRegex(pd).FindStringSubmatch(md.TagBranch)
		return match[3]
	}

	return strings.Replace(strings.TrimPrefix(md.TagBranch, "refs/heads/"), "%", "_", -1)
}
