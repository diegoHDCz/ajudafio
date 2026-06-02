# ADR-003: Migração de Keycloak para Clerk como Identity Provider

**Status:** aceito  
**Data:** 2026-06-02  
**Decisores:** @Diego Hernan Demitto Czajka  
**Tags:** [arquitetura, segurança, autenticação, infra, clerk]  
**Supersede:** [ADR-002 — Keycloak como Identity Provider](ADR-002-keycloak-identity-provider.md)

---

## Contexto

O [ADR-002](ADR-002-keycloak-identity-provider.md) adotou o Keycloak como Identity Provider. A decisão era tecnicamente sólida, mas ao longo do desenvolvimento os custos operacionais do Keycloak mostraram-se incompatíveis com o estágio atual do projeto:

- **Custo de infraestrutura real:** Keycloak requer JVM com ~500 MB de RAM mínimos e ~60 segundos de boot em `start-dev`. Em um ambiente de desenvolvimento local com banco, API e frontend já rodando, isso dobra o uso de recursos da máquina.
- **Complexidade de setup:** subir o Keycloak corretamente exige configurar realm, client, scopes, mappers, tema customizado e exportar/versionar o `realm-export.json` a cada mudança de configuração. Qualquer novo colaborador precisa replicar esse processo antes de rodar o projeto.
- **Docker Compose frágil:** o serviço Keycloak no `docker-compose.dev.yml` dependia de health checks com tempo de espera longo, volumes para o tema e variáveis de ambiente de inicialização — qualquer mudança de versão do Keycloak ou do tema quebrava o ambiente.
- **`start-dev` não é produção:** o modo de desenvolvimento usa banco H2 embutido e desabilita TLS. Produção exigiria `start` com PostgreSQL externo, certificados e um processo de deploy separado — custo de infra significativo para um projeto early-stage.
- **Fluxo OIDC desnecessariamente complexo:** o Authorization Code Flow com redirect funcionava, mas adicionava uma callback route, gerenciamento de `state` CSRF e troca de `code` por token — toda essa complexidade vivia no backend sem valor de negócio adicional.

O frontend já gerencia sessão de usuário; o backend precisa apenas **validar que o token recebido é legítimo**.

## Decisão

Substituímos o Keycloak pelo **Clerk** como Identity Provider, delegando completamente a ele o gerenciamento de identidade (login, registro, recuperação de senha, sessão). O backend passa a ser um **validador stateless de JWT**: recebe o Bearer token emitido pelo Clerk, valida a assinatura via JWKS e extrai as claims.

O fluxo de autenticação passa a ser:

```
Usuário → Frontend (Clerk SDK)
               ↓ login/register gerenciado pelo Clerk
          Clerk emite JWT (session token)
               ↓ Frontend envia Bearer <token> nas requisições
          Middleware (backend Go)
               ↓ busca chaves públicas em CLERK_JWKS_URL
               ↓ valida assinatura + expiração
               ↓ extrai claims (email, name, role)
          Handler / Use Case
```

A stack de integração no backend se reduziu a:

- **`github.com/MicahParks/keyfunc/v3`** — busca e cacheia as chaves públicas JWKS do Clerk
- **`github.com/golang-jwt/jwt/v5`** — parsing e validação do JWT com as chaves JWKS
- **`CLERK_JWKS_URL`** — única variável de ambiente necessária para auth (ex: `https://clerk.seu-app.dev/.well-known/jwks.json`)

O `AuthMiddleware` inicializa as chaves JWKS uma vez via `keyfunc.NewDefaultCtx` e as renova automaticamente em background. As claims relevantes (`email`, `name`, `role`) são extraídas para `domain.JWTClaims` e injetadas no `context.Context` da requisição — a interface do middleware com o resto do sistema não mudou.

Não há callback route, não há troca de código, não há gerenciamento de sessão no backend.

## Consequências

### Positivas
- **Zero infra para auth:** Clerk é SaaS — sem container, sem JVM, sem volume Docker, sem `realm-export.json` para versionar. O ambiente de desenvolvimento sobe com um único serviço a menos.
- **Setup de um collaborador novo:** clonar o repo, copiar `.env.example` com a `CLERK_JWKS_URL` e rodar — auth funciona. Antes eram necessários ~20 minutos para configurar o realm no Keycloak.
- **Backend mais simples:** o middleware é um validador stateless de ~50 linhas sem dependências OIDC. A callback route, o gerenciamento de `state` e a troca de code por token foram removidos completamente.
- **Token seguro sem esforço:** o Clerk assina os JWTs com RS256 e rotaciona as chaves; o `keyfunc` renova o JWKS em background sem reinicialização do servidor.
- **Fluxos de produto prontos:** login social (Google, GitHub), MFA, magic links e passkeys são configuração no dashboard do Clerk, não código.
- **Princípio do adapter substituível mantido:** `internal/auth/middleware/requestAuth.go` e `internal/auth/domain/claims.go` são os únicos pontos de acoplamento ao Clerk. Uma troca futura (Auth0, Zitadel) toca apenas esses dois arquivos.

### Negativas / Trade-offs
- **Vendor lock-in SaaS:** os dados de usuário (email, hash de senha, sessões) vivem nos servidores do Clerk. Migrar exige exportar usuários e importar em outro provider — processo não trivial em produção.
- **Custo em escala:** o plano gratuito do Clerk suporta 10.000 MAU; acima disso há custo por usuário ativo. Para o estágio atual do projeto isso é irrelevante, mas deve ser reavaliado ao escalar.
- **Sem controle sobre fluxos de UI de auth:** telas de login/registro são gerenciadas pelo Clerk (componentes ou hosted pages). Customização profunda exige o plano pago.
- **Debug de claims depende do Clerk:** se um campo não aparecer no JWT (ex: `role` ausente), a causa está na configuração de session claims no dashboard do Clerk — não no código Go.

### Neutras / Notas de implementação
- A `CLERK_JWKS_URL` segue o padrão `https://<frontend-api>.clerk.accounts.dev/.well-known/jwks.json` e é encontrada no dashboard do Clerk em **API Keys**.
- O campo `role` nas claims deve ser adicionado via **session token customization** no dashboard do Clerk (JWT Templates), mapeando o metadata do usuário para o claim `role`.
- O `keyfunc.NewDefaultCtx` recebe o `context.Context` do `main` e encerra o refresh de chaves quando o contexto é cancelado — o shutdown gracioso já está coberto.
- Em produção, a `CLERK_JWKS_URL` deve apontar para o ambiente de produção do Clerk (domínio próprio), diferente do ambiente de desenvolvimento.

## Alternativas consideradas

| Alternativa | Por que foi descartada |
|---|---|
| Manter Keycloak e resolver a complexidade de infra | O custo operacional é estrutural (JVM, config management, boot lento) — não resolve com ajuste de configuração. Para o estágio atual do projeto, o investimento não se justifica. |
| Auth próprio (JWT + bcrypt no Go) | Mantém todos os problemas de segurança e de cobertura de features que motivaram o ADR-002. |
| Auth0 | Funcionalidade equivalente ao Clerk, mas pricing menos favorável no free tier e SDK menos ergonômico para frontend moderno. |
| Zitadel | Open-source e moderno, mas reintroduz o problema de infra (container próprio) que motivou a saída do Keycloak. Pode ser reavaliado se controle total dos dados se tornar um requisito regulatório. |
| Supabase Auth | Acoplamento forte ao ecossistema Supabase — contrário ao princípio de adapter substituível do ADR-001. |

## Referências

- [Clerk Documentation](https://clerk.com/docs)
- [Clerk — Verifying JWTs manually (backend)](https://clerk.com/docs/backend-requests/handling/manual-jwt)
- [MicahParks/keyfunc — JWKS client for Go](https://github.com/MicahParks/keyfunc)
- [golang-jwt/jwt v5](https://github.com/golang-jwt/jwt)
- [[ADR-001-architecture-combined-hexagonal-slicing]] — Arquitetura Hexagonal + Vertical Slicing (contexto de adapter substituível)
- [[ADR-002-keycloak-identity-provider]] — Decisão original que esta ADR supersede

---
