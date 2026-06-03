Contexto:

- Tarefa de modelagem de DB
- seguir o modelo que hoje encontra na pasta migration.
- Hoje temos as domains, contract,availability,professional,auth,user(client),address e gostaria de adicionar appointment
- Nunca deve fazer DROP TABLE em migrations sempre ALTER se necessário
  Requisitos:
- Ao realizar o appointment deve considerar o availability (quer diz os dias e horarios e turnos que o profissional abriu para atendimento)
- O appointment deve ter o endereço de atendimento
- O appointment tem relacionamento com o contrato que será atendido
- O appointment deve usar horas de inicio e de fim e não o turno.
  Requisitos do dominio
  Relacionamento com contract: Cada appointment deve possuir obrigatoriamente um contract_id (Relacionamento 1:N - Um contrato pode ter múltiplos agendamentos). O agendamento representa a execução física/temporal do que foi acordado no contrato.

Vínculo com professional e client: Deve conter chaves estrangeiras para identificar o profissional prestador e o cliente (podendo herdar ou validar através do contrato vinculado).

Definição de Tempo Precisa: O agendamento não usa turnos. Ele deve registrar a data exata (date) e os horários específicos de início (start_time) e fim (end_time).

Localização: O agendamento deve conter campos para o endereço completo de atendimento onde o serviço será prestado.

Ciclo de Vida (Status): A tabela deve conter um campo status (ex: PENDING, CONFIRMED, CANCELLED, COMPLETED) para controlar o fluxo do agendamento.

Regra de Negócio com availability (Validação): A IA deve entender que o appointment é uma instância concreta de tempo que precisa respeitar as regras abstratas da tabela availability (ex: se o agendamento for para uma segunda-feira às 14h, a tabela de disponibilidade do profissional deve prever atendimento nas segundas-feiras à tarde).

Evitar Concorrência (Double Booking): Incluir um mecanismo de controle (como um campo version para lock otimista ou índice único composto, se aplicável) para garantir que dois agendamentos não ocupem o mesmo slot de tempo para o mesmo profissional.
