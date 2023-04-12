# Imagem base com Go 1.16 instalado
FROM golang:1.16-alpine

# Define o diretório de trabalho
WORKDIR /app

# Copia o código para o diretório de trabalho
COPY . .

# Instala a biblioteca godotenv
RUN go get github.com/joho/godotenv

# Compila o projeto
RUN go build -o main .

# Define a porta em que a aplicação escuta
ENV PORT=8080

# Expõe a porta
EXPOSE $PORT

# Define o comando para iniciar a aplicação
CMD ["./main"]
