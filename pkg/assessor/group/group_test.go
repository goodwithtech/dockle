package group

import (
	"os"
	"testing"

	"github.com/goodwithtech/dockle/pkg/log"
	"github.com/goodwithtech/dockle/pkg/types"
)

func TestMain(m *testing.M) {
	if err := log.InitLogger(false, true); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestAssessMalformedLines(t *testing.T) {
	// etc/group body with a valid row, a blank line, a comment line and a
	// short/truncated entry. The short and blank lines used to panic on data[2].
	body := "root:x:0:\n\n# a comment\ndaemon:x:1:\nbroken\nalsobroken:x\n"
	fileMap := types.FileMap{
		"etc/group": types.FileData{Body: []byte(body)},
	}

	assessor := GroupAssessor{}
	if _, err := assessor.Assess(fileMap); err != nil {
		t.Fatalf("Assess returned an error: %v", err)
	}
}

func TestAssessDuplicateGID(t *testing.T) {
	// two rows share GID 0, so a duplicate-group assessment is expected.
	body := "root:x:0:\ndup:x:0:\n"
	fileMap := types.FileMap{
		"etc/group": types.FileData{Body: []byte(body)},
	}

	assessor := GroupAssessor{}
	assesses, err := assessor.Assess(fileMap)
	if err != nil {
		t.Fatalf("Assess returned an error: %v", err)
	}
	found := false
	for _, a := range assesses {
		if a.Code == types.AvoidDuplicateUserGroup && a.Filename == "etc/group" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a duplicate group assessment, got %v", assesses)
	}
}
