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

O endpoint /r/{alias} é utilizado para acessar a URL original associada a um alias.

Parâmetros

- alias: a URL encurtada.

### Exemplo

```bash
GET http://localhost:8080/r/yt

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

### /mais_acessadas

O endpoint /mais_acessadas é utilizado para vizualizar as 10 Urls mais acessadas

### Exemplo

```bash
GET http://localhost:8080/mais_acessadas

```

### Resposta

```json
[
  { "alias": "Facebook", "url": "https://www.facebook.com", "acessos": 1000 },
  { "alias": "Google", "url": "https://www.google.com", "acessos": 750 },
  { "alias": "YouTube", "url": "https://www.youtube.com", "acessos": 500 },
  { "alias": "Twitter", "url": "https://www.twitter.com", "acessos": 400 },
  { "alias": "Instagram", "url": "https://www.instagram.com", "acessos": 350 },
  { "alias": "LinkedIn", "url": "https://www.linkedin.com", "acessos": 250 },
  { "alias": "Amazon", "url": "https://www.amazon.com", "acessos": 200 },
  { "alias": "Wikipedia", "url": "https://www.wikipedia.org", "acessos": 150 },
  { "alias": "Reddit", "url": "https://www.reddit.com", "acessos": 100 },
  { "alias": "Netflix", "url": "https://www.netflix.com", "acessos": 50 }
]
```
