# Encoder

## SOBRE

## Estrutura do repositório

### Arquivos e pastas

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

### Explicação da arquitetura optada
O projeto foi desenvolido utilizando os princípios de Arquitetura Limpa e Arquitetura Hexagonal, onde separamos os nossos domínios, regras de negócio, casos de uso, e aplicação web em pastas seraradas. 



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