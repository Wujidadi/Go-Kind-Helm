# Go-Kind-Helm

本專案示範如何使用 Go 語言應用程式搭配 Kind + Helm + Ingress 於本機環境建立 K8s 部署。

---

## 目錄結構

```
.
├── Dockerfile                      # 定義 Go 應用容器環境的 Dockerfile
├── README.md                       # 本專案說明文件
├── charts                          # 儲存 Helm charts 的資料夾
│   ├── go                          # Go 應用程式專用的 Helm chart
│   │   ├── Chart.yaml              # Helm chart 的基本描述與版本資訊
│   │   ├── charts                  # 子 chart（通常可忽略，空資料夾）
│   │   ├── templates               # Helm chart 的 K8s YAML 範本檔
│   │   │   ├── NOTES.txt           # Helm 安裝後提示訊息
│   │   │   ├── _helpers.tpl        # 自定義函式與共用 template
│   │   │   ├── deployment.yaml     # 部署設定檔 (Deployment)
│   │   │   ├── hpa.yaml            # 自動水平擴展設定 (HPA)
│   │   │   ├── ingress.yaml        # Ingress 設定檔
│   │   │   ├── pv.yaml             # Persistent Volume 設定檔
│   │   │   ├── pvc.yaml            # Persistent Volume Claim 設定檔
│   │   │   ├── service.yaml        # Service 設定檔
│   │   │   ├── serviceaccount.yaml # ServiceAccount 設定檔
│   │   │   └── tests               # 測試用途目錄（可放 Helm 測試檔）
│   │   ├── values-dev.yaml         # 開發環境使用的參數設定檔
│   │   ├── values-prod.yaml        # 生產環境使用的參數設定檔
│   │   └── values.yaml             # 預設參數設定檔（可被 dev/prod 覆寫）
│   └── postgresql                  # PostgreSQL Helm chart (Bitnami 官方)
├── golang                          # 掛載 Go GOPATH 的本機卷宗資料夾
│   ├── bin                         # Go 安裝模組產生的可執行檔 (容器中 /go/bin)
│   ├── pkg                         # Go 模組快取 (容器中 /go/pkg)
│   └── src                         # Go 模組原始碼目錄 (容器中 /go/src)
├── kind-config.yaml                # 本地 Kind Kubernetes 集群設定檔
├── makefile                        # 提供開發與部署自動化指令的 makefile
├── src
│   └── app                         # Go 主應用程式原始碼放置目錄
│       ├── go.mod                  # Go module 設定檔
│       ├── go.sum                  # Go module 檢查碼 (鎖定依賴版本)
│       └── main.go                 # Go 主程式進入點
└── zsh
    └── root.zsh_history            # 容器內 Zsh 指令歷史紀錄檔（掛載用）
```

---

## Kind

| 容器埠號 | 宿主埠號 |
| :------: | :------: |
|    80    |   840    |
|   443    |   843    |

---

## Helm

```bash
helm search repo bitnami/postgresql # 查看 bitnami/postgresql 最新版
helm search repo bitnami/redis # 查看 bitnami/redis 最新版
```

其他指令都寫在 `make helm` 了（見 `makefile`）

---

## Kubernetes Secrets

寫在 `make k8s-secret` 了（見 `makefile`）
> 跑生產模式前一定要先產生 Kubernetes Secrets (已加到 `make deploy-redis-prod` 前面)

---

## 啟動服務及檢測

```bash
# 若未建置過 Helm 須先執行
make helm

# 一鍵啟動開發環境
make startup-dev  # 也可以用 make startup

# 檢查 Go App 是否正常
kubectl get pods -n go-helm
# 若此處不正常，可先檢查 log
kubectl logs $(kubectl get pods -n go-helm -o jsonpath="{.items[0].metadata.name}") -n go-helm
# 若 log 無輸出，再用 describe 查看詳情
kubectl describe pod $(kubectl get pods -n go-helm -o jsonpath="{.items[0].metadata.name}") -n go-helm

# 檢查 PostgreSQL 是否正常
kubectl get pods -n postgres-go-helm

# 檢查 Redis 是否正常
kubectl get pods -n redis-go-helm
kubectl get svc -n redis-go-helm

# 查看 Ingress
kubectl get ingress -n go-helm
make show-all | grep ingress
```

---

## makefile 支援指令

| 指令名稱                    | 說明                                                    |
| --------------------------- | ------------------------------------------------------- |
| `make helm`                 | Helm 環境建置                                           |
| `make k8s-secret`           | 使用 Kubernetes Secrets 建立 Redis 生產環境所需的密碼   |
| `make save-tag`             | 產生 build tag                                          |
| `make zsh-history`          | 建立 Zsh 歷史記錄掛載點                                 |
| `make build`                | 建立 Docker 映像檔                                      |
| `make load`                 | 載入映像檔到 Kind                                       |
| `make deploy`               | 預設部署 Go 應用（使用開發環境設定）                    |
| `make deploy-dev`           | 部署 Go 應用（開發環境）                                |
| `make deploy-prod`          | 部署 Go 應用（生產環境）                                |
| `make deploy-postgres-dev`  | 部署 PostgreSQL（開發環境）                             |
| `make deploy-postgres-prod` | 部署 PostgreSQL（生產環境）                             |
| `make deploy-redis-dev`     | 部署 Redis（開發環境）                                  |
| `make deploy-redis-prod`    | 部署 Redis（生產環境）                                  |
| `make create-cluster`       | 建立本地 K8s Cluster                                    |
| `make delete-cluster`       | 刪除 Cluster                                            |
| `make install-ingress`      | 安裝 Ingress-nginx 並等候就緒，結束後提示 URL 範本      |
| `make uninstall-ingress`    | 清除 Ingress-nginx                                      |
| `make reload`               | 重建映像、更新並部署                                    |
| `make restart`              | 滾動重啟 Pod                                            |
| `make startup`              | 一鍵完整重部署（預設模式）                              |
| `make startup-dev`          | 一鍵完整重部署（開發環境版）                            |
| `make startup-prod`         | 一鍵完整重部署（生產環境版）                            |
| `make clean`                | 清除 Go Helm release                                    |
| `make clean-postgres`       | 清除 PostgreSQL Helm release                            |
| `make clean-redis`          | 清除 Redis Helm release                                 |
| `make show-all`             | 一鍵顯示所有 Helm Release、Pods、Services、Ingress 狀態 |
| `make forward-db`           | 將 PostgreSQL DB 埠暫時暴露給宿主                       |
| `make forward-redis`        | 將 Redis 埠暫時暴露給宿主                               |
| `make forward-all`          | 同時 Forward DB 及 Redis                                |
| `make list-clusters`        | 查看 Kubernetes 集群狀態（通用指令）                    |
| `make list-pods`            | 查看 Kubernetes Pod 狀態                                |
| `make get-containers`       | 列出 Pod 的中容器名稱                                   |
| `make shell`                | 自動偵測 Pod + Container 並進入 shell                   |
| `make shell-pick`           | 多 Pod 手動選擇 shell                                   |
| `make shutdown`             | 關閉並刪除所有 Helm release 與 Cluster                  |

---

## 建議訪問路由

由於本地 Ingress Controller 的 Service 非 LoadBalancer 或 HostNetwork，必須使用域名 + 埠號的方式訪問：

| 環境     | 訪問網址                                                              |
| -------- | --------------------------------------------------------------------- |
| 開發環境 | http://go-app-dev.localhost:840<br />https://go-app-dev.localhost:843 |
| 生產環境 | https://go-app.localhost:843<br />http://go-app.localhost:840         |

> 特別注意：頂級域名 `local` 會有 mDNS (Multicast DNS) 阻塞問題導致緩慢，在 macOS 上載入頁面可能會慢至 10 秒以上，即使是本機測試環境亦不應使用。

### 驗證

訪問根路由應顯示以下訊息：

```
✅ Go server is running.
```

訪問 `/healthz` 應顯示以下訊息：

```
ok
```

訪問 `/db` (開發環境) 或 `/db?prod` (生產環境) 應顯示以下訊息：

```
✅ 成功連線到 PostgreSQL！
PostgreSQL 版本：PostgreSQL 17.4 on aarch64-unknown-linux-gnu, compiled by gcc (Debian 12.2.0-14) 12.2.0, 64-bit
```

訪問 `/cache` (開發環境) 或 `/cache?prod` (生產環境) 應顯示以下訊息：

```
✅ Redis 測試成功！取得值: HaHaHa Redis! I am Taras!
```

---

## 小提醒

- 請在 `/etc/hosts` 加入對應域名指向 127.0.0.1
- 請確認 Docker、Kind、Helm 均已安裝
- 如需要 PostgreSQL 客戶端，可在容器或宿主上安裝 `psql`
- 若一開始的啟動命令尚未掛載具體服務而是像 `tail -f /dev/null`，
  > 比方說，在你的 `charts/go/values.yaml` 或 `charts/go/templates/deployment.yaml` 中這樣設定：
  > ```yaml
  > command: ["tail"]
  > args: ["-f", "/dev/null"]
  > ```
  要先把 `charts/go/values.yaml` 中的 `livenessProbe` 和 `readinessProbe` 兩個區塊註解或刪除掉（即使設成空物件 `{}`
  ）也不行  
  因為它們會一直監聽容器內的 HTTP 服務，此時會因為不斷監聽失敗，最終導致 `CrashLoopBackOff` 狀態。
