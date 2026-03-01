@echo off
REM -------------------------------------------------------------------
REM Notaria 178 API — Race Detector (Windows)
REM Ejecuta la aplicacion con el detector de carreras de Go habilitado.
REM Uso: run_with_race_detector.bat
REM -------------------------------------------------------------------

echo === Compilando con -race ===
go build -race -o notaria178_race.exe .
if %ERRORLEVEL% neq 0 (
    echo Error en la compilacion con -race.
    exit /b 1
)

echo === Ejecutando con Race Detector habilitado ===
notaria178_race.exe
