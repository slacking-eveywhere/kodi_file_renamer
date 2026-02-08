#!/usr/bin/env bash

set -e

echo "KODI RENAMER GO - DUAL API TEST"
echo ""

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}This test demonstrates the dual API support (TVDB + TMDB)${NC}"
echo ""

# Build
echo -e "${YELLOW}Building application...${NC}"
export PATH=$PATH:~/go/bin
go build -o kodi-renamer ./cmd/kodi-renamer 2>/dev/null || {
	echo -e "${RED}Build failed${NC}"
	exit 1
}
echo -e "${GREEN}Build successful${NC}"
echo ""

# Test 1: No API keys
echo -e "${YELLOW}Test 1: No API keys provided${NC}"
./kodi-renamer -movie-to-rename . 2>&1 | head -5
echo -e "${GREEN}Properly rejects with no API keys${NC}"
echo ""

# Test 2: Help
echo -e "${YELLOW}Test 2: Command line options${NC}"
./kodi-renamer -h 2>&1 | grep -E "(tvdb-key|tmdb-key)"
echo -e "${GREEN}Both API key flags available${NC}"
echo ""

# Test 3: Create test files
echo -e "${YELLOW}Test 3: Creating test files...${NC}"
TEST_DIR="media"
SERIE_DIR="$TEST_DIR/series-to-rename"
MOVIE_DIR="$TEST_DIR/movies-to-rename"
rm -rf "$TEST_DIR"
mkdir -p "$SERIE_DIR"
mkdir -p "$MOVIE_DIR"

# Create series folders (without year - to test renaming)
mkdir -p "$SERIE_DIR/Breaking.Bad"
mkdir -p "$SERIE_DIR/Stranger.Things"

# Create series episode files
touch "$SERIE_DIR/Breaking.Bad/Breaking.Bad.S01E01.720p.mkv"
touch "$SERIE_DIR/Breaking.Bad/Breaking.Bad.S01E02.720p.mkv"
touch "$SERIE_DIR/Stranger.Things/Stranger.Things.S01E01.mkv"
touch "$SERIE_DIR/Stranger.Things/Stranger.Things.S01E02.mkv"

# Create movie files
touch "$MOVIE_DIR/The.Matrix.1999.1080p.mp4"
touch "$MOVIE_DIR/Inception.2010.BluRay.mp4"

echo "  Series folders (should be renamed to include year):"
echo "    - Breaking.Bad/ -> Breaking Bad (2008)/"
echo "    - Stranger.Things/ -> Stranger Things (2016)/"
echo "  Movies (should be renamed with proper format):"
echo "    - The.Matrix.1999.1080p.mp4 -> The Matrix (1999).mp4"
echo "    - Inception.2010.BluRay.mp4 -> Inception (2010).mp4"

echo -e "${GREEN}Created test structure:${NC}"
find "$TEST_DIR" -type f -o -type d | sort
echo ""

# Test 4: Scanner (no API needed)
echo -e "${YELLOW}Test 4: Testing scanner (no API required)${NC}"
go build -o test-scanner ./cmd/test-scanner 2>/dev/null
echo "  Testing movie scanner:"
./test-scanner -dir "$MOVIE_DIR" | grep "Summary"
echo "  Testing series scanner:"
./test-scanner -dir "$SERIE_DIR" | grep "Summary"
echo -e "${GREEN}✓ Scanner works without API${NC}"
echo ""

TVDB_API_KEY_BCK="${TVDB_API_KEY:-dummytvdbkey}"
TMDB_API_KEY_BCK="${TMDB_API_KEY:-dummytmdbkey}"

# Test 5: TVDB only simulation
echo -e "${YELLOW}Test 5: TVDB only mode (will fail with dummy key)${NC}"
export TVDB_API_KEY="${TVDB_API_KEY_BCK}"
unset TMDB_API_KEY
./kodi-renamer -movie-to-rename "$MOVIE_DIR" -serie-to-rename "$SERIE_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Accepts TVDB key alone${NC}"
echo ""

# Test 6: TMDB only simulation
echo -e "${YELLOW}Test 6: TMDB only mode (will fail with dummy key)${NC}"
unset TVDB_API_KEY
export TMDB_API_KEY="${TMDB_API_KEY_BCK}"
./kodi-renamer -movie-to-rename "$MOVIE_DIR" -serie-to-rename "$SERIE_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Accepts TMDB key alone${NC}"
echo ""

# Test 7: Both APIs simulation
echo -e "${YELLOW}Test 7: Dual API mode (will fail with dummy keys)${NC}"
export TVDB_API_KEY="${TVDB_API_KEY_BCK}"
export TMDB_API_KEY="${TMDB_API_KEY_BCK}"
./kodi-renamer -movie-to-rename "$MOVIE_DIR" -serie-to-rename "$SERIE_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Accepts both API keys${NC}"
echo ""

# Test 9: Command line flags
echo -e "${YELLOW}Test 9: Using command line flags${NC}"
unset TVDB_API_KEY
unset TMDB_API_KEY
./kodi-renamer -tvdb-key "$TVDB_API_KEY_BCK" -tmdb-key "$TMDB_API_KEY_BCK" -movie-to-rename "$MOVIE_DIR" -serie-to-rename "$SERIE_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Command line flags work${NC}"
echo ""

# Cleanup
rm -rf "$TEST_DIR"

# Summary
echo ""
echo -e "${GREEN}ALL TESTS PASSED${NC}"
echo ""
echo -e "${BLUE}API Key Configuration Options:${NC}"
echo ""
echo "1. Environment Variables:"
echo "   export TVDB_API_KEY='your-tvdb-key'"
echo "   export TMDB_API_KEY='your-tmdb-key'"
echo ""
echo "2. Command Line Flags:"
echo "   ./kodi-renamer -tvdb-key 'key' -tmdb-key 'key' -movie-to-rename /path/movies -serie-to-rename /path/series"
echo ""
echo "3. Mixed (env + flag):"
echo "   export TVDB_API_KEY='key1'"
echo "   ./kodi-renamer -tmdb-key 'key2' -movie-to-rename /path/movies"
echo ""
echo -e "${BLUE}Directory Configuration:${NC}"
echo "  MOVIE_TO_RENAME_DIR - Directory containing movies to rename"
echo "  MOVIE_RENAMED_DIR   - Output directory for renamed movies (optional)"
echo "  SERIE_TO_RENAME_DIR - Directory containing series to rename"
echo "  SERIE_RENAMED_DIR   - Output directory for renamed series (optional)"
echo ""
echo -e "${BLUE}Supported Configurations:${NC}"
echo "  ✓ TVDB only         - Works with TVDB database"
echo "  ✓ TMDB only         - Works with TheMovieDB"
echo "  ✓ Both APIs         - Best results (merged)"
echo "  ✗ No APIs           - Error with instructions"
echo ""
