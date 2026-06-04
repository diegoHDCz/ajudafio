-- =============================================================================
-- Seed data for integration testing
-- IDs are deterministic UUIDs for consistent cross-table references
-- Run: psql $DATABASE_URL -f seed.sql
-- =============================================================================

-- Users: 1 admin, 3 clients, 4 professionals
INSERT INTO users (id, name, email, phone, role, created_at, updated_at) VALUES
  ('00000000-0000-0000-0000-000000000001', 'Admin Master',       'admin@ajudafio.com',       '11999990001', 'ADMIN',        NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000002', 'João Silva',         'joao.silva@email.com',     '11999990002', 'CLIENT',       NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000003', 'Maria Oliveira',     'maria.oliveira@email.com', '11999990003', 'CLIENT',       NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000004', 'Carlos Souza',       'carlos.souza@email.com',   '11999990004', 'CLIENT',       NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000005', 'Ana Paula Ferreira', 'ana.ferreira@email.com',   '11999990005', 'PROFESSIONAL', NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000006', 'Roberto Mendes',     'roberto.mendes@email.com', '11999990006', 'PROFESSIONAL', NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000007', 'Fernanda Lima',      'fernanda.lima@email.com',  '11999990007', 'PROFESSIONAL', NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000008', 'Lucas Alves',        'lucas.alves@email.com',    '11999990008', 'PROFESSIONAL', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Professionals (one per professional user)
INSERT INTO professionals (id, user_id, license_number, category, years_of_experience, verified, resume, metadata, created_at, updated_at) VALUES
  (
    '00000000-0000-0000-0001-000000000001',
    '00000000-0000-0000-0000-000000000005',
    'COREN-SP 123456', 'NURSE', 8, true,
    'Enfermeira com 8 anos de experiência em UTI e home care. Especialização em cuidados paliativos.',
    '{"languages": ["Português", "Inglês"], "specialties": ["UTI", "Home Care", "Cuidados Paliativos"]}',
    NOW(), NOW()
  ),
  (
    '00000000-0000-0000-0001-000000000002',
    '00000000-0000-0000-0000-000000000006',
    NULL, 'ELDERLY_CAREGIVER', 5, true,
    'Cuidador de idosos com formação em gerontologia. Experiência com pacientes de Alzheimer e Parkinson.',
    '{"languages": ["Português"], "specialties": ["Alzheimer", "Parkinson", "Geriatria"]}',
    NOW(), NOW()
  ),
  (
    '00000000-0000-0000-0001-000000000003',
    '00000000-0000-0000-0000-000000000007',
    'CREFITO-SP 78910', 'PHYSIOTHERAPIST', 10, false,
    'Fisioterapeuta especialista em reabilitação neurológica e motora.',
    '{"languages": ["Português", "Espanhol"], "specialties": ["Reabilitação Neurológica", "Fisioterapia Motora"]}',
    NOW(), NOW()
  ),
  (
    '00000000-0000-0000-0001-000000000004',
    '00000000-0000-0000-0000-000000000008',
    NULL, 'HOSPITAL_COMPANION', 3, true,
    'Acompanhante hospitalar com experiência em oncologia e pós-operatório.',
    '{"languages": ["Português"], "specialties": ["Oncologia", "Pós-operatório"]}',
    NOW(), NOW()
  )
ON CONFLICT (id) DO NOTHING;

-- Availabilities (arrays por migration 000003)
INSERT INTO availabilities (id, professional_id, day_of_week, shift, start_hour, end_hour) VALUES
  -- Ana Paula (Nurse): seg-sex manhã+tarde
  ('00000000-0000-0000-0002-000000000001', '00000000-0000-0000-0001-000000000001',
   ARRAY['MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY']::VARCHAR(20)[],
   ARRAY['MORNING','AFTERNOON']::VARCHAR(20)[], '07:00', '19:00'),
  -- Ana Paula: sábado manhã
  ('00000000-0000-0000-0002-000000000002', '00000000-0000-0000-0001-000000000001',
   ARRAY['SATURDAY']::VARCHAR(20)[],
   ARRAY['MORNING']::VARCHAR(20)[], '08:00', '12:00'),
  -- Roberto (Elderly Caregiver): todos os dias, dia inteiro
  ('00000000-0000-0000-0002-000000000003', '00000000-0000-0000-0001-000000000002',
   ARRAY['MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY','SUNDAY']::VARCHAR(20)[],
   ARRAY['FULL_DAY']::VARCHAR(20)[], '07:00', '19:00'),
  -- Fernanda (Physiotherapist): ter, qui, sáb manhã+tarde
  ('00000000-0000-0000-0002-000000000004', '00000000-0000-0000-0001-000000000003',
   ARRAY['TUESDAY','THURSDAY','SATURDAY']::VARCHAR(20)[],
   ARRAY['MORNING','AFTERNOON']::VARCHAR(20)[], '08:00', '18:00'),
  -- Lucas (Hospital Companion): sáb+dom noturno
  ('00000000-0000-0000-0002-000000000005', '00000000-0000-0000-0001-000000000004',
   ARRAY['SATURDAY','SUNDAY']::VARCHAR(20)[],
   ARRAY['NIGHT']::VARCHAR(20)[], '19:00', '07:00'),
  -- Lucas: seg, qua, sex tarde
  ('00000000-0000-0000-0002-000000000006', '00000000-0000-0000-0001-000000000004',
   ARRAY['MONDAY','WEDNESDAY','FRIDAY']::VARCHAR(20)[],
   ARRAY['AFTERNOON']::VARCHAR(20)[], '13:00', '19:00')
ON CONFLICT (id) DO NOTHING;

-- Contracts (valores em centavos de real)
-- hour_rate: R$/hora, total_amount: total do contrato
INSERT INTO contracts (id, client_id, professional_id, status, hour_rate, total_amount, details, week_days, shift, start_time, hours_per_day, total_hours, created_at) VALUES
  -- João contrata Ana Paula (Nurse) - ACTIVE
  (
    '00000000-0000-0000-0003-000000000001',
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0001-000000000001',
    'ACTIVE', 8000, 320000,
    '{"description": "Cuidados de enfermagem home care para paciente pós-cirúrgico", "observations": "Paciente com restrição de mobilidade"}',
    ARRAY['MONDAY','WEDNESDAY','FRIDAY']::VARCHAR(20)[],
    'MORNING', '08:00:00', 8, 40, NOW() - INTERVAL '30 days'
  ),
  -- Maria contrata Roberto (Elderly Caregiver) - COMPLETED
  (
    '00000000-0000-0000-0003-000000000002',
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0001-000000000002',
    'COMPLETED', 5000, 150000,
    '{"description": "Cuidado integral para idosa com Alzheimer", "observations": "Paciente necessita de supervisão constante"}',
    ARRAY['MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY']::VARCHAR(20)[],
    'FULL_DAY', '07:00:00', 10, 30, NOW() - INTERVAL '90 days'
  ),
  -- Carlos contrata Fernanda (Physiotherapist) - PENDING
  (
    '00000000-0000-0000-0003-000000000003',
    '00000000-0000-0000-0000-000000000004',
    '00000000-0000-0000-0001-000000000003',
    'PENDING', 12000, 120000,
    '{"description": "Fisioterapia motora para reabilitação pós-AVC", "observations": "10 sessões de 1 hora"}',
    ARRAY['TUESDAY','THURSDAY']::VARCHAR(20)[],
    'MORNING', '09:00:00', 1, 10, NOW() - INTERVAL '2 days'
  ),
  -- João contrata Lucas (Hospital Companion) - CANCELLED
  (
    '00000000-0000-0000-0003-000000000004',
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0001-000000000004',
    'CANCELLED', 4500, 36000,
    '{"description": "Acompanhamento hospitalar para procedimento cirúrgico", "observations": "Cirurgia remarcada", "cancellation_reason": "Procedimento adiado pelo hospital"}',
    ARRAY['SATURDAY']::VARCHAR(20)[],
    'AFTERNOON', '13:00:00', 8, 8, NOW() - INTERVAL '15 days'
  ),
  -- Maria contrata Ana Paula (Nurse) - ACTIVE (segundo contrato)
  (
    '00000000-0000-0000-0003-000000000005',
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0001-000000000001',
    'ACTIVE', 8000, 480000,
    '{"description": "Administração de medicamentos e curativos diários", "observations": "Paciente diabético com úlcera venosa"}',
    ARRAY['MONDAY','TUESDAY','WEDNESDAY','THURSDAY','FRIDAY','SATURDAY']::VARCHAR(20)[],
    'MORNING', '07:00:00', 2, 60, NOW() - INTERVAL '10 days'
  )
ON CONFLICT (id) DO NOTHING;

-- Addresses (endereços de usuários e contratos)
INSERT INTO addresses (id, user_id, contract_id, zip_code, address_line, number, complement, district, city, state, reference, created_at, updated_at) VALUES
  -- João Silva (SP)
  ('00000000-0000-0000-0004-000000000001', '00000000-0000-0000-0000-000000000002', NULL,
   '01310-100', 'Avenida Paulista', '1000', 'Apto 42', 'Bela Vista', 'São Paulo', 'SP', 'Próximo ao MASP', NOW(), NOW()),
  -- Maria Oliveira (RJ)
  ('00000000-0000-0000-0004-000000000002', '00000000-0000-0000-0000-000000000003', NULL,
   '22041-001', 'Rua Barão da Torre', '500', NULL, 'Ipanema', 'Rio de Janeiro', 'RJ', NULL, NOW(), NOW()),
  -- Carlos Souza (MG)
  ('00000000-0000-0000-0004-000000000003', '00000000-0000-0000-0000-000000000004', NULL,
   '30112-000', 'Avenida Afonso Pena', '2000', 'Casa', 'Centro', 'Belo Horizonte', 'MG', 'Em frente ao parque municipal', NOW(), NOW()),
  -- Ana Paula Ferreira (SP)
  ('00000000-0000-0000-0004-000000000004', '00000000-0000-0000-0000-000000000005', NULL,
   '01415-001', 'Rua Augusta', '200', 'Sala 5', 'Consolação', 'São Paulo', 'SP', NULL, NOW(), NOW()),
  -- Endereço do contrato João/Ana Paula (ACTIVE)
  ('00000000-0000-0000-0004-000000000005', NULL, '00000000-0000-0000-0003-000000000001',
   '01310-100', 'Avenida Paulista', '1000', 'Apto 42', 'Bela Vista', 'São Paulo', 'SP', 'Interfone 42', NOW(), NOW()),
  -- Endereço do contrato Maria/Roberto (COMPLETED)
  ('00000000-0000-0000-0004-000000000006', NULL, '00000000-0000-0000-0003-000000000002',
   '22041-001', 'Rua Barão da Torre', '500', NULL, 'Ipanema', 'Rio de Janeiro', 'RJ', NULL, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Reviews (somente para contratos COMPLETED)
INSERT INTO reviews (id, client_id, professional_id, contract_id, rating, comment, created_at, updated_at) VALUES
  (
    '00000000-0000-0000-0005-000000000001',
    '00000000-0000-0000-0000-000000000003',
    '00000000-0000-0000-0001-000000000002',
    '00000000-0000-0000-0003-000000000002',
    5,
    'Roberto foi excelente! Muito atencioso com minha mãe, sempre pontual e profissional. Recomendo muito!',
    NOW(), NOW()
  )
ON CONFLICT (id) DO NOTHING;
