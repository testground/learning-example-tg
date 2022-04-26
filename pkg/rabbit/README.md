# RabbitMq test plan

The test cases within this test plan use a RabbitMq broker to transmit messages from producers to consumers.
All test cases will fail (ie. the testground runner will fail to connect to the RabbitMq broker) if you haven't 
started a RabbitMq docker image in the `testground-control` network.

## Getting started

```
Navigate to the project root directory

# Start RabbitMq:
make start_docker_rabbit

# Import the plan
testground plan import --from . --name tg-learning

# Run a test case
testground run single --plan tg-learning --testcase rabbit-4to1 --builder docker:generic --runner local:docker --instances 4
```