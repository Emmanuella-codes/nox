UPDATE personas
SET category = 'fan'
WHERE category = 'patron';

DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  FOR constraint_name IN
    SELECT con.conname
    FROM pg_constraint con
    WHERE con.conrelid = 'personas'::regclass
      AND con.contype = 'c'
      AND pg_get_constraintdef(con.oid) LIKE '%category%'
  LOOP
    EXECUTE format('ALTER TABLE personas DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END LOOP;
END $$;

ALTER TABLE personas
  ALTER COLUMN category SET DEFAULT 'fan',
  ADD CONSTRAINT personas_category_check CHECK (category IN ('fan', 'dj', 'organizer', 'creator'));

UPDATE stories
SET contribution_mode = 'followers'
WHERE contribution_mode = 'private';

DO $$
DECLARE
  constraint_name TEXT;
BEGIN
  FOR constraint_name IN
    SELECT con.conname
    FROM pg_constraint con
    WHERE con.conrelid = 'stories'::regclass
      AND con.contype = 'c'
      AND pg_get_constraintdef(con.oid) LIKE '%contribution_mode%'
  LOOP
    EXECUTE format('ALTER TABLE stories DROP CONSTRAINT IF EXISTS %I', constraint_name);
  END LOOP;
END $$;

ALTER TABLE stories
  ADD CONSTRAINT stories_contribution_mode_check CHECK (contribution_mode IN ('public', 'followers'));
