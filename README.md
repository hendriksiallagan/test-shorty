# test-shorty
Amarha Shorten Test Using GO

## Installation
- Download & Install docker
- git clone git@github.com:hendriksiallagan/test-shorty.git
- Run command
```sh
docker-compose up --build
```
- Import API Collection on file Amartha Test Shorten.postman_collection.json in your postman
- Enjoy test the endpoints

## Endpoint List

### POST /shorten
```sh
curl --location --request POST 'http://localhost:7777/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url":"https://blog.trello.com/navigate-communication-styles-difficult-times",
    "shortcode":"tes123"
}'
```

### GET /shortcode/:shortcode
```sh
curl --location --request GET 'http://localhost:7777/shortcode/tes123'
```

### GET /shortcode/stats/:shortcode
```sh
curl --location --request GET 'http://localhost:7777/shortcode/tes123'
```

