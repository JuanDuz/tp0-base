### Ejercicio N°7:

En este ejercicio se aprovecho y se hizo un refactor en cliente y servidor para apuntar a una arcitectura CLEAN

### Lado del Servidor
#### Arquitectura
Utilizamos una arquitectura tipo Clean Architecture que separa responsabilidades en capas:

- MessageHandler: router de mensajes que delega según el tipo.

- BetController: parsea, valida y responde con errores del protocolo.

- BetService: maneja la lógica de dominio (almacenamiento, sorteo, cálculo de ganadores).

#### Lógica del sorteo
El servidor guarda los IDs de las agencias que notificaron finalización.

Una vez que se notifican las N agencias activas (obtenidas dinámicamente desde el entorno en tiempo de ejecución), se realiza el sorteo:

1. Se invoca load_bets() para cargar todas las apuestas.

2. Se filtran las apuestas ganadoras con has_won(...).

3. Se almacenan los ganadores en un diccionario por agencia.

#### Consulta de ganadores
Antes del sorteo, el servidor responde con ERROR_LOTTERY_HASNT_ENDED.

Luego del sorteo, cada consulta se responde únicamente con los ganadores de la agencia correspondiente, respetando el requerimiento de no realizar broadcast.


### Lado del Cliente
#### Arquitectura

- Application: punto de entrada, realiza la orquestación e inyección de dependencias.
- SendBetsUseCase: se encarga del envío en batches respetando el límite de 8KB y el maxAmount pasado por config.
- PollWinnersUseCase: consulta periódicamente los ganadores hasta obtener respuesta exitosa. Respetando el loopTime pasado por config
- NetworkClient: maneja las conexiones TCP y el protocolo de comunicación.

#### Flujo de ejecución
1. El cliente carga sus apuestas desde el archivo CSV y las envía por batches.
2. Al alcanzar el final del archivo, notifica automáticamente al servidor.
3. Luego entra en modo "polling" hasta que el sorteo se haya realizado.
4. Una vez recibe la lista de ganadores, loguea la cantidad de winners que obtuvo

### Protocolos y Validaciones
El cliente y servidor usan un protocolo de mensajes con ```length\nmessage_body```

Se validan errores como:

- Consulta antes del sorteo.
- Agencia inválida.
- Consulta malformada.