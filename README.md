# Account service


## Pre requisites

- Docker
- Golang v1.14+
 
 
## Running App 

1. Export DB_PASSWORD
`export DB_PASSWORD="S3cretP@ssw0rd"` 

2. Bring up the mysql container using:

`make infra-local`

3. Run the migrations on the local DB 
 
`make setup`

4. Build and run the app container.

`make app`

5. Inspect logs using docker 

`docker logs account-service-go -f`

## Verifying the Functionality

Add Balance  Request
```shell script

curl -X POST \
  http://localhost:8888/account/credit \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  -d '{ "userId":"a900c144-9f25-4324-994f-451a7ac9d46d", "amount":2, "type":"subscription", "priority":2, "expiry":1605688837 }'
```

GET Credit logs request
```shell script
curl -X GET \
  'http://localhost:8888/account/logs?activity=Credit' \
 
```


