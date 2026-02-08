# Quick Start Guide

## Getting Started in 5 Minutes

### 1. Get API Keys (Required)

You need at least one API key (TVDB or TMDB, or both for best results):

**TVDB**: https://thetvdb.com/api-information  
**TMDB**: https://www.themoviedb.org/settings/api

### 2. Build the Application

```bash
# Using Makefile (recommended)
make build

# Or using Go directly
go build -o kodi-renamer ./cmd/kodi-renamer
```

### 3. Set Your API Keys

```bash
export TVDB_API_KEY='your-tvdb-api-key'
export TMDB_API_KEY='your-tmdb-api-key'
```

### 4. Configure Directories

```bash
# For movies only
export MOVIE_TO_RENAME_DIR='/path/to/movies'

# For TV series only
export SERIE_TO_RENAME_DIR='/path/to/series'

# Or configure both
export MOVIE_TO_RENAME_DIR='/path/to/movies'
export SERIE_TO_RENAME_DIR='/path/to/series'

# Optional: Separate output directories
export MOVIE_RENAMED_DIR='/path/to/renamed/movies'
export SERIE_RENAMED_DIR='/path/to/renamed/series'
```

### 5. Run Your First Rename

```bash
# Preview changes (dry-run mode)
./kodi-renamer -dry-run

# Actually rename files
./kodi-renamer

# Auto mode (no prompts)
./kodi-renamer -auto

# Or use Makefile targets
make run-movies MOVIE_TO_RENAME_DIR='/path/to/movies'
make run-series SERIE_TO_RENAME_DIR='/path/to/series'
make dry-run MOVIE_TO_RENAME_DIR='/path/to/movies' SERIE_TO_RENAME_DIR='/path/to/series'
```

## What to Expect

### For TV Series

**Before:**
```
Breaking.Bad/
├── Breaking.Bad.S01E01.720p.mkv
├── Breaking.Bad.S01E02.1080p.mkv
└── Breaking.Bad.S01E03.mkv
```

**After:**
```
Breaking Bad (2008)/
├── Breaking Bad S01E01 - Pilot.mkv
├── Breaking Bad S01E02 - Cat's in the Bag....mkv
└── Breaking Bad S01E03 - ...And the Bag's in the River.mkv
```

**What Happens:**
1. You select the series ONCE for the entire folder
2. See a table of ALL episodes that will be renamed
3. Confirm ONCE for all episodes
4. Folder and all files renamed automatically
5. Summary: "✓ Successfully renamed 3/3 episodes"

### For Movies

**Before:**
```
The.Matrix.1999.1080p.BluRay.mkv
Inception.2010.720p.BRRip.mp4
```

**After:**
```
The Matrix (1999).mkv
Inception (2010).mp4
```

**What Happens:**
1. You select the correct movie from a table
2. Confirm the rename
3. File renamed with clean format

## Interactive Mode (Default)

```bash
./kodi-renamer -dir /media/shows
```

**Example Session:**
```
Processing series folder: /media/Breaking.Bad
Searching for series: 'Breaking Bad'

Select TV series:
#    Name                  Year  Status      Genres                    Source
1    Breaking Bad          2008  Ended       Drama, Crime, Thriller    tmdb
2    Breaking Bad: Making  2009  Ended       Documentary               tvdb
3    Skip / None

Select an option (number): 1

Series: Breaking Bad (2008)
Status: Ended
Genres: Drama, Crime, Thriller

Fetching episode details...

=== Pending Renames for: Breaking Bad ===

Episode   Current Name                     New Name                         Episode Title           Status
S01E01    Breaking.Bad.S01E01.720p.mkv    Breaking Bad S01E01 - Pilot...   Pilot                  ✓ OK
S01E02    Breaking.Bad.S01E02.1080p.mkv   Breaking Bad S01E02 - Cat's...   Cat's in the Bag...    ✓ OK
S01E03    Breaking.Bad.S01E03.mkv         Breaking Bad S01E03 - ...And...  ...And the Bag's...    ✓ OK

Series folder will be renamed to: Breaking Bad (2008)

Proceed with renaming 3 episode(s)? (y/n): y

✓ Successfully renamed 3/3 episodes
```

## Auto Mode (No Prompts)

```bash
./kodi-renamer -dir /media/shows -auto
```

- Automatically selects first search result
- Processes entire folders at once
- Perfect for trusted, well-named files

## Dry-Run Mode (Preview Only)

```bash
./kodi-renamer -dir /media/shows -dry-run
```

- Shows what WOULD be renamed
- Makes NO actual changes
- Safe for testing

**Combine with auto:**
```bash
./kodi-renamer -dir /media/shows -auto -dry-run
```

## Key Features

### Batch Processing
- ✅ Process entire series folders at once
- ✅ One selection for all episodes
- ✅ See ALL changes before confirming
- ✅ 90% fewer prompts, 10x faster

### Folder Renaming
- ✅ Automatically renames series folders
- ✅ Format: `Series Name (Year)/`
- ✅ Kodi-compatible structure

### Error Prevention
- ✅ Unknown episodes detected BEFORE rename
- ✅ Clear warnings in table
- ✅ Won't proceed with errors
- ✅ Maintains consistency

### Table Displays
- ✅ Beautiful formatted tables
- ✅ Shows all relevant info (year, genres, runtime, etc.)
- ✅ Status indicators (✓ OK, ✗ NOT FOUND)
- ✅ Easy to read and verify

## Common Commands

```bash
# Basic usage with both APIs
export TVDB_API_KEY='your-key'
export TMDB_API_KEY='your-key'
./kodi-renamer -dir /path/to/media

# Using command-line flags
./kodi-renamer -tvdb-key 'key1' -tmdb-key 'key2' -dir /path/to/media

# TVDB only
export TVDB_API_KEY='your-key'
./kodi-renamer -dir /path/to/media

# TMDB only
export TMDB_API_KEY='your-key'
./kodi-renamer -dir /path/to/media

# Preview in auto mode
./kodi-renamer -dir /path/to/media -auto -dry-run

# Process without prompts
./kodi-renamer -dir /path/to/media -auto

# Docker usage
make docker-build
export TVDB_API_KEY='your-key'
export TMDB_API_KEY='your-key'
make docker-run
```

## Makefile Shortcuts

```bash
# Build
make build

# Run with dry-run
make dry-run

# Run with auto mode
make auto

# Docker
make docker-build
make docker-run

# Tests
make test              # Go unit tests
make test-dual-api     # API configuration tests
make test-integration  # Full integration tests (requires real keys)
make test-all          # All tests

# Help
make help
```

## Tips & Best Practices

### For Best Results
1. **Use both APIs** - TVDB + TMDB gives most complete results
2. **Preview first** - Always try `-dry-run` before actual rename
3. **Organize folders** - One series per folder works best
4. **Check episode numbers** - Ensure S##E## format in filenames
5. **Clean folders** - Remove non-video files before processing

### File Organization
```
/media/
├── Breaking.Bad/           ← One series per folder
│   ├── S01E01.mkv
│   ├── S01E02.mkv
│   └── S01E03.mkv
├── The.Matrix.1999.mkv     ← Movies in root or dedicated folder
└── Inception.2010.mp4
```

### If Episodes Show "NOT FOUND"
- Check episode number in filename is correct
- Verify correct series was selected
- Check if episode exists in TVDB/TMDB databases
- Fix filename and re-run

## Troubleshooting

### "No API keys found"
**Solution:** Set environment variables
```bash
export TVDB_API_KEY='your-key'
export TMDB_API_KEY='your-key'
```

### "Episode not found"
**Solution:** Verify episode exists in database or fix episode number

### "File already exists"
**Solution:** Remove duplicate file or rename existing file

### "Directory already exists"
**Solution:** Remove or rename existing folder with target name

## What Gets Renamed

### TV Series
- ✅ Folder renamed to: `Series Name (YYYY)/`
- ✅ Files renamed to: `Series Name S##E## - Episode Title.ext`
- ✅ Quality tags removed (720p, 1080p, WEB-DL, etc.)
- ✅ Special characters cleaned

### Movies
- ✅ Files renamed to: `Movie Title (YYYY).ext`
- ✅ Quality tags removed
- ✅ Special characters cleaned
- ✅ Proper year from database

## Next Steps

1. **Read full documentation:**
   - `docs/SERIES_BATCH_RENAMING.md` - Detailed batch processing guide
   - `docs/TESTING.md` - Testing documentation
   - `CHANGELOG.md` - All changes and features

2. **Try it out:**
   - Start with `-dry-run` mode
   - Test on a small subset first
   - Use `-auto` once comfortable

3. **Get help:**
   - Run `./kodi-renamer -h` for command options
   - Run `make help` for Makefile targets
   - Check documentation in `docs/` folder

## Success!

You're now ready to organize your media library efficiently with batch processing, automatic folder renaming, and full preview capabilities!

**Remember:** Always backup important files before batch operations. Use `-dry-run` first to preview changes.

Happy renaming! 🎬📺