# ADR-004: Migração de Clerk para autenticação JWT self-hosted

**Status:** aceito
**Data:** 2026-07-19
**Decisores:** @Diego Hernan Demitto Czajka
**Tags:** [arquitetura, segurança, autenticação, infra, jwt]
**Supersede:** [ADR-003 — Migração de Keycloak para Clerk](ADR-003-clerk-auth-migration.md)

---

## Contexto

O [ADR-003](ADR-003-clerk-auth-migration.md) adotou o Clerk como Identity Provider SaaS, delegando login, registro e emissão de JWT a um serviço de terceiros. A decisão resolveu os problemas operacionais do Keycloak, mas introduziu uma dependência de terceiro que o projeto não quer mais carregar neste estágio:

- **Dependência de third-party API para funcionalidade crítica:** todo login do sistema passa a depender da disponibilidade e dos termos de serviço do Clerk. Um outage do Clerk é um outage de autenticação do AjudaFio, sem alternativa local.
- **`CLERK_JWKS_URL` como variável obrigatória no boot:** `config.Load()` fazia `mustGetEnv("CLERK_JWKS_URL")` — o servidor não sobe sem uma URL de um serviço externo específico.
- **Dados de identidade fora do domínio da aplicação:** o backend já mantém sua própria tabela `users` (id, name, email, role) para as regras de negócio; a senha e a sessão, porém, viviam inteiramente no Clerk, criando duas fontes de verdade de identidade ligadas apenas por e-mail.
- **Custo em escala e vendor lock-in**, já sinalizados como trade-off no próprio ADR-003.

## Decisão

Substituímos o Clerk por um módulo de autenticação **self-hosted**, implementado em `internal/auth`, seguindo o padrão de `golang-jwt/jwt/v5` com assinatura HMAC (HS256) e chave simétrica própria (`JWT_SECRET`). O backend deixa de ser um validador stateless de token de terceiro e passa a ser o **emissor e validador** do próprio JWT.

Fluxo de autenticação:

```
Usuário → POST /auth/register ou /auth/login (backend)
               ↓ senha verificada com bcrypt (users.password_hash)
          Backend emite access token (15 min) + refresh token (7 dias)
               ↓ Frontend envia Bearer <access token> nas requisições
          Middleware (backend Go)
               ↓ valida assinatura HMAC com JWT_SECRET
               ↓ valida algoritmo (recusa alg != HMAC)
               ↓ extrai claims (user_id, email, name, role)
          Handler / Use Case
               ↓ (mesmo padrão admin-ou-dono de antes, inalterado)

POST /auth/refresh → valida refresh token, rotaciona (revoga o antigo, emite novo par)
POST /auth/logout  → revoga o refresh token
```

Principais mudanças técnicas:

- **`internal/auth`** (pacote raiz) ganha `jwt.go` (geração/validação HS256, com keyfunc que recusa qualquer algoritmo não-HMAC), `password.go` (bcrypt), `hash.go` (fingerprint SHA-256 do refresh token para armazenamento), e `service.go` (`Register`, `Login`, `Refresh`, `Logout`).
- **`internal/auth/adapters/postgres`** ganha queries sqlc próprias contra `users.password_hash` e a nova tabela `refresh_tokens` — sem tocar `internal/user/ports`, evitando quebrar os mocks de `UserService` já usados em 4+ módulos de teste.
- **`internal/auth/middleware/requestAuth.go`** trocou `keyfunc`/JWKS remoto por validação local com `JWT_SECRET`; a interface pública (`RequestAuth`, `GetClaims`, `WithClaims`, `IsAdmin`) não mudou.
- **`internal/auth/domain/claims.go`** ganhou o campo `UserID` (subject), preservando `Email`/`Name`/`Role` — `ValidateSameUserID` (`internal/shared/validate.go`) continua resolvendo identidade por e-mail, sem alteração.
- **`POST /auth/register`** substitui o uso público de `POST /users` como via de cadastro: cria o usuário, grava `password_hash` e já retorna o par de tokens. `POST /users` agora exige autenticação (uso administrativo).
- Corrigido, no mesmo commit, o `IsAdmin` (comparava `"admin"` minúsculo contra o enum do banco `"ADMIN"`) e uma escalação de privilégio latente em `PATCH /users/{id}` (usuário conseguia alterar o próprio `role` — inofensivo sob Clerk porque a claim vinha de fora do banco, mas explorável assim que o backend passou a emitir `role` a partir da própria tabela `users`).

## Consequências

### Positivas
- **Zero dependência de terceiro para autenticação:** o sistema funciona totalmente offline de qualquer SaaS de identidade; outage de auth só pode vir do próprio backend/banco.
- **Fonte única de verdade de identidade:** `users` (com `password_hash`) é dona completa do dado de credencial, eliminando a duplicidade Clerk↔DB que existia via e-mail.
- **Revogação real de sessão:** a tabela `refresh_tokens` com rotação e `revoked_at` permite logout efetivo — algo que o Clerk resolvia como caixa-preta.
- **RBAC e regras de admin-ou-dono preservados sem redesenho:** `IsAdmin`/`ValidateSameUserID` mantêm assinatura e call sites; o único ajuste foi corrigir o case do valor do claim `role`.

### Negativas / Trade-offs
- **Responsabilidade de segurança volta para o backend:** hashing de senha, geração/rotação de token e proteção contra `alg: none` agora são código nosso, não mais delegado a um provider especializado.
- **Sem MFA, login social, magic link ou passkeys** prontos — recursos que o Clerk oferecia via dashboard agora exigiriam implementação própria caso se tornem necessários.
- **`JWT_SECRET` é um segredo único e simétrico:** seu vazamento compromete todos os tokens emitidos; rotação de secret invalida todas as sessões ativas (não há chave assimétrica/JWKS para rotação gradual).

### Neutras / Notas de implementação
- `JWT_SECRET` deve ter no mínimo 32 bytes aleatórios (`openssl rand -base64 32`); o middleware recusa secrets menores no boot.
- `password_hash` é nullable na tabela `users` — linhas criadas na era Clerk não têm senha e precisarão de um fluxo de "definir senha" fora do escopo desta migração.
- Qualquer frontend que ainda use o SDK do Clerk quebra integralmente com esta troca; a migração do lado do frontend é coordenação separada, fora do alcance deste ADR.

## Alternativas consideradas

| Alternativa | Por que foi descartada |
|---|---|
| Manter Clerk e só reduzir uso | Não resolve a exigência explícita do projeto de eliminar dependência de terceiro para autenticação. |
| Auth0 / outro IdP SaaS | Mesmo trade-off de vendor lock-in que motivou esta ADR; não atende ao requisito de não depender de third-party API. |
| Zitadel self-hosted | Reintroduz custo de infraestrutura própria (container, banco, config) — complexidade descartada desde o ADR-002/003 para o estágio atual do projeto. |
| Sessão em cookie opaco + store server-side (ao invés de JWT) | Também viável e sem terceiro, mas descartada porque o artigo de referência e o padrão já em uso no projeto (`golang-jwt/jwt/v5`, já presente no `go.mod`) apontavam para JWT com o menor diff possível em relação ao middleware existente. |

## Referências

- [golang.com.br — Autenticação JWT em Go](https://golang.com.br/aprenda/golang-jwt-autenticacao/)
- [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)
- [golang.org/x/crypto/bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [[ADR-001-architecture-combined-hexagonal-slicing]] — Arquitetura Hexagonal + Vertical Slicing (contexto de adapter substituível)
- [[ADR-003-clerk-auth-migration]] — Decisão original que esta ADR supersede

---
