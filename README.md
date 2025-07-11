# qomoboro - Canonical Hours Task Manager

A systematic task management tool built with Go and Charmbracelet TUI components, designed around canonical hours and gamified productivity tracking.

## Evolution

**qomoboro** started as a simple pomodoro timer experiment but has evolved into a comprehensive task management system aligned with canonical hours methodology and systematic productivity tracking.

## Core Concepts

### Canonical Hours Framework
- **Structured Time Blocks**: Based on traditional canonical hours (Matins, Lauds, Prime, etc.)
- **Flexible Scheduling**: Adaptable to modern work patterns while maintaining systematic structure
- **Natural Rhythms**: Aligns with human energy cycles and attention patterns

### Triple-Metric Scoring System
Each activity is scored across three dimensions (0-5 scale):
- **Work**: Business/professional productivity value
- **Play**: Recreation/leisure/enjoyment value  
- **Learn**: Educational/skill development value

Activities can score in multiple categories (e.g., coding tutorial: Work 4/5, Play 2/5, Learn 5/5)

### Gamification Elements
- **Daily Scores**: Visual progress tracking across all three metrics
- **Pattern Recognition**: Identify optimal scheduling patterns
- **Achievement Tracking**: Long-term productivity insights
- **Non-Intrusive**: Gentle reminders without disrupting flow

## Features

### Core Functionality
- **Task Management**: Create, schedule, and track activities
- **Time Tracking**: Flexible block-based timing system
- **Multi-Metric Scoring**: Work/Play/Learn evaluation for each task
- **Schedule Visualization**: Clean TUI calendar/timeline display
- **Progress Tracking**: Daily, weekly, and monthly insights

### TUI Interface
- **Minimal Design**: Clean, focused terminal interface
- **Keyboard Navigation**: Efficient interaction without mouse dependency
- **Responsive Layout**: Adapts to different terminal sizes
- **Color Coding**: Visual indicators for different activity types and scores

### Smart Features
- **Pattern Learning**: Recognize optimal scheduling patterns
- **Score Suggestions**: AI-assisted activity categorization
- **Habit Tracking**: Long-term productivity pattern analysis
- **Flexible Blocks**: Adapt canonical hours to personal schedule

## Technical Stack

- **Language**: Go 1.23+
- **TUI Framework**: Charmbracelet Bubble Tea
- **UI Components**: Charmbracelet Huh (forms), Lipgloss (styling)
- **Data Storage**: Local file-based persistence
- **Audio**: Beep library for optional notifications

## Installation

```bash
git clone https://github.com/QRY91/qomoboro.git
cd qomoboro
go mod tidy
go build -o qomoboro
```

## Usage

### Basic Operation
```bash
# Start the task manager
./qomoboro

# Quick task entry
./qomoboro add "Research AI patterns" --work 4 --learn 5 --play 2

# View today's schedule
./qomoboro today

# Weekly overview
./qomoboro week
```

### Canonical Hours Setup
The system supports flexible canonical hours configuration:
- **Matins** (06:00-07:30): Deep work, planning
- **Lauds** (07:30-09:00): Administrative tasks
- **Prime** (09:00-12:00): High-focus work blocks
- **Terce** (12:00-13:30): Meetings, collaboration
- **Sext** (13:30-15:00): Lunch, recovery
- **None** (15:00-16:30): Creative work, experimentation
- **Vespers** (16:30-18:00): Learning, documentation
- **Compline** (18:00-20:00): Planning, reflection

## Configuration

### Settings File
```yaml
# ~/.config/qomoboro/config.yaml
canonical_hours:
  matins: "06:00-07:30"
  lauds: "07:30-09:00"
  prime: "09:00-12:00"
  # ... etc

scoring:
  work_weight: 1.0
  play_weight: 1.0
  learn_weight: 1.0

ui:
  theme: "default"
  compact_mode: false
  show_notifications: true
```

## Integration with QRY Ecosystem

### Data Export
- **JSON**: Export for analysis and backup
- **CSV**: Compatible with spreadsheet analysis
- **Markdown**: Human-readable summaries

### AI Collaboration
- **Pattern Recognition**: Identify optimal productivity patterns
- **Score Suggestions**: Machine learning-assisted activity scoring
- **Schedule Optimization**: AI-driven schedule recommendations

### Cross-Tool Integration
- **Uroboro**: Work documentation and acknowledgment
- **Wherewasi**: Context generation for AI collaboration
- **QRY Zone**: Website content and project tracking

## Development

### Architecture
```
qomoboro/
├── cmd/           # CLI commands and entry points
├── internal/      # Core business logic
│   ├── models/    # Data structures
│   ├── storage/   # Persistence layer
│   ├── scoring/   # Scoring system logic
│   └── ui/        # TUI components
├── pkg/          # Reusable packages
├── assets/       # Sounds, themes, etc.
└── docs/         # Documentation
```

### Key Components
- **Task Model**: Core data structure for activities
- **Scheduler**: Canonical hours management
- **Scorer**: Multi-metric evaluation system
- **TUI Manager**: Bubble Tea application state
- **Storage Engine**: Local data persistence

## Contributing

### Development Setup
```bash
# Install dependencies
go mod tidy

# Run tests
go test ./...

# Run with development config
go run . --dev

# Build for release
make build
```

### Code Style
- Follow Go conventions
- Use structured logging
- Maintain TUI responsiveness
- Document public APIs

## Philosophy

**qomoboro** embodies the QRY methodology (Query, Refine, Yield) applied to personal productivity:

- **Query**: What activities truly contribute to your goals?
- **Refine**: How can you optimize your schedule for maximum impact?
- **Yield**: What insights emerge from systematic tracking?

The tool is designed to work with human psychology rather than against it, providing structure without rigidity, and insights without intrusion.

## License

MIT License - See LICENSE file for details

## Acknowledgments

- Charmbracelet team for excellent TUI libraries
- Canonical hours tradition for timeless scheduling wisdom
- QRY community for systematic productivity insights