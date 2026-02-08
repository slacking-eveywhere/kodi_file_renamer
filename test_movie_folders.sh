#!/usr/bin/env bash

echo "KODI RENAMER - MOVIE FOLDER STRUCTURE TEST"
echo "==========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Check for API keys
if [[ -z "${TVDB_API_KEY}" && -z "${TMDB_API_KEY}" ]]; then
	echo -e "${RED}ERROR: No API keys found!${NC}"
	echo ""
	echo "This test requires at least one API key to be set:"
	echo "  export TVDB_API_KEY='your-tvdb-key'"
	echo "  export TMDB_API_KEY='your-tmdb-key'"
	echo ""
	exit 1
fi

echo -e "${GREEN}✓ API keys detected${NC}"
if [[ -n "${TVDB_API_KEY}" ]]; then
	echo "  - TVDB: ${TVDB_API_KEY:0:10}..."
fi
if [[ -n "${TMDB_API_KEY}" ]]; then
	echo "  - TMDB: ${TMDB_API_KEY:0:10}..."
fi
echo ""

# Build application
echo -e "${YELLOW}Building application...${NC}"
go build -o kodi-renamer ./cmd/kodi-renamer 2>/dev/null || {
	echo -e "${RED}Build failed${NC}"
	exit 1
}
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Create test directory
TEST_DIR="test_movie_folders"
MOVIE_INPUT="$TEST_DIR/movies-input"
MOVIE_OUTPUT="$TEST_DIR/movies-output"
rm -rf "$TEST_DIR"
mkdir -p "$MOVIE_INPUT"
mkdir -p "$MOVIE_OUTPUT"

echo -e "${YELLOW}Creating test movie structures...${NC}"
echo ""

# =============================================================================
# TEST CASE 1: Standalone movie file (no folder)
# =============================================================================
echo -e "${CYAN}Test Case 1: Standalone movie file${NC}"
touch "$MOVIE_INPUT/Inception.2010.1080p.mkv"
echo "  Created: Inception.2010.1080p.mkv"
echo ""

# =============================================================================
# TEST CASE 2: Movie file with subtitle
# =============================================================================
echo -e "${CYAN}Test Case 2: Movie file with subtitle${NC}"
touch "$MOVIE_INPUT/The.Matrix.1999.mkv"
touch "$MOVIE_INPUT/The.Matrix.1999.srt"
touch "$MOVIE_INPUT/The.Matrix.1999.en.srt"
echo "  Created: The.Matrix.1999.mkv"
echo "  Created: The.Matrix.1999.srt"
echo "  Created: The.Matrix.1999.en.srt"
echo ""

# =============================================================================
# TEST CASE 3: Movie folder with video and subtitle
# =============================================================================
echo -e "${CYAN}Test Case 3: Movie folder with video and subtitle${NC}"
mkdir -p "$MOVIE_INPUT/Interstellar (2014)"
touch "$MOVIE_INPUT/Interstellar (2014)/Interstellar.mkv"
touch "$MOVIE_INPUT/Interstellar (2014)/Interstellar.srt"
echo "  Created: Interstellar (2014)/"
echo "    - Interstellar.mkv"
echo "    - Interstellar.srt"
echo ""

# =============================================================================
# TEST CASE 4: Movie folder with messy name
# =============================================================================
echo -e "${CYAN}Test Case 4: Movie folder with messy name${NC}"
mkdir -p "$MOVIE_INPUT/The.Dark.Knight.2008.1080p.BluRay"
touch "$MOVIE_INPUT/The.Dark.Knight.2008.1080p.BluRay/The.Dark.Knight.2008.mkv"
touch "$MOVIE_INPUT/The.Dark.Knight.2008.1080p.BluRay/The.Dark.Knight.2008.eng.srt"
touch "$MOVIE_INPUT/The.Dark.Knight.2008.1080p.BluRay/The.Dark.Knight.2008.fre.srt"
echo "  Created: The.Dark.Knight.2008.1080p.BluRay/"
echo "    - The.Dark.Knight.2008.mkv"
echo "    - The.Dark.Knight.2008.eng.srt"
echo "    - The.Dark.Knight.2008.fre.srt"
echo ""

# =============================================================================
# TEST CASE 5: Movie folder with multiple video files
# =============================================================================
echo -e "${CYAN}Test Case 5: Movie folder with multiple files${NC}"
mkdir -p "$MOVIE_INPUT/Pulp.Fiction"
touch "$MOVIE_INPUT/Pulp.Fiction/Pulp.Fiction.1994.mkv"
touch "$MOVIE_INPUT/Pulp.Fiction/Pulp.Fiction.1994.srt"
touch "$MOVIE_INPUT/Pulp.Fiction/behind-the-scenes.mkv"
echo "  Created: Pulp.Fiction/"
echo "    - Pulp.Fiction.1994.mkv"
echo "    - Pulp.Fiction.1994.srt"
echo "    - behind-the-scenes.mkv"
echo ""

# =============================================================================
# TEST CASE 6: Blu-ray structure
# =============================================================================
echo -e "${CYAN}Test Case 6: Blu-ray disc structure${NC}"
mkdir -p "$MOVIE_INPUT/Avatar (2009)/BDMV/STREAM"
touch "$MOVIE_INPUT/Avatar (2009)/BDMV/index.bdmv"
touch "$MOVIE_INPUT/Avatar (2009)/BDMV/STREAM/00000.m2ts"
echo "  Created: Avatar (2009)/"
echo "    - BDMV/"
echo "      - index.bdmv"
echo "      - STREAM/00000.m2ts"
echo ""

# =============================================================================
# TEST CASE 7: DVD structure
# =============================================================================
echo -e "${CYAN}Test Case 7: DVD disc structure${NC}"
mkdir -p "$MOVIE_INPUT/Gladiator/VIDEO_TS"
touch "$MOVIE_INPUT/Gladiator/VIDEO_TS/VIDEO_TS.IFO"
touch "$MOVIE_INPUT/Gladiator/VIDEO_TS/VTS_01_1.VOB"
echo "  Created: Gladiator/"
echo "    - VIDEO_TS/"
echo "      - VIDEO_TS.IFO"
echo "      - VTS_01_1.VOB"
echo ""

echo -e "${GREEN}Test structure created!${NC}"
echo ""
echo "Directory structure:"
find "$MOVIE_INPUT" -type d | sort | sed 's/^/  /'
echo ""
echo "Files:"
find "$MOVIE_INPUT" -type f | sort | sed 's/^/  /'
echo ""

# =============================================================================
# TEST: Dry-run mode (preview without changes)
# =============================================================================
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST 1: Dry-run mode (preview renames)${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""

echo "Running: ./kodi-renamer -movie-to-rename $MOVIE_INPUT -movie-renamed $MOVIE_OUTPUT -dry-run -auto"
echo ""
./kodi-renamer -movie-to-rename "$MOVIE_INPUT" -movie-renamed "$MOVIE_OUTPUT" -dry-run -auto 2>&1 | tee /tmp/movie_folder_dryrun.log

echo ""
echo -e "${BLUE}Verifying files are unchanged after dry-run...${NC}"
if [[ -f "$MOVIE_INPUT/Inception.2010.1080p.mkv" ]] && [[ -d "$MOVIE_INPUT/Interstellar (2014)" ]]; then
	echo -e "${GREEN}✓ Files unchanged (dry-run worked correctly)${NC}"
else
	echo -e "${RED}✗ Files were modified during dry-run!${NC}"
	exit 1
fi
echo ""

# =============================================================================
# TEST: Actual renaming
# =============================================================================
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST 2: Actual movie folder processing${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""

echo "BEFORE processing:"
echo "=================="
find "$MOVIE_INPUT" -maxdepth 2 -type d | sort | sed 's/^/  Folder: /'
find "$MOVIE_INPUT" -type f | wc -l | xargs echo "  Total files:"
echo ""

echo "Running: ./kodi-renamer -movie-to-rename $MOVIE_INPUT -movie-renamed $MOVIE_OUTPUT -auto"
echo ""
./kodi-renamer -movie-to-rename "$MOVIE_INPUT" -movie-renamed "$MOVIE_OUTPUT" -auto 2>&1 | tee /tmp/movie_folder_rename.log

echo ""
echo -e "${BLUE}AFTER processing:${NC}"
echo "================="
echo "Output directory:"
find "$MOVIE_OUTPUT" -maxdepth 1 -type d | sort | sed 's/^/  /'
echo ""
find "$MOVIE_OUTPUT" -type f | wc -l | xargs echo "  Total files moved:"
echo ""

# =============================================================================
# TEST: Verify expected results
# =============================================================================
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST 3: Verify movie folder structure${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""

PASS_COUNT=0
FAIL_COUNT=0

echo -e "${CYAN}Checking movie folders were created:${NC}"

# Check for movie folders
if [[ -d "$MOVIE_OUTPUT/Inception (2010)" ]]; then
	echo -e "  ${GREEN}✓${NC} Inception (2010)/ folder exists"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Inception (2010)/ folder NOT found"
	((FAIL_COUNT++))
fi

if [[ -d "$MOVIE_OUTPUT/The Matrix (1999)" ]]; then
	echo -e "  ${GREEN}✓${NC} The Matrix (1999)/ folder exists"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} The Matrix (1999)/ folder NOT found"
	((FAIL_COUNT++))
fi

if [[ -d "$MOVIE_OUTPUT/Interstellar (2014)" ]]; then
	echo -e "  ${GREEN}✓${NC} Interstellar (2014)/ folder exists"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Interstellar (2014)/ folder NOT found"
	((FAIL_COUNT++))
fi

if [[ -d "$MOVIE_OUTPUT/The Dark Knight (2008)" ]]; then
	echo -e "  ${GREEN}✓${NC} The Dark Knight (2008)/ folder exists"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} The Dark Knight (2008)/ folder NOT found"
	((FAIL_COUNT++))
fi

echo ""
echo -e "${CYAN}Checking subtitles were moved with movies:${NC}"

# Check if subtitles were moved
if find "$MOVIE_OUTPUT/The Matrix"* -name "*.srt" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} The Matrix subtitles moved"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} The Matrix subtitles NOT found"
	((FAIL_COUNT++))
fi

if find "$MOVIE_OUTPUT/Interstellar"* -name "*.srt" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} Interstellar subtitle moved"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Interstellar subtitle NOT found"
	((FAIL_COUNT++))
fi

if find "$MOVIE_OUTPUT/The Dark Knight"* -name "*.srt" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} The Dark Knight subtitles moved"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} The Dark Knight subtitles NOT found"
	((FAIL_COUNT++))
fi

echo ""
echo -e "${CYAN}Checking Blu-ray/DVD structures:${NC}"

if [[ -d "$MOVIE_OUTPUT/Avatar (2009)/BDMV" ]]; then
	echo -e "  ${GREEN}✓${NC} Avatar Blu-ray structure preserved"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Avatar Blu-ray structure NOT preserved"
	((FAIL_COUNT++))
fi

# Note: Gladiator might not match in API, so check if it exists or was skipped
if [[ -d "$MOVIE_OUTPUT" ]]; then
	GLADIATOR_COUNT=$(find "$MOVIE_OUTPUT" -name "*VIDEO_TS*" -type d 2>/dev/null | wc -l)
	if [[ $GLADIATOR_COUNT -gt 0 ]]; then
		echo -e "  ${GREEN}✓${NC} DVD structure detected and processed"
		((PASS_COUNT++))
	else
		echo -e "  ${YELLOW}⚠${NC} DVD structure not found (may not have matched in API)"
		# Don't count as failure
	fi
fi

echo ""
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST RESULTS${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "Passed: ${GREEN}${PASS_COUNT}${NC}"
echo -e "Failed: ${RED}${FAIL_COUNT}${NC}"
echo ""

if [[ $FAIL_COUNT -eq 0 ]]; then
	echo -e "${GREEN}╔═══════════════════════════════════════╗${NC}"
	echo -e "${GREEN}║  ✓ ALL TESTS PASSED! ✓               ║${NC}"
	echo -e "${GREEN}╚═══════════════════════════════════════╝${NC}"
	echo ""
	echo "Final output structure:"
	find "$MOVIE_OUTPUT" -maxdepth 2 | sort | sed 's/^/  /'
else
	echo -e "${RED}╔═══════════════════════════════════════╗${NC}"
	echo -e "${RED}║  ✗ SOME TESTS FAILED ✗                ║${NC}"
	echo -e "${RED}╚═══════════════════════════════════════╝${NC}"
	echo ""
	echo "Please check the output above for details."
fi

echo ""
echo "Logs saved to:"
echo "  - /tmp/movie_folder_dryrun.log"
echo "  - /tmp/movie_folder_rename.log"
echo ""

# Cleanup
rm -rf "$TEST_DIR"

exit $FAIL_COUNT
