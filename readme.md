# CodeCowboy 🤠

This program executes GitHub Classroom repositories and outputs usable CSV for Canvas/SIS grade imports.

## Overview

Usage comes to three steps:

1. Set up a database of student identifiers mapping to Canvas identifiers
2. Clone classroom (helper provided)
3. Run grader
4. Import result file

### Database

The database used throughout this project is [Charm KV](https://github.com/charmbracelet/charm#charm-kv). You should self-host your own Charm server: [Charm Self Hosting](https://github.com/charmbracelet/charm#self-hosting). In addition, you may want [Skate](https://github.com/charmbracelet/skate) to inspect your database.

## Web

Requirements:

- Linux/MacOS shell (with `/bin/sh`)
- [GitHub CLI](https://cli.github.com) `gh` utility
- Authenticate with `gh auth login`
- CLI [Classroom extension](https://docs.github.com/en/education/manage-coursework-with-github-classroom/teach-with-github-classroom/using-github-classroom-with-github-cli)
  - Install with `gh extension install github/gh-classroom`
- Canvas student export
- GitHub classroom student export

Run with cmd/web. This is an unauthenticated service automatically running unverified student code. Be careful.

## Commandline

### Creating a classroom

Requirements:

- GitHub Classroom roster (download from the students tab of a classroom)
- Canvas grade export
- A build version of `cmd/mkclassroom`

Flags:

- `course`: The key in the database for this course, anything you want
- `canvaspath`: The path to your Canvas CSV export (CSV)
- `ghpath`: The path to your GitHub Classroom export (CSV)

Example:

`mkclassroom -course=cowboytest -debug -canvaspath=CowboyTest_canvas.csv -ghpath=classroom_roster.csv`

### Clone repositories

Requirements:

- [GitHub CLI](https://cli.github.com) `gh` utility
  - Authenticate with `gh auth login`
- CLI [Classroom extension](https://docs.github.com/en/education/manage-coursework-with-github-classroom/teach-with-github-classroom/using-github-classroom-with-github-cli)
  - Install with `gh extension install github/gh-classroom`

You may use `scripts/clone.sh` to automatically attempt to match course names and assignment names.

`scripts/clone.sh --dest ./student-repositories -v MyCourseName Assignment1`

### Running graders

Requirements:

- Cloned student repositories
- An assignment ID from Canvas
  - This ID will be the trailing number in the assignment's URL
  - The ID will be used to match the column for importing grades

Flags:

- `dir`: The directory of student repositories
- `course`: The database identifier for your course/student roster
- `assignment`: The assignment name for Canvas import
  - Include your ID in this in parens, e.g. "Assignment 1 (ID HERE)"
- `type`: Language we are grading 
  - Supports .NET, Java, and Go

Example:

`-dir=./student-repositories/assignment1-submissions -course=cowboytest -assignment="Assignment 1 (123)" -type=net`

## Tasks

### Generate

Generates templates.

```
go generate
```

### Web

Build the webserver.

requires: Generate

```
go build -o server cmd/web/main.go
```

### Serve

Run the webserver.

requires: Web

```
./server
```

### MkClassroom

Build the classroom import tool.

```
go build -o mkclassroom cmd/mkclassroom/main.go
```

### Import

Imports example classroom assuming example files exist.

requires: MkClassroom

```
./mkclassroom -canvaspath=example/CowboyTest_canvas.csv -ghpath=example/classroom_roster.csv -assignments=example/cowboytest_assign.json -course=cowboytest
```

## Is it good?

I think it's sufficient.