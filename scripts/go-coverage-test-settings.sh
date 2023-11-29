# Custom settings for go-coverage-test pre-commit hook
# Threshold for minimum code coverage percentage (default: 85)
THRESHOLD=100
# Directory to store the coverage output file (default: build)
BUILD_DIR=build
# Directories to exclude from code coverage
EXCLUDE_DIRS=(
    "internal"
    "vendor"
)
# Files to exclude from code coverage
EXCLUDE_FILES=(
    "mocks.go"
)