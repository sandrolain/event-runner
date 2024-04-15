# Event Runner

Event Runner is a simple, distributed event processing system. It consists of several components:

* An event input, which reads events from a source such as HTTP, gRPC, NATS
* An event runner, which executes user-provided code (JavaScript, WASM, ...) against each event
* An event output, which writes the results of the event runner to a destination such as HTTP, gRPC, NATS

The system is designed to be scalable and fault-tolerant, with support for multiple event inputs, runners, and outputs. The system also includes support for buffering events in case of failure, and for handling events in parallel to improve performance.

The system is written in Go, and is intended to be run as a Docker container. It is designed to be easy to use, with a simple configuration file that specifies the inputs, runners, and outputs to use.

Future improvements
-------------------

* [ ] Support for more input and output sources:
  * [x] HTTP
  * [ ] gRPC
  * [x] NATS
    * [ ] NATS pub/sub
    * [ ] NATS streams
  * [ ] Redis
    * [ ] Redis streams
  * [ ] Kafka
  * [ ] Other messaging systems (e.g. RabbitMQ)
* [ ] Support for more runners:
  * [x] JavaScript (ES5 with [Goja](https://github.com/dop251/goja))
  * [ ] WebAssembly (WASM with [Wazero](https://github.com/tetratelabs/wazero))
  * [ ] Go (Go scripting with [Risor](https://risor.io/))
  * [ ] Other languages (e.g. PHP?)
* [ ] Other improvements:
  * [ ] Test coverage (unit tests, integration tests, etc.)
  * [ ] Versioning management (semver, tags)
  * [ ] Image build (docker, podman, etc.)
  * [ ] Kubernetes resource definition (deployment, service, etc.)
  * [ ] TLS and mTLS
  * [ ] Observability
  * [ ] Integrated cache options (e.g. NATS, Redis, in-memory)
