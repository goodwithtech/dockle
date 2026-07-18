package user

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
	// etc/passwd body with a valid row, a blank line, a comment line and a
	// short/truncated entry. The short and blank lines used to panic on data[2].
	body := "root:x:0:0:root:/root:/bin/sh\n\n# a comment\ndaemon:x:1:1:daemon:/:/sbin/nologin\nbroken\nalsobroken:x\n"
	fileMap := types.FileMap{
		"etc/passwd": types.FileData{Body: []byte(body)},
	}

	assessor := UserAssessor{}
	if _, err := assessor.Assess(fileMap); err != nil {
		t.Fatalf("Assess returned an error: %v", err)
	}
}

func TestAssessDuplicateUID(t *testing.T) {
	// two rows share UID 0, so a duplicate-user assessment is expected.
	body := "root:x:0:0:root:/root:/bin/sh\ndup:x:0:0:dup:/:/bin/sh\n"
	fileMap := types.FileMap{
		"etc/passwd": types.FileData{Body: []byte(body)},
	}

	assessor := UserAssessor{}
	assesses, err := assessor.Assess(fileMap)
	if err != nil {
		t.Fatalf("Assess returned an error: %v", err)
	}
	found := false
	for _, a := range assesses {
		if a.Code == types.AvoidDuplicateUserGroup && a.Filename == "etc/passwd" {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a duplicate user assessment, got %v", assesses)
	}
}
