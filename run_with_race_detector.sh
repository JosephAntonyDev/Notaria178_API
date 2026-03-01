#!/bin/bash
# -------------------------------------------------------------------
# Notaria 178 API — Race Detector
# Ejecuta la aplicación con el detector de carreras de Go habilitado.
# Uso: chmod +x run_with_race_detector.sh && ./run_with_race_detector.sh
# -------------------------------------------------------------------

set -euo pipefail

echo "=== Compilando con -race ==="
go build -race -o notaria178_race ./...

echo "=== Ejecutando con Race Detector habilitado ==="
./notaria178_race
