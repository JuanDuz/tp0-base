### Ejercicio N°4:

Se modificó **el cliente y el servidor** para que respondan de manera adecuada a la señal `SIGTERM` y finalicen sus procesos de forma **graceful**. Esto implica:

- Cierre correcto de **file descriptors**
- Loguear el cierre de los recursos afectados.
- Terminar el proceso principal **solo después** de cerrar los recursos correctamente.

#### Cliente
- Se utilizó el paquete `os/signal` y `context` de Go para capturar señales `SIGTERM` y `SIGINT`.
- Al recibir la señal, se ejecuta el close del cliente:
    - Se termina la iteración del loop actual cerrando la conexión.
    - Se sale del loop.
    - Se loguea el shutdown

#### Servidor

- Se utilizó el módulo `signal` de Python para capturar `SIGTERM` y `SIGINT`.
- Al recibir la señal:
    - Se cierra el socket del servidor.
    - Se marca el flag was_killed como true para salir del loop de aceptar clientes.
    - Se loguea el cierre del socket y el shutdown con los mensajes correspondientes.
    - El proceso termina con `sys.exit(0)` para indicar una salida exitosa.