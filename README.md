### Ejercicio N°3:

Se levanta un container momentaneo conectado a la misma network que server, con imagen de busybox para correr nc y luego --rm

1. Dar permisos al script

```
chmod +x validar-echo-server.sh
```

2. Correr build and up

```
make docker-compose-up
```

3. Ver logs

```
make docker-compose-logs
```

4. En nueva terminal correr y visualizar resultado

```
./validar-echo-server.sh
```
