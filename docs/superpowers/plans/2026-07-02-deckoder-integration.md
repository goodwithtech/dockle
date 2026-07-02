# deckoder統合 実装プラン

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 外部依存 `github.com/goodwithtech/deckoder` を dockle 本体（pkg/配下）へ取り込み、依存を廃止する。

**Architecture:** deckoder のソースを module cache から dockle の `pkg/` 配下へコピーし、内部 import パスを機械的に書き換える。deckoder の `types` は dockle の `pkg/types` へ完全マージし、`deckodertypes` エイリアスを全廃する。`types` を葉とする一方向依存のため循環importは発生しない。既存テスト・モックも移送して回帰を担保する。

**Tech Stack:** Go 1.22、`containers/image/v5`、`urfave/cli` v1、`stretchr/testify`（モック）、`knqyf263/nested`。

## 実装上の重要な性質（TDDの扱い）

本作業は**テスト済みの既存コードの逐語的ポート**である。新規ロジックは書かない。したがって各タスクの検証は「移送したパッケージが `go build` し、移送した既存テストが `go test` で通り、`go vet` が通る」ことで行う。古典的な「失敗するテストを先に書く」サイクルは適用しない（テストは deckoder から移送されるため既に存在する）。

## Global Constraints

- Go バージョン: `go 1.22.10`（`go.mod` の `go` ディレクティブを変更しない）。
- CI と同条件で検証すること: `CGO_ENABLED=0 go test ./...`。
- コードは逐語コピーを原則とし、抽出ロジックのリファクタリング・機能追加・スタイル変更（`ioutil`→`os` 置換等）は**行わない**。差分は import パスと `deckodertypes`→`types` の置換に限定する。
- 取り込むコードは Apache-2.0 扱い。deckoder の `LICENSE`（AGPL-3.0）ファイルは dockle へ**コピーしない**。deckoder の `.go` にライセンスヘッダは無いため、ヘッダ除去作業は不要（確認のみ）。
- deckoder の CLI（`cmd/deckoder/main.go`）は取り込まない。
- module cache のソースは `$(go env GOMODCACHE)/github.com/goodwithtech/deckoder@v0.0.6` にある。コピー先ファイルは cache から取得後 `chmod u+w` すること（cache は read-only 0444）。

### import パス書き換えルール（全コピーファイル共通）

コピーした `.go` に対し、以下の 4 置換を適用する。`deckoder/extractor` はプレフィックス置換で `deckoder/extractor/docker` 等のサブパスも同時に正しく変換される。適用順は問わない（プレフィックスが重複しないため）。

| 置換前 | 置換後 |
|---|---|
| `github.com/goodwithtech/deckoder/types` | `github.com/goodwithtech/dockle/pkg/types` |
| `github.com/goodwithtech/deckoder/utils` | `github.com/goodwithtech/dockle/pkg/deckoder/utils` |
| `github.com/goodwithtech/deckoder/analyzer` | `github.com/goodwithtech/dockle/pkg/analyzer` |
| `github.com/goodwithtech/deckoder/extractor` | `github.com/goodwithtech/dockle/pkg/extractor` |

sed ワンライナー（各コピー先ファイルに適用）:

```bash
sed -i '' \
  -e 's#github.com/goodwithtech/deckoder/types#github.com/goodwithtech/dockle/pkg/types#g' \
  -e 's#github.com/goodwithtech/deckoder/utils#github.com/goodwithtech/dockle/pkg/deckoder/utils#g' \
  -e 's#github.com/goodwithtech/deckoder/analyzer#github.com/goodwithtech/dockle/pkg/analyzer#g' \
  -e 's#github.com/goodwithtech/deckoder/extractor#github.com/goodwithtech/dockle/pkg/extractor#g' \
  "$TARGET"
```

（macOS の `sed -i ''` 形式。上の順序で `types` を先に処理しても `extractor` 置換と衝突しない。）

---

## Task 1: deckoder types を pkg/types へマージ

**Files:**
- Create: `pkg/types/filemap.go`
- Create: `pkg/types/docker_option.go`
- Create: `pkg/types/imageref.go`
- Create: `pkg/types/const.go`
- Modify: `pkg/types/error.go`（2 変数を追記）

**Interfaces:**
- Produces（後続タスク・全アセッサが依存）:
  - `type FilterFunc func(*tar.Header) (bool, error)`
  - `type FileMap map[string]FileData`
  - `type FileData struct { Body []byte; FileMode os.FileMode }`
  - `type DockerOption struct { ... }`（下記コード参照）
  - `type FilePath string`
  - `type ImageReference struct { Name string; ID digest.Digest; LayerIDs []string }`
  - `type ImageDetail struct { Files FileMap }`
  - `type ImageInfo struct { ... }`
  - `type LayerInfo struct { ID digest.Digest; SchemaVersion int; TargetFiles FileMap; OpaqueDirs []string; WhiteoutFiles []string }`
  - `const ImageJSONSchemaVersion = 1`, `const LayerJSONSchemaVersion = 1`
  - `var InvalidURLPattern`, `var ErrNoRpmCmd`
- すべて package `types`（`github.com/goodwithtech/dockle/pkg/types`）。

- [ ] **Step 1: `pkg/types/filemap.go` を作成**

```go
package types

import (
	"archive/tar"
	"os"
)

type FilterFunc func(*tar.Header) (bool, error)

type FileMap map[string]FileData
type FileData struct {
	Body     []byte
	FileMode os.FileMode
}
```

- [ ] **Step 2: `pkg/types/docker_option.go` を作成**

```go
package types

import "time"

type DockerOption struct {
	// Auth
	UserName string
	Password string

	// ECR
	AwsAccessKey    string
	AwsSecretKey    string
	AwsSessionToken string
	AwsRegion       string

	// GCP
	GcpCredPath string

	// Docker daemon
	DockerDaemonCertPath string
	DockerDaemonHost     string

	InsecureSkipTLSVerify bool
	SkipPing              bool
	Timeout               time.Duration
}
```

- [ ] **Step 3: `pkg/types/imageref.go` を作成**

```go
package types

import (
	"time"

	digest "github.com/opencontainers/go-digest"
)

type FilePath string

type ImageReference struct {
	Name     string // image name or tar file name
	ID       digest.Digest
	LayerIDs []string
}

type ImageDetail struct {
	Files FileMap
}

// ImageInfo is stored in cache
type ImageInfo struct {
	SchemaVersion int
	Architecture  string
	Created       time.Time
	DockerVersion string
	OS            string
}

// LayerInfo is stored in cache
type LayerInfo struct {
	ID            digest.Digest `json:",omitempty"`
	SchemaVersion int
	TargetFiles   FileMap
	OpaqueDirs    []string `json:",omitempty"`
	WhiteoutFiles []string `json:",omitempty"`
}
```

- [ ] **Step 4: `pkg/types/const.go` を作成**

```go
package types

const (
	ImageJSONSchemaVersion = 1
	LayerJSONSchemaVersion = 1
)
```

- [ ] **Step 5: `pkg/types/error.go` に 2 変数を追記**

既存ファイルは以下（`errors` を使用）。deckoder の 2 変数を `errors.New` へ変換して同一 `var` ブロックに統合する（`xerrors` は導入しない）。編集後の全文:

```go
package types

import "errors"

var (
	ErrSetImageOrFile = errors.New("image name or image file must be specified")
	InvalidURLPattern = errors.New("invalid url pattern")
	ErrNoRpmCmd       = errors.New("no rpm command")
)
```

- [ ] **Step 6: pkg/types をビルド・vet して確認**

Run: `go build ./pkg/types/ && go vet ./pkg/types/`
Expected: エラーなし（既存の `pkg/types` テストも壊さないこと）。

- [ ] **Step 7: pkg/types のテストを実行**

Run: `CGO_ENABLED=0 go test ./pkg/types/`
Expected: PASS（既存の `assessment_test.go` が通る）。

- [ ] **Step 8: コミット**

```bash
git add pkg/types/
git commit -m "[機能] deckoder types を pkg/types へマージ"
```

---

## Task 2: utils を pkg/deckoder/utils へ移送

**Files:**
- Create: `pkg/deckoder/utils/utils.go`
- Create: `pkg/deckoder/utils/utils_test.go`

**Interfaces:**
- Consumes: `types.FilterFunc`（Task 1）。
- Produces（cache アセッサ・docker テストが依存）:
  - `func CacheDir() string`
  - `func StringInSlice(a string, list []string) bool`
  - `func IsCommandAvailable(name string) bool`
  - `func IsGzip(f *bufio.Reader) bool`
  - `func CreateFilterPathFunc(filenames []string) types.FilterFunc`
  - `var PathSeparator string`
- package `utils`。

- [ ] **Step 1: ソースをコピーして import 書き換え**

```bash
DECK="$(go env GOMODCACHE)/github.com/goodwithtech/deckoder@v0.0.6"
mkdir -p pkg/deckoder/utils
cp "$DECK/utils/utils.go" pkg/deckoder/utils/utils.go
cp "$DECK/utils/utils_test.go" pkg/deckoder/utils/utils_test.go
chmod u+w pkg/deckoder/utils/*.go
for TARGET in pkg/deckoder/utils/utils.go pkg/deckoder/utils/utils_test.go; do
  sed -i '' \
    -e 's#github.com/goodwithtech/deckoder/types#github.com/goodwithtech/dockle/pkg/types#g' \
    -e 's#github.com/goodwithtech/deckoder/utils#github.com/goodwithtech/dockle/pkg/deckoder/utils#g' \
    -e 's#github.com/goodwithtech/deckoder/analyzer#github.com/goodwithtech/dockle/pkg/analyzer#g' \
    -e 's#github.com/goodwithtech/deckoder/extractor#github.com/goodwithtech/dockle/pkg/extractor#g' \
    "$TARGET"
done
```

- [ ] **Step 2: import 書き換えの確認**

Run: `grep -rn "goodwithtech/deckoder" pkg/deckoder/utils/`
Expected: 出力なし（deckoder 参照が残っていないこと）。

- [ ] **Step 3: ビルド・vet・テスト**

Run: `CGO_ENABLED=0 go test ./pkg/deckoder/utils/ && go vet ./pkg/deckoder/utils/`
Expected: PASS。

- [ ] **Step 4: コミット**

```bash
git add pkg/deckoder/utils/
git commit -m "[機能] deckoder utils を pkg/deckoder/utils へ移送"
```

---

## Task 3: extractor / image / token(ecr,gcr) を pkg/extractor へ移送

`extractor.go`・`image/**`・`token/ecr`・`token/gcr` はいずれも `types` のみに依存する（相互依存なし）。まとめて移送する。

**Files:**
- Create: `pkg/extractor/extractor.go`
- Create: `pkg/extractor/image/image.go`
- Create: `pkg/extractor/image/token.go`
- Create: `pkg/extractor/image/mock_image.go`
- Create: `pkg/extractor/image/mock_image_closer.go`
- Create: `pkg/extractor/image/mock_image_source.go`
- Create: `pkg/extractor/image/mock_registry.go`
- Create: `pkg/extractor/image/image_test.go`
- Create: `pkg/extractor/image/token_test.go`
- Create: `pkg/extractor/image/token/ecr/ecr.go`
- Create: `pkg/extractor/image/token/ecr/ecr_test.go`
- Create: `pkg/extractor/image/token/gcr/gcr.go`
- Create: `pkg/extractor/image/token/gcr/gcr_test.go`

**Interfaces:**
- Consumes: `types.DockerOption`, `types.FilterFunc`, `types.FileMap`, `types.InvalidURLPattern`（Task 1）。
- Produces（docker が依存）:
  - `type Extractor interface { ImageName() string; ImageID() digest.Digest; ConfigBlob(ctx) ([]byte, error); LayerIDs() []string; ExtractLayerFiles(ctx, digest.Digest, types.FilterFunc) (types.FileMap, []string, []string, error) }`（package `extractor`）
  - `type Image interface`, `type Reference struct { Name string; IsFile bool }`, `func NewImage(...)`, `func RegisterRegistry(Registry)`, `func GetToken(...)`, `type Registry interface`（package `image`）
  - `type ECR struct{}`（package `ecr`）, `type GCR struct{}`（package `gcr`）
- 各ディレクトリの package 名: `extractor` / `image` / `ecr` / `gcr`。

- [ ] **Step 1: ディレクトリ作成とコピー・import 書き換え**

```bash
DECK="$(go env GOMODCACHE)/github.com/goodwithtech/deckoder@v0.0.6"
mkdir -p pkg/extractor/image/token/ecr pkg/extractor/image/token/gcr

cp "$DECK/extractor/extractor.go" pkg/extractor/extractor.go
cp "$DECK"/extractor/image/*.go pkg/extractor/image/
cp "$DECK"/extractor/image/token/ecr/*.go pkg/extractor/image/token/ecr/
cp "$DECK"/extractor/image/token/gcr/*.go pkg/extractor/image/token/gcr/

find pkg/extractor -name '*.go' -exec chmod u+w {} +
find pkg/extractor -name '*.go' | while read -r TARGET; do
  sed -i '' \
    -e 's#github.com/goodwithtech/deckoder/types#github.com/goodwithtech/dockle/pkg/types#g' \
    -e 's#github.com/goodwithtech/deckoder/utils#github.com/goodwithtech/dockle/pkg/deckoder/utils#g' \
    -e 's#github.com/goodwithtech/deckoder/analyzer#github.com/goodwithtech/dockle/pkg/analyzer#g' \
    -e 's#github.com/goodwithtech/deckoder/extractor#github.com/goodwithtech/dockle/pkg/extractor#g' \
    "$TARGET"
done
```

（注: `mock_registry.go` は `import types "…deckoder/types"` を持つ。上記 sed で `…dockle/pkg/types` に変換され、エイリアス `types` はそのまま機能する。`mock_image.go` 等の `import types "…containers/image/v5/types"` は deckoder 参照ではないため影響を受けない。）

- [ ] **Step 2: deckoder 参照が残っていないか確認**

Run: `grep -rn "goodwithtech/deckoder" pkg/extractor/`
Expected: 出力なし。

- [ ] **Step 3: ビルド・vet・テスト**

Run: `CGO_ENABLED=0 go test ./pkg/extractor/... && go vet ./pkg/extractor/...`
Expected: PASS。

- [ ] **Step 4: コミット**

```bash
git add pkg/extractor/
git commit -m "[機能] deckoder extractor/image/token を pkg/extractor へ移送"
```

---

## Task 4: extractor/docker を pkg/extractor/docker へ移送

`docker.go` は `extractor/image`・`token/ecr`・`token/gcr`・`types`・`knqyf263/nested` に依存する。Task 3 完了後に移送する。

**Files:**
- Create: `pkg/extractor/docker/docker.go`
- Create: `pkg/extractor/docker/docker_test.go`

**Interfaces:**
- Consumes: `image.Reference`, `image.NewImage`, `image.RegisterRegistry`, `ecr.ECR`, `gcr.GCR`, `types.*`（Task 3, 1）, `utils.*`（Task 2、テストが使用）。
- Produces（scanner が依存）:
  - `func NewDockerExtractor(ctx, imageName string, option types.DockerOption) (Extractor, func(), error)`
  - `func NewDockerArchiveExtractor(ctx, fileName string, option types.DockerOption) (Extractor, func(), error)`
  - `func ApplyLayers(layers []types.LayerInfo) types.FileMap`
  - `type Extractor struct{}`（`extractor.Extractor` インターフェースを満たす）
- package `docker`。

- [ ] **Step 1: コピー・import 書き換え**

```bash
DECK="$(go env GOMODCACHE)/github.com/goodwithtech/deckoder@v0.0.6"
mkdir -p pkg/extractor/docker
cp "$DECK"/extractor/docker/*.go pkg/extractor/docker/
find pkg/extractor/docker -name '*.go' -exec chmod u+w {} +
find pkg/extractor/docker -name '*.go' | while read -r TARGET; do
  sed -i '' \
    -e 's#github.com/goodwithtech/deckoder/types#github.com/goodwithtech/dockle/pkg/types#g' \
    -e 's#github.com/goodwithtech/deckoder/utils#github.com/goodwithtech/dockle/pkg/deckoder/utils#g' \
    -e 's#github.com/goodwithtech/deckoder/analyzer#github.com/goodwithtech/dockle/pkg/analyzer#g' \
    -e 's#github.com/goodwithtech/deckoder/extractor#github.com/goodwithtech/dockle/pkg/extractor#g' \
    "$TARGET"
done
```

- [ ] **Step 2: deckoder 参照が残っていないか確認**

Run: `grep -rn "goodwithtech/deckoder" pkg/extractor/docker/`
Expected: 出力なし。

- [ ] **Step 3: ビルド・vet・テスト**

Run: `CGO_ENABLED=0 go test ./pkg/extractor/docker/ && go vet ./pkg/extractor/docker/`
Expected: PASS。

- [ ] **Step 4: コミット**

```bash
git add pkg/extractor/docker/
git commit -m "[機能] deckoder extractor/docker を pkg/extractor/docker へ移送"
```

---

## Task 5: analyzer を pkg/analyzer へ移送

`analyzer.go` は `extractor`・`extractor/docker`・`types` に依存する。Task 4 完了後に移送する。

**Files:**
- Create: `pkg/analyzer/analyzer.go`

**Interfaces:**
- Consumes: `extractor.Extractor`, `docker.ApplyLayers`, `types.*`（Task 3, 4, 1）。
- Produces（scanner が依存）:
  - `func New(ext extractor.Extractor) Config`
  - `type Config struct { Extractor extractor.Extractor }`
  - `func (ac Config) Analyze(ctx context.Context, filterFunc types.FilterFunc) (types.FileMap, error)`
- package `analyzer`。

- [ ] **Step 1: コピー・import 書き換え**

```bash
DECK="$(go env GOMODCACHE)/github.com/goodwithtech/deckoder@v0.0.6"
mkdir -p pkg/analyzer
cp "$DECK/analyzer/analyzer.go" pkg/analyzer/analyzer.go
chmod u+w pkg/analyzer/analyzer.go
sed -i '' \
  -e 's#github.com/goodwithtech/deckoder/types#github.com/goodwithtech/dockle/pkg/types#g' \
  -e 's#github.com/goodwithtech/deckoder/utils#github.com/goodwithtech/dockle/pkg/deckoder/utils#g' \
  -e 's#github.com/goodwithtech/deckoder/analyzer#github.com/goodwithtech/dockle/pkg/analyzer#g' \
  -e 's#github.com/goodwithtech/deckoder/extractor#github.com/goodwithtech/dockle/pkg/extractor#g' \
  pkg/analyzer/analyzer.go
```

- [ ] **Step 2: deckoder 参照が残っていないか確認**

Run: `grep -rn "goodwithtech/deckoder" pkg/analyzer/`
Expected: 出力なし。

- [ ] **Step 3: ビルド・vet**

Run: `go build ./pkg/analyzer/ && go vet ./pkg/analyzer/`
Expected: エラーなし。

- [ ] **Step 4: コミット**

```bash
git add pkg/analyzer/
git commit -m "[機能] deckoder analyzer を pkg/analyzer へ移送"
```

---

## Task 6: dockle 側の import を新パスへ切り替え

移送が完了したので、dockle の呼び出し側を新パスへ切り替え、`deckodertypes` エイリアスを全廃する。この時点で外部 deckoder への参照がコードから消える。

**Files:**
- Modify: `pkg/run.go`
- Modify: `pkg/scanner/scan.go`
- Modify: `pkg/scanner/scan_test.go`
- Modify: `pkg/assessor/assessor.go`
- Modify: `pkg/assessor/cache/cache.go`
- Modify: `pkg/assessor/contentTrust/contentTrust.go`
- Modify: `pkg/assessor/credential/credential.go`
- Modify: `pkg/assessor/manifest/manifest.go`
- Modify: `pkg/assessor/privilege/suid.go`
- Modify: `pkg/assessor/group/group.go`
- Modify: `pkg/assessor/hosts/hosts.go`
- Modify: `pkg/assessor/user/user.go`
- Modify: `pkg/assessor/passwd/passwd.go`

**Interfaces:**
- Consumes: Task 1–5 の全 Produces。

- [ ] **Step 1: scan.go の deckoder import を新パスへ**

`pkg/scanner/scan.go` の import ブロックを次のとおり書き換える。`deckodertypes` エイリアスは廃し、dockle の `types` に一本化する（`types` は既に import 済み）。

変更前:
```go
	"github.com/goodwithtech/deckoder/analyzer"
	"github.com/goodwithtech/deckoder/extractor"
	"github.com/goodwithtech/deckoder/extractor/docker"
	deckodertypes "github.com/goodwithtech/deckoder/types"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/goodwithtech/dockle/pkg/assessor"
```
変更後:
```go
	"github.com/goodwithtech/dockle/pkg/analyzer"
	"github.com/goodwithtech/dockle/pkg/extractor"
	"github.com/goodwithtech/dockle/pkg/extractor/docker"

	"github.com/goodwithtech/dockle/pkg/types"

	"github.com/goodwithtech/dockle/pkg/assessor"
```
続いて本文の `deckodertypes.DockerOption` → `types.DockerOption`、`deckodertypes.FileMap` → `types.FileMap`、`deckodertypes.FilterFunc` → `types.FilterFunc` を置換する（該当箇所: `ScanImage` の引数 `dockerOption deckodertypes.DockerOption`、`var files deckodertypes.FileMap`、`createPathPermissionFilterFunc` の戻り値 `deckodertypes.FilterFunc`）。

- [ ] **Step 2: run.go の deckoder import を置換**

`pkg/run.go`: import 行 `deckodertypes "github.com/goodwithtech/deckoder/types"` を削除。`pkg/run.go` は既に dockle の `types` を import 済みか確認し、未 import なら `"github.com/goodwithtech/dockle/pkg/types"` を追加する。本文 `deckodertypes.DockerOption{` → `types.DockerOption{`。

Run（確認）: `grep -n '"github.com/goodwithtech/dockle/pkg/types"' pkg/run.go`
Expected: import が存在すること（なければ追加）。

- [ ] **Step 3: assessor.go と全アセッサの deckodertypes を置換**

以下のコマンドで、`deckodertypes` エイリアスを使う全ファイルについて (a) alias import 行を削除し (b) `deckodertypes.` を `types.` に置換する。これらのファイルはすべて `[]*types.Assessment` を返すため dockle の `types` を既に import しており、alias 削除後も `types` 参照は解決する。

```bash
FILES="pkg/assessor/assessor.go \
pkg/assessor/cache/cache.go \
pkg/assessor/contentTrust/contentTrust.go \
pkg/assessor/credential/credential.go \
pkg/assessor/manifest/manifest.go \
pkg/assessor/privilege/suid.go \
pkg/assessor/group/group.go \
pkg/assessor/hosts/hosts.go \
pkg/assessor/user/user.go \
pkg/assessor/passwd/passwd.go"
for f in $FILES; do
  sed -i '' \
    -e '/deckodertypes "github.com\/goodwithtech\/deckoder\/types"/d' \
    -e 's#deckodertypes\.#types.#g' \
    "$f"
done
```

- [ ] **Step 4: cache.go の utils import を新パスへ**

`pkg/assessor/cache/cache.go` の `"github.com/goodwithtech/deckoder/utils"` を `"github.com/goodwithtech/dockle/pkg/deckoder/utils"` へ置換する（Step 3 の sed は utils import 行を対象にしていないため個別対応）。

```bash
sed -i '' \
  -e 's#github.com/goodwithtech/deckoder/utils#github.com/goodwithtech/dockle/pkg/deckoder/utils#g' \
  pkg/assessor/cache/cache.go
```

- [ ] **Step 5: scan_test.go の deckodertypes を置換**

`pkg/scanner/scan_test.go`: alias import 行 `deckodertypes "github.com/goodwithtech/deckoder/types"` を削除し、`deckodertypes.DockerOption` → `types.DockerOption`（2 箇所: `option deckodertypes.DockerOption` フィールドと `deckodertypes.DockerOption{Timeout: time.Minute}`）へ置換する。`types` は既に import 済み。

```bash
sed -i '' \
  -e '/deckodertypes "github.com\/goodwithtech\/deckoder\/types"/d' \
  -e 's#deckodertypes\.#types.#g' \
  pkg/scanner/scan_test.go
```

- [ ] **Step 6: コード中に deckoder 参照が残っていないか確認**

Run: `grep -rn "goodwithtech/deckoder" pkg/ cmd/`
Expected: 出力なし。

- [ ] **Step 7: gofmt で import 整形**

Run: `gofmt -w pkg/`
Expected: エラーなし（alias 削除で空行が残った import ブロックを整形）。

- [ ] **Step 8: 全体ビルド・vet**

Run: `go build ./... && go vet ./...`
Expected: エラーなし。

- [ ] **Step 9: コミット**

```bash
git add pkg/
git commit -m "[機能] dockle の deckoder import を pkg/ 内部パスへ切り替え"
```

---

## Task 7: go.mod から deckoder を削除し全体検証・ライセンス確認

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: deckoder require を削除して tidy**

```bash
go mod edit -droprequire github.com/goodwithtech/deckoder
CGO_ENABLED=0 go mod tidy
```

`go mod tidy` が deckoder の direct 依存（`knqyf263/nested`, `GoogleCloudPlatform/docker-credential-gcr/v2`, `aws/aws-sdk-go`, `distribution/reference`, `opencontainers/go-digest`, `golang.org/x/xerrors` 等）を dockle の require へ昇格・正規化する。

- [ ] **Step 2: go.mod に deckoder が残っていないか確認**

Run: `grep -n "goodwithtech/deckoder" go.mod go.sum`
Expected: 出力なし。

- [ ] **Step 3: deckoder の LICENSE を持ち込んでいないか確認**

Run: `grep -rl "AFFERO" . --include="LICENSE*" 2>/dev/null; ls`
Expected: dockle リポジトリ内に AGPL の `LICENSE` が無いこと（dockle の `LICENSE` は Apache-2.0 のまま。deckoder の LICENSE はコピーしていないので新規出現しないはず）。

- [ ] **Step 4: CI 同条件で全テスト**

Run: `CGO_ENABLED=0 go test ./...`
Expected: 全パッケージ PASS。

- [ ] **Step 5: バイナリのビルド確認**

Run: `go build -o dockle cmd/dockle/main.go && ./dockle --version`
Expected: ビルド成功しバージョンが表示される。

- [ ] **Step 6: コミット**

```bash
git add go.mod go.sum
git commit -m "[機能] go.mod から deckoder 依存を削除し tidy"
```

---

## 完了条件

- `grep -rn "goodwithtech/deckoder" .`（`.git` 除く）が空。
- `CGO_ENABLED=0 go test ./...` が全 PASS。
- `go build -o dockle cmd/dockle/main.go` 成功。
- dockle の `LICENSE`（Apache-2.0）が維持され、AGPL 記述が持ち込まれていない。
