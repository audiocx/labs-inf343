docker-rabbit:
	sudo docker compose -f docker-compose.yml up

docker-central:
	sudo docker compose -f Central/docker-compose.yml up

HOST = $(shell hostname)
docker-regional:
ifeq ($(HOST),dist025)
	sudo docker compose -f America/docker-compose.yml up 
endif
ifeq ($(HOST),dist026)
	sudo docker compose -f Asia/docker-compose.yml up
endif
ifeq ($(HOST),dist027)
	sudo docker compose -f Europa/docker-compose.yml up
endif
ifeq ($(HOST),dist028)
	sudo docker compose -f Oceania/docker-compose.yml up
else
	echo "Hecho"
endif