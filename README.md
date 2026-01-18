# sact

さくらのクラウドの TUI アプリケーション｡さくっとサーバーを管理するよ｡

基本的に､サーバーの追加･削除は terraform でする前提で､オペレーションとしては停止と開始とかメインとする｡

## 構成要素

 * bubbletea
 * golang

## 使い方

### クレデンシャル

usacloud のプロファイルから自動的に読み込みます。事前に usacloud で認証設定を行ってください:

```bash
usacloud config
```

環境変数での指定も可能です:

```bash
export SAKURACLOUD_ACCESS_TOKEN=your_token
export SAKURACLOUD_ACCESS_TOKEN_SECRET=your_secret
```

詳細は [usacloud のドキュメント](https://docs.usacloud.jp/usacloud/installation/start_guide/#configuration) を参照してください。

### 実行

```bash
# 通常起動（ログは標準エラー出力）
./sact

# ログをファイルに出力
./sact --log=/path/to/logfile
```

### 操作

- `t`: リソースタイプ切り替え (Server, Switch, DNS, ELB, GSLB, DB)
- `z`: ゾーン切り替え (tk1a, tk1b, is1a, is1b, is1c)
- `r`: 一覧の再読み込み
- `Enter`: 詳細表示
- `/`: 検索
- `n`/`N`: 次/前の検索結果
- `j`/`k` または `↑`/`↓`: カーソル移動
- `q` または `Ctrl+C`: 終了

### デフォルトゾーン

デフォルトゾーンは usacloud プロファイルの `Zone` フィールドで設定できます。
プロファイルの設定ファイルは `~/.usacloud/{プロファイル名}/config.json` に保存されます。

デフォルトゾーンが設定されていない場合は `tk1b` が使用されます。

## 実装方針

 * サーバー一覧の表示機能
 * ゾーンの切り替え機能｡tk1b, tk1a, is1a, is1b, is1c に対応
 * https://github.com/sacloud/iaas-api-go を使う

ここまでできてから､他のコンポーネント例えば switch+router, switch, dns, DBアプライアンス などのリソースも対応していく｡

