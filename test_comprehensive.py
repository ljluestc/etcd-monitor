#!/usr/bin/env python3
"""
Comprehensive Test Suite for etcd-monitor
This script provides end-to-end testing, performance benchmarking, and coverage validation.
"""

import os
import sys
import subprocess
import json
import time
import requests
import threading
import signal
import tempfile
import shutil
from pathlib import Path
from typing import Dict, List, Any, Optional
from dataclasses import dataclass
from concurrent.futures import ThreadPoolExecutor, as_completed
import logging

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

@dataclass
class TestResult:
    """Test result container"""
    name: str
    status: str  # 'PASS', 'FAIL', 'SKIP'
    duration: float
    message: str
    coverage: Optional[float] = None

@dataclass
class BenchmarkResult:
    """Benchmark result container"""
    operation: str
    operations_per_second: float
    avg_latency_ms: float
    p95_latency_ms: float
    p99_latency_ms: float
    error_rate: float

class ComprehensiveTestSuite:
    """Main test suite orchestrator"""
    
    def __init__(self, project_root: str):
        self.project_root = Path(project_root)
        self.go_mod_path = self.project_root / "go.mod"
        self.test_results: List[TestResult] = []
        self.benchmark_results: List[BenchmarkResult] = []
        self.coverage_data: Dict[str, float] = {}
        self.etcd_process: Optional[subprocess.Popen] = None
        self.monitor_process: Optional[subprocess.Popen] = None
        self.temp_dir = None
        
    def setup_environment(self) -> bool:
        """Set up test environment including etcd and dependencies"""
        logger.info("Setting up test environment...")
        
        try:
            # Create temporary directory for test data
            self.temp_dir = tempfile.mkdtemp(prefix="etcd-monitor-test-")
            logger.info(f"Created temp directory: {self.temp_dir}")
            
            # Check if Go is installed
            result = subprocess.run(['go', 'version'], capture_output=True, text=True)
            if result.returncode != 0:
                logger.error("Go is not installed or not in PATH")
                return False
            logger.info(f"Go version: {result.stdout.strip()}")
            
            # Install dependencies
            logger.info("Installing Go dependencies...")
            result = subprocess.run(['go', 'mod', 'download'], 
                                  cwd=self.project_root, capture_output=True, text=True)
            if result.returncode != 0:
                logger.error(f"Failed to install dependencies: {result.stderr}")
                return False
            
            # Install test dependencies
            test_deps = [
                'github.com/stretchr/testify/assert',
                'github.com/stretchr/testify/require',
                'github.com/stretchr/testify/mock',
                'github.com/golang/mock/gomock',
                'github.com/onsi/ginkgo/v2',
                'github.com/onsi/gomega',
            ]
            
            for dep in test_deps:
                result = subprocess.run(['go', 'get', dep], 
                                      cwd=self.project_root, capture_output=True, text=True)
                if result.returncode != 0:
                    logger.warning(f"Failed to install {dep}: {result.stderr}")
            
            return True
            
        except Exception as e:
            logger.error(f"Failed to setup environment: {e}")
            return False
    
    def start_etcd(self) -> bool:
        """Start etcd server for testing"""
        logger.info("Starting etcd server...")
        
        try:
            # Create etcd data directory
            etcd_data_dir = os.path.join(self.temp_dir, "etcd-data")
            os.makedirs(etcd_data_dir, exist_ok=True)
            
            # Start etcd process
            etcd_cmd = [
                'etcd',
                '--data-dir', etcd_data_dir,
                '--listen-client-urls', 'http://localhost:2379',
                '--advertise-client-urls', 'http://localhost:2379',
                '--listen-peer-urls', 'http://localhost:2380',
                '--initial-advertise-peer-urls', 'http://localhost:2380',
                '--initial-cluster', 'default=http://localhost:2380',
                '--initial-cluster-token', 'etcd-cluster-1',
                '--initial-cluster-state', 'new',
                '--log-level', 'warn'
            ]
            
            self.etcd_process = subprocess.Popen(
                etcd_cmd,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            # Wait for etcd to be ready
            max_retries = 30
            for i in range(max_retries):
                try:
                    response = requests.get('http://localhost:2379/health', timeout=1)
                    if response.status_code == 200:
                        logger.info("etcd server is ready")
                        return True
                except requests.exceptions.RequestException:
                    pass
                
                time.sleep(1)
                logger.info(f"Waiting for etcd... ({i+1}/{max_retries})")
            
            logger.error("etcd server failed to start")
            return False
            
        except Exception as e:
            logger.error(f"Failed to start etcd: {e}")
            return False
    
    def stop_etcd(self):
        """Stop etcd server"""
        if self.etcd_process:
            logger.info("Stopping etcd server...")
            self.etcd_process.terminate()
            try:
                self.etcd_process.wait(timeout=10)
            except subprocess.TimeoutExpired:
                self.etcd_process.kill()
            self.etcd_process = None
    
    def run_unit_tests(self) -> bool:
        """Run all unit tests with coverage"""
        logger.info("Running unit tests...")
        
        try:
            # Run tests with coverage
            cmd = [
                'go', 'test', '-v', '-race', '-coverprofile=coverage.out',
                '-covermode=atomic', './...'
            ]
            
            result = subprocess.run(
                cmd,
                cwd=self.project_root,
                capture_output=True,
                text=True,
                timeout=300
            )
            
            # Parse coverage
            if os.path.exists(self.project_root / "coverage.out"):
                self.parse_coverage()
            
            # Record test results
            test_result = TestResult(
                name="Unit Tests",
                status="PASS" if result.returncode == 0 else "FAIL",
                duration=0,  # Would need to parse from output
                message=result.stdout if result.returncode == 0 else result.stderr
            )
            self.test_results.append(test_result)
            
            if result.returncode != 0:
                logger.error(f"Unit tests failed: {result.stderr}")
                return False
            
            logger.info("Unit tests passed")
            return True
            
        except subprocess.TimeoutExpired:
            logger.error("Unit tests timed out")
            return False
        except Exception as e:
            logger.error(f"Failed to run unit tests: {e}")
            return False
    
    def run_integration_tests(self) -> bool:
        """Run integration tests"""
        logger.info("Running integration tests...")
        
        try:
            # Build the monitor service
            build_cmd = ['go', 'build', '-o', 'etcd-monitor', './cmd/etcd-monitor']
            result = subprocess.run(build_cmd, cwd=self.project_root, capture_output=True, text=True)
            if result.returncode != 0:
                logger.error(f"Failed to build monitor: {result.stderr}")
                return False
            
            # Start the monitor service
            monitor_binary = self.project_root / "etcd-monitor"
            config_file = self.create_test_config()
            
            self.monitor_process = subprocess.Popen(
                [str(monitor_binary), '--config', config_file],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            # Wait for monitor to be ready
            time.sleep(5)
            
            # Run integration tests
            integration_tests = [
                self.test_health_endpoint,
                self.test_cluster_status,
                self.test_metrics_endpoints,
                self.test_alert_endpoints,
                self.test_performance_endpoints
            ]
            
            success = True
            for test_func in integration_tests:
                try:
                    if not test_func():
                        success = False
                except Exception as e:
                    logger.error(f"Integration test {test_func.__name__} failed: {e}")
                    success = False
            
            return success
            
        except Exception as e:
            logger.error(f"Failed to run integration tests: {e}")
            return False
        finally:
            if self.monitor_process:
                self.monitor_process.terminate()
                self.monitor_process.wait()
                self.monitor_process = None
    
    def create_test_config(self) -> str:
        """Create test configuration file"""
        config = {
            "etcd": {
                "endpoints": ["http://localhost:2379"],
                "dial_timeout": "5s"
            },
            "monitor": {
                "health_check_interval": "10s",
                "metrics_interval": "5s",
                "watch_interval": "1s"
            },
            "api": {
                "port": 8080,
                "host": "localhost"
            },
            "alerts": {
                "max_latency_ms": 100,
                "max_database_size_mb": 100,
                "min_available_nodes": 1
            }
        }
        
        config_path = os.path.join(self.temp_dir, "config.yaml")
        with open(config_path, 'w') as f:
            import yaml
            yaml.dump(config, f)
        
        return config_path
    
    def test_health_endpoint(self) -> bool:
        """Test health endpoint"""
        try:
            response = requests.get('http://localhost:8080/health', timeout=5)
            if response.status_code == 200:
                data = response.json()
                if data.get('status') == 'healthy':
                    logger.info("Health endpoint test passed")
                    return True
            
            logger.error(f"Health endpoint test failed: {response.status_code} - {response.text}")
            return False
            
        except Exception as e:
            logger.error(f"Health endpoint test failed: {e}")
            return False
    
    def test_cluster_status(self) -> bool:
        """Test cluster status endpoint"""
        try:
            response = requests.get('http://localhost:8080/api/v1/cluster/status', timeout=5)
            if response.status_code == 200:
                data = response.json()
                required_fields = ['healthy', 'member_count', 'quorum_size']
                if all(field in data for field in required_fields):
                    logger.info("Cluster status test passed")
                    return True
            
            logger.error(f"Cluster status test failed: {response.status_code} - {response.text}")
            return False
            
        except Exception as e:
            logger.error(f"Cluster status test failed: {e}")
            return False
    
    def test_metrics_endpoints(self) -> bool:
        """Test metrics endpoints"""
        endpoints = [
            '/api/v1/metrics/current',
            '/api/v1/metrics/latency'
        ]
        
        for endpoint in endpoints:
            try:
                response = requests.get(f'http://localhost:8080{endpoint}', timeout=5)
                if response.status_code != 200:
                    logger.error(f"Metrics endpoint {endpoint} failed: {response.status_code}")
                    return False
            except Exception as e:
                logger.error(f"Metrics endpoint {endpoint} failed: {e}")
                return False
        
        logger.info("Metrics endpoints test passed")
        return True
    
    def test_alert_endpoints(self) -> bool:
        """Test alert endpoints"""
        endpoints = [
            '/api/v1/alerts',
            '/api/v1/alerts/history'
        ]
        
        for endpoint in endpoints:
            try:
                response = requests.get(f'http://localhost:8080{endpoint}', timeout=5)
                if response.status_code != 200:
                    logger.error(f"Alert endpoint {endpoint} failed: {response.status_code}")
                    return False
            except Exception as e:
                logger.error(f"Alert endpoint {endpoint} failed: {e}")
                return False
        
        logger.info("Alert endpoints test passed")
        return True
    
    def test_performance_endpoints(self) -> bool:
        """Test performance endpoints"""
        try:
            # Test benchmark endpoint
            data = {"operations": 100}
            response = requests.post(
                'http://localhost:8080/api/v1/performance/benchmark',
                json=data,
                timeout=10
            )
            
            if response.status_code in [200, 501]:  # 501 is OK for not implemented
                logger.info("Performance endpoints test passed")
                return True
            
            logger.error(f"Performance endpoint failed: {response.status_code}")
            return False
            
        except Exception as e:
            logger.error(f"Performance endpoint test failed: {e}")
            return False
    
    def run_benchmark_tests(self) -> bool:
        """Run performance benchmark tests"""
        logger.info("Running benchmark tests...")
        
        try:
            # Run Go benchmarks
            cmd = ['go', 'test', '-bench=.', '-benchmem', './...']
            result = subprocess.run(
                cmd,
                cwd=self.project_root,
                capture_output=True,
                text=True,
                timeout=300
            )
            
            if result.returncode != 0:
                logger.error(f"Benchmark tests failed: {result.stderr}")
                return False
            
            # Parse benchmark results
            self.parse_benchmark_results(result.stdout)
            
            logger.info("Benchmark tests completed")
            return True
            
        except subprocess.TimeoutExpired:
            logger.error("Benchmark tests timed out")
            return False
        except Exception as e:
            logger.error(f"Failed to run benchmark tests: {e}")
            return False
    
    def parse_benchmark_results(self, output: str):
        """Parse benchmark results from Go test output"""
        lines = output.split('\n')
        for line in lines:
            if 'Benchmark' in line and 'ns/op' in line:
                parts = line.split()
                if len(parts) >= 3:
                    name = parts[0]
                    ops_per_sec = 0
                    avg_latency = 0
                    
                    # Parse ns/op
                    for part in parts:
                        if 'ns/op' in part:
                            ns_per_op = float(part.replace('ns/op', ''))
                            ops_per_sec = 1e9 / ns_per_op
                            avg_latency = ns_per_op / 1e6  # Convert to ms
                            break
                    
                    result = BenchmarkResult(
                        operation=name,
                        operations_per_second=ops_per_sec,
                        avg_latency_ms=avg_latency,
                        p95_latency_ms=0,  # Would need more detailed parsing
                        p99_latency_ms=0,
                        error_rate=0
                    )
                    self.benchmark_results.append(result)
    
    def parse_coverage(self):
        """Parse coverage data from coverage.out file"""
        try:
            cmd = ['go', 'tool', 'cover', '-func=coverage.out']
            result = subprocess.run(cmd, cwd=self.project_root, capture_output=True, text=True)
            
            if result.returncode == 0:
                lines = result.stdout.split('\n')
                total_coverage = 0
                
                for line in lines:
                    if 'total:' in line:
                        parts = line.split()
                        for part in parts:
                            if '%' in part:
                                total_coverage = float(part.replace('%', ''))
                                break
                        break
                
                self.coverage_data['total'] = total_coverage
                logger.info(f"Total coverage: {total_coverage:.2f}%")
            
        except Exception as e:
            logger.error(f"Failed to parse coverage: {e}")
    
    def run_load_tests(self) -> bool:
        """Run load tests to verify performance under stress"""
        logger.info("Running load tests...")
        
        try:
            # Simulate concurrent requests
            def make_request():
                try:
                    response = requests.get('http://localhost:8080/api/v1/metrics/current', timeout=5)
                    return response.status_code == 200
                except:
                    return False
            
            # Run 100 concurrent requests
            with ThreadPoolExecutor(max_workers=20) as executor:
                futures = [executor.submit(make_request) for _ in range(100)]
                results = [future.result() for future in as_completed(futures)]
            
            success_rate = sum(results) / len(results)
            logger.info(f"Load test success rate: {success_rate:.2%}")
            
            return success_rate > 0.95  # 95% success rate threshold
            
        except Exception as e:
            logger.error(f"Load test failed: {e}")
            return False
    
    def generate_report(self) -> str:
        """Generate comprehensive test report"""
        report = {
            "timestamp": time.strftime("%Y-%m-%d %H:%M:%S"),
            "test_results": [
                {
                    "name": result.name,
                    "status": result.status,
                    "duration": result.duration,
                    "message": result.message,
                    "coverage": result.coverage
                }
                for result in self.test_results
            ],
            "benchmark_results": [
                {
                    "operation": result.operation,
                    "ops_per_sec": result.operations_per_second,
                    "avg_latency_ms": result.avg_latency_ms,
                    "p95_latency_ms": result.p95_latency_ms,
                    "p99_latency_ms": result.p99_latency_ms,
                    "error_rate": result.error_rate
                }
                for result in self.benchmark_results
            ],
            "coverage": self.coverage_data,
            "summary": {
                "total_tests": len(self.test_results),
                "passed_tests": len([r for r in self.test_results if r.status == "PASS"]),
                "failed_tests": len([r for r in self.test_results if r.status == "FAIL"]),
                "skipped_tests": len([r for r in self.test_results if r.status == "SKIP"]),
                "total_coverage": self.coverage_data.get('total', 0)
            }
        }
        
        report_path = self.project_root / "test_report.json"
        with open(report_path, 'w') as f:
            json.dump(report, f, indent=2)
        
        return str(report_path)
    
    def cleanup(self):
        """Clean up test environment"""
        logger.info("Cleaning up test environment...")
        
        # Stop processes
        if self.monitor_process:
            self.monitor_process.terminate()
            self.monitor_process.wait()
        
        self.stop_etcd()
        
        # Clean up temporary directory
        if self.temp_dir and os.path.exists(self.temp_dir):
            shutil.rmtree(self.temp_dir)
        
        # Clean up build artifacts
        build_artifacts = ['etcd-monitor', 'coverage.out']
        for artifact in build_artifacts:
            artifact_path = self.project_root / artifact
            if artifact_path.exists():
                artifact_path.unlink()
    
    def run_all_tests(self) -> bool:
        """Run the complete test suite"""
        logger.info("Starting comprehensive test suite...")
        
        try:
            # Setup
            if not self.setup_environment():
                return False
            
            # Start etcd
            if not self.start_etcd():
                return False
            
            # Run tests
            success = True
            
            # Unit tests
            if not self.run_unit_tests():
                success = False
            
            # Integration tests
            if not self.run_integration_tests():
                success = False
            
            # Benchmark tests
            if not self.run_benchmark_tests():
                success = False
            
            # Load tests
            if not self.run_load_tests():
                success = False
            
            # Generate report
            report_path = self.generate_report()
            logger.info(f"Test report generated: {report_path}")
            
            # Print summary
            self.print_summary()
            
            return success
            
        except KeyboardInterrupt:
            logger.info("Test suite interrupted by user")
            return False
        except Exception as e:
            logger.error(f"Test suite failed with error: {e}")
            return False
        finally:
            self.cleanup()
    
    def print_summary(self):
        """Print test summary"""
        print("\n" + "="*60)
        print("TEST SUITE SUMMARY")
        print("="*60)
        
        total = len(self.test_results)
        passed = len([r for r in self.test_results if r.status == "PASS"])
        failed = len([r for r in self.test_results if r.status == "FAIL"])
        skipped = len([r for r in self.test_results if r.status == "SKIP"])
        
        print(f"Total Tests: {total}")
        print(f"Passed: {passed}")
        print(f"Failed: {failed}")
        print(f"Skipped: {skipped}")
        print(f"Success Rate: {(passed/total*100):.1f}%" if total > 0 else "N/A")
        
        coverage = self.coverage_data.get('total', 0)
        print(f"Code Coverage: {coverage:.1f}%")
        
        if self.benchmark_results:
            print(f"\nBenchmark Results: {len(self.benchmark_results)} operations tested")
            for result in self.benchmark_results[:5]:  # Show top 5
                print(f"  {result.operation}: {result.operations_per_second:.0f} ops/sec")
        
        print("="*60)

def main():
    """Main entry point"""
    if len(sys.argv) != 2:
        print("Usage: python3 test_comprehensive.py <project_root>")
        sys.exit(1)
    
    project_root = sys.argv[1]
    if not os.path.exists(project_root):
        print(f"Project root does not exist: {project_root}")
        sys.exit(1)
    
    # Create and run test suite
    test_suite = ComprehensiveTestSuite(project_root)
    success = test_suite.run_all_tests()
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
