-- Droping tables in reverse order of dependencies
DROP TABLE IF EXISTS timetable_entries CASCADE;
DROP TABLE IF EXISTS versions CASCADE;
DROP TABLE IF EXISTS teacher_assignments CASCADE;
DROP TABLE IF EXISTS subjects CASCADE;
DROP TABLE IF EXISTS teachers CASCADE;
DROP TABLE IF EXISTS classes CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS projects CASCADE;
DROP TABLE IF EXISTS user_tokens CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Droping custom trigger functions
DROP FUNCTION IF EXISTS bump_assignments_version();
DROP FUNCTION IF EXISTS bump_subjects_version();
DROP FUNCTION IF EXISTS bump_teachers_version();
DROP FUNCTION IF EXISTS bump_classes_version();
DROP FUNCTION IF EXISTS bump_rooms_version();
DROP FUNCTION IF EXISTS update_modified_column() CASCADE;

-- Drop enums
DROP TYPE IF EXISTS week_days;