## Infra

```
# rabbit
docker run -it --rm  --name some-rabbit -p 5672:5672 -p 8080:15672 rabbitmq:3-management

# es
docker run -p 9200:9200 --rm --name elasticsearch -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.5.0

# kibana
docker run --rm  --link elasticsearch:elasticsearch -p 5601:5601 docker.elastic.co/kibana/kibana:7.5.0
```

## Usage

```
./mq-cli -h
Usage of ./mq-cli:
  -channel string
        Channel name (default "default")
  -connection-string string
        The connection string used to interact with the queue. (default "amqp://guest:guest@localhost:5672/")
  -mode string
        Mode: [publish, consume] (default "publish")
```

```
# produce
~/Downloads/subfinder-darwin-amd64 -d uber.com | ./mq-cli

# consume
./mq-cli --mode consume | ~/Downloads/httprobe-max | while read line; do echo $line |  http POST :9200/pages/_doc;  done
```