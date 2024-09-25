# Encoder

## SOBRE

Esse projeto é um encoder de vídeos que converte arquivos .MP4 em fragmentos no formato .MPEG-DASH, permitindo sua utilização em sistemas de streaming. A aplicação monitora uma fila do RabbitMQ chamada "videos.new", da qual recebe o caminho do vídeo armazenado em um bucket do Google Cloud Platform. O processo inclui o download do vídeo para a máquina local, a conversão e a fragmentação para o formato MPEG-DASH utilizando a biblioteca Bento4. Em seguida, todos os arquivos gerados são enviados de volta para o bucket, e os arquivos locais são deletados. O resultado da operação — sucesso ou falha — é enviado para outra fila do RabbitMQ chamada "videos.converted", que será monitorada por um serviço distinto. Em caso de erro, a mensagem original é redirecionada para uma fila com uma exchange chamada "dlx" do tipo fanout, para monitoramento dos erros. A aplicação é escalável, permitindo ajustar a quantidade de workers através de uma variável de ambiente. O Go facilita essa implementação com goroutines, e também lida com erros de race condition.

Na pasta domain, temos as entidades DAO que refletem os objetos do banco, onde aplicamos as regras de negócio e garantimos a integridade dos dados. A entidade Video representa os inputs da fila "videos.new", enquanto a entidade Job representa a tarefa a ser executada pela aplicação.

A pasta application serve como uma camada acima do domain no modelo de Arquitetura Hexagonal. A pasta repositories abstrai a lógica de acesso a dados, utilizando a biblioteca Go chamada GORM para realizar as queries no banco de dados. Se um dia for necessário trocar essa biblioteca por SQL puro, basta modificar essa classe, sem afetar o restante do sistema.

Ainda na pasta de application, temos a camada service, em que é responsável por organizar a lógica da aplicação e orquestrar as operações entre entidades e repositórios. Essa camada atua como um intermediário, garantindo que cada componente se concentre em sua responsabilidade específica.

O job_manager é a primeira parte dessa camada, gerenciando todos os jobs que a aplicação precisa processar. Ele centraliza informações essenciais, como os dados dos jobs, a conexão com o banco de dados e a fila RabbitMQ. O job_manager também determina quantos workers serão utilizados para processar os jobs em paralelo, permitindo que a aplicação opere de forma eficiente. Quando um job é concluído, ele notifica o RabbitMQ sobre o sucesso ou erro do processamento.

Em seguida, o job_worker entra em ação. Essa função é responsável por processar as mensagens recebidas da fila. Ela transforma as mensagens em formato JSON, verifica se estão corretas e completas, e valida os dados. Após a validação, o job_worker cria um objeto Job, insere-o no banco de dados e chama a função do job_service para iniciar o processo necessário para o vídeo.

O job_service organiza o fluxo das operações que envolvem o video_service e o upload_manager. Ele gerencia diversas tarefas, como o download do vídeo, sua fragmentação, o encode e a conexão com o bucket de armazenamento. Essa organização garante que todas as etapas sejam executadas de maneira sequencial e eficiente.

Por fim, o video_service foca nas operações relacionadas ao objeto vídeo. Ele é responsável pelo download do vídeo do bucket, pela fragmentação em partes e pela execução do encode usando a biblioteca Bento4. Após o processamento, o video_service também se encarrega de remover os arquivos locais temporários, mantendo o ambiente limpo.

O upload_manager, por sua vez, finaliza o fluxo enviando todos os arquivos convertidos e fracionados de volta para o bucket de armazenamento. Ele otimiza o processo de upload ao usar goroutines, permitindo que múltiplos arquivos sejam enviados simultaneamente, o que melhora a eficiência e a performance da aplicação.


## Estrutura do repositório

```bash

encoder
├── application
│   ├── repositories
│   │   ├── job_repository.go
│   │   ├── job_repository_test.go
│   │   ├── video_repository.go
│   │   └── video_repository_test.go
│   └── services
│       ├── job_manager.go
│       ├── job_service.go
│       ├── job_worker.go
│       ├── upload_manager.go
│       ├── upload_manager_test.go
│       ├── video_service.go
│       └── video_service_test.go
├── bucket-credentials.json
├── docker-compose.yaml
├── Dockerfile
├── domain
│   ├── job.go
│   ├── job_test.go
│   ├── video.go
│   └── video_test.go
├── framework
│   ├── cmd
│   │   └── server
│   │       └── server.go
│   ├── database
│   │   └── db.go
│   ├── queue
│   │   └── queue.go
│   └── utils
│       ├── utils.go
│       └── utils_test.go
├── go.mod
├── go.sum
└── README.md

```


## Como rodar esse projeto

```bash
# clone este repositorio
git clone https://github.com/Igorcand/encoder

# Entre na pasta
cd encoder

# Rode os serviços
docker-compose up --build

```

## RabbitMQ
Acesse o dashboard do RabbitMQ no endpoint "http://127.0.0.1:15672" e faça login com as credenciais: Username = rabbitmq | Password = rabbitmq

### Fila videos.new
Esta fila será consumida pelo encoder para saber quais videos precisa converter, o formato das mensagens tem que ser na estrutura:

```

{   
    "resource_id": "260310fc-53f1-4d92-8b81-81b25613537f.VIDEO", 
    "file_path": "videos/260310fc-53f1-4d92-8b81-81b25613537f/video.mp4"
}

```

### Fila videos.converted
Esta fila será consumida pelo administrador do catálogo de videos, onde o encoder envia as mensagens de conversão, o formato das mensagens tem a estrutura:

```

{
    "error": "",
    "video": {
        "resource_id": "24846fe1-6218-46bb-96a8-d6d4534e0885.VIDEO",
        "encoded_video_folder": "/path/to/encoded/video"
    },
    "status": "COMPLETED"
}

```

### Fila videos.rejected
Esta fila tem o bind com a exchange "dlx" e é do tipo fanout, onde todo erro de gerado na fila videos.converted será enviado a mensagem original para ser analisada.

```

{   
    "resource_id": "260310fc-53f1-4d92-8b81-81b25613537f.VIDEO", 
    "file_path": "videos/260310fc-53f1-4d92-8b81-81b25613537f/video.mp4"
}

```


# Tecnologias Usadas

## Back end
- GO
- Gin
- Docker

## Database
- SQLite
  
## Infra
- RabbitMQ
- GCP

# Author

Igor Cândido Rodrigues

https://www.linkedin.com/in/igorc%C3%A2ndido/