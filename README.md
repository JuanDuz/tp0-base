## 🧵 Ejercicio 8 

### ✅ Objetivo

El objetivo de este ejercicio fue adaptar el servidor implementado en Python para que sea capaz de manejar múltiples conexiones de manera concurrente, garantizando la **consistencia de los datos** y evitando **condiciones de carrera**.

---

### ⚙️ Solución implementada

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

