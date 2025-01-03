# Base para el proyecto de Bases de Datos 2 usando Cassandra

Esta base hace uso de Docker CLI, en Windows eso se traduce a instalar Docker Desktop
(el cual necesita tener la vitualización habilitada, ya sea Hyper-V o AMD-V)
o instalar el paquete `docker` usando WSL (aqui un [articulo][articulo-wsl] para facilitar el proceso).

En Linux es tan fácil como instalar el paquete `docker` y empezar a usarlo,
los siguientes comandos se pueden utilizar para mejorar el uso de la herramienta en la terminal:
```bash
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
```
-----

## Comandos implementados en el makefile

- `make` o `make run`: Se utiliza para ejecutar el código escrito en _Go_ el cual se comunica con el contenedor docker de Cassandra.

- `make get-cassandra`: Para obtener la imagen del contenedor de Cassandra si aun no se ha descargado.

- `make set-cassandra`: Para crear la red que vamos a utilizar para conectarnos con Cassandra desde fuera del contenedor.

- `make run-cassandra`: Para crear el nodo principal de Cassandra en un contenedor llamado _cassandra1_.

- `make run-cassandra2`: Para crear un nodo secundario de Cassandra en un contenedor llamado _cassandra2_.
Otros nodos se pueden crear con una sintaxis similar al comando descrito en el makefile si es necesario.

- `make run-cqlsh`: Para ejecutar el shell desde donde podemos ejecutar comandos de cql en la base de datos directamente.

- `make inspect-ip`: Para visualizar la dirección ip del contenedor principal, esto es necesario para conectarse desde fuera del contenedor.

- `make cleanup`: Para deshacer lo que los comandos anteriores han hecho y cerrar los contenedores docker.
Si se han creado nodos adicionales a los 2 disponibles en el makefile se deben detener y eliminar manualmente.


#### El orden de ejecución de los comandos es el siguiente:

1. `make get-cassandra` en caso de que no se haya descargado la imagen de Cassandra.
2. `make set-cassandra` para crear la red.
3. `make run-cassandra` para iniciar el nodo principal de Cassandra en un contenedor.
4. `make inspect-ip` para obtener la dirección ip del contenedor que puede no ser la misma para cada uno.
5. Una vez obtenida la dirección ip podemos insertar este valor en el archivo **.env** con la llave `CASSANDRA_IPADDRESS`.
6. Por ultimo podemos ejecutar `make` o `make run` para ejecutar el código escrito en _Go_.
Si los pasos se han seguido correctamente el programa debería ejecutarse sin problema.
Es posible que las primeras al correr este comando de algún error inmediatamente después de crear el contenedor de Cassandra,
para evitar esto es mejor esperar 5 o 10 minutos antes de ejecutar este comando.


[articulo-wsl]: https://dev.to/julianlasso/how-to-install-docker-cli-on-windows-without-docker-desktop-and-not-die-trying-4033
