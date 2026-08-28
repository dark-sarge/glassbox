// Copyright 2026 Erst Users
// SPDX-License-Identifier: Apache-2.0

//! Resource limits and canonical units for cross-language agreement.
//!
//! This module defines canonical resource limits, units, and rounding rules
//! that must match between the Rust simulator and Go-side reporting to ensure
//! consistent resource accounting and limit exhaustion detection.
//!
//! Source: Soroban protocol specifications and soroban-env-host defaults.

/// Default CPU instruction limit for Soroban contracts.
/// Unit: CPU instructions (count)
/// Value: 100,000,000 instructions (100M)
/// Protocol: Soroban Protocol 21+
pub const DEFAULT_CPU_LIMIT: u64 = 100_000_000;

/// Default memory byte limit for Soroban contracts.
/// Unit: Bytes
/// Value: 50,000,000 bytes (50 MiB)
/// Protocol: Soroban Protocol 21+
pub const DEFAULT_MEMORY_LIMIT: u64 = 50_000_000;

/// Maximum number of operations allowed per transaction.
/// Unit: Operation count
/// Value: 100 operations
/// Protocol: Soroban Protocol 21+
pub const MAX_OPERATIONS_LIMIT: usize = 100;

/// Canonical resource units and rounding rules.
#[derive(Debug, Clone)]
pub struct CanonicalResourceUnits {
    /// Canonical unit for CPU instructions.
    /// Unit: "cpu_instructions" (raw count, no conversion)
    pub cpu_unit: String,

    /// Canonical unit for memory.
    /// Unit: "bytes" (raw byte count, no conversion)
    pub memory_unit: String,

    /// Canonical unit for operations count.
    /// Unit: "operations" (raw count, no conversion)
    pub operations_unit: String,

    /// Rounding rule for CPU values.
    /// Rule: "none" (use exact integer values from soroban-env-host)
    pub cpu_rounding_rule: String,

    /// Rounding rule for memory values.
    /// Rule: "none" (use exact integer values from soroban-env-host)
    pub memory_rounding_rule: String,

    /// Decimal places for percentage calculations.
    /// Rule: 1 decimal place (e.g., 50.5%)
    pub percentage_precision: u32,
}

impl Default for CanonicalResourceUnits {
    fn default() -> Self {
        Self {
            cpu_unit: "cpu_instructions".to_string(),
            memory_unit: "bytes".to_string(),
            operations_unit: "operations".to_string(),
            cpu_rounding_rule: "none".to_string(),
            memory_rounding_rule: "none".to_string(),
            percentage_precision: 1,
        }
    }
}

impl CanonicalResourceUnits {
    /// Returns the canonical resource units and rounding rules.
    pub fn get_canonical_units() -> Self {
        Self::default()
    }
}

/// Resource limit fixture for testing with known limits and expected behavior.
#[derive(Debug, Clone)]
pub struct ResourceLimitFixture {
    /// Name identifying this fixture
    pub name: String,

    /// Soroban protocol version
    pub protocol_version: u32,

    /// CPU instruction limit for this protocol
    pub cpu_limit: u64,

    /// Memory byte limit for this protocol
    pub memory_limit: u64,

    /// Operations count limit
    pub operations_limit: usize,

    /// Expected CPU consumption for the fixture
    pub expected_cpu_usage: u64,

    /// Expected memory consumption for the fixture
    pub expected_memory_usage: u64,

    /// Expected operations count
    pub expected_operations_count: usize,

    /// Whether this fixture should exceed CPU limits
    pub should_exceed_cpu: bool,

    /// Whether this fixture should exceed memory limits
    pub should_exceed_memory: bool,

    /// Whether this fixture should exceed operations limits
    pub should_exceed_operations: bool,
}

impl ResourceLimitFixture {
    /// Returns the canonical resource limit fixtures for testing.
    pub fn get_canonical_fixtures() -> Vec<Self> {
        vec![
            Self {
                name: "exact-boundary-cpu".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: DEFAULT_CPU_LIMIT,        // Exactly at limit
                expected_memory_usage: 25_000_000,             // 50% of limit
                expected_operations_count: 10,
                should_exceed_cpu: false,
                should_exceed_memory: false,
                should_exceed_operations: false,
            },
            Self {
                name: "exact-boundary-memory".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: 50_000_000,             // 50% of limit
                expected_memory_usage: DEFAULT_MEMORY_LIMIT,     // Exactly at limit
                expected_operations_count: 10,
                should_exceed_cpu: false,
                should_exceed_memory: false,
                should_exceed_operations: false,
            },
            Self {
                name: "one-over-cpu".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: DEFAULT_CPU_LIMIT + 1,    // One instruction over limit
                expected_memory_usage: 25_000_000,             // 50% of limit
                expected_operations_count: 10,
                should_exceed_cpu: true,
                should_exceed_memory: false,
                should_exceed_operations: false,
            },
            Self {
                name: "one-over-memory".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: 50_000_000,             // 50% of limit
                expected_memory_usage: DEFAULT_MEMORY_LIMIT + 1,  // One byte over limit
                expected_operations_count: 10,
                should_exceed_cpu: false,
                should_exceed_memory: true,
                should_exceed_operations: false,
            },
            Self {
                name: "one-over-operations".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: 50_000_000,             // 50% of limit
                expected_memory_usage: 25_000_000,             // 50% of limit
                expected_operations_count: MAX_OPERATIONS_LIMIT + 1, // One operation over limit
                should_exceed_cpu: false,
                should_exceed_memory: false,
                should_exceed_operations: true,
            },
            Self {
                name: "low-usage".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: 5_000_000,              // 5% of limit
                expected_memory_usage: 2_500_000,              // 5% of limit
                expected_operations_count: 5,
                should_exceed_cpu: false,
                should_exceed_memory: false,
                should_exceed_operations: false,
            },
            Self {
                name: "high-usage-within-limits".to_string(),
                protocol_version: 21,
                cpu_limit: DEFAULT_CPU_LIMIT,
                memory_limit: DEFAULT_MEMORY_LIMIT,
                operations_limit: MAX_OPERATIONS_LIMIT,
                expected_cpu_usage: 95_000_000,             // 95% of limit
                expected_memory_usage: 47_500_000,             // 95% of limit
                expected_operations_count: 95,
                should_exceed_cpu: false,
                should_exceed_memory: false,
                should_exceed_operations: false,
            },
        ]
    }
}

/// Diagnostic information about resource limit disagreements.
#[derive(Debug, Clone)]
pub struct ResourceMismatchDiagnostic {
    /// Canonical units used for measurement
    pub canonical_units: CanonicalResourceUnits,

    /// CPU consumption reported by the simulator
    pub simulator_cpu: u64,

    /// Memory consumption reported by the simulator
    pub simulator_memory: u64,

    /// Operations count reported by the simulator
    pub simulator_ops: usize,

    /// CPU limit from the transaction metadata
    pub transaction_cpu_limit: u64,

    /// Memory limit from the transaction metadata
    pub transaction_memory_limit: u64,

    /// Operations limit from the transaction metadata
    pub transaction_ops_limit: usize,

    /// Whether CPU limits were exceeded
    pub cpu_exceeded: bool,

    /// Amount by which CPU exceeded the limit (if exceeded)
    pub cpu_exceeded_by: u64,

    /// Whether memory limits were exceeded
    pub memory_exceeded: bool,

    /// Amount by which memory exceeded the limit (if exceeded)
    pub memory_exceeded_by: u64,

    /// Whether operations limits were exceeded
    pub operations_exceeded: bool,

    /// Amount by which operations exceeded the limit (if exceeded)
    pub operations_exceeded_by: usize,

    /// First resource that exceeded its limit (priority: CPU > Memory > Operations)
    pub first_exceeded: String,
}

impl ResourceMismatchDiagnostic {
    /// Compares simulator counters with decoded transaction limits.
    /// Returns a diagnostic if values disagree, naming both observed and expected values.
    pub fn validate_resource_agreement(
        simulator_cpu: u64,
        simulator_memory: u64,
        simulator_ops: usize,
        transaction_cpu_limit: u64,
        transaction_memory_limit: u64,
        transaction_ops_limit: usize,
    ) -> Self {
        let units = CanonicalResourceUnits::get_canonical_units();

        let mut diagnostics = Self {
            canonical_units: units,
            simulator_cpu,
            simulator_memory,
            simulator_ops,
            transaction_cpu_limit,
            transaction_memory_limit,
            transaction_ops_limit,
            cpu_exceeded: false,
            cpu_exceeded_by: 0,
            memory_exceeded: false,
            memory_exceeded_by: 0,
            operations_exceeded: false,
            operations_exceeded_by: 0,
            first_exceeded: String::new(),
        };

        // Check CPU agreement
        if transaction_cpu_limit > 0 && simulator_cpu > transaction_cpu_limit {
            diagnostics.cpu_exceeded = true;
            diagnostics.cpu_exceeded_by = simulator_cpu - transaction_cpu_limit;
            diagnostics.first_exceeded = "cpu".to_string();
        }

        // Check Memory agreement
        if transaction_memory_limit > 0 && simulator_memory > transaction_memory_limit {
            diagnostics.memory_exceeded = true;
            diagnostics.memory_exceeded_by = simulator_memory - transaction_memory_limit;
            if diagnostics.first_exceeded.is_empty() {
                diagnostics.first_exceeded = "memory".to_string();
            }
        }

        // Check Operations agreement
        if transaction_ops_limit > 0 && simulator_ops > transaction_ops_limit {
            diagnostics.operations_exceeded = true;
            diagnostics.operations_exceeded_by = simulator_ops - transaction_ops_limit;
            if diagnostics.first_exceeded.is_empty() {
                diagnostics.first_exceeded = "operations".to_string();
            }
        }

        diagnostics
    }

    /// Returns a human-readable summary of the resource mismatch diagnostic.
    pub fn to_summary(&self) -> String {
        if !self.cpu_exceeded && !self.memory_exceeded && !self.operations_exceeded {
            return "Resource agreement: all limits within bounds".to_string();
        }

        let mut summary = String::new();

        if self.cpu_exceeded {
            summary.push_str(&format!(
                "CPU exceeded by {} {} (observed: {}, limit: {})",
                self.cpu_exceeded_by, self.canonical_units.cpu_unit,
                self.simulator_cpu, self.transaction_cpu_limit
            ));
        }

        if self.memory_exceeded {
            if !summary.is_empty() {
                summary.push_str("; ");
            }
            summary.push_str(&format!(
                "Memory exceeded by {} {} (observed: {}, limit: {})",
                self.memory_exceeded_by, self.canonical_units.memory_unit,
                self.simulator_memory, self.transaction_memory_limit
            ));
        }

        if self.operations_exceeded {
            if !summary.is_empty() {
                summary.push_str("; ");
            }
            summary.push_str(&format!(
                "Operations exceeded by {} {} (observed: {}, limit: {})",
                self.operations_exceeded_by, self.canonical_units.operations_unit,
                self.simulator_ops, self.transaction_ops_limit
            ));
        }

        if !self.first_exceeded.is_empty() {
            summary.push_str(&format!(" | First exceeded: {}", self.first_exceeded));
        }

        summary
    }
}

/// Failure classification type.
#[derive(Debug, Clone, PartialEq)]
pub enum FailureClassification {
    /// Failure due to resource limit exhaustion
    LimitExhaustion,
    /// Failure due to contract logic (not resource limits)
    ContractLogic,
    /// Failure classification could not be determined
    Unknown,
}

/// Classifies whether a failure occurred from limit exhaustion or contract logic.
/// Returns the failure classification and a diagnostic message.
pub fn classify_failure(
    diagnostic: &ResourceMismatchDiagnostic,
    contract_error: Option<&str>,
) -> (FailureClassification, String) {
    if diagnostic.cpu_exceeded || diagnostic.memory_exceeded || diagnostic.operations_exceeded {
        return (FailureClassification::LimitExhaustion, diagnostic.to_summary());
    }

    if let Some(error) = contract_error {
        return (FailureClassification::ContractLogic, format!("Contract logic failure: {}", error));
    }

    (FailureClassification::Unknown, "Unable to classify failure: no limit exhaustion or contract error detected".to_string())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_canonical_units_default() {
        let units = CanonicalResourceUnits::get_canonical_units();
        assert_eq!(units.cpu_unit, "cpu_instructions");
        assert_eq!(units.memory_unit, "bytes");
        assert_eq!(units.operations_unit, "operations");
        assert_eq!(units.cpu_rounding_rule, "none");
        assert_eq!(units.memory_rounding_rule, "none");
        assert_eq!(units.percentage_precision, 1);
    }

    #[test]
    fn test_resource_limit_fixtures() {
        let fixtures = ResourceLimitFixture::get_canonical_fixtures();
        assert_eq!(fixtures.len(), 7);

        // Check exact boundary CPU fixture
        let cpu_boundary = fixtures.iter().find(|f| f.name == "exact-boundary-cpu").unwrap();
        assert_eq!(cpu_boundary.expected_cpu_usage, DEFAULT_CPU_LIMIT);
        assert!(!cpu_boundary.should_exceed_cpu);

        // Check one over CPU fixture
        let cpu_over = fixtures.iter().find(|f| f.name == "one-over-cpu").unwrap();
        assert_eq!(cpu_over.expected_cpu_usage, DEFAULT_CPU_LIMIT + 1);
        assert!(cpu_over.should_exceed_cpu);
    }

    #[test]
    fn test_validate_resource_agreement_within_limits() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            50_000_000,  // 50% CPU
            25_000_000,  // 50% Memory
            10,          // 10 ops
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        assert!(!diagnostic.cpu_exceeded);
        assert!(!diagnostic.memory_exceeded);
        assert!(!diagnostic.operations_exceeded);
        assert_eq!(diagnostic.to_summary(), "Resource agreement: all limits within bounds");
    }

    #[test]
    fn test_validate_resource_agreement_cpu_exceeded() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            DEFAULT_CPU_LIMIT + 1,  // One over CPU
            25_000_000,
            10,
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        assert!(diagnostic.cpu_exceeded);
        assert_eq!(diagnostic.cpu_exceeded_by, 1);
        assert_eq!(diagnostic.first_exceeded, "cpu");
        assert!(diagnostic.to_summary().contains("CPU exceeded by 1"));
    }

    #[test]
    fn test_validate_resource_agreement_memory_exceeded() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            50_000_000,
            DEFAULT_MEMORY_LIMIT + 1,  // One over Memory
            10,
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        assert!(diagnostic.memory_exceeded);
        assert_eq!(diagnostic.memory_exceeded_by, 1);
        assert_eq!(diagnostic.first_exceeded, "memory");
        assert!(diagnostic.to_summary().contains("Memory exceeded by 1"));
    }

    #[test]
    fn test_validate_resource_agreement_operations_exceeded() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            50_000_000,
            25_000_000,
            MAX_OPERATIONS_LIMIT + 1,  // One over operations
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        assert!(diagnostic.operations_exceeded);
        assert_eq!(diagnostic.operations_exceeded_by, 1);
        assert_eq!(diagnostic.first_exceeded, "operations");
        assert!(diagnostic.to_summary().contains("Operations exceeded by 1"));
    }

    #[test]
    fn test_classify_failure_limit_exhaustion() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            DEFAULT_CPU_LIMIT + 1,
            25_000_000,
            10,
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        let (classification, message) = classify_failure(&diagnostic, None);
        assert_eq!(classification, FailureClassification::LimitExhaustion);
        assert!(message.contains("CPU exceeded"));
    }

    #[test]
    fn test_classify_failure_contract_logic() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            50_000_000,
            25_000_000,
            10,
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        let (classification, message) = classify_failure(&diagnostic, Some("insufficient balance"));
        assert_eq!(classification, FailureClassification::ContractLogic);
        assert!(message.contains("insufficient balance"));
    }

    #[test]
    fn test_classify_failure_unknown() {
        let diagnostic = ResourceMismatchDiagnostic::validate_resource_agreement(
            50_000_000,
            25_000_000,
            10,
            DEFAULT_CPU_LIMIT,
            DEFAULT_MEMORY_LIMIT,
            MAX_OPERATIONS_LIMIT,
        );

        let (classification, message) = classify_failure(&diagnostic, None);
        assert_eq!(classification, FailureClassification::Unknown);
        assert!(message.contains("Unable to classify failure"));
    }
}
