package manifest

import (
	"encoding/json"
	"io"
	"os"
	"sort"
	"testing"

	"github.com/d4l3k/messagediff"
	"github.com/goodwithtech/dockle/pkg/types"
)

func TestAssess(t *testing.T) {
	var tests = map[string]struct {
		path     string
		assesses []*types.Assessment
	}{
		"RootDefault": {
			path: "./testdata/root_default.json",
			assesses: []*types.Assessment{
				{
					Code:     types.AvoidRootDefault,
					Filename: ConfigFileName,
				},
				{
					Code:     types.AddHealthcheck,
					Filename: ConfigFileName,
				},
			},
		},
		"ApkCached": {
			path: "./testdata/apk_cache.json",

			assesses: []*types.Assessment{
				{
					Code:     types.UseApkAddNoCache,
					Filename: ConfigFileName,
				},
				{
					Code:     types.AvoidRootDefault,
					Filename: ConfigFileName,
				},
				{
					Code:     types.AddHealthcheck,
					Filename: ConfigFileName,
				},
				{
					Code:     types.UseCOPY,
					Filename: ConfigFileName,
				},
			},
		},
		"AptUpdateUpgrade": {
			path: "./testdata/apt_update_upgrade.json",

			assesses: []*types.Assessment{
				{
					Code:     types.AvoidRootDefault,
					Filename: ConfigFileName,
				},
				{
					Code:     types.MinimizeAptGet,
					Filename: ConfigFileName,
				},
				{
					Code:     types.AddHealthcheck,
					Filename: ConfigFileName,
				},
			},
		},
		"ADDStatementNotFirst": {
			path: "./testdata/add_with_arg_statement.json",

			assesses: []*types.Assessment{
				{
					Code:     types.AvoidRootDefault,
					Filename: ConfigFileName,
				},
				{
					Code:     types.AddHealthcheck,
					Filename: ConfigFileName,
				},
			},
		},
		"MultiADDStatements": {
			path: "./testdata/multi_add.json",

			assesses: []*types.Assessment{
				{
					Code:     types.AvoidRootDefault,
					Filename: ConfigFileName,
				},
				{
					Code:     types.AddHealthcheck,
					Filename: ConfigFileName,
				},
				{
					Code:     types.UseCOPY,
					Filename: ConfigFileName,
				},
				{
					Code:     types.UseCOPY,
					Filename: ConfigFileName,
				},
			},
		},
	}

	for testname, v := range tests {
		d, err := loadImageFromFile(v.path)

		if err != nil {
			t.Errorf("%s : can't open file %s", testname, v.path)
			continue
		}
		actual, err := checkAssessments(d)
		if err != nil {
			t.Errorf("%s : catch the error : %v", testname, err)
		}

		diff, equal := messagediff.PrettyDiff(
			sortByType(v.assesses),
			sortByType(actual),
			messagediff.IgnoreStructField("Desc"),
		)
		if !equal {
			t.Errorf("%s diff : %v", testname, diff)
		}
	}
}

func TestSplitByCommands(t *testing.T) {
	var tests = map[string]struct {
		path     string
		index    int
		expected map[int][]string
	}{
		"RootDefault": {
			path:  "./testdata/root_default.json",
			index: 1,
			expected: map[int][]string{
				0: {"/bin/sh", "-c", "#(nop)", "CMD", "[\"/bin/sh\"]"},
			},
		},
		"Nginx": {
			path:  "./testdata/nginx.json",
			index: 6,
			expected: map[int][]string{
				0:  {"/bin/sh", "-c", "set", "-x"},
				1:  {"addgroup", "--system", "--gid", "101", "nginx"},
				2:  {"adduser", "--system", "--disabled-login", "--ingroup", "nginx", "--no-create-home", "--home", "/nonexistent", "--gecos", "\"nginx", "user\"", "--shell", "/bin/false", "--uid", "101", "nginx"},
				3:  {"apt-get", "update"},
				4:  {"apt-get", "install", "--no-install-recommends", "--no-install-suggests", "-y", "gnupg1", "apt-transport-https", "ca-certificates"},
				5:  {"NGINX_GPGKEY=573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62;", "found='';", "for", "server", "in", "ha.pool.sks-keyservers.net", "hkp://keyserver.ubuntu.com:80", "hkp://p80.pool.sks-keyservers.net:80", "pgp.mit.edu", ";", "do", "echo", "\"Fetching", "GPG", "key", "$NGINX_GPGKEY", "from", "$server\";", "apt-key", "adv", "--keyserver", "\"$server\"", "--keyserver-options", "timeout=10", "--recv-keys", "\"$NGINX_GPGKEY\""},
				6:  {"found=yes"},
				7:  {"break;", "done;", "test", "-z", "\"$found\""},
				8:  {"echo", ">&2", "\"error:", "failed", "to", "fetch", "GPG", "key", "$NGINX_GPGKEY\""},
				9:  {"exit", "1;", "apt-get", "remove", "--purge", "--auto-remove", "-y", "gnupg1"},
				10: {"rm", "-rf", "/var/lib/apt/lists/*"},
				11: {"dpkgArch=\"$(dpkg", "--print-architecture)\""},
				12: {"nginxPackages=\"", "nginx=${NGINX_VERSION}-${PKG_RELEASE}", "nginx-module-xslt=${NGINX_VERSION}-${PKG_RELEASE}", "nginx-module-geoip=${NGINX_VERSION}-${PKG_RELEASE}", "nginx-module-image-filter=${NGINX_VERSION}-${PKG_RELEASE}", "nginx-module-njs=${NGINX_VERSION}.${NJS_VERSION}-${PKG_RELEASE}", "\""},
				13: {"case", "\"$dpkgArch\"", "in", "amd64|i386)", "echo", "\"deb", "https://nginx.org/packages/mainline/debian/", "stretch", "nginx\"", ">>", "/etc/apt/sources.list.d/nginx.list"},
				14: {"apt-get", "update", ";;", "*)", "echo", "\"deb-src", "https://nginx.org/packages/mainline/debian/", "stretch", "nginx\"", ">>", "/etc/apt/sources.list.d/nginx.list"},
				15: {"tempDir=\"$(mktemp", "-d)\""},
				16: {"chmod", "777", "\"$tempDir\""},
				17: {"savedAptMark=\"$(apt-mark", "showmanual)\""},
				18: {"apt-get", "update"},
				19: {"apt-get", "build-dep", "-y", "$nginxPackages"},
				20: {"(", "cd", "\"$tempDir\""},
				21: {"DEB_BUILD_OPTIONS=\"nocheck", "parallel=$(nproc)\"", "apt-get", "source", "--compile", "$nginxPackages", ")"},
				22: {"apt-mark", "showmanual", "|", "xargs", "apt-mark", "auto", ">", "/dev/null"},
				23: {"{", "[", "-z", "\"$savedAptMark\"", "]", "||", "apt-mark", "manual", "$savedAptMark;", "}"},
				24: {"ls", "-lAFh", "\"$tempDir\""},
				25: {"(", "cd", "\"$tempDir\""},
				26: {"dpkg-scanpackages", ".", ">", "Packages", ")"},
				27: {"grep", "'^Package:", "'", "\"$tempDir/Packages\""},
				28: {"echo", "\"deb", "[", "trusted=yes", "]", "file://$tempDir", "./\"", ">", "/etc/apt/sources.list.d/temp.list"},
				29: {"apt-get", "-o", "Acquire::GzipIndexes=false", "update", ";;", "esac"},
				30: {"apt-get", "install", "--no-install-recommends", "--no-install-suggests", "-y", "$nginxPackages", "gettext-base"},
				31: {"apt-get", "remove", "--purge", "--auto-remove", "-y", "apt-transport-https", "ca-certificates"},
				32: {"rm", "-rf", "/var/lib/apt/lists/*", "/etc/apt/sources.list.d/nginx.list"},
				33: {"if", "[", "-n", "\"$tempDir\"", "];", "then", "apt-get", "purge", "-y", "--auto-remove"},
				34: {"rm", "-rf", "\"$tempDir\"", "/etc/apt/sources.list.d/temp.list;", "fi"},
			},
		},
	}

	for testname, v := range tests {
		d, err := loadImageFromFile(v.path)
		if err != nil {
			t.Errorf("%s : can't open file %s", testname, v.path)
			continue
		}
		cmd := d.History[v.index]
		actual := splitByCommands(cmd.CreatedBy)
		diff, equal := messagediff.PrettyDiff(
			v.expected,
			actual,
		)
		if !equal {
			t.Errorf("%s diff : %v", testname, diff)
		}
	}
}

func TestReducableApkAdd(t *testing.T) {
	var tests = map[string]struct {
		cmdSlices map[int][]string
		expected  bool
	}{
		"Reducable": {
			cmdSlices: map[int][]string{
				0: {
					"apk", "add", "git",
				},
			},
			expected: true,
		},
		"UnReducable": {
			cmdSlices: map[int][]string{
				0: {
					"apk", "add", "--no-cache", "git",
				},
			},
			expected: false,
		},
	}
	for testname, v := range tests {
		actual := reducableApkAdd(v.cmdSlices)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}
}

func TestReducableAptGetUpdate(t *testing.T) {
	var tests = map[string]struct {
		cmdSlices map[int][]string
		expected  bool
	}{
		"Reducable": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "update",
				},
				1: {
					"apt-get", "purge",
				},
			},
			expected: true,
		},
		"NoUpdate": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "install",
				},
				1: {
					"apt-get", "purge",
				},
			},
			expected: false,
		},
		"UnReducable": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "update",
				},
				1: {
					"apt-get", "-y", "--no-install-recommends", "install",
				},
			},
			expected: false,
		},
		"UpdateAfterInstalled": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "-y", "--no-install-recommends", "install",
				},
				1: {
					"apt-get", "update",
				},
			},
			expected: true,
		},
		"CheckAptCommand": {
			cmdSlices: map[int][]string{
				0: {
					"apt", "update",
				},
				1: {
					"apt", "-y", "--no-install-recommends", "install",
				},
			},
			expected: false,
		},
		"LongInvalidCommand": {
			// https://github.com/docker-library/golang/blob/3f2c52653043f067156ce4f41182c2a758c4c857/1.17/alpine3.14/Dockerfile#L20-L107
			// Issue: https://github.com/goodwithtech/dockle/issues/151
			cmdSlices: map[int][]string{
				0: {
					"/bin/sh", "-c", "set", "-eux;", "apk", "add", "--no-cache", "--virtual", ".fetch-deps", "gnupg;", "arch=$(apk", "--print-arch);", "url=;",
					"case", "$arch", "in", "'x86_64')", "export", "GOARCH='amd64'", "GOOS='linux';", ";;", "'armhf')", "export", "GOARCH='arm'", "GOARM='6'", "GOOS='linux';", ";;", "'armv7')",
					"export", "GOARCH='arm'", "GOARM='7'", "GOOS='linux';", ";;", "'aarch64')", "export", "GOARCH='arm64'", "GOOS='linux';", ";;", "'x86')", "export", "GO386='softfloat'", "GOARCH='386'",
					"GOOS='linux';", ";;", "'ppc64le')", "export", "GOARCH='ppc64le'", "GOOS='linux';", ";;", "'s390x')", "export", "GOARCH='s390x'", "GOOS='linux';", ";;", "*)", "echo", ">&2", "error:",
					"unsupported", "architecture", "'$arch'", "(likely", "packaging", "update", "needed);", "exit", "1", ";;", "esac;", "build=;", "if", "[", "-z", "$url", "];", "then", "build=1;",
					"url='https://dl.google.com/go/go1.17.1.src.tar.gz';", "sha256='49dc08339770acd5613312db8c141eaf61779995577b89d93b541ef83067e5b1';", "fi;", "wget", "-O", "go.tgz.asc", "$url.asc;", "wget",
					"-O", "go.tgz", "$url;", "echo", "$sha256", "*go.tgz", "|", "sha256sum", "-c", "-;", "GNUPGHOME=$(mktemp", "-d);", "export", "GNUPGHOME;", "gpg", "--batch", "--keyserver", "keyserver.ubuntu.com",
					"--recv-keys", "'EB4C", "1BFD", "4F04", "2F6D", "DDCC", "EC91", "7721", "F63B", "D38B", "4796';", "gpg", "--batch", "--verify", "go.tgz.asc", "go.tgz;", "gpgconf", "--kill", "all;", "rm", "-rf", "$GNUPGHOME",
					"go.tgz.asc;", "tar", "-C", "/usr/local", "-xzf", "go.tgz;", "rm", "go.tgz;", "if", "[", "-n", "$build", "];", "then", "apk", "add", "--no-cache", "--virtual", ".build-deps", "bash", "gcc", "go", "musl-dev",
					";", "(", "cd", "/usr/local/go/src;", "export", "GOROOT_BOOTSTRAP=$(go", "env", "GOROOT)", "GOHOSTOS=$GOOS", "GOHOSTARCH=$GOARCH;", "./make.bash;", ");", "apk", "del", "--no-network", ".build-deps;", "go", "install",
					"std;", "rm", "-rf", "/usr/local/go/pkg/*/cmd", "/usr/local/go/pkg/bootstrap", "/usr/local/go/pkg/obj", "/usr/local/go/pkg/tool/*/api", "/usr/local/go/pkg/tool/*/go_bootstrap",
					"/usr/local/go/src/cmd/dist/dist", ";", "fi;", "apk", "del", "--no-network", ".fetch-deps;", "go", "version",
				},
			},
			expected: false,
		},
	}
	for testname, v := range tests {
		actual := reducableAptGetUpdate(v.cmdSlices)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}
}

func TestReducableAptGetInstall(t *testing.T) {
	var tests = map[string]struct {
		cmdSlices map[int][]string
		expected  bool
	}{
		"Reducable": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "-y", "install",
				},
				1: {
					"apt-get", "update",
				},
			},
			expected: true,
		},
		"OnlyUpdate": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "update",
				},
				1: {
					"apt-get", "purge",
				},
			},
			expected: true,
		},
		"NoUpdateInstall": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "purge",
				},
			},
			expected: false,
		},
		"UnReducable": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "install",
				},
				1: {
					"rm", "-fR", "/var/lib/apt/lists/*",
				},
			},
			expected: false,
		},
		"UnReducable2": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "install", "-y", "git",
				},
				1: {
					"rm", "-rf", "/var/lib/apt/lists",
				},
			},
			expected: false,
		},
		"UnReducable3": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "install", "-y", "git",
				},
				1: {
					"rm", "-r", "/var/lib/apt/lists",
				},
			},
			expected: false,
		},
	}
	for testname, v := range tests {
		actual := reducableAptGetInstall(v.cmdSlices)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}
}

func TestAddStatement(t *testing.T) {
	var tests = map[string]struct {
		cmdSlices map[int][]string
		expected  bool
	}{
		"UseADD": {
			cmdSlices: map[int][]string{
				0: {
					"/bin/sh", "-c", "#(nop)", "ADD", "file:2e3a37883f56a4a278bec2931fc9f91fb9ebdaa9047540fe8fde419b84a1701b", "in", "/cmd",
				},
			},
			expected: true,
		},
		"NotADD": {
			cmdSlices: map[int][]string{
				0: {"/bin/sh", "-c", "set", "-x"},
				1: {"addgroup", "--system", "--gid", "101", "nginx"},
			},
			expected: false,
		},
		"UseADDR": {
			cmdSlices: map[int][]string{
				0: {"/bin/sh", "-c", "set", "-x"},
				1: {"/bin/sh", "-c", "RETHINKDB_CLUSTER_IP_ADDR"},
			},
			expected: false,
		},
	}
	for testname, v := range tests {
		actual := useADDstatement(v.cmdSlices)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}
}

func TestSensitiveVars(t *testing.T) {
	if err := compileSensitivePatterns(); err != nil {
		t.Fatalf("compile sensitive var patterns: %s", err)
	}
	var tests = map[string]struct {
		cmd      string
		expected bool
	}{
		"basic":              {cmd: "/bin/sh -c #(nop) ENV PASS=ADMIN", expected: true},
		"empty value":        {cmd: "/bin/sh -c #(nop) ENV PASS=", expected: false},
		"mixed cases":        {cmd: "/bin/sh -c #(nop) ENV PasS=ADMIN", expected: true},
		"two vars":           {cmd: "/bin/sh -c #(nop) ENV abc=hello password=sensibledata", expected: true},
		"empty two value":    {cmd: "/bin/sh -c #(nop) ENV ABC=hello PASS= ", expected: false},
		"run command":        {cmd: `/bin/sh -c SECRET_API_KEY=63AF7AA15067C05616FDDD88A3A2E8F226F0BC06 echo "data"`, expected: true},
		"run false positive": {cmd: `/bin/sh -c HELLO="PASS=\"notThis\"" echo "false positive"`, expected: false},
		"run command 2":      {cmd: `/bin/sh -c SECRET=myLittleSecret VAR2=VALUE2 VAR3=VALUE3 echo "Do something"`, expected: true},
		"secret with space":  {cmd: `/bin/sh -c SECRET="hello world"`, expected: true},
		"skip space key":     {cmd: `/bin/sh -c echo 'secret = foo;' > test.conf`, expected: false},
		// TODO : expected must be false
		//"skip echo string": {cmd: `/bin/sh -c echo 'secret=foo;' > test.conf`, expected: true},
	}
	for testname, v := range tests {
		actual, _, _ := sensitiveVars(v.cmd)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}

}

func TestUseDistUpgrade(t *testing.T) {
	var tests = map[string]struct {
		cmdSlices map[int][]string
		expected  bool
	}{
		"UseUpgrade": {
			cmdSlices: map[int][]string{
				0: {
					"apt-get", "upgrade",
				},
			},
			expected: false,
		},
		"UseAptUpgrade": {
			cmdSlices: map[int][]string{
				0: {"apt", "upgrade"},
				1: {"addgroup", "--system", "--gid", "101", "nginx"},
			},
			expected: false,
		},
		"UseDistUpgrade": {
			cmdSlices: map[int][]string{
				0: {"apt-get", "dist-upgrade"},
			},
			expected: true,
		},
		"UseAptDistUpgrade": {
			cmdSlices: map[int][]string{
				0: {"apt", "dist-upgrade"},
			},
			expected: true,
		},

		"NoAptDistUpgrade": {
			cmdSlices: map[int][]string{
				0: {"somecommand", "dist-upgrade", "pip", "setuptools"},
			},
			expected: false,
		},
	}
	for testname, v := range tests {
		actual := useDistUpgrade(v.cmdSlices)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}
}

func TestContainsThreshold(t *testing.T) {
	var tests = map[string]struct {
		heystack  []string
		needles   []string
		threshold int
		expected  bool
	}{
		"SimpleSuccess2": {
			heystack:  []string{"a", "b", "c", "d"},
			needles:   []string{"a", "b", "x", "z"},
			threshold: 2,
			expected:  true,
		},
		"SimpleFail2": {
			heystack:  []string{"a", "b", "c", "d"},
			needles:   []string{"a", "x", "y", "z"},
			threshold: 3,
			expected:  false,
		},
		"SimpleSuccess3": {
			heystack:  []string{"a", "b", "c", "d"},
			needles:   []string{"a", "b", "c", "z"},
			threshold: 3,
			expected:  true,
		},
		"SimpleFail3": {
			heystack:  []string{"a", "b", "d", "f"},
			needles:   []string{"a", "b", "c", "z"},
			threshold: 3,
			expected:  false,
		},
		"DuplicateHeystackSuccess": {
			heystack:  []string{"a", "a", "b", "c"},
			needles:   []string{"a", "a", "y", "z"},
			threshold: 2,
			expected:  false,
		},
		"DuplicateHeystackFail": {
			heystack:  []string{"a", "a", "b", "c"},
			needles:   []string{"a", "x", "y", "z"},
			threshold: 2,
			expected:  false,
		},
	}
	for testname, v := range tests {
		actual := containsThreshold(v.heystack, v.needles, v.threshold)
		if actual != v.expected {
			t.Errorf("%s want: %t, got %t", testname, v.expected, actual)
		}
	}
}

func loadImageFromFile(path string) (config types.Image, err error) {
	read, err := os.Open(path)
	if err != nil {
		return config, err
	}
	filebytes, err := io.ReadAll(read)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(filebytes, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func sortByType(assesses []*types.Assessment) []*types.Assessment {
	sort.Slice(assesses, func(i, j int) bool {
		if assesses[i].Code != assesses[j].Code {
			return assesses[i].Code < assesses[j].Code
		}
		return assesses[i].Code < assesses[j].Code
	})
	return assesses
}
