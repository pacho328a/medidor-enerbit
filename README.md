# Prueba tecnica enerbit Golang

Ing Francisco Anacona  

## Requisitos
- Docker  
- go
- preferiblemente ambiente Linux

## Ejecutar con docker

**Primero** clonar este repositorio.

`git clone https://github.com/pacho328a/medidor-enerbit`

**Segundo** en la raíz del proyecto donde está el archivo main.go

`docker-compose up` o `docker compose up`

## gRPC

importar en postman el archivo **medidorgRPC.proto** está en el path gRPC/medidorRPC/medidorgRPC.prot

El servicio queda escuchando en el puerto 50001 como se ve en la imagen, así que seria localhost:50001

![enter image description here](https://i.ibb.co/LxR6kcf/Captura-de-pantalla-de-2022-12-25-18-58-08.png)

para generar los datos a llenar de manera automática los datos hay que dar clic donde dice **Generate Example Messager** hay que tener en cuenta el número de lines

### Lista de los procedimientos
Aca se encuentran todos los procedimientos con la **lógica de negocio** y el **CRUD**

![enter image description here](https://i.ibb.co/Gfw2CZt/Captura-de-pantalla-de-2022-12-25-18-53-33.png)

## RESTful api

El servicio queda escuchando en el puerto 5000 la url **localhost:5000**
  

![enter image description here](https://i.ibb.co/0DW9qrW/Captura-de-pantalla-de-2022-12-25-19-45-14.png)

  
Dentro de este repositorio hay un archivo **EnerBit.postman_collection.json** que contiene los request CRUD de medidores, en este parte servicio **no tiene la lógica del negocio**, solo la lógica del CRUD

## SWAGGER
Para ingresar al Swagger (openapi) se debe ingresar a la siguiente ruta en el navegador:
**http://localhost:5000/docs/index.html**
![enter image description here](https://i.ibb.co/JyzqJM0/Sw.png)

Acá se encuentran los diferentes puntos de acceso a la API 

#### Ejemplo/Example
Para crear un Medidor damos click en la opción POST:
![enter image description here](https://i.ibb.co/JtR56r7/create-Medidor.png)
**advertencia**
- El valor de lines debe ser mayor a 0 

**Respuesta** 

![enter image description here](https://i.ibb.co/P1ZrbW1/response-Sw.png)

## Adminer Cli postgres
Para ingresar al Adminer Cli  se debe ingresar a la siguiente ruta en el navegador:
**http://localhost:8080**
**Credenciales:**
![enter image description here](https://i.ibb.co/kxqvZZp/cli-db.png)
**contraseña:** example
Acá podemos ver los datos que están almacenados en la base de datos.
![enter image description here](https://i.ibb.co/SvzFFT5/CreateDb.png)

## redis-commander
Para ingresar al redis-commander se debe ingresar a la siguiente ruta en el navegador:
**http://localhost:8081**

![enter image description here](https://i.ibb.co/TBWxz13/redis.png)

Acá podemos ver los identificadores de los medidores que están almacenados en el stream de redis.