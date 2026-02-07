# progress

Package `progress` provides terminal progress bar rendering with stage support for CLI utilities. It handles ANSI escape codes for line clearing and overwriting to show updating progress without scrolling the terminal.

## Features

- Single-stage and multi-stage progress tracking
- Unicode and ASCII progress bar rendering
- Terminal line overwriting for smooth updates
- Configurable bar and text widths

## Installation

```go
import "github.com/grokify/mogo/fmt/progress"
```

## Usage

### SingleStageRenderer

Use `SingleStageRenderer` for simple progress tracking with a single operation:

```go
package main

import (
    "os"
    "time"

    "github.com/grokify/mogo/fmt/progress"
)

func main() {
    // Create renderer writing to stdout
    renderer := progress.NewSingleStageRenderer(os.Stdout).
        WithBarWidth(40).
        WithTextWidth(30)

    items := []string{"file1.txt", "file2.txt", "file3.txt"}
    total := len(items)

    for i, item := range items {
        renderer.Update(i+1, total, item)
        time.Sleep(500 * time.Millisecond) // simulate work
    }

    // Clear progress and show completion message
    renderer.Done("Processing complete!")
}
```

Output while running:

```
[████████████████░░░░░░░░░░░░░░░░░░░░░░░░]  40% (2/5) file2.txt
```

### MultiStageRenderer

Use `MultiStageRenderer` for operations with multiple distinct phases:

```go
package main

import (
    "os"

    "github.com/grokify/mogo/fmt/progress"
)

func main() {
    renderer := progress.NewMultiStageRenderer(os.Stderr).
        WithBarWidth(20).
        WithDescWidth(34)

    totalStages := 3

    // Stage 1: Fetching data
    for i := 0; i <= 10; i++ {
        renderer.Update(progress.StageInfo{
            Stage:       1,
            TotalStages: totalStages,
            Description: "Fetching repositories",
            Current:     i,
            Total:       10,
        })
    }
    renderer.Update(progress.StageInfo{
        Stage:       1,
        TotalStages: totalStages,
        Description: "Fetching repositories",
        Done:        true,
    })

    // Stage 2: Processing
    renderer.Update(progress.StageInfo{
        Stage:       2,
        TotalStages: totalStages,
        Description: "Processing data",
        Done:        true,
    })

    // Stage 3: Writing output
    renderer.Update(progress.StageInfo{
        Stage:       3,
        TotalStages: totalStages,
        Description: "Writing results",
        Done:        true,
    })
}
```

Output:

```
[1/3] Fetching repositories              [████████████████████] 100%
[2/3] Processing data                    [████████████████████] 100%
[3/3] Writing results                    [████████████████████] 100%
```

### Progress Bar Rendering

The package provides two bar rendering functions:

```go
// Unicode bar (default): [████████░░░░░░░░░░░░]
bar := progress.RenderBar(50, 20)

// ASCII bar: [==========----------]
bar := progress.RenderBarASCII(50, 20)
```

## Types

### StageInfo

```go
type StageInfo struct {
    Stage       int    // Current stage number (1-based)
    TotalStages int    // Total number of stages
    Description string // What's happening in this stage
    Current     int    // Current item within stage (0 if not applicable)
    Total       int    // Total items in stage (0 if not applicable)
    Done        bool   // True if this stage is complete
    Text        string // Optional extra text (e.g., current item name)
}
```

The `Percent()` method returns the completion percentage for the stage.

## Integration Pattern

A common pattern is to pass a progress callback to long-running functions:

```go
type ProgressFunc func(current, total int, name string)

func ProcessItems(items []Item, progress ProgressFunc) error {
    for i, item := range items {
        if progress != nil {
            progress(i+1, len(items), item.Name)
        }
        // process item...
    }
    return nil
}

// Usage
renderer := progress.NewSingleStageRenderer(os.Stdout)
err := ProcessItems(items, func(current, total int, name string) {
    renderer.Update(current, total, name)
})
renderer.Done("")
```
