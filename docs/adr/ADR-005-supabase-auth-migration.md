# ADR-005: Migração de JWT self-hosted para Supabase Auth

**Status:** aceito
**Data:** 2026-08-29
**Decisores:** @Diego Hernan Demitto Czajka
**Tags:** [arquitetura, segurança, autenticação, autorização, rbac, infra, supabase]
**Supersede:** [ADR-004 — Migração de Clerk para autenticação JWT self-hosted](ADR-004-self-hosted-jwt-auth.md)

---

## Contexto

O [ADR-004](ADR-004-self-hosted-jwt-auth.md) trouxe a responsabilidade de emitir e validar JWTs, hashear senhas (bcrypt) e gerenciar `refresh_tokens` para dentro do backend Go, eliminando a dependência do Clerk. Isso resolveu o problema de vendor lock-in de identidade, mas reintroduziu um custo que o projeto não quer mais carregar ao iniciar o desenvolvimento mobile (Android KMP + iOS nativo):

- **Sem login social, MFA ou magic link prontos** — features que exigiriam implementação própria no backend (§61 do ADR-004), justamente quando o projeto precisa de login social (Google, Apple) para o MVP mobile.
- **Toda a responsabilidade de segurança de credenciais** (hashing, rotação de refresh token, proteção contra `alg: none`) é código próprio, sem superfície de auditoria de terceiro especializado.
- **O projeto já usa Supabase** como Postgres gerenciado e, desde a migração de storage (`feat(storage): migrate from AWS S3 to Supabase Storage configuration`), também como storage — adotar Supabase Auth consolida o número de fornecedores em vez de aumentá-lo, ao contrário do que ocorreria com qualquer outro IdP externo.
- **O MVP não terá domínio próprio** — Supabase Auth suporta deep link/custom URL scheme sem exigir um domínio HTTPS público, o que Keycloak/Clerk também suportariam, mas com mais peças móveis.
- **RBAC precisa evoluir**: o modelo atual (`users.role` como enum único: `CLIENT`/`PROFESSIONAL`/`ADMIN`) não separa papel de status (um profissional pendente de aprovação já é `PROFESSIONAL`, sem forma de bloquear ações antes da aprovação) e não contempla o novo papel `FINANCIAL_SPONSOR`.

## Decisão

Adotamos o **Supabase Auth** como Identity Provider, delegando a ele cadastro, login, recuperação de senha, verificação de e-mail, login social (Google, Apple) e emissão de JWT. O backend Go deixa de ser emissor de token e volta a ser um **validador stateless de JWT de terceiro** — como era na era Clerk (ADR-003) — mas com autorização (RBAC) mais robusta e dados de perfil por tipo de usuário mantidos no Postgres da própria aplicação (que, notavelmente, é o mesmo Postgres hospedado no Supabase).

Fluxo de autenticação:

```
Usuário → Mobile (KMP / iOS nativo)
               ↓ login/registro via Supabase Auth SDK
          Supabase emite JWT (assinatura assimétrica, chaves via JWKS)
               ↓ Mobile envia Bearer <access_token> nas requisições
          Middleware (backend Go)
               ↓ busca chaves públicas em SUPABASE_JWKS_URL (cache local)
               ↓ valida assinatura, issuer, audience e expiração
               ↓ extrai `sub` (auth_user_id) e demais claims
          Handler / Use Case
               ↓ resolve o usuário de aplicação por auth_user_id (cria no primeiro acesso)
               ↓ carrega role(s) em user_roles
               ↓ RBAC decide permissão
```

Principais mudanças técnicas:

- **`internal/auth`** perde `service.go` (`Register/Login/Refresh/Logout/GoogleLogin`), `jwt.go` (emissão de token), `password.go`, `hash.go` e `google_verifier.go` — o backend não implementa mais uma segunda via de login.
- **`internal/auth/middleware/requestAuth.go`** passa a validar JWTs do Supabase via JWKS (`github.com/MicahParks/keyfunc/v3`, já usado na era Clerk) em vez de HMAC local.
- **`internal/auth/domain/claims.go`** passa a refletir o formato de claim do Supabase (`sub`, `email`, `role` = `authenticated`) em vez de claims próprias.
- **Novo pacote `internal/auth/rbac`** mapeia role → permissões (`internal/auth/rbac/permissions.go`), usado por `GET /me/permissions` e pelos handlers que precisam de checagem de permissão além de "é dono ou admin".
- **`users`** ganha `auth_user_id` (UUID do Supabase) e `onboarding_status`; `role` é renomeado para o novo enum (`FAMILY_CLIENT`, `HEALTH_CAREPROVIDER`, `FINANCIAL_SPONSOR`, `PLATFORM_ADMIN`) e passa a ser espelhado em uma nova tabela `user_roles` (papel e status continuam conceitualmente separados: `users.status` já existia como conceito de bloqueio/ativo, `provider` — hoje `internal/professional` — já carrega seu próprio `verified`/status de aprovação).
- **Reuso deliberado:** o feat doc pede uma tabela `provider_profiles` com `professional_data`/`verification_status`; o backend já tem exatamente isso em `internal/professional` (tabela `professionals`, com `category`, `verified`, `resume`, `metadata` por `user_id`). Criar uma segunda tabela duplicaria a fonte de verdade do profissional — `internal/professional` **é** o perfil de `HEALTH_CAREPROVIDER`. Apenas `family_profiles` e `financial_profiles` são novas, por não terem equivalente hoje.
- **`refresh_tokens`, `identities`, `users.password_hash`** são removidos (migração 000020) — sessão e refresh passam a ser geridos inteiramente pelo SDK do Supabase no cliente; o backend nunca vê senha nem refresh token.
- **`POST /auth/register|login|refresh|logout|google` são removidos**, não mantidos em paralelo — o mobile passa a falar diretamente com o Supabase Auth SDK. Isso é um corte, não uma migração incremental: o release do backend e o release do app mobile precisam ser coordenados.

## Consequências

### Positivas
- **Login social, MFA e recuperação de senha "de graça"**, via Supabase, sem código próprio de segurança de credencial no backend.
- **Consolidação de fornecedores**: Postgres, Storage e Auth no mesmo provedor (Supabase), reduzindo a superfície operacional total do projeto em vez de aumentá-la.
- **RBAC mais expressivo**: papel e status desacoplados, papel armazenado em tabela própria (`user_roles`), pronto para múltiplos papéis por usuário no futuro sem migração de schema adicional.
- **Sem reinvenção**: perfil de profissional de saúde reaproveita `internal/professional` em vez de duplicar dado.
- **Princípio do adapter substituível mantido**: os pontos de acoplamento ao Supabase ficam concentrados em `internal/auth/middleware` e `internal/auth/domain/claims.go`, como já era o padrão desde o ADR-003.

### Negativas / Trade-offs
- **Vendor lock-in de identidade**: dados de credencial (senha, sessão) voltam a viver fora do backend, no Supabase — mesmo trade-off já aceito no ADR-003 e revertido no ADR-004, agora reaceito porque o projeto já depende do Supabase para banco e storage.
- **Corte, não transição suave**: os endpoints de auth self-hosted somem nesta mesma entrega; o app mobile precisa migrar para o SDK do Supabase Auth no mesmo release do backend.
- **Sem visibilidade de login/logout no backend**: como o Supabase Auth cuida do login diretamente com o cliente, os eventos de auditoria `USER_LOGIN`/`USER_LOGOUT` (§31 do feat doc) exigem um webhook do Supabase Auth para o backend — não são mais observados naturalmente pelo código Go como eram no fluxo self-hosted. Esse webhook fica fora do escopo desta primeira entrega e é uma lacuna conhecida.
- **Migração de usuários existentes**: contas criadas na era self-hosted têm `password_hash` bcrypt, que não é portável para o Supabase Auth. Requer um script de importação (Supabase Admin API, `email_confirm: true` + fluxo de redefinição de senha) fora do escopo desta ADR.

### Neutras / Notas de implementação
- Este projeto escolheu **JWKS/assimétrico** em vez de segredo compartilhado HS256 para validar os JWTs do Supabase — nenhum segredo de assinatura precisa viver no backend, e o padrão já foi usado (e removido) na era Clerk, então a equipe já tem familiaridade com `keyfunc`.
- `users.id` permanece a chave primária interna, distinta de `auth_user_id` (UUID do Supabase) — todas as FKs existentes continuam apontando para `users.id`, sem migração em cascata.

## Alternativas consideradas

| Alternativa | Por que foi descartada |
|---|---|
| Manter JWT self-hosted e adicionar login social manualmente | Reintroduz o trabalho de implementar OAuth com Google/Apple do zero no backend — exatamente o custo que Supabase Auth elimina, com o mesmo risco de segurança de "código próprio" já aceito como trade-off negativo no ADR-004. |
| Segredo HS256 compartilhado (`SUPABASE_JWT_SECRET`) em vez de JWKS | Funcional e com diff menor em relação ao `HMACKeyfunc` existente, mas exige um segredo simétrico vivo no backend; descartado a favor de chaves assimétricas já que o projeto não tem hoje nenhuma decisão tomada sobre o modo de assinatura do projeto Supabase e JWKS é o padrão mais novo recomendado pelo próprio Supabase. |
| Criar uma tabela `provider_profiles` nova, como o feat doc descreve literalmente | Duplicaria `internal/professional`/`professionals`, que já cobre precisamente os mesmos dados (`category`, `verified`, `resume`, `metadata` por `user_id`). Reaproveitar evita duas fontes de verdade para o mesmo profissional. |
| Manter os endpoints self-hosted em paralelo durante uma janela de transição | Avaliado e descartado por decisão explícita do projeto — o feat doc já deixa claro que o backend não deve manter uma segunda implementação de login (§32), e a duplicação de rotas de auth por um período aumentaria a superfície de ataque sem benefício líquido dado que o app mobile ainda está em desenvolvimento inicial. |

## Referências

- [Supabase Auth — JWT signing keys](https://supabase.com/docs/guides/auth/signing-keys)
- [MicahParks/keyfunc — JWKS client for Go](https://github.com/MicahParks/keyfunc)
- [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)
- [[ADR-001-architecture-combined-hexagonal-slicing]] — Arquitetura Hexagonal + Vertical Slicing (contexto de adapter substituível)
- [[ADR-004-self-hosted-jwt-auth]] — Decisão original que esta ADR supersede
- [docs/feat/refactor-auth.md](../feat/refactor-auth.md) — Especificação funcional que motivou esta migração
- [docs/result/refactor-auth.md](../result/refactor-auth.md) — Plano de implementação derivado desta ADR

---
