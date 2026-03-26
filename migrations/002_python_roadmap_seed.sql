-- Insert Python Learning Program
INSERT INTO programs (name, description, is_active) 
VALUES (
    'Python: Zero to Senior Developer',
    'Complete Python learning roadmap from beginner to senior developer level',
    true
) ON CONFLICT DO NOTHING;

-- Get the program ID (we'll use a subquery in each insert)
-- Phase 0: Foundations of Programming
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 0: Foundations of Programming', 'Beginner Mindset - What is programming and setting up your environment', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

-- Phase 0 Subtopics
INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'What is Programming?', 'How computers execute code, Interpreted vs compiled languages, Why Python is powerful', true
FROM topics WHERE name = 'Phase 0: Foundations of Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Setting Up Environment', 'Installing Python, Using VS Code/PyCharm, Running scripts, Virtual environments', true
FROM topics WHERE name = 'Phase 0: Foundations of Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'First Program', 'print(), Comments, Basic syntax rules', true
FROM topics WHERE name = 'Phase 0: Foundations of Programming' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 1: Core Python Basics
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 1: Core Python Basics', 'Variables, operators, input/output, conditionals, loops', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Variables & Data Types', 'Integers, Floats, Strings, Booleans, Type casting, Dynamic typing', true
FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Operators', 'Arithmetic, Comparison, Logical, Assignment operators', true
FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Input/Output', 'input() function, Formatting strings (f-strings)', true
FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Conditional Statements', 'if, elif, else, Nested conditions', true
FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Loops', 'for, while, break, continue, pass', true
FROM topics WHERE name = 'Phase 1: Core Python Basics' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 2: Data Structures
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 2: Data Structures', 'Core data structures - Lists, Tuples, Sets, Dictionaries, Strings', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Strings Deep Dive', 'Indexing, slicing, String methods', true
FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Lists', 'CRUD operations, List slicing, List comprehensions', true
FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Tuples', 'Immutability, Packing/unpacking', true
FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Sets', 'Unique elements, Set operations', true
FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Dictionaries', 'Key-value structure, Nested dictionaries', true
FROM topics WHERE name = 'Phase 2: Data Structures' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 3: Functions & Modular Code
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 3: Functions & Modular Code', 'Functions, Lambda, Recursion, Modules & Packages', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Functions Basics', 'Defining & calling functions, Arguments, Keyword arguments', true
FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Advanced Arguments', '*args, **kwargs, Default parameters', true
FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Return Values', 'Multiple returns, Returning functions', true
FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Lambda Functions', 'Anonymous functions, Using lambda in map/filter', true
FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Recursion', 'Base cases, Call stack, Optimization', true
FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Modules & Packages', 'Importing modules, Creating your own modules', true
FROM topics WHERE name = 'Phase 3: Functions & Modular Code' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 4: Object-Oriented Programming (OOP)
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 4: Object-Oriented Programming', 'Classes, Inheritance, Polymorphism, Encapsulation', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Classes & Objects', 'Attributes & methods, Creating instances', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Constructors', '__init__ method, Initialization', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Encapsulation', 'Private variables, Getters/Setters', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Inheritance', 'Single & multiple inheritance, super()', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Polymorphism', 'Method overriding, Duck typing', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Magic Methods', '__str__, __repr__, __eq__, __lt__', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Abstract Classes', 'Interface design, ABC module', true
FROM topics WHERE name = 'Phase 4: Object-Oriented Programming' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 5: Error Handling & Debugging
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 5: Error Handling & Debugging', 'Exceptions, Debugging techniques, Custom exceptions', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Exceptions', 'try, except, finally, else blocks', true
FROM topics WHERE name = 'Phase 5: Error Handling & Debugging' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Custom Exceptions', 'Creating custom exception classes', true
FROM topics WHERE name = 'Phase 5: Error Handling & Debugging' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Debugging Techniques', 'Using breakpoints, Reading stack traces, Logging', true
FROM topics WHERE name = 'Phase 5: Error Handling & Debugging' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 6: File Handling & Data Processing
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 6: File Handling & Data Processing', 'Files, JSON, CSV, Logging', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'File Operations', 'Read/write files, File modes, Context managers', true
FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Working with JSON', 'json.load, json.dump, Parsing JSON', true
FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'CSV Handling', 'csv module, pandas basics', true
FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Logging', 'logging module, Log levels, Handlers', true
FROM topics WHERE name = 'Phase 6: File Handling & Data Processing' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 7: Advanced Python Concepts
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 7: Advanced Python Concepts', 'Iterators, Generators, Decorators, Context Managers', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Iterators & Generators', 'yield, generator functions, Iterator protocol', true
FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Decorators', 'Function wrapping, @decorator syntax, Parameterized decorators', true
FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Context Managers', 'with statement, __enter__ and __exit__', true
FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Closures', 'Nested functions, Variable capture', true
FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Memory Management', 'Garbage collection, Reference counting', true
FROM topics WHERE name = 'Phase 7: Advanced Python Concepts' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 8: Working with APIs & Networking
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 8: Working with APIs & Networking', 'HTTP, REST, requests library, Authentication', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'HTTP Basics', 'GET, POST, PUT, DELETE, Status codes', true
FROM topics WHERE name = 'Phase 8: Working with APIs & Networking' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Using APIs', 'requests library, JSON responses, Error handling', true
FROM topics WHERE name = 'Phase 8: Working with APIs & Networking' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'REST Concepts', 'RESTful principles, HTTP methods, Status codes', true
FROM topics WHERE name = 'Phase 8: Working with APIs & Networking' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Authentication', 'API keys, JWT basics, OAuth concepts', true
FROM topics WHERE name = 'Phase 8: Working with APIs & Networking' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 9: Concurrency & Performance
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 9: Concurrency & Performance', 'Threading, Multiprocessing, Async, Optimization', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Multithreading', 'Thread class, Thread synchronization, Locks', true
FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Multiprocessing', 'Process class, Process pools, IPC', true
FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Async Programming', 'async/await, asyncio, Event loops', true
FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Performance Optimization', 'Profiling, Caching, Optimization techniques', true
FROM topics WHERE name = 'Phase 9: Concurrency & Performance' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 10: Databases
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 10: Databases', 'SQL, PostgreSQL, ORM, Transactions', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'SQL Basics', 'SELECT, INSERT, UPDATE, DELETE, Joins', true
FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'PostgreSQL Integration', 'psycopg2, Connection strings, Query execution', true
FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'ORM', 'SQLAlchemy, Django ORM, Model definitions', true
FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Transactions', 'ACID properties, Commit/Rollback', true
FROM topics WHERE name = 'Phase 10: Databases' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 11: Testing & Code Quality
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 11: Testing & Code Quality', 'Unit Testing, Mocking, Coverage, Linting', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Unit Testing', 'unittest, pytest, Writing test cases', true
FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Mocking', 'unittest.mock, Mocking external dependencies', true
FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Code Coverage', 'Coverage.py, Measuring code coverage', true
FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Linting & Formatting', 'flake8, black, Code style', true
FROM topics WHERE name = 'Phase 11: Testing & Code Quality' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 12: Backend Development
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 12: Backend Development', 'Web frameworks, APIs, Authentication, Middleware', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Web Basics', 'HTTP lifecycle, Request/Response, Headers', true
FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Frameworks', 'Flask, FastAPI, Django basics', true
FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Building APIs', 'RESTful APIs, Request validation, Response formatting', true
FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Authentication Systems', 'JWT, Sessions, OAuth integration', true
FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Middleware', 'Custom middleware, CORS, Rate limiting', true
FROM topics WHERE name = 'Phase 12: Backend Development' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 13: System Design
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 13: System Design with Python', 'Architecture, Design Patterns, Scaling, Caching', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Project Structure', 'Folder organization, Modularity, Configuration', true
FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Design Patterns', 'Singleton, Factory, Observer, Strategy', true
FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Scaling Applications', 'Horizontal scaling, Load balancing, Microservices', true
FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Caching', 'Redis, Memcached, Caching strategies', true
FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Message Queues', 'Kafka, RabbitMQ, Event-driven architecture', true
FROM topics WHERE name = 'Phase 13: System Design with Python' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 14: DevOps & Deployment
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 14: DevOps & Deployment', 'Git, Docker, CI/CD, Cloud Deployment', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Git & GitHub', 'Version control, Branching, Pull requests', true
FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Docker', 'Containers, Images, Docker Compose', true
FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'CI/CD', 'GitHub Actions, Jenkins, Automated testing', true
FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Cloud Deployment', 'AWS, GCP, Heroku, Render basics', true
FROM topics WHERE name = 'Phase 14: DevOps & Deployment' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 15: Specialization Tracks
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 15: Specialization Tracks', 'Choose your path - Backend, Data Engineering, AI/ML, or Automation', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Backend Engineer Track', 'FastAPI mastery, Microservices, Distributed systems, Rate limiting', true
FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Data Engineering Track', 'Pandas, NumPy, ETL pipelines, Airflow, Data warehousing', true
FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'AI/ML Track', 'Scikit-learn, TensorFlow, Deep learning, Model deployment', true
FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Automation/Scripting Track', 'Web scraping, Task automation, Scheduling, System administration', true
FROM topics WHERE name = 'Phase 15: Specialization Tracks' LIMIT 1
ON CONFLICT DO NOTHING;

-- Phase 16: Senior-Level Skills
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Phase 16: Senior-Level Skills', 'Architecture, System Design, Code Reviews, Mentoring', true
FROM programs WHERE name = 'Python: Zero to Senior Developer'
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Code Architecture', 'Clean code, SOLID principles, Design patterns mastery', true
FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'System Design Interviews', 'Designing scalable systems, Trade-offs, Architecture decisions', true
FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Open Source Contribution', 'Contributing to open source, Building communities', true
FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Mentoring & Code Reviews', 'Effective code reviews, Mentoring junior developers', true
FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1
ON CONFLICT DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active)
SELECT id, 'Building Production Systems', 'Reliability, Performance, Monitoring, Incident response', true
FROM topics WHERE name = 'Phase 16: Senior-Level Skills' LIMIT 1
ON CONFLICT DO NOTHING;
