-- Return a bool if enum exists
CREATE OR REPLACE FUNCTION enum_exists(enum_name text, enum_schema text DEFAULT 'public')
RETURNS boolean
LANGUAGE sql
AS $$
    SELECT EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_namespace n ON t.typnamespace = n.oid
        WHERE t.typname = enum_name
          AND n.nspname = enum_schema
    );
$$;


-- Create Roles Enum
DO $$
BEGIN
IF NOT enum_exists('app_roles') THEN
  CREATE TYPE public.app_roles AS ENUM ('user', 'admin', 'super');
  -- ELSE
    -- ALTER TYPE public.app_roles RENAME VALUE 'super' TO 'super-user';
    -- ALTER TYPE app_role ADD VALUE 'people-places-things';

    -- ALTER TYPE app_role RENAME VALUE 'pending' TO 'waiting';

    -- DROP TYPE app_roles;
END IF;
END $$;

-- Create Permission Enum
-- Permissions style 'action:resource:[optional]modifier', should just be strings tho ...
-- | =================================== |
-- | Modifider in use:                   |
-- | self    : resources owned by a user |
-- | friends : friends can see this data |
-- | global  : open to the public        |
-- | =================================== |
DO $$
BEGIN
IF NOT enum_exists('app_permissions') THEN
    CREATE TYPE public.app_permissions AS ENUM (
      'create:post:friends' -- (?)
      ,'create:post:global' -- (?)
      ,'read:post'
      ,'update:post:self'
      ,'delete:post:self'
    );
  -- ElSE
  -- DROP TYPE public.app_permissions;
END IF;
END $$;
