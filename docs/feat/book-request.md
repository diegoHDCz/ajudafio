Contexto:

    Tarefa: Modelagem de Banco de Dados (Drizzle ORM / SQL) especificamente para a entidade de solicitações/propostas.

    Padrão: Seguir estritamente o modelo de código, nomenclatura (camelCase/snake_case) e estrutura que hoje se encontra na pasta migration.

    Arquitetura Atual: Já temos os domínios professional, user/client, auth, contract e availability.

    Nova Entidade a Adicionar: booking_request (ou booking_proposal).

    Restrição Crítica de Migration: Nunca gerar instruções de DROP TABLE. Se houver alterações em tabelas existentes, use sempre ALTER TABLE ou adicione novas colunas com segurança.

Requisitos do Domínio booking_request:

    Objetivo do Domínio: Registrar a proposta inicial de agendamento feita pelo Cliente para o Profissional. Esta tabela guarda os dados de datas, horas, endereço e valores antes do aceite do profissional e antes da criação do contrato formal. Ela é o coração do funil de conversão e métricas do sistema.

    Relacionamentos Dinâmicos:

        Deve possuir client_id (vinculado à tabela de usuários/clientes).

        Deve possuir professional_id (vinculado à tabela de profissionais).

        Nota de Regra de Negócio: Ela não possui contract_id no nascimento, pois o contrato só passa a existir no banco de dados se esta proposta for aceita.

    Dados da Proposta Financeira e Local:

        proposed_value (decimal/numeric): O valor total ou valor/hora proposto para a prestação desse serviço.

        address (text ou campos separados de endereço): O local exato onde o cliente solicita o atendimento.

    Dados do Cronograma Solicitado (schedule_details):

        Um campo do tipo jsonb para armazenar a intenção de agenda do cliente. Esse JSON deve ser capaz de estruturar múltiplos dias e horários (ex: Segundas e Quartas das 13:00 às 17:00, iniciando na data X e terminando na data Y). Isso evita criar linhas prematuras na tabela de agendamentos.

    Controle de Estados e Métricas (O Core da Tabela):

        status: Um enum contendo os estados PENDING (Aguardando profissional), ACCEPTED (Aceito, gera contrato), REJECTED (Negado pelo profissional), EXPIRED (Tempo limite excedido) ou CANCELLED (Cliente desistiu antes do aceite).

        rejection_reason (text/nullable): Campo que se torna obrigatório caso o status mude para REJECTED, registrando a justificativa/motivo do profissional ter negado o pedido.

    Auditoria para SLAs e Métricas de Resposta:

        created_at (timestamp): Data e hora que o cliente enviou o pedido.

        responded_at (timestamp/nullable): Data e hora exata em que o profissional aceitou ou rejeitou a proposta. (Crucial para calcular o tempo médio de resposta dos profissionais).