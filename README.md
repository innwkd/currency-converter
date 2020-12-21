## Currency converter
Not for production. For practise purposes only

#### Supported currency providers
- fake
- https://exchangeratesapi.io/

#### Available env configs
- `DEBUG`= 0 | 1
- `APP_REST_PORT`= 12345
- `APP_RPC_PORT`= 4444
- `APP_CACHE_DURATION_MIN`= 60
- `APP_RATE_PROVIDER`= fake | exchangeratesapi.io
- `STORAGE_TYPE`= redis | memory
- `STORAGE_REDIS_HOST`= redis
- `STORAGE_REDIS_PORT`= 6379

#### Run
`docker-compose up --scale backend=3`

#### Requests
##### REST
- Convert
```
curl --request GET --url 'http://localhost:10000/convert?amount=10&from=USD&to=EUR'
```
```
{
   "result":"8.869179601",
   "currency_rate":{
      "pair":{
         "from":"USD",
         "to":"EUR"
      },
      "value":"0.8869179601",
      "provider":"exchangeratesapi.io",
      "updated_at":"2019-03-13T16:49:05.315733+03:00"
   }
}
```

- Stats
```
curl --request GET --url http://localhost:10001/stat
```
```
{
   "available_pair":[
      {
         "from":"EUR",
         "to":"USD"
      },
      {
         "from":"USD",
         "to":"EUR"
      }
   ],
   "cached_rates":[
      {
         "pair":{
            "from":"USD",
            "to":"EUR"
         },
         "value":"0.8869179601",
         "provider":"exchangeratesapi.io",
         "updated_at":"2019-03-13T16:49:05.315733+03:00"
      }
   ],
   "cache_duration":3600
}
```
##### jRPC
- Convert
```
curl --request POST \
  --url http://localhost:11001/rpc \
  --header 'Content-Type: application/json' \
  --data '{"method":"Converter.Convert","params":[{"pair":{"from":"USD","to":"EUR"},"amount":"10"}],"id":1}'
```
```
{
   "result":{
      "result":"8.869179601",
      "currency_rate":{
         "pair":{
            "from":"USD",
            "to":"EUR"
         },
         "value":"0.8869179601",
         "provider":"exchangeratesapi.io",
         "updated_at":"2019-03-13T16:49:05.315733+03:00"
      }
   },
   "error":null,
   "id":1
}
```
