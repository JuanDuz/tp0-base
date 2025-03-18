### Ejercicio N°1:

1. Brindar permisos al script

```
chmod +x generar-compose.sh
```

2. Correr el script con el nombre `docker-compose-dev.yaml` y cantidad de clientes requeridos

```
./generar-compose.sh docker-compose-dev.yaml {n}
```

3. Ejecutar

```
make docker-compose-up
```

4. Revisar logs

```
make docker-compose-logs
```
