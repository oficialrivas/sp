# Usa la imagen oficial de Go como base
FROM golang:1.18-alpine

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiar el archivo go.mod y go.sum
COPY go.mod go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar el código fuente de la aplicación
COPY . .

# Compilar la aplicación
RUN go build -o main .

# Especificar el puerto en el que la aplicación escuchará
EXPOSE 8080

# Copiar el archivo .env al contenedor
COPY .env .env

# Comando para ejecutar la aplicación
CMD ["./main"]
