#!/usr/bin/env bash

# Note: Don't use set -e because arithmetic operations like ((PASS_COUNT++))
# can return non-zero and cause early exit. took me age to fix this shit

echo "KODI RENAMER - INTEGRATION TEST FOR FILE RENAMING"
echo "================================================="
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

# Create test directories with new 4-directory structure
TEST_DIR="test_rename_integration"
MOVIE_DIR="$TEST_DIR/movies-to-rename"
MOVIE_RENAMED_DIR="$TEST_DIR/movies-renamed"
SERIE_DIR="$TEST_DIR/series-to-rename"
SERIE_RENAMED_DIR="$TEST_DIR/series-renamed"
rm -rf "$TEST_DIR"
mkdir -p "$MOVIE_DIR" "$SERIE_DIR" "$MOVIE_RENAMED_DIR" "$SERIE_RENAMED_DIR"

echo -e "${YELLOW}Creating test file structure...${NC}"
echo ""

# Create series folders with messy names
mkdir -p "$SERIE_DIR/Breaking.Bad"
mkdir -p "$SERIE_DIR/Stranger.Things"

# Create series episode files with typical messy naming
# Put quality tags AFTER episode numbers so scanner can properly extract series name
touch "$SERIE_DIR/Breaking.Bad/Breaking.Bad.S01E01.mkv"
touch "$SERIE_DIR/Breaking.Bad/Breaking.Bad.S01E02.mkv"
touch "$SERIE_DIR/Breaking.Bad/Breaking.Bad.S01E03.mkv"

touch "$SERIE_DIR/Stranger.Things/Stranger.Things.S01E01.mkv"
touch "$SERIE_DIR/Stranger.Things/Stranger.Things.S01E02.mkv"

# Create movie files with messy names
# Note: Year is included in search by cleanMovieName, some may not match
touch "$MOVIE_DIR/The.Dark.Knight.2008.mkv"
mkdir "$MOVIE_DIR/The Matrix/"
touch "$MOVIE_DIR/The Matrix/The Matrix.mkv"

echo -e "${CYAN}Test structure created:${NC}"
echo ""
echo "Movie directory: $MOVIE_DIR"
find "$MOVIE_DIR" -type f | sort | sed 's/^/  /'
echo ""
echo "Series directory: $SERIE_DIR"
find "$SERIE_DIR" -type d -mindepth 1 | sort | sed 's/^/  Folder: /'
find "$SERIE_DIR" -type f | sort | sed 's/^/  File:   /'
echo ""

# =============================================================================
# TEST 1: Dry-run mode (preview without changes)
# =============================================================================
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST 1: Dry-run mode (preview renames)${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""

echo "Running: ./kodi-renamer -movie-to-rename $MOVIE_DIR -serie-to-rename $SERIE_DIR -dry-run -auto"
echo ""
./kodi-renamer \
	-movie-to-rename "$MOVIE_DIR" \
	-movie-renamed "$MOVIE_RENAMED_DIR" \
	-serie-to-rename "$SERIE_DIR" \
	-serie-renamed "$SERIE_RENAMED_DIR" \
	-dry-run \
	-auto 2>&1 | tee /tmp/dryrun_output.log

echo ""
echo -e "${BLUE}Verifying files are unchanged after dry-run...${NC}"
if [[ -f "$SERIE_DIR/Breaking.Bad/Breaking.Bad.S01E01.mkv" ]] && [[ -f "$MOVIE_DIR/The.Dark.Knight.2008.mkv" ]]; then
	echo -e "${GREEN}✓ Files unchanged (dry-run worked correctly)${NC}"
else
	echo -e "${RED}✗ Files were modified during dry-run!${NC}"
	exit 1
fi
echo ""

# =============================================================================
# TEST 2: Actual renaming with auto mode
# =============================================================================
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST 2: Actual file renaming (auto mode)${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""

echo "BEFORE renaming:"
echo "================"
echo "Movies:"
find "$MOVIE_DIR" -type f | sort | sed 's/^/  /'
echo "Series:"
find "$SERIE_DIR" -type d -mindepth 1 | sort | sed 's/^/  Folder: /'
find "$SERIE_DIR" -type f | sort | sed 's/^/  File:   /'
echo ""

echo "Running: ./kodi-renamer -movie-to-rename $MOVIE_DIR -serie-to-rename $SERIE_DIR -auto"
echo ""
./kodi-renamer \
	-movie-to-rename "$MOVIE_DIR" \
	-movie-renamed "$MOVIE_RENAMED_DIR" \
	-serie-to-rename "$SERIE_DIR" \
	-serie-renamed "$SERIE_RENAMED_DIR" \
	-auto 2>&1 | tee /tmp/rename_output.log

echo ""
echo -e "${BLUE}AFTER renaming:${NC}"
echo "==============="
echo "Movies:"
find "$MOVIE_RENAMED_DIR" -type f | sort | sed 's/^/  /'
echo "Series:"
find "$SERIE_RENAMED_DIR" -type d -mindepth 1 | sort | sed 's/^/  Folder: /'
find "$SERIE_RENAMED_DIR" -type f | sort | sed 's/^/  File:   /'
echo ""

# =============================================================================
# TEST 3: Verify expected renames
# =============================================================================
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo -e "${YELLOW}TEST 3: Verify renames are correct${NC}"
echo -e "${YELLOW}═══════════════════════════════════════════════════════${NC}"
echo ""

PASS_COUNT=0
FAIL_COUNT=0

# Check series folders were renamed with year
echo -e "${CYAN}Checking series folder renames:${NC}"

if [[ -d "$SERIE_RENAMED_DIR/Breaking Bad (2008)" ]]; then
	echo -e "  ${GREEN}✓${NC} Breaking Bad (2008)/ folder exists"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Breaking Bad (2008)/ folder NOT found"
	((FAIL_COUNT++))
fi

if [[ -d "$SERIE_RENAMED_DIR/Stranger Things (2016)" ]]; then
	echo -e "  ${GREEN}✓${NC} Stranger Things (2016)/ folder exists"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Stranger Things (2016)/ folder NOT found"
	((FAIL_COUNT++))
fi

echo ""
echo -e "${CYAN}Checking series episode renames:${NC}"

# Check Breaking Bad episodes
if find "$SERIE_RENAMED_DIR/Breaking Bad"* -name "Breaking Bad S01E01 - Pilot.mkv" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} Breaking Bad S01E01 - Pilot.mkv"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Breaking Bad S01E01 - Pilot.mkv NOT found"
	((FAIL_COUNT++))
fi

if find "$SERIE_RENAMED_DIR/Breaking Bad"* -name "Breaking Bad S01E02*.mkv" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} Breaking Bad S01E02 episode renamed"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Breaking Bad S01E02 NOT renamed"
	((FAIL_COUNT++))
fi

# Check Stranger Things episodes
if find "$SERIE_RENAMED_DIR/Stranger Things"* -name "Stranger Things S01E01 - *.mkv" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} Stranger Things S01E01 episode renamed"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Stranger Things S01E01 NOT renamed"
	((FAIL_COUNT++))
fi

if find "$SERIE_RENAMED_DIR/Stranger Things"* -name "Stranger Things S01E02 - *.mkv" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} Stranger Things S01E02 episode renamed"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Stranger Things S01E02 NOT renamed"
	((FAIL_COUNT++))
fi

echo ""
echo -e "${CYAN}Checking movie renames:${NC}"

# Check movies
if [[ -f "$MOVIE_RENAMED_DIR/The Dark Knight (2008)/The Dark Knight (2008).mkv" ]]; then
	echo -e "  ${GREEN}✓${NC} The Dark Knight (2008).mkv"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} The Dark Knight (2008).mkv NOT found"
	((FAIL_COUNT++))
fi

if [[ -f "$MOVIE_RENAMED_DIR/The Matrix (1999)/The Matrix (1999).mkv" ]]; then
	echo -e "  ${GREEN}✓${NC} Matrix (1999).mkv"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} The Matrix (1999).mkv NOT found"
	((FAIL_COUNT++))
fi

echo ""
echo -e "${CYAN}Checking old files are gone:${NC}"

if [[ ! -d "$SERIE_RENAMED_DIR/Breaking.Bad" ]] && [[ ! -d "$SERIE_RENAMED_DIR/Stranger.Things" ]]; then
	echo -e "  ${GREEN}✓${NC} Old series folders removed"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Some old series folders still exist"
	((FAIL_COUNT++))
fi

# Check that files follow Kodi naming convention
if find "$SERIE_RENAMED_DIR" -name "*S[0-9][0-9]E[0-9][0-9] - *.mkv" 2>/dev/null | grep -q .; then
	echo -e "  ${GREEN}✓${NC} Episode files follow Kodi naming convention"
	((PASS_COUNT++))
else
	echo -e "  ${RED}✗${NC} Episode files not properly formatted"
	((FAIL_COUNT++))
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
	echo -e "${GREEN}║  ✓ ALL INTEGRATION TESTS PASSED! ✓   ║${NC}"
	echo -e "${GREEN}╚═══════════════════════════════════════╝${NC}"
	echo ""
	echo "Final structure:"
	find "$TEST_DIR" -type f -o -type d | sort | sed 's/^/  /'
else
	echo -e "${RED}╔═══════════════════════════════════════╗${NC}"
	echo -e "${RED}║  ✗ SOME TESTS FAILED ✗                ║${NC}"
	echo -e "${RED}╚═══════════════════════════════════════╝${NC}"
	echo ""
	echo "Please check the output above for details."
fi
# Cleanup
read -p "Clean up test files? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
	rm -rf "$TEST_DIR"
	echo -e "${GREEN}Test files cleaned up!${NC}"
fi
echo ""
echo "Logs saved to:"
echo "  - /tmp/dryrun_output.log"
echo "  - /tmp/rename_output.log"
echo ""

exit $FAIL_COUNT
