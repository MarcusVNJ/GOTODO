# MCP — Database Schema Server

Servidor MCP (Model Context Protocol) que permite ao agente Planner acessar o schema do banco de dados PostgreSQL e os modelos de domínio do projeto em tempo real.

## Ferramentas Disponíveis

| Ferramenta | Descrição |
|---|---|
| `list_tables` | Lista todas as tabelas do banco com seus schemas |
| `describe_table` | Descreve a estrutura completa de uma tabela (colunas, tipos, constraints, defaults, nullable, PK) |
| `list_enums` | Lista todos os tipos ENUM do PostgreSQL com seus valores |
| `list_indexes` | Lista todos os indexes de uma tabela com colunas e tipo (unique ou não) |
| `list_foreign_keys` | Lista todas as foreign keys de uma tabela com a tabela e coluna referenciada |
| `database_schema_overview` | Visão geral completa: todas as tabelas com colunas, enums, flags de soft delete |
| `read_domain_models` | Lê todos os modelos de domínio em `internal/core/models/` (campos, tipos, métodos, factories) |
| `read_domain_enums` | Lê todos os enums em `internal/core/enums/` com seus valores |
| `read_error_codes` | Lê todos os códigos de erro em `internal/core/exceptions/codes/` |

## Configuração

### Variáveis de Ambiente

| Variável | Default | Descrição |
|---|---|---|
| `DB_HOST` | `localhost` | Host do PostgreSQL |
| `DB_PORT` | `5432` | Porta do PostgreSQL |
| `DB_USER` | `postgres` | Usuário do PostgreSQL |
| `DB_PASSWORD` | `root` | Senha do PostgreSQL |
| `DB_NAME` | `todo_db` | Nome do banco de dados |
| `DB_SSLMODE` | `disable` | SSL mode do PostgreSQL |

### Uso com opencode / Claude Desktop

Adicione ao seu arquivo de configuração MCP (`mcp.json` ou `claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "gotodo-db": {
      "command": "go",
      "args": ["run", ".agent/mcp/db_schema_server.go"],
      "env": {
        "DB_HOST": "localhost",
        "DB_PORT": "5432",
        "DB_USER": "postgres",
        "DB_PASSWORD": "root",
        "DB_NAME": "todo_db",
        "DB_SSLMODE": "disable"
      }
    }
  }
}
```

> **Nota**: Ajuste as variáveis de ambiente conforme a configuração do seu banco. O MCP usa as mesmas credenciais do `.env` do projeto.

## Pré-requisitos

- Go 1.26+
- PostgreSQL rodando e acessível
- Dependências Go: `github.com/mark3labs/mcp-go` e `github.com/lib/pq`

### Instalar dependências

```bash
cd .agent/mcp
go mod init gotodo-mcp
go mod tidy
```

## Exemplo de Uso pelo Agente Planner

1. **Antes de criar uma nova entidade**: Usar `read_domain_models` para verificar modelos existentes e `list_tables` para ver tabelas no banco
2. **Ao definir o schema de migration**: Usar `database_schema_overview` para ver o estado atual do banco e `list_enums` para verificar enums existentes
3. **Ao criar foreign keys**: Usar `list_foreign_keys` para ver relacionamentos atuais
4. **Ao definir error codes**: Usar `read_error_codes` para não duplicar códigos
5. **Ao estender uma entidade existente**: Usar `describe_table` + `list_indexes` para ver a estrutura completa