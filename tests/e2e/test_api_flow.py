"""
Suite E2E — Notaría 178 API
============================
Ejecutar con:
    pytest tests/e2e/test_api_flow.py -v --tb=short

Requiere la API corriendo en localhost:8080 y un usuario SUPER_ADMIN
previamente registrado en la BD. Configura las credenciales abajo.
"""

import os
import uuid
import pytest
import requests

# ─── Configuración ──────────────────────────────────────────────────────────

BASE_URL = os.getenv("API_BASE_URL", "http://localhost:8080/api/v1")

# Credenciales del SUPER_ADMIN existente en la BD
ADMIN_EMAIL = os.getenv("ADMIN_EMAIL", "admin@notaria178.com")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD", "admin123")


# ─── Estado compartido entre tests (orden importa) ─────────────────────────

class SharedState:
    """Almacena datos entre tests que se ejecutan en orden."""
    token: str = ""
    branch_id: str = ""
    client_id: str = ""
    act_id: str = ""
    work_id: str = ""


state = SharedState()


# ─── Helpers ────────────────────────────────────────────────────────────────

def auth_headers() -> dict:
    return {"Authorization": f"Bearer {state.token}"}


# ─── TEST 1: Login como SUPER_ADMIN ────────────────────────────────────────

class TestT01Login:
    """Autentica al usuario administrador y obtiene el JWT."""

    def test_login_super_admin(self):
        resp = requests.post(f"{BASE_URL}/users/login", json={
            "email": ADMIN_EMAIL,
            "password": ADMIN_PASSWORD,
        })
        assert resp.status_code == 200, f"Login falló: {resp.text}"
        data = resp.json()
        assert "token" in data, "La respuesta no contiene el campo 'token'"
        state.token = data["token"]


# ─── TEST 2: Middleware — acceso sin token ──────────────────────────────────

class TestT02Middleware:
    """Verifica que las rutas protegidas rechacen peticiones sin token."""

    def test_unauthorized_without_token(self):
        resp = requests.get(f"{BASE_URL}/users/profile")
        assert resp.status_code == 401, (
            f"Se esperaba 401 Unauthorized, se obtuvo {resp.status_code}"
        )


# ─── TEST 3: Crear Cliente ─────────────────────────────────────────────────

class TestT03CreateClient:
    """Crea un cliente de prueba y guarda su ID."""

    def test_create_client(self):
        unique = uuid.uuid4().hex[:8]
        resp = requests.post(
            f"{BASE_URL}/clients/create",
            headers=auth_headers(),
            json={
                "full_name": f"Cliente E2E {unique}",
                "rfc": f"RFC{unique[:10].upper()}",
                "phone": "5551234567",
                "email": f"e2e_{unique}@test.com",
            },
        )
        assert resp.status_code in (200, 201), f"Crear cliente falló: {resp.text}"
        data = resp.json()
        # La respuesta puede venir en data.id o data.data.id
        client_id = data.get("id") or data.get("data", {}).get("id")
        assert client_id, f"No se obtuvo client_id. Respuesta: {data}"
        state.client_id = client_id


# ─── TEST 4: Crear Expediente (Work) ───────────────────────────────────────

class TestT04CreateWork:
    """Crea un expediente de prueba. Necesita branch_id y act_id existentes."""

    def _get_or_create_branch(self):
        """Obtiene la primera sucursal o crea una."""
        resp = requests.get(
            f"{BASE_URL}/branches/search",
            headers=auth_headers(),
        )
        assert resp.status_code == 200
        data = resp.json()
        items = data.get("data", [])
        if items:
            state.branch_id = items[0]["id"]
            return

        # Crear sucursal si no existe
        resp = requests.post(
            f"{BASE_URL}/branches/create",
            headers=auth_headers(),
            json={"name": "Sucursal E2E", "address": "Dirección de prueba"},
        )
        assert resp.status_code in (200, 201), f"Crear branch falló: {resp.text}"
        branch_data = resp.json()
        state.branch_id = branch_data.get("id") or branch_data.get("data", {}).get("id")

    def _get_or_create_act(self):
        """Obtiene el primer acto del catálogo o crea uno."""
        resp = requests.get(
            f"{BASE_URL}/acts/search",
            headers=auth_headers(),
        )
        assert resp.status_code == 200
        data = resp.json()
        items = data.get("data", [])
        if items:
            state.act_id = items[0]["id"]
            return

        resp = requests.post(
            f"{BASE_URL}/acts/create",
            headers=auth_headers(),
            json={"name": f"Acto E2E {uuid.uuid4().hex[:6]}", "description": "Acto de prueba"},
        )
        assert resp.status_code in (200, 201), f"Crear act falló: {resp.text}"
        act_data = resp.json()
        state.act_id = act_data.get("id") or act_data.get("data", {}).get("id")

    def test_create_work(self):
        self._get_or_create_branch()
        self._get_or_create_act()

        assert state.branch_id, "No se pudo obtener branch_id"
        assert state.act_id, "No se pudo obtener act_id"
        assert state.client_id, "No se pudo obtener client_id (test anterior falló)"

        resp = requests.post(
            f"{BASE_URL}/works/create",
            headers=auth_headers(),
            json={
                "branch_id": state.branch_id,
                "client_id": state.client_id,
                "act_ids": [state.act_id],
            },
        )
        assert resp.status_code in (200, 201), f"Crear work falló: {resp.text}"
        data = resp.json()
        work_id = data.get("id") or data.get("data", {}).get("id")
        assert work_id, f"No se obtuvo work_id. Respuesta: {data}"
        state.work_id = work_id


# ─── TEST 5: Cambiar Status del Expediente ─────────────────────────────────

class TestT05UpdateWorkStatus:
    """Cambia el estado de PENDING → IN_PROGRESS. Dispara auditoría en backend."""

    def test_update_status(self):
        assert state.work_id, "No se pudo obtener work_id (test anterior falló)"

        resp = requests.patch(
            f"{BASE_URL}/works/status/{state.work_id}",
            headers=auth_headers(),
            json={"status": "IN_PROGRESS"},
        )
        assert resp.status_code == 200, f"Cambiar status falló: {resp.text}"
        data = resp.json()
        # Verificar que el status se actualizó
        new_status = data.get("status") or data.get("data", {}).get("status")
        assert new_status == "IN_PROGRESS", (
            f"Se esperaba IN_PROGRESS, se obtuvo: {new_status}"
        )


# ─── TEST 6: Verificar Auditoría ───────────────────────────────────────────

class TestT06VerifyAudit:
    """Consulta los logs de auditoría y verifica que exista el STATUS_CHANGE."""

    def test_audit_log_exists(self):
        assert state.work_id, "No se pudo obtener work_id (test anterior falló)"

        resp = requests.get(
            f"{BASE_URL}/audit/search",
            headers=auth_headers(),
            params={
                "entity_id": state.work_id,
                "action": "STATUS_CHANGE",
                "limit": 10,
            },
        )
        assert resp.status_code == 200, f"Buscar audit falló: {resp.text}"
        data = resp.json()
        logs = data.get("data", [])
        assert len(logs) > 0, (
            f"No se encontró registro de auditoría STATUS_CHANGE para work {state.work_id}"
        )

        # Verificar contenido del primer log
        log = logs[0]
        assert log["action"] == "STATUS_CHANGE"
        assert log["entity"] == "WORK"
        assert log["entity_id"] == state.work_id
