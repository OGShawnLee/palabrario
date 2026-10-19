-- ==============================================================================
-- PALABRARIO - POSTGRESQL DDL COMPLETO CON FULL-TEXT SEARCH (FTS)
-- Motor: PostgreSQL 17
-- Extensiones: unaccent, pg_trgm (búsqueda difusa/trigramas)
-- Configuración de idioma: 'spanish' con tsvector indexado en GIN
-- ==============================================================================

-- ------------------------------------------------------------------------------
-- 1. EXTENSIONES INDISPENSABLES
-- ------------------------------------------------------------------------------
-- Soporte para ignorar acentos/tildes en español (ej. 'fácil' encuentra 'facil')
CREATE EXTENSION IF NOT EXISTS unaccent;
-- Soporte para autocompletado y tolerancia a errores tipográficos (fuzzy search)
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Configuración FTS en español insensible a acentos
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_ts_config WHERE cfgname = 'spanish_unaccent'
  ) THEN
    CREATE TEXT SEARCH CONFIGURATION spanish_unaccent (COPY = spanish);
    ALTER TEXT SEARCH CONFIGURATION spanish_unaccent
      ALTER MAPPING FOR hword, hword_part, word
      WITH unaccent, spanish_stem;
  END IF;
END $$;

-- ------------------------------------------------------------------------------
-- 2. TIPOS ENUMERADOS
-- ------------------------------------------------------------------------------
CREATE TYPE account_role AS ENUM ('LEARNER', 'EDITOR', 'ADMIN');
CREATE TYPE request_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');
CREATE TYPE word_gender AS ENUM ('MALE', 'FEMALE', 'NEUTRAL', 'AMBIGUOUS');
CREATE TYPE cefr_level AS ENUM ('A1', 'A2', 'B1', 'B2', 'C1', 'C2', 'NATIVO');
CREATE TYPE usage_frequency AS ENUM ('EVERYDAY', 'COMMON', 'RARE', 'SLANG', 'ARCHAIC');
CREATE TYPE review_state AS ENUM ('DRAFT', 'APPROVED', 'HIDDEN');
CREATE TYPE review_decision AS ENUM ('APPROVE', 'REQUEST_CHANGES', 'REJECT');
CREATE TYPE report_target_type AS ENUM ('WORD', 'SENSE', 'EXAMPLE');
CREATE TYPE report_reason AS ENUM ('INACCURATE', 'OFFENSIVE_UNMARKED', 'SPAM_OR_VANDALISM', 'DUPLICATE', 'OTHER');
CREATE TYPE report_status AS ENUM ('PENDING', 'RESOLVED_HIDDEN', 'RESOLVED_AMENDED', 'DISMISSED');
CREATE TYPE mastery_stage AS ENUM ('NEW', 'LEARNING', 'REVIEWING', 'MASTERED');
CREATE TYPE xp_reason AS ENUM ('QUIZ_CORRECT', 'MATCH_WIN', 'DAILY_STREAK', 'DECK_COMPLETED', 'ACHIEVEMENT_UNLOCKED');
CREATE TYPE achievement_tier AS ENUM ('BRONZE', 'SILVER', 'GOLD', 'PLATINUM', 'DIAMOND');
CREATE TYPE achievement_category AS ENUM ('STREAK', 'COMPETITION', 'MASTERY', 'DIALECT_EXPLORER');
CREATE TYPE match_status AS ENUM ('WAITING', 'IN_PROGRESS', 'FINISHED');
CREATE TYPE question_type AS ENUM (
  'DEFINITION_TO_WORD',
  'WORD_TO_DEFINITION',
  'SENTENCE_COMPLETION',
  'REGIONAL_EQUIVALENT',
  'TONE_OR_REGISTER'
);

-- ------------------------------------------------------------------------------
-- 3. CUENTAS, USUARIOS Y GOBERNANZA DE ACCESO (CU01, CU02, CU03, CU04)
-- ------------------------------------------------------------------------------
CREATE TABLE accounts (
  id_account UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role account_role NOT NULL DEFAULT 'LEARNER',
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
  id_user UUID PRIMARY KEY REFERENCES accounts(id_account) ON DELETE CASCADE,
  name VARCHAR(100) NOT NULL,
  user_name VARCHAR(50) NOT NULL UNIQUE,
  avatar_url TEXT,
  anchor_dialect VARCHAR(10) NOT NULL DEFAULT 'MX',

  -- Gamificación principal
  total_xp INTEGER NOT NULL DEFAULT 0,
  current_level INTEGER NOT NULL DEFAULT 1,
  current_streak INTEGER NOT NULL DEFAULT 0,
  longest_streak INTEGER NOT NULL DEFAULT 0,
  last_activity_at TIMESTAMPTZ NOT NULL DEFAULT now(),

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- CU02 & CU03: Solicitudes de acreditación para ascenso a Editor
CREATE TABLE editor_requests (
  id_request UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  institution_name VARCHAR(150) NOT NULL,
  justification TEXT NOT NULL,
  credentials_url TEXT NOT NULL,
  status request_status NOT NULL DEFAULT 'PENDING',
  reviewed_by UUID REFERENCES users(id_user) ON DELETE SET NULL,
  review_notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  resolved_at TIMESTAMPTZ
);

-- RN02: Un Aprendiz solo puede mantener una solicitud en estado PENDING a la vez
CREATE UNIQUE INDEX idx_unique_pending_editor_request
ON editor_requests(id_user)
WHERE status = 'PENDING';

-- CU04: Boletos de invitación institucionales o personales
CREATE TABLE invitation_tickets (
  id_ticket UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(50) NOT NULL UNIQUE,
  created_by UUID NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,
  granted_role account_role NOT NULL DEFAULT 'LEARNER',
  institution_name VARCHAR(150),
  max_uses INTEGER NOT NULL DEFAULT 1 CHECK (max_uses >= 1),
  current_uses INTEGER NOT NULL DEFAULT 0 CHECK (current_uses >= 0 AND current_uses <= max_uses),
  is_active BOOLEAN NOT NULL DEFAULT true,
  expires_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE ticket_redemptions (
  id_ticket UUID NOT NULL REFERENCES invitation_tickets(id_ticket) ON DELETE CASCADE,
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  redeemed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id_ticket, id_user)
);

-- ------------------------------------------------------------------------------
-- 4. REGIONES Y DIALECTOS
-- ------------------------------------------------------------------------------
CREATE TABLE regions (
  id_region UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  code VARCHAR(10) NOT NULL UNIQUE,         -- 'MX', 'ES', 'AR', 'CO', 'CL', 'CU', 'CAM', 'PAN'
  name VARCHAR(100) NOT NULL,
  parent_zone VARCHAR(100) NOT NULL,        -- 'Norteamérica', 'Rioplatense', 'Andina', 'Peninsular', 'Caribe'
  country_code VARCHAR(2),                  -- 'MX', 'ES', 'AR'
  flag_emoji VARCHAR(10) NOT NULL
);

-- ------------------------------------------------------------------------------
-- 5. NÚCLEO LÉXICO: WORDS, SENSES & EXAMPLES (CU05, CU06)
-- ------------------------------------------------------------------------------
CREATE TABLE words (
  id_word SERIAL PRIMARY KEY,
  lemma VARCHAR(100) NOT NULL,
  normalized_lemma VARCHAR(100) NOT NULL,   -- Minúsculas y sin acentos
  part_of_speech VARCHAR(50) NOT NULL,      -- 'sustantivo', 'verbo', 'modismo'
  gender word_gender DEFAULT 'NEUTRAL',
  etymology TEXT,

  -- El blurb estilo Vocabulary.com a nivel de palabra
  friendly_blurb TEXT,

  -- Región de origen endémica (null si es español panhispánico/general)
  id_origin_region UUID REFERENCES regions(id_region) ON DELETE SET NULL,

  -- Vector FTS generado automáticamente para búsquedas instantáneas
  search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('spanish_unaccent', coalesce(lemma, '')), 'A') ||
    setweight(to_tsvector('spanish_unaccent', coalesce(friendly_blurb, '')), 'B')
  ) STORED,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Índices de búsqueda para autocompletado y FTS
CREATE INDEX idx_words_search_vector ON words USING GIN(search_vector);
CREATE INDEX idx_words_normalized_trgm ON words USING GIN(normalized_lemma gin_trgm_ops);
CREATE INDEX idx_words_lemma_prefix ON words (normalized_lemma varchar_pattern_ops);

CREATE TABLE senses (
  id_sense SERIAL PRIMARY KEY,
  id_word INTEGER NOT NULL REFERENCES words(id_word) ON DELETE CASCADE,
  id_region UUID NOT NULL REFERENCES regions(id_region) ON DELETE RESTRICT,
  id_editor UUID NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,

  short_definition VARCHAR(255) NOT NULL,   -- Glosa rápida mostrada en el dropdown de búsqueda
  definition TEXT NOT NULL,
  is_primary BOOLEAN NOT NULL DEFAULT true,  -- Utilizado por el buscador y las trivias

  is_taboo_or_vulgar BOOLEAN NOT NULL DEFAULT false,
  cultural_note TEXT,

  cefr_level cefr_level DEFAULT 'B1',
  usage_frequency usage_frequency DEFAULT 'COMMON',
  state review_state NOT NULL DEFAULT 'DRAFT',

  -- Vector FTS para buscar directamente dentro de las definiciones de los sentidos
  sense_search_vector tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('spanish_unaccent', coalesce(short_definition, '')), 'A') ||
    setweight(to_tsvector('spanish_unaccent', coalesce(definition, '')), 'B') ||
    setweight(to_tsvector('spanish_unaccent', coalesce(cultural_note, '')), 'C')
  ) STORED,

  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_senses_search_vector ON senses USING GIN(sense_search_vector);
CREATE INDEX idx_senses_primary_lookup ON senses(id_word, is_primary) WHERE state = 'APPROVED';

CREATE TABLE examples (
  id_example SERIAL PRIMARY KEY,
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  id_editor UUID NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,
  sentence TEXT NOT NULL,
  source VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ------------------------------------------------------------------------------
-- 6. EQUIVALENCIAS TRANSDIALECTALES (Heteronimia / ¿Cómo se dice allá?)
-- ------------------------------------------------------------------------------
CREATE TABLE semantic_concepts (
  id_concept SERIAL PRIMARY KEY,
  concept_key VARCHAR(100) NOT NULL UNIQUE, -- 'drinking_straw', 'popcorn', 'avocado'
  description TEXT NOT NULL
);

CREATE TABLE concept_mappings (
  id_concept INTEGER NOT NULL REFERENCES semantic_concepts(id_concept) ON DELETE CASCADE,
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  PRIMARY KEY (id_concept, id_sense)
);

-- ------------------------------------------------------------------------------
-- 7. GOBERNANZA COLABORATIVA Y MODERACIÓN (CU07, CU08, CU09)
-- ------------------------------------------------------------------------------
CREATE TABLE sense_approvals (
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  id_editor UUID NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,
  decision review_decision NOT NULL DEFAULT 'APPROVE',
  feedback_notes TEXT,
  approved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id_sense, id_editor)
);

-- RN07: Un editor no puede votar sobre su propia propuesta
CREATE OR REPLACE FUNCTION fn_prevent_self_approval()
RETURNS TRIGGER AS $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM senses
    WHERE id_sense = NEW.id_sense AND id_editor = NEW.id_editor
  ) THEN
    RAISE EXCEPTION 'Regla RN07: Un editor no puede emitir votos sobre su propia propuesta.';
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_self_approval
BEFORE INSERT OR UPDATE ON sense_approvals
FOR EACH ROW EXECUTE FUNCTION fn_prevent_self_approval();

-- RN06: Tras 3 aprobaciones de editores distintos, se publica automáticamente
CREATE OR REPLACE FUNCTION fn_check_sense_approvals()
RETURNS TRIGGER AS $$
DECLARE
  approval_count INTEGER;
BEGIN
  SELECT count(*) INTO approval_count
  FROM sense_approvals
  WHERE id_sense = NEW.id_sense AND decision = 'APPROVE';

  IF approval_count >= 3 THEN
    UPDATE senses
    SET state = 'APPROVED', updated_at = now()
    WHERE id_sense = NEW.id_sense;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_auto_approve_sense
AFTER INSERT OR UPDATE ON sense_approvals
FOR EACH ROW EXECUTE FUNCTION fn_check_sense_approvals();

-- CU09: Reportes de contenido léxico emitidos por la comunidad
CREATE TABLE content_reports (
  id_report UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  id_reporter UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  target_type report_target_type NOT NULL DEFAULT 'SENSE',
  id_word INTEGER REFERENCES words(id_word) ON DELETE CASCADE,
  id_sense INTEGER REFERENCES senses(id_sense) ON DELETE CASCADE,
  id_example INTEGER REFERENCES examples(id_example) ON DELETE CASCADE,
  reason report_reason NOT NULL,
  details TEXT NOT NULL,
  status report_status NOT NULL DEFAULT 'PENDING',
  resolved_by UUID REFERENCES users(id_user) ON DELETE SET NULL,
  resolution_notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  resolved_at TIMESTAMPTZ
);

CREATE INDEX idx_content_reports_status ON content_reports(status);

-- ------------------------------------------------------------------------------
-- 8. STUDIO: DECKS, COLECCIONES Y ELEMENTOS GUARDADOS (CU10, CU11)
-- ------------------------------------------------------------------------------
CREATE TABLE decks (
  id_deck SERIAL PRIMARY KEY,
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  title VARCHAR(150) NOT NULL,
  description TEXT,
  is_public BOOLEAN NOT NULL DEFAULT true,
  is_official BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE deck_items (
  id_deck INTEGER NOT NULL REFERENCES decks(id_deck) ON DELETE CASCADE,
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (id_deck, id_sense)
);

-- Suscripción/Guardado de Barajas Públicas en el Studio de otros miembros
CREATE TABLE saved_decks (
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  id_deck INTEGER NOT NULL REFERENCES decks(id_deck) ON DELETE CASCADE,
  saved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id_user, id_deck)
);

CREATE TABLE collections (
  id_collection SERIAL PRIMARY KEY,
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  title VARCHAR(150) NOT NULL,
  description TEXT,
  is_public BOOLEAN NOT NULL DEFAULT true,
  is_featured BOOLEAN NOT NULL DEFAULT false,
  is_official BOOLEAN NOT NULL DEFAULT false,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE collection_decks (
  id_collection INTEGER NOT NULL REFERENCES collections(id_collection) ON DELETE CASCADE,
  id_deck INTEGER NOT NULL REFERENCES decks(id_deck) ON DELETE CASCADE,
  sort_order INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (id_collection, id_deck)
);

-- Suscripción/Guardado de Colecciones Públicas en el Studio de otros miembros (CU11 1a)
CREATE TABLE saved_collections (
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  id_collection INTEGER NOT NULL REFERENCES collections(id_collection) ON DELETE CASCADE,
  saved_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id_user, id_collection)
);

-- ------------------------------------------------------------------------------
-- 9. MOTOR DE APRENDIZAJE (Mastery / SM-2) (CU12, CU13)
-- ------------------------------------------------------------------------------
CREATE TABLE user_mastery (
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  progress_stage mastery_stage NOT NULL DEFAULT 'NEW',
  mastery_percentage REAL NOT NULL DEFAULT 0.0,
  consecutive_correct INTEGER NOT NULL DEFAULT 0,
  times_reviewed INTEGER NOT NULL DEFAULT 0,
  
  -- Algoritmo matemático SuperMemo SM-2
  easiness_factor REAL NOT NULL DEFAULT 2.5,
  current_interval_days INTEGER NOT NULL DEFAULT 1,
  
  next_review_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_reviewed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id_user, id_sense)
);

CREATE INDEX idx_user_mastery_due_review ON user_mastery(id_user, next_review_at);

-- ------------------------------------------------------------------------------
-- 10. GAMIFICACIÓN: LEDGER DE XP Y LOGROS (CU12, CU13, CU16)
-- ------------------------------------------------------------------------------
CREATE TABLE xp_transactions (
  id_transaction UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  amount_xp INTEGER NOT NULL,
  reason xp_reason NOT NULL,
  metadata_id VARCHAR(100),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_xp_transactions_user ON xp_transactions(id_user, created_at DESC);

CREATE TABLE achievements (
  id_achievement SERIAL PRIMARY KEY,
  slug VARCHAR(100) NOT NULL UNIQUE,
  title VARCHAR(150) NOT NULL,
  description TEXT NOT NULL,
  badge_icon_url TEXT,
  xp_reward INTEGER NOT NULL DEFAULT 100,
  tier achievement_tier NOT NULL DEFAULT 'BRONZE',
  category achievement_category NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_achievements (
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  id_achievement INTEGER NOT NULL REFERENCES achievements(id_achievement) ON DELETE CASCADE,
  unlocked_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (id_user, id_achievement)
);

-- ------------------------------------------------------------------------------
-- 11. MULTIJUGADOR EN TIEMPO REAL: SALAS Y RESULTADOS EN LOTE (CU14, CU15, CU16)
-- ------------------------------------------------------------------------------
CREATE TABLE matches (
  id_match UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  id_deck INTEGER NOT NULL REFERENCES decks(id_deck) ON DELETE CASCADE,
  host_user_id UUID REFERENCES users(id_user) ON DELETE SET NULL,
  room_code VARCHAR(10) NOT NULL UNIQUE,
  status match_status NOT NULL DEFAULT 'WAITING',
  winner_user_id UUID REFERENCES users(id_user) ON DELETE SET NULL,
  is_public BOOLEAN NOT NULL DEFAULT true,
  total_rounds INTEGER NOT NULL DEFAULT 10,
  round_time_sec INTEGER NOT NULL DEFAULT 15,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ
);

CREATE TABLE match_participants (
  id_match UUID NOT NULL REFERENCES matches(id_match) ON DELETE CASCADE,
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  score INTEGER NOT NULL DEFAULT 0,
  final_rank INTEGER NOT NULL DEFAULT 0,
  correct_answers INTEGER NOT NULL DEFAULT 0,
  total_time_spent_ms INTEGER NOT NULL DEFAULT 0,
  xp_earned INTEGER NOT NULL DEFAULT 0,
  PRIMARY KEY (id_match, id_user)
);

-- ------------------------------------------------------------------------------
-- 12. BANCO DE PREGUNTAS CURADAS E INTENTOS INDIVIDUALES
-- ------------------------------------------------------------------------------
CREATE TABLE questions (
  id_question SERIAL PRIMARY KEY,
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  id_editor UUID NOT NULL REFERENCES users(id_user) ON DELETE RESTRICT,
  type question_type NOT NULL,
  prompt TEXT NOT NULL,
  custom_distractors JSONB NOT NULL,       -- Array JSON de 3 distractores: ["pitillo", "bombilla", "pajita"]
  explanation TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE quiz_attempts (
  id_attempt UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  id_user UUID NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
  id_sense INTEGER NOT NULL REFERENCES senses(id_sense) ON DELETE CASCADE,
  id_question INTEGER REFERENCES questions(id_question) ON DELETE SET NULL,
  selected_option VARCHAR(255) NOT NULL,
  is_correct BOOLEAN NOT NULL,
  quality_rating SMALLINT CHECK (quality_rating BETWEEN 0 AND 5),
  time_spent_ms INTEGER NOT NULL,
  answered_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==============================================================================
-- 13. VISTA Y FUNCIÓN OPTIMIZADA DE BÚSQUEDA (ESTILO VOCABULARY.COM)
-- ==============================================================================

-- Función rápida para el dropdown del autocompletado Master-Detail (sin duplicados)
CREATE OR REPLACE FUNCTION search_vocabulary_dropdown(
  search_query TEXT,
  user_dialect VARCHAR DEFAULT 'MX',
  max_results INT DEFAULT 6
)
RETURNS TABLE (
  id_word INT,
  lemma VARCHAR,
  short_definition VARCHAR,
  flag_emoji VARCHAR,
  similarity REAL
) AS $$
  WITH matched_words AS (
    SELECT
      w.id_word,
      w.lemma,
      w.normalized_lemma,
      similarity(w.normalized_lemma, unaccent(lower(search_query))) AS sim
    FROM words w
    WHERE
      search_query IS NOT NULL AND length(trim(search_query)) > 0 AND (
        w.normalized_lemma ILIKE unaccent(lower(search_query)) || '%'
        OR w.search_vector @@ plainto_tsquery('spanish_unaccent', search_query)
        OR similarity(w.normalized_lemma, unaccent(lower(search_query))) > 0.25
      )
    ORDER BY
      w.normalized_lemma LIKE unaccent(lower(search_query)) || '%' DESC,
      sim DESC,
      w.lemma ASC
    LIMIT 15
  ),
  distinct_words AS (
    SELECT DISTINCT ON (mw.id_word)
      mw.id_word,
      mw.lemma,
      mw.normalized_lemma,
      coalesce(s.short_definition, 'Sin definición corta') AS short_definition,
      coalesce(r.flag_emoji, '🌐') AS flag_emoji,
      mw.sim AS similarity
    FROM matched_words mw
    LEFT JOIN senses s ON s.id_word = mw.id_word AND s.state = 'APPROVED'
    LEFT JOIN regions r ON r.id_region = s.id_region
    ORDER BY
      mw.id_word,
      (r.code = user_dialect) DESC NULLS LAST,
      s.is_primary DESC
  )
  SELECT
    dw.id_word,
    dw.lemma,
    dw.short_definition,
    dw.flag_emoji,
    dw.similarity
  FROM distinct_words dw
  ORDER BY
    dw.normalized_lemma LIKE unaccent(lower(search_query)) || '%' DESC,
    dw.similarity DESC,
    dw.lemma ASC
  LIMIT max_results;
$$ LANGUAGE sql STABLE;

-- ==============================================================================
-- 14. DATOS SEMILLA DETERMINÍSTICOS (SEEDS)
-- ==============================================================================

-- 1. Regiones Dialectales
INSERT INTO regions (id_region, code, name, parent_zone, country_code, flag_emoji) VALUES
('a1111111-0000-0000-0000-000000000001', 'MX', 'México', 'Norteamérica', 'MX', '🇲🇽'),
('a1111111-0000-0000-0000-000000000002', 'ES', 'España', 'Peninsular', 'ES', '🇪🇸'),
('a1111111-0000-0000-0000-000000000003', 'AR', 'Argentina', 'Rioplatense', 'AR', '🇦🇷'),
('a1111111-0000-0000-0000-000000000004', 'CO', 'Colombia', 'Andina', 'CO', '🇨🇴'),
('a1111111-0000-0000-0000-000000000005', 'CL', 'Chile', 'Cono Sur', 'CL', '🇨🇱'),
('a1111111-0000-0000-0000-000000000006', 'CU', 'Cuba', 'Caribe', 'CU', '🇨🇺'),
('a1111111-0000-0000-0000-000000000007', 'CAM', 'Centroamérica', 'Centroamérica', NULL, '🌎'),
('a1111111-0000-0000-0000-000000000008', 'PAN', 'Panhispánico', 'General', NULL, '🌐')
ON CONFLICT (code) DO NOTHING;

-- 2. Cuentas y Usuarios demo (Contraseña: 'password123')
INSERT INTO accounts (id_account, email, password_hash, role) VALUES
('b1111111-0000-0000-0000-000000000001', 'admin@palabrario.org', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8.3dtOfsvAPGNGDGqOSwhNGEiccrdUHFKnByP0S9PsEfs6c0i', 'ADMIN'),
('b1111111-0000-0000-0000-000000000002', 'silas@palabrario.org', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8.3dtOfsvAPGNGDGqOSwhNGEiccrdUHFKnByP0S9PsEfs6c0i', 'EDITOR'),
('b1111111-0000-0000-0000-000000000003', 'gael@palabrario.org', '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8.3dtOfsvAPGNGDGqOSwhNGEiccrdUHFKnByP0S9PsEfs6c0i', 'LEARNER')
ON CONFLICT (email) DO NOTHING;

INSERT INTO users (id_user, name, user_name, anchor_dialect, total_xp, current_level) VALUES
('b1111111-0000-0000-0000-000000000001', 'Administrador Palabrario', 'admin', 'MX', 15000, 10),
('b1111111-0000-0000-0000-000000000002', 'Silas Lingüista', 'silas_editor', 'MX', 8500, 7),
('b1111111-0000-0000-0000-000000000003', 'Gael Aprendiz', 'gael_learner', 'MX', 4850, 4)
ON CONFLICT (user_name) DO NOTHING;

-- 3. Núcleo Léxico: Palabras canónicas
INSERT INTO words (id_word, lemma, normalized_lemma, part_of_speech, gender, etymology, friendly_blurb, id_origin_region) VALUES
(1, 'popote', 'popote', 'sustantivo', 'MALE', 'Del náhuatl popotl (tallo hueco o paja seca).',
 'Si vas a una taquería en México y pides refresco, te darán un popote. Pero si cruzas fronteras, prepárate: en Colombia es "pitillo", en Argentina "sorbete" o "pajita", y en Chile "bombilla". Una misma pajita, ¡media docena de identidades!',
 'a1111111-0000-0000-0000-000000000001'),

(2, 'apapachar', 'apapachar', 'verbo', 'NEUTRAL', 'Del náhuatl papatzoa (ablandar con cariño o masajes afectuosos).',
 'Catalogada muchas veces como la palabra más hermosa del español mexicano. Los lingüistas y poetas la definen como "acariciar con el alma". Cuando alguien tuvo un día pesado, no le des consejos: dale un buen apapacho.',
 'a1111111-0000-0000-0000-000000000001'),

(3, 'guagua', 'guagua', 'sustantivo', 'NEUTRAL', 'Poligenético: del quechua wawa (bebé) en Sudamérica; origen caribeño para el transporte público.',
 'Una de las palabras trampa más fascinantes del idioma. Si estás en La Habana o en Tenerife y te dicen "corre a la guagua", súbete al camión. Pero si estás en Santiago de Chile o Quito, ¡guagua es un tierno bebé de brazos!',
 'a1111111-0000-0000-0000-000000000008'),

(4, 'chamaco', 'chamaco', 'sustantivo', 'MALE', 'Del maya chamak (niño, muchacho).',
 'Manera entrañable y coloquial de llamar a los niños y adolescentes en México y Centroamérica. Transmite picardía infantil, energía desbordante y cariño popular.',
 'a1111111-0000-0000-0000-000000000001'),

(5, 'coger', 'coger', 'verbo', 'NEUTRAL', 'Del latín colligĕre (juntar, recoger, tomar).',
 'El falso amigo supremo del mundo hispanohablante. Mientras en Madrid tomas el autobús diciendo "voy a coger el bus", en México y Buenos Aires ese verbo es tabú y tiene una fuerte carga sexual explícita. ¡En América di siempre "tomar" o "agarrar"!',
 'a1111111-0000-0000-0000-000000000002'),

(6, 'tianguis', 'tianguis', 'sustantivo', 'MALE', 'Del náhuatl tiyānquiztli (mercado público al aire libre).',
 'El corazón de los barrios mexicanos. No es un supermercado ni un centro comercial: es un laberinto colorido de lonas rosas y amarillas, marchantes gritando ofertas y puestos de comida callejera que se instalan un día a la semana.',
 'a1111111-0000-0000-0000-000000000001'),

(7, 'chamba', 'chamba', 'sustantivo', 'FEMALE', 'Probablemente de "chamber" (cámara de comercio) en el contexto de braceros mexicanos.',
 'Término informal por excelencia para referirse al empleo o al trabajo arduo. Si un amigo te dice "ando en la chamba", significa que está ocupado ganándose el pan.',
 'a1111111-0000-0000-0000-000000000001')
ON CONFLICT (id_word) DO NOTHING;

-- 4. Sentidos aprobados (Senses)
INSERT INTO senses (id_sense, id_word, id_region, id_editor, short_definition, definition, is_primary, is_taboo_or_vulgar, cultural_note, cefr_level, usage_frequency, state) VALUES
(101, 1, 'a1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000002',
 'Tubo hueco de plástico o papel para sorber bebidas.',
 'Cilindro hueco y delgado utilizado comúnmente para extraer líquidos de un recipiente hacia la boca.',
 true, false, 'En años recientes se ha impulsado el popote biodegradable de semilla de aguacate.', 'A2', 'EVERYDAY', 'APPROVED'),

(102, 2, 'a1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000002',
 'Acariciar o consentir con mucho afecto y ternura profunda.',
 'Dar muestras de cariño especialmente cálidas, brindando confort emocional y cercanía protectora.',
 true, false, 'Traspasó fronteras hacia Centroamérica y Colombia por su resonancia poética.', 'B1', 'COMMON', 'APPROVED'),

(103, 3, 'a1111111-0000-0000-0000-000000000006', 'b1111111-0000-0000-0000-000000000002',
 'Autobús o vehículo de transporte colectivo público.',
 'Vehículo automotor de gran capacidad destinado al transporte regular de pasajeros.',
 true, false, 'En las Antillas hispanas y las Islas Canarias es el término estándar de diario.', 'B1', 'EVERYDAY', 'APPROVED'),

(104, 3, 'a1111111-0000-0000-0000-000000000005', 'b1111111-0000-0000-0000-000000000002',
 'Niño pequeño, infante o criatura recién nacida.',
 'Bebé de muy corta edad o que aún no camina.',
 false, false, 'Proviene del quechua; uso generalizado en Chile, Bolivia, Ecuador y Perú.', 'B1', 'COMMON', 'APPROVED'),

(105, 4, 'a1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000002',
 'Niño, infante o muchacho joven.',
 'Persona que está en la etapa de la niñez o adolescencia temprana.',
 true, false, 'Muy común en el habla cotidiana urbana y de provincia en México.', 'A2', 'EVERYDAY', 'APPROVED'),

(106, 5, 'a1111111-0000-0000-0000-000000000002', 'b1111111-0000-0000-0000-000000000002',
 'Agarrar, asir, tomar algo con la mano o subirse a un transporte.',
 'Tomar o asir una cosa física o abordar un medio de transporte colectivo.',
 true, false, 'Uso perfectamente formal y cotidiano en toda España peninsular.', 'A1', 'EVERYDAY', 'APPROVED'),

(107, 5, 'a1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000002',
 'Practicar el coito / relación sexual (Tabú coloquial).',
 'Término vulgar de uso informal para denotar el acto sexual.',
 false, true, 'En México y Argentina se considera lenguaje vulgar en contextos formales; usar tomar o agarrar.', 'B2', 'SLANG', 'APPROVED'),

(108, 6, 'a1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000002',
 'Mercado ambulante tradicional instalado al aire libre.',
 'Bazar o mercado sobre ruedas que se instala temporalmente en calles un día fijo.',
 true, false, 'Patrimonio vivo de la cultura comercial prehispánica mexicana.', 'B1', 'EVERYDAY', 'APPROVED'),

(109, 7, 'a1111111-0000-0000-0000-000000000001', 'b1111111-0000-0000-0000-000000000002',
 'Empleo, trabajo u ocupación remunerada.',
 'Actividad laboral cotidiana que genera sustento económico.',
 true, false, 'Expresión muy arraigada de orgullo por el esfuerzo honesto.', 'A2', 'EVERYDAY', 'APPROVED')
ON CONFLICT (id_sense) DO NOTHING;

-- 5. Ejemplos de uso
INSERT INTO examples (id_example, id_sense, id_editor, sentence, source) VALUES
(1, 101, 'b1111111-0000-0000-0000-000000000002', '¿Me regalas un popote para el agua de horchata, por favor?', 'Conversación cotidiana'),
(2, 102, 'b1111111-0000-0000-0000-000000000002', 'Mi abuela siempre me apapachaba preparándome un atole caliente cuando llovía.', 'Narrativa familiar'),
(3, 103, 'b1111111-0000-0000-0000-000000000002', 'Esperamos media hora en la parada pero la guagua venía completamente repleta.', 'Prensa de La Habana'),
(4, 104, 'b1111111-0000-0000-0000-000000000002', 'Abrigaron a la guagüita con una mantita de lana para sacarla al jardín.', 'Literatura chilena'),
(5, 106, 'b1111111-0000-0000-0000-000000000002', 'Voy a coger el metro de la línea 1 para llegar a tiempo a Sol.', 'Madrid cotidiano'),
(6, 108, 'b1111111-0000-0000-0000-000000000002', 'Los domingos vamos al tianguis a comer barbacoa y comprar verdura fresca.', 'Tradición mexicana'),
(7, 109, 'b1111111-0000-0000-0000-000000000002', 'Apenas salí de la chamba y ya hay tráfico en periférico.', 'Habla urbana')
ON CONFLICT (id_example) DO NOTHING;

-- 6. Conceptos Semánticos y Heteronimia
INSERT INTO semantic_concepts (id_concept, concept_key, description) VALUES
(1, 'drinking_straw', 'Tubo delgado de plástico o papel para sorber líquidos.'),
(2, 'infant_baby', 'Niño pequeño o recién nacido que aún no camina.'),
(3, 'public_bus', 'Vehículo automotor de gran capacidad para transporte público.')
ON CONFLICT (concept_key) DO NOTHING;

INSERT INTO concept_mappings (id_concept, id_sense) VALUES
(1, 101), -- popote (MX)
(2, 104), -- guagua (CL)
(3, 103)  -- guagua (CU)
ON CONFLICT (id_concept, id_sense) DO NOTHING;

-- 7. Decks Oficiales Iniciales
INSERT INTO decks (id_deck, id_user, title, description, is_public, is_official) VALUES
(1, 'b1111111-0000-0000-0000-000000000001', 'Mexicanismos Indispensables', 'Vocabulario esencial y modismos de México que todo viajero o estudiante debe dominar.', true, true)
ON CONFLICT (id_deck) DO NOTHING;

INSERT INTO deck_items (id_deck, id_sense, sort_order) VALUES
(1, 101, 1), -- popote
(1, 102, 2), -- apapachar
(1, 105, 3), -- chamaco
(1, 108, 4), -- tianguis
(1, 109, 5)  -- chamba
ON CONFLICT (id_deck, id_sense) DO NOTHING;

-- 8. Catálogo de Logros (Achievements)
INSERT INTO achievements (id_achievement, slug, title, description, badge_icon_url, xp_reward, tier, category) VALUES
(1, 'primeros_pasos', 'Primeros Pasos', 'Consulta tus primeras 5 palabras en el diccionario panhispánico.', 'badge_steps', 50, 'BRONZE', 'MASTERY'),
(2, 'chilango_master', 'Chilango Honorario', 'Domina 10 palabras de la jerga tradicional de la Ciudad de México.', 'badge_mexico', 150, 'SILVER', 'DIALECT_EXPLORER'),
(3, 'ojo_critico', 'Guateque Lingüístico', 'Participa en tu primera partida multijugador en vivo en la Arena.', 'badge_arena', 200, 'GOLD', 'COMPETITION')
ON CONFLICT (slug) DO NOTHING;

-- Ajustar las secuencias automáticas tras insertar datos semilla explícitos
SELECT setval('words_id_word_seq', COALESCE((SELECT MAX(id_word) FROM words), 1));
SELECT setval('senses_id_sense_seq', COALESCE((SELECT MAX(id_sense) FROM senses), 1));
SELECT setval('examples_id_example_seq', COALESCE((SELECT MAX(id_example) FROM examples), 1));
SELECT setval('decks_id_deck_seq', COALESCE((SELECT MAX(id_deck) FROM decks), 1));
SELECT setval('semantic_concepts_id_concept_seq', COALESCE((SELECT MAX(id_concept) FROM semantic_concepts), 1));
SELECT setval('achievements_id_achievement_seq', COALESCE((SELECT MAX(id_achievement) FROM achievements), 1));
