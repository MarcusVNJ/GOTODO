Arquitetura do Projeto GOTODO

1. Visão Geral

O GOTODO é uma API RESTful para o gerenciamento de tarefas (estilo Kanban) desenvolvida em Go. Para garantir que o
sistema seja escalável, altamente testável e resistente à deterioração do código ao longo do tempo, adotamos a
Arquitetura Hexagonal (também conhecida como Ports and Adapters).

O princípio fundamental desta arquitetura é que as regras de negócio devem ser o centro da aplicação. As dependências
devem sempre apontar "de fora para dentro". O domínio (Core) não sabe que existe um banco de dados PostgreSQL ou que
está respondendo a requisições HTTP; ele se comunica com o mundo externo puramente através de interfaces (Ports).

2. Estrutura de Diretórios e Fronteiras

A separação de responsabilidades reflete rigorosamente nossa árvore de diretórios:

    cmd/api/: O ponto de entrada da aplicação e nossa Raiz de Composição. O main.go e o subpacote router detêm a responsabilidade exclusiva de inicializar dependências e acoplar os adaptadores ao Core.

    internal/config/: Lida com a leitura tipada de variáveis de ambiente (12-Factor App).

    internal/core/ (O Hexágono): O coração do sistema. Contém a lógica pura.

        base/: Estruturas genéricas de base (ex: abstrações de UseCases).

        enums/: Definições estáticas do negócio (ex: status das tarefas).

        models/: Entidades de negócio ricas (ex: Task), contendo invariantes e regras de transição de status do Kanban.

        ports/: Interfaces que definem os contratos de entrada (UseCases) e saída (Repositories).

        usecase/: Orquestradores da lógica da aplicação.

        exceptions/: Erros de negócio tipados (ex: BusinessException), totalmente agnósticos ao protocolo HTTP.

    internal/adapters/ (A Infraestrutura): Tradutores entre o mundo externo e o Core.

        in/http/: Adaptadores primários (Handlers/Resources, DTOs de request/response e Middlewares).

        out/infrastructure/: Adaptadores secundários. Contém entity (modelos exclusivos de DB), mappers (tradutores de Domain <-> Entity), query_builder (geradores de SQL) e implementações do repositório PostgreSQL.

    migrations/: Scripts SQL puros (Up/Down) processados pelo golang-migrate, mantidos na raiz para garantir que a infraestrutura de banco de dados não polua o código interno.

3. Padrões de Engenharia e Design
   3.1. Tratamento Centralizado de Erros (HandlerException)

Implementamos um wrapper genérico (`middlewares.HandlerException`) que envelopa todos os controladores da API. Como utilizamos o ecossistema Huma, nossos Resources implementam assinaturas limpas no formato `func(ctx context.Context, input *I) (*O, error)`. O middleware intercepta qualquer `error` retornado:

    Erros do Cliente (4xx): Falhas de validação de input (ex: constraints de tamanho ou formatação JSON) são tratadas nativamente pelo Huma. Erros de negócio gerados pelo Core (BusinessException) são formatados para a RFC 7807 via interface `huma.StatusError`. Não geram logs de nível Error, evitando falsos-positivos na observabilidade.

    Erros de Servidor (5xx): Falhas técnicas ou inesperadas são encapsuladas pela `UnexpectedException`. O wrapper aciona um log assíncrono (em uma *goroutine* apartada) utilizando `slog`, contendo a stack trace completa rica gerada pela biblioteca `oops`, e o cliente final recebe apenas uma mensagem genérica de erro interno.

3.2. Segregação de Modelos (DTO -> Domain -> Entity)

Para proteger o encapsulamento, utilizamos modelos diferentes para cada fase da jornada dos dados:

    DTO (adapters/in/http/dto): Usado unicamente para parsing e validação de payloads HTTP JSON.

    Domain (core/models): Usado para a lógica de negócio pura. Não possui tags de serialização (json ou db).

    Entity (adapters/out/infrastructure/entity): Usado pelo repositório e pelo query_builder para espelhar as tabelas do banco de dados (contendo tags db). A conversão ocorre no pacote mappers.

6. O Fluxo de uma Requisição

Ao executar um fluxo (ex: Criar Tarefa), a requisição atravessa as camadas da seguinte forma:

    Entrada: A requisição atinge o framework `go-chi`, passa pela serialização do `huma` e atinge o wrapper `HandlerException`.

    Transporte (adapters/in): O Huma valida automaticamente o corpo e os parâmetros contra o schema OpenAPI usando os DTOs (structs `Input`). Se válido, invoca o Resource. O Resource mapeia para o modelo do Core, e o UseCase é acionado.

    Lógica (core/usecase): O UseCase orquestra as regras. Se todas as validações do modelo de negócio passarem, ele invoca a interface do Repositório (Porta de saída).

    Persistência (adapters/out): A implementação concreta (task_repository_postgres_impl.go) recebe o modelo de domínio, usa o Mapper para convertê-lo em TaskEntity, aciona o TaskQueryBuilder para montar o SQL e executa no pool do Postgres.

    Retorno: O sucesso (ou falha) borbulha de volta pela pilha até o Resource ou o ExceptionHandler formartar a resposta HTTP correspondente.

7. 🛡️ Tratamento Centralizado de Erros no GOTODO

O projeto GOTODO utiliza a abordagem de Middleware/Wrapper atrelada à declaração de rotas do Huma para garantir um tratamento de exceções centralizado, uniforme e rápido.

A mecânica rigorosa de tratamento de erros funciona da seguinte forma:

    Roteamento Padronizado (Huma): A declaração de rotas ocorre via `huma.Register`. Isso garante que o input seja validado dinamicamente antes de bater no core da aplicação. No momento do registro, injetamos a nossa função `middlewares.HandlerException`.

    Assinatura Limpa nos Handlers: Os manipuladores de rota (Resources) estão completamente limpos. Eles não mexem em `http.ResponseWriter` ou definem Status Code manualmente. Em caso de falha, apenas retornam `nil, err`.

    O Interceptador e Logging Assíncrono (`HandlerException`): Sua única responsabilidade é capturar os erros, realizar os castings necessários (para `BusinessException` ou `UnexpectedException`), disparar em uma Thread separada (goroutine) o registro no sistema de logs (`slog`), e devolver o erro estruturado para o Huma formatar de acordo com a RFC 7807 (Problem Details).

1. Erros do Cliente e Regras de Negócio (Status 4xx)

   Este fluxo trata falhas de requisição ou exceções geradas pelo Core (ex: `BusinessException`). Graças à interface `GetStatus() int`, a nossa exceção customizada indica ao Huma exatamente qual status HTTP devolver (como 400, 404, 409).

   O Huma gera um JSON RFC 7807. Esses erros são logados em modo assíncrono na severidade `INFO` (já que são fluxos de negócio naturais), protegendo os alertas da equipe.

2. Erros Técnicos e Falhas de Infraestrutura (Status 5xx)

   Problemas não mapeados disparam o `UnexpectedException`. Para garantir segurança total, os detalhes originais e a stack trace (`oops`) são logados localmente via `slog.Error` dentro de uma *goroutine*. O cliente, no entanto, visualiza apenas o genérico 500, bloqueando o vazamento de detalhes estruturais.