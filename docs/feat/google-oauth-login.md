Contexto:

    Tarefa: Adicionar "Login/Sign up with Google" (OAuth 2.0 / OpenID Connect) na arquitetura de autenticação atual, sem quebrar o fluxo existente de email+senha.

    Motivação: Reduzir fricção no cadastro/login (menos formulário, menos senha pra lembrar). O próprio ADR-004 (que removeu o Clerk e trouxe o auth self-hosted) já registrou essa lacuna: "Sem MFA, login social, magic link ou passkeys prontos — recursos que o Clerk oferecia via dashboard agora exigiriam implementação própria caso se tornem necessários." Este documento é essa implementação própria.

    Arquitetura Atual (`internal/auth`, ver ADR-004): auth 100% self-hosted, sem Clerk/Keycloak.
        - JWT HS256 (`golang-jwt/v5`) com `JWT_SECRET` único (`internal/auth/jwt.go`), access token 15min, refresh token 7 dias.
        - Senha com bcrypt (`internal/auth/password.go`).
        - `internal/auth/service.go`: `Register`, `Login`, `Refresh`, `Logout`, todos emitindo par de tokens via `issueTokenPair`.
        - Refresh tokens persistidos (hash SHA-256) na tabela `refresh_tokens`, rotacionados a cada refresh.
        - Rotas públicas em `POST /auth/{register,login,refresh,logout}` (`internal/auth/adapters/http/handler.go`, montadas em `cmd/main.go`).
        - Tabela `users`: `id, name, email UNIQUE, phone, role, avatar_url, password_hash (nullable), created_at, updated_at`. Sem nenhuma coluna de OAuth hoje.
        - `internal/user` (módulo separado) expõe `UserService.Create/GetByEmail/GetByID`, que o `auth.service` já usa para achar/criar usuário — o fluxo Google vai reaproveitar essas mesmas chamadas.

Decisão de fluxo OAuth:

    Existem duas formas usuais de implementar "Sign in with Google" numa API separada do frontend (SPA). Escolher uma:

    Opção A — Verificação de ID Token no backend (recomendada):
        - Frontend usa o Google Identity Services (GIS, `accounts.google.com/gsi/client`) para autenticar o usuário direto no navegador e obter um `id_token` (JWT assinado pelo Google).
        - Frontend manda esse `id_token` pro backend: `POST /auth/google { "id_token": "..." }`.
        - Backend valida a assinatura/issuer/audience/expiração do token (biblioteca oficial `google.golang.org/api/idtoken`, que já cuida do JWKS do Google e do cache de chaves), extrai `sub` (ID único do usuário no Google), `email`, `email_verified`, `name`, `picture`.
        - Backend acha ou cria o usuário local e emite o PAR DE TOKENS PRÓPRIO de sempre (mesmo `issueTokenPair` do fluxo de senha). O frontend nunca precisa saber que o backend fala com o Google — ele só troca `id_token` do Google por `access_token`/`refresh_token` do ajudafio, igual já faz em `/auth/login`.
        - Prós: backend fica stateless (sem sessão de redirect, sem `state`/`nonce` pra guardar), reaproveita 100% da infra de resposta/erro que já existe em `/auth/login`, é o padrão mais comum quando já existe um frontend SPA consumindo a API via JSON.
        - Contras: depende do frontend carregar o script do GIS.

    Opção B — Authorization Code Flow tradicional no backend:
        - Backend expõe `GET /auth/google/redirect` (redireciona pro Google) e `GET /auth/google/callback` (troca `code` por token, usando `golang.org/x/oauth2/google`).
        - Precisa de `state`/PKCE, cookie ou storage temporário de sessão durante o redirect, e um passo extra de redirecionar de volta pro frontend com os tokens (geralmente via cookie httpOnly ou query param + página intermediária).
        - Só compensa se não houver frontend JS controlando o fluxo (ex.: app mobile nativo com WebView, ou backend servindo HTML diretamente).

    Recomendação: Opção A. É bem mais simples de implementar e testar, evita gerenciar estado de redirect no backend, e se encaixa no formato atual da API (JSON request/response, sem sessão/cookie).

Mudanças no banco de dados:

    Nova migration (ex.: `migrations/000013_add_identities.up.sql`), criando uma tabela separada em vez de colunas soltas em `users` — assim já fica pronto pra outros provedores no futuro (Apple, Facebook, etc.) e pra um mesmo usuário linkar múltiplos provedores:

        CREATE TABLE identities (
          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
          user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
          provider TEXT NOT NULL,              -- 'google'
          provider_user_id TEXT NOT NULL,      -- 'sub' do Google
          email TEXT NOT NULL,                 -- email retornado pelo Google no momento do login
          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
          UNIQUE (provider, provider_user_id)
        );
        CREATE INDEX idx_identities_user_id ON identities(user_id);

    `password_hash` em `users` já é nullable (migration 000011) — usuários que só usam Google não precisam de senha, o que já funciona sem mudança extra.

    Ponto de atenção / decisão de produto a confirmar: se já existe um usuário cadastrado por email+senha e ele faz login com Google usando o mesmo email, o comportamento deve ser linkar automaticamente a identidade Google a esse usuário (recomendado, desde que `email_verified=true` no token do Google) ou bloquear/pedir confirmação? Este documento assume "linkar automaticamente por email verificado", mas vale validar com o time.

Mudanças no domínio/portas (`internal/auth`):

    `internal/auth/domain/`: novo arquivo (ex. `identity.go`) com a struct de identidade OAuth e/ou os claims relevantes do Google (não precisa reimplementar validação de assinatura — isso fica a cargo da lib `idtoken`).

    `internal/auth/ports/repository.go`: novos métodos na interface `AuthRepository`:
        - `GetIdentityByProvider(ctx, provider, providerUserID string) (*domain.Identity, error)`
        - `CreateIdentity(ctx, userID uuid.UUID, provider, providerUserID, email string) error`

    `internal/auth/ports/service.go`: novo método na interface `AuthService`:
        - `GoogleLogin(ctx context.Context, idToken string) (*domain.TokenPair, error)`

    `internal/auth/service.go`: implementar `GoogleLogin`:
        1. Validar o `id_token` via `idtoken.Validate(ctx, idToken, googleClientID)` — rejeita se `aud` não bater com `GOOGLE_CLIENT_ID` configurado, ou se assinatura/expiração forem inválidas.
        2. Exigir `email_verified == true` nos claims (senão retornar erro — não confiar em email não verificado).
        3. `repo.GetIdentityByProvider(ctx, "google", sub)`:
            - Se existir → já sabe o `user_id`, pula direto pra emitir tokens.
            - Se não existir → `userSvc.GetByEmail(ctx, email)`:
                - Se usuário já existe (cadastro por senha) → cria a `identity` linkando esse `user_id` (ver decisão de produto acima).
                - Se não existe → cria usuário novo via `userSvc.Create` (mesmo padrão do `Register`: `Role: RoleClient` por padrão, `AvatarURL` pode ser populado com a `picture` do Google) e cria a `identity`.
        4. Emitir tokens com o `issueTokenPair` já existente — nenhuma mudança na emissão/formato do JWT próprio.

    `internal/auth/adapters/postgres/`: nova query sqlc (`queries/auth.sql`) para `GetIdentityByProvider` e `CreateIdentity`, regenerar código sqlc.

Config (`internal/infra/config/config.go`):

    Novos campos, seguindo o padrão atual (`mustGetEnv`/`getEnv`):
        - `GoogleClientID` (`mustGetEnv("GOOGLE_CLIENT_ID")` se a feature for obrigatória, ou `getEnv` com vazio se quiser feature-flag opcional)
    Atualizar `.env.example` com `GOOGLE_CLIENT_ID=...` e um comentário de onde obter (Google Cloud Console → Credentials → OAuth Client ID, tipo "Web application").
    Não precisa de `GOOGLE_CLIENT_SECRET` na Opção A (verificação de ID token é pública, não troca segredo com o Google) — só seria necessário na Opção B.

Rotas HTTP (`internal/auth/adapters/http/handler.go`):

    Novo handler `GoogleLogin`, nova rota pública (sem `authMW`, igual `/auth/login`):
        `POST /auth/google` — body `{ "id_token": "string" }` → resposta idêntica à de `/auth/login` (par de tokens).
    Adicionar ao `NewRouter` existente, ao lado de `register`/`login`/`refresh`/`logout`.
    Atualizar comentários Swagger (`docs.go`/anotações) e o arquivo `rest/auth.rest` (hoje desatualizado, ainda referencia Keycloak — bom momento de limpar).

Dependências novas (`go.mod`):

    `google.golang.org/api` (pacote `idtoken`) — validação de ID token do Google com verificação de assinatura/JWKS/cache embutidas. Evita reimplementar parsing de JWKS na mão.
    Não precisa de `golang.org/x/oauth2` na Opção A (só seria necessário na Opção B, para o Authorization Code Flow).

Segurança:

    Sempre checar `email_verified` nos claims do Google antes de aceitar o email como confiável.
    `aud` do token deve bater exatamente com o `GOOGLE_CLIENT_ID` configurado — a lib `idtoken.Validate` já faz isso, mas vale um teste explícito pra não deixar passar `aud` vazio/errado por engano de config.
    Rate limiting em `/auth/google` (se já não houver algo global) pelo mesmo motivo de `/auth/login`.
    CORS: já precisa liberar o domínio do frontend para as rotas de auth existentes; nada novo aqui além de garantir que o script do GIS carrega no domínio certo (config no Google Cloud Console, não no backend).

Testes:

    Unit test de `service.GoogleLogin` com um `idtoken.Validate` abstraído atrás de uma interface pequena (ex. `GoogleTokenVerifier`) pra poder mockar nos testes — evita bater na rede real do Google nos testes.
    Casos a cobrir: usuário novo (cria user + identity), usuário existente com identity já linkada (login direto), usuário existente por senha fazendo primeiro login com Google (linka identity), `email_verified=false` (rejeita), token inválido/expirado (rejeita).
    Handler test padrão (like os demais handlers do módulo) para o parsing do body e mapeamento de erros pra status HTTP.

Passos de implementação (ordem sugerida):

    1. Criar app OAuth no Google Cloud Console, obter `GOOGLE_CLIENT_ID`, configurar origens JS autorizadas com o domínio do frontend.
    2. Adicionar `google.golang.org/api/idtoken` ao `go.mod`.
    3. Migration `000013_add_identities.up.sql` (+ `.down.sql`) e rodar.
    4. Query sqlc + regenerar (`GetIdentityByProvider`, `CreateIdentity`).
    5. `ports`: adicionar métodos em `AuthRepository`/`AuthService`.
    6. `domain`: struct de identidade (se necessário).
    7. `service.go`: implementar `GoogleLogin` (com `GoogleTokenVerifier` mockável).
    8. `config.go` + `.env.example`: `GOOGLE_CLIENT_ID`.
    9. Handler + rota `POST /auth/google`, wiring em `cmd/main.go`.
    10. Testes unitários (service + handler).
    11. Atualizar Swagger e `rest/auth.rest`.
    12. (Frontend, fora deste doc) integrar Google Identity Services e chamar `POST /auth/google`.

Perguntas em aberto (confirmar com o time antes de implementar):

    - Linkar automaticamente por email verificado ou exigir confirmação explícita do usuário?
    - Popular `avatar_url` do usuário com a foto do Google no primeiro login? Sobrescrever em logins seguintes ou só na criação?
    - Login com Google deve permitir escolher role (ex. virar PROFESSIONAL depois) ou sempre entra como CLIENT, igual ao `Register` atual?
    - Vale a pena um endpoint `POST /auth/logout-all` (usando `RevokeAllRefreshTokensForUser`, que já existe na porta mas está sem uso) pro caso de usuário linkar/desvincular Google e querer invalidar sessões antigas?
