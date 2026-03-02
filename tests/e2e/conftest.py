"""
conftest.py — Configuración global de pytest + Generador de Reporte PDF
========================================================================
Recolecta resultados de cada test y genera un PDF ejecutivo al finalizar.
"""

import os
import datetime
from pathlib import Path

import pytest

# ─── Recolector de resultados ───────────────────────────────────────────────

_results: list[dict] = []

RESULTS_DIR = Path(__file__).parent / "results"


@pytest.hookimpl(tryfirst=True, hookwrapper=True)
def pytest_runtest_makereport(item, call):
    """Captura el resultado de cada fase del test (setup/call/teardown)."""
    outcome = yield
    report = outcome.get_result()

    # Solo nos interesa la fase "call" (ejecución real del test)
    if report.when == "call":
        _results.append({
            "name": item.name,
            "nodeid": report.nodeid,
            "outcome": report.outcome.upper(),  # PASSED / FAILED / SKIPPED
            "duration": round(report.duration, 3),
            "message": str(report.longrepr)[:200] if report.failed else "",
        })


def pytest_sessionfinish(session, exitstatus):
    """Al terminar la sesión de pytest, genera el reporte PDF."""
    _generate_pdf_report(_results, exitstatus)


# ─── Generador PDF ─────────────────────────────────────────────────────────

def _generate_pdf_report(results: list[dict], exit_status: int):
    try:
        from fpdf import FPDF
    except ImportError:
        print("\n[WARN] fpdf2 no instalado. Omitiendo generación de PDF.")
        print("       Instala con: pip install fpdf2")
        return

    RESULTS_DIR.mkdir(parents=True, exist_ok=True)
    pdf_path = RESULTS_DIR / "Test_Report_Notaria178.pdf"

    now = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
    total = len(results)
    passed = sum(1 for r in results if r["outcome"] == "PASSED")
    failed = sum(1 for r in results if r["outcome"] == "FAILED")
    skipped = sum(1 for r in results if r["outcome"] == "SKIPPED")

    pdf = FPDF()
    pdf.set_auto_page_break(auto=True, margin=20)
    pdf.add_page()

    # ── Título ──
    pdf.set_font("Helvetica", "B", 22)
    pdf.set_text_color(30, 60, 110)
    pdf.cell(0, 15, "Notaria 178 API", new_x="LMARGIN", new_y="NEXT", align="C")

    pdf.set_font("Helvetica", "", 14)
    pdf.set_text_color(80, 80, 80)
    pdf.cell(0, 8, "Reporte de Pruebas E2E", new_x="LMARGIN", new_y="NEXT", align="C")

    pdf.ln(5)

    # ── Línea separadora ──
    pdf.set_draw_color(30, 60, 110)
    pdf.set_line_width(0.8)
    pdf.line(15, pdf.get_y(), 195, pdf.get_y())
    pdf.ln(8)

    # ── Información general ──
    pdf.set_font("Helvetica", "", 11)
    pdf.set_text_color(40, 40, 40)
    pdf.cell(0, 7, f"Fecha de ejecucion:  {now}", new_x="LMARGIN", new_y="NEXT")
    pdf.cell(0, 7, f"Exit status pytest:  {exit_status}", new_x="LMARGIN", new_y="NEXT")
    pdf.ln(5)

    # ── Resumen ejecutivo ──
    pdf.set_font("Helvetica", "B", 13)
    pdf.set_text_color(30, 60, 110)
    pdf.cell(0, 10, "Resumen Ejecutivo", new_x="LMARGIN", new_y="NEXT")

    pdf.set_font("Helvetica", "", 11)
    pdf.set_text_color(40, 40, 40)
    pdf.cell(0, 7, f"Total de pruebas:   {total}", new_x="LMARGIN", new_y="NEXT")
    pdf.cell(0, 7, f"Exitosas (PASS):    {passed}", new_x="LMARGIN", new_y="NEXT")
    pdf.cell(0, 7, f"Fallidas (FAIL):    {failed}", new_x="LMARGIN", new_y="NEXT")
    pdf.cell(0, 7, f"Omitidas (SKIP):    {skipped}", new_x="LMARGIN", new_y="NEXT")
    pdf.ln(3)

    # ── Barra de progreso visual ──
    bar_w = 160
    bar_h = 8
    x_start = 25
    y_bar = pdf.get_y()

    if total > 0:
        pass_w = (passed / total) * bar_w
        fail_w = (failed / total) * bar_w
        skip_w = bar_w - pass_w - fail_w
    else:
        pass_w = fail_w = skip_w = 0

    # Verde (passed)
    pdf.set_fill_color(46, 160, 67)
    pdf.rect(x_start, y_bar, pass_w, bar_h, "F")
    # Rojo (failed)
    pdf.set_fill_color(207, 34, 46)
    pdf.rect(x_start + pass_w, y_bar, fail_w, bar_h, "F")
    # Gris (skipped)
    pdf.set_fill_color(180, 180, 180)
    pdf.rect(x_start + pass_w + fail_w, y_bar, skip_w, bar_h, "F")
    # Borde
    pdf.set_draw_color(100, 100, 100)
    pdf.rect(x_start, y_bar, bar_w, bar_h, "D")
    pdf.ln(14)

    # ── Middlewares probados ──
    pdf.set_font("Helvetica", "B", 13)
    pdf.set_text_color(30, 60, 110)
    pdf.cell(0, 10, "Middlewares Verificados", new_x="LMARGIN", new_y="NEXT")

    pdf.set_font("Helvetica", "", 10)
    pdf.set_text_color(40, 40, 40)
    middlewares = [
        "AuthMiddleware (JWT) — Verifica token Bearer en rutas protegidas",
        "RequireRoles — Restringe acceso por rol (SUPER_ADMIN, LOCAL_ADMIN)",
        "CORS — Configurado globalmente para peticiones cross-origin",
    ]
    for mw in middlewares:
        pdf.cell(5, 6, "-")
        pdf.cell(0, 6, f"  {mw}".replace("—", "-"), new_x="LMARGIN", new_y="NEXT")
    pdf.ln(5)

    # ── Tabla de resultados ──
    pdf.set_font("Helvetica", "B", 13)
    pdf.set_text_color(30, 60, 110)
    pdf.cell(0, 10, "Detalle de Pruebas", new_x="LMARGIN", new_y="NEXT")

    # Header de tabla
    col_widths = [10, 90, 25, 25, 30]
    headers = ["#", "Test", "Estado", "Tiempo", "Detalle"]
    pdf.set_font("Helvetica", "B", 9)
    pdf.set_fill_color(30, 60, 110)
    pdf.set_text_color(255, 255, 255)
    for i, h in enumerate(headers):
        pdf.cell(col_widths[i], 8, h, border=1, fill=True, align="C")
    pdf.ln()

    # Filas
    pdf.set_font("Helvetica", "", 9)
    for idx, r in enumerate(results, 1):
        # Colores según resultado
        if r["outcome"] == "PASSED":
            pdf.set_text_color(46, 160, 67)
            status_txt = "PASS"
        elif r["outcome"] == "FAILED":
            pdf.set_text_color(207, 34, 46)
            status_txt = "FAIL"
        else:
            pdf.set_text_color(150, 150, 150)
            status_txt = "SKIP"

        pdf.set_fill_color(245, 245, 250) if idx % 2 == 0 else pdf.set_fill_color(255, 255, 255)

        pdf.cell(col_widths[0], 7, str(idx), border=1, fill=True, align="C")

        pdf.set_text_color(40, 40, 40)
        name = r["name"][:40] + "..." if len(r["name"]) > 40 else r["name"]
        pdf.cell(col_widths[1], 7, name, border=1, fill=True)

        if r["outcome"] == "PASSED":
            pdf.set_text_color(46, 160, 67)
        elif r["outcome"] == "FAILED":
            pdf.set_text_color(207, 34, 46)
        else:
            pdf.set_text_color(150, 150, 150)
        pdf.cell(col_widths[2], 7, status_txt, border=1, fill=True, align="C")

        pdf.set_text_color(40, 40, 40)
        pdf.cell(col_widths[3], 7, f"{r['duration']}s", border=1, fill=True, align="C")

        detail = r["message"][:20] + "..." if r["message"] else "-"
        pdf.cell(col_widths[4], 7, detail, border=1, fill=True, align="C")
        pdf.ln()

    # ── Footer ──
    pdf.ln(10)
    pdf.set_font("Helvetica", "I", 8)
    pdf.set_text_color(150, 150, 150)
    pdf.cell(0, 5, f"Generado automaticamente por pytest + fpdf2  |  {now}", align="C")

    pdf.output(str(pdf_path))
    print(f"\n{'='*60}")
    print(f"  PDF generado: {pdf_path}")
    print(f"{'='*60}")
