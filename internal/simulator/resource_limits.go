// Copyright 2026 Glassbox Users
// SPDX-License-Identifier: Apache-2.0

package simulator

import (
	"fmt"
)

// ResourceLimits defines canonical resource limits for Soroban protocol versions.
// These constants must match between Rust simulator and Go-side reporting to ensure
// consistent resource accounting and limit exhaustion detection.
//
// Source: Soroban protocol specifications and soroban-env-host defaults.
const (
	// DefaultCPULimit is the default CPU instruction limit for Soroban contracts.
	// Unit: CPU instructions (count)
	// Value: 100,000,000 instructions (100M)
	// Protocol: Soroban Protocol 21+
	DefaultCPULimit uint64 = 100_000_000

	// DefaultMemoryLimit is the default memory byte limit for Soroban contracts.
	// Unit: Bytes
	// Value: 50,000,000 bytes (50 MiB)
	// Protocol: Soroban Protocol 21+
	DefaultMemoryLimit uint64 = 50_000_000

	// MaxOperationsLimit is the maximum number of operations allowed per transaction.
	// Unit: Operation count
	// Value: 100 operations
	// Protocol: Soroban Protocol 21+
	MaxOperationsLimit int = 100
)

// CanonicalResourceUnits defines the canonical units and rounding rules for resource metrics.
type CanonicalResourceUnits struct {
	// CPUUnit is the canonical unit for CPU instructions.
	// Unit: "cpu_instructions" (raw count, no conversion)
	CPUUnit string

	// MemoryUnit is the canonical unit for memory.
	// Unit: "bytes" (raw byte count, no conversion)
	MemoryUnit string

	// OperationsUnit is the canonical unit for operations count.
	// Unit: "operations" (raw count, no conversion)
	OperationsUnit string

	// CPURoundingRule specifies how CPU values should be rounded.
	// Rule: "none" (use exact integer values from soroban-env-host)
	CPURoundingRule string

	// MemoryRoundingRule specifies how memory values should be rounded.
	// Rule: "none" (use exact integer values from soroban-env-host)
	MemoryRoundingRule string

	// PercentagePrecision specifies decimal places for percentage calculations.
	// Rule: 1 decimal place (e.g., 50.5%)
	PercentagePrecision int
}

// GetCanonicalUnits returns the canonical resource units and rounding rules.
func GetCanonicalUnits() *CanonicalResourceUnits {
	return &CanonicalResourceUnits{
		CPUUnit:           "cpu_instructions",
		MemoryUnit:        "bytes",
		OperationsUnit:    "operations",
		CPURoundingRule:   "none",
		MemoryRoundingRule: "none",
		PercentagePrecision: 1,
	}
}

// ResourceLimitFixture represents a test fixture with known resource limits and expected behavior.
type ResourceLimitFixture struct {
	// Name identifies this fixture
	Name string

	// ProtocolVersion is the Soroban protocol version
	ProtocolVersion uint32

	// CPULimit is the CPU instruction limit for this protocol
	CPULimit uint64

	// MemoryLimit is the memory byte limit for this protocol
	MemoryLimit uint64

	// OperationsLimit is the operations count limit
	OperationsLimit int

	// ExpectedCPUUsage is the expected CPU consumption for the fixture
	ExpectedCPUUsage uint64

	// ExpectedMemoryUsage is the expected memory consumption for the fixture
	ExpectedMemoryUsage uint64

	// ExpectedOperationsCount is the expected operations count
	ExpectedOperationsCount int

	// ShouldExceedCPU indicates whether this fixture should exceed CPU limits
	ShouldExceedCPU bool

	// ShouldExceedMemory indicates whether this fixture should exceed memory limits
	ShouldExceedMemory bool

	// ShouldExceedOperations indicates whether this fixture should exceed operations limits
	ShouldExceedOperations bool
}

// GetCanonicalFixtures returns the canonical resource limit fixtures for testing.
func GetCanonicalFixtures() []*ResourceLimitFixture {
	return []*ResourceLimitFixture{
		{
			Name:                  "exact-boundary-cpu",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      DefaultCPULimit,        // Exactly at limit
			ExpectedMemoryUsage:   25_000_000,             // 50% of limit
			ExpectedOperationsCount: 10,
			ShouldExceedCPU:       false,
			ShouldExceedMemory:    false,
			ShouldExceedOperations: false,
		},
		{
			Name:                  "exact-boundary-memory",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      50_000_000,             // 50% of limit
			ExpectedMemoryUsage:   DefaultMemoryLimit,     // Exactly at limit
			ExpectedOperationsCount: 10,
			ShouldExceedCPU:       false,
			ShouldExceedMemory:    false,
			ShouldExceedOperations: false,
		},
		{
			Name:                  "one-over-cpu",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      DefaultCPULimit + 1,    // One instruction over limit
			ExpectedMemoryUsage:   25_000_000,             // 50% of limit
			ExpectedOperationsCount: 10,
			ShouldExceedCPU:       true,
			ShouldExceedMemory:    false,
			ShouldExceedOperations: false,
		},
		{
			Name:                  "one-over-memory",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      50_000_000,             // 50% of limit
			ExpectedMemoryUsage:   DefaultMemoryLimit + 1,  // One byte over limit
			ExpectedOperationsCount: 10,
			ShouldExceedCPU:       false,
			ShouldExceedMemory:    true,
			ShouldExceedOperations: false,
		},
		{
			Name:                  "one-over-operations",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      50_000_000,             // 50% of limit
			ExpectedMemoryUsage:   25_000_000,             // 50% of limit
			ExpectedOperationsCount: MaxOperationsLimit + 1, // One operation over limit
			ShouldExceedCPU:       false,
			ShouldExceedMemory:    false,
			ShouldExceedOperations: true,
		},
		{
			Name:                  "low-usage",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      5_000_000,              // 5% of limit
			ExpectedMemoryUsage:   2_500_000,              // 5% of limit
			ExpectedOperationsCount: 5,
			ShouldExceedCPU:       false,
			ShouldExceedMemory:    false,
			ShouldExceedOperations: false,
		},
		{
			Name:                  "high-usage-within-limits",
			ProtocolVersion:       21,
			CPULimit:              DefaultCPULimit,
			MemoryLimit:           DefaultMemoryLimit,
			OperationsLimit:       MaxOperationsLimit,
			ExpectedCPUUsage:      95_000_000,             // 95% of limit
			ExpectedMemoryUsage:   47_500_000,             // 95% of limit
			ExpectedOperationsCount: 95,
			ShouldExceedCPU:       false,
			ShouldExceedMemory:    false,
			ShouldExceedOperations: false,
		},
	}
}

// ValidateResourceAgreement compares simulator counters with decoded transaction limits.
// Returns a diagnostic if values disagree, naming both observed and expected values.
func ValidateResourceAgreement(
	simulatorCPU uint64,
	simulatorMemory uint64,
	simulatorOps int,
	transactionCPULimit uint64,
	transactionMemoryLimit uint64,
	transactionOpsLimit int,
) *ResourceMismatchDiagnostic {
	units := GetCanonicalUnits()
	
	diagnostics := &ResourceMismatchDiagnostic{
		CanonicalUnits: units,
		SimulatorCPU: simulatorCPU,
		SimulatorMemory: simulatorMemory,
		SimulatorOps: simulatorOps,
		TransactionCPULimit: transactionCPULimit,
		TransactionMemoryLimit: transactionMemoryLimit,
		TransactionOpsLimit: transactionOpsLimit,
	}

	// Check CPU agreement
	if transactionCPULimit > 0 && simulatorCPU > transactionCPULimit {
		diagnostics.CPUExceeded = true
		diagnostics.CPUExceededBy = simulatorCPU - transactionCPULimit
		diagnostics.FirstExceeded = "cpu"
	}

	// Check Memory agreement
	if transactionMemoryLimit > 0 && simulatorMemory > transactionMemoryLimit {
		diagnostics.MemoryExceeded = true
		diagnostics.MemoryExceededBy = simulatorMemory - transactionMemoryLimit
		if diagnostics.FirstExceeded == "" {
			diagnostics.FirstExceeded = "memory"
		}
	}

	// Check Operations agreement
	if transactionOpsLimit > 0 && simulatorOps > transactionOpsLimit {
		diagnostics.OperationsExceeded = true
		diagnostics.OperationsExceededBy = simulatorOps - transactionOpsLimit
		if diagnostics.FirstExceeded == "" {
			diagnostics.FirstExceeded = "operations"
		}
	}

	return diagnostics
}

// ResourceMismatchDiagnostic provides detailed information about resource limit disagreements.
type ResourceMismatchDiagnostic struct {
	// CanonicalUnits specifies the units used for measurement
	CanonicalUnits *CanonicalResourceUnits

	// SimulatorCPU is the CPU consumption reported by the simulator
	SimulatorCPU uint64

	// SimulatorMemory is the memory consumption reported by the simulator
	SimulatorMemory uint64

	// SimulatorOps is the operations count reported by the simulator
	SimulatorOps int

	// TransactionCPULimit is the CPU limit from the transaction metadata
	TransactionCPULimit uint64

	// TransactionMemoryLimit is the memory limit from the transaction metadata
	TransactionMemoryLimit uint64

	// TransactionOpsLimit is the operations limit from the transaction metadata
	TransactionOpsLimit int

	// CPUExceeded indicates whether CPU limits were exceeded
	CPUExceeded bool

	// CPUExceededBy is the amount by which CPU exceeded the limit (if exceeded)
	CPUExceededBy uint64

	// MemoryExceeded indicates whether memory limits were exceeded
	MemoryExceeded bool

	// MemoryExceededBy is the amount by which memory exceeded the limit (if exceeded)
	MemoryExceededBy uint64

	// OperationsExceeded indicates whether operations limits were exceeded
	OperationsExceeded bool

	// OperationsExceededBy is the amount by which operations exceeded the limit (if exceeded)
	OperationsExceededBy int

	// FirstExceeded names the first resource that exceeded its limit (priority: CPU > Memory > Operations)
	FirstExceeded string
}

// ToSummary returns a human-readable summary of the resource mismatch diagnostic.
func (d *ResourceMismatchDiagnostic) ToSummary() string {
	if d == nil {
		return "Resource agreement: validation not performed"
	}

	if !d.CPUExceeded && !d.MemoryExceeded && !d.OperationsExceeded {
		return "Resource agreement: all limits within bounds"
	}

	var summary string
	if d.CPUExceeded {
		summary += fmt.Sprintf("CPU exceeded by %d %s (observed: %d, limit: %d)",
			d.CPUExceededBy, d.CanonicalUnits.CPUUnit, d.SimulatorCPU, d.TransactionCPULimit)
	}
	if d.MemoryExceeded {
		if summary != "" {
			summary += "; "
		}
		summary += fmt.Sprintf("Memory exceeded by %d %s (observed: %d, limit: %d)",
			d.MemoryExceededBy, d.CanonicalUnits.MemoryUnit, d.SimulatorMemory, d.TransactionMemoryLimit)
	}
	if d.OperationsExceeded {
		if summary != "" {
			summary += "; "
		}
		summary += fmt.Sprintf("Operations exceeded by %d %s (observed: %d, limit: %d)",
			d.OperationsExceededBy, d.CanonicalUnits.OperationsUnit, d.SimulatorOps, d.TransactionOpsLimit)
	}

	if d.FirstExceeded != "" {
		summary += fmt.Sprintf(" | First exceeded: %s", d.FirstExceeded)
	}

	return summary
}

// ClassifyFailure determines whether a failure occurred from limit exhaustion or contract logic.
// Returns the failure classification and a diagnostic message.
func ClassifyFailure(
	diagnostic *ResourceMismatchDiagnostic,
	contractError string,
) (FailureClassification, string) {
	if diagnostic == nil {
		return FailureClassificationUnknown, "Unable to classify failure: no resource diagnostic available"
	}

	if diagnostic.CPUExceeded || diagnostic.MemoryExceeded || diagnostic.OperationsExceeded {
		return FailureClassificationLimitExhaustion, diagnostic.ToSummary()
	}

	if contractError != "" {
		return FailureClassificationContractLogic, fmt.Sprintf("Contract logic failure: %s", contractError)
	}

	return FailureClassificationUnknown, "Unable to classify failure: no limit exhaustion or contract error detected"
}

// FailureClassification represents the type of failure that occurred.
type FailureClassification string

const (
	// FailureClassificationLimitExhaustion indicates failure due to resource limit exhaustion
	FailureClassificationLimitExhaustion FailureClassification = "limit_exhaustion"
	
	// FailureClassificationContractLogic indicates failure due to contract logic (not resource limits)
	FailureClassificationContractLogic FailureClassification = "contract_logic"
	
	// FailureClassificationUnknown indicates failure classification could not be determined
	FailureClassificationUnknown FailureClassification = "unknown"
)
