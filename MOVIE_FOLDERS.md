# Movie Folder Support

## Overview

The Kodi File Renamer now supports movies as folders, not just standalone video files. This allows proper handling of:

- Movies with subtitle files
- Blu-ray disc structures (BDMV)
- DVD disc structures (VIDEO_TS)
- Movies with extras and behind-the-scenes content
- Multi-file movie releases

## Supported Movie Structures

### 1. Standalone Movie File

The simplest case - a single video file without accompanying files.

**Input:**
```
movies/
└── Inception.2010.1080p.mkv
```

**Output (with output directory):**
```
library/
└── Inception (2010)/
    └── Inception (2010).mkv
```

**Output (in-place rename):**
```
movies/
└── Inception (2010).mkv
```

### 2. Movie File with Subtitles

A video file with one or more subtitle files (same base name).

**Input:**
```
movies/
├── The.Matrix.1999.mkv
├── The.Matrix.1999.srt
└── The.Matrix.1999.en.srt
```

**Output:**
```
library/
└── The Matrix (1999)/
    ├── The Matrix (1999).mkv
    ├── The Matrix (1999).srt
    └── The Matrix (1999).en.srt
```

### 3. Movie Folder with Video and Subtitles

A folder containing a movie and its subtitle files.

**Input:**
```
movies/
└── Interstellar (2014)/
    ├── Interstellar.mkv
    └── Interstellar.srt
```

**Output:**
```
library/
└── Interstellar (2014)/
    ├── Interstellar (2014).mkv
    └── Interstellar (2014).srt
```

### 4. Movie Folder with Messy Name

A folder with quality tags and messy naming.

**Input:**
```
movies/
└── The.Dark.Knight.2008.1080p.BluRay/
    ├── The.Dark.Knight.2008.mkv
    ├── The.Dark.Knight.2008.eng.srt
    └── The.Dark.Knight.2008.fre.srt
```

**Output:**
```
library/
└── The Dark Knight (2008)/
    ├── The Dark Knight (2008).mkv
    ├── The Dark Knight (2008).eng.srt
    └── The Dark Knight (2008).fre.srt
```

### 5. Blu-ray Disc Structure

Full Blu-ray disc rip with BDMV folder structure.

**Input:**
```
movies/
└── Avatar (2009)/
    └── BDMV/
        ├── index.bdmv
        ├── MovieObject.bdmv
        └── STREAM/
            ├── 00000.m2ts
            └── 00001.m2ts
```

**Output:**
```
library/
└── Avatar (2009)/
    └── BDMV/
        ├── index.bdmv
        ├── MovieObject.bdmv
        └── STREAM/
            ├── 00000.m2ts
            └── 00001.m2ts
```

The entire folder structure is preserved and moved as-is.

### 6. DVD Disc Structure

Full DVD disc rip with VIDEO_TS folder structure.

**Input:**
```
movies/
└── Gladiator/
    └── VIDEO_TS/
        ├── VIDEO_TS.IFO
        ├── VIDEO_TS.VOB
        ├── VTS_01_0.IFO
        └── VTS_01_1.VOB
```

**Output:**
```
library/
└── Gladiator (2000)/
    └── VIDEO_TS/
        ├── VIDEO_TS.IFO
        ├── VIDEO_TS.VOB
        ├── VTS_01_0.IFO
        └── VTS_01_1.VOB
```

The entire folder structure is preserved and moved as-is.

## Movie Folder Detection Logic

The scanner identifies a folder as a movie folder if it contains:

1. **Blu-ray Structure**: Has a `BDMV/` directory
2. **DVD Structure**: Has a `VIDEO_TS/` directory
3. **Video + Subtitles**: Video files AND subtitle files
4. **Single Video with Year**: One video file and folder name contains a year
5. **Matching Names**: Folder name matches the video filename

### Detection Examples

**Detected as Movie Folder:**
```
✓ Avatar (2009)/BDMV/              # Blu-ray structure
✓ Gladiator/VIDEO_TS/              # DVD structure
✓ Movie/movie.mkv + movie.srt      # Video + subtitle
✓ The Matrix (1999)/matrix.mkv     # Has year in folder name
✓ Inception/Inception.mkv          # Folder matches video name
```

**NOT Detected as Movie Folder:**
```
✗ RandomFolder/                    # No video files
✗ Movies/movie1.mkv movie2.mkv     # Multiple unrelated videos
✗ Downloads/some-video.mkv         # No year, no matching folder
```

These are processed as standalone video files instead.

## Supported Subtitle Formats

- `.srt` - SubRip
- `.sub` - MicroDVD, SubViewer
- `.ass` - Advanced SubStation Alpha
- `.ssa` - SubStation Alpha
- `.vtt` - WebVTT

## Processing Behavior

### With Output Directory Specified

```bash
./kodi-renamer \
  -movie-to-rename /downloads/movies \
  -movie-renamed /library/movies
```

**Behavior:**
1. Scans input directory for movies
2. Identifies standalone files and movie folders
3. Searches API for movie metadata
4. Creates properly named folder in output directory
5. Moves all files (video + subtitles) to new folder
6. Preserves Blu-ray/DVD structures

**Result:**
- Input directory remains (possibly empty after processing)
- Output directory contains clean, organized movie folders
- All subtitles preserved with correct naming

### Without Output Directory (In-Place)

```bash
./kodi-renamer -movie-to-rename /downloads/movies
```

**Behavior:**
1. Scans directory for movies
2. Renames folders in place
3. Renames video files (with subtitles) in place
4. Preserves folder structures

**Result:**
- Movies renamed in their current location
- Standalone files → Stay as files (with renamed subtitles)
- Movie folders → Renamed in place

## Kodi Compatibility

All output structures are compatible with Kodi's movie naming conventions:

### Preferred Structure (with output directory)
```
Movies/
├── Avatar (2009)/
│   └── Avatar (2009).mkv
├── The Matrix (1999)/
│   ├── The Matrix (1999).mkv
│   └── The Matrix (1999).srt
└── Inception (2010)/
    └── BDMV/
        └── ...
```

### Also Valid (in-place, no folders)
```
Movies/
├── Avatar (2009).mkv
├── The Matrix (1999).mkv
└── Inception (2010).mkv
```

Both structures are recognized by Kodi's movie scraper.

## Examples

### Example 1: Process Movies with Subtitles

```bash
# Input structure
downloads/movies/
├── inception.2010.mkv
├── inception.2010.srt
├── matrix.1999.mkv
└── matrix.1999.en.srt

# Command
./kodi-renamer \
  -movie-to-rename downloads/movies \
  -movie-renamed library/movies \
  -auto

# Output structure
library/movies/
├── Inception (2010)/
│   ├── Inception (2010).mkv
│   └── Inception (2010).srt
└── The Matrix (1999)/
    ├── The Matrix (1999).mkv
    └── The Matrix (1999).en.srt
```

### Example 2: Process Blu-ray Collection

```bash
# Input structure
blurays/
├── Avatar/
│   └── BDMV/
│       └── STREAM/
│           └── 00000.m2ts
└── Inception/
    └── BDMV/
        └── STREAM/
            └── 00000.m2ts

# Command
./kodi-renamer \
  -movie-to-rename blurays \
  -movie-renamed library/movies

# Output structure
library/movies/
├── Avatar (2009)/
│   └── BDMV/
│       └── STREAM/
│           └── 00000.m2ts
└── Inception (2010)/
    └── BDMV/
        └── STREAM/
            └── 00000.m2ts
```

### Example 3: Mixed Movie Types

```bash
# Input structure
downloads/
├── standalone.movie.2020.mkv           # Standalone file
├── Movie.With.Subs/                    # Folder with subs
│   ├── movie.mkv
│   └── movie.srt
└── BluRay.Movie/                       # Blu-ray structure
    └── BDMV/

# Command
./kodi-renamer \
  -movie-to-rename downloads \
  -movie-renamed library/movies

# Output structure
library/movies/
├── Standalone Movie (2020)/            # File → Folder
│   └── Standalone Movie (2020).mkv
├── Movie with Subs (2021)/             # Folder → Renamed folder
│   ├── Movie with Subs (2021).mkv
│   └── Movie with Subs (2021).srt
└── BluRay Movie (2019)/                # Blu-ray → Renamed folder
    └── BDMV/
```

## Dry-Run Mode

Always test with dry-run first to preview changes:

```bash
./kodi-renamer \
  -movie-to-rename /downloads/movies \
  -movie-renamed /library/movies \
  -dry-run
```

**Output shows:**
```
Processing movie folder: 'Inception' (1 video files, 1 subtitles)
Searching for: 'Inception'
...
[DRY RUN] Would create movie folder: /library/movies/Inception (2010)
[DRY RUN] Would move file to: /library/movies/Inception (2010)/Inception (2010).mkv
[DRY RUN] Would move subtitle to: /library/movies/Inception (2010)/Inception (2010).srt
```

## Subtitle Handling

### Automatic Subtitle Matching

Subtitles are matched by filename prefix:

```
Movie.mkv          → Movie (2020).mkv
Movie.srt          ✓ Movie (2020).srt
Movie.en.srt       ✓ Movie (2020).en.srt
Movie.eng.srt      ✓ Movie (2020).eng.srt
Movie.forced.srt   ✓ Movie (2020).forced.srt
Other.srt          ✗ Not matched (different name)
```

### Multiple Subtitle Tracks

All matching subtitles are renamed and moved:

```bash
# Input
The.Matrix.1999.mkv
The.Matrix.1999.srt           # Default
The.Matrix.1999.en.srt        # English
The.Matrix.1999.fr.srt        # French
The.Matrix.1999.en.forced.srt # English forced

# Output
The Matrix (1999)/
├── The Matrix (1999).mkv
├── The Matrix (1999).srt
├── The Matrix (1999).en.srt
├── The Matrix (1999).fr.srt
└── The Matrix (1999).en.forced.srt
```

## Special Cases

### Multiple Video Files in Folder

If a folder contains multiple video files:
- Only recognized as movie folder if has Blu-ray/DVD structure OR subtitles
- Otherwise, each video file processed individually

### Extras and Bonus Content

Files with different names are NOT automatically moved:

```
Movie (2020)/
├── movie.mkv              ✓ Moved/renamed
├── movie.srt              ✓ Moved/renamed
└── behind-the-scenes.mkv  ✗ Different name, not moved
```

To include extras, ensure they're in a properly structured folder.

### Already Correctly Named

If a folder is already correctly named, it's still processed:

```bash
# Input (already correct name)
The Matrix (1999)/
└── The Matrix (1999).mkv

# Still moved to output directory if specified
library/movies/The Matrix (1999)/
└── The Matrix (1999).mkv
```

## Troubleshooting

### Issue: Movie folder not detected

**Symptoms:** Folder is processed as series or skipped

**Solutions:**
1. Add year to folder name: `Movie/` → `Movie (2020)/`
2. Ensure folder name matches video filename
3. Add subtitle files to trigger folder detection
4. Check folder contains video files

### Issue: Subtitles not moved

**Symptoms:** Video moved but subtitles left behind

**Solutions:**
1. Ensure subtitle filename matches video filename (prefix)
2. Check subtitle has supported extension (.srt, .sub, .ass, .ssa, .vtt)
3. Verify subtitle is in same directory as video

### Issue: Blu-ray structure not preserved

**Symptoms:** BDMV folder not moved correctly

**Solutions:**
1. Ensure BDMV folder is directly under movie folder
2. Check permissions on source and destination
3. Verify enough disk space for large Blu-ray rips

### Issue: Multiple videos processed separately

**Symptoms:** Each video in folder treated as separate movie

**Solutions:**
1. This is expected if folder doesn't meet movie folder criteria
2. Add year to folder name or subtitles to trigger folder detection
3. Process manually if videos are unrelated

## Best Practices

1. **Use Output Directory**: Separate input/output keeps originals safe
   ```bash
   -movie-to-rename /downloads -movie-renamed /library
   ```

2. **Test with Dry-Run**: Always preview changes first
   ```bash
   -dry-run
   ```

3. **Organize Before Processing**: Use consistent folder structure
   ```
   downloads/movies/
   ├── Movie1 (2020)/
   ├── Movie2 (2021)/
   └── Movie3 (2022)/
   ```

4. **Include Year in Folder Names**: Helps with detection
   ```
   ✓ Inception (2010)/
   ✗ Inception/
   ```

5. **Keep Subtitles with Videos**: Same directory, matching names
   ```
   ✓ movie.mkv + movie.srt
   ✗ movie.mkv in folder1, movie.srt in folder2
   ```

6. **Backup Important Collections**: Especially for Blu-ray/DVD rips
   ```bash
   rsync -av /blurays /backup/blurays
   ```

## Testing

Run the movie folder test suite:

```bash
# Requires API keys
export TMDB_API_KEY='your-key'

# Run test
./test_movie_folders.sh
```

Tests cover:
- Standalone movie files
- Movies with subtitles
- Movie folders
- Blu-ray structures
- DVD structures
- Mixed content

## Related Documentation

- [CONFIGURATION.md](CONFIGURATION.md) - Directory configuration
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [README](README) - General information
- [MIGRATION_GUIDE.md](MIGRATION_GUIDE.md) - Upgrading from older versions

---

**Note:** Movie folder support is available in version 2.1.0+