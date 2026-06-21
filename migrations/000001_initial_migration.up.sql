CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    disabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE, -- SHA-256 hash of the token sent to the user (WE USE IT TO FIND THE ROW)
    is_revoked BOOLEAN DEFAULT FALSE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE week_days AS ENUM ('Mon', 'Tue', 'Wed', 'Thu', 'Fri');

CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    slots INT NOT NULL,
    days week_days[] DEFAULT '{}',
    rooms_version INT DEFAULT 1,
    classes_version INT DEFAULT 1,
    teachers_version INT DEFAULT 1,
    subjects_version INT DEFAULT 1,
    assignments_version INT DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    is_lab BOOLEAN DEFAULT FALSE,
    constraints JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    constraints JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS teachers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    constraints JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS subjects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    constraints JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS teacher_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    teacher_id UUID NOT NULL REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id UUID NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
    target_room_id UUID REFERENCES rooms(id) ON DELETE CASCADE,
    constraints JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    is_public BOOLEAN DEFAULT FALSE,
    is_live BOOLEAN DEFAULT FALSE,
    violations JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS timetable_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version_id UUID REFERENCES versions(id) ON DELETE CASCADE,
    slot INT NOT NULL,
    day week_days NOT NULL,
    assignment_id UUID REFERENCES teacher_assignments(id),
    -- add properties (snap shots details)
    -- Add the details needed for the entry showing
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';


DO $$
DECLARE
    t_name text;
BEGIN
    FOR t_name IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') 
    LOOP
        EXECUTE format('
            CREATE TRIGGER update_%I_modtime 
            BEFORE UPDATE ON %I 
            FOR EACH ROW EXECUTE FUNCTION update_modified_column();
        ', t_name, t_name);
    END LOOP;
END $$;


-- For Rooms
CREATE OR REPLACE FUNCTION bump_rooms_version() RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects SET rooms_version = rooms_version + 1 
    WHERE id = COALESCE(NEW.project_id, OLD.project_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- For Classes
CREATE OR REPLACE FUNCTION bump_classes_version() RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects SET classes_version = classes_version + 1 
    WHERE id = COALESCE(NEW.project_id, OLD.project_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- For Teachers
CREATE OR REPLACE FUNCTION bump_teachers_version() RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects SET teachers_version = teachers_version + 1 
    WHERE id = COALESCE(NEW.project_id, OLD.project_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- For Subjects
CREATE OR REPLACE FUNCTION bump_subjects_version() RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects SET subjects_version = subjects_version + 1 
    WHERE id = COALESCE(NEW.project_id, OLD.project_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;


-- For Assignments
CREATE OR REPLACE FUNCTION bump_assignments_version() RETURNS TRIGGER AS $$
BEGIN
    UPDATE projects SET assignments_version = assignments_version + 1 
    WHERE id = COALESCE(NEW.project_id, OLD.project_id);
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;



CREATE TRIGGER trigger_bump_rooms
AFTER INSERT OR UPDATE OR DELETE ON rooms
FOR EACH ROW EXECUTE FUNCTION bump_rooms_version();

CREATE TRIGGER trigger_bump_classes
AFTER INSERT OR UPDATE OR DELETE ON classes
FOR EACH ROW EXECUTE FUNCTION bump_classes_version();

CREATE TRIGGER trigger_bump_teachers
AFTER INSERT OR UPDATE OR DELETE ON teachers
FOR EACH ROW EXECUTE FUNCTION bump_teachers_version();

CREATE TRIGGER trigger_bump_subjects
AFTER INSERT OR UPDATE OR DELETE ON subjects
FOR EACH ROW EXECUTE FUNCTION bump_subjects_version();

CREATE TRIGGER trigger_bump_assignments
AFTER INSERT OR UPDATE OR DELETE ON teacher_assignments
FOR EACH ROW EXECUTE FUNCTION bump_assignments_version();


CREATE INDEX IF NOT EXISTS idx_user_tokens_user_id ON user_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_user_id ON projects(user_id);

CREATE INDEX IF NOT EXISTS idx_rooms_project_id ON rooms(project_id);
CREATE INDEX IF NOT EXISTS idx_classes_project_id ON classes(project_id);
CREATE INDEX IF NOT EXISTS idx_teachers_project_id ON teachers(project_id);
CREATE INDEX IF NOT EXISTS idx_subjects_project_id ON subjects(project_id);

CREATE INDEX IF NOT EXISTS idx_assignments_project_id ON teacher_assignments(project_id);
CREATE INDEX IF NOT EXISTS idx_assignments_class_id ON teacher_assignments(class_id);
CREATE INDEX IF NOT EXISTS idx_assignments_teacher_id ON teacher_assignments(teacher_id);
CREATE INDEX IF NOT EXISTS idx_assignments_subject_id ON teacher_assignments(subject_id);

CREATE INDEX IF NOT EXISTS idx_versions_project_id ON versions(project_id);
CREATE INDEX IF NOT EXISTS idx_entries_version_id ON timetable_entries(version_id);
CREATE INDEX IF NOT EXISTS idx_entries_assignment_id ON timetable_entries(assignment_id);