-- +goose Up

CREATE TABLE retros (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    retro_id UUID NOT NULL REFERENCES retros(id) ON DELETE CASCADE,
    column_type TEXT NOT NULL CHECK (column_type IN ('went_well', 'didnt_go_well', 'puzzling', 'action_item')),
    content TEXT NOT NULL,
    author_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_cards_retro_id ON cards (retro_id);

CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    retro_id UUID NOT NULL UNIQUE REFERENCES retros(id) ON DELETE CASCADE,
    result JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_analyses_retro_id ON analyses (retro_id);

-- +goose Down

DROP TABLE IF EXISTS analyses;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS retros;
