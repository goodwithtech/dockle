# deckoder統合 設計書

作成日: 2026-07-02
対象: `github.com/goodwithtech/deckoder` を dockle 本体へ統合し、外部依存を廃止する

## 背景と目的

dockle は画像抽出（Dockerデーモン／リモートレジストリ／tarアーカイブからのファイル取得）にコンパニオンライブラリ `github.com/goodwithtech/deckoder` を利用している。deckoder は実質 dockle 専用に維持されており、独立パッケージとして分離する必要性が薄い。別リポジトリ管理の運用コストを削減するため、deckoder を dockle 本体へ取り込み、外部依存を廃止する。

deckoder はテストを除き約1527行の Go コードで、`analyzer` / `extractor` / `types` / `utils` の構成。依存グラフは `types` を葉とする綺麗な木構造で、循環importの懸念はない。

```
types            (葉: 内部依存なし)
utils            → types
extractor        → types
extractor/image  → types
extractor/image/token/ecr → types
extractor/image/token/gcr → types
extractor/docker → extractor/image, token/ecr, token/gcr, types
analyzer         → extractor, extractor/docker, types
```

## 決定事項

1. **配置**: pkg/ 配下へ分散統合する（`internal/` ではなく `pkg/`）。
2. **types**: deckoder の `types` を dockle の `pkg/types` へ完全マージし、`deckodertypes` エイリアスを廃止する。
3. **ライセンス**: deckoder は AGPL-3.0、dockle は Apache-2.0。著作権者（goodwithtech）による再ライセンスとして、取り込むコードを Apache-2.0 扱いとする。dockle 側のライセンスは現状（Apache-2.0）を維持する。
4. **進め方**: 本設計書 → 実装プラン（writing-plans）→ 実装のフルプロセス。

## 配置マッピング

| deckoder（元） | 統合先（dockle） |
|---|---|
| `types/const.go`, `types/docker.go`, `types/error.go`, `types/filter.go`, `types/image.go` | `pkg/types/` へマージ |
| `utils/` | `pkg/deckoder/utils/` |
| `analyzer/` | `pkg/analyzer/` |
| `extractor/` | `pkg/extractor/` |
| `extractor/docker/` | `pkg/extractor/docker/` |
| `extractor/image/`（`token/ecr`, `token/gcr`, `mock_*` 含む） | `pkg/extractor/image/...` |
| `cmd/deckoder/main.go` | 破棄（不要なCLIエントリ） |

### pkg/types へのマージ詳細

deckoder types が提供する宣言を dockle `pkg/types` に取り込む。dockle 既存の `pkg/types` とはファイル名が衝突するため、以下のようにリネームして配置する（型名の衝突はない）。

| deckoder の宣言 | 新規ファイル | 備考 |
|---|---|---|
| `FileMap`, `FileData`, `FilterFunc` | `pkg/types/filemap.go` | 旧 `types/filter.go` |
| `DockerOption` | `pkg/types/docker_option.go` | 旧 `types/docker.go`。dockle に `docker` 関連ファイルはないが命名を明確化 |
| `FilePath`, `ImageReference`, `ImageInfo`, `LayerInfo`, `ImageDetail` | `pkg/types/imageref.go` | 旧 `types/image.go`（dockle 既存 `image.go` と衝突するためリネーム） |
| `ImageJSONSchemaVersion`, `LayerJSONSchemaVersion` | `pkg/types/const.go` | 旧 `types/const.go` |
| `InvalidURLPattern`, `ErrNoRpmCmd` | `pkg/types/error.go` へ追記 | dockle 既存 `error.go` に統合。dockle は `errors`、deckoder は `xerrors` を使用 → `errors` に統一するか判断（実装プランで確定） |

**注意点**:
- dockle 既存 `pkg/types/image.go` は Docker 画像 config スキーマ（`Image`, `V1Image`, `Config`, `HealthConfig`, `History`）を定義。deckoder の `ImageReference` 等とは別物であり型名衝突はない。
- `pkg/types` は葉パッケージのまま（deckoder types は dockle 内の何もimportしない）。`pkg/extractor` / `pkg/analyzer` が `pkg/types` をimportする一方向依存となり、循環importは発生しない。

## import 書き換え

`github.com/goodwithtech/deckoder/*` を新パスへ一括置換する。影響を受ける dockle 側の呼び出し箇所:

- `pkg/run.go` — `deckodertypes.DockerOption` → `types.DockerOption`
- `pkg/scanner/scan.go` — `analyzer`, `extractor`, `extractor/docker`, `deckodertypes` → 新パス
- `pkg/scanner/scan_test.go` — `deckodertypes` → `types`
- `pkg/assessor/assessor.go` — `deckodertypes.FileMap` → `types.FileMap`
- 各アセッサ（`contentTrust`, `cache`, `manifest`, `passwd`, `credential`, `group`, `hosts`, `privilege`, `user`）— `deckodertypes.FileMap` → `types.FileMap`
- `pkg/assessor/cache/cache.go` — `deckoder/utils` → `pkg/deckoder/utils`

`deckodertypes` エイリアスは全廃し、dockle 既存の `types`（`github.com/goodwithtech/dockle/pkg/types`）に一本化する。これにより各アセッサの `Assess(fileMap deckodertypes.FileMap)` が `Assess(fileMap types.FileMap)` となり、`FileMap` と `Assessment` を単一パッケージから参照できる。

## go.mod / 依存整理

- `require github.com/goodwithtech/deckoder v0.0.6` を削除。
- deckoder の direct 依存を dockle の require へ昇格させる（現状 dockle の indirect にほぼ含まれる）:
  - `github.com/knqyf263/nested`（utils で使用、現在 indirect）
  - `github.com/GoogleCloudPlatform/docker-credential-gcr/v2`, `github.com/aws/aws-sdk-go`, `github.com/distribution/reference`, `github.com/opencontainers/go-digest`, `golang.org/x/xerrors`（extractor 系で使用）
  - `github.com/containers/image/v5` は既に dockle の direct 依存。
- 最終的に `go mod tidy` で direct/indirect を正規化する。

## テスト

deckoder のテスト・モックも移送し、回帰を担保する:

- `extractor/image/mock_image.go`, `mock_image_closer.go`, `mock_image_source.go`, `mock_registry.go` → `pkg/extractor/image/`
- `extractor/image/image_test.go`, `token_test.go` → `pkg/extractor/image/`
- `extractor/image/token/ecr/ecr_test.go`, `token/gcr/gcr_test.go` → 対応先
- `extractor/docker/docker_test.go` → `pkg/extractor/docker/`
- `utils/utils_test.go` → `pkg/deckoder/utils/`

検証コマンド:

```bash
go build -o dockle cmd/dockle/main.go
CGO_ENABLED=0 go test ./...
```

CI 同等の `CGO_ENABLED=0 go test ./...` が全て通ることを完了条件とする。

## ライセンス作業

- deckoder コードを Apache-2.0 として取り込む方針に従い、AGPL 由来の記述（`LICENSE` ファイル等）は dockle へ持ち込まない。
- 取り込んだファイルに AGPL ライセンスヘッダがあれば除去する（deckoder のソースにヘッダがあるか実装時に確認）。
- dockle の `LICENSE`（Apache-2.0）は変更しない。

## スコープ外（YAGNI）

- deckoder の CLI（`cmd/deckoder`）は取り込まない。
- deckoder リポジトリ側のアーカイブ／削除は本作業の対象外（統合完了後に別途判断）。
- 抽出ロジック自体のリファクタリング・機能追加は行わない（コピー＋import書き換えに徹する）。

## リスクと対処

- **依存の抜け漏れ**: `go mod tidy` と `go build` / `go test` で機械的に検出。
- **モックの可視性**: deckoder 内でパッケージをまたぐモック参照がある場合、移送後のパッケージ境界で export 状態を確認する。
- **xerrors vs errors**: `pkg/types/error.go` 統合時に import が二重化しないよう、`errors` へ寄せるか `xerrors` を残すか実装プランで確定する。
