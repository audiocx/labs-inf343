docker-ONU:
	sudo docker compose -f ONU/docker-compose.yml run --rm onu

docker-OMS:
	sudo docker compose -f OMS/docker-compose.yml up

HOST = $(shell hostname)
docker-datanode:
ifeq ($(HOST),dist027)
	sudo docker compose -f DN1/docker-compose.yml up
endif
ifeq ($(HOST),dist028)
	sudo docker compose -f DN2/docker-compose.yml up
else
	echo "Hecho"
endif

docker-continentes:
ifeq ($(HOST),dist025)
	sudo docker compose -f AS/docker-compose.yml up 
endif
ifeq ($(HOST),dist026)
	sudo docker compose -f AU/docker-compose.yml up
endif
ifeq ($(HOST),dist027)
	sudo docker compose -f EU/docker-compose.yml up
endif
ifeq ($(HOST),dist028)
	sudo docker compose -f LA/docker-compose.yml up
else
	echo "Hecho"
endif