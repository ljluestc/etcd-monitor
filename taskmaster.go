package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

// Task status constants
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusBlocked    = "blocked"
)

// Priority constants
const (
	PriorityP0 = "P0"
	PriorityP1 = "P1"
	PriorityP2 = "P2"
)

// Subtask represents a subtask within a task
type Subtask struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

// Task represents a single task
type Task struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Priority        string    `json:"priority"`
	EstimatedHours  int       `json:"estimatedHours"`
	Status          string    `json:"status"`
	Assignee        *string   `json:"assignee"`
	Subtasks        []Subtask `json:"subtasks"`
	StartedAt       *string   `json:"startedAt,omitempty"`
	CompletedAt     *string   `json:"completedAt,omitempty"`
	ActualHours     *int      `json:"actualHours,omitempty"`
	BlockedReason   *string   `json:"blockedReason,omitempty"`
}

// Phase represents a project phase
type Phase struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Timeline string `json:"timeline"`
	Status   string `json:"status"`
	Tasks    []Task `json:"tasks"`
}

// Summary contains project statistics
type Summary struct {
	TotalPhases       int            `json:"totalPhases"`
	TotalTasks        int            `json:"totalTasks"`
	TotalSubtasks     int            `json:"totalSubtasks"`
	EstimatedTotalHours int          `json:"estimatedTotalHours"`
	Priorities        map[string]int `json:"priorities"`
}

// Project represents the entire project structure
type Project struct {
	ProjectName string  `json:"project"`
	Version     string  `json:"version"`
	Phases      []Phase `json:"phases"`
	Summary     Summary `json:"summary"`
}

// TaskMaster manages task state and operations
type TaskMaster struct {
	project      *Project
	filePath     string
	stateFilePath string
}

// StateHistory tracks state changes
type StateHistory struct {
	Timestamp string `json:"timestamp"`
	TaskID    string `json:"taskId"`
	OldStatus string `json:"oldStatus"`
	NewStatus string `json:"newStatus"`
	Notes     string `json:"notes,omitempty"`
}

// State tracks overall project state
type State struct {
	LastUpdated string         `json:"lastUpdated"`
	History     []StateHistory `json:"history"`
}

// NewTaskMaster creates a new TaskMaster instance
func NewTaskMaster(filePath string) (*TaskMaster, error) {
	tm := &TaskMaster{
		filePath:      filePath,
		stateFilePath: strings.TrimSuffix(filePath, ".json") + ".state.json",
	}

	if err := tm.load(); err != nil {
		return nil, err
	}

	return tm, nil
}

// load reads the tasks from the JSON file
func (tm *TaskMaster) load() error {
	data, err := ioutil.ReadFile(tm.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	tm.project = &Project{}
	if err := json.Unmarshal(data, tm.project); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	return nil
}

// save writes the tasks back to the JSON file
func (tm *TaskMaster) save() error {
	data, err := json.MarshalIndent(tm.project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := ioutil.WriteFile(tm.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// loadState loads the state history
func (tm *TaskMaster) loadState() (*State, error) {
	state := &State{
		History: make([]StateHistory, 0),
	}

	data, err := ioutil.ReadFile(tm.stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return state, nil
		}
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	if err := json.Unmarshal(data, state); err != nil {
		return nil, fmt.Errorf("failed to parse state JSON: %w", err)
	}

	return state, nil
}

// saveState saves the state history
func (tm *TaskMaster) saveState(state *State) error {
	state.LastUpdated = time.Now().Format(time.RFC3339)

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state JSON: %w", err)
	}

	if err := ioutil.WriteFile(tm.stateFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// findTask finds a task by ID across all phases
func (tm *TaskMaster) findTask(taskID string) (*Task, *Phase, error) {
	for i := range tm.project.Phases {
		for j := range tm.project.Phases[i].Tasks {
			if tm.project.Phases[i].Tasks[j].ID == taskID {
				return &tm.project.Phases[i].Tasks[j], &tm.project.Phases[i], nil
			}
		}
	}
	return nil, nil, fmt.Errorf("task not found: %s", taskID)
}

// UpdateTaskStatus updates the status of a task
func (tm *TaskMaster) UpdateTaskStatus(taskID, newStatus, notes string) error {
	task, _, err := tm.findTask(taskID)
	if err != nil {
		return err
	}

	oldStatus := task.Status

	// Load state
	state, err := tm.loadState()
	if err != nil {
		return err
	}

	// Update task status
	task.Status = newStatus

	// Track timestamps
	now := time.Now().Format(time.RFC3339)
	if newStatus == StatusInProgress && task.StartedAt == nil {
		task.StartedAt = &now
	}
	if newStatus == StatusCompleted && task.CompletedAt == nil {
		task.CompletedAt = &now
	}

	// Add to history
	state.History = append(state.History, StateHistory{
		Timestamp: now,
		TaskID:    taskID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Notes:     notes,
	})

	// Save both files
	if err := tm.save(); err != nil {
		return err
	}
	if err := tm.saveState(state); err != nil {
		return err
	}

	fmt.Printf("✓ Updated task %s: %s -> %s\n", taskID, oldStatus, newStatus)
	return nil
}

// UpdateSubtaskStatus updates the status of a subtask
func (tm *TaskMaster) UpdateSubtaskStatus(taskID, subtaskID, newStatus string) error {
	task, _, err := tm.findTask(taskID)
	if err != nil {
		return err
	}

	for i := range task.Subtasks {
		if task.Subtasks[i].ID == subtaskID {
			oldStatus := task.Subtasks[i].Status
			task.Subtasks[i].Status = newStatus

			// Load and update state
			state, err := tm.loadState()
			if err != nil {
				return err
			}

			state.History = append(state.History, StateHistory{
				Timestamp: time.Now().Format(time.RFC3339),
				TaskID:    subtaskID,
				OldStatus: oldStatus,
				NewStatus: newStatus,
			})

			if err := tm.save(); err != nil {
				return err
			}
			if err := tm.saveState(state); err != nil {
				return err
			}

			fmt.Printf("✓ Updated subtask %s: %s -> %s\n", subtaskID, oldStatus, newStatus)
			return nil
		}
	}

	return fmt.Errorf("subtask not found: %s", subtaskID)
}

// AssignTask assigns a task to someone
func (tm *TaskMaster) AssignTask(taskID, assignee string) error {
	task, _, err := tm.findTask(taskID)
	if err != nil {
		return err
	}

	task.Assignee = &assignee

	if err := tm.save(); err != nil {
		return err
	}

	fmt.Printf("✓ Assigned task %s to %s\n", taskID, assignee)
	return nil
}

// GetProgress calculates project progress
func (tm *TaskMaster) GetProgress() map[string]interface{} {
	totalTasks := 0
	completedTasks := 0
	inProgressTasks := 0
	blockedTasks := 0
	pendingTasks := 0

	totalSubtasks := 0
	completedSubtasks := 0

	for _, phase := range tm.project.Phases {
		for _, task := range phase.Tasks {
			totalTasks++
			switch task.Status {
			case StatusCompleted:
				completedTasks++
			case StatusInProgress:
				inProgressTasks++
			case StatusBlocked:
				blockedTasks++
			case StatusPending:
				pendingTasks++
			}

			for _, subtask := range task.Subtasks {
				totalSubtasks++
				if subtask.Status == StatusCompleted {
					completedSubtasks++
				}
			}
		}
	}

	taskProgress := 0.0
	if totalTasks > 0 {
		taskProgress = float64(completedTasks) / float64(totalTasks) * 100
	}

	subtaskProgress := 0.0
	if totalSubtasks > 0 {
		subtaskProgress = float64(completedSubtasks) / float64(totalSubtasks) * 100
	}

	return map[string]interface{}{
		"totalTasks":        totalTasks,
		"completedTasks":    completedTasks,
		"inProgressTasks":   inProgressTasks,
		"blockedTasks":      blockedTasks,
		"pendingTasks":      pendingTasks,
		"taskProgress":      fmt.Sprintf("%.1f%%", taskProgress),
		"totalSubtasks":     totalSubtasks,
		"completedSubtasks": completedSubtasks,
		"subtaskProgress":   fmt.Sprintf("%.1f%%", subtaskProgress),
	}
}

// PrintStatus prints the current status of all tasks
func (tm *TaskMaster) PrintStatus() {
	fmt.Printf("\n=== Project: %s (v%s) ===\n\n", tm.project.ProjectName, tm.project.Version)

	for _, phase := range tm.project.Phases {
		fmt.Printf("📋 Phase: %s (%s)\n", phase.Name, phase.Timeline)
		fmt.Printf("   Status: %s\n\n", phase.Status)

		for _, task := range phase.Tasks {
			statusIcon := tm.getStatusIcon(task.Status)
			fmt.Printf("   %s [%s] %s - %s\n", statusIcon, task.ID, task.Title, task.Priority)
			fmt.Printf("      Status: %s | Est. Hours: %d\n", task.Status, task.EstimatedHours)

			if task.Assignee != nil {
				fmt.Printf("      Assignee: %s\n", *task.Assignee)
			}

			if len(task.Subtasks) > 0 {
				completed := 0
				for _, st := range task.Subtasks {
					if st.Status == StatusCompleted {
						completed++
					}
				}
				fmt.Printf("      Subtasks: %d/%d completed\n", completed, len(task.Subtasks))
			}
			fmt.Println()
		}
	}

	// Print progress summary
	progress := tm.GetProgress()
	fmt.Println("=== Progress Summary ===")
	fmt.Printf("Tasks: %v completed, %v in progress, %v pending, %v blocked\n",
		progress["completedTasks"], progress["inProgressTasks"],
		progress["pendingTasks"], progress["blockedTasks"])
	fmt.Printf("Overall Progress: %s\n", progress["taskProgress"])
	fmt.Printf("Subtask Progress: %s\n\n", progress["subtaskProgress"])
}

// PrintTaskDetails prints detailed information about a specific task
func (tm *TaskMaster) PrintTaskDetails(taskID string) error {
	task, phase, err := tm.findTask(taskID)
	if err != nil {
		return err
	}

	fmt.Printf("\n=== Task Details ===\n")
	fmt.Printf("ID: %s\n", task.ID)
	fmt.Printf("Title: %s\n", task.Title)
	fmt.Printf("Description: %s\n", task.Description)
	fmt.Printf("Phase: %s (%s)\n", phase.Name, phase.Timeline)
	fmt.Printf("Priority: %s\n", task.Priority)
	fmt.Printf("Status: %s\n", task.Status)
	fmt.Printf("Estimated Hours: %d\n", task.EstimatedHours)

	if task.Assignee != nil {
		fmt.Printf("Assignee: %s\n", *task.Assignee)
	}

	if task.StartedAt != nil {
		fmt.Printf("Started At: %s\n", *task.StartedAt)
	}

	if task.CompletedAt != nil {
		fmt.Printf("Completed At: %s\n", *task.CompletedAt)
	}

	if len(task.Subtasks) > 0 {
		fmt.Printf("\nSubtasks (%d):\n", len(task.Subtasks))
		for _, st := range task.Subtasks {
			icon := tm.getStatusIcon(st.Status)
			fmt.Printf("  %s [%s] %s\n", icon, st.ID, st.Title)
		}
	}

	fmt.Println()
	return nil
}

// getStatusIcon returns an icon for the status
func (tm *TaskMaster) getStatusIcon(status string) string {
	switch status {
	case StatusCompleted:
		return "✓"
	case StatusInProgress:
		return "▶"
	case StatusBlocked:
		return "⊗"
	case StatusPending:
		return "○"
	default:
		return "?"
	}
}

// PrintHistory prints the state change history
func (tm *TaskMaster) PrintHistory(limit int) error {
	state, err := tm.loadState()
	if err != nil {
		return err
	}

	fmt.Printf("\n=== State History ===\n")
	fmt.Printf("Last Updated: %s\n\n", state.LastUpdated)

	start := len(state.History) - limit
	if start < 0 {
		start = 0
	}

	for i := len(state.History) - 1; i >= start; i-- {
		h := state.History[i]
		fmt.Printf("[%s] %s: %s -> %s", h.Timestamp, h.TaskID, h.OldStatus, h.NewStatus)
		if h.Notes != "" {
			fmt.Printf(" (%s)", h.Notes)
		}
		fmt.Println()
	}
	fmt.Println()

	return nil
}

// PRDFeature represents a parsed feature from the PRD
type PRDFeature struct {
	Number       string
	Title        string
	Priority     string
	Requirements []string
	EstimatedHrs int
}

// PRDPhase represents a phase from the PRD
type PRDPhase struct {
	Number   string
	Name     string
	Timeline string
	Features []string
}

// ParsePRD parses the PRD markdown file and extracts features
func ParsePRD(prdPath string) ([]PRDPhase, map[string]*PRDFeature, error) {
	file, err := os.Open(prdPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open PRD: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	phases := []PRDPhase{}
	features := make(map[string]*PRDFeature)

	var currentFeature *PRDFeature
	var currentPhase *PRDPhase
	inRequirements := false
	inPhases := false

	featureRegex := regexp.MustCompile(`^###\s+(\d+\.\d+)\s+(.+)$`)
	priorityRegex := regexp.MustCompile(`^\*\*Priority:\*\*\s+(P\d)`)
	phaseRegex := regexp.MustCompile(`^###\s+Phase\s+(\d+):\s+(.+?)\s+\((.+?)\)`)
	requirementRegex := regexp.MustCompile(`^\s*-\s+(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Parse phases
		if strings.Contains(line, "## 6. Implementation Phases") {
			inPhases = true
			continue
		}

		if inPhases {
			if matches := phaseRegex.FindStringSubmatch(line); matches != nil {
				if currentPhase != nil {
					phases = append(phases, *currentPhase)
				}
				currentPhase = &PRDPhase{
					Number:   matches[1],
					Name:     matches[2],
					Timeline: matches[3],
					Features: []string{},
				}
			}

			if strings.HasPrefix(line, "## ") && !strings.Contains(line, "Phase") {
				if currentPhase != nil {
					phases = append(phases, *currentPhase)
					currentPhase = nil
				}
				inPhases = false
			}

			if currentPhase != nil && strings.HasPrefix(line, "- ") {
				currentPhase.Features = append(currentPhase.Features, strings.TrimPrefix(line, "- "))
			}
		}

		// Parse features
		if matches := featureRegex.FindStringSubmatch(line); matches != nil {
			if currentFeature != nil {
				features[currentFeature.Number] = currentFeature
			}
			currentFeature = &PRDFeature{
				Number:       matches[1],
				Title:        matches[2],
				Requirements: []string{},
				EstimatedHrs: 20, // Default
			}
			inRequirements = false
		}

		if currentFeature != nil {
			if matches := priorityRegex.FindStringSubmatch(line); matches != nil {
				currentFeature.Priority = matches[1]
			}

			if strings.Contains(line, "#### Requirements") {
				inRequirements = true
				continue
			}

			if strings.HasPrefix(line, "####") && !strings.Contains(line, "Requirements") {
				inRequirements = false
			}

			if inRequirements {
				if matches := requirementRegex.FindStringSubmatch(line); matches != nil {
					currentFeature.Requirements = append(currentFeature.Requirements, matches[1])
				}
			}
		}
	}

	if currentFeature != nil {
		features[currentFeature.Number] = currentFeature
	}
	if currentPhase != nil {
		phases = append(phases, *currentPhase)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("error reading PRD: %w", err)
	}

	return phases, features, nil
}

// GenerateTasksFromPRD generates a tasks.json structure from the PRD
func GenerateTasksFromPRD(prdPath, projectName, version string) (*Project, error) {
	phases, features, err := ParsePRD(prdPath)
	if err != nil {
		return nil, err
	}

	project := &Project{
		ProjectName: projectName,
		Version:     version,
		Phases:      []Phase{},
		Summary: Summary{
			TotalPhases:         len(phases),
			TotalTasks:          0,
			TotalSubtasks:       0,
			EstimatedTotalHours: 0,
			Priorities:          make(map[string]int),
		},
	}

	taskCounter := 0
	subtaskCounter := 0
	totalHours := 0

	for _, prdPhase := range phases {
		phase := Phase{
			ID:       fmt.Sprintf("phase-%s", prdPhase.Number),
			Name:     prdPhase.Name,
			Timeline: prdPhase.Timeline,
			Status:   StatusPending,
			Tasks:    []Task{},
		}

		for featureNum, feature := range features {
			// Match features to phases based on number prefix
			if strings.HasPrefix(featureNum, "3.") { // Core features go to appropriate phases
				taskCounter++
				task := Task{
					ID:             fmt.Sprintf("task-%s.%d", prdPhase.Number, taskCounter),
					Title:          feature.Title,
					Description:    fmt.Sprintf("Implement %s as per PRD section %s", feature.Title, featureNum),
					Priority:       feature.Priority,
					EstimatedHours: feature.EstimatedHrs,
					Status:         StatusPending,
					Assignee:       nil,
					Subtasks:       []Subtask{},
				}

				// Create subtasks from requirements
				for i, req := range feature.Requirements {
					if !strings.Contains(req, ":") { // Skip nested lists
						subtaskCounter++
						subtask := Subtask{
							ID:     fmt.Sprintf("task-%s.%d.%d", prdPhase.Number, taskCounter, i+1),
							Title:  req,
							Status: StatusPending,
						}
						task.Subtasks = append(task.Subtasks, subtask)
					}
				}

				project.Summary.Priorities[feature.Priority]++
				totalHours += feature.EstimatedHrs
				phase.Tasks = append(phase.Tasks, task)
			}
		}

		if len(phase.Tasks) > 0 {
			project.Phases = append(project.Phases, phase)
		}
	}

	project.Summary.TotalTasks = taskCounter
	project.Summary.TotalSubtasks = subtaskCounter
	project.Summary.EstimatedTotalHours = totalHours

	return project, nil
}

// InitFromPRD initializes tasks.json from PRD.md
func InitFromPRD(prdPath, outputPath, projectName, version string) error {
	fmt.Printf("Parsing PRD from %s...\n", prdPath)

	project, err := GenerateTasksFromPRD(prdPath, projectName, version)
	if err != nil {
		return fmt.Errorf("failed to generate tasks: %w", err)
	}

	fmt.Printf("Generated %d phases, %d tasks, %d subtasks\n",
		project.Summary.TotalPhases,
		project.Summary.TotalTasks,
		project.Summary.TotalSubtasks)

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := ioutil.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("✓ Successfully generated %s\n", outputPath)
	return nil
}

// Main CLI interface
func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	tm, err := NewTaskMaster("tasks.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "status":
		tm.PrintStatus()

	case "progress":
		progress := tm.GetProgress()
		data, _ := json.MarshalIndent(progress, "", "  ")
		fmt.Println(string(data))

	case "update":
		if len(os.Args) < 4 {
			fmt.Println("Usage: taskmaster update <task-id> <new-status> [notes]")
			os.Exit(1)
		}
		taskID := os.Args[2]
		newStatus := os.Args[3]
		notes := ""
		if len(os.Args) > 4 {
			notes = strings.Join(os.Args[4:], " ")
		}
		if err := tm.UpdateTaskStatus(taskID, newStatus, notes); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "update-subtask":
		if len(os.Args) < 5 {
			fmt.Println("Usage: taskmaster update-subtask <task-id> <subtask-id> <new-status>")
			os.Exit(1)
		}
		if err := tm.UpdateSubtaskStatus(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "assign":
		if len(os.Args) < 4 {
			fmt.Println("Usage: taskmaster assign <task-id> <assignee>")
			os.Exit(1)
		}
		if err := tm.AssignTask(os.Args[2], os.Args[3]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "details":
		if len(os.Args) < 3 {
			fmt.Println("Usage: taskmaster details <task-id>")
			os.Exit(1)
		}
		if err := tm.PrintTaskDetails(os.Args[2]); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "history":
		limit := 20
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &limit)
		}
		if err := tm.PrintHistory(limit); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	case "init":
		prdPath := "PRD.md"
		outputPath := "tasks.json"
		projectName := "etcd-monitor"
		version := "1.0.0"

		if len(os.Args) > 2 {
			prdPath = os.Args[2]
		}
		if len(os.Args) > 3 {
			outputPath = os.Args[3]
		}
		if len(os.Args) > 4 {
			projectName = os.Args[4]
		}
		if len(os.Args) > 5 {
			version = os.Args[5]
		}

		if err := InitFromPRD(prdPath, outputPath, projectName, version); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("TaskMaster - Project Task Management CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  taskmaster status                                  Show all tasks and their status")
	fmt.Println("  taskmaster progress                                Show progress statistics")
	fmt.Println("  taskmaster details <task-id>                       Show detailed task information")
	fmt.Println("  taskmaster update <task-id> <status> [notes]       Update task status")
	fmt.Println("  taskmaster update-subtask <task-id> <subtask-id> <status>")
	fmt.Println("  taskmaster assign <task-id> <assignee>             Assign task to someone")
	fmt.Println("  taskmaster history [limit]                         Show state change history")
	fmt.Println("  taskmaster init [prd-file] [output] [name] [ver]   Generate tasks.json from PRD.md")
	fmt.Println("\nStatus values: pending, in_progress, completed, blocked")
	fmt.Println("\nInit command defaults:")
	fmt.Println("  prd-file: PRD.md")
	fmt.Println("  output:   tasks.json")
	fmt.Println("  name:     etcd-monitor")
	fmt.Println("  version:  1.0.0")
}
