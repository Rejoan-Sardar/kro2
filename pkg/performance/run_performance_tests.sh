#!/bin/bash

# Default values
DURATION="10s"
WORKERS=4
RESOURCES=100
SIMULATION=true
NAMESPACE="default"
BENCHMARK_TYPE="all"
CEL_COMPLEXITY="all"
GRAPH_COMPLEXITY="all"
NODES=30
EDGES=45
OUTPUT_DIR="./results"
VERBOSE=false
DEMO_MODE=false

# Process command line arguments
while [[ $# -gt 0 ]]; do
  case $1 in
    --duration)
      DURATION="$2"
      shift 2
      ;;
    --workers)
      WORKERS="$2"
      shift 2
      ;;
    --resources)
      RESOURCES="$2"
      shift 2
      ;;
    --no-simulation)
      SIMULATION=false
      shift
      ;;
    --namespace)
      NAMESPACE="$2"
      shift 2
      ;;
    --kubeconfig)
      KUBECONFIG="$2"
      shift 2
      ;;
    --type)
      BENCHMARK_TYPE="$2"
      shift 2
      ;;
    --complexity)
      CEL_COMPLEXITY="$2"
      shift 2
      ;;
    --graph-complexity)
      GRAPH_COMPLEXITY="$2"
      shift 2
      ;;
    --nodes)
      NODES="$2"
      shift 2
      ;;
    --edges)
      EDGES="$2"
      shift 2
      ;;
    --output)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --verbose)
      VERBOSE=true
      shift
      ;;
    --demo)
      DEMO_MODE=true
      shift
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Set up the common flags
COMMON_FLAGS="--duration $DURATION --workers $WORKERS --resources $RESOURCES --namespace $NAMESPACE"
if [ "$SIMULATION" = false ]; then
  COMMON_FLAGS="$COMMON_FLAGS"
else
  COMMON_FLAGS="$COMMON_FLAGS --simulation"
fi

if [ -n "$KUBECONFIG" ]; then
  COMMON_FLAGS="$COMMON_FLAGS --kubeconfig $KUBECONFIG"
fi

if [ "$VERBOSE" = true ]; then
  COMMON_FLAGS="$COMMON_FLAGS --verbose"
fi

# Run in demo mode if requested
if [ "$DEMO_MODE" = true ]; then
  echo "Running in demo mode - showing simulated results"
  echo "==============================================="
  echo ""
  
  # Generate sample data directory
  mkdir -p "$OUTPUT_DIR/visualizations"
  
  echo -e "\033[1;32m===== KRO Performance Test Results =====\033[0m"
  echo ""
  echo -e "\033[1;33mCRUD Operations Performance:\033[0m"
  echo "  Create: 424 ops/sec (16.8ms P95 latency)"
  echo "  Read:   892 ops/sec (9.2ms P95 latency)"
  echo "  Update: 582 ops/sec (14.5ms P95 latency)"
  echo "  Delete: 726 ops/sec (11.4ms P95 latency)"
  echo "  Overall: 658 ops/sec (12.5ms P95 latency)"
  echo ""
  
  echo -e "\033[1;33mCEL Expression Performance:\033[0m"
  echo "  Simple expressions:       9821 ops/sec (0.9ms P95 latency)"
  echo "  Medium expressions:       4876 ops/sec (1.8ms P95 latency)"
  echo "  Complex expressions:      1246 ops/sec (7.2ms P95 latency)"
  echo "  Very complex expressions: 789 ops/sec (11.6ms P95 latency)"
  echo "  Overall: 5433 ops/sec (3.2ms P95 latency)"
  echo ""
  
  echo -e "\033[1;33mResourceGraph Performance:\033[0m"
  echo "  Small graphs (10 nodes):  246 ops/sec (18.2ms P95 latency)"
  echo "  Medium graphs (30 nodes): 148 ops/sec (28.6ms P95 latency)"
  echo "  Large graphs (50 nodes):  92 ops/sec (42.1ms P95 latency)"
  echo "  Overall: 148 ops/sec (28.6ms P95 latency)"
  echo ""
  
  echo -e "\033[1;33mResource Usage Patterns:\033[0m"
  echo "  CRUD operations:     45.2% CPU, 128.5 MB memory"
  echo "  CEL evaluation:      62.0% CPU, 95.7 MB memory"
  echo "  ResourceGraph:       79.0% CPU, 257.2 MB memory"
  echo ""
  
  echo -e "\033[1;33mScaling Behavior:\033[0m"
  echo "  1 worker:  219 ops/sec (21.8% CPU)"
  echo "  2 workers: 412 ops/sec (38.5% CPU)"
  echo "  4 workers: 658 ops/sec (45.2% CPU)"
  echo "  8 workers: 872 ops/sec (72.6% CPU)"
  echo "  16 workers: 983 ops/sec (92.4% CPU)"
  echo ""
  
  echo -e "\033[1;33mKey Insights:\033[0m"
  echo "  1. Read operations (892 ops/sec) are 2.1x faster than create operations (424 ops/sec)"
  echo "  2. Simple CEL expressions (9821 ops/sec) are 7.8x faster than complex ones (1246 ops/sec)"
  echo "  3. ResourceGraph operations are CPU and memory intensive (79% CPU, 257MB)"
  echo "  4. Performance scales well up to 4 workers, with diminishing returns beyond that"
  echo "  5. CEL evaluation is CPU-efficient compared to other operations"
  echo "  6. Large resource graphs (50+ nodes) should be used carefully due to performance impact"
  echo ""
  
  echo -e "\033[1;32m=========================================\033[0m"
  
  exit 0
fi

# Run the benchmark(s) based on the requested type
case "$BENCHMARK_TYPE" in
  "crud")
    echo "Running CRUD benchmarks..."
    go run main.go benchmark $COMMON_FLAGS --type crud --output "$OUTPUT_DIR"
    ;;
  "cel")
    echo "Running CEL benchmarks..."
    go run main.go benchmark $COMMON_FLAGS --type cel --complexity "$CEL_COMPLEXITY" --output "$OUTPUT_DIR"
    ;;
  "resourcegraph")
    echo "Running ResourceGraph benchmarks..."
    go run main.go benchmark $COMMON_FLAGS --type resourcegraph --graph-complexity "$GRAPH_COMPLEXITY" --nodes "$NODES" --edges "$EDGES" --output "$OUTPUT_DIR"
    ;;
  "all")
    echo "Running all benchmarks..."
    go run main.go benchmark $COMMON_FLAGS --type all --complexity "$CEL_COMPLEXITY" --graph-complexity "$GRAPH_COMPLEXITY" --nodes "$NODES" --edges "$EDGES" --output "$OUTPUT_DIR"
    ;;
  *)
    echo "Unknown benchmark type: $BENCHMARK_TYPE"
    exit 1
    ;;
esac

# Run analysis on the results if we're doing all benchmarks
if [ "$BENCHMARK_TYPE" = "all" ]; then
  echo "Analyzing benchmark results..."
  go run main.go analyze --input "$OUTPUT_DIR/all_results.json" --output "$OUTPUT_DIR/analysis.json"
  
  echo "Generating visualizations..."
  go run main.go visualize --input "$OUTPUT_DIR/analysis.json" --output "$OUTPUT_DIR/visualizations" --charts all
fi

echo "Performance testing completed. Results are in $OUTPUT_DIR"
