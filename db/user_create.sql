CREATE USER 'snippetbox'@'%';
GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'snippetbox'@'%';
-- Important: Make sure to swap 'pass' with a password of your own choosing.
ALTER USER 'snippetbox'@'%' IDENTIFIED  BY 'root';