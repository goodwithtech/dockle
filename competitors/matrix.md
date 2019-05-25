主にセキュリティまわり。

# 単語の定義

- イメージ : コンテナイメージ, docker images ででてくるやつ
- コンテナ : 動いているコンテナ, docker ps ででてくるやつ


# 資料
- [公式ページ](https://www.docker.com/legal/security)
- [Introduction to Container Security](https://d3oypxn00j2a10.cloudfront.net/assets/img/Docker%20Security/WP_Intro_to_container_security_03.20.2015.pdf)(pdf)


# マトリクス

## 検証できる内容

|  | 実行環境チェック(Docker Engine) | OSパッケージチェック | App Dependencyチェック | ウイルス | コンテナ内の設定 |
|-- | -- | -- | -- | -- | --|
|Docker Bench for Security | ○ |   |   |   |  
|Dagda | ○ | ○ |   | ○ |  |
|Lynis | ○ |   |   |   | ○ |
|Hadolint |   |   |   |   | ○ |
|Clair |   | ○ |   |   |  |
|Trivy |   | ○ | ○ |   |  |
|Anchore |   | ○ | ○ | ○ | ○ |
|OpenSCAP |   | ○ |   |   |  |
|ClamAV |   |   |   | ○ |  |
|Dockscan |   | ○ |   |   |  |

## 検知に必要なもの


|  | 実行環境で動作 | イメージのみ | docker run | Dockerfile |
|-- | -- | -- | -- | --|
|Docker Bench for Security | ○ |  ○ | ○ | ○ |
|Dagda | ○ | ○ | ○ |  |
|Lynis | ○ |   | ○ | ○ |
|Hadolint |   |   |   | ○ |
|Clair |   | ○ | ○ |  |
|Trivy |   | ○ |   |  |
|Anchore |   | ○ | ○ | ○ | 
|OpenSCAP | ○ | ○ |   | | 
|ClamAV |   | ○ | ○ |  |
|Dockscan |   | ○ |   | | 


# 各ツールの詳細

## [Docker Bench for Security](https://github.com/docker/docker-bench-security)
実行環境で動作
Docker Engine側の脆弱性/設定をチェック

## [Dogda](https://github.com/eliasgranderubio/dagda)
python製。イメージ/コンテナスキャンに対し、CVE, ウイルスが検知でき、実行環境の検証もできる。

## Lyon : 作成中
イメージのみスキャン
- セキュリティ(rootで実行してないか)
- CMD, ENTRYPOINTが複数ないか
- コンテナのベストプラクティスに乗っているか

## [Lynis](https://github.com/CISOfy/lynis)
DockerfileをLint。
実行環境の設定も確認できる。
docker run内で走らせれば、一応コンテナ内の環境もチェックできる。

## [Hadolint](https://github.com/hadolint/hadolint)
DockerfileをLint。

## [Clair]()
パッケージスキャン。コンテナとイメージ対象。

## [Anchore](https://anchore.com/)
幅広く対応できる。
どこまでできるかよくわからない。

## [Aqua]()
イメージのみスキャン

## [Trivy]()
イメージのみスキャン

## [OpenSCAP]()
コンテナ, イメージスキャン。
RedHatに強い?

## [Vuls]()
パッケージの検証。
コンテナスキャンも可能。
もうすぐイメージスキャン対応。

## [Dockscan](https://github.com/kost/dockscan)
ruby製。コンテナスキャン

## [Dagda](https://github.com/eliasgranderubio/dagda)
python製。イメージ/コンテナスキャンに対し、CVE, ウイルスが検知でき、実行環境の検証もできる。
ウイルススキャンは内部でClamAVを利用。
高機能っぽいけど、あまり伸びてない?

## [ClamAV](https://www.clamav.net/)
イメージのtarファイルをもとにスキャン可能。
