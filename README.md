# Testground learning example

This project is intended as a practical example of a "real" testground project. The test plans contained within are fairly straightforward and easy to understand, and demonstrate how to solve issues you may encounter when you are writing test plans for your projects.

## Getting started

Note: this assumes you have already installed [Testground]
```
testground daemon

# => open a different console in your client
git clone https://github.com/testground/learning-example-tg
cd learning-example-tg

# import all test plans from the repository
$ testground plan import --from . --name tg-learning

# run two instances of the `simple` test case
# building with generic:docker, running with local:docker
$ testground run single --plan tg-learning --testcase tgsync-1to1 --builder docker:generic --runner local:docker --instances 2
```

## Featured test plans

- Tests which rely on learning-example project
- Tests which require additional docker containers to be run (e.g. RabbitMq)

## Tech

- [Docker] - Basic Testground dependency
- [Learning-Example] - The project used as an example of a "real" application with its own logic

## License

Dual-licensed: [MIT](./LICENSE-MIT), [Apache Software License v2](./LICENSE-APACHE), by way of the
[Permissive License Stack](https://protocol.ai/blog/announcing-the-permissive-license-stack/).



[//]: # (Reference links)

   [Learning-Example]: <https://github.com/testground/learning-example>
   [Docker]: <https://www.docker.com/>
   [Testground]: <https://github.com/testground/testground>