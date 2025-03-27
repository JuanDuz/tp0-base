## 🧵 Ejercicio 8 

### Importante
Correr
```
./generar-compose.sh docker-compose-dev.yaml
```
hace unzip de los datasets, requerido para levantar los dockers, ya que si no no existiran los csv

### ✅ Objetivo

El objetivo de este ejercicio fue adaptar el servidor implementado en Python para que sea capaz de manejar múltiples conexiones de manera concurrente, garantizando la **consistencia de los datos** y evitando **condiciones de carrera**.

---

### ⚙️ Solución implementada

### 📡 Protocolo de Comunicación

Para la interacción entre el cliente (agencia) y el servidor (Lotería Nacional), se definió un protocolo **basado en longitud + mensaje**, con serialización **texto plano delimitado por `|`**. Esto permite estructurar los mensajes de manera sencilla y robusta para su lectura y escritura en sockets TCP.

---


### 📥 Formato General de Envío de Mensajes

Cada mensaje se envía de la siguiente forma:
```
<longitud>\n<mensaje>
```
- `<longitud>`: cantidad de bytes del mensaje (no incluye el `\n`).
- `<mensaje>`: contenido del mensaje, serializado como string.

### 🧾 Tipos de Mensaje

#### ✅ Carga de Apuestas

El cliente envía apuestas en formato texto plano, una por línea:
```
<nombre>|<apellido>|<documento>|<fecha_nacimiento>|<número_apostado>|<id_agencia>
```
ejemplo:
```
Juan|Pérez|12345678|1990-05-23|7574|1
```

Un batch de apuestas puede enviarse concatenando múltiples líneas separadas por `\n`.

---

#### ❓ Consulta de Ganadores

Una vez que el cliente termino de enviar las apuestas pide los ganadores,
si aun todas las agencias no pidieron los ganadores, entonces recibira un error especifico para seguir polleando.
Hasta que reciba los winners.

```
GET_WINNERS|<id_agencia>
```
Ejemplo:
```
GET_WINNERS|3
```

---

### 🧾 Respuestas del Servidor

- Si el sorteo no ha finalizado:
```
ERROR_LOTTERY_HASNT_ENDED
```
- Si el formato del pedido es inválido:
```
ERROR_INVALID_GET_WINNERS
```

- Si el pedido es válido y hay ganadores:

Se devuelven uno o más registros de apuestas ganadoras, uno por línea, usando el mismo formato que en la carga:
```
Juan|Pérez|12345678|1990-05-23|7574|1
```

---

#### ✅ Confirmación del Servidor

El servidor puede responder con: 
Por ejemplo para el caso que el cliente manda un batch de apuestas.
El cliente quedara esperando el ACK del server de que pudo parsearlas correctamente.
```
ACK
```
para confirmar la recepción y procesamiento exitoso de un mensaje.

---

### 🔄 Ejemplo de Intercambio
```
# Cliente 1 envía primer batch (2 apuestas)

Cliente1 → Servidor:
87\nJuan|Pérez|12345678|1990-05-23|7574|1\nMaría|Gómez|87654321|1985-07-10|2311|1

Servidor → Cliente1:
ACK

# Cliente 2 envía primer batch (2 apuestas)

Cliente2 → Servidor:
92\nCarlos|López|99887766|1982-02-17|7574|2\nLucía|Martínez|11223344|1995-09-30|5421|2

Servidor → Cliente2:
ACK

# Cliente 1 consulta ganadores (aún no terminó el sorteo)

Cliente1 → Servidor:
18\nGET_WINNERS|1

Servidor → Cliente1:
24\nERROR_LOTTERY_HASNT_ENDED

# Cliente 2 consulta ganadores (ahora sí termina el sorteo)

Cliente2 → Servidor:
18\nGET_WINNERS|2

Servidor → Cliente2:
94\nJuan|Pérez|12345678|1990-05-23|7574|1\nCarlos|López|99887766|1982-02-17|7574|2

# Cliente 1 vuelve a consultar ganadores

Cliente1 → Servidor:
18\nGET_WINNERS|1

Servidor → Cliente1:
94\nJuan|Pérez|12345678|1990-05-23|7574|1\nCarlos|López|99887766|1982-02-17|7574|2
```

---

#### 🧠 Modelo de concurrencia

Cada conexión entrante desde una agencia se maneja en un **proceso independiente** utilizando el módulo `multiprocessing`. Esto permite que múltiples agencias interactúen simultáneamente con el servidor sin bloquearse entre sí.

#### 📌 Estado compartido

Como los procesos no comparten memoria por defecto, se utilizó `multiprocessing.Manager()` para crear estructuras de datos **compartidas entre procesos**:

- `agencies_ready`: lista compartida para registrar qué agencias solicitaron los resultados.
- `winners`: lista compartida que contiene todas las apuestas ganadoras.
- `lottery_ended`: valor booleano compartido que indica si el sorteo ya fue realizado.

#### 🔒 Secciones críticas

Para evitar **condiciones de carrera** al acceder/modificar estas estructuras compartidas, se utilizó un **`Lock`** también generado con el `Manager`. Este lock protege las siguientes operaciones:

- Verificación y modificación de `agencies_ready`.
- Sorteo de los ganadores (`__draw_lottery`).
- Lectura de resultados si el sorteo ya fue ejecutado.

#### 📁 Acceso al archivo de apuestas

El acceso concurrente al archivo donde se almacenan las apuestas (`bets.csv`) se protege con un **monitor** llamado `BetsFileMonitor`, que encapsula las operaciones de lectura y escritura utilizando un `Lock` propio.

---

### 🧪 Resultado

- Se garantiza que el sorteo solo se realiza una vez, exactamente cuando todas las agencias han solicitado los resultados.
- Los ganadores son almacenados de forma compartida y accesibles por cualquier conexión futura.
- Se evita la corrupción del archivo o duplicación de apuestas.
- Todos los procesos terminan correctamente y el servidor finaliza cuando ya no quedan procesos vivos.

---

### 🔮 Futuras mejoras

- Mantener una **única conexión TCP por cliente** en lugar de abrir una nueva por cada mensaje. Esto evitaría la creación y destrucción de procesos para cada conexión.
- **Cerrar el cliente automáticamente** una vez que recibe la respuesta con los ganadores.
- Al cerrar todos los sockets de los clientes, se podría detectar esta condición y finalizar automáticamente el servidor.

