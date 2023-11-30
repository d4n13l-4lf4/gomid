#!/bin/bash
set -e

# Parse arguments
while [[ $# -gt 0 ]]
do
key="$1"
case $key in
    --settings-file)
    SETTINGS_FILE="$2"
    shift # past argument
    shift # past value
    ;;
    --settings-file=*)
    SETTINGS_FILE="${key#*=}"
    shift # past argument
    ;;
    --build-dir)
    BUILD_DIR="$2"
    shift # past argument
    shift # past value
    ;;
    --build-dir=*)
    BUILD_DIR="${key#*=}"
    shift # past argument
    ;;
    --threshold)
    THRESHOLD="$2"
    shift # past argument
    shift # past value
    ;;
    --threshold=*)
    THRESHOLD="${key#*=}"
    shift # past argument
    ;;
    --exclude-dirs)
    EXCLUDE_DIRS=($(echo "$2" | tr ' ' ' '))
    shift # past argument
    shift # past value
    ;;
    --exclude-dirs=*)
    EXCLUDE_DIRS=($(echo "${key#*=}" | tr ' ' ' '))
    shift # past argument
    ;;
    --check-coverage)
    CHECK_COVERAGE=($(echo "$2" | tr ' ' ' '))
    shift # past argument
    shift # past value
    ;;
    --exclude-files)
    EXCLUDE_FILES=($(echo "$2" | tr ' ' ' '))
    shift # past argument
    shift # past value
    ;;
    --exclude-files=*)
    EXCLUDE_FILES=($(echo "${key#*=}" | tr ' ' ' '))
    shift # past argument
    ;;
    *)
    shift # past argument
    ;;
esac
done
THRESHOLD="${THRESHOLD:-85}"
BUILD_DIR="${BUILD_DIR:-build}"
CHECK_COVERAGE=${CHECK_COVERAGE:-true}
# If a settings file is provided, source it
if [[ -n $SETTINGS_FILE && -f $SETTINGS_FILE ]]; then
    source "$SETTINGS_FILE"
fi
EXCLUDE=$(
    IFS="|"
    echo "${EXCLUDE_DIRS[*]//\//\\/}"
)
# Initialize PKGS
PKGS=$(go list ./...)
# Exclude directories only if EXCLUDE is not empty
if [[ -n $EXCLUDE ]]; then
    PKGS=$(echo "$PKGS" | grep -v -E "${EXCLUDE}")
fi
COVERAGE_FILE="${BUILD_DIR}/coverage.out"
# Create build dir
mkdir -p "${BUILD_DIR}"
# Step 1: Run tests only
echo "=== Running tests only"
go test $PKGS -v -race -count=2
# Step 2: Run tests with coverage
echo "=== Running tests with coverage"
go test $PKGS -race -coverprofile="${COVERAGE_FILE}" -covermode=atomic
# Step 3: Show coverage results
go tool cover -func="${COVERAGE_FILE}"
# Step 4: Exclude certain files from coverage
for p in "${EXCLUDE_FILES[@]}"; do
    sed -i.bak "/${p//\//\\/}/d" "${COVERAGE_FILE}"
done
# Check if it should check coverage
if [ $CHECK_COVERAGE = "false" ]; then
    exit $(echo $?)
fi
# Step 5: Check coverage and compare with the threshold
COVERAGE=$(go tool cover -func="${COVERAGE_FILE}" | grep total | awk '{print substr($3, 1, length($3)-1)}')
echo "Total coverage: ${COVERAGE}%. Threshold: ${THRESHOLD}%."
if (($(echo "${COVERAGE} < ${THRESHOLD}" | bc -l))); then
    echo "Code coverage is less than ${THRESHOLD}%. Failing..."
    exit 1
else
    echo "Code coverage is above the threshold. Passing..."
    exit 0
fi