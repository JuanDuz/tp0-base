### Ejercicio N°6:

El objetivo del ejercicio 6 fue extender el protocolo de comunicación cliente-servidor implementado previamente, incorporando la posibilidad de que el cliente envíe múltiples apuestas en un solo mensaje (batch), en lugar de enviar una apuesta por vez. Además, se debía:

Limitar la cantidad máxima de apuestas por batch, configurable mediante config.yaml bajo la clave batch: maxAmount.

Garantizar que el tamaño del mensaje nunca supere los 8kB, ajustando el valor por defecto del maxAmount en consecuencia.

Asegurar que si al menos una apuesta del batch es inválida, el servidor rechaza todo el batch y responde con un error.

Separar claramente las responsabilidades entre la lógica del protocolo de comunicación y la lógica de dominio de las apuestas.

### Protocolo
Anteriormente, el protocolo permitía que el cliente enviara una apuesta como un string con longitud prefijada. Para admitir batches, mantuvimos el mismo protocolo base (longitud\nmensaje) y simplemente concatenamos múltiples apuestas separadas por saltos de línea (\n) en un solo string.
Ej:
```
78\n
A|B|00000000|2000-01-01|0|1
A|B|00000001|2000-01-01|1|1
```

### Manejo de errores y corner cases
Archivo vacío o sin apuestas válidas: Se responde con ERROR_EMPTY_BATCH.

Apuesta malformada: Todo el batch se descarta. ERROR_INVALID_BATCH

Batch muy grande: El cliente lo limita antes de enviarlo.

Errores en conexión: Ambos lados manejan errores de socket con mensajes claros y logs correspondientes.