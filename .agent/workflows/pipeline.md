# Workflow: Pipeline de Desenvolvimento Completo

## Visão Geral

Pipeline sequencial que transforma demandas brutas (histórias de usuário, bugs, specs) em código implementado, testado e documentado. Cada agente consome a saída do anterior e produz artefatos para o próximo.

```
Entrada (histórias, protótipos, specs)
    │
    ▼
┌─────────────────┐     .agent/tasks/output/task01_plan.md
│   1. Planner    │──────────────────────────────────────────┐
└─────────────────┘                                          │
                                                             ▼
┌──────────────────────┐     Código implementado + testes
│ 2. Software Engineer │─────────────────────────────────────┐
└──────────────────────┘                                     │
                                                             ▼
┌─────────────────┐     Testes de unidade criados
│ 3. Test Engineer│─────────────────────────────────────────┐
└─────────────────┘                                        │
                                                           ▼
┌───────────────┐     .agent/knowledge-base/*.md + docs do projeto atualizados
│ 4. Doc Updater│
└───────────────┘
```

## Ordem de Execução

```
planner → software-engineer → test-engineer → doc-updater
```

Cada etapa **deve completar com sucesso** antes da próxima iniciar. Se uma etapa falhar, o pipeline para e o erro deve ser corrigido antes de prosseguir.

---

## Etapa 1: Planner

### Agente
`@planner`

### Entrada
- Histórias de usuário (texto)
- Protótipos de tela (imagens)
- Especificações técnicas (documentos .md)
- Bugs reportados (texto)
- Qualquer anexo em `.agent/tasks/`

### O que faz
1. Analisa a demanda e identifica entidades, regras de negócio e dependências
2. Consulta as bases de conhecimento (especialmente `api-resources.md`)
3. Decompõe em tarefas acionáveis com tipo `[FEAT]`, `[FIX]`, `[REFACTOR]`, etc.
4. Gera planos detalhados em `.agent/tasks/output/`

### Saída
- `.agent/tasks/output/task01_plan.md` (e subsequentes se houver divisão)
- Cada plano contém: contexto, regras de negócio, specs técnicas, endpoints, error codes, critérios de aceite, ordem de implementação

### Critério de Sucesso
- Plano gerado em `.agent/tasks/output/` com todas as seções relevantes preenchidas
- Regras de negócio numeradas (RN01, RN02...)
- Endpoints especificados com DTOs, error codes e status HTTP
- Ordem de implementação definida (Core → Outbound → Inbound → Orchestrator)

### Como Invocar
```
@planner [descrever a história de usuário, bug ou feature]
```
Ou, com anexos:
```
@planner Analise o arquivo .agent/tasks/spec.md e os protótipos em .agent/tasks/ para criar o plano de implementação.
```

---

## Etapa 2: Software Engineer

### Agente
`@software-engineer`

### Entrada
- Planos gerados pelo Planner em `.agent/tasks/output/`
- Bases de conhecimento em `.agent/knowledge-base/`
- Código existente no repositório

### O que faz
1. Lê o plano da tarefa (ex: `task01_plan.md`)
2. Consulta as bases de conhecimento (arquitetura, convenções, DI, etc.)
3. Implementa na ordem: Core → Outbound → Inbound → Orchestrator
4. Cria migrations se necessário
5. Executa `go build ./...` e `go test ./...`

### Saída
- Código implementado em todas as camadas
- Migrations criadas e executadas (se necessário)
- `go build ./...` passando
- `go test ./...` passando (testes existentes continuam passando)

### Critério de Sucesso
- `go build ./...` compila sem erros
- `go test ./...` passa (todos os testes existentes)
- Código segue Arquitetura Hexagonal (Core sem imports de infra)
- Handlers não definem HTTP status code diretamente
- Soft delete em todas as operações de deleção

### Como Invocar
```
@software-engineer Implemente a tarefa task01_plan.md
```
Ou, para a tarefa mais recente:
```
@software-engineer Implemente a tarefa mais recente em .agent/tasks/output/
```

---

## Etapa 3: Test Engineer

### Agente
`@test-engineer`

### Entrada
- Planos do Planner em `.agent/tasks/output/` (para entender as regras de negócio)
- Código implementado pelo Software Engineer
- Bases de conhecimento (especialmente `testing.md`, `domain-model.md`, `api-resources.md`)

### O que faz
1. Lê o plano da tarefa para extrair regras de negócio
2. Lê o código-fonte dos Models e UseCases para identificar TODAS as regras
3. Verifica/atualiza o Mock do Repository se necessário
4. Cria testes de unidade focados em regras de negócio (padrão AAA com testify)
5. Executa `go test ./internal/core/usecase/... -v`

### Saída
- Arquivos de teste em `internal/core/usecase/<entity>/*_test.go`
- Mock atualizado em `internal/core/usecase/<entity>/<entity>_repository_mock_test.go` (se necessário)
- Todos os testes passando

### Critério de Sucesso
- `go test ./internal/core/usecase/<entity>/... -v` passa
- `go build ./...` compila sem erros
- Todas as regras de negócio identificadas têm testes correspondentes
- Error codes de BusinessException verificados com `assert.Equal(t, codes.XXX.Code(), bizErr.Code)`
- Package de teste é `<entity>_test` (black-box)

### Como Invocar
```
@test-engineer Crie testes para as regras de negócio da tarefa task01_plan.md
```
Ou, para a tarefa mais recente:
```
@test-engineer Crie testes de unidade para as regras de negócio implementadas
```

---

## Etapa 4: Doc Updater

### Agente
`@doc-updater`

### Entrada
- Mudanças de código (detectadas via `git diff`)
- Planos do Planner em `.agent/tasks/output/` (contexto do que foi implementado)
- Saídas dos agentes anteriores
- Bases de conhecimento atuais em `.agent/knowledge-base/`
- Documentação do projeto na raiz (`AGENTS.md`, `ARCHITECTURE.md`, `ONBOARDING.md`, `README.md`)

### O que faz
1. Executa `git diff` para detectar todas as mudanças no código
2. Categoriza as mudanças por camada (Domain, Inbound, Outbound, etc.)
3. Determina quais bases de conhecimento e documentos do projeto precisam ser atualizados
4. Atualiza incrementalmente cada documento afetado
5. Atualiza `AGENTS.md`, `ARCHITECTURE.md`, `ONBOARDING.md` e `README.md` se necessário

### Saída
- Bases de conhecimento atualizadas em `.agent/knowledge-base/`
- Documentação do projeto atualizada na raiz (quando aplicável)
- Resumo das mudanças documentadas

### Critério de Sucesso
- Todas as KBs afetadas pelas mudanças foram atualizadas
- `api-resources.md` reflete todos os endpoints existentes
- `domain-model.md` reflete todos os models, enums e UseCases
- `codebase-map.md` reflete a estrutura real de diretórios
- `AGENTS.md` atualizado se houve mudança em Architecture, Key Patterns, Testing ou Common Pitfalls
- `ARCHITECTURE.md` atualizado se houve mudança estrutural
- `ONBOARDING.md` atualizado se houve mudança de padrões ou conceitos
- `README.md` atualizado se houve mudança de dependências ou estrutura

### Como Invocar
```
@doc-updater Atualize a documentação após a implementação da tarefa task01_plan.md
```
Ou, genérico:
```
@doc-updater Atualize a documentação com base nas mudanças recentes
```

---

## Fluxo Completo: Exemplo Prático

### 1. Invocar o Planner
```
@planner
Como usuário, eu quero organizar minhas tarefas em pastas, para que eu possa agrupar tarefas por projeto ou contexto.

Critérios de Aceite:
- Dado que o usuário está na tela principal, quando clica em "Nova Pasta", então deve aparecer um formulário para criar a pasta
- Uma pasta deve ter nome obrigatório (1-150 caracteres)
- Uma pasta pode ter uma descrição opcional (0-500 caracteres)
- Pastas podem ser soft-deletadas
- Tarefas podem ser movidas entre pastas
```

**Saída**: `.agent/tasks/output/task01_plan.md` com regras de negócio, endpoints, DTOs, error codes, ordem de implementação

### 2. Invocar o Software Engineer
```
@software-engineer Implemente a tarefa task01_plan.md
```

**Saída**: Código implementado em 4 camadas, migrations executadas, build e testes passando

### 3. Invocar o Test Engineer
```
@test-engineer Crie testes para as regras de negócio da pasta
```

**Saída**: Testes de unidade cobrindo: caminho feliz, validações (nome vazio, descrição longa), not found, regras de estado, erros de repositório

### 4. Invocar o Doc Updater
```
@doc-updater Atualize a documentação após implementação de pastas
```

**Saída**: KBs atualizadas (api-resources.md com novos endpoints, domain-model.md com Folder, database.md com migration, etc.), AGENTS.md atualizado, ARCHITECTURE.md atualizado

---

## Fluxos Alternativos

### Apenas Implementação (sem planejamento)
Se o plano já existe e foi revisado:
```
@software-engineer → @test-engineer → @doc-updater
```

### Apenas Testes (código já implementado)
```
@test-engineer → @doc-updater
```

### Apenas Documentação (algo mudou manualmente)
```
@doc-updater
```

### Bug Fix (fluxo rápido)
```
@planner [FIX] → @software-engineer → @test-engineer → @doc-updater
```

### Refatoração
```
@planner [REFACTOR] → @software-engineer → @test-engineer → @doc-updater
```

---

## Notas Importantes

1. **Pré-requisito**: A base de conhecimento deve existir antes de iniciar o pipeline. Se `.agent/knowledge-base/` não existe, executar `@knowledge-base-creator` primeiro.
2. **Falhas bloqueantes**: Se qualquer etapa falhar (build, testes, etc.), corrigir antes de prosseguir para a próxima.
3. **Idempotência**: Cada agente pode ser invocado independentemente. O pipeline é uma convenção de ordem, não uma dependência rígida.
4. **Saidas são cumulativas**: O Software Engineer não remove o plano do Planner. O Test Engineer não remove o código. Cada agente adiciona ao que já existe.
5. **Git diff como entrada**: O Doc Updater usa `git diff` como fonte primária para detectar mudanças, complementado pelo plano do Planner e pelas saídas dos agentes anteriores.