# Go API Server with 〒 郵便番号検索
このリポジトリは、Go 言語で実装された郵便番号検索 API サーバーのサンプルです。外部 API (HeartRails Geo API) を利用して郵便番号から住所情報を取得し、そのアクセス履歴を MySQL にログとして保存する機能を持っています。Docker / Docker Compose を使用することで、開発環境を簡単に構築・起動できます。Features/ (GET): 郵便番号を指定して住所情報を取得/hello (GET): シンプルなサンプル API/address/access_logs (GET): 郵便番号ごとのアクセスログを集計して表示アクセスログ管理: MySQL へのアクセスログ保存機能簡単起動: Docker によるコンテナ化と Docker Compose による一括起動 Requirementsこのプロジェクトをローカルで実行するには、以下のツールが必要です。Docker (v20+ 推奨)Docker Compose (v1.29+ 推奨)Go (v1.25+): ローカルでビルド・実行する場合Note: Docker Compose を使用する場合、MySQL サーバーはコンテナ内で自動的に立ち上がります。Usage1. リポジトリのクローンBashgit clone <リポジトリURL>
cd <リポジトリ名>
2. 環境変数ファイル (.env) の作成.env.example を参考に、プロジェクトルートに .env ファイルを作成し、環境変数を設定します。Bash# .env ファイルの内容 (例)

# MySQL 接続情報
MYSQL_ROOT_PASSWORD=root123
MYSQL_PASSWORD=pass1234
# ...その他の環境変数
注意:公開リポジトリには .env.example のみを置き、.env ファイルは .gitignore に追加してコミットしないでください。本番環境では、必ずこれらのパスワードを強力でセキュアな値に変更してください。3. Docker Compose での起動以下のコマンドで、Go API サーバーと MySQL データベースを同時にビルド・起動します。Bashdocker-compose up --build
サービス名ポート (ホスト側)概要app (Go API)8081API サーバー本体。DB 接続は DB_DSN 環境変数で設定されます。db (MySQL)3307データベース。初期 DB 名は mydb、ユーザー名は go_user です。4. API の実行例サーバー起動後、以下の curl コマンドで API をテストできます。Hello APIBashcurl http://localhost:8081/hello
郵便番号検索 API郵便番号 5016121 の住所情報を取得し、アクセスログを記録します。Bashcurl "http://localhost:8081/?postal=5016121"
アクセスログ集計 API記録されたアクセスログを集計して表示します。Bashcurl http://localhost:8081/address/access_logs
# Project Structure.
├── main.go               # Go API サーバーのメイン処理
├── go.mod                # Go Modules 定義ファイル
├── go.sum                # Go Modules チェックサムファイル
├── Dockerfile            # Go API サーバーの Docker イメージ定義
├── docker-compose.yml    # Docker Compose 設定ファイル
├── init.sql              # MySQL コンテナ起動時に実行される初期化 SQL
├── .env.example          # 環境変数設定ファイルのサンプル
└── README.md             # このファイル
 Security NotesDB パスワードの管理: データベースのパスワードは .env ファイルで管理し、GitHub などの公開リポジトリには絶対に含めないようにしてください。本番環境での設定: 本番環境にデプロイする際は、.env の値を必ず、外部から推測されにくいよりセキュアな値に変更してください。