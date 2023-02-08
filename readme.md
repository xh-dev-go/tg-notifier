### Build

#### Build from go
```shell
go build .
```

#### Build from docker
```shell
docker build -t tg-notifier .
```

#### Build as docker compose
```shell
docker compose up --build
```

### Usage 
When the tg-notifier is up.


##### Event body
```json
{
  "receiver": "{receiver tg id}",
  "encoding": "{none | base64}",
  "message": "{message content in plaintext or base64}"
}
```
Push message (`Event body`) to the source topic, the request will get processed and send to result topic.
The message body is in [markdownV2](https://core.telegram.org/bots/api#markdownv2-style) format.