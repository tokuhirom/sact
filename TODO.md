# sact リソース対応状況

iaas-api-go v1.24.1 で利用可能なリソースの対応状況です。

## 実装済み (18リソース)

| リソース | ファイル | 説明 |
|---------|---------|------|
| Server | `internal/server.go` | サーバー |
| Switch | `internal/switch.go` | スイッチ |
| DNS | `internal/dns.go` | DNS |
| ProxyLB (ELB) | `internal/elb.go` | エンハンスドロードバランサ |
| GSLB | `internal/gslb.go` | 広域負荷分散 |
| Database | `internal/db.go` | データベースアプライアンス |
| Disk | `internal/disk.go` | ディスク |
| Archive | `internal/archive.go` | アーカイブ |
| Internet | `internal/internet.go` | ルーター |
| VPCRouter | `internal/vpcrouter.go` | VPCルーター |
| PacketFilter | `internal/packetfilter.go` | パケットフィルタ |
| LoadBalancer | `internal/loadbalancer.go` | 標準ロードバランサ |
| NFS | `internal/nfs.go` | NFSアプライアンス |
| SSHKey | `internal/sshkey.go` | SSH公開鍵 |
| AutoBackup | `internal/autobackup.go` | 自動バックアップ |
| SimpleMonitor | `internal/simplemonitor.go` | シンプル監視 |
| Bridge | `internal/bridge.go` | ブリッジ接続 |
| ContainerRegistry | `internal/containerregistry.go` | コンテナレジストリ |

## 未実装 - 中優先度 (特定用途で使うリソース)

| リソース | API名 | 説明 | ゾーン依存 |
|---------|------|------|-----------|
| CDROM | CDROMAPI | ISOイメージ | Yes |
| LocalRouter | LocalRouterAPI | ローカルルーター | No |
| MobileGateway | MobileGatewayAPI | モバイルゲートウェイ | Yes |
| SIM | SIMAPI | SIM | No |
| PrivateHost | PrivateHostAPI | 専有ホスト | Yes |
| Note | NoteAPI | スタートアップスクリプト | No |
| Interface | InterfaceAPI | NIC | Yes |
| EnhancedDB | EnhancedDBAPI | エンハンスドデータベース (TiDB) | No |
| AutoScale | AutoScaleAPI | オートスケール | No |
| CertificateAuthority | CertificateAuthorityAPI | マネージドPKI | No |
| ESME | ESMEAPI | 2要素認証 (SMS) | No |

## 未実装 - 低優先度 (参照専用・課金系・IP管理)

| リソース | API名 | 説明 | ゾーン依存 |
|---------|------|------|-----------|
| IPAddress | IPAddressAPI | IPv4アドレス管理 | Yes |
| IPv6Net | IPv6NetAPI | IPv6ネットワーク | Yes |
| IPv6Addr | IPv6AddrAPI | IPv6アドレス | Yes |
| Subnet | SubnetAPI | サブネット | Yes |
| License | LicenseAPI | ライセンス (Windows等) | No |
| Icon | IconAPI | アイコン | No |
| Bill | BillAPI | 請求情報 | No |
| Coupon | CouponAPI | クーポン情報 | No |
| SimpleNotificationDestination | SimpleNotificationDestinationAPI | 通知先 | No |
| SimpleNotificationGroup | SimpleNotificationGroupAPI | 通知グループ | No |

## 次のステップ

1. [x] LoadBalancer - 標準ロードバランサの対応
2. [x] NFS - NFSアプライアンスの対応
3. [x] PacketFilter - パケットフィルタの対応
4. [x] SSHKey - SSH公開鍵の対応
5. [x] AutoBackup - 自動バックアップの対応
6. [x] SimpleMonitor - シンプル監視の対応
7. [x] Bridge - ブリッジ接続の対応
8. [x] ContainerRegistry - コンテナレジストリの対応

## 実装パターン

新しいリソースを追加する際の基本パターン:

1. `internal/<resource>.go` にリソース構造体と API ラッパーを実装
2. `internal/client.go` の `ResourceType` に追加
3. `internal/model.go` でリソース一覧取得とビュー表示を実装
4. `internal/render.go` でテーブル表示ロジックを追加
