CREATE TABLE IF NOT EXISTS changes_history (
    id UInt64,
    project_id Int32, 
    name String, 
    description String, 
    priority UInt64, 
    removed UInt8, 
    event_time DateTime DEFAULT now()
) ENGINE = NATS SETTINGS
    nats_url = '127.0.0.1:4222',
    nats_subjects = 'logs',
    nats_format = 'JSONEachRow'