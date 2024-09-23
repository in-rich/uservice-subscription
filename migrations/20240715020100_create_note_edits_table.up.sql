CREATE TYPE note_target AS ENUM ('user', 'company');

--bun:split

CREATE TABLE note_edits (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    author_id         VARCHAR(255) NOT NULL,

    public_identifier VARCHAR(255) NOT NULL,
    target            note_target NOT NULL,

    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

--bun:split

CREATE INDEX edits_per_author ON note_edits (author_id);
CREATE INDEX edits_per_author_per_note ON note_edits (author_id, public_identifier, target);
