# Hire.me

O projeto é um serviço que permite encurtar URLs, tornando-as mais fáceis de serem compartilhadas e lembradas.

## Utilização

Para utilizar o projeto, é necessário ter o Go instalado em sua máquina. Em seguida, execute o seguinte comando:

```bash
go run main.go
```

## Utilização

### /create

O endpoint /create é utilizado para criar uma nova URL encurtada a partir de uma URL original.

Parâmetros

- url: a URL original a ser encurtada.

- CUSTOM_ALIAS (opcional): um alias personalizado para a URL encurtada.

### Exemplo

```bash
PUT http://localhost:8080/create?url=https://youtube.com&CUSTOM_ALIAS=yt

```

### Resposta

```json
{
  "alias": "yt",
  "url": "https://youtube.com",
  "statistics": {
    "time_taken": "359ms"
  }
}
```

### /{alias}

O endpoint /{alias} é utilizado para acessar a URL original associada a um alias.

Parâmetros

- alias: a URL encurtada.

### Exemplo

```bash
GET http://localhost:8080/yt

```

### Resposta

Caso a URL esteja registrada no banco, automaticamente será redirecionado ao site destino, caso não uma mensagem de erro será retornada

```json
{
  "alias": "yt",
  "err_code": "002",
  "description": "SHORTENED URL NOT FOUND IN DATABASE"
}
```
