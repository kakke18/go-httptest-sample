# go-httptest-sample
アプリケーションサーバと認証サーバがある場合のアプリケーションサーバのテストサンプル

## 実行方法
### アプリケーションサーバと認証サーバをどちらも起動させる場合
```shell
# アプリケーションサーバと認証サーバを起動
docker compose up -d

# アプリケーションサーバのテストを実行
cd app/httptest
go test -v .
```

### アプリケーションサーバのみ起動させる場合
```shell
# アプリケーションサーバを起動
cd app
go run main.go

# アプリケーションサーバのテストを実行
cd httptest
go test -v . # 認証サーバが立ち上がってないので、失敗するはず
```

### アプリケーションサーバをENV=testで起動させる場合
```shell
# アプリケーションサーバをENV=testで起動
cd app
ENV=test go run main.go

# アプリケーションサーバのテストを実行
cd httptest
go test -v . # 認証サーバはmockを使ってるので、成功するはず
```

## 備考
テスト結果がキャッシュされる恐れがあるので、テストの前には以下のコマンドを実行してキャッシュ削除した方が良い
```shell
go clean -testcache
```
