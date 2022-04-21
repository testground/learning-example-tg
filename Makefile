# Starts docker containers which are required for certain test plans (e.g. RabbitMq tests)
start_docker_deps:
	start_docker_rabbit

# Starts the RabbitMq container. The image used already has an admin panel integrated, and 
# is accessible on port 15672 (e.g. localhost:15672 from the host machine)
start_docker_rabbit:
	docker run -p 5672:5672 -p 15672:15672 --name rabbitmq --network testground-control rabbitmq:3-management