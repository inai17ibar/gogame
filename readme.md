
# gogame

TechTrain MISSION ゲームAPI入門仕様

[Swagger](https://editor.swagger.io/)に、api-document.yamlをドラッグすると、API仕様が確認できる。

# Requirements

* go(version go1.15.5 darwin/amd64)
* docker(version 20.10.0, build 7287ab3)
* docker-compose(version 1.27.4, build 40524192)

# Install

MacOSX
```bash
brew install go

brew install docker
#1.25.4 の部分は任意のバージョン
sudo curl -L https://github.com/docker/compose/releases/download/1.27.4/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
# macならDocker Desktopをいれたらcomposeもはいる

# go libraries
go get -u github.com/gorilla/mux
go get -u github.com/go-sql-driver/mysql
```

# Build
```bash
go build
# buildしてそのまま実行
go run
```

# Usage
MySQLサーバの起動
```bash
cd docker-compose
# start myaql server with docker(一回目は時間がかかる)
docker-compose up -d
# コンテナがある場合
docker-compose start

# serviceNameはmysql. Dockerコンテナ内でコマンド実行したい場合
docker-compose exec mysql /bin/bash 
mysql -u root -p 

# 動作確認。docker-composeにPasswordを記載
mysql -u root -p -h 0.0.0.0 -P 3306 sample
```

goのWebAPIサーバの起動
```bash
# start web-api server
./gogame
```

終了時
```bash
# コンテナを残す場合
docker-compose stop

#cnfの内容やSQLを修正した場合
docker-compose down
docker-compose up -d

#./gogameは、Ctrl+Cで終了
```

Access Front webpage (仮)
[http;//localhost:8080](http;//localhost:8080)


