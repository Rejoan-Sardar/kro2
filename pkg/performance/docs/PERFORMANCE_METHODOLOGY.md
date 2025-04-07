# KRO Performance Testing Methodology

This document describes the methodology used to test KRO performance and generate reproducible, reliable metrics.

## Core Principles

1. **Reproducibility**: All tests should be reproducible with minimal variance between runs
2. **Realistic Workloads**: Tests should simulate real-world usage patterns
3. **Isolation**: Each benchmark should isolate specific components for accurate measurement
4. **Scalability Testing**: Tests should evaluate performance across different scales and loads
5. **Resource Efficiency**: Measure resource consumption (CPU, memory) in addition to latency

## Test Environment

### Simulation Mode

Simulation mode uses in-memory data structures instead of a real Kubernetes cluster. This provides:
- Consistent, reproducible results unaffected by cluster conditions
- Ability to run tests without a Kubernetes cluster
- Isolation of KRO's internal performance from Kubernetes API overhead

### Kubernetes Mode

Tests can also run against a real Kubernetes cluster to measure end-to-end performance:
- Tests run in the specified namespace
- Results include real-world network and API latency
- Resource usage metrics include the complete stack

## Benchmark Implementation

### CRUD Benchmarks

The CRUD benchmarks measure:
1. **Throughput**: Operations per second for each CRUD operation type
2. **Latency**: Average, median, p95, and p99 latency per operation
3. **Resource Consumption**: CPU and memory usage during operations

The test creates a pool of resources and performs operations with configurable:
- Number of concurrent workers
- Number of resources
- Test duration
- Operation mix (% of each operation type)

### CEL Expression Benchmarks

The CEL benchmarks evaluate:
1. **Expression Parsing**: Time to parse and compile expressions
2. **Expression Evaluation**: Time to evaluate expressions against data
3. **Scaling Characteristics**: How performance scales with expression complexity

Tests use expressions of varying complexity:
- Simple: `resource.metadata.name == "test"`
- Medium: `resource.metadata.name.startsWith("test") && resource.spec.replicas > 1`
- Complex: `resource.metadata.labels.exists(l, l.startsWith("env")) && resource.spec.containers.all(c, c.image.contains("latest") || c.resources.limits.exists(r, r > 100))`
- Very Complex: Multi-conditional expressions with functions, regex, and complex logic

### ResourceGraph Benchmarks

The ResourceGraph benchmarks measure:
1. **Graph Building**: Time to construct resource graphs
2. **Graph Traversal**: Performance of traversal operations
3. **Relationship Resolution**: Efficiency of resolving resource relationships

Tests use graphs of varying complexity:
- Small: ~10 nodes, ~15 edges
- Medium: ~20 nodes, ~30 edges
- Large: ~30+ nodes, ~45+ edges

## Metrics Collection

The framework collects:
1. **Timing Metrics**: Using high-precision timers
2. **Operation Counts**: Total operations and operations per second
3. **Resource Usage**: CPU and memory consumption
4. **Errors**: Rate of errors or failures

## Analysis Methodology

Results are analyzed to determine:
1. **Baseline Performance**: Establish baseline metrics for comparison
2. **Scaling Behavior**: How performance scales with load and complexity
3. **Bottlenecks**: Identify performance bottlenecks
4. **Optimization Opportunities**: Areas with potential for improvement

## Visualization

Results are visualized through:
1. **Time Series**: Performance over time
2. **Comparative Charts**: Performance across different configurations
3. **Heatmaps**: Identifying resource usage patterns
4. **Resource Profiles**: CPU and memory usage visualization

## Interpretation Guidelines

When interpreting results:
1. **Relative Performance**: Focus on relative performance between test runs
2. **Trend Analysis**: Look for performance trends over time
3. **Variance Analysis**: Consider the variance between test runs
4. **Context**: Interpret results in the context of real-world workloads
