#!/usr/bin/env bash
# Generate Markdown release notes from git log, grouped by Conventional Commits prefixes.
# Usage: gen-release-notes.sh <current_tag>
set -euo pipefail

CURRENT_TAG="${1:-}"
if [[ -z "$CURRENT_TAG" ]]; then
  CURRENT_TAG="$(git describe --tags --abbrev=0 2>/dev/null || echo HEAD)"
fi

# 上一个 tag（排除当前 tag 自身）；若无则回溯到首个 commit
PREV_TAG="$(git describe --tags --abbrev=0 "${CURRENT_TAG}^" 2>/dev/null || true)"
if [[ -n "$PREV_TAG" ]]; then
  RANGE="${PREV_TAG}..${CURRENT_TAG}"
  COMPARE_URL_LINE=""
  if [[ -n "${GITHUB_REPOSITORY:-}" ]]; then
    COMPARE_URL_LINE="**Full Changelog**: https://github.com/${GITHUB_REPOSITORY}/compare/${PREV_TAG}...${CURRENT_TAG}"
  fi
else
  RANGE="${CURRENT_TAG}"
  COMPARE_URL_LINE=""
fi

# 没传 tag（workflow_dispatch 且未传入），直接取最近 20 条
if [[ "$CURRENT_TAG" == "HEAD" ]]; then
  RANGE="-20"
fi

declare -A GROUPS=(
  [feat]="### 新功能"
  [fix]="### Bug 修复"
  [perf]="### 性能优化"
  [refactor]="### 重构"
  [docs]="### 文档"
  [test]="### 测试"
  [build]="### 构建"
  [ci]="### CI"
  [chore]="### 杂项"
  [style]="### 样式"
  [revert]="### 回滚"
)
ORDER=(feat fix perf refactor docs test build ci chore style revert)

# 读入提交：格式 "<short-sha>|<subject>"
mapfile -t LINES < <(git log --no-merges --pretty=format:'%h|%s' "$RANGE")

declare -A BUCKETS
OTHERS=()

for line in "${LINES[@]}"; do
  sha="${line%%|*}"
  msg="${line#*|}"
  # 解析前缀：type(scope)?!?:
  if [[ "$msg" =~ ^([a-zA-Z]+)(\([^)]+\))?!?:\ (.+)$ ]]; then
    type="${BASH_REMATCH[1],,}"
    scope="${BASH_REMATCH[2]}"
    desc="${BASH_REMATCH[3]}"
    item="- ${scope:+**${scope//[()]/}**: }${desc} (${sha})"
    if [[ -n "${GROUPS[$type]:-}" ]]; then
      BUCKETS[$type]+="${item}"$'\n'
    else
      OTHERS+=("- ${msg} (${sha})")
    fi
  else
    OTHERS+=("- ${msg} (${sha})")
  fi
done

{
  echo "## ${CURRENT_TAG}"
  echo
  for t in "${ORDER[@]}"; do
    if [[ -n "${BUCKETS[$t]:-}" ]]; then
      echo "${GROUPS[$t]}"
      echo
      echo -n "${BUCKETS[$t]}"
      echo
    fi
  done
  if (( ${#OTHERS[@]} > 0 )); then
    echo "### 其他"
    echo
    printf '%s\n' "${OTHERS[@]}"
    echo
  fi
  if [[ -n "$COMPARE_URL_LINE" ]]; then
    echo "$COMPARE_URL_LINE"
  fi
}
