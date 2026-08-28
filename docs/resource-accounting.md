# Resource Accounting Agreement

This document describes the canonical resource accounting agreement between the Rust simulator and Go-side reporting to ensure consistent CPU, memory, operations, and budget limits across both implementations.

## Overview

The Rust simulator and Go-side reporting must agree on resource metrics to prevent divergence that could cause a locally successful replay to disagree with on-chain execution. This document defines:

- Canonical units and rounding rules
- Shared resource limit constants
- Test fixtures for boundary conditions
- Diagnostic reporting for limit exhaustion

## Canonical Resource Limits

### Constants

Both implementations use these canonical limits for Soroban Protocol 21+:

| Resource | Limit | Unit | Protocol |
|----------|-------|------|----------|
| CPU Instructions | 100,000,000 | cpu_instructions | Soroban Protocol 21+ |
| Memory | 50,000,000 | bytes | Soroban Protocol 21+ |
| Operations | 100 | operations | Soroban Protocol 21+ |

**Go Implementation:** `internal/simulator/resource_limits.go`
```go
const (
    DefaultCPULimit      uint64 = 100_000_000
    DefaultMemoryLimit   uint64 = 50_000_000
    MaxOperationsLimit   int    = 100
)
```

**Rust Implementation:** `simulator/src/resource_limits.rs`
```rust
pub const DEFAULT_CPU_LIMIT: u64 = 100_000_000;
pub const DEFAULT_MEMORY_LIMIT: u64 = 50_000_000;
pub const MAX_OPERATIONS_LIMIT: usize = 100;
```

## Canonical Units and Rounding Rules

### Units

- **CPU Unit:** `cpu_instructions` (raw count, no conversion)
- **Memory Unit:** `bytes` (raw byte count, no conversion)  
- **Operations Unit:** `operations` (raw count, no conversion)

### Rounding Rules

- **CPU Rounding:** `none` (use exact integer values from soroban-env-host)
- **Memory Rounding:** `none` (use exact integer values from soroban-env-host)
- **Percentage Precision:** 1 decimal place (e.g., 50.5%)

### Implementation

**Go:**
```go
type CanonicalResourceUnits struct {
    CPUUnit           string
    MemoryUnit        string
    OperationsUnit    string
    CPURoundingRule   string
    MemoryRoundingRule string
    PercentagePrecision int
}
```

**Rust:**
```rust
pub struct CanonicalResourceUnits {
    pub cpu_unit: String,
    pub memory_unit: String,
    pub operations_unit: String,
    pub cpu_rounding_rule: String,
    pub memory_rounding_rule: String,
    pub percentage_precision: u32,
}
```

## Resource Limit Fixtures

Canonical test fixtures ensure both implementations agree on boundary conditions:

### Fixture Categories

1. **Exact Boundary Tests:** Verify behavior at exactly the limit
2. **One-Over-Limit Tests:** Verify detection of limit exhaustion
3. **Low Usage Tests:** Verify normal operation with minimal resource usage
4. **High Usage Tests:** Verify behavior near but within limits

### Available Fixtures

| Fixture Name | CPU Usage | Memory Usage | Operations | Expected Behavior |
|--------------|-----------|-------------|------------|-------------------|
| `exact-boundary-cpu` | 100,000,000 | 25,000,000 | 10 | Within limits |
| `exact-boundary-memory` | 50,000,000 | 50,000,000 | 10 | Within limits |
| `one-over-cpu` | 100,000,001 | 25,000,000 | 10 | CPU exceeded |
| `one-over-memory` | 50,000,000 | 50,000,001 | 10 | Memory exceeded |
| `one-over-operations` | 50,000,000 | 25,000,000 | 101 | Operations exceeded |
| `low-usage` | 5,000,000 | 2,500,000 | 5 | Within limits |
| `high-usage-within-limits` | 95,000,000 | 47,500,000 | 95 | Within limits |

## Resource Validation

### Go Implementation

```go
diagnostic := ValidateResourceAgreement(
    simulatorCPU,      // uint64 from simulator
    simulatorMemory,   // uint64 from simulator
    simulatorOps,      // int from simulator
    transactionCPULimit,    // uint64 from transaction metadata
    transactionMemoryLimit, // uint64 from transaction metadata
    transactionOpsLimit,    // int from transaction metadata
)

summary := diagnostic.ToSummary()
classification, message := ClassifyFailure(diagnostic, contractError)
```

### Rust Implementation

```rust
let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
    simulator_cpu,
    simulator_memory,
    simulator_ops,
    transaction_cpu_limit,
    transaction_memory_limit,
    transaction_ops_limit,
);

let summary = diagnostic.to_summary();
let (classification, message) = classify_failure(&diagnostic, contract_error);
```

## Event Schema Updates

Both event schemas now include resource diagnostic information:

### Go DiagnosticEvent

```go
type DiagnosticEvent struct {
    // ... existing fields ...
    ResourceDiagnostic *ResourceDiagnostic `json:"resource_diagnostic,omitempty"`
}

type ResourceDiagnostic struct {
    CPUInstructions uint64 `json:"cpu_instructions"`
    MemoryBytes     uint64 `json:"memory_bytes"`
    CPULimit        uint64 `json:"cpu_limit"`
    MemoryLimit     uint64 `json:"memory_limit"`
    CPUExceeded     bool   `json:"cpu_exceeded"`
    MemoryExceeded  bool   `json:"memory_exceeded"`
    FirstExceeded   *string `json:"first_exceeded,omitempty"`
}
```

### Rust DiagnosticEvent

```rust
pub struct DiagnosticEvent {
    // ... existing fields ...
    pub resource_diagnostic: Option<ResourceDiagnostic>,
}

pub struct ResourceDiagnostic {
    pub cpu_instructions: u64,
    pub memory_bytes: u64,
    pub cpu_limit: u64,
    pub memory_limit: u64,
    pub cpu_exceeded: bool,
    pub memory_exceeded: bool,
    pub first_exceeded: Option<String>,
}
```

## Failure Classification

Failures are classified into three categories:

1. **LimitExhaustion:** Failure due to resource limit exhaustion (CPU, memory, or operations)
2. **ContractLogic:** Failure due to contract logic (not resource limits)
3. **Unknown:** Failure classification could not be determined

### Priority Order

When multiple resources exceed limits, the first exceeded resource follows this priority:
1. CPU (highest priority)
2. Memory
3. Operations (lowest priority)

## Testing

### Running Tests

**Go Tests:**
```bash
cd internal/simulator
go test -v -run TestResourceLimits
```

**Rust Tests:**
```bash
cd simulator
cargo test --lib resource_limits
```

### Test Coverage

- Canonical units validation
- Fixture agreement tests
- Exact boundary tests
- One-over-limit tests
- Resource validation logic
- Failure classification
- Event schema compatibility

## Troubleshooting

### Common Issues

#### 1. Divergent CPU Counts

**Symptom:** Rust simulator reports different CPU instructions than Go-side decoding.

**Diagnosis:**
```go
diagnostic := ValidateResourceAgreement(rustCPU, rustMemory, rustOps, goCPULimit, goMemoryLimit, goOpsLimit)
fmt.Println(diagnostic.ToSummary())
```

**Solution:** Ensure both use `soroban-env-host::Host.budget_cloned().get_cpu_insns_consumed()` (Rust) and equivalent Go-side budget extraction.

#### 2. Memory Unit Mismatch

**Symptom:** Memory values differ by factor of 1024 (KiB vs bytes).

**Diagnosis:** Check that both implementations use raw bytes, not KiB/MiB units.

**Solution:** Verify canonical units specify "bytes" with rounding rule "none".

#### 3. Operations Count Disagreement

**Symptom:** Operations count differs between implementations.

**Diagnosis:** Rust counts from envelope operations slice; Go may use different source.

**Solution:** Ensure both count from the same source (transaction envelope operations).

#### 4. Limit Exhaustion Not Detected

**Symptom:** Transaction exceeds limits but no exhaustion diagnostic is generated.

**Diagnosis:** Check that limit values are being passed correctly to validation function.

**Solution:** Verify protocol version is correctly set and limits are being applied from transaction metadata.

## Migration Guide

### Updating Existing Code

#### Rust Simulator

1. Import canonical limits:
```rust
use crate::resource_limits::{DEFAULT_CPU_LIMIT, DEFAULT_MEMORY_LIMIT, MAX_OPERATIONS_LIMIT};
```

2. Replace hardcoded constants:
```rust
// Before
const CPU_LIMIT: u64 = 100_000_000;

// After
use crate::resource_limits::DEFAULT_CPU_LIMIT as CPU_LIMIT;
```

3. Add resource diagnostics to events:
```rust
let diagnostic = ResourceDiagnostic {
    cpu_instructions: cpu_insns,
    memory_bytes: mem_bytes,
    cpu_limit: DEFAULT_CPU_LIMIT,
    memory_limit: DEFAULT_MEMORY_LIMIT,
    cpu_exceeded: cpu_insns > DEFAULT_CPU_LIMIT,
    memory_exceeded: mem_bytes > DEFAULT_MEMORY_LIMIT,
    first_exceeded: if cpu_insns > DEFAULT_CPU_LIMIT {
        Some("cpu".to_string())
    } else if mem_bytes > DEFAULT_MEMORY_LIMIT {
        Some("memory".to_string())
    } else {
        None
    },
};
```

#### Go-Side Reporting

1. Import canonical limits:
```go
import "github.com/dotandev/glassbox/internal/simulator"
```

2. Use validation function:
```go
diagnostic := simulator.ValidateResourceAgreement(
    resp.BudgetUsage.CPUInstructions,
    resp.BudgetUsage.MemoryBytes,
    resp.BudgetUsage.OperationsCount,
    resp.BudgetUsage.CPULimit,
    resp.BudgetUsage.MemoryLimit,
    resp.BudgetUsage.OperationsLimit,
)
```

3. Add resource diagnostics to events:
```go
event.ResourceDiagnostic = &simulator.ResourceDiagnostic{
    CPUInstructions: budget.CPUInstructions,
    MemoryBytes:     budget.MemoryBytes,
    CPULimit:        budget.CPULimit,
    MemoryLimit:     budget.MemoryLimit,
    CPUExceeded:     budget.CPUInstructions >= budget.CPULimit,
    MemoryExceeded:  budget.MemoryBytes >= budget.MemoryLimit,
    FirstExceeded:   getFirstExceeded(budget),
}
```

## References

- Soroban Protocol Specifications: https://github.com/stellar/stellar-protocol
- soroban-env-host Documentation: https://docs.rs/soroban-env-host
- Resource Limits Implementation: `internal/simulator/resource_limits.go`
- Rust Resource Limits: `simulator/src/resource_limits.rs`
