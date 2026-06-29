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

GROUP_FEAT="### 新功能"
GROUP_FIX="### Bug 修复"
GROUP_PERF="### 性能优化"
GROUP_REFACTOR="### 重构"
GROUP_DOCS="### 文档"
GROUP_TEST="### 测试"
GROUP_BUILD="### 构建"
GROUP_CI="### CI"
GROUP_CHORE="### 杂项"
GROUP_STYLE="### 样式"
GROUP_REVERT="### 回滚"
GROUP_OTHER="### 其他"

# 读入提交：格式 "<short-sha>|<subject>"
mapfile -t LINES < <(git log --no-merges --pretty=format:'%h|%s' "$RANGE")

BUCKET_FEAT=""
BUCKET_FIX=""
BUCKET_PERF=""
BUCKET_REFACTOR=""
BUCKET_DOCS=""
BUCKET_TEST=""
BUCKET_BUILD=""
BUCKET_CI=""
BUCKET_CHORE=""
BUCKET_STYLE=""
BUCKET_REVERT=""
OTHERS=""

for line in "${LINES[@]}"; do
  sha="${line%%|*}"
  msg="${line#*|}"
  # 解析前缀：type(scope)?!?:
  if [[ "$msg" =~ ^([a-zA-Z]+)(\([^)]+\))?!?:\ (.+)$ ]]; then
    type="${BASH_REMATCH[1],,}"
    scope="${BASH_REMATCH[2]}"
    desc="${BASH_REMATCH[3]}"
    item="- ${scope:+**${scope//[()]/}**: }${desc} (${sha})"
    case "$type" in
      feat)     BUCKET_FEAT+="${item}"$'\n' ;;
      fix)      BUCKET_FIX+="${item}"$'\n' ;;
      perf)     BUCKET_PERF+="${item}"$'\n' ;;
      refactor) BUCKET_REFACTOR+="${item}"$'\n' ;;
      docs)     BUCKET_DOCS+="${item}"$'\n' ;;
      test)     BUCKET_TEST+="${item}"$'\n' ;;
      build)    BUCKET_BUILD+="${item}"$'\n' ;;
      ci)       BUCKET_CI+="${item}"$'\n' ;;
      chore)    BUCKET_CHORE+="${item}"$'\n' ;;
      style)    BUCKET_STYLE+="${item}"$'\n' ;;
      revert)   BUCKET_REVERT+="${item}"$'\n' ;;
      *)        OTHERS+="- ${msg} (${sha})"$'\n' ;;
    esac
  else
    OTHERS+="- ${msg} (${sha})"$'\n'
  fi
done

{
  echo "## ${CURRENT_TAG}"
  echo
  if [[ -n "$BUCKET_FEAT" ]]; then     echo "$GROUP_FEAT"; echo; echo -n "$BUCKET_FEAT"; echo; fi
  if [[ -n "$BUCKET_FIX" ]]; then      echo "$GROUP_FIX"; echo; echo -n "$BUCKET_FIX"; echo; fi
  if [[ -n "$BUCKET_PERF" ]]; then     echo "$GROUP_PERF"; echo; echo -n "$BUCKET_PERF"; echo; fi
  if [[ -n "$BUCKET_REFACTOR" ]]; then echo "$GROUP_REFACTOR"; echo; echo -n "$BUCKET_REFACTOR"; echo; fi
  if [[ -n "$BUCKET_DOCS" ]]; then     echo "$GROUP_DOCS"; echo; echo -n "$BUCKET_DOCS"; echo; fi
  if [[ -n "$BUCKET_TEST" ]]; then     echo "$GROUP_TEST"; echo; echo -n "$BUCKET_TEST"; echo; fi
  if [[ -n "$BUCKET_BUILD" ]]; then    echo "$GROUP_BUILD"; echo; echo -n "$BUCKET_BUILD"; echo; fi
  if [[ -n "$BUCKET_CI" ]]; then       echo "$GROUP_CI"; echo; echo -n "$BUCKET_CI"; echo; fi
  if [[ -n "$BUCKET_CHORE" ]]; then    echo "$GROUP_CHORE"; echo; echo -n "$BUCKET_CHORE"; echo; fi
  if [[ -n "$BUCKET_STYLE" ]]; then    echo "$GROUP_STYLE"; echo; echo -n "$BUCKET_STYLE"; echo; fi
  if [[ -n "$BUCKET_REVERT" ]]; then   echo "$GROUP_REVERT"; echo; echo -n "$BUCKET_REVERT"; echo; fi
  if [[ -n "$OTHERS" ]]; then          echo "$GROUP_OTHER"; echo; echo -n "$OTHERS"; echo; fi
  if [[ -n "$COMPARE_URL_LINE" ]]; then
    echo "$COMPARE_URL_LINE"
  fi
}
