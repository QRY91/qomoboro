# qomoboro Usage Guide

## Quick Start

### Installation
```bash
# Clone and build
git clone https://github.com/QRY91/qomoboro.git
cd qomoboro
make build

# Run the application
./qomoboro
```

### Create Your First Task
1. Launch qomoboro
2. Press `c` to create a new task
3. Fill in:
   - **Title**: What you need to do
   - **Description**: Optional context
   - **Work Score**: Business value (0-5)
   - **Play Score**: Enjoyment factor (0-5)
   - **Learn Score**: Learning value (0-5)
4. Press Enter to save

## Interface Overview

### Main Menu
- `[t]` Tasks - View and manage your task list
- `[s]` Schedule - View canonical hours
- `[d]` Statistics - View productivity stats
- `[c]` Create - Add new task
- `[g]` Settings - Configure app
- `[r]` Refresh - Reload data
- `[q]` Quit - Exit application

### Task List Navigation
- `↑/↓` or `k/j` - Navigate tasks
- `Enter` - View task details
- `Space` - Toggle task completion
- `c` - Create new task
- `d` - Delete selected task
- `q/Esc` - Back to main menu

## Task Management

### Task States
- **Pending**: Not yet started
- **Active**: Currently being worked on
- **Completed**: Finished
- **Paused**: Temporarily stopped
- **Cancelled**: Abandoned

### Task Operations
- **Create**: Press `c` from main menu or task list
- **Complete**: Press `Space` on any task
- **Delete**: Press `d` on selected task
- **View Details**: Press `Enter` on selected task

## Scoring System

Each task gets three scores (0-5 scale):

### Work Score
- **5**: Critical business/professional task
- **4**: Important work contribution
- **3**: Moderate work value
- **2**: Minor work benefit
- **1**: Minimal work impact
- **0**: No work value

### Play Score
- **5**: Extremely enjoyable/fun
- **4**: Very enjoyable
- **3**: Moderately enjoyable
- **2**: Slightly enjoyable
- **1**: Minimal enjoyment
- **0**: No enjoyment

### Learn Score
- **5**: Major learning opportunity
- **4**: Significant learning
- **3**: Moderate learning
- **2**: Some learning
- **1**: Minimal learning
- **0**: No learning

### Multi-Dimensional Scoring
Tasks can score high in multiple areas:
- Coding tutorial: Work 4, Play 2, Learn 5
- Team meeting: Work 3, Play 2, Learn 2
- Personal project: Work 2, Play 5, Learn 4

## Canonical Hours

Traditional time blocks adapted for modern productivity:

### Schedule Overview
- **Matins** (06:00-07:30): Deep work, planning
- **Lauds** (07:30-09:00): Administrative tasks
- **Prime** (09:00-12:00): High-focus work blocks
- **Terce** (12:00-13:30): Meetings, collaboration
- **Sext** (13:30-15:00): Lunch, recovery
- **None** (15:00-16:30): Creative work, experimentation
- **Vespers** (16:30-18:00): Learning, documentation
- **Compline** (18:00-20:00): Planning, reflection

### Using Canonical Hours
- Schedule tasks during appropriate hours
- View current canonical hour on main screen
- Use suggested scoring for each time block
- Adapt schedule to your personal rhythm

## Statistics & Analytics

### Daily Stats
- Total tasks created/completed
- Work/Play/Learn score totals
- Time spent on activities
- Completion rate

### Weekly Overview
- 7-day trend analysis
- Average daily scores
- Productivity patterns
- Goal tracking

### Accessing Stats
- Press `d` from main menu
- View today's summary
- Compare with previous days
- Identify productivity patterns

## Data Management

### Data Location
- Linux/macOS: `~/.local/share/qomoboro/`
- Or `$XDG_DATA_HOME/qomoboro/` if set

### Backup
```bash
# Create backup
./qomoboro --backup

# Show data directory
./qomoboro --data-dir
```

### File Structure
```
~/.local/share/qomoboro/
├── tasks.json          # All tasks
├── schedule.json       # Canonical hours config
├── stats/              # Daily statistics
│   ├── 2024-01-01.json
│   └── 2024-01-02.json
└── backups/            # Backup files
```

## Command Line Usage

### Basic Commands
```bash
# Show help
./qomoboro --help

# Show version
./qomoboro --version

# Show data directory
./qomoboro --data-dir

# Create backup
./qomoboro --backup
```

### Development Commands
```bash
# Quick build and run
make quick

# Run with development data
make run-dev

# Run tests
make test

# Format code
make format
```

## Tips & Best Practices

### Effective Task Management
1. **Be specific** with task titles
2. **Score honestly** - it helps with patterns
3. **Use descriptions** for context
4. **Regular review** - check completed tasks
5. **Adapt scoring** as you learn your patterns

### Canonical Hours Optimization
1. **Match energy levels** to hour purposes
2. **Batch similar tasks** in appropriate hours
3. **Respect recovery time** (Sext)
4. **Plan ahead** during Compline
5. **Deep work** during Matins/Prime

### Scoring Strategy
1. **Consider long-term impact** for Work scores
2. **Rate immediate enjoyment** for Play scores
3. **Think about skill development** for Learn scores
4. **Multi-dimensional is good** - embrace it
5. **Track patterns** over time

### Workflow Integration
1. **Start each day** by reviewing schedule
2. **Create tasks** for each canonical hour
3. **Update status** throughout the day
4. **Review statistics** weekly
5. **Adjust approach** based on data

## Troubleshooting

### Common Issues

**Application won't start**
- Check if binary is built: `make build`
- Verify data directory permissions
- Check for conflicting processes

**Tasks not saving**
- Verify write permissions to data directory
- Check disk space
- Review error messages

**Statistics not updating**
- Complete tasks to generate stats
- Check date/time settings
- Manually refresh with `r`

### Getting Help
- Use `./qomoboro --help` for command reference
- Check the main README for setup instructions
- Review source code for advanced configuration
- Submit issues to the GitHub repository

## Advanced Features

### Customization
- Edit `schedule.json` to modify canonical hours
- Adjust time blocks to your schedule
- Modify default scoring suggestions
- Create custom task categories

### Integration
- Export data as JSON for analysis
- Use with other productivity tools
- Integrate with time tracking systems
- Connect with calendar applications

### Development
- Fork the repository for custom features
- Contribute improvements back to the project
- Use the Makefile for development workflows
- Follow Go best practices for contributions

## Philosophy

qomoboro embodies the QRY methodology:
- **Query**: What activities truly contribute to your goals?
- **Refine**: How can you optimize your schedule for maximum impact?
- **Yield**: What insights emerge from systematic tracking?

The tool works with human psychology rather than against it, providing structure without rigidity, and insights without intrusion.