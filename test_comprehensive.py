#!/usr/bin/env python3
"""
Comprehensive Test Orchestration Script for etcd-monitor
========================================================

This script provides a unified interface to run all tests, generate coverage
reports, execute benchmarks, and aggregate results for the etcd-monitor project.

Usage:
    python3 test_comprehensive.py --all          # Run all tests
    python3 test_comprehensive.py --unit         # Unit tests only
    python3 test_comprehensive.py --integration  # Integration tests only
    python3 test_comprehensive.py --benchmark    # Benchmarks only
    python3 test_comprehensive.py --coverage     # Generate coverage report
    python3 test_comprehensive.py --race         # Run with race detection
    python3 test_comprehensive.py --verbose      # Verbose output
    python3 test_comprehensive.py --package <pkg>  # Test specific package

Author: etcd-monitor Team
Date: 2025-10-18
"""

import argparse
import subprocess
import sys
import os
import json
import time
from datetime import datetime
from pathlib import Path
from typing import List, Dict, Tuple, Optional

# ANSI color codes for terminal output
class Colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    ENDC = '\033[0m'
    BOLD = '\033[1m'
    UNDERLINE = '\033[4m'

class TestRunner:
    """Main test orchestration class"""

    def __init__(self, args):
        self.args = args
        self.project_root = Path(__file__).parent.resolve()
        self.coverage_dir = self.project_root / "coverage"
        self.results = {
            "unit_tests": {},
            "integration_tests": {},
            "benchmarks": {},
            "coverage": {},
            "start_time": datetime.now().isoformat(),
            "success": False
        }

        # Ensure coverage directory exists
        self.coverage_dir.mkdir(exist_ok=True)

    def run(self) -> int:
        """Main entry point - orchestrate all test execution"""
        print(f"{Colors.HEADER}{Colors.BOLD}")
        print("=" * 80)
        print("  etcd-monitor Comprehensive Test Suite")
        print("=" * 80)
        print(f"{Colors.ENDC}")
        print(f"Project Root: {self.project_root}")
        print(f"Coverage Dir: {self.coverage_dir}")
        print()

        start_time = time.time()

        try:
            # Validate environment
            if not self._validate_environment():
                return 1

            # Run tests based on arguments
            if self.args.all or self.args.unit:
                if not self._run_unit_tests():
                    return 1

            if self.args.all or self.args.integration:
                if not self._run_integration_tests():
                    return 1

            if self.args.all or self.args.benchmark:
                self._run_benchmarks()

            if self.args.coverage or self.args.all:
                self._generate_coverage_report()

            # Generate final report
            self._generate_summary_report()

            # Mark as successful
            self.results["success"] = True

            # Save results to JSON
            self._save_results()

            elapsed = time.time() - start_time
            print(f"\n{Colors.GREEN}{Colors.BOLD}✓ All tests completed successfully in {elapsed:.2f}s{Colors.ENDC}\n")
            return 0

        except KeyboardInterrupt:
            print(f"\n{Colors.YELLOW}Test execution interrupted by user{Colors.ENDC}")
            return 130
        except Exception as e:
            print(f"\n{Colors.RED}{Colors.BOLD}✗ Test execution failed: {e}{Colors.ENDC}")
            import traceback
            traceback.print_exc()
            return 1

    def _validate_environment(self) -> bool:
        """Validate that required tools are available"""
        print(f"{Colors.CYAN}Validating environment...{Colors.ENDC}")

        # Check for Go
        if not self._check_command("go version"):
            print(f"{Colors.RED}✗ Go not found. Please install Go 1.19+{Colors.ENDC}")
            return False

        # Check Go version
        result = subprocess.run(["go", "version"], capture_output=True, text=True)
        print(f"  ✓ {result.stdout.strip()}")

        # Check for git (optional but recommended)
        if self._check_command("git --version"):
            result = subprocess.run(["git", "--version"], capture_output=True, text=True)
            print(f"  ✓ {result.stdout.strip()}")

        # Check for golangci-lint (optional)
        if self._check_command("golangci-lint --version"):
            result = subprocess.run(["golangci-lint", "--version"], capture_output=True, text=True)
            print(f"  ✓ golangci-lint found")
        else:
            print(f"  ⚠ golangci-lint not found (optional)")

        print()
        return True

    def _check_command(self, cmd: str) -> bool:
        """Check if a command exists"""
        try:
            subprocess.run(cmd.split(), capture_output=True, check=True)
            return True
        except (subprocess.CalledProcessError, FileNotFoundError):
            return False

    def _run_unit_tests(self) -> bool:
        """Run all unit tests"""
        print(f"{Colors.CYAN}{Colors.BOLD}Running Unit Tests...{Colors.ENDC}")
        print()

        # Build test command
        cmd = ["go", "test"]

        # Add package selector
        if self.args.package:
            cmd.append(f"./{self.args.package}/...")
        else:
            cmd.append("./...")

        # Add flags
        cmd.extend([
            "-v" if self.args.verbose else "-v",
            "-coverprofile=" + str(self.coverage_dir / "coverage.out"),
            "-covermode=atomic",
        ])

        # Add race detection if requested
        if self.args.race or self.args.all:
            cmd.append("-race")
            print(f"  {Colors.YELLOW}Race detection enabled{Colors.ENDC}")

        # Add timeout
        cmd.append("-timeout=10m")

        # Add short flag if quick mode
        if self.args.quick:
            cmd.append("-short")

        # Run tests
        print(f"  Command: {' '.join(cmd)}")
        print()

        start_time = time.time()
        result = subprocess.run(cmd, cwd=self.project_root)
        elapsed = time.time() - start_time

        # Store results
        self.results["unit_tests"] = {
            "command": " ".join(cmd),
            "exit_code": result.returncode,
            "duration": elapsed,
            "success": result.returncode == 0
        }

        if result.returncode == 0:
            print(f"\n{Colors.GREEN}✓ Unit tests passed in {elapsed:.2f}s{Colors.ENDC}\n")
            return True
        else:
            print(f"\n{Colors.RED}✗ Unit tests failed with exit code {result.returncode}{Colors.ENDC}\n")
            return False

    def _run_integration_tests(self) -> bool:
        """Run integration tests"""
        print(f"{Colors.CYAN}{Colors.BOLD}Running Integration Tests...{Colors.ENDC}")
        print()

        # Integration tests are in integration_tests/ directory
        integration_dir = self.project_root / "integration_tests"

        if not integration_dir.exists():
            print(f"  {Colors.YELLOW}No integration tests directory found, skipping...{Colors.ENDC}\n")
            self.results["integration_tests"]["skipped"] = True
            return True

        # Build test command
        cmd = [
            "go", "test",
            "-v",
            "-tags=integration",
            "-timeout=30m",
            "./integration_tests/..."
        ]

        if self.args.race:
            cmd.append("-race")

        print(f"  Command: {' '.join(cmd)}")
        print()

        start_time = time.time()
        result = subprocess.run(cmd, cwd=self.project_root)
        elapsed = time.time() - start_time

        # Store results
        self.results["integration_tests"] = {
            "command": " ".join(cmd),
            "exit_code": result.returncode,
            "duration": elapsed,
            "success": result.returncode == 0
        }

        if result.returncode == 0:
            print(f"\n{Colors.GREEN}✓ Integration tests passed in {elapsed:.2f}s{Colors.ENDC}\n")
            return True
        else:
            print(f"\n{Colors.RED}✗ Integration tests failed with exit code {result.returncode}{Colors.ENDC}\n")
            return False

    def _run_benchmarks(self) -> bool:
        """Run Go benchmarks"""
        print(f"{Colors.CYAN}{Colors.BOLD}Running Benchmarks...{Colors.ENDC}")
        print()

        # Build benchmark command
        cmd = [
            "go", "test",
            "-bench=.",
            "-benchmem",
            "-benchtime=5s",
            "-run=^$",  # Don't run regular tests
        ]

        if self.args.package:
            cmd.append(f"./{self.args.package}/...")
        else:
            cmd.append("./...")

        print(f"  Command: {' '.join(cmd)}")
        print()

        start_time = time.time()
        result = subprocess.run(cmd, cwd=self.project_root, capture_output=True, text=True)
        elapsed = time.time() - start_time

        # Store results
        self.results["benchmarks"] = {
            "command": " ".join(cmd),
            "exit_code": result.returncode,
            "duration": elapsed,
            "output": result.stdout,
            "success": result.returncode == 0
        }

        # Print benchmark output
        print(result.stdout)

        if result.returncode == 0:
            print(f"\n{Colors.GREEN}✓ Benchmarks completed in {elapsed:.2f}s{Colors.ENDC}\n")

            # Save benchmark output to file
            benchmark_file = self.coverage_dir / "benchmarks.txt"
            benchmark_file.write_text(result.stdout)
            print(f"  Benchmark results saved to: {benchmark_file}")
            print()

            return True
        else:
            print(f"\n{Colors.YELLOW}⚠ Benchmarks completed with warnings{Colors.ENDC}\n")
            return True  # Don't fail on benchmark issues

    def _generate_coverage_report(self):
        """Generate HTML coverage report"""
        print(f"{Colors.CYAN}{Colors.BOLD}Generating Coverage Report...{Colors.ENDC}")
        print()

        coverage_file = self.coverage_dir / "coverage.out"

        if not coverage_file.exists():
            print(f"  {Colors.YELLOW}No coverage data found, skipping...{Colors.ENDC}\n")
            return

        # Generate HTML report
        html_file = self.coverage_dir / "coverage.html"
        cmd = ["go", "tool", "cover", f"-html={coverage_file}", f"-o={html_file}"]

        result = subprocess.run(cmd, cwd=self.project_root, capture_output=True)

        if result.returncode == 0:
            print(f"  {Colors.GREEN}✓ HTML coverage report generated: {html_file}{Colors.ENDC}")

        # Calculate coverage percentage
        cmd = ["go", "tool", "cover", f"-func={coverage_file}"]
        result = subprocess.run(cmd, cwd=self.project_root, capture_output=True, text=True)

        if result.returncode == 0:
            # Parse coverage output
            lines = result.stdout.strip().split('\n')
            total_line = lines[-1]  # Last line has total coverage

            # Extract percentage from "total:  (statements)  XX.X%"
            if "%" in total_line:
                coverage_pct = total_line.split()[-1].rstrip('%')
                try:
                    coverage_value = float(coverage_pct)
                    self.results["coverage"] = {
                        "percentage": coverage_value,
                        "report": str(html_file),
                        "raw_output": result.stdout
                    }

                    # Color code the coverage percentage
                    if coverage_value >= 95:
                        color = Colors.GREEN
                        status = "EXCELLENT"
                    elif coverage_value >= 80:
                        color = Colors.CYAN
                        status = "GOOD"
                    elif coverage_value >= 60:
                        color = Colors.YELLOW
                        status = "NEEDS IMPROVEMENT"
                    else:
                        color = Colors.RED
                        status = "CRITICAL"

                    print(f"\n  {color}{Colors.BOLD}Overall Coverage: {coverage_pct}% ({status}){Colors.ENDC}")

                    # Save coverage summary
                    summary_file = self.coverage_dir / "coverage-summary.txt"
                    summary_file.write_text(result.stdout)
                    print(f"  Coverage summary saved to: {summary_file}")

                except ValueError:
                    print(f"  {Colors.YELLOW}Could not parse coverage percentage{Colors.ENDC}")

        print()

    def _generate_summary_report(self):
        """Generate and display summary report"""
        print(f"{Colors.HEADER}{Colors.BOLD}")
        print("=" * 80)
        print("  Test Execution Summary")
        print("=" * 80)
        print(f"{Colors.ENDC}")

        # Unit tests
        if "unit_tests" in self.results and not self.results["unit_tests"].get("skipped"):
            unit = self.results["unit_tests"]
            status = "✓ PASSED" if unit.get("success") else "✗ FAILED"
            color = Colors.GREEN if unit.get("success") else Colors.RED
            print(f"\n{Colors.BOLD}Unit Tests:{Colors.ENDC}")
            print(f"  Status: {color}{status}{Colors.ENDC}")
            if "duration" in unit:
                print(f"  Duration: {unit['duration']:.2f}s")

        # Integration tests
        if "integration_tests" in self.results and not self.results["integration_tests"].get("skipped"):
            integration = self.results["integration_tests"]
            status = "✓ PASSED" if integration.get("success") else "✗ FAILED"
            color = Colors.GREEN if integration.get("success") else Colors.RED
            print(f"\n{Colors.BOLD}Integration Tests:{Colors.ENDC}")
            print(f"  Status: {color}{status}{Colors.ENDC}")
            if "duration" in integration:
                print(f"  Duration: {integration['duration']:.2f}s")

        # Benchmarks
        if "benchmarks" in self.results and not self.results["benchmarks"].get("skipped"):
            bench = self.results["benchmarks"]
            status = "✓ COMPLETED" if bench.get("success") else "⚠ WARNINGS"
            color = Colors.GREEN if bench.get("success") else Colors.YELLOW
            print(f"\n{Colors.BOLD}Benchmarks:{Colors.ENDC}")
            print(f"  Status: {color}{status}{Colors.ENDC}")
            if "duration" in bench:
                print(f"  Duration: {bench['duration']:.2f}s")

        # Coverage
        if "coverage" in self.results and "percentage" in self.results["coverage"]:
            cov = self.results["coverage"]
            pct = cov["percentage"]

            if pct >= 95:
                color = Colors.GREEN
                status = "EXCELLENT"
            elif pct >= 80:
                color = Colors.CYAN
                status = "GOOD"
            elif pct >= 60:
                color = Colors.YELLOW
                status = "NEEDS IMPROVEMENT"
            else:
                color = Colors.RED
                status = "CRITICAL"

            print(f"\n{Colors.BOLD}Test Coverage:{Colors.ENDC}")
            print(f"  Overall: {color}{pct:.1f}% ({status}){Colors.ENDC}")
            print(f"  Report: {cov['report']}")

        print()
        print("=" * 80)
        print()

    def _save_results(self):
        """Save test results to JSON file"""
        self.results["end_time"] = datetime.now().isoformat()

        results_file = self.coverage_dir / "test-results.json"
        with open(results_file, 'w') as f:
            json.dump(self.results, f, indent=2)

        print(f"Test results saved to: {results_file}\n")


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description="Comprehensive test orchestration for etcd-monitor",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s --all              Run all tests with coverage
  %(prog)s --unit --race      Run unit tests with race detection
  %(prog)s --integration      Run integration tests only
  %(prog)s --benchmark        Run benchmarks only
  %(prog)s --coverage         Generate coverage report from existing data
  %(prog)s --package monitor  Test only the monitor package
        """
    )

    # Test selection
    parser.add_argument("--all", action="store_true",
                       help="Run all tests (unit, integration, benchmarks)")
    parser.add_argument("--unit", action="store_true",
                       help="Run unit tests only")
    parser.add_argument("--integration", action="store_true",
                       help="Run integration tests only")
    parser.add_argument("--benchmark", action="store_true",
                       help="Run benchmarks only")
    parser.add_argument("--coverage", action="store_true",
                       help="Generate coverage report")

    # Test options
    parser.add_argument("--race", action="store_true",
                       help="Enable race detection")
    parser.add_argument("--verbose", "-v", action="store_true",
                       help="Verbose output")
    parser.add_argument("--quick", action="store_true",
                       help="Run tests in short mode")
    parser.add_argument("--package", type=str,
                       help="Test specific package (e.g., 'pkg/monitor')")

    args = parser.parse_args()

    # If no test type specified, default to --all
    if not any([args.all, args.unit, args.integration, args.benchmark, args.coverage]):
        args.all = True

    # Create and run test runner
    runner = TestRunner(args)
    sys.exit(runner.run())


if __name__ == "__main__":
    main()
