Build your image:
`make build-dev`
OR 
`make docker-build`

Set up your kafka  with the docker compose:
`make docker-kafka-start`

Inject your metrics ( by default it inject to `metrics` topic ): 
`echo "mymetric 42 000000000" | make docker-inject`

Run your image:
`make run`
OR 
`make docker-start` 

Stop kafka:
`make docker-kafka-stop`

Stop your image by sighup or by `make docker-stop`
