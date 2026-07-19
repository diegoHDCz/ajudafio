Contexto:

    Tarefa: Adaptar a tela/listagem de profissionais no React para consumir a nova resposta paginada do endpoint `GET /professionals`.

    Motivação: O backend deixou de retornar um array simples de profissionais nesse endpoint e passou a retornar um objeto de página (`items` + metadados de paginação), além de ter enriquecido cada profissional com dados do usuário vinculado (nome, avatar, e-mail, role).

    Arquitetura Atual: A listagem hoje provavelmente busca o array direto da resposta e renderiza a lista sem paginação real (ou com paginação client-side). Isso precisa ser ajustado para paginação server-side.

Mudanças no contrato da API (`GET /professionals`):

    Novos query params, todos opcionais:

        page (int, padrão 1)

        page_size (int, padrão 20, máximo 100 — valores fora do range são normalizados pelo backend)

    Os filtros existentes continuam os mesmos: city, state, day_of_week (multi), shift (multi).

    Formato da resposta (200 OK), antes era um array de profissionais e agora é um objeto:

        {
          "items": ProfessionalResponse[],
          "total": number,        // total de registros que casam com o filtro
          "page": number,         // página atual (normalizada pelo backend)
          "page_size": number,    // tamanho de página efetivo (normalizado pelo backend)
          "total_pages": number   // total de páginas calculado a partir de total/page_size
        }

    Cada item de `items` (ProfessionalResponse) ganhou 4 novos campos, todos opcionais/nuláveis (podem vir ausentes ou `null` quando o usuário vinculado não tiver o dado preenchido):

        name (string | null) — nome do usuário dono do perfil profissional

        avatar_url (string | null) — URL assinada (S3) do avatar do usuário, com expiração; não deve ser cacheada por muito tempo

        email (string | null)

        role (string | null)

Requisitos da implementação no React:

    Tipos/Interfaces:

        Atualizar o tipo do profissional para incluir `name`, `avatarUrl`, `email`, `role` como opcionais/nuláveis.

        Criar (ou atualizar) um tipo de página, ex: `ProfessionalPage { items, total, page, pageSize, totalPages }`, espelhando a resposta do backend (atenção ao mapeamento snake_case -> camelCase se for o padrão do projeto).

    Camada de dados (hook/service que chama a API):

        Atualizar a função/hook que busca profissionais para enviar `page` e `page_size` como query params e parsear o novo formato de resposta (objeto com `items` ao invés de array).

        Manter compatibilidade com os filtros já existentes (city, state, day_of_week, shift).

    UI da listagem:

        Adicionar controles de paginação (ex.: anterior/próxima página, ou indicador de página atual / total de páginas), usando `total`, `page`, `page_size` e `total_pages` retornados pela API.

        Exibir os novos dados do usuário no card/linha do profissional: nome (`name`) e avatar (`avatar_url`), com fallback visual (placeholder) para quando vierem `null`.

        Garantir que a troca de página/filtros reseta a paginação para a página 1 quando os filtros mudam.

    Estados de carregamento e borda:

        Tratar lista vazia (`items: []`, `total: 0`) com uma mensagem apropriada.

        Tratar o caso de `avatar_url` nulo sem quebrar o layout (placeholder/avatar padrão).
