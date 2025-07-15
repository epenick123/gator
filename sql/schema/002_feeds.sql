CREATE TABLE feeds (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP not null,
    updated_at TIMESTAMP not null,
    name TEXT unique not null,
    url TEXT unique not null,
    user_id UUID not null REFERENCES users (id) ON DELETE CASCADE  
);