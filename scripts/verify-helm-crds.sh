#!/usr/bin/env bash

set -euo pipefail

chart_dir="${1:-charts/adcs-issuer}"
rendered="$(mktemp)"
disabled="$(mktemp)"
trap 'rm -f "$rendered" "$disabled"' EXIT

helm template adcs-issuer "$chart_dir" >"$rendered"

for crd in \
  adcsissuers.adcs.certmanager.csf.nokia.com \
  adcsrequests.adcs.certmanager.csf.nokia.com \
  clusteradcsissuers.adcs.certmanager.csf.nokia.com; do
  if ! grep -q "name: ${crd}" "$rendered"; then
    echo "missing rendered CRD: ${crd}" >&2
    exit 1
  fi
done

if [[ "$(grep -c '^kind: CustomResourceDefinition$' "$rendered")" -ne 3 ]]; then
  echo "expected exactly three rendered ADCS CRDs" >&2
  exit 1
fi

if [[ "$(grep -c '^              caBundleRef:' "$rendered")" -ne 2 ]]; then
  echo "expected caBundleRef in both issuer CRD schemas" >&2
  exit 1
fi

helm template adcs-issuer "$chart_dir" --set crd.install=false >"$disabled"
if grep -q '^kind: CustomResourceDefinition$' "$disabled"; then
  echo "CRDs rendered with crd.install=false" >&2
  exit 1
fi
