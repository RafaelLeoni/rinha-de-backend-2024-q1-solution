# Desafio Rinha de Backend, Segunda Edição: 2024/Q1 - Controle de Concorrência
Esta é minha solução para o desafio [Rinha da Backend: 2024/Q1](https://github.com/zanfranceschi/rinha-de-backend-2024-q1) <br />
Solução feita com:
- <img src="https://upload.wikimedia.org/wikipedia/commons/1/1e/Traefik_Logo.svg" alt="logo traefik" width="15" height="auto">  `traefik` como load balancer
- <img src="https://upload.wikimedia.org/wikipedia/commons/2/29/Postgresql_elephant.svg" alt="logo postgres" width="15" height="auto"> `postgres` como banco de dados
- <img src="https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg" alt="logo golang" width="15" height="auto"> `golang` para api com as libs `github.com/gorilla/mux` e `github.com/lib/pq`

# Como Executar
Para executar basta rodar o comando a seguir:
```shell
docker-compose up -d
```

Caso queira construir a imagem em ambiente local, rodar o script a seguir:
```shell
./build.sh
```

Para executar os testes de carga em ambiente local, rodar o script a seguir após ter iniciado os containers:
```shell
./executar-teste-local.sh
```