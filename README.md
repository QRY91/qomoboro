# qomoboro - Canonical Hours Task Manager

A systematic task management tool designed around canonical hours and gamified productivity tracking, built with Go and a CLI-first approach.

## Philosophy

**qomoboro** embodies the QRY methodology (Query, Refine, Yield) applied to personal productivity:

- **Query**: What activities truly contribute to your goals?
- **Refine**: How can you optimize your schedule for maximum impact?
- **Yield**: What insights emerge from systematic tracking?

The tool works with human psychology rather than against it, providing structure without rigidity, and insights without intrusion.

## Quick Start

### Installation

```bash
# Clone and build
git clone https://github.com/QRY91/qomoboro.git
cd qomoboro
make install

# Verify installation
qomoboro version
```

### Basic Usage

```bash
# Add a task with title, description, and scores (Work/Play/Learn 0-5)
qomoboro add "Fix API bug" "Memory leak in handler" 4 1 3

# List all tasks
qomoboro list

# Check current status and canonical hour
qomoboro status

# Complete a task
qomoboro complete 1

# View canonical hours schedule
qomoboro schedule
```

## Core Concepts

### Canonical Hours Framework

Traditional time blocks adapted for modern productivity:

- **Matins** (06:00-07:30): Deep work, planning
- **Lauds** (07:30-09:00): Administrative tasks
- **Prime** (09:00-12:00): High-focus work blocks
- **Terce** (12:00-13:30): Meetings, collaboration
- **Sext** (13:30-15:00): Lunch, recovery
- **None** (15:00-16:30): Creative work, experimentation
- **Vespers** (16:30-18:00): Learning, documentation
- **Compline** (18:00-20:00): Planning, reflection

### Triple-Metric Scoring System

Each task receives three independent scores (0-5 scale):

- **Work**: Business/professional productivity value
- **Play**: Recreation/leisure/enjoyment factor
- **Learn**: Educational/skill development potential

Tasks can excel in multiple dimensions:
- Code review: Work 4, Play 2, Learn 3
- Personal project: Work 2, Play 5, Learn 4
- Team meeting: Work 3, Play 2, Learn 2

## Commands Reference

### Task Management
```bash
# Create tasks
qomoboro add "Task title"                           # Basic task
qomoboro add "Fix bug" "Details here" 4 2 3         # With scores
qomoboro add "Research" "" 2 3 5                    # Skip description

# List and manage
qomoboro list                                       # Show all tasks
qomoboro complete 1                                 # Complete task #1
qomoboro delete 2                                   # Delete task #2
```

### Productivity Insights
```bash
qomoboro status                                     # Current hour + summary
qomoboro schedule                                   # Show canonical hours
qomoboro stats                                      # Today's statistics
```

### System Management
```bash
qomoboro backup                                     # Create data backup
qomoboro data-dir                                   # Show data location
qomoboro version                                    # Show version
qomoboro help                                       # Show full help
```

## Features

### Systematic Task Management
- Create, track, and complete tasks with multi-dimensional scoring
- Automatic canonical hour detection and suggestions
- Clean CLI interface with ASCII art and emojis
- Persistent local storage (no cloud dependencies)

### Productivity Analytics
- Daily completion rates and score totals
- Canonical hour productivity patterns
- Work/Play/Learn balance tracking
- Time-based insights and recommendations

### Data Management
- Local JSON file storage for privacy and ownership
- Automatic backup functionality
- Human-readable data format for easy inspection/export
- Cross-platform data directory handling

## Project Structure

```
qomoboro/
├── internal/
│   ├── models/         # Core data structures
│   ├── storage/        # File-based persistence
│   └── ui/            # TUI components (future)
├── docs/              # Documentation
├── assets/           # Audio files and resources
├── Makefile          # Development workflow
└── qomo.go           # CLI application entry point
```

## Development

### Prerequisites
- Go 1.23+
- Make (for build automation)

### Development Workflow
```bash
# Quick development cycle
make quick              # Format, build, and run

# Individual commands
make build              # Build binary
make test               # Run test suite
make format             # Format code
make clean              # Clean build artifacts

# Installation
make install            # Install system-wide
make uninstall          # Remove from system
```

### Testing
```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Development with live data
make run-dev
```

## Data Location

Configuration and data are stored locally:
- **Linux/macOS**: `~/.local/share/qomoboro/`
- **With XDG**: `$XDG_DATA_HOME/qomoboro/`

### File Structure
```
~/.local/share/qomoboro/
├── tasks.json          # All tasks and their data
├── schedule.json       # Canonical hours configuration
├── stats/              # Daily statistics
│   ├── 2024-01-01.json
│   └── 2024-01-02.json
└── backups/            # Automated backups
```

## Evolution

**qomoboro** started as a simple pomodoro timer experiment but evolved into a comprehensive task management system aligned with canonical hours methodology and systematic productivity tracking.

The CLI-first approach prioritizes immediate functionality and scriptability over complex interfaces, making it perfect for systematic users who value efficiency and data ownership.

## Future Development

### Planned Features
- Interactive TUI mode (`qomoboro --interactive`)
- Custom canonical hour configurations
- Task templates and categories
- Export functionality for analysis
- Team collaboration features

### Integration Opportunities
- **Uroboro**: Capture development insights
- **Calendar apps**: Schedule synchronization
- **Time tracking**: Detailed productivity analysis
- **AI tools**: Smart task classification and scheduling

## Philosophy in Practice

qomoboro demonstrates systematic thinking applied to personal productivity:

- **Local-first**: Your data belongs to you
- **Systematic tracking**: Measure what matters across multiple dimensions
- **Human-centered**: Work with natural rhythms, not against them
- **Scriptable**: Integrate with your existing workflow
- **Transparent**: Simple data formats you can inspect and modify

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make ci` to verify
5. Submit pull request

Follow Go conventions and maintain the systematic philosophy of simple, effective tools.

## License

MIT License - See LICENSE file for details

---

**Built with systematic methodology | Available for immediate use | Part of the QRY ecosystem**

*"Three commands beats seventeen commands" - Focus on what works, systematically.*