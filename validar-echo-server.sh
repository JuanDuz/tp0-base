#!/bin/bash

# Configuración
SERVER_IP="server"
SERVER_PORT=12345
TEST_MESSAGE="HelloServer"
TIMEOUT=2

# Enviar mensaje con printf y recibir respuesta
RESPONSE=$(docker run --rm --network=tp0_testing_net busybox sh -c "printf '$TEST_MESSAGE' | nc -w $TIMEOUT $SERVER_IP $SERVER_PORT")

# Verificar si la respuesta es igual al mensaje enviado
if [ "$RESPONSE" == "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
    exit 0
else
    echo "action: test_echo_server | result: fail"
    exit 1
fi
