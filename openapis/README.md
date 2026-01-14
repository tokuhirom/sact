* モニタリングスイートの openapi 定義が実際のレスポンスと違う部分があるため書き換えている｡
* モニタリングスイートは ogen でクライアントを生成出来なかったため､oapi-codegen を利用している

## monitoring-suite.json の変更点

オリジナルは `monitoring-suite.orig.json` として保存している。

### `id` フィールド: string → integer

API は数値を返すが、スキーマでは string と定義されていた。

対象スキーマ:
- LogStorage, MetricsStorage, TraceStorage
- WrappedLogStorage, WrappedMetricsStorage, WrappedTraceStorage
- LogRouting, MetricsRouting, NotificationRouting
- PatchedLogRouting, PatchedMetricsRouting, PatchedNotificationRouting
- WrappedLogRouting, WrappedMetricsRouting

### `resource_id` フィールド: integer → string

API は文字列を返すが、スキーマでは integer と定義されていた。

対象スキーマ:
- LogStorage, MetricsStorage, TraceStorage
- WrappedLogStorage, WrappedMetricsStorage, WrappedTraceStorage
- LogRouting, MetricsRouting
- 各種パスパラメータ・クエリパラメータ

### `log_storage_id`, `metrics_storage_id` フィールド

- integer → string に変更
- `writeOnly: true` を削除（API レスポンスに含まれるため）

対象スキーマ:
- LogRouting, MetricsRouting
- PatchedLogRouting, PatchedMetricsRouting
- WrappedLogRouting, WrappedMetricsRouting
