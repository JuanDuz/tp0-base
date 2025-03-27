#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"

unzip -n .data/dataset.zip -d .data

python3 ./generar_compose.py $1 $2