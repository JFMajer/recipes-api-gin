start_mongo:
	docker run -d --name mongodb -v /home/kuba/data:/data/db -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3

start_redis:
	docker run -d -v ${CURDIR}/redis-conf:/usr/local/etc/redis --name redis -p 6379:6379 redis:6.2.3 redis-server /usr/local/etc/redis/redis.conf

start_redis_insight:
	docker run -d --name redisinsight --link redis -p 8001:8001 redislabs/redisinsight:latest

deletecontainers:
	docker stop $$(docker ps -a -q) && docker container rm $$(docker ps -aq) || true