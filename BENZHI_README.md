# SignFlow

SignFlow 是电子合同签署与归档平台：法务/合同命名空间隔离管理合同台账，合同按版本
维护正文与修订记录，签署方按约定顺序依次签署并落盘签名证据，全部签署完成后合同
生效并可归档，归档文件带内容哈希防篡改，变更全程留痕，合同总量受命名空间配额
约束，签署与归档事件全部写入审计流水。

## 运行

```bash
go build -mod=vendor -o signflow ./cmd/signflow
./signflow -addr :8080 -root ./data
```

打开 http://localhost:8080/ 查看合同台账，/signing、/archive、/audit 分别查看
顺序签署、归档中心与审计变更页面。

## 测试

```bash
go test -mod=vendor ./...
go vet -mod=vendor ./...
```

## Docker

```bash
bash build_benzhi_docker.sh signflow linux/amd64
docker run --rm -p 8080:8080 signflow bash -c 'go run ./cmd/signflow -addr :8080 -root /tmp/data'
```
