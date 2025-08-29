# Desafio Clean Architecture

## Descrição

Olá devs!

Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o `usecase de listagem das orders`.  
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL

Não esqueça de criar as migrações necessárias e o arquivo `api.http` com a request para `criar` e `listar as orders`.  
Para a criação do banco de dados, utilize o `Docker (Dockerfile / docker-compose.yaml)`, com isso ao rodar o comando `docker compose up` tudo deverá subir, preparando o banco de dados.  
Inclua um `README.md` com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

## Como baixar a aplicação

1. Rode no terminal o seguinte comando  
    ```bash
    $ git clone git@github.com:Berchon/Clean-Architecture.git
    ```

## Como executar a aplicação

1. Subir o MySQL e o RabbitMQ  
    ```bash
    $ docker-compose up -d
    ```
    1.1. Caso precise remover todas as instâncias ativas no docker  
    ```bash
    $ docker rm -f $(docker ps -a -q)
    ```

2. Levantar os servidores  
    ```bash
    $ cd cmd/ordersystem/
    $ go run main.go wire_gen.go
    ```
- O servidor REST roda na porta `8000`
- O servidor GraphQL roda na porta `8080`
- O servidor GRPC roda na porta `50051`

## Como efetuar as chamadas aos servidores

1. As chamadas `REST` podem ser feitas através do arquivo `./api/api.http`.
2. As chamadas `GRPC` podem ser feitas utilizando o `evans`
    ```protobuf
    $ evans -r repl

    127.0.0.1:50051> package pb

    pb@127.0.0.1:50051> service OrderService
    ```

    2.1. Para criar uma nova `order`:
    ```protobuf
    pb.OrderService@127.0.0.1:50051> call CreateOrder
    id (TYPE_STRING) => id-2
    price (TYPE_FLOAT) => 10
    tax (TYPE_FLOAT) => 1
    {
      "finalPrice": 11,
      "id": "id-2",
      "price": 10,
      "tax": 1
    }
    ```

    2.2. Para listar todas as `orders`:
    ```protobuf
    pb.OrderService@127.0.0.1:50051> call ListOrders
    {
      "orders": [
        {
          "finalPrice": 101,
          "id": "id-1",
          "price": 100.5,
          "tax": 0.5
        },
        {
          "finalPrice": 11,
          "id": "id-2",
          "price": 10,
          "tax": 1
        }
      ]
    }
    ```

3. As chamadas `GraphQL` podem ser feitas utilizando o navegador
    3.1. Digite na barra de endereço do navegador
    ```bash
    localhost:8080
    ```

    3.2. Para criar uma `Order` utilize a `mutation createOrder` abaixo
    ```properties
    mutation createOrder {
      createOrder(input:{id: "id-3", Price: 12.2, Tax: 2.0}) {
        id
        Price
        Tax
        FinalPrice
      }
    }
    ```

    3.3. Para visualizar todas as `Orders` criadas utilize a `query listOrders` abaixo
    ```properties
    query listOrders {
      orders {
        id
        Price
        Tax
        FinalPrice
      }
    }
    ```

## Como configurar o RabbitMQ
1. Digite na barra de endereço do navegador
    ```bash
    http://localhost:15672/
    ```

2. Entre com o `usuário` e `senha`
    ```bash
    usuário: guest
    senha: guest
    ```

3. Adicione uma fila chamada `orders` na aba `Queues and Streams`, ou seja, faça um `Bind` da `exchanges` à fila `orders`

4. Na aba `Exchanges` vincule o `amq.direct` com a fila `orders`

## Como acessar o Banco de Dados MySQL

Executar o comando `bash` dentro do container de serviço `MySQL`.
```bash
$ docker-compose exec mysql bash
```
Entrar no CLI do MySQL (senha root)
```bash
bash-4.2# mysql -uroot -p
```
Listar todas as base de dados
```bash
mysql> show databases;
```
Vai existar uma base de dados chamada `orders`. Para entrar nessa base de dados
```bash
mysql> use orders
```
Para visualizar as tabelas da base de dados `orders`
```bash
mysql> show tables;
```
Caso não exista nenhuma tabela, pode-se criar a tabela usada na aplicação com o comando abaixo. **OBS.: Ao rodar o main.go a aplicação se encaregará de criar a tabela**
```sql
mysql>  CREATE TABLE IF NOT EXISTS orders (
          id varchar(255) NOT NULL, 
          price float NOT NULL, 
          tax float NOT NULL, 
          final_price float NOT NULL, 
          PRIMARY KEY (id)
        );
```
## Como gerar o wire inject
```bash
  $ cd cmd/ordersystem
  $ wire
```

## Como gerar os arquivos do Protol Buffers para o GRPC
```bash
  $ protoc --go_out=. --go-grpc_out=. ./internal/infra/grpc/protofiles/order.proto
```
## Como gerar os arquivos do GraphQL a partir de um `Schema`
```bash
  $ go run github.com/99designs/gqlgen generate
```

