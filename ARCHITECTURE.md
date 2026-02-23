Arquitetura do Projeto GOTODO
1. Visão Geral

O GOTODO é uma API RESTful para o gerenciamento de tarefas (estilo Kanban) desenvolvida em Go. Para garantir que o sistema seja escalável, altamente testável e resistente à deterioração do código ao longo do tempo, adotamos a Arquitetura Hexagonal (também conhecida como Ports and Adapters).

O princípio fundamental desta arquitetura é que as regras de negócio devem ser o centro da aplicação. As dependências devem sempre apontar "de fora para dentro". O domínio (Core) não sabe que existe um banco de dados PostgreSQL ou que está respondendo a requisições HTTP; ele se comunica com o mundo externo puramente através de interfaces (Ports).
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
3.1. Tratamento Centralizado de Erros (Exception Handler)

Implementamos um adaptador de manipulador customizado (ExceptionHandler). Nossos Resources implementam a assinatura func(w, r) error. O middleware intercepta o erro retornado:

    Erros do Cliente (4xx): Falhas de conversão de JSON (ex: string enviada em campo int) ou erros de negócio gerados pelo Core são interceptados, convertidos para o Status HTTP correspondente (400, 404, 409) e retornados como JSON estruturado. Não geram logs de nível Error, evitando falsos-positivos na observabilidade.

    Erros de Servidor (5xx): Falhas técnicas são logadas pelo slog com todo o contexto (stack trace, path, método) e o cliente recebe apenas uma mensagem sanitizada genérica.

3.2. Segregação de Modelos (DTO -> Domain -> Entity)

Para proteger o encapsulamento, utilizamos modelos diferentes para cada fase da jornada dos dados:

    DTO (adapters/in/http/dto): Usado unicamente para parsing e validação de payloads HTTP JSON.

    Domain (core/models): Usado para a lógica de negócio pura. Não possui tags de serialização (json ou db).

    Entity (adapters/out/infrastructure/entity): Usado pelo repositório e pelo query_builder para espelhar as tabelas do banco de dados (contendo tags db). A conversão ocorre no pacote mappers.

4. O Fluxo de uma Requisição

Ao executar um fluxo (ex: Criar Tarefa), a requisição atravessa as camadas da seguinte forma:

    Entrada: A requisição atinge o framework go-chi, passando por middlewares globais e caindo no nosso AppRouter customizado.

    Transporte (adapters/in): O Resource executa o decode do JSON para o DTO. Falhas de parse retornam imediatamente. Ocorrendo sucesso, o DTO é mapeado para o modelo do Core, e o UseCase é acionado.

    Lógica (core/usecase): O UseCase orquestra as regras. Se todas as validações do modelo de negócio passarem, ele invoca a interface do Repositório (Porta de saída).

    Persistência (adapters/out): A implementação concreta (task_repository_postgres_impl.go) recebe o modelo de domínio, usa o Mapper para convertê-lo em TaskEntity, aciona o TaskQueryBuilder para montar o SQL e executa no pool do Postgres.

    Retorno: O sucesso (ou falha) borbulha de volta pela pilha até o Resource ou o ExceptionHandler formartar a resposta HTTP correspondente.