// Copyright 2026 Glassbox Users
// SPDX-License-Identifier: Apache-2.0

package simulator

import (
	"testing"
)

func TestCanonicalUnits(t *testing.T) {
	units := GetCanonicalUnits()

	if units.CPUUnit != "cpu_instructions" {
		t.Errorf("Expected CPUUnit 'cpu_instructions', got '%s'", units.CPUUnit)
	}

	if units.MemoryUnit != "bytes" {
		t.Errorf("Expected MemoryUnit 'bytes', got '%s'", units.MemoryUnit)
	}

	if units.OperationsUnit != "operations" {
		t.Errorf("Expected OperationsUnit 'operations', got '%s'", units.OperationsUnit)
	}

	if units.CPURoundingRule != "none" {
		t.Errorf("Expected CPURoundingRule 'none', got '%s'", units.CPURoundingRule)
	}

	if units.MemoryRoundingRule != "none" {
		t.Errorf("Expected MemoryRoundingRule 'none', got '%s'", units.MemoryRoundingRule)
	}

	if units.PercentagePrecision != 1 {
		t.Errorf("Expected PercentagePrecision 1, got %d", units.PercentagePrecision)
	}
}

func TestCanonicalFixtures(t *testing.T) {
	fixtures := GetCanonicalFixtures()

	if len(fixtures) != 7 {
		t.Errorf("Expected 7 fixtures, got %d", len(fixtures))
	}

	// Test exact boundary CPU fixture
	cpuBoundary := findFixture(fixtures, "exact-boundary-cpu")
	if cpuBoundary == nil {
		t.Fatal("exact-boundary-cpu fixture not found")
	}

	if cpuBoundary.ExpectedCPUUsage != DefaultCPULimit {
		t.Errorf("Expected CPU usage %d, got %d", DefaultCPULimit, cpuBoundary.ExpectedCPUUsage)
	}

	if cpuBoundary.ShouldExceedCPU {
		t.Error("Expected ShouldExceedCPU to be false for exact boundary")
	}

	// Test one over CPU fixture
	cpuOver := findFixture(fixtures, "one-over-cpu")
	if cpuOver == nil {
		t.Fatal("one-over-cpu fixture not found")
	}

	if cpuOver.ExpectedCPUUsage != DefaultCPULimit+1 {
		t.Errorf("Expected CPU usage %d, got %d", DefaultCPULimit+1, cpuOver.ExpectedCPUUsage)
	}

	if !cpuOver.ShouldExceedCPU {
		t.Error("Expected ShouldExceedCPU to be true for one over limit")
	}

	// Test exact boundary memory fixture
	memBoundary := findFixture(fixtures, "exact-boundary-memory")
	if memBoundary == nil {
		t.Fatal("exact-boundary-memory fixture not found")
	}

	if memBoundary.ExpectedMemoryUsage != DefaultMemoryLimit {
		t.Errorf("Expected memory usage %d, got %d", DefaultMemoryLimit, memBoundary.ExpectedMemoryUsage)
	}

	if memBoundary.ShouldExceedMemory {
		t.Error("Expected ShouldExceedMemory to be false for exact boundary")
	}

	// Test one over memory fixture
	memOver := findFixture(fixtures, "one-over-memory")
	if memOver == nil {
		t.Fatal("one-over-memory fixture not found")
	}

	if memOver.ExpectedMemoryUsage != DefaultMemoryLimit+1 {
		t.Errorf("Expected memory usage %d, got %d", DefaultMemoryLimit+1, memOver.ExpectedMemoryUsage)
	}

	if !memOver.ShouldExceedMemory {
		t.Error("Expected ShouldExceedMemory to be true for one over limit")
	}

	// Test one over operations fixture
	opsOver := findFixture(fixtures, "one-over-operations")
	if opsOver == nil {
		t.Fatal("one-over-operations fixture not found")
	}

	if opsOver.ExpectedOperationsCount != MaxOperationsLimit+1 {
		t.Errorf("Expected operations count %d, got %d", MaxOperationsLimit+1, opsOver.ExpectedOperationsCount)
	}

	if !opsOver.ShouldExceedOperations {
		t.Error("Expected ShouldExceedOperations to be true for one over limit")
	}
}

func findFixture(fixtures []*ResourceLimitFixture, name string) *ResourceLimitFixture {
	for _, f := range fixtures {
		if f.Name == name {
			return f
		}
	}
	return nil
}

func TestValidateResourceAgreementWithinLimits(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		50_000_000,  // 50% CPU
		25_000_000,  // 50% Memory
		10,          // 10 ops
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	if diagnostic.CPUExceeded {
		t.Error("Expected CPUExceeded to be false within limits")
	}

	if diagnostic.MemoryExceeded {
		t.Error("Expected MemoryExceeded to be false within limits")
	}

	if diagnostic.OperationsExceeded {
		t.Error("Expected OperationsExceeded to be false within limits")
	}

	summary := diagnostic.ToSummary()
	if summary != "Resource agreement: all limits within bounds" {
		t.Errorf("Expected 'all limits within bounds' summary, got '%s'", summary)
	}
}

func TestValidateResourceAgreementCPUExceeded(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		DefaultCPULimit+1,  // One over CPU
		25_000_000,
		10,
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	if !diagnostic.CPUExceeded {
		t.Error("Expected CPUExceeded to be true")
	}

	if diagnostic.CPUExceededBy != 1 {
		t.Errorf("Expected CPUExceededBy 1, got %d", diagnostic.CPUExceededBy)
	}

	if diagnostic.FirstExceeded != "cpu" {
		t.Errorf("Expected FirstExceeded 'cpu', got '%s'", diagnostic.FirstExceeded)
	}

	summary := diagnostic.ToSummary()
	if !containsString(summary, "CPU exceeded by 1") {
		t.Errorf("Expected summary to contain 'CPU exceeded by 1', got '%s'", summary)
	}
}

func TestValidateResourceAgreementMemoryExceeded(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		50_000_000,
		DefaultMemoryLimit+1,  // One over Memory
		10,
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	if !diagnostic.MemoryExceeded {
		t.Error("Expected MemoryExceeded to be true")
	}

	if diagnostic.MemoryExceededBy != 1 {
		t.Errorf("Expected MemoryExceededBy 1, got %d", diagnostic.MemoryExceededBy)
	}

	if diagnostic.FirstExceeded != "memory" {
		t.Errorf("Expected FirstExceeded 'memory', got '%s'", diagnostic.FirstExceeded)
	}

	summary := diagnostic.ToSummary()
	if !containsString(summary, "Memory exceeded by 1") {
		t.Errorf("Expected summary to contain 'Memory exceeded by 1', got '%s'", summary)
	}
}

func TestValidateResourceAgreementOperationsExceeded(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		50_000_000,
		25_000_000,
		MaxOperationsLimit+1,  // One over operations
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	if !diagnostic.OperationsExceeded {
		t.Error("Expected OperationsExceeded to be true")
	}

	if diagnostic.OperationsExceededBy != 1 {
		t.Errorf("Expected OperationsExceededBy 1, got %d", diagnostic.OperationsExceededBy)
	}

	if diagnostic.FirstExceeded != "operations" {
		t.Errorf("Expected FirstExceeded 'operations', got '%s'", diagnostic.FirstExceeded)
	}

	summary := diagnostic.ToSummary()
	if !containsString(summary, "Operations exceeded by 1") {
		t.Errorf("Expected summary to contain 'Operations exceeded by 1', got '%s'", summary)
	}
}

func TestValidateResourceAgreementMultipleExceeded(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		DefaultCPULimit+1,  // CPU exceeded
		DefaultMemoryLimit+1,  // Memory exceeded
		MaxOperationsLimit+1,  // Operations exceeded
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	if !diagnostic.CPUExceeded {
		t.Error("Expected CPUExceeded to be true")
	}

	if !diagnostic.MemoryExceeded {
		t.Error("Expected MemoryExceeded to be true")
	}

	if !diagnostic.OperationsExceeded {
		t.Error("Expected OperationsExceeded to be true")
	}

	// CPU should be first exceeded due to priority
	if diagnostic.FirstExceeded != "cpu" {
		t.Errorf("Expected FirstExceeded 'cpu' (priority), got '%s'", diagnostic.FirstExceeded)
	}
}

func TestClassifyFailureLimitExhaustion(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		DefaultCPULimit+1,
		25_000_000,
		10,
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	classification, message := ClassifyFailure(diagnostic, "")

	if classification != FailureClassificationLimitExhaustion {
		t.Errorf("Expected LimitExhaustion classification, got '%s'", classification)
	}

	if !containsString(message, "CPU exceeded") {
		t.Errorf("Expected message to contain 'CPU exceeded', got '%s'", message)
	}
}

func TestClassifyFailureContractLogic(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		50_000_000,
		25_000_000,
		10,
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	classification, message := ClassifyFailure(diagnostic, "insufficient balance")

	if classification != FailureClassificationContractLogic {
		t.Errorf("Expected ContractLogic classification, got '%s'", classification)
	}

	if !containsString(message, "insufficient balance") {
		t.Errorf("Expected message to contain 'insufficient balance', got '%s'", message)
	}
}

func TestClassifyFailureUnknown(t *testing.T) {
	diagnostic := ValidateResourceAgreement(
		50_000_000,
		25_000_000,
		10,
		DefaultCPULimit,
		DefaultMemoryLimit,
		MaxOperationsLimit,
	)

	classification, message := ClassifyFailure(diagnostic, "")

	if classification != FailureClassificationUnknown {
		t.Errorf("Expected Unknown classification, got '%s'", classification)
	}

	if !containsString(message, "Unable to classify failure") {
		t.Errorf("Expected message to contain 'Unable to classify failure', got '%s'", message)
	}
}

func TestResourceLimitConstants(t *testing.T) {
	// Verify constants match canonical values
	if DefaultCPULimit != 100_000_000 {
		t.Errorf("Expected DefaultCPULimit 100_000_000, got %d", DefaultCPULimit)
	}

	if DefaultMemoryLimit != 50_000_000 {
		t.Errorf("Expected DefaultMemoryLimit 50_000_000, got %d", DefaultMemoryLimit)
	}

	if MaxOperationsLimit != 100 {
		t.Errorf("Expected MaxOperationsLimit 100, got %d", MaxOperationsLimit)
	}
}

func TestExactBoundaryFixtureAgreement(t *testing.T) {
	fixtures := GetCanonicalFixtures()

	for _, fixture := range fixtures {
		if fixture.Name == "exact-boundary-cpu" {
			// Simulate exact boundary CPU usage
			diagnostic := ValidateResourceAgreement(
				fixture.ExpectedCPUUsage,
				fixture.ExpectedMemoryUsage,
				fixture.ExpectedOperationsCount,
				fixture.CPULimit,
				fixture.MemoryLimit,
				fixture.OperationsLimit,
			)

			if diagnostic.CPUExceeded {
				t.Errorf("Exact boundary CPU should not be marked as exceeded, but got CPUExceeded=true")
			}

			if diagnostic.ToSummary() != "Resource agreement: all limits within bounds" {
				t.Errorf("Expected 'all limits within bounds' for exact boundary, got '%s'", diagnostic.ToSummary())
			}
		}

		if fixture.Name == "exact-boundary-memory" {
			// Simulate exact boundary memory usage
			diagnostic := ValidateResourceAgreement(
				fixture.ExpectedCPUUsage,
				fixture.ExpectedMemoryUsage,
				fixture.ExpectedOperationsCount,
				fixture.CPULimit,
				fixture.MemoryLimit,
				fixture.OperationsLimit,
			)

			if diagnostic.MemoryExceeded {
				t.Errorf("Exact boundary memory should not be marked as exceeded, but got MemoryExceeded=true")
			}

			if diagnostic.ToSummary() != "Resource agreement: all limits within bounds" {
				t.Errorf("Expected 'all limits within bounds' for exact boundary, got '%s'", diagnostic.ToSummary())
			}
		}
	}
}

func TestOneOverLimitFixtureAgreement(t *testing.T) {
	fixtures := GetCanonicalFixtures()

	for _, fixture := range fixtures {
		if fixture.Name == "one-over-cpu" {
			diagnostic := ValidateResourceAgreement(
				fixture.ExpectedCPUUsage,
				fixture.ExpectedMemoryUsage,
				fixture.ExpectedOperationsCount,
				fixture.CPULimit,
				fixture.MemoryLimit,
				fixture.OperationsLimit,
			)

			if !diagnostic.CPUExceeded {
				t.Error("One over CPU should be marked as exceeded")
			}

			if diagnostic.CPUExceededBy != 1 {
				t.Errorf("Expected CPUExceededBy 1, got %d", diagnostic.CPUExceededBy)
			}
		}

		if fixture.Name == "one-over-memory" {
			diagnostic := ValidateResourceAgreement(
				fixture.ExpectedCPUUsage,
				fixture.ExpectedMemoryUsage,
				fixture.ExpectedOperationsCount,
				fixture.CPULimit,
				fixture.MemoryLimit,
				fixture.OperationsLimit,
			)

			if !diagnostic.MemoryExceeded {
				t.Error("One over memory should be marked as exceeded")
			}

			if diagnostic.MemoryExceededBy != 1 {
				t.Errorf("Expected MemoryExceededBy 1, got %d", diagnostic.MemoryExceededBy)
			}
		}

		if fixture.Name == "one-over-operations" {
			diagnostic := ValidateResourceAgreement(
				fixture.ExpectedCPUUsage,
				fixture.ExpectedMemoryUsage,
				fixture.ExpectedOperationsCount,
				fixture.CPULimit,
				fixture.MemoryLimit,
				fixture.OperationsLimit,
			)

			if !diagnostic.OperationsExceeded {
				t.Error("One over operations should be marked as exceeded")
			}

			if diagnostic.OperationsExceededBy != 1 {
				t.Errorf("Expected OperationsExceededBy 1, got %d", diagnostic.OperationsExceededBy)
			}
		}
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
