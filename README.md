# Test-minio-server

# run minio podman (requis pour DEV ou PROD)

Terminal 1
```
sudo podman run -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
```

L'interface web se trouve ici: http://127.0.0.1:9001

# DEV

Dans config.yaml, changer "Endpoint" par : ip:port

Terminal 2
```
cd test-minio-server
go run main.go
```

# PROD

Terminal 2
```
docker-compose up -d
```



Terminal 3
```
docker logs --follow <containerID_app1>
```

Terminal 4
```
docker logs --follow <containerID_app2>
```

# Requête HTTP

StoreFile
```
curl -v -X POST http://127.0.0.1:1323/storeFile -H 'Content-Type: application/json' -H "Accept: " -d '{"bucketName":"testbucket","location":"pamatai","s3Filename":"testS3","data":"ceci est un test."}'
```

GetFile
```
curl -v -X POST http://127.0.0.1:1323/getFile -H 'Content-Type: application/json'  -H 'Accept: text/plain' -d '{"bucketName":"testbucket","s3Filename":"testS3"}'
```
