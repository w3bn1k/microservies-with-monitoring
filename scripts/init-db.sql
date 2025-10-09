-- Инициализация базы данных для микросервисного приложения

-- Создание базы данных (если не существует)
-- CREATE DATABASE microservices;

-- Подключение к базе данных
\c microservices;

-- Создание расширений
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Создание таблицы транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    kafka_status VARCHAR(10) NOT NULL CHECK (kafka_status IN ('ok', 'bad')),
    redis_status VARCHAR(10) NOT NULL CHECK (redis_status IN ('ok', 'bad')),
    duration_ms BIGINT NOT NULL CHECK (duration_ms >= 0),
    service VARCHAR(50) NOT NULL,
    event_id VARCHAR(100),
    error_msg TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Создание индексов для оптимизации запросов
CREATE INDEX IF NOT EXISTS idx_transactions_timestamp ON transactions(timestamp);
CREATE INDEX IF NOT EXISTS idx_transactions_service ON transactions(service);
CREATE INDEX IF NOT EXISTS idx_transactions_kafka_status ON transactions(kafka_status);
CREATE INDEX IF NOT EXISTS idx_transactions_redis_status ON transactions(redis_status);
CREATE INDEX IF NOT EXISTS idx_transactions_event_id ON transactions(event_id);

-- Создание таблицы событий (для истории)
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(50) NOT NULL,
    user_id VARCHAR(100),
    data JSONB,
    source VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processed', 'failed')),
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для таблицы событий
CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
CREATE INDEX IF NOT EXISTS idx_events_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_status ON events(status);
CREATE INDEX IF NOT EXISTS idx_events_source ON events(source);

-- Создание таблицы метрик
CREATE TABLE IF NOT EXISTS metrics (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(50) NOT NULL,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(15,4) NOT NULL,
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для таблицы метрик
CREATE INDEX IF NOT EXISTS idx_metrics_service ON metrics(service_name);
CREATE INDEX IF NOT EXISTS idx_metrics_name ON metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp);

-- Создание таблицы health checks
CREATE TABLE IF NOT EXISTS health_checks (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('healthy', 'unhealthy', 'degraded')),
    response_time_ms INTEGER,
    error_message TEXT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для таблицы health checks
CREATE INDEX IF NOT EXISTS idx_health_checks_service ON health_checks(service_name);
CREATE INDEX IF NOT EXISTS idx_health_checks_timestamp ON health_checks(timestamp);
CREATE INDEX IF NOT EXISTS idx_health_checks_status ON health_checks(status);

-- Создание функции для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создание триггера для автоматического обновления updated_at
CREATE TRIGGER update_transactions_updated_at 
    BEFORE UPDATE ON transactions 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Создание представления для статистики транзакций
CREATE OR REPLACE VIEW transaction_stats AS
SELECT 
    service,
    DATE_TRUNC('hour', timestamp) as hour,
    COUNT(*) as total_transactions,
    COUNT(CASE WHEN kafka_status = 'ok' THEN 1 END) as kafka_success,
    COUNT(CASE WHEN kafka_status = 'bad' THEN 1 END) as kafka_failures,
    COUNT(CASE WHEN redis_status = 'ok' THEN 1 END) as redis_success,
    COUNT(CASE WHEN redis_status = 'bad' THEN 1 END) as redis_failures,
    AVG(duration_ms) as avg_duration_ms,
    MAX(duration_ms) as max_duration_ms,
    MIN(duration_ms) as min_duration_ms
FROM transactions
GROUP BY service, DATE_TRUNC('hour', timestamp)
ORDER BY hour DESC, service;

-- Создание представления для статистики событий
CREATE OR REPLACE VIEW event_stats AS
SELECT 
    event_type,
    source,
    DATE_TRUNC('hour', timestamp) as hour,
    COUNT(*) as total_events,
    COUNT(CASE WHEN status = 'processed' THEN 1 END) as processed_events,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_events,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_events
FROM events
GROUP BY event_type, source, DATE_TRUNC('hour', timestamp)
ORDER BY hour DESC, event_type, source;

-- Создание представления для метрик сервисов
CREATE OR REPLACE VIEW service_metrics AS
SELECT 
    service_name,
    metric_name,
    DATE_TRUNC('minute', timestamp) as minute,
    AVG(metric_value) as avg_value,
    MAX(metric_value) as max_value,
    MIN(metric_value) as min_value,
    COUNT(*) as sample_count
FROM metrics
GROUP BY service_name, metric_name, DATE_TRUNC('minute', timestamp)
ORDER BY minute DESC, service_name, metric_name;

-- Создание пользователя для приложения (опционально)
-- CREATE USER app_user WITH PASSWORD 'app_password';
-- GRANT CONNECT ON DATABASE microservices TO app_user;
-- GRANT USAGE ON SCHEMA public TO app_user;
-- GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
-- GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO app_user;

-- Вставка тестовых данных (опционально)
INSERT INTO transactions (kafka_status, redis_status, duration_ms, service, event_id) VALUES
('ok', 'ok', 150, 'producer', 'test-event-1'),
('ok', 'bad', 200, 'consumer', 'test-event-2'),
('bad', 'ok', 100, 'monitor', 'test-event-3'),
('ok', 'ok', 300, 'api-gateway', 'test-event-4')
ON CONFLICT DO NOTHING;

-- Создание партиций для больших таблиц (для production)
-- Пример для таблицы transactions по месяцам
-- CREATE TABLE transactions_y2024m01 PARTITION OF transactions
-- FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Создание политик безопасности (Row Level Security)
-- ALTER TABLE transactions ENABLE ROW LEVEL SECURITY;
-- CREATE POLICY transactions_policy ON transactions
--     FOR ALL TO app_user
--     USING (true);

-- Настройка логирования медленных запросов
-- ALTER SYSTEM SET log_min_duration_statement = 1000;
-- ALTER SYSTEM SET log_statement = 'mod';
-- SELECT pg_reload_conf();

-- Создание backup скрипта
-- CREATE OR REPLACE FUNCTION backup_transactions()
-- RETURNS void AS $$
-- BEGIN
--     COPY transactions TO '/tmp/transactions_backup.csv' WITH CSV HEADER;
-- END;
-- $$ LANGUAGE plpgsql;

-- Комментарии к таблицам
COMMENT ON TABLE transactions IS 'Таблица для хранения транзакций между сервисами';
COMMENT ON TABLE events IS 'Таблица для хранения событий системы';
COMMENT ON TABLE metrics IS 'Таблица для хранения метрик сервисов';
COMMENT ON TABLE health_checks IS 'Таблица для хранения результатов health checks';

COMMENT ON COLUMN transactions.kafka_status IS 'Статус отправки в Kafka: ok или bad';
COMMENT ON COLUMN transactions.redis_status IS 'Статус операции с Redis: ok или bad';
COMMENT ON COLUMN transactions.duration_ms IS 'Длительность операции в миллисекундах';
COMMENT ON COLUMN transactions.service IS 'Название сервиса, выполнившего операцию';

-- Вывод информации о созданных объектах
SELECT 'Database initialization completed successfully!' as status;
SELECT 'Created tables: ' || string_agg(tablename, ', ') as tables
FROM pg_tables 
WHERE schemaname = 'public' AND tablename IN ('transactions', 'events', 'metrics', 'health_checks');
