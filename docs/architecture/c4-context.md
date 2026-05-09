# C4 – Diagrama de contexto

```mermaid
C4Context
    title Health Contracts Platform – Contexto do sistema

    Person(clinicAdmin, "Admin da Clínica", "Gestão de contratos e profissionais")
    Person(professional, "Profissional de Saúde", "Visualiza e assina contratos")
    Person(auditor, "Auditor / Compliance", "Revisão de trilhas de auditoria")

    System(platform, "Health Contracts Platform", "Gestão de contratos, profissionais, auditoria e notificações")

    System_Ext(emailProvider, "Sendgrid / SES", "Envio de e-mails e notificações")
    System_Ext(smsProvider, "Twilio", "Notificações SMS para profissionais")
    System_Ext(idProvider, "Identity Provider (Keycloak/Auth0)", "Autenticação e controle de acesso")
    System_Ext(storageProvider, "S3 / GCS", "Armazenamento de documentos e contratos PDF")

    Rel(clinicAdmin, platform, "Gerencia contratos e cadastros", "HTTPS")
    Rel(professional, platform, "Visualiza e assina contratos", "HTTPS")
    Rel(auditor, platform, "Consulta trilha de auditoria", "HTTPS")
    Rel(platform, emailProvider, "Envia notificações", "HTTPS/SMTP")
    Rel(platform, smsProvider, "Envia SMS", "HTTPS")
    Rel(platform, idProvider, "Valida tokens JWT", "HTTPS")
    Rel(platform, storageProvider, "Armazena documentos", "HTTPS")
```

---

## C4 – Diagrama de container

```mermaid
C4Container
    title Health Contracts Platform – Containers

    Person(user, "Usuário", "Admin, Profissional ou Auditor")

    Container(api, "API Server", "Go / Chi", "HTTP REST API – cmd/server")
    Container(worker, "Outbox Worker", "Go", "Processa events da outbox e publica no broker – cmd/worker")
    ContainerDb(postgres, "PostgreSQL", "PostgreSQL 16", "Dados transacionais + tabela outbox_messages")
    Container(broker, "Message Broker", "RabbitMQ / SQS", "Entrega events entre slices de forma assíncrona")

    Rel(user, api, "Usa", "HTTPS")
    Rel(api, postgres, "Lê e escreve", "pgx")
    Rel(worker, postgres, "Lê outbox pendentes", "pgx")
    Rel(worker, broker, "Publica domain events", "AMQP / AWS SDK")
    Rel(api, broker, "Consome eventos (opcional)", "AMQP / AWS SDK")
```

---

## Legenda de cores por slice (nos diagramas de componente)

| Cor | Slice |
|---|---|
| Azul | `contracts/` |
| Verde | `professionals/` |
| Laranja | `audit/` |
| Roxo | `notifications/` |
| Cinza | `shared/` e `infra/` |
