# Tech Stack

## Linguagem

- **Go 1.26+**
- **Arquitetura**: Hexagonal (Ports & Adapters) com 4 camadas: Core → App → Adapters.In → Adapters.Out
- **DI**: uber-go/fx (apenas nos adapters e composition root; Core e App são imaculados)

## Frameworks e Bibliotecas

### HTTP & API
| Biblioteca | Propósito |
|---|---|
| `go-chi/chi/v5` | Roteador HTTP subjacente |
| `danielgtaylor/huma/v2` | Framework REST API com geração automática de OpenAPI/Swagger, validação de input, serialização |

### Injeção de Dependência
| Biblioteca | Propósito |
|---|---|
| `uber-go/fx` | DI framework com lifecycle management, Provide vs Invoke, Value Groups, Named Instances |

### Banco de Dados
| Biblioteca | Propósito |
|---|---|
| `jackc/pgx/v5` | Driver e connection pool PostgreSQL |
| `Masterminds/squirrel` | Query builder fluente para SQL programático |
| `golang-migrate/migrate` | Versionamento do esquema do banco de dados (CLI) |
| `lib/pq` | Compatibilidade com driver stdlib |

### Configuração e Utilidades
| Biblioteca | Propósito |
|---|---|
| `kelseyhightower/envconfig` | Mapeamento tipado de variáveis de ambiente (Fail-Fast) |
| `joho/godotenv` | Carregamento de arquivo `.env` |
| `rs/xid` | Geração de IDs únicos globais |
| `samber/oops` | Error wrapping com contexto e stack traces ricas |
| `log/slog` | Logging estruturado nativo do Go |

### Testes
| Biblioteca | Propósito |
|---|---|
| `stretchr/testify` | Assertions e mock framework para testes unitários |

## Infraestrutura

| Componente | Detalhe |
|---|---|
| **Banco de Dados** | PostgreSQL (via Docker, imagem alpine) |
| **Migrations** | Scripts SQL puros (Up/Down) gerenciados por golang-migrate |

## Variáveis de Ambiente

| Variável | Tipo | Default | Descrição |
|---|---|---|---|
| `ENVIRONMENT` | string | `"development"` | Ambiente de execução |
| `PORT` | string | `"8080"` | Porta do servidor HTTP |
| `DATABASE_URL` | string | (obrigatório) | Connection string PostgreSQL |
| `ENABLE_DOCS` | bool | `true` | Habilita Swagger/OpenAPI |

## Endpoints

A definição de rotas segue o padrão `/api/<recurso>` com métodos RESTful. Cada entidade possui:
- `POST /api/<recurso>` — Criar (201)
- `GET /api/<recurso>/{id}` — Buscar por ID (200)
- `PUT /api/<recurso>` — Atualizar (200)
- `DELETE /api/<recurso>/{id}` — Excluir (204)

## Documentação Automática

- **Swagger UI**: `http://localhost:<PORT>/docs`
- **OpenAPI JSON**: `http://localhost:<PORT>/openapi.json`
- Requer `ENABLE_DOCS=true` no ambiente