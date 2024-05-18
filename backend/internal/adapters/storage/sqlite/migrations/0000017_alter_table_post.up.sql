ALTER TABLE post ADD COLUMN privacy_default TEXT CHECK (privacy_default IN ('public', 'private', 'almost private')) DEFAULT 'public';

UPDATE post SET privacy_default = privacy;

ALTER TABLE post DROP COLUMN privacy;

ALTER TABLE post RENAME COLUMN privacy_default TO privacy;
