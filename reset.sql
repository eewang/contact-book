-- Remove foreign keys first
-- ALTER TABLE students DROP FOREIGN KEY students_ibfk_1;
-- ALTER TABLE classes DROP FOREIGN KEY classes_ibfk_1;

-- Then drop the tables
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS groups;
