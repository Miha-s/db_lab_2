
CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    telephone VARCHAR(20),
    employment_date DATE NOT NULL,
    firing_date DATE,
    current_position VARCHAR(100),
    current_salary DECIMAL(10, 2)
);

alter table employees
    owner to postgres;

CREATE TABLE employee_skills (
    employee_id INT REFERENCES employees(id),
    skill VARCHAR(100),
    PRIMARY KEY (employee_id, skill)
);

alter table employee_skills
    owner to postgres;

CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    start_date DATE,
    end_date DATE,
    CONSTRAINT valid_dates CHECK (start_date <= end_date)
);

alter table projects
    owner to postgres;

CREATE TABLE employee_projects (
    employee_id INT REFERENCES employees(id),
    project_id INT REFERENCES projects(project_id),
    PRIMARY KEY (employee_id, project_id)
);


alter table employee_projects
    owner to postgres;
