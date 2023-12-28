

# Highly available, fault tolerant, distributed file storage

Pipcas is highly available distributed file storage with eventual consistency and high throughput

## Features

- storing files on disk across multiple shards instead of one single storage unit
- Hash based sharding
- Shards are essentially deployable storage units 
- Optimized read write throughput because of high concurrency
- Pipcas pipelines file and object store process along with pipcas-lb and pipcas-client
- Shards can be deployed and run in isolation, which ensures that downtime of one shard wont affect other shards of the node 

## Deployment Support

Pipcas supports flexible deployments(e.q shards can be deployed in one node and also across various nodes in the cluster). Depending upon the usecases business can choose their deployment strategy as it does not have any bounded rule for deployment. It also supports deployments across spectrum of infra such as on premise and cloud and also hybrid approach.


# Quickstart

## Deploy the binaries on your local machine directly onto your host os.

make sure Go 1.18 is installed and it's path variable is configured

clone the repository
```
    git clone https://github.com/raja-dettex/pipcas
```
```cd pipcas```

```go mod tidy```

```make build```

```make run```

## Deploy onto a docker cotainer

after cloning the repository as mentioined previously,

Build a Docker image from the Docker file

```docker build -t <your image name> .```

Start the docker container from the built image

```docker run -d -p <port to access contiainer process>:<port according to listen address> -e  LISTEN_ADDR=<port> <your image name>000```

## latest release: 
    pipcas:1.0

### N.B : About upcoming release
    docker images will be available directly on docker hub on the next release