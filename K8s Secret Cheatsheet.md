# Kubernetes Secret 管理指令速查表（Cheat Sheet）

## 📜 查看 Secret

| 指令                                                                         | 說明                                    |
| ---------------------------------------------------------------------------- | --------------------------------------- |
| `kubectl get secrets`                                                        | 列出所有 Secret                         |
| `kubectl describe secret <name>`                                             | 查看指定 Secret 詳細資訊（含 Metadata） |
| `kubectl get secret <name> -o yaml`                                          | 以 YAML 格式查看 Secret 原始內容        |
| `kubectl get secret <name> -o jsonpath="{.data.<key>}" \| base64 -d && echo` | 將指定 key 的值解碼後顯示               |

---

## ✏️ 編輯 Secret（僅限 metadata 建議使用）

```bash
kubectl edit secret <name>
```

> ⚠ 若要修改內容（data），建議刪除重建或使用 patch 指令。

---

## 🔄 更新（Patch） Secret

```bash
kubectl patch secret <name> \
  --type='json' \
  -p='[{"op": "replace", "path": "/data/<key>", "value": "'$(echo -n 新值 | base64)'"}]'
```

範例：
```bash
kubectl patch secret redis-secret \
  --type='json' \
  -p='[{"op": "replace", "path": "/data/redis-password", "value": "'$(echo -n newpassword | base64)'"}]'
```

---

## ❌ 刪除 Secret

```bash
kubectl delete secret <name>
```

---

## ✅ 小提示
- Secret 的值在 K8s 中預設是 base64 編碼。
- 可以搭配 `-n <namespace>` 查看/管理指定 Namespace 下的 Secret。
- 機密資料最好不要直接存在 Git 倉庫，建議透過 secret manager 工具或 CI/CD 注入。

---

> 建議：如果 secret 頻繁變更，可以透過 Helm 的 `--set` 或 `--set-file` 參數搭配 secret 模板動態建立。
