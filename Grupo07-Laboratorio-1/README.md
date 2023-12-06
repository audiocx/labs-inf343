# Grupo07-Laboratorio-1

## Integrantes
1. Constanza Beatriz Alvarado Valenzuela - 201973521-7
1. Exequiel Orlando Perez Lopez - 201873555-8
1. Claudio Isaac Inal Quintrel  - 19795016-1

## Instrucciones de ejecucion:

1. En maquina 26, revisar que rabbitmq no este siendo usado: sudo lsof -i:15672 y luego sudo kill -9 <PID> si es que esta siendo usado.
1. Ejecutar make docker-rabbit en maquina 26
1. Ejecutar make docker-regional en todas las maquinas, 25, 26, 27 y 28
1. Ejecutar make docker-central en maquina 25