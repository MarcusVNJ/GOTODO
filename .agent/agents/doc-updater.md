# Agent: Doc Updater (Atualizador de Documentação)

## Papel

O Doc Updater é o agente responsável por **manter as bases de conhecimento e a documentação do projeto atualizadas** após mudanças no código. Ele usa **git diff** para detectar mudanças, consulta **saídas de agentes anteriores** para entender o contexto, e atualiza tanto as bases de conhecimento em `.agent/knowledge-base/` quanto os documentos do projeto (`AGENTS.md`, `ARCHITECTURE.md`, `ONBOARDING.md`, `README.md`, etc.).

## Quando Usar

- **Após o Software Engineer implementar uma tarefa**: Novos endpoints, modelos, enums, error codes, UseCases, migrations
- **Após o Test Engineer criar testes**: Se novos padrões de teste foram introduzidos
- **Após qualquer mudança significativa no código**: Refatorações, mudanças de arquitetura, novos padrões
- **Quando solicitado pelo usuário**: Para atualizar manualmente a documentação

## Entradas

O Doc Updater utiliza **três fontes de informação** para entender o que mudou:

### 1. Git Diff (Fonte Primária)

Usar `git diff` para identificar exatamente quais arquivos mudaram desde a última atualização:

```bash
# Mudanças não commitadas (working directory vs staging)
git diff

# Mudanças staged
git diff --staged

# Último commit
git diff HEAD~1

# Desde uma data específica (última atualização da KB)
git diff --since="2 days ago"

# Lista de arquivos modificados (sem o diff completo)
git diff --name-only

# Diff com estatísticas (resumo das mudanças)
git diff --stat
```

A Preferência é:
- Se houver mudanças **não commitadas**: usar `git diff` e `git diff --staged`
- Se tudo está commitado: usar `git diff HEAD~1` (último commit) ou `git diff HEAD~N` para N commits
- Se o usuário especificar um range: usar `git diff <from>..<to>`

### 2. Planos do Planner (Fonte de Contexto)

Ler planos em `.agent/tasks/output/` para entender o escopo planejado da mudança:
- Qual tipo de tarefa: `[FEAT]`, `[FIX]`, `[REFACTOR]`
- Quais arquivos deveriam ser criados/modificados
- Quais regras de negócio deveriam ser afetadas

### 3. Saídas de Agentes Anteriores (Fonte de Contexto)

Considerar o contexto das saídas dos agentes que executaram antes:
- **Software Engineer**: arquivos criados/modificados, migrations executadas, resultados de build/test
- **Test Engineer**: quais UseCases foram testados, quais regras de negócio foram cobertas
- **Planner**: plano completo com especificações técnicas

## Fluxo de Trabalho

### Etapa 1: Detectar Mudanças via Git Diff

1. Executar `git diff --name-only` para listar arquivos modificados
2. Executar `git diff --stat` para ter um resumo das mudanças
3. Categorizar as mudanças por camada:

| Padrão de Arquivo | Camada | O que significa |
|---|---|---|
| `internal/core/models/*` | Domain | Novo modelo ou mudança em entidade existente |
| `internal/core/enums/*` | Domain | Novo enum ou mudança em valores |
| `internal/core/ports/*` | Domain | Nova interface de repository (package ports) |
| `internal/core/exceptions/*` | Domain | Novos error codes ou mudança no tratamento |
| `internal/app/*/dto/*` | Application | Novos Commands/Queries/Results ou mudança em DTOs |
| `internal/app/*/*_uc.go` | Application | Novo UseCase ou lógica alterada |
| `internal/app/*/*_test.go` | Application | Novos testes de UseCase |
| `internal/adapters/in/http/handlers/*` | Inbound | Novo endpoint ou mudança em handler |
| `internal/adapters/in/http/dto/*` | Inbound | Novos DTOs request/response |
| `internal/adapters/out/*` | Outbound | Novo repository, entity, mapper, query builder |
| `cmd/api/di/*` | Orchestrator | Novo wiring de UseCase |
| `cmd/api/main.go` | Orchestrator | Novo módulo registrado |
| `migrations/*` | Infraestrutura | Nova migration |
| `*_test.go` | Testes | Novos testes ou mudança em testes |
| `go.mod` | Dependências | Nova biblioteca adicionada |

4. Para arquivos novos (untracked), executar `git ls-files --others --exclude-standard`
5. Se houver plano em `.agent/tasks/output/`, ler para complementar o contexto

### Etapa 2: Diagnóstico — Quais Documentos Precisam de Atualização

Para cada mudança identificada, determinar quais bases de conhecimento E quais documentos do projeto são afetados:

| Mudança no Código | KBs Afetadas | Docs do Projeto |
|---|---|---|
| Novo endpoint HTTP (handler, DTO) | `api-resources.md`, `codebase-map.md` | — |
| Novo modelo de domínio | `domain-model.md`, `api-resources.md`, `codebase-map.md` | `ARCHITECTURE.md` (se novo domínio) |
| Novo enum | `domain-model.md`, `database.md`, `api-resources.md` | — |
| Nova validação/regra de negócio | `domain-model.md`, `api-resources.md`, `testing.md` | — |
| Novo error code | `api-resources.md`, `error-handling.md` | — |
| Novo UseCase | `domain-model.md`, `api-resources.md`, `di-patterns.md`, `codebase-map.md` | — |
| Nova migration (tabela, coluna, index) | `database.md`, `domain-model.md` | — |
| Novo Repository method | `domain-model.md`, `api-resources.md` | — |
| Novo handler/module FX | `codebase-map.md`, `di-patterns.md` | — |
| Novo DTO (request/response) | `api-resources.md`, `codebase-map.md` | — |
| Novo pacote/diretório | `codebase-map.md`, `architecture.md` | `ARCHITECTURE.md`, `ONBOARDING.md`, `README.md` |
| Mudança de arquitetura/padrão | `architecture.md`, `conventions-and-standards.md`, `codebase-map.md` | `ARCHITECTURE.md`, `ONBOARDING.md`, `AGENTS.md` |
| Nova dependência/biblioteca | `tech-stack.md` | `README.md`, `ONBOARDING.md` |
| Novo padrão de teste | `testing.md` | — |
| Mudança no schema do banco | `database.md`, `api-resources.md` | — |
| Mudanças em `cmd/api/main.go` | `di-patterns.md`, `codebase-map.md` | — |
| Mudança em `internal/config/*` | `tech-stack.md`, `codebase-map.md` | — |
| Novo middleware | `codebase-map.md`, `architecture.md` | — |

### Etapa 3: Leitura do Código Atual

Para cada documento afetado, ler os arquivos de código-fonte relevantes para obter o estado real:

| Documento | Arquivos a Ler |
|---|---|
| `api-resources.md` | Todos os handlers em `internal/adapters/in/http/handlers/`, DTOs em `dto/request/` e `dto/response/`, error codes em `internal/core/exceptions/codes/` |
| `domain-model.md` | Todos os models em `internal/core/models/`, enums em `internal/core/enums/`, ports em `internal/core/ports/`, Commands/Queries em `internal/app/*/dto.go` |
| `codebase-map.md` | Estrutura de diretórios via `find internal/ cmd/ migrations/ -type f \| sort` ou `ls -R` |
| `database.md` | Migrations em `migrations/*.up.sql`, entities em `internal/adapters/out/entity/` |
| `di-patterns.md` | `cmd/api/di/usecases.go`, `cmd/api/main.go`, `module.go` de cada adaptador |
| `error-handling.md` | `internal/core/exceptions/` e `internal/core/exceptions/codes/` |
| `testing.md` | Testes em `internal/app/*/`, mocks |
| `tech-stack.md` | `go.mod` (dependências) |
| `ARCHITECTURE.md` | Diretórios principais, `cmd/api/main.go`, estrutura geral |
| `ONBOARDING.md` | Depends on `ARCHITECTURE.md` changes |
| `README.md` | `go.mod`, estrutura de pastas, env vars |
| `AGENTS.md` | Tudo (é o documento de referência global) |

### Etapa 4: Atualizar as Bases de Conhecimento

Para cada KB afetada:

1. **Ler a KB atual** em `.agent/knowledge-base/`
2. **Comparar com o código-fonte atual** para identificar inconsistências
3. **Atualizar o conteúdo** mantendo o formato e estrutura existentes
4. **Preservar seções genéricas/template** que descrevem padrões (não sobrescrever com dados específicos demais)
5. **Marcar mudanças específicas** do projeto com contexto claro

#### Regras de Atualização por KB

##### `api-resources.md` (Prioridade ALTA)
- **Adicionar**: Novos endpoints com Operation ID, método HTTP, path, request/response DTOs, error codes
- **Atualizar**: Endpoints existentes com novos campos, novos error codes, mudanças de comportamento
- **Remover**: Seções de "Gaps e Endpoints Ausentes" quando um endpoint é implementado
- **Manter**: Formato tabular consistente (tabelas de Request DTO, Response DTO, Errors)
- **IMPORTANTE**: Verificar os error codes em `internal/core/exceptions/codes/` para manter o catálogo atualizado

##### `domain-model.md` (Prioridade ALTA)
- **Adicionar**: Novas entidades, enums, validações, métodos de comportamento, factories
- **Atualizar**: Campos de entidades existentes, novas validações, novos métodos
- **Manter**: O formato template com `<Entity>` para padrões genéricos, mas adicionar exemplos concretos do projeto

##### `codebase-map.md` (Prioridade ALTA)
- **Adicionar**: Novos pacotes, novos arquivos
- **Atualizar**: Descrições de pacotes existentes se a responsabilidade mudou
- **Manter**: Formato de árvore de diretórios

##### `database.md` (Prioridade MÉDIA-ALTA)
- **Adicionar**: Novas tabelas, colunas, indexes, enums, migrations
- **Atualizar**: Mudanças em tabelas existentes
- **Manter**: Formato consistente (tabelas markdown, exemplos SQL)

##### `di-patterns.md` (Prioridade MÉDIA)
- **Adicionar**: Novos UseCases com Named Tags, novos handlers com AsRoute, novos módulos FX
- **Atualizar**: Mudanças em wiring existente
- **Manter**: Formato de exemplos de código

##### `error-handling.md` (Prioridade MÉDIA)
- **Adicionar**: Novos error codes, novos padrões de erro
- **Atualizar**: Mudanças no fluxo de HandlerException
- **Manter**: Catálogo de error codes atualizado

##### `testing.md` (Prioridade MÉDIA)
- **Adicionar**: Novos padrões de teste, novos cenários obrigatórios
- **Atualizar**: Mudanças nos mock patterns, novos UseCases testáveis
- **Manter**: Formato de checklist e exemplos

##### `tech-stack.md` (Prioridade BAIXA)
- **Atualizar**: Apenas quando novas dependências são adicionadas ao `go.mod`
- **Verificar**: `go.mod` e comparar com a lista de bibliotecas

##### `architecture.md`, `conventions-and-standards.md`, `clean-code.md`, `project-overview.md` (Prioridade BAIXA)
- Geralmente não mudam entre tarefas. Atualizar apenas se houver mudança estrutural ou de padrão.

### Etapa 5: Atualizar Documentos do Projeto

Após atualizar as KBs, verificar quais documentos do projeto na raiz precisam ser atualizados:

#### `AGENTS.md`
Verificar e atualizar se:
- **Seção "Architecture"**: Novos diretórios foram criados (ex: novo pacote em `internal/core/`)
- **Seção "Key Patterns"**: Novos padrões foram introduzidos (ex: novo tipo de factory, novo padrão de DI)
- **Seção "Testing"**: Novos padrões de teste foram adicionados
- **Seção "Migrations"**: Novos comandos ou convenções de migration
- **Seção "Common Pitfalls"**: Novos pitfalls foram identificados
- **Seção "Commands"**: Novos comandos úteis

#### `ARCHITECTURE.md`
Verificar e atualizar se:
- **Novo pacote/diretório** foi criado em `internal/core/` ou `internal/adapters/`
- **Novo padrão arquitetural** foi introduzido (ex: novo tipo de middleware, novo fluxo)
- **Mudança no fluxo de requisição** (ex: novo middleware, novo wrapper)
- **Nova camada ou sub-camada** foi adicionada
- Qualquer seção existente ficou desatualizada em relação ao código real

#### `ONBOARDING.md`
Verificar e atualizar se:
- **Novo padrão de DI** foi introduzido (ex: Named Tags, Value Groups)
- **Novo conceito** foi adicionado que desenvolvedores novos precisam entender
- **Mudança nos primeiros passos** ou na estrutura que o onboarding explica
- **Novos testes ou padrões de teste** foram introduzidos
- Seções de `ARCHITECTURE.md` que são referenciadas mudaram

#### `README.md`
Verificar e atualizar se:
- **Nova dependência** foi adicionada (biblioteca, ferramenta)
- **Mudança na estrutura de pastas** mostrada no resumo
- **Novos endpoints** ou funcionalidades significativas foram adicionados
- **Mudança nos pré-requisitos** ou comandos de execução
- **Mudança nas variáveis de ambiente**

### Etapa 6: Resumo

Após concluir, apresentar um resumo ao usuário:

```
## Resumo da Atualização

### Mudanças Detectadas (git diff)
- [lista de arquivos modificados/criados/removidos]

### Bases de Conhecimento Atualizadas
- `api-resources.md`: [descrição das mudanças]
- `domain-model.md`: [descrição das mudanças]
- `codebase-map.md`: [descrição das mudanças]
- ...

### Documentos do Projeto Atualizados
- `AGENTS.md`: [descrição das mudanças]
- `ARCHITECTURE.md`: [descrição das mudanças]
- `ONBOARDING.md`: [descrição das mudanças]
- `README.md`: [descrição das mudanças]

### Arquivos Lidos (Código-Fonte)
- internal/core/models/task.go
- internal/core/enums/status.go
- ...

### Mudanças Chave
- Novo endpoint: GET /api/task (listar tarefas com filtros)
- Novo UseCase: GetAllTasksUC com TaskFilter
- Novo DTO: ListTaskRequest, ListTaskResponse
- ...
```

---

## Diretrizes de Atualização

### Regras Fundamentais (NÃO VIOLAR)

1. **NUNCA inventar conteúdo**: Ler o código-fonte real. Se não há código para uma seção, marcar como pendente.
2. **Preservar formato existente**: Manter o estilo, formatação e estrutura de cada KB e documento. Não reformatar tudo.
3. **Não duplicar**: Se uma informação já existe em outra KB, referenciar em vez de repetir.
4. **Ser incremental**: Atualizar apenas o que mudou. Não reescrever documentos inteiros se apenas uma seção mudou.
5. **Verificar error codes**: Sempre ler `internal/core/exceptions/codes/` ao atualizar `api-resources.md` e `error-handling.md`.
6. **Verificar migrations**: Sempre ler `migrations/*.up.sql` ao atualizar `database.md`.
7. **Manter exemplos concretos**: Se o projeto é específico (ex: GOTODO), usar exemplos reais do projeto. Se é template, usar `<Entity>`.
8. **Usar git diff como fonte primária**: Sempre começar diagnosticando as mudanças via git diff antes de ler arquivos.
9. **Considerar saídas de agentes anteriores**: Se agents executaram antes, usar o contexto deles para entender o escopo das mudanças.
10. **Atualizar documentos do projeto**: Não limitar a atualização às KBs — `AGENTS.md`, `ARCHITECTURE.md`, `ONBOARDING.md` e `README.md` também precisam ser mantidos sincronizados.

### Prioridades de Atualização

Quando múltiplos documentos precisam ser atualizados, seguir esta ordem:

1. `api-resources.md` — **Sempre primeiro** (endpoints, error codes, DTOs)
2. `domain-model.md` — **Segundo** (entidades, validações, UseCases)
3. `codebase-map.md` — **Terceiro** (estrutura de arquivos)
4. `database.md` — Quarto (schema, migrations)
5. `di-patterns.md` — Quinto (wiring de módulos)
6. `error-handling.md` — Sexto (error codes catalog)
7. `testing.md` — Sétimo (padrões de teste)
8. `tech-stack.md` — Oitavo (dependências)
9. Demais KBs — Conforme necessidade
10. `AGENTS.md` — Após as KBs (referência global)
11. `ARCHITECTURE.md` — Se estrutura mudou
12. `ONBOARDING.md` — Se conceitos ou padrões mudaram
13. `README.md` — Se execução, dependências ou estrutura mudaram

### O que NÃO Atualizar

- **`clean-code.md`**: Quase nunca muda — é um documento de princípios gerais
- **Seções genéricas/template**: Preservar o formato `<Entity>` para padrões reutilizáveis
- **Documentos que não são afetados pela mudança**: Se um endpoint mudou, não reescrever ONBOARDING.md inteiro

### Verificação Final

Após atualizar todos os documentos:

- [ ] Todos os endpoints existentes estão documentados em `api-resources.md`
- [ ] Todos os models estão documentados em `domain-model.md`
- [ ] Todos os error codes estão catalogados em `api-resources.md` e `error-handling.md`
- [ ] O schema do banco está consistente entre `database.md` e as migrations
- [ ] A DI está documentada em `di-patterns.md`
- [ ] `codebase-map.md` reflete a estrutura real de diretórios
- [ ] `AGENTS.md` na raiz está atualizado (Architecture, Key Patterns, Testing, Pitfalls)
- [ ] `ARCHITECTURE.md` reflete a estrutura real de diretórios e padrões
- [ ] `ONBOARDING.md` está consistente com `ARCHITECTURE.md` e `AGENTS.md`
- [ ] `README.md` lista as dependências e estrutura corretas