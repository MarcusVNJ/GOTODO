GOTODO - API de Gestão de Tarefas (Kanban)

Este projeto foi construído com o propósito de demonstrar conhecimentos avançados na linguagem Go (Golang), aplicando os princípios rigorosos da Arquitetura Hexagonal (Ports and Adapters) e diversas técnicas modernas de Engenharia de Software.

O domínio principal da aplicação é um sistema de gerenciamento de tarefas baseado no fluxo Kanban, suportando controle de prioridades, auditoria e transições rigorosas de estado.
📚 Documentação do Projeto

Para manter o código limpo e o time alinhado, documentamos nossas decisões técnicas e fundamentos. Se você é novo no projeto, comece por aqui:

    🚀 Guia de Onboarding: Leitura obrigatória. Explica os fundamentos de memória do Go (Stack vs Heap) e os padrões específicos de fluxo e injeção de dependências do nosso código.

    🏛️ (./ARCHITECTURE.md): Detalha nossa estrutura de diretórios Hexagonal, o isolamento absoluto do domínio (Core) e como lidamos com a conversão de DTOs e Entities.

✨ Principais Tecnologias e Padrões

    Linguagem: Go 1.21+

    Arquitetura: Hexagonal (Ports & Adapters)

    Roteamento: go-chi/chi (Roteador leve, idiomático e compatível com net/http)

    Tratamento de Erros: Padrão Handler Adapter (tratamento centralizado de exceções em rotas web)

    Logs Estruturados: pacote nativo log/slog para rastreabilidade contextual e proteção contra vazamento de dados sensíveis.

    Configuração: kelseyhightower/envconfig para injeção tipada de variáveis de ambiente (Fail-Fast).

    Banco de Dados: PostgreSQL (via pgxpool).

    Query Builder: Masterminds/squirrel acoplado ao repositório para escrita programática e segura de SQL.

    Migrations: golang-migrate/migrate para versionamento do esquema de banco de dados.

🚀 Como executar o projeto localmente
1. Pré-requisitos

    Go instalado (versão 1.21+)

    Instância do PostgreSQL rodando localmente (ou via Docker)

    golang-migrate CLI instalada.

2. Configuração do Ambiente

Crie um arquivo .env na raiz do projeto baseando-se nas chaves abaixo:env
ENVIRONMENT=development
PORT=8080
DATABASE_URL=postgres://seu_user:sua_senha@localhost:5432/todo_db?sslmode=disable


### 3. Migrações de Banco de Dados
Antes de iniciar a aplicação, crie o esquema do banco de dados executando as migrations:
```bash
migrate -database "postgres://seu_user:sua_senha@localhost:5432/todo_db?sslmode=disable" -path migrations up

4. Rodando a Aplicação

Baixe as dependências e inicie o servidor a partir da Raiz de Composição:
Bash

go mod tidy
go run cmd/api/main.go

O servidor estará disponível e escutando na porta configurada.
📂 Estrutura de Pastas de Alto Nível

GOTODO/
├── cmd/
│   └── api/                  # Entry point e Injeção de Dependências (Composition Root)
├── internal/
│   ├── config/               # Load de variáveis de ambiente
│   ├── core/                 # REGRAS DE NEGÓCIO PURAS (Agnóstico à Web/DB)
│   │   ├── base/
│   │   ├── enums/
│   │   ├── exceptions/
│   │   ├── models/
│   │   ├── ports/
│   │   └── usecase/
│   └── adapters/             # INTEGRAÇÃO EXTERNA
│       ├── in/http/          # Handlers, Middlewares e DTOs (Driver Adapters)
│       └── out/infrastructure# Repositórios PG, Mappers, Entities e Query Builder
├── migrations/               # Scripts SQL (Up/Down) do golang-migrate
├──.env
├── go.mod
└── README.md

Desenvolvido como demonstração de engenharia de software robusta, escalável e testável em Go.