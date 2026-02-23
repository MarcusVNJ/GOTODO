Guia de Integração para Desenvolvedores - Projeto GOTODO

Bem-vindo ao projeto GOTODO! Este é o nosso sistema de gerenciamento de tarefas (estilo Kanban). Para garantir que o código permaneça escalável, limpo e de fácil manutenção, nós adotamos uma série de padrões avançados de engenharia de software e recursos específicos da linguagem Go.

Este documento foi criado para acelerar o seu entendimento sobre como o Go funciona por debaixo dos panos, como estruturamos nossa Arquitetura Hexagonal e como lidamos de forma elegante com os fluxos HTTP, banco de dados e erros.
1. Fundamentos Computacionais em Go: Memória e Desempenho

Diferente de linguagens interpretadas ou que rodam em máquinas virtuais densas, Go é compilado para código de máquina e nos dá grande controle sobre a memória. Para escrever um código eficiente aqui, você precisa entender como o Go gerencia a alocação de dados.
1.1. Stack vs. Heap

Toda vez que uma variável é criada, o Go precisa decidir onde guardá-la na memória RAM do servidor. Ele tem dois lugares principais: a Stack (Pilha) e o Heap (Monte).

    A Stack (A Mochila): É uma região de memória incrivelmente rápida que opera no formato LIFO (Last In, First Out). É como sua mochila para uma viagem de um dia: rápida para acessar, mas com tamanho limitado e esvaziada automaticamente assim que a função que você estava executando termina.

    O Heap (A Mala de Porão): É uma região de memória maior, usada para armazenamento dinâmico. É como uma mala grande que fica guardada mesmo depois que a "viagem" da função acaba. É mais lento e exige que o Garbage Collector (GC) do Go trabalhe periodicamente para limpar o que não está mais sendo usado.

1.2. Escape Analysis (Análise de Fuga)

A regra de ouro do Go é: o tamanho da variável não importa tanto quanto o tempo de vida dela. Durante a compilação, o Go rastreia o caminho de cada variável (Escape Analysis). Se o compilador provar que a variável só será usada dentro daquela função, ela fica na Stack. Mas se a variável "escapar" (ex: você retornar um ponteiro para ela, passá-la para uma goroutine ou para uma interface vazia any), o Go a joga para o Heap.
1.3. Ponteiros vs. Valores (Quando usar o quê?)

É comum que novos desenvolvedores Go passem tudo como ponteiros (*models.Task) achando que isso economiza memória. Isso é um mito.

    Passe por Valor: Para structs pequenas ou médias (como DTOs), o custo de copiar os dados é quase zero e mantém a variável na Stack, aliviando o trabalho do Garbage Collector. Nós preferimos a passagem por valor por padrão.

    Passe por Ponteiro: Use ponteiros apenas quando precisar modificar/mutar a estrutura original ou quando a struct for excepcionalmente gigantesca.

2. A Arquitetura Hexagonal (Ports and Adapters)

O projeto GOTODO utiliza a Arquitetura Hexagonal. O objetivo é garantir que as nossas regras de negócio nunca dependam de detalhes técnicos como bancos de dados (Postgres) ou frameworks web (Chi).
2.1. O Core (Domínio e UseCases)

É o coração da nossa aplicação, localizado em internal/core.

    Regra Absoluta: Nenhuma interface, struct ou função dentro do core pode importar bibliotecas de infraestrutura (como net/http, pgx ou json). O domínio é cego para o mundo externo.

    Possuímos pacotes como enums (para os status Kanban), models (entidades ricas com regras de negócio) e base (interfaces genéricas de caso de uso).

2.2. Ports (As Portas)

Para que o core consiga conversar com o banco de dados sem conhecê-lo, ele define Interfaces em internal/core/ports.

    O domínio apenas diz: "Eu preciso de alguém que implemente a função Save(task *models.Task)". Ele não se importa se quem vai salvar é o Postgres ou a memória.

2.3. Adapters (Os Adaptadores)

Ficam em internal/adapters. Eles são os tradutores entre o mundo externo e o nosso domínio.

    Inbound (in/http): Nossos Handlers que pegam o JSON da requisição, validam via DTOs e convertem para o modelo puro antes de chamar os Casos de Uso.

    Outbound (out/infrastructure): Nossa implementação do Postgres. Aqui usamos mappers para converter a entidade de domínio em uma TaskEntity e o query_builder para gerar o SQL.

3. Padrões Específicos do Projeto GOTODO
3.1. Raiz de Composição (Composition Root)

Não espalhamos injeção de dependência pelo código. A "colagem" das peças acontece exclusivamente no diretório cmd/api/router e main.go. É lá que instanciamos o banco, passamos pro repositório, que vai pro UseCase, que vai pro Handler.
3.2. Tratamento Centralizado de Erros

Nossos Handlers de rota retornam um error (ex: func(w http.ResponseWriter, r *http.Request) error).
Nós criamos um middleware/adaptador central (ExceptionHandler) que envolve todas as rotas. Se houver falha de validação de input (ex: enviar uma string num campo de inteiro gerando json.UnmarshalTypeError), ou um erro de negócio (CodeInvalidData), você apenas retorna o erro. O ExceptionHandler o traduz para um JSON padronizado com Status 400. Se for erro de banco, ele devolve um erro 500 genérico ao cliente e loga o detalhe técnico no servidor.
3.3. Configuração Orientada a Ambiente

Lemos configurações via variáveis de ambiente usando o pacote internal/config, populando uma struct fortemente tipada (AppConfig). Se faltar uma variável obrigatória, o sistema aplica o padrão Fail-Fast e não inicia.
3.4. Versionamento de Banco de Dados (Migrations)

Nunca criamos tabelas manualmente. Utilizamos a biblioteca golang-migrate/migrate.

    Onde ficam: Todos os scripts SQL ficam na pasta migrations/ na raiz do projeto, fora do internal (já que o Core não deve saber sobre SQL).

    Up e Down: Para cada alteração, há um arquivo .up.sql (aplica a mudança) e um .down.sql (reverte a mudança). Isso garante que o esquema do banco evolua de forma segura e rastreável junto com a API.

Primeiros Passos

    Abra a pasta internal/core/models e analise as regras de negócio do Kanban.

    Veja as rotas em cmd/api/router para entender como as dependências são injetadas.

    Observe o internal/adapters/in/http/handlers para ver como o fluxo HTTP é limpo graças ao nosso ExceptionHandler.