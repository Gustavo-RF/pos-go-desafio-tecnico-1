# Rate limiter
O objetivo desse desafio é desenvolver um rate limiter, onde, dado configurações iniciais, bloqueará requisições por IP ou por um token no header de um determinado client.

## Funcionamento
O rate limiter injeta o rate limiter como middleware no servidor http. Como padrão, utilizo o redis para arquivar as informações de limiter de um determinado client, sendo o IP ou o token API_TOKEN no header.

Para cada requisição, ao executar o tratamento, a aplicação verifica a middleware que valida se aquele client está bloqueado. Se estiver, é retornado um erro 429 no formato:
```
{
  "message": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

Caso o client não esteja bloqueado, a execução da chamada da requisição continua e é adicionado um novo incremento no redis daquele client.

## Configurando
Na raiz do projeto existe um arquivo .env:
```
PORT=8080                    // Porta padrão da aplicação
REDIS_ADDRESS=localhost      // endereço padrão do redis
REDIS_PORT=6379              // porta padrão do redis
REDIS_PASSWORD=              // Senha para acesso ao redis

RATE_LIMITER=IP              // Identifica o tipo de bloqueio do rate limiter. Pode ser IP, padrão, onde é validado o IP do client ou pode ser TOKEN, onde valido o header na request
REQUESTS_PER_SECOND=5        // Configura qual o limite aceitável de requisições por segundo
BLOCKED_TIME_IN_SECONDS=60   // Configura quanto tempo o client ficará bloqueado em segundos
```

## Executando o projeto
Para executar o projeto, suba o docker compose:
```
docker compose up -d --build
```

Ao subir, temos acesso ao servidor em localhost:8080 e o redis commander em localhost:8081

## Testes
Para validar as configurações, utilizaremos o [vegeta](https://github.com/tsenart/vegeta). Ele é uma ferramenta de teste de carga, onde podemos mandar várias requisições simultâneas e testar a resposta do servidor.
Exemplo de teste com o vegeta:
```
echo "GET http://localhost:8080" | vegeta attack -rate=100 -duration=10s | vegeta report
```
ou, caso seja por token:
```
echo "
GET http://localhost:8080
Api-Key: teste2
" | vegeta attack -rate=10 -duration=3s | vegeta report
```

Podemos ver que configuramos 100 requisições por segundo durante 10 segundos.

Ou pode ser executado um curl para o endereço:
```
curl --request GET --url http://localhost:8080/
```

Caso seja configuração por TOKEN, ele deve ser informado na requisição
```
curl --request GET --url http://localhost:8080/ --header 'Api-Key: teste'
```

## Resultados
O vegeta nos mostra as informações das tentativas separadas por http status:
![image](https://github.com/user-attachments/assets/2762b715-c33c-441d-8e77-570be4601d6b)


