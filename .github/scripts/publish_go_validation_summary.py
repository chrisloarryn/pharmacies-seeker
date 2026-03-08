#!/usr/bin/env python3

from __future__ import annotations

import argparse
import json
import re
from pathlib import Path


def append_summary(summary_file: str | None, content: str) -> None:
    if not summary_file:
        return
    with Path(summary_file).open("a", encoding="utf-8") as handle:
        handle.write(content.rstrip() + "\n")


def write_output(output_file: str | None, key: str, value: str) -> None:
    if not output_file:
        return
    with Path(output_file).open("a", encoding="utf-8") as handle:
        handle.write(f"{key}={value}\n")


def read_lines(path: str | None) -> list[str]:
    if not path:
        return []
    file_path = Path(path)
    if not file_path.exists():
        return []
    return file_path.read_text(encoding="utf-8", errors="replace").splitlines()


def format_duration(seconds: float) -> str:
    return f"{seconds:.2f}s"


def format_optional_number(value: float | None, suffix: str = "") -> str:
    if value is None:
        return "n/a"
    if float(value).is_integer():
        return f"{int(value)}{suffix}"
    return f"{value:.2f}{suffix}"


def parse_bool(value: str | None) -> bool:
    return str(value).strip().lower() in {"1", "true", "yes", "y", "on"}


def parse_go_test_report(report_file: str | None) -> tuple[dict[str, float | int], list[tuple[str, str, float]], list[str]]:
    metrics: dict[str, float | int] = {
        "tests_run": 0,
        "failures": 0,
        "skipped": 0,
        "packages_total": 0,
        "packages_failed": 0,
        "duration_seconds": 0.0,
    }
    finalized_tests: set[tuple[str, str]] = set()
    package_results: dict[str, tuple[str, float]] = {}
    failed_tests: list[str] = []

    for raw_line in read_lines(report_file):
        try:
            event = json.loads(raw_line)
        except json.JSONDecodeError:
            continue

        package_name = event.get("Package")
        test_name = event.get("Test")
        action = event.get("Action")
        elapsed = float(event.get("Elapsed", 0.0) or 0.0)

        if package_name and action in {"pass", "fail", "skip"} and not test_name:
            previous = package_results.get(package_name)
            if previous is None or action == "fail":
                package_results[package_name] = (action, elapsed)

        if not package_name or not test_name or action not in {"pass", "fail", "skip"}:
            continue

        test_key = (package_name, test_name)
        if test_key in finalized_tests:
            continue

        finalized_tests.add(test_key)
        metrics["tests_run"] += 1
        if action == "fail":
            metrics["failures"] += 1
            failed_tests.append(f"{package_name}:{test_name}")
        elif action == "skip":
            metrics["skipped"] += 1

    metrics["packages_total"] = len(package_results)
    metrics["packages_failed"] = sum(1 for action, _ in package_results.values() if action == "fail")
    metrics["duration_seconds"] = sum(duration for _, duration in package_results.values())

    ordered_packages = sorted(
        ((package_name, action, duration) for package_name, (action, duration) in package_results.items()),
        key=lambda item: (-item[2], item[0]),
    )

    return metrics, ordered_packages, failed_tests


def publish_go_test_summary(
    *,
    title: str,
    report_file: str | None,
    summary_file: str | None,
    output_file: str | None,
    detected: bool,
    empty_message: str,
) -> int:
    if not detected:
        write_output(output_file, "status", "skipped")
        for key in ("tests_run", "failures", "skipped", "packages_total", "packages_failed", "duration_seconds"):
            write_output(output_file, key, "0")
        append_summary(
            summary_file,
            "\n".join(
                [
                    f"## {title}",
                    "",
                    empty_message,
                    "",
                ]
            ),
        )
        return 0

    metrics, packages, failed_tests = parse_go_test_report(report_file)
    status = "failure" if metrics["failures"] or metrics["packages_failed"] else "success"
    if metrics["tests_run"] == 0 and metrics["packages_total"] == 0:
        status = "empty"

    for key, value in metrics.items():
        write_output(output_file, key, str(value))
    write_output(output_file, "status", status)
    write_output(output_file, "failed_tests_count", str(len(failed_tests)))

    package_rows = "\n".join(
        f"| `{package_name}` | {action.upper()} | {format_duration(duration)} |"
        for package_name, action, duration in packages[:10]
    ) or "| No packages reported | n/a | 0.00s |"

    failing_rows = "\n".join(f"| `{test_name}` |" for test_name in failed_tests[:10]) or "| None |"
    if status == "empty":
        headline = "The suite was detected but no Go tests were executed."
    elif status == "failure":
        headline = "At least one test or package failed."
    else:
        headline = "The suite completed successfully."

    append_summary(
        summary_file,
        "\n".join(
            [
                f"## {title}",
                "",
                headline,
                "",
                "| Metric | Value |",
                "| --- | --- |",
                f"| Tests run | {metrics['tests_run']} |",
                f"| Failures | {metrics['failures']} |",
                f"| Skipped | {metrics['skipped']} |",
                f"| Packages reported | {metrics['packages_total']} |",
                f"| Failed packages | {metrics['packages_failed']} |",
                f"| Duration | {format_duration(float(metrics['duration_seconds']))} |",
                "",
                "| Package | Result | Duration |",
                "| --- | --- | ---: |",
                package_rows,
                "",
                "| Failing tests |",
                "| --- |",
                failing_rows,
                "",
            ]
        ),
    )
    return 0


BENCHMARK_LINE_RE = re.compile(
    r"^(Benchmark\S+)\s+\d+\s+([0-9.]+)\s+ns/op(?:\s+([0-9.]+)\s+B/op)?(?:\s+([0-9.]+)\s+allocs/op)?$"
)


def publish_benchmark_summary(args: argparse.Namespace) -> int:
    detected = parse_bool(args.detected)
    if not detected:
        write_output(args.output_file, "status", "skipped")
        for key in ("benchmarks_total", "slowest_ns_per_op", "largest_alloc_bytes", "largest_allocs_per_op"):
            write_output(args.output_file, key, "0")
        write_output(args.output_file, "slowest_name", "")
        append_summary(
            args.summary_file,
            "\n".join(
                [
                    "## Go Performance Benchmarks",
                    "",
                    "No Go benchmark suites were detected. Add `Benchmark...` functions in `*_test.go` files to enable this stage.",
                    "",
                ]
            ),
        )
        return 0

    current_package = "n/a"
    benchmarks: list[dict[str, float | str]] = []

    for raw_line in read_lines(args.log_file):
        line = raw_line.strip()
        if line.startswith("pkg: "):
            current_package = line.removeprefix("pkg: ").strip()
            continue

        match = BENCHMARK_LINE_RE.match(line)
        if not match:
            continue

        benchmarks.append(
            {
                "name": match.group(1),
                "package": current_package,
                "ns_per_op": float(match.group(2)),
                "bytes_per_op": float(match.group(3)) if match.group(3) else 0.0,
                "allocs_per_op": float(match.group(4)) if match.group(4) else 0.0,
            }
        )

    if not benchmarks:
        status = "empty"
        slowest = {
            "name": "",
            "ns_per_op": 0.0,
            "bytes_per_op": 0.0,
            "allocs_per_op": 0.0,
        }
    else:
        status = "success"
        slowest = max(benchmarks, key=lambda item: float(item["ns_per_op"]))

    write_output(args.output_file, "status", status)
    write_output(args.output_file, "benchmarks_total", str(len(benchmarks)))
    write_output(args.output_file, "slowest_name", str(slowest["name"]))
    write_output(args.output_file, "slowest_ns_per_op", str(slowest["ns_per_op"]))
    write_output(args.output_file, "largest_alloc_bytes", str(max((float(item["bytes_per_op"]) for item in benchmarks), default=0.0)))
    write_output(args.output_file, "largest_allocs_per_op", str(max((float(item["allocs_per_op"]) for item in benchmarks), default=0.0)))

    slowest_rows = "\n".join(
        f"| `{item['package']}` | `{item['name']}` | {format_optional_number(float(item['ns_per_op']))} | {format_optional_number(float(item['bytes_per_op']))} | {format_optional_number(float(item['allocs_per_op']))} |"
        for item in sorted(benchmarks, key=lambda item: float(item["ns_per_op"]), reverse=True)[:10]
    ) or "| n/a | No benchmark lines parsed | 0 | 0 | 0 |"

    headline = (
        "The benchmark suite completed successfully."
        if status == "success"
        else "Benchmark files were found, but no benchmark results were emitted."
    )

    append_summary(
        args.summary_file,
        "\n".join(
            [
                "## Go Performance Benchmarks",
                "",
                headline,
                "",
                "| Metric | Value |",
                "| --- | --- |",
                f"| Benchmarks parsed | {len(benchmarks)} |",
                f"| Slowest benchmark | `{slowest['name'] or 'n/a'}` |",
                f"| Slowest time | {format_optional_number(float(slowest['ns_per_op']), ' ns/op')} |",
                f"| Largest allocation | {format_optional_number(max((float(item['bytes_per_op']) for item in benchmarks), default=0.0), ' B/op')} |",
                f"| Largest allocs | {format_optional_number(max((float(item['allocs_per_op']) for item in benchmarks), default=0.0), ' allocs/op')} |",
                "",
                "| Package | Benchmark | ns/op | B/op | allocs/op |",
                "| --- | --- | ---: | ---: | ---: |",
                slowest_rows,
                "",
            ]
        ),
    )
    return 0


TOTAL_COVERAGE_RE = re.compile(r"^total:\s+\(statements\)\s+([0-9.]+)%$")


def publish_coverage_summary(args: argparse.Namespace) -> int:
    threshold_pct = float(args.threshold_pct)
    coverage_pct = None

    for raw_line in read_lines(args.report_file):
        match = TOTAL_COVERAGE_RE.match(raw_line.strip())
        if match:
            coverage_pct = float(match.group(1))
            break

    status = "skipped"
    below_threshold = False
    if coverage_pct is not None:
        status = "success"
        if threshold_pct > 0 and coverage_pct < threshold_pct:
            status = "failure"
            below_threshold = True

    write_output(args.output_file, "status", status)
    write_output(args.output_file, "coverage_pct", str(coverage_pct if coverage_pct is not None else 0.0))
    write_output(args.output_file, "threshold_pct", str(threshold_pct))
    write_output(args.output_file, "below_threshold", str(below_threshold).lower())

    if coverage_pct is None:
        message = "No coverage summary file was generated."
    elif below_threshold:
        message = "Coverage completed, but the configured threshold was not met."
    else:
        message = "Coverage completed successfully."

    append_summary(
        args.summary_file,
        "\n".join(
            [
                "## Coverage Report",
                "",
                message,
                "",
                "| Metric | Value |",
                "| --- | --- |",
                f"| Total coverage | {format_optional_number(coverage_pct, '%')} |",
                f"| Threshold | {format_optional_number(threshold_pct, '%')} |",
                "",
            ]
        ),
    )
    return 0


def build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(dest="command", required=True)

    unit_parser = subparsers.add_parser("unit")
    unit_parser.add_argument("--report-file", required=True)
    unit_parser.add_argument("--summary-file")
    unit_parser.add_argument("--output-file")

    api_parser = subparsers.add_parser("api")
    api_parser.add_argument("--report-file", required=True)
    api_parser.add_argument("--summary-file")
    api_parser.add_argument("--output-file")
    api_parser.add_argument("--detected", default="false")

    benchmark_parser = subparsers.add_parser("benchmark")
    benchmark_parser.add_argument("--log-file", required=True)
    benchmark_parser.add_argument("--summary-file")
    benchmark_parser.add_argument("--output-file")
    benchmark_parser.add_argument("--detected", default="false")

    coverage_parser = subparsers.add_parser("coverage")
    coverage_parser.add_argument("--report-file", required=True)
    coverage_parser.add_argument("--summary-file")
    coverage_parser.add_argument("--output-file")
    coverage_parser.add_argument("--threshold-pct", default="0")

    return parser


def main() -> int:
    parser = build_parser()
    args = parser.parse_args()

    if args.command == "unit":
        return publish_go_test_summary(
            title="Unit Tests",
            report_file=args.report_file,
            summary_file=args.summary_file,
            output_file=args.output_file,
            detected=True,
            empty_message="No Go unit test report was generated.",
        )

    if args.command == "api":
        return publish_go_test_summary(
            title="API Behavior / BDD Tests",
            report_file=args.report_file,
            summary_file=args.summary_file,
            output_file=args.output_file,
            detected=parse_bool(args.detected),
            empty_message=(
                "No API behavior suite was detected. This stage looks for files such as "
                "`*_contract_test.go`, `*_api_test.go`, `*_bdd_test.go`, `*_godog_test.go`, "
                "or Go build tags like `contract`, `bdd`, `api`, or `integration`."
            ),
        )

    if args.command == "benchmark":
        return publish_benchmark_summary(args)

    if args.command == "coverage":
        return publish_coverage_summary(args)

    return 1


if __name__ == "__main__":
    raise SystemExit(main())
