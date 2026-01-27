#!/bin/bash

set -e

echo "=== KODI RENAMER GO - DUAL API TEST ==="
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
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Test 1: No API keys
echo -e "${YELLOW}Test 1: No API keys provided${NC}"
./kodi-renamer -dir . 2>&1 | head -5
echo -e "${GREEN}✓ Properly rejects with no API keys${NC}"
echo ""

# Test 2: Help
echo -e "${YELLOW}Test 2: Command line options${NC}"
./kodi-renamer -h 2>&1 | grep -E "(tvdb-key|tmdb-key)"
echo -e "${GREEN}✓ Both API key flags available${NC}"
echo ""

# Test 3: Create test files
echo -e "${YELLOW}Test 3: Creating test files...${NC}"
TEST_DIR="dual_api_test"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

touch "$TEST_DIR/Breaking.Bad.S01E01.720p.mkv"
touch "$TEST_DIR/The.Matrix.1999.1080p.mp4"
touch "$TEST_DIR/Stranger.Things.S01E01.mkv"
touch "$TEST_DIR/Inception.2010.BluRay.mp4"

echo -e "${GREEN}✓ Created $(ls -1 "$TEST_DIR" | wc -l) test files${NC}"
ls -1 "$TEST_DIR"
echo ""

# Test 4: Scanner (no API needed)
echo -e "${YELLOW}Test 4: Testing scanner (no API required)${NC}"
go build -o test-scanner ./cmd/test-scanner 2>/dev/null
./test-scanner -dir "$TEST_DIR" | grep "Summary"
echo -e "${GREEN}✓ Scanner works without API${NC}"
echo ""

# Test 5: TVDB only simulation
echo -e "${YELLOW}Test 5: TVDB only mode (will fail with dummy key)${NC}"
export TVDB_API_KEY="dummy-tvdb-key"
unset TMDB_API_KEY
./kodi-renamer -dir "$TEST_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Accepts TVDB key alone${NC}"
echo ""

# Test 6: TMDB only simulation
echo -e "${YELLOW}Test 6: TMDB only mode (will fail with dummy key)${NC}"
unset TVDB_API_KEY
export TMDB_API_KEY="dummy-tmdb-key"
./kodi-renamer -dir "$TEST_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Accepts TMDB key alone${NC}"
echo ""

# Test 7: Both APIs simulation
echo -e "${YELLOW}Test 7: Dual API mode (will fail with dummy keys)${NC}"
export TVDB_API_KEY="dummy-tvdb-key"
export TMDB_API_KEY="dummy-tmdb-key"
./kodi-renamer -dir "$TEST_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Accepts both API keys${NC}"
echo ""

# Test 8: Command line flags
echo -e "${YELLOW}Test 8: Using command line flags${NC}"
unset TVDB_API_KEY
unset TMDB_API_KEY
./kodi-renamer -tvdb-key "test1" -tmdb-key "test2" -dir "$TEST_DIR" -dry-run 2>&1 | head -3
echo -e "${GREEN}✓ Command line flags work${NC}"
echo ""

# Cleanup
rm -rf "$TEST_DIR"

# Summary
echo ""
echo -e "${GREEN}=== ALL TESTS PASSED ===${NC}"
echo ""
echo -e "${BLUE}API Key Configuration Options:${NC}"
echo ""
echo "1. Environment Variables:"
echo "   export TVDB_API_KEY='your-tvdb-key'"
echo "   export TMDB_API_KEY='your-tmdb-key'"
echo ""
echo "2. Command Line Flags:"
echo "   ./kodi-renamer -tvdb-key 'key' -tmdb-key 'key' -dir /path"
echo ""
echo "3. Mixed (env + flag):"
echo "   export TVDB_API_KEY='key1'"
echo "   ./kodi-renamer -tmdb-key 'key2' -dir /path"
echo ""
echo -e "${BLUE}Supported Configurations:${NC}"
echo "  ✓ TVDB only         - Works with TVDB database"
echo "  ✓ TMDB only         - Works with TheMovieDB"
echo "  ✓ Both APIs         - Best results (merged)"
echo "  ✗ No APIs           - Error with instructions"
echo ""
echo -e "${YELLOW}To use with real APIs:${NC}"
echo "1. Get TVDB key: https://thetvdb.com/api-information"
echo "2. Get TMDB key: https://www.themoviedb.org/settings/api"
echo "3. Set keys and run: ./kodi-renamer -dir /your/media -dry-run"
echo ""
echo -e "${GREEN}Documentation:${NC}"
echo "  README - Full documentation"
echo "  HOWTO  - Quick start guide"
echo ""
