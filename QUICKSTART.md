# TaskMaster Quick Start Guide

## Prerequisites

Install Go 1.21 or later from https://golang.org/dl/

## Build TaskMaster

```bash
# Build the CLI tool
go build -o taskmaster taskmaster.go

# On Windows
go build -o taskmaster.exe taskmaster.go
```

## Quick Examples

### 1. Initialize Tasks from PRD (New!)

```bash
./taskmaster init
```

**Output:**
```
Parsing PRD from PRD.md...
Generated 4 phases, 33 tasks, 192 subtasks
✓ Successfully generated tasks.json
```

This command automatically parses the PRD.md file and generates a complete tasks.json with:
- All features organized by phase
- Automatic task ID generation
- Subtasks created from requirements
- Priority and time estimates

You can also specify custom parameters:
```bash
./taskmaster init custom-prd.md my-tasks.json my-project 2.0.0
```

### 2. View Current Project Status

```bash
./taskmaster status
```

**Output:**
```
=== Project: etcd-monitor (v1.0.0) ===

📋 Phase: MVP (Weeks 1-6)
   Status: pending

   ○ [task-1.1] Setup Project Infrastructure - P0
      Status: pending | Est. Hours: 8
      Subtasks: 0/5 completed

   ○ [task-1.2] Implement etcd Client Connection - P0
      Status: pending | Est. Hours: 12
      Subtasks: 0/5 completed

   ...

=== Progress Summary ===
Tasks: 0 completed, 0 in progress, 29 pending, 0 blocked
Overall Progress: 0.0%
Subtask Progress: 0.0%
```

### 3. Start Working on a Task

```bash
./taskmaster update task-1.1 in_progress "Beginning project setup"
```

**Output:**
```
✓ Updated task task-1.1: pending -> in_progress
```

### 4. Complete a Subtask

```bash
./taskmaster update-subtask task-1.1 task-1.1.1 completed
```

**Output:**
```
✓ Updated subtask task-1.1.1: pending -> completed
```

### 5. Assign a Task

```bash
./taskmaster assign task-1.2 "Alice"
```

**Output:**
```
✓ Assigned task task-1.2 to Alice
```

### 6. View Task Details

```bash
./taskmaster details task-1.1
```

**Output:**
```
=== Task Details ===
ID: task-1.1
Title: Setup Project Infrastructure
Description: Initialize project structure, dependencies, and development environment
Phase: MVP (Weeks 1-6)
Priority: P0
Status: in_progress
Estimated Hours: 8
Started At: 2025-10-10T14:30:00Z

Subtasks (5):
  ✓ [task-1.1.1] Initialize Go backend project
  ○ [task-1.1.2] Setup React/TypeScript frontend
  ○ [task-1.1.3] Configure PostgreSQL + TimescaleDB
  ○ [task-1.1.4] Setup Redis for messaging
  ○ [task-1.1.5] Configure Docker Compose for local dev
```

### 7. Check Progress

```bash
./taskmaster progress
```

**Output:**
```json
{
  "blockedTasks": 0,
  "completedSubtasks": 1,
  "completedTasks": 0,
  "inProgressTasks": 1,
  "pendingTasks": 28,
  "subtaskProgress": "0.6%",
  "taskProgress": "0.0%",
  "totalSubtasks": 166,
  "totalTasks": 29
}
```

### 8. View Change History

```bash
./taskmaster history 10
```

**Output:**
```
=== State History ===
Last Updated: 2025-10-10T14:35:22Z

[2025-10-10T14:32:15Z] task-1.1.1: pending -> completed
[2025-10-10T14:31:45Z] task-1.2: pending -> pending
[2025-10-10T14:30:12Z] task-1.1: pending -> in_progress (Beginning project setup)
```

### 9. Mark Task as Blocked

```bash
./taskmaster update task-1.3 blocked "Waiting for etcd cluster setup"
```

## Workflow Example

Here's a typical workflow for managing tasks:

```bash
# 1. Check what needs to be done
./taskmaster status

# 2. Pick a task and start working
./taskmaster update task-1.1 in_progress "Starting infrastructure setup"

# 3. Complete subtasks as you go
./taskmaster update-subtask task-1.1 task-1.1.1 completed
./taskmaster update-subtask task-1.1 task-1.1.2 completed
./taskmaster update-subtask task-1.1 task-1.1.3 completed
./taskmaster update-subtask task-1.1 task-1.1.4 completed
./taskmaster update-subtask task-1.1 task-1.1.5 completed

# 4. Complete the main task
./taskmaster update task-1.1 completed "All infrastructure components ready"

# 5. Check overall progress
./taskmaster progress

# 6. Move to next task
./taskmaster update task-1.2 in_progress
```

## State Management

TaskMaster automatically maintains two files:

1. **tasks.json** - Current state of all tasks
2. **tasks.state.json** - History of all state changes

Both files are automatically updated with each command. The state history provides a complete audit trail of all task transitions.

## Integration with Development

You can integrate TaskMaster into your development workflow:

### Git Hooks

Add to `.git/hooks/pre-commit`:
```bash
#!/bin/bash
# Automatically log task progress
echo "Recent task changes:" > .task-summary
./taskmaster history 5 >> .task-summary
```

### CI/CD Pipeline

```yaml
# Example GitHub Actions workflow
- name: Check Task Progress
  run: |
    ./taskmaster progress > progress.json
    echo "::set-output name=progress::$(cat progress.json)"
```

### Team Collaboration

Export current status for team meetings:
```bash
# Generate status report
./taskmaster status > weekly-status.txt
./taskmaster progress > progress.json
```

## Tips

1. **Update frequently** - Keep task status current for accurate progress tracking
2. **Use notes** - Add context when changing status (especially for blocked tasks)
3. **Review history** - Check history to understand project velocity
4. **Track subtasks** - Complete subtasks as you finish them for better granularity
5. **Assign tasks** - Distribute work across team members

## Next Steps

1. Read the full [PRD](PRD.md) to understand requirements
2. Review [tasks.json](tasks.json) for detailed task breakdown
3. Start with Phase 1 tasks (MVP)
4. Use TaskMaster daily to track your progress
