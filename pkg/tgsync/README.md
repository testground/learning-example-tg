# Testground sync test plan

This test plan utilizes Testground's sync logic to transmit messages from producers to consumers.
It does not depend on any additional dependencies, other than the ones already started by Testground.

## Getting started

```
Navigate to the project root directory

# Import the plan
testground plan import --from . --name tg-learning

# Run a test case
testground run single --plan tg-learning --testcase tgsync-1to1 --builder docker:generic --runner local:docker --instances 4
```