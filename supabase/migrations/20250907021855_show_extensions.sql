select
  name, comment, default_version, installed_version
from
  pg_available_extensions
order by
  name asc;
