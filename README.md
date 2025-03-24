### Ejercicio N°5:

El cliente se encarga de enviar apuestas a través de una conexión TCP, mientras que el servidor las recibe, almacena y responde con una confirmación (ACK).

Tanto en cliente como en servidor se crea un:
- protocol
- BetFormatter
- BetClient que utiliza el formatter y protocol para mandar y recibir mensajes con el protocolo
- El cliente y servidor tienen su implementacion de un BetClient

El protocolo consta de mandar strings, donde primero se indica la longitud del mensaje con un \n.
A partir de la longitud se puede leer el resto del mensaje.

```
<length>\n<message>
```

El mensaje para transmitir una apuesta se envia y recibe como sus campos ordenados de la siguiente manera, separados por '|'
```
<first_name>|<last_name>|<document>|<birthdate>|<number>|<agency_id>
```

Para revisar ejecución, setear las variables de entorno
```
export NOMBRE='nombre'
export APELLIDO='apellido'
export DOCUMENT0='123456'
export NACIMIENTO='1999-10-10'
export NUMERO='1234'
```

levantar el docker y ver los logs
```
make docker-compose-up
make docker-compose-logs
```