#!/usr/bin/env python3
import json
import os
import sys
from pathlib import Path


def load(path: Path):
    with path.open("r", encoding="utf-8") as f:
        return json.load(f)


def ratio(hits: int, tested: int) -> float:
    if tested <= 0:
        return 0.0
    return (hits / tested) * 100.0


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: summarize_results.py <run_dir>", file=sys.stderr)
        return 1

    run_dir = Path(sys.argv[1]).resolve()
    if not run_dir.is_dir():
        print(f"run_dir not found: {run_dir}", file=sys.stderr)
        return 1

    scan_html = load(run_dir / "scan_html.json")
    scan_attr = load(run_dir / "scan_attr.json")
    scan_script = load(run_dir / "scan_script.json")
    scan_dom = load(run_dir / "scan_dom.json")
    scan_waf = load(run_dir / "scan_waf.json")
    fuzz_vuln = load(run_dir / "fuzzer_vuln.json")
    fuzz_sanitized = load(run_dir / "fuzzer_sanitized.json")
    bruteforce = load(run_dir / "bruteforce_public.json")

    vuln_ratio = ratio(fuzz_vuln.get("hits", 0), fuzz_vuln.get("tested", 0))
    sanitized_ratio = ratio(fuzz_sanitized.get("hits", 0), fuzz_sanitized.get("tested", 0))

    strengths = []
    gaps = []

    if scan_html.get("reflected", 0) > 0 and scan_attr.get("reflected", 0) > 0 and scan_script.get("reflected", 0) > 0:
        strengths.append("反射型上下文（HTML/属性/脚本）均可命中基础反射检测。")
    else:
        gaps.append("部分反射上下文未命中，需增强参数注入与响应匹配策略。")

    if scan_waf.get("waf", {}).get("detected", False):
        strengths.append("WAF 探测链路在模拟 Cloudflare 阻断场景下可识别。")
    else:
        gaps.append("WAF 探测未命中，需检查特征匹配与探针有效性。")

    dom_findings = len(scan_dom.get("dom", {}).get("findings", []))
    if dom_findings > 0:
        strengths.append(f"DOM 分析可识别 source/sink 线索（当前 {dom_findings} 条）。")
    else:
        gaps.append("DOM 分析未产出线索，需增强 sink/source 规则。")

    if vuln_ratio >= 80.0:
        strengths.append(f"公开语料在可利用端点上反射命中率较高（{vuln_ratio:.2f}%）。")
    else:
        gaps.append(f"公开语料在可利用端点命中率偏低（{vuln_ratio:.2f}%），需要优化 payload 变换与请求策略。")

    if sanitized_ratio <= 5.0:
        strengths.append(f"在已转义端点上误报率可控（{sanitized_ratio:.2f}%）。")
    else:
        gaps.append(f"在已转义端点上仍有明显命中（{sanitized_ratio:.2f}%），需要降低误报。")

    summary = []
    summary.append("# XSStrike-go Benchmark Summary")
    summary.append("")
    summary.append(f"- Run directory: `{run_dir}`")
    summary.append(f"- Public corpus size: `{fuzz_vuln.get('tested', 0)}` payload tests")
    summary.append("")
    summary.append("## Metrics")
    summary.append("")
    summary.append(f"- Fuzzer (vulnerable endpoint): `{fuzz_vuln.get('hits', 0)}/{fuzz_vuln.get('tested', 0)}` (`{vuln_ratio:.2f}%`)")
    summary.append(f"- Fuzzer (sanitized endpoint): `{fuzz_sanitized.get('hits', 0)}/{fuzz_sanitized.get('tested', 0)}` (`{sanitized_ratio:.2f}%`)")
    summary.append(f"- Bruteforce tested: `{bruteforce.get('tested', 0)}`, hits: `{len(bruteforce.get('hits', []))}`")
    summary.append(f"- Scan generated candidates: `html={scan_html.get('generated_candidates', 0)}`, `attr={scan_attr.get('generated_candidates', 0)}`, `script={scan_script.get('generated_candidates', 0)}`")
    summary.append(f"- WAF detected: `{scan_waf.get('waf', {}).get('detected', False)}` (`{scan_waf.get('waf', {}).get('name', '')}`)")
    summary.append(f"- DOM findings: `{dom_findings}`")
    summary.append("")
    summary.append("## Strengths")
    summary.append("")
    if strengths:
        for item in strengths:
            summary.append(f"- {item}")
    else:
        summary.append("- 暂无明显优势，需补充更多检测与评测能力。")
    summary.append("")
    summary.append("## Gaps")
    summary.append("")
    if gaps:
        for item in gaps:
            summary.append(f"- {item}")
    else:
        summary.append("- 当前评测未发现明显短板。")
    summary.append("")
    summary.append("## Notes")
    summary.append("")
    summary.append("- 本评测聚焦静态响应与反射链路，不包含浏览器真实执行验证。")
    summary.append("- 建议结合真实目标环境与回归样本持续比较。")
    summary.append("")

    out = run_dir / "summary.md"
    out.write_text("\n".join(summary), encoding="utf-8")
    print(out)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
