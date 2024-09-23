DROP INDEX IF EXISTS edits_per_author;
DROP INDEX IF EXISTS edits_per_author_per_note;

--bun:split

DROP TABLE IF EXISTS note_edits;

--bun:split

DROP TYPE IF EXISTS note_target;
