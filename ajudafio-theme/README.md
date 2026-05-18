# Tema Ajuda Fio — Keycloak

Tema personalizado do Keycloak com a identidade visual da **Ajuda Fio**.
Inclui as páginas: **Login**, **Registro** e **Esqueci Senha**.

---

## Estrutura de Arquivos

```
ajudafio-theme/
├── docker-compose.yml
└── theme/
    └── ajudafio/
        └── login/
            ├── theme.properties
            ├── login.ftl                  # Página de login
            ├── register.ftl               # Página de registro
            ├── login-reset-password.ftl   # Esqueci minha senha
            └── resources/
                ├── css/
                │   └── login.css
                ├── js/
                │   └── login.js
                └── img/
                    └── logo.png
```

---

## Como Instalar

### 1. Subir o Keycloak com Docker

```bash
# Na pasta ajudafio-theme/
docker compose up -d
```

Aguarde o Keycloak iniciar (pode levar ~30s na primeira vez).

### 2. Acessar o Painel Admin

Abra: http://localhost:8080/admin  
Usuário: `admin`  
Senha: `admin123`

### 3. Selecionar o Tema no Realm

1. No painel admin, clique no nome do seu **Realm** (ou crie um novo)
2. Vá em **Realm Settings** → aba **Themes**
3. Em **Login Theme**, selecione `ajudafio`
4. Clique em **Save**

### 4. Testar as Páginas

- **Login:** http://localhost:8080/realms/master/protocol/openid-connect/auth?client_id=account
- **Registro:** http://localhost:8080/realms/master/protocol/openid-connect/registrations?...
- **Esqueci senha:** Clique no link "Esqueci minha senha" na tela de login

---

## Se Já Tem o Keycloak Rodando

Se o seu Keycloak já está em execução via Docker (sem docker-compose), copie o tema:

```bash
# Descobrir o nome/ID do container
docker ps

# Copiar o tema para dentro do container
docker cp ./theme/ajudafio SEU_CONTAINER:/opt/keycloak/themes/ajudafio
```

Depois ative o tema no painel admin conforme o passo 3 acima.

---

## Personalizar o Tema

| Arquivo | O que editar |
|---------|-------------|
| `resources/css/login.css` | Cores, fontes, espaçamento |
| `resources/img/logo.png` | Substituir pelo logo em outra versão |
| `login.ftl` | Textos e estrutura da tela de login |
| `register.ftl` | Campos e textos do cadastro |
| `login-reset-password.ftl` | Tela de recuperação de senha |

### Trocar as Cores Principais

No topo do `login.css`, edite as variáveis CSS:

```css
:root {
  --af-blue: #1565C0;   /* Azul principal (botões, labels) */
  --af-green: #2E7D32;  /* Verde (destaques) */
}
```

---

## Versões Suportadas

Testado com Keycloak **23.x** e **24.x**.
Para versões anteriores (< 20), substituir `messagesPerField` pelos equivalentes da versão.
