# ADR-001: Arquitetura de Monolito Go com Vertical Slicing + Hexagonal

**Status:** aceito  
**Data:** 2026-05-08  
**Decisores:** @Diego Hernan Demitto Czajka
**Tags:** [arquitetura, domínio, infra, go]

---

## Contexto

O projeto é um monolito escrito em Go que precisa crescer de forma sustentável sem se tornar um *big ball of mud*. As principais pressões são:

- **Time pequeno** com necessidade de autonomia por feature, evitando conflitos de merge constantes em camadas técnicas compartilhadas.
- **Domínio em evolução** — regras de negócio mudam frequentemente e precisam ser isoladas de detalhes de infraestrutura (HTTP, banco de dados, filas).
- **Possibilidade futura de extração de serviços** — o monolito pode eventualmente ter slices promovidos a microsserviços independentes.
- **Go idiomático** — a linguagem favorece interfaces implícitas e composição simples; a arquitetura deve respeitar isso em vez de lutar contra.

A ausência de uma estrutura definida levaria a handlers acessando o banco diretamente, regras de negócio espalhadas, e dificuldade crescente para testar e evoluir o código.

## Decisão

Adotaremos **Vertical Slicing como estrutura de pastas** combinado com **Arquitetura Hexagonal (Ports & Adapters) dentro de cada slice**, usando a seguinte stack:

- **Router HTTP:** Chi (`github.com/go-chi/chi`)
- **Banco de dados:** sqlc + pgx (PostgreSQL)
- **Migrações:** golang-migrate
- **Logs:** slog (stdlib Go 1.21+)
- **Injeção de dependência:** manual (constructor injection)

A estrutura de diretórios padrão do projeto será:

```
/cmd
  /api
    main.go                  → wire manual: config → db → repos → services → handlers → router

/internal
  /infra
    /database
      postgres.go            → abre conexão e pool pgx, ping de healthcheck
      migrations.go          → executa golang-migrate no boot da aplicação
    /config
      config.go              → struct Config, leitura de env vars

  /shared                    → tipos primitivos reutilizáveis entre slices (manter mínimo)
    id.go                    → ex: type UserID uuid.UUID
    money.go                 → ex: type Money int64

  /<slice>                   → um diretório por domínio/feature (ex: user, order, payment)
    /domain
      <entity>.go            → entidade pura, zero imports externos (só stdlib)
    /ports
      repository.go          → interface do repositório (outbound port)
      service.go             → interface do serviço / use cases (inbound port)
    /adapters
      /http
        handler.go           → Chi handler, traduz HTTP ↔ domínio
      /postgres
        repository.go        → implementação sqlc/pgx da interface de repositório
    service.go               → implementação dos use cases (só conhece domain e ports)

/migrations                  → arquivos SQL versionados (golang-migrate)
```

Exemplo concreto com dois slices:

```
/internal
  /infra
    /database
      postgres.go
      migrations.go
    /config
      config.go

  /shared
    id.go
    money.go

  /user
    /domain
      user.go
    /ports
      repository.go          → interface UserRepository
      service.go             → interface UserService
    /adapters
      /http
        handler.go
      /postgres
        repository.go
    service.go               → CreateUser, GetUser, DeactivateUser...

  /order
    /domain
      order.go
    /ports
      repository.go          → interface OrderRepository
      service.go             → interface OrderService
    /adapters
      /http
        handler.go
      /postgres
        repository.go
    service.go               → CreateOrder, CancelOrder...
```

O `main.go` é o único lugar que conhece todas as dependências e as conecta:

```go
// cmd/api/main.go
cfg  := config.Load()
db   := database.Connect(cfg.DatabaseURL)

userRepo    := userpostgres.NewRepository(db)
userSvc     := user.NewService(userRepo)
userHandler := userhttp.NewHandler(userSvc)

orderRepo    := orderpostgres.NewRepository(db)
orderSvc     := order.NewService(orderRepo)
orderHandler := orderhttp.NewHandler(orderSvc)

r := chi.NewRouter()
r.Mount("/users",  user.NewRouter(userHandler))
r.Mount("/orders", order.NewRouter(orderHandler))
```

O fluxo de dependência segue estritamente a direção:

```
handler → (port) Service → (port) Repository → postgres
   ↑                                               ↑
adapter/http                               adapter/postgres
```

Nenhuma camada interna (domain, service) conhece adapters. Apenas o `main.go` enxerga a aplicação inteira.

## Consequências

### Positivas
- **Isolamento real do domínio** — entidades em `/domain` não importam `net/http`, `database/sql` ou qualquer lib externa; são testáveis puras.
- **Substituição de adapters sem impacto no domínio** — trocar PostgreSQL por outra engine ou trocar HTTP por gRPC exige apenas um novo adapter, sem tocar em `service.go` ou `/domain`.
- **Autonomia por feature** — times ou devs diferentes podem trabalhar em slices distintos com mínimo de conflito.
- **Go idiomático** — interfaces implícitas do Go satisfazem os ports naturalmente, sem anotações ou frameworks de DI.
- **Caminho claro para extração** — cada slice já carrega domínio + ports + adapters; extrair como microsserviço é remover o diretório e criar um novo repositório.

### Negativas / Trade-offs
- **Mais arquivos por feature** — um CRUD simples já gera 4–5 arquivos. Para features muito pequenas pode parecer over-engineering.
- **Curva de aprendizado inicial** — devs acostumados com estrutura MVC (controllers/models/services flat) precisam internalizar o modelo de ports.
- **Shared Kernel requer disciplina** — tipos compartilhados entre slices (ex: `Money`, `UserID`) devem viver em `/internal/shared` e não criar dependência circular entre slices.
- **sqlc requer regeneração** — mudanças no SQL exigem rodar `sqlc generate`; precisa estar no pipeline de desenvolvimento (Makefile/script).

### Neutras / Notas de implementação
- Slices **não devem importar pacotes internos uns dos outros diretamente**. Comunicação entre slices deve ocorrer via interface (port) ou eventos de domínio.
- `/internal/infra` é infraestrutura **compartilhada e técnica** (conexão de banco, config); não contém regra de negócio.
- `/internal/shared` abriga apenas tipos primitivos de domínio reutilizáveis (IDs, value objects). Deve ser mantido mínimo e estável — crescimento excessivo é sinal de que algo pertence a um slice específico.
- `service.go` na raiz do slice **é** o use case. Só cria `/usecases` como subpacote se o arquivo crescer demais; nesse caso `service.go` vira um facade que delega aos use cases individuais.
- A injeção de dependência é feita via construtores no `main.go`. Para projetos que crescerem muito, considerar Wire (ADR futuro).
- Testes de unidade ficam junto ao código (`service_test.go`, `handler_test.go`); testes de integração em `/test/integration`.
- O diretório `/domain` de cada slice **nunca** deve importar nada fora da stdlib Go.

## Alternativas consideradas

| Alternativa | Por que foi descartada |
|---|---|
| Estrutura flat por camada técnica (`/handlers`, `/services`, `/repositories`) | Causa acoplamento implícito entre features; conflitos de merge crescem com o time; dificulta extração futura de serviços. |
| Hexagonal puro sem Vertical Slicing | Agrupa por camada técnica globalmente (ex: `/adapters/http` com todos os handlers juntos); perde a autonomia por feature e dificulta navegação. |
| Vertical Slicing sem Hexagonal interno | Sem a separação de ports, adapters tendem a vazar para o domínio; handlers acessam o banco diretamente ao longo do tempo. |
| Fiber como framework HTTP | Alta performance, mas **não é compatível com `net/http`**; cria lock-in e impede uso de middlewares padrão da comunidade Go. |
| GORM no lugar de sqlc + pgx | Maior abstração, mas queries inesperadas em runtime, N+1 silencioso e menor performance; sqlc garante SQL explícito e tipagem em compile time. |

## Referências

- [Chi Router](https://github.com/go-chi/chi)
- [sqlc — Generate type-safe Go from SQL](https://sqlc.dev/)
- [pgx — PostgreSQL Driver](https://github.com/jackc/pgx)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Hexagonal Architecture — Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Vertical Slice Architecture — Jimmy Bogard](https://www.jimmybogard.com/vertical-slice-architecture/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

---

