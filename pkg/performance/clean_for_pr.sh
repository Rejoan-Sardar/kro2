#!/bin/bash

# Script to prepare the KRO Performance package for a Pull Request
# This creates all necessary directories and files required by the GitHub workflow

echo "============ KRO Performance PR Preparation Tool ============"
echo "Current date and time: $(date)"
echo ""

# Create required directories
echo "Creating required directories..."
for dir in "visualizations" "results"; do
  if [ -d "$dir" ]; then
    echo "✓ $dir already exists"
  else
    mkdir -p "$dir"
    echo "✓ $dir created"
    
    # Create .keep file to ensure directory is tracked by Git
    touch "$dir/.keep"
    echo "  Added .keep file to $dir"
  fi
done

# Clean up any previous test results or artifacts
echo "Cleaning up previous artifacts..."
rm -f benchmark_results.json analysis_report.json
rm -f visualizations/*.png

# Make sure run_performance_tests.sh is executable
echo "Setting executable permissions on scripts..."
chmod +x run_performance_tests.sh
if [ -f "./clean_for_pr.sh" ]; then
  chmod +x clean_for_pr.sh
fi

echo ""
echo "PR preparation completed"
echo "The repository is now ready for a Pull Request"
echo "=========================================================="
