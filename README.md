## Currency converter
Not for production. For practise purposes only

## Supported currency providers
- fake
- https://exchangeratesapi.io/

#### Run
`docker build --tag yddmat/currency-converter .`

`docker run -it --rm -p 12345:8080 yddmat/currency-converter`

#### Requests
- Convert
```
curl --request GET --url 'http://localhost:12345/convert?amount=10&from=USD&to=EUR'
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
curl --request GET --url http://localhost:12345/stat
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
