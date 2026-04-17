# Notaria 178 API

Backend REST del sistema privado de gestion notarial de la Notaria 178.

Esta API esta construida en Go con Gin, PostgreSQL y una organizacion por modulos con enfoque hexagonal (`domain`, `app`, `infra`). Hoy ya no solo cubre autenticacion, usuarios y expedientes: tambien incluye cache con Redis, comentarios en tiempo real por WebSocket, centro de notificaciones por SSE, push notifications opcionales con Firebase FCM, auditoria, dashboard y manejo de documentos con optimizacion de PDFs.

## Nota de nomenclatura

En negocio y en la UI el concepto actual es **oficinas**. Sin embargo, por compatibilidad tecnica el backend todavia usa nombres como `branch`, `branches` y `branch_id` en:

- tablas SQL
- endpoints `/branches/*`
- DTOs y modelos internos

Cuando en este README veas `branches`, leelo como **oficinas**.

## Stack actual

- Go `1.25.5`
- Gin `1.11.0`
- PostgreSQL
- Redis opcional para cache
- JWT para autenticacion
- BCrypt para hashing
- Gorilla WebSocket para comentarios en vivo
- SSE para notificaciones en tiempo real
- Firebase Admin SDK para push notifications opcionales
- pdfcpu para optimizacion de PDFs
- Docker Compose para levantar DB + Redis + API

## Variables de entorno

El proyecto usa `.env`. El nombre real de la variable de conexion es **`DB_URL`**, no `DATABASE_URL`.

```env
PORT=8080
DB_URL=postgres://usuario:password@localhost:5432/notaria178_db?sslmode=disable
JWT_SECRET=tu_secreto_jwt

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

FIREBASE_CREDENTIALS_PATH=./ruta_a_tu_service_account.json
```

### Variables usadas hoy

| Variable | Obligatoria | Uso |
| --- | --- | --- |
| `PORT` | No | Puerto HTTP. Default: `8080`. |
| `DB_URL` | Si | Conexion a PostgreSQL. |
| `JWT_SECRET` | Si | Firma y validacion de JWT. |
| `REDIS_ADDR` | No | Activa cache Redis si existe. |
| `REDIS_PASSWORD` | No | Password de Redis. |
| `REDIS_DB` | No | Indice de DB en Redis. |
| `FIREBASE_CREDENTIALS_PATH` | No | Habilita push notifications FCM si apunta a un service account valido. |

## Como correrlo

### Local

```bash
go mod tidy
go run .
```

La API levanta en `http://localhost:8080` salvo que cambies `PORT`.

### Docker

```bash
docker compose up -d --build
```

Servicios incluidos:

- `db` con `postgres:18-alpine`
- `redis` con `redis:7-alpine`
- `api` con la aplicacion Go

El `docker-compose.yml` ya monta:

- `./schema.sql` como init script de PostgreSQL
- `./uploads` para persistir archivos
- `./notaria178-firebase-admin.json` como credencial de Firebase dentro del contenedor

## Estructura actual

```text
Notaria178_API/
  main.go
  schema.sql
  docker-compose.yml
  Dockerfile
  cmd/
    migrator/
  internal/
    act/
    attendance/
    audit/
    branch/
    client/
    core/
      cache/
      dtos/
    dashboard/
    document/
    integration/
      adapters/
    messaging/
    middleware/
    notification/
    user/
    work/
  tests/
    e2e/
  uploads/
```

### Convencion por modulo

Cada modulo sigue esta forma:

```text
module/
  domain/   # entidades, puertos, reglas puras
  app/      # casos de uso
  infra/    # repositorios, controladores, rutas, adapters
```

## Modulos que existen hoy

### `user`

- login JWT
- perfil propio
- creacion y edicion de personal
- busqueda de usuarios

### `attendance`

- check-in / check-out
- historial propio
- consulta administrativa de asistencias

### `act`

- catalogo de actos
- categorias
- requisitos por acto
- soft delete / desactivacion cuando hay dependencias historicas
- cache Redis para busquedas

### `client`

- busqueda, alta y edicion de clientes

### `branch`

- catalogo tecnico de oficinas (`branches`)
- usado por filtros, asignaciones y dashboard

### `work`

Es el modulo central del sistema:

- crear y editar expedientes
- cambiar estado del trabajo
- asociar y quitar actos
- agregar y quitar colaboradores
- comentarios por expediente
- requisitos ad-hoc por expediente
- detalle completo del expediente

Estados manejados:

`PENDING -> IN_PROGRESS -> READY_FOR_REVIEW -> APPROVED / REJECTED`

### `document`

- subida de documentos
- listado por expediente
- descarga autenticada
- borrado de documentos
- soporte para documentos ligados a requisitos por `requirement_id` + `requirement_source`
- optimizacion automatica de PDFs con `pdfcpu`

### `notification`

- listado de notificaciones del usuario
- contador de no leidas
- marcar una o todas como leidas
- SSE en `/notifications/stream`
- registro de `device_token` para push
- integracion opcional con Firebase FCM

### `messaging`

- chat en vivo de comentarios de expedientes por WebSocket
- endpoint `/ws/comments`
- rooms por `work_id`
- reutiliza los mismos casos de uso de comentarios del modulo `work`

### `audit`

- busqueda de auditoria
- metricas de auditoria
- registro de acciones desde otros modulos via adapters

### `dashboard`

Agregaciones para panel administrativo:

- KPIs
- tendencia
- distribucion por estado
- actividad reciente
- top proyectistas
- top actos

Redis se usa aqui como cache opcional.

## Integraciones en tiempo real

Hoy existen dos canales distintos:

### SSE para notificaciones

- endpoint: `GET /notifications/stream`
- sirve para centro de notificaciones en tiempo real

### WebSocket para comentarios

- endpoint: `GET /ws/comments`
- se usa para entrar/salir de rooms por expediente y enviar comentarios en vivo

### Push notifications opcionales

Si `FIREBASE_CREDENTIALS_PATH` esta configurado y la credencial es valida:

- se registran tokens de dispositivo en `user_device_tokens`
- se pueden enviar push notifications de comentarios y eventos

Si Firebase no esta configurado, el sistema sigue funcionando con notificaciones in-app + SSE.

## Base de datos actual

`schema.sql` ya contempla, entre otras, estas tablas:

- `branches`
- `users`
- `attendances`
- `clients`
- `act_catalogs`
- `act_requirements`
- `works`
- `work_acts`
- `work_collaborators`
- `documents`
- `work_comments`
- `user_device_tokens`
- `notifications`
- `audit_logs`

Tambien incluye seed inicial de oficinas, usuario `SUPER_ADMIN`, proyectistas y catalogo base de actos/requisitos.

## Resumen de endpoints

Hoy el backend expone **55 endpoints/rutas** en total:

- `users`: login, perfil, creacion, actualizacion y busqueda
- `attendance`: check, historial propio e historial admin
- `acts`: CRUD del catalogo y requisitos
- `clients`: buscar, crear y actualizar
- `branches`: buscar, crear y actualizar oficinas
- `works`: busqueda, detalle, estados, comentarios, colaboradores, actos y requisitos
- `documents`: subir, listar, descargar y borrar
- `notifications`: listado, SSE, unread count, read, read-all y registro de device token
- `audit`: search y metrics
- `dashboard`: kpis, trend, distribution, activity, top-drafters y top-acts
- `ws/comments`: chat en vivo de comentarios

## Endpoints por modulo

### Users

- `POST /users/login`
- `GET /users/profile`
- `PATCH /users/profile`
- `POST /users/create`
- `PATCH /users/update/:id`
- `GET /users/search`

### Attendance

- `POST /attendance/check`
- `GET /attendance/history`
- `GET /attendance/admin/history/:id`

### Acts

- `GET /acts/search`
- `GET /acts/:id/requirements`
- `POST /acts/create`
- `PATCH /acts/update/:id`
- `PATCH /acts/status/:id`
- `DELETE /acts/:id`
- `POST /acts/:id/requirements`
- `DELETE /acts/:id/requirements/:req_id`

### Clients

- `GET /clients/search`
- `POST /clients/create`
- `PATCH /clients/update/:id`

### Offices (`branches`)

- `GET /branches/search`
- `POST /branches/create`
- `PATCH /branches/update/:id`

### Works

- `GET /works/search`
- `GET /works/:id`
- `POST /works/create`
- `PATCH /works/update/:id`
- `PATCH /works/status/:id`
- `GET /works/:id/comments`
- `POST /works/:id/comments`
- `POST /works/:id/collaborators`
- `DELETE /works/:id/collaborators/:userId`
- `POST /works/:id/acts`
- `DELETE /works/:id/acts/:actId`
- `POST /works/:id/requirements`
- `DELETE /works/:id/requirements/:reqId`

### Documents

- `POST /documents/upload`
- `GET /documents/work/:work_id`
- `GET /documents/download/:id`
- `DELETE /documents/:id`

### Notifications

- `GET /notifications`
- `GET /notifications/stream`
- `GET /notifications/unread-count`
- `PATCH /notifications/:id/read`
- `PATCH /notifications/read-all`
- `PUT /notifications/device-token`

### Audit

- `GET /audit/search`
- `GET /audit/metrics`

### Dashboard

- `GET /dashboard/kpis`
- `GET /dashboard/trend`
- `GET /dashboard/distribution`
- `GET /dashboard/activity`
- `GET /dashboard/top-drafters`
- `GET /dashboard/top-acts`

### Messaging

- `GET /ws/comments`

## Cache actual

Redis es opcional.

Si `REDIS_ADDR` no existe o Redis falla:

- la API sigue levantando
- actos y dashboard responden directo desde PostgreSQL

Hoy el cache se usa principalmente en:

- busquedas del modulo de actos
- agregaciones del dashboard

## Documentos y requisitos

La relacion actual ya es por IDs, no por nombre de archivo.

Campos relevantes en `documents`:

- `requirement_id`
- `requirement_source` con valores `ACT` o `WORK`

Esto permite:

- ligar archivos a requisitos de catalogo
- ligar archivos a requisitos ad-hoc del expediente
- reemplazar documentos sin depender de nombres "magicos"

## Calidad y soporte

### Race detector

```bash
./run_with_race_detector.sh
```

o en Windows:

```cmd
run_with_race_detector.bat
```

### E2E

```bash
pip install -r requirements-qa.txt
pytest tests/e2e/ -v
```

La suite actual vive en `tests/e2e/`.

### Migraciones auxiliares

Existen archivos de apoyo historicos:

- `migrate_acts.sql`
- `migrate_epic.sql`
- `cmd/migrator/main.go`

`cmd/migrator` es un helper local y hoy tiene una conexion hardcodeada para ejecutar `migrate_acts.sql`, asi que no reemplaza una herramienta formal de migraciones.

## Relacion con el frontend

Este backend sirve a `notaria178_frontend`.

Puntos importantes de integracion actuales:

- REST en `http://localhost:8080`
- SSE para notificaciones
- WebSocket en `ws://localhost:8080/ws/comments`
- Firebase FCM para push opcional
