## My FIRST GO-WEB APP

sqlc docker
```bash
docker run --rm -v D:\project\study\go\project_bank:/src -w /src kjconroy/sqlc generate

```

create migrate
```bash
 migrate create -ext sql -dir db/migration -seq <name>
```
viper 配置环境的读取
```bash
https://github.com/deepin-community/golang-github-spf13-viper

go get -U github.com/spf13/viper
```

gin  go web
```bash
https://github.com/gin-gonic/gin
go get - U github.com/gin-gonic/gin
```


dbdocs
```bash
npm install -g dbdocs 
dbdocs build doc/db.dbml
```
