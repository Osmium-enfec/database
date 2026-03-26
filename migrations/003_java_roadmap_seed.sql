-- Java Programming Roadmap
-- Organized by topics and subtopics

INSERT INTO programs (name, description, is_active) 
VALUES ('Java Programming', 'Comprehensive Java learning from basics to advanced Java 8 features', true)
ON CONFLICT (name) DO NOTHING;

-- ============================================================================
-- TOPIC 1: Core Concepts
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Core Concepts', 'Java fundamentals, data types, operators, and control flow', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Basics & Installation', 'Intro, Installation, program execution flow, types of code editors, Compilation and Execution', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Variables and Data Types', 'Primitives (byte, short, int, long, double, float, char, boolean) and Non-Primitives (arrays, Strings, classes)', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Identifiers, Rules & Keywords', 'Identifier rules and Java keywords', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Comments and Naming Conventions', 'Code documentation and naming best practices', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Scanner Class', 'Reading input from keyboard using Scanner class', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Operators', 'Arithmetic, logical, relational, new, assignment, and instanceof operators', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Control Statements', 'if, if-else, if-else ladder, switch, while, do-while, for, for-each loops, and break/continue/return', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Arrays & Its Types', '1D and 2D arrays', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Strings', 'String creation, manipulation methods, StringBuffer, StringBuilder, StringTokenizer', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Type Casting', 'Widening/Narrowing with primitives, UpCasting/DownCasting with non-primitives', true),
((SELECT id FROM topics WHERE name = 'Core Concepts' LIMIT 1), 'Wrapper Classes', 'Wrapper class uses and importance', true);

-- ============================================================================
-- TOPIC 2: Object-Oriented Programming
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Object-Oriented Programming', 'Classes, objects, inheritance, polymorphism, encapsulation, abstraction', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Classes & Objects', 'Class definition and object creation', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Packages & Types', 'Package organization and types', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Inheritance', 'Single, multiple, multilevel, and hierarchical inheritance', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'SUPER Keyword', 'Using super keyword for parent class access', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Encapsulation', 'Data hiding and access control', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Polymorphism', 'Method overloading and overriding', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Abstraction', 'Abstract concept implementation', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Access Modifiers', 'public, private, protected, default access levels', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'FINAL Keyword', 'Using final keyword for classes, methods, and variables', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Abstract Classes', 'Abstract class definition and usage', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Object Class', 'Object class methods and functionality', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Interfaces', 'Interface definition, implementation, abstract methods, marker interfaces', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Types of Variables', 'Instance, local, and static variables', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Methods', 'Method combinations and types', true),
((SELECT id FROM topics WHERE name = 'Object-Oriented Programming' LIMIT 1), 'Constructors', 'Constructor types and THIS keyword usage', true);

-- ============================================================================
-- TOPIC 3: Exception Handling
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Exception Handling', 'Checked/unchecked exceptions, try-catch-finally, custom exceptions', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Exception Types', 'Checked and unchecked exceptions', true),
((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Try-Catch-Finally', 'try, catch, finally, throw, throws keywords', true),
((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Try with Resources', 'Automatic resource management', true),
((SELECT id FROM topics WHERE name = 'Exception Handling' LIMIT 1), 'Custom Exceptions', 'Creating and using user-defined exceptions', true);

-- ============================================================================
-- TOPIC 4: Collections Framework
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Collections Framework', 'List, Set, Map, Queues, and related interfaces', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Cursors & Interfaces', 'Iterators, Comparable, Comparator interfaces', true),
((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'List Interface', 'ArrayList, LinkedList, Vector, and Stack', true),
((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Queue Interface', 'Queue implementations and operations', true),
((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Set Interface', 'HashSet, LinkedHashSet, SortedSet, TreeSet', true),
((SELECT id FROM topics WHERE name = 'Collections Framework' LIMIT 1), 'Map Interface', 'HashMap, LinkedHashMap, TreeMap, Hashtable, IdentityHashMap, WeakHashMap', true);

-- ============================================================================
-- TOPIC 5: Multi-Threading
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Multi-Threading', 'Threads, synchronization, concurrent programming', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Thread Basics', 'Default threads and user-defined threads', true),
((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Thread Life Cycle', 'Thread states and transitions', true),
((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Callable & Executor', 'Callable interface and ExecutorService', true),
((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Daemon Threads', 'Daemon thread creation and usage', true),
((SELECT id FROM topics WHERE name = 'Multi-Threading' LIMIT 1), 'Synchronization', 'Synchronization techniques and methods', true);

-- ============================================================================
-- TOPIC 6: File IO & Serialization
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'File IO & Serialization', 'File operations, streams, serialization/deserialization', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'File IO & Serialization' LIMIT 1), 'File Handling', 'Create, write, read, delete file operations', true),
((SELECT id FROM topics WHERE name = 'File IO & Serialization' LIMIT 1), 'IO Streams', 'FileWriter, FileReader, and stream operations', true),
((SELECT id FROM topics WHERE name = 'File IO & Serialization' LIMIT 1), 'Serialization', 'SerialVersionUID, transient keyword, serialization/deserialization', true);

-- ============================================================================
-- TOPIC 7: Generics
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Generics', 'Generic types, wildcards, type parameters, bounded types', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Generics' LIMIT 1), 'Generic Basics', 'Generic type parameters and type safety', true),
((SELECT id FROM topics WHERE name = 'Generics' LIMIT 1), 'Wildcards', 'Wildcard types and bounded types', true),
((SELECT id FROM topics WHERE name = 'Generics' LIMIT 1), 'Generic Methods', 'Generic method definition and usage', true);

-- ============================================================================
-- TOPIC 8: Java 8 Features
-- ============================================================================
INSERT INTO topics (program_id, name, description, is_active)
SELECT id, 'Java 8 Features', 'Lambda expressions, functional interfaces, streams, and new APIs', true
FROM programs WHERE name = 'Java Programming' LIMIT 1
ON CONFLICT (program_id, name) DO NOTHING;

INSERT INTO subtopics (topic_id, name, description, is_active) VALUES
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Interface Changes', 'default and static methods in interfaces', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Lambda Expressions', 'Lambda syntax and functional programming', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Functional Interfaces', 'Consumer, Supplier, Predicate, Function', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Stream API', 'Creation, filtering, mapping, slicing, matching, finding, collectors', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Parallel Streams', 'Parallel stream processing', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Date & Time API', 'New date/time API changes', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Optional Class', 'Optional class and its methods', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Method References', 'Method and constructor references', true),
((SELECT id FROM topics WHERE name = 'Java 8 Features' LIMIT 1), 'Utilities', 'forEach() method, Spliterator, StringJoiner', true);
