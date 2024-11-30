# Proyecto de Chat en Tiempo Real Cliente-Servidor

Este proyecto implementa un sistema de chat en tiempo real utilizando una arquitectura cliente-servidor. El chat se maneja completamente a través de la terminal, permitiendo que los usuarios se conecten, envíen y reciban mensajes de texto en tiempo real.

## Tecnologías Utilizadas

- **Go**: Backend del servidor y cliente, utilizando canales (channels) para la distribución de mensajes en tiempo real.
- **SQL Server**: Base de datos para almacenar los usuarios y mensajes.
- **HTTP**: Protocolo utilizado para la comunicación entre el cliente y el servidor, que permite el envío y la recepción de mensajes.

## Funcionalidades

1. **Conexión en Tiempo Real**: Los mensajes se distribuyen a todos los clientes conectados mediante WebSockets, sin necesidad de recargar la terminal.
2. **Interfaz de Línea de Comandos**: El cliente interactúa con el usuario a través de la terminal, permitiendo enviar y recibir mensajes de forma directa.
3. **Almacenamiento de Mensajes**: Los mensajes y usuarios se almacenan en una base de datos SQL Server para persistencia.
4. **Manejo de Múltiples Clientes**: El servidor puede manejar múltiples clientes conectados simultáneamente.

## Arquitectura

### Cliente

El cliente se conecta al servidor mediante solicitudes HTTP a través de la terminal. El cliente puede:

- **Ingresar su apodo**: El apodo se solicita al inicio.
- **Enviar mensajes**: Los mensajes se envían a través de la terminal.
- **Recibir mensajes**: Los mensajes de otros usuarios se reciben en tiempo real en la terminal.

### Servidor

El servidor se encarga de manejar las solicitudes de los clientes, distribuir los mensajes entre ellos, y almacenar los mensajes y usuarios en la base de datos SQL Server.

- **Manejo de mensajes**: Los mensajes enviados por los clientes son procesados y almacenados en la base de datos.
- **Distribución en tiempo real**: Los mensajes se envían a todos los clientes conectados mediante canales (channels) en Go.

## Instalación

### Requisitos previos

- Go (versión mínima recomendada: 1.19)
- SQL Server (o Docker para ejecutar una instancia local de SQL Server)

### Instrucciones

1. **Clona el repositorio**:

    ```bash
    git clone https://github.com/tu-usuario/proyecto-chat.git
    ```

2. **Configura la base de datos**:

    - Asegúrate de tener un servidor SQL Server disponible.
    - Crea una base de datos llamada `CHAT`.
    - Crea las tablas `Usuarios` y `Mensajes` con las siguientes estructuras:

    ```sql
    CREATE TABLE Usuarios (
        ID INT PRIMARY KEY IDENTITY,
        APODO VARCHAR(255) NOT NULL
    );

    CREATE TABLE Mensajes (
        ID INT PRIMARY KEY IDENTITY,
        USUARIO_ID INT,
        CONTENIDO TEXT,
        FOREIGN KEY (USUARIO_ID) REFERENCES Usuarios(ID)
    );
    ```

3. **Configura el archivo de conexión a la base de datos**:

    En el archivo `server.go`, ajusta la cadena de conexión para que apunte a tu base de datos:

    ```go
    connString := "sqlserver://@LESLIE:1433?database=CHAT"
    ```

4. **Instala las dependencias de Go**:

    Dirígete al directorio del backend y ejecuta:

    ```bash
    go mod tidy
    ```

5. **Inicia el servidor**:

    En el directorio raíz del proyecto, ejecuta:

    ```bash
    go run server.go
    ```

6. **Inicia el cliente**:

    En el directorio del cliente, ejecuta:

    ```bash
    go run cliente.go
    ```

    Después, ingresa tu apodo cuando se te solicite para unirte al chat y comienza a enviar y recibir mensajes desde la terminal.

## Licencia

Este proyecto está bajo la Licencia MIT - consulta el archivo [LICENSE](LICENSE) para más detalles.
