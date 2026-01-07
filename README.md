# sact

さくらのクラウドの TUI アプリケーション｡さくっとサーバーを管理するよ｡

基本的に､サーバーの追加･削除は terraform でする前提で､オペレーションとしては停止と開始とかメインとする｡

## 構成要素

 * bubbletea
 * golang

## 使い方

### 環境変数

クレデンシャル情報は環境変数から取得します:

```bash
export SAKURA_ACCESS_TOKEN=your_token
export SAKURA_ACCESS_TOKEN_SECRET=your_secret
```

### 実行

```bash
# 通常起動（ログは標準エラー出力）
./sact

# ログをファイルに出力
./sact --log=/path/to/logfile
```

### 操作

- `z`: ゾーン切り替え (tk1a, tk1b, is1a, is1b, is1c)
- `r`: サーバー一覧の再読み込み
- `q` または `Ctrl+C`: 終了

### 設定ファイル

`~/.config/sact/config.toml` でデフォルトゾーンを設定できます:

```toml
default_zone = "tk1b"
```

## 実装方針

 * サーバー一覧の表示機能
 * ゾーンの切り替え機能｡tk1b, tk1a, is1a, is1b, is1c に対応
 * https://github.com/sacloud/iaas-api-go を使う

ここまでできてから､他のコンポーネント例えば switch+router, switch, dns, DBアプライアンス などのリソースも対応していく｡

