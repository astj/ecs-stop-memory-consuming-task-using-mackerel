# ecs-stop-memory-consuming-task-scheduler の要件

## 機能

特定の ECS Service に属するタスクのうち、もっともメモリ消費量の多いタスクを終了する、ということを最終的には実現したい。

## 技術的手段

「特定の ECS Service に属するタスクのうち、もっともメモリ消費量の多いタスクを選定する」という部分について、 mackerel.io のメトリックを利用する。

すなわち、指定された Mackerel.io の特定の Role に属するホスト (ECS タスクと 1:1 対応する) のうち、指定された名前(メモリ消費量に対応させる)のメトリックの値が一番大きい物を選定する。
このプロセスにより Mackerel 上のホストのメタデータから ECS タスクの task ID を取得できるので、後は AWS ECS の API 呼び出しによって当該のタスクを停止する、とする。

## インターフェース

（デーモンなどではなく）上記の作業を1回行って終了する CLI として実装する。
以下のようなパラメータを、環境変数とコマンドラインオプションの両方から指定できるようにする。

- dry_run かどうか
- mackerel 上のロール
- mackerel 上のメトリック名
- AWS の Profile, Region
- Mackerel の API キー
    - この環境変数名は、 conventional として `MACKEREL_APIKEY` とします

## 技術選定基準

- 標準ライブラリの採用を優先する
- 標準ライブラリの機能に不備がある場合は、一般的に知名度の高いライブラリを選定する

## 作業ログコーナー

ここに Copilot が作業ログを書く

### 2025年8月7日 - パラメータ処理の実装

- main.go にコマンドラインパラメータと環境変数の両方から設定を読み込む機能を実装
- 実装したパラメータ:
  - `dry-run` / `DRY_RUN`: ドライランモード
  - `mackerel-service` / `MACKEREL_SERVICE`: Mackerel上のサービス名
  - `mackerel-role` / `MACKEREL_ROLE`: Mackerel上のロール名
  - `mackerel-metric` / `MACKEREL_METRIC`: メモリ消費量に対応するメトリック名
  - `aws-profile` / `AWS_PROFILE`: AWSプロファイル名
  - `aws-region` / `AWS_REGION`: AWSリージョン
  - `mackerel-api-key` / `MACKEREL_API_KEY`: Mackerel APIキー
- コマンドラインオプションが環境変数より優先される仕様
- 必須パラメータのバリデーション機能を追加（mackerel-service, mackerel-role, mackerel-metric, mackerel-api-key, aws-region）
- APIキーのマスク表示機能を実装（セキュリティ考慮）
- 標準ライブラリの `flag` パッケージを使用（技術選定基準に準拠）
