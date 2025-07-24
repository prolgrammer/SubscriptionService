CREATE TABLE IF NOT EXISTS subscriptions
(
    id UUID default gen_random_uuid() primary key,
    service_name VARCHAR(255) not null,
    price INTEGER not null check (price > 0),
    user_id UUID not null,
    start_date DATE not null,
    end_date DATE
);

CREATE INDEX IF NOT EXISTS idx_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_service_name ON subscriptions(service_name);
