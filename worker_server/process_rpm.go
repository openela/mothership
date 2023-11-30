package mothership_worker_server

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	mothershippb "github.com/openela/mothership/proto/v1"
	"github.com/openela/mothership/worker_server/srpm_import"
	"github.com/pkg/errors"
	"github.com/sassoftware/go-rpmutils"
	"go.temporal.io/sdk/temporal"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// VerifyResourceExists verifies that the resource exists.
// This is a Temporal activity.
func (w *Worker) VerifyResourceExists(uri string) error {
	canRead, err := w.storage.CanReadURI(uri)
	if err != nil {
		return errors.Wrap(err, "failed to check if resource URI can be read")
	}

	if !canRead {
		return temporal.NewNonRetryableApplicationError(
			"cannot read resource URI",
			"cannotReadResourceURI",
			errors.New("client submitted a resource URI that cannot be read by server"),
		)
	}

	object, err := getObjectPath(uri)
	if err != nil {
		return err
	}

	exists, err := w.storage.Exists(object)
	if err != nil {
		return errors.Wrap(err, "failed to check if resource exists")
	}

	if !exists {
		// Since the client can trigger this activity before uploading the resource,
		// we should not return a non-retryable error.
		// The parent workflow should handle the retry arrangements up to 2 hours
		// per the spec.
		return errors.New("resource does not exist")
	}

	return nil
}

// ImportRPM imports an RPM into the database.
// This is a Temporal activity.
func (w *Worker) ImportRPM(uri string, checksumSha256 string, osRelease string) (*mothershippb.ImportRPMResponse, error) {
	tempDir, err := os.MkdirTemp("", "mothership-worker-server-import-rpm-*")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	// Parse uri
	object, err := getObjectPath(uri)
	if err != nil {
		return nil, err
	}

	// Download the resource to the temporary directory
	err = w.storage.Download(object, filepath.Join(tempDir, "resource.rpm"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to download resource")
	}

	// Verify checksum
	hash := sha256.New()
	f, err := os.Open(filepath.Join(tempDir, "resource.rpm"))
	if err != nil {
		return nil, errors.Wrap(err, "failed to open resource")
	}
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
		return nil, errors.Wrap(err, "failed to hash resource")
	}
	if hex.EncodeToString(hash.Sum(nil)) != checksumSha256 {
		return nil, temporal.NewNonRetryableApplicationError(
			"checksum does not match",
			"checksumDoesNotMatch",
			errors.New("client submitted a checksum that does not match the resource"),
		)
	}

	// Read the RPM headers
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return nil, errors.Wrap(err, "failed to seek resource")
	}
	rpm, err := rpmutils.ReadRpm(f)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read RPM headers")
	}

	nevra, err := rpm.Header.GetNEVRA()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get RPM NEVRA")
	}

	// Ensure repository exists
	repoName := nevra.Name

	// First ensure that the repo exists.
	authenticator, err := w.forge.GetAuthenticator()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get forge authenticator")
	}

	err = w.forge.EnsureRepositoryExists(authenticator, repoName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure repository exists")
	}

	// Then do an import
	srpmState, err := srpm_import.FromFile(filepath.Join(tempDir, "resource.rpm"), w.rolling, w.gpgKeys...)
	if err != nil {
		if strings.Contains(err.Error(), "failed to verify RPM") {
			return nil, temporal.NewNonRetryableApplicationError(
				"failed to verify RPM",
				"failedToVerifyRPM",
				err,
			)
		}
		return nil, errors.Wrap(err, "failed to import SRPM")
	}
	defer srpmState.Close()
	srpmState.SetAuthor(authenticator.AuthorName, authenticator.AuthorEmail)

	cloneOpts := &git.CloneOptions{
		URL:  w.forge.GetRemote(repoName),
		Auth: authenticator.AuthMethod,
	}
	storer := memory.NewStorage()
	fs := memfs.New()
	importOut, err := srpmState.Import(cloneOpts, storer, fs, w.storage, osRelease)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import SRPM")
	}

	commitURI := w.forge.GetCommitViewerURL(repoName, importOut.Commit.Hash.String())

	return &mothershippb.ImportRPMResponse{
		CommitHash:   importOut.Commit.Hash.String(),
		CommitUri:    commitURI,
		CommitBranch: importOut.Branch,
		CommitTag:    importOut.Tag,
		Nevra:        nevra.String(),
		Pkg:          nevra.Name,
	}, nil
}
