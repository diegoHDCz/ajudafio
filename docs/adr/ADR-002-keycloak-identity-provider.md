# ADR-002: Keycloak como Identity Provider para Autenticação

**Status:** supersedido por [ADR-003 — Migração para Clerk](ADR-003-clerk-auth-migration.md)  
**Data:** 2026-05-17  
**Decisores:** @Diego Hernan Demitto Czajka  
**Tags:** [arquitetura, segurança, autenticação, infra, keycloak]

---

## Contexto

O projeto **ajudafio** precisa de autenticação de usuários contemplando os fluxos de:

- **Login** com e-mail e senha
- **Registro** de novos usuários
- **Recuperação de senha** (forgot password / reset)
- **Proteção de rotas** na API com verificação de identidade e papéis (RBAC)

As pressões que moldaram a decisão foram:

- **Time pequeno** — implementar auth do zero (hashing de senhas, tokens, refresh, MFA, e-mail de reset) seria semanas de trabalho de baixo valor de negócio.
- **Segurança por padrão** — qualquer implementação caseira de auth concentra risco: basta um bug de bcrypt mal configurado ou um token sem expiração para comprometer toda a base.
- **Extensibilidade futura** — o produto pode precisar de login social (Google, Facebook), MFA ou SSO corporativo; construir uma segunda vez do zero não é uma opção.
- **Go idiomático** — a linguagem não possui um framework de auth consolidado como o Spring Security ou o Devise; a escolha de uma solução externa é natural no ecossistema Go.
- **Alinhamento com ADR-001** — a arquitetura Hexagonal exige que a infraestrutura de autenticação seja um adapter substituível, não lógica acoplada ao domínio.

## Decisão

Usaremos **Keycloak 26.x como Identity Provider**, delegando ao Keycloak os fluxos de login, registro e recuperação de senha via **OAuth2 Authorization Code Flow + OIDC**. O backend Go valida tokens JWT emitidos pelo Keycloak via middleware OIDC, sem armazenar credenciais ou gerenciar sessões.

O fluxo de autenticação é:

```
Usuário → Frontend → Keycloak (login/register/forgot-password)
                         ↓ Authorization Code
                  Backend /auth/callback
                         ↓ troca code → id_token + access_token
                  Middleware valida JWT via OIDC Verifier
                         ↓ claims (sub, email, roles)
                  Handler / Use Case
```

A stack de integração no backend:

- **`github.com/coreos/go-oidc`** — discovery OIDC e verificação de ID Token
- **`golang.org/x/oauth2`** — Authorization Code Flow (troca de code por token)
- **Realm:** `ajudafio` (isolado, sem interferir no realm `master`)
- **Client:** `app-ajudafio` com `confidential access type`
- **Tema customizado:** `ajudafio-theme` montado no container para manter identidade visual consistente

O middleware `RequestAuth` intercepta todas as rotas protegidas, valida o Bearer JWT contra o OIDC verifier e injeta as claims (`sub`, `email`, `roles`) no `context.Context` da requisição. Helpers `HasRealmRole` e `HasClientRole` permitem RBAC granular nos handlers sem duplicar lógica.

## Consequências

### Positivas
- **Zero código de auth sensível no domínio** — hashing de senha, geração de token, e-mail de reset, bloqueio por tentativas, MFA: tudo é responsabilidade do Keycloak.
- **Fluxos prontos e testados** — login, registro com verificação de e-mail, forgot/reset password funcionam out-of-the-box; apenas o tema visual precisou de customização.
- **RBAC nativo** — papéis de realm e de client são propagados no JWT; novos papéis não exigem migração de banco no backend.
- **Extensível sem reescrever auth** — adicionar login social, MFA ou SSO no futuro é configuração no Keycloak, não código Go.
- **Adapter substituível** — `internal/auth/adapters/keycloak/repository.go` é o único ponto de acoplamento ao Keycloak; se a decisão mudar (Auth0, Zitadel), apenas esse adapter e o middleware precisam ser trocados.
- **Tema customizado desacoplado** — o diretório `ajudafio-theme/` monta o tema via volume Docker sem modificar a imagem oficial; atualizações do Keycloak não requerem rebuild da imagem.

### Negativas / Trade-offs
- **Dependência de infra adicional** — Keycloak requer JVM, consome ~500 MB de RAM e tem boot lento (~60 s em `start-dev`); o ambiente de desenvolvimento fica mais pesado.
- **Curva de aprendizado** — conceitos como realm, client, scopes, mappers e flows do Keycloak têm terminologia própria e documentação extensa.
- **Config não versionada por padrão** — as configurações do realm vivem no banco interno do Keycloak; é preciso disciplina para exportar e versionar o `realm-export.json` a cada mudança.
- **Debug de tokens mais opaco** — erros de claim ou role incorreta aparecem como `403 Forbidden` genérico; exige inspecionar o JWT manualmente durante desenvolvimento.
- **`start-dev` não é adequado para produção** — o modo de desenvolvimento desabilita TLS e usa banco H2 embutido; produção requer `start` com banco externo e certificados configurados.

### Neutras / Notas de implementação
- O `clientSecret` **não deve** ser hardcoded no source; mover para variável de ambiente (`KEYCLOAK_CLIENT_SECRET`) antes de qualquer deploy.
- O parâmetro `state` do OAuth2 deve ser um valor aleatório por sessão para prevenir CSRF; o valor fixo atual (`"exemplo"`) é adequado apenas para desenvolvimento local.
- O endpoint do OIDC discovery (`/realms/ajudafio`) é resolvido automaticamente pelo `go-oidc`; não é necessário hardcodar URLs de token/jwks.
- Em produção, o Keycloak deve ser provisionado com banco PostgreSQL externo (pode ser o mesmo `ajudafio_postgres` em schemas separados ou uma instância dedicada).
- O tema customizado fica em `ajudafio-theme/theme/ajudafio/` e é montado via volume no `docker-compose`; cada tela (login, register, email) pode ser sobrescrita individualmente com templates FreeMarker.
- O `realm-export.json` deve ser regenerado via Admin Console após cada mudança de configuração e commitado em `config/keycloak/`.

## Alternativas consideradas

| Alternativa | Por que foi descartada |
|---|---|
| Auth próprio (JWT + bcrypt no Go) | Semanas de implementação para cobrir o mesmo conjunto de features; alto risco de vulnerabilidades em código de segurança escrito do zero e não auditado. |
| Auth0 / Clerk / Okta (SaaS) | Dependência de vendor externo e custo por MAU (Monthly Active User); inaceitável para um projeto que pode ter volumes variáveis e precisa de controle total sobre os dados de usuário. |
| Zitadel | Alternativa open-source moderna e com boa API; descartada por ter comunidade menor e menos exemplos Go disponíveis no momento da decisão — pode ser reavaliada no futuro. |
| Firebase Authentication | Ecossistema Google; lock-in forte, sem hospedagem on-premise e sem suporte nativo a fluxos customizados de UI sem SDK proprietário. |
| Supabase Auth | Fortemente acoplado ao ecossistema Supabase/Postgres; adicionar esse coupling seria contrário ao princípio de adapter substituível do ADR-001. |
| Passport.js / Next-Auth | Soluções JavaScript; incompatíveis com o backend Go — consideradas apenas se houvesse um BFF em Node. |

## Referências

- [Keycloak Documentation](https://www.keycloak.org/documentation)
- [go-oidc — CoreOS OIDC client for Go](https://github.com/coreos/go-oidc)
- [golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)
- [OAuth 2.0 Authorization Code Flow — RFC 6749](https://datatracker.ietf.org/doc/html/rfc6749#section-4.1)
- [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)
- [Keycloak Themes — Customization Guide](https://www.keycloak.org/docs/latest/server_development/#_themes)
- [[ADR-001-architecture-combined-hexagonal-slicing]] — Arquitetura Hexagonal + Vertical Slicing (contexto de adapter substituível)

---
