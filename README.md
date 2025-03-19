### Ejercicio N°2:

Se hace uso de los volumenes en docker-compose para evitar tener que buildear tras cambios en config.
Se elimina del script py generador del docker-compose yaml la linea que fozaba loggin_level debug

Corriendo

```
docker-compose -f docker-compose-dev.yaml up
```

se puede probar realizar cambios a la config y volver a correr

```
docker-compose -f docker-compose-dev.yaml up
```

para chequear si los cambios impactaron sin hacer docker-compose -f docker-compose-dev.yaml build
