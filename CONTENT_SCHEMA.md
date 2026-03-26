# Content Creation Schema - UI Input Guide

## Overview
There are 3 types of content you can create:
1. **Question** (MCQ, Multiple Select, Fill in the blank, Short answer)
2. **Code Problem** (Coding exercise)
3. **Documentation** (Markdown content)

---

## 1. CREATE QUESTION (Quiz Type)

### Request to Backend
```json
{
  "type": "question",
  "program_id": "550e8400-e29b-41d4-a716-446655440000",
  "topic_id": "550e8400-e29b-41d4-a716-446655440001",
  "subtopic_id": "550e8400-e29b-41d4-a716-446655440002",
  "difficulty": "medium",
  "estimated_time_minutes": 5,
  "tags": ["algebra", "linear-equations", "beginner"],
  "data": {
    "type": "question",
    "title": "Solve Linear Equation",
    "description": "Find the value of x in the given equation",
    "question_type": "mcq",
    "question_text": "What is the value of x if 2x + 5 = 15?",
    "hints": [
      "Subtract 5 from both sides of the equation",
      "After subtraction you should have 2x = 10",
      "Divide both sides by 2 to isolate x"
    ],
    "options": [
      "3",
      "5",
      "7",
      "10"
    ],
    "correct_options": [1]
  }
}
```

### Question Types & Structure

#### MCQ (Multiple Choice - Single Answer)
```json
{
  "question_type": "mcq",
  "question_text": "What is 2+2?",
  "hints": ["Think about basic arithmetic", "Count on your fingers", "The answer is between 3 and 5"],
  "options": ["3", "4", "5", "6"],
  "correct_options": [1]
}
```

#### MSQ (Multiple Select - Multiple Answers)
```json
{
  "question_type": "msq",
  "question_text": "Which of these are prime numbers?",
  "hints": ["Prime numbers are divisible only by 1 and themselves", "Consider 1, 2, 3, 4, 5", "Check each number carefully"],
  "options": ["1", "2", "3", "4", "5"],
  "correct_options": [1, 2, 4]
}
```

#### Fill in the Blank
```json
{
  "question_type": "fill",
  "question_text": "The capital of France is _______",
  "hints": ["It's the most visited city in France", "Home to the Eiffel Tower", "Starts with 'P'"],
  "options": ["Paris", "Lyon", "Marseille"],
  "correct_options": [0]
}
```

#### Short Answer
```json
{
  "question_type": "short",
  "question_text": "Name the largest planet in our solar system",
  "hints": ["It starts with 'J'", "It's a gas giant", "Named after the king of Roman gods"],
  "options": ["Jupiter"],
  "correct_options": [0]
}
```

### Input Fields for Questions

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| type | string | Yes | Must be "question" | "question" |
| program_id | UUID | Yes | Educational program ID | "550e8400-..." |
| topic_id | UUID | Yes | Topic within program | "550e8400-..." |
| subtopic_id | UUID | Yes | Subtopic within topic | "550e8400-..." |
| difficulty | string | Yes | One of: easy, medium, hard | "medium" |
| estimated_time_minutes | number | Yes | Time to complete (1-1440) | 5 |
| tags | array | No | Up to 20 tags | ["algebra", "beginner"] |
| data.title | string | Yes | Question title | "Solve Linear Equation" |
| data.description | string | Yes | Question description | "Find the value of x..." |
| data.question_type | string | Yes | One of: mcq, msq, fill, short | "mcq" |
| data.question_text | string | Yes | The actual question | "What is 2+2?" |
| data.hints | array | Yes | **At least 3 hints** | ["Hint 1", "Hint 2", "Hint 3"] |
| data.options | array | Yes | Answer choices (1-10 options) | ["A", "B", "C", "D"] |
| data.correct_options | array | Yes | Indices of correct answers | [1] or [1, 3] for MSQ |

---

## 2. CREATE CODE PROBLEM

### Request to Backend
```json
{
  "type": "code_problem",
  "program_id": "550e8400-e29b-41d4-a716-446655440000",
  "topic_id": "550e8400-e29b-41d4-a716-446655440001",
  "subtopic_id": "550e8400-e29b-41d4-a716-446655440002",
  "difficulty": "hard",
  "estimated_time_minutes": 30,
  "tags": ["python", "sorting", "algorithms"],
  "data": {
    "type": "code_problem",
    "title": "Implement Bubble Sort",
    "description": "Write a function to sort an array using the bubble sort algorithm",
    "code_problem_data": {
      "starter_code": "def bubble_sort(arr):\n    # Write your code here\n    pass",
      "solution_code": "def bubble_sort(arr):\n    n = len(arr)\n    for i in range(n):\n        for j in range(0, n-i-1):\n            if arr[j] > arr[j+1]:\n                arr[j], arr[j+1] = arr[j+1], arr[j]\n    return arr",
      "execution_template": "result = bubble_sort({input})\nprint(result)",
      "test_cases": [
        {
          "input": "[64, 34, 25, 12, 22, 11, 90]",
          "expected_output": "[11, 12, 22, 25, 34, 64, 90]",
          "is_hidden": false
        },
        {
          "input": "[5, 2, 8, 1, 9]",
          "expected_output": "[1, 2, 5, 8, 9]",
          "is_hidden": false
        },
        {
          "input": "[3, 3, 3, 3]",
          "expected_output": "[3, 3, 3, 3]",
          "is_hidden": true
        }
      ]
    }
  }
}
```

### Input Fields for Code Problems

| Field | Type | Required | Description | Example |
|-------|------|----------|-------------|---------|
| type | string | Yes | Must be "code_problem" | "code_problem" |
| program_id | UUID | Yes | Educational program ID | "550e8400-..." |
| topic_id | UUID | Yes | Topic within program | "550e8400-..." |
| subtopic_id | UUID | Yes | Subtopic within topic | "550e8400-..." |
| difficulty | string | Yes | One of: easy, medium, hard | "hard" |
| estimated_time_minutes | number | Yes | Time to complete (1-1440) | 30 |
| tags | array | No | Up to 20 tags | ["python", "sorting"] |
| data.title | string | Yes | Problem title | "Implement Bubble Sort" |
| data.description | string | Yes | Problem description | "Write a function to..." |
| code_problem_data.starter_code | string | Yes | Initial code template | "def bubble_sort(arr):\n    pass" |
| code_problem_data.solution_code | string | Yes | Complete solution | "def bubble_sort(arr):\n    # full code" |
| code_problem_data.execution_template | string | Yes | How to run (use {input}) | "result = bubble_sort({input})\nprint(result)" |
| code_problem_data.test_cases | array | Yes | At least 1 test case | [...] |
| test_case.input | string | Yes | Input for test case | "[1, 5, 3]" |
| test_case.expected_output | string | Yes | Expected output | "[1, 3, 5]" |
| test_case.is_hidden | boolean | Yes | Hide from user or show | false |

---

## 3. CREATE DOCUMENTATION

### Request to Backend
```json
{
  "type": "documentation",
  "program_id": "550e8400-e29b-41d4-a716-446655440000",
  "topic_id": "550e8400-e29b-41d4-a716-446655440001",
  "subtopic_id": "550e8400-e29b-41d4-a716-446655440002",
  "difficulty": "easy",
  "estimated_time_minutes": 10,
  "tags": ["markdown", "guide", "reference"],
  "data": {
    "type": "documentation",
    "title": "Python Lists - Complete Guide",
    "description": "Learn everything about Python lists",
    "documentation_data": {
      "markdown_content": "# Python Lists\n\n## Introduction\nLists are ordered collections of items...\n\n## Creating Lists\n```python\nmy_list = [1, 2, 3]\n```\n\n## Methods\n- `append()` - Add item\n- `remove()` - Remove item\n\n## Examples\n1. Basic list creation\n2. List operations\n"
    }
  }
}
```

### Input Fields for Documentation

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| type | string | Yes | Must be "documentation" |
| program_id | UUID | Yes | Educational program ID |
| topic_id | UUID | Yes | Topic within program |
| subtopic_id | UUID | Yes | Subtopic within topic |
| difficulty | string | Yes | One of: easy, medium, hard |
| estimated_time_minutes | number | Yes | Time to read (1-1440) |
| tags | array | No | Up to 20 tags |
| data.title | string | Yes | Documentation title |
| data.description | string | Yes | Short description |
| markdown_content | string | Yes | Full markdown content |

---

## Key Points for UI Development

### Questions (Quiz)
✅ **Required**: At least 3 hints per question
✅ **Validation**: Hints array length >= 3
✅ **Question Types**: MCQ, MSQ, Fill in blank, Short answer
✅ **Correct Options**: Array of indices (0-based) that are correct

### Code Problems
✅ **Visible Tests**: Show to user (is_hidden: false)
✅ **Hidden Tests**: Use for grading only (is_hidden: true)
✅ **Execution**: Use {input} placeholder in template
✅ **At least 1 visible test case required**

### Common Fields
✅ **Difficulty**: easy, medium, hard
✅ **Time**: 1-1440 minutes
✅ **Tags**: Max 20 tags per content
✅ **All content requires**: program_id, topic_id, subtopic_id

---

## Full Create Content Request Structure

```json
{
  "type": "question | code_problem | documentation",
  "program_id": "uuid",
  "topic_id": "uuid",
  "subtopic_id": "uuid",
  "difficulty": "easy|medium|hard",
  "estimated_time_minutes": 5,
  "tags": ["tag1", "tag2"],
  "data": {
    "type": "question | code_problem | documentation",
    "title": "Content title",
    "description": "Content description",
    "question_type": "mcq|msq|fill|short",
    "question_text": "The question/problem text",
    "hints": ["hint1", "hint2", "hint3"],
    "options": ["option1", "option2"],
    "correct_options": [0, 1],
    "code_problem_data": {
      "starter_code": "code",
      "solution_code": "code",
      "execution_template": "template",
      "test_cases": [{"input": "", "expected_output": "", "is_hidden": false}]
    },
    "documentation_data": {
      "markdown_content": "# Markdown"
    }
  }
}
```

---

## Next Steps for UI

1. **Question Form**: Create form with MCQ/MSQ/Fill/Short tabs
2. **Hints Section**: Add field for 3+ hints (expandable)
3. **Options Input**: Dynamic field to add/remove options
4. **Code Problem**: Code editor fields + test case manager
5. **Documentation**: Markdown editor with preview
6. **Validation**: Enforce 3 hints minimum for questions

