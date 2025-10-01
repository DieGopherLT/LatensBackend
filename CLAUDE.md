# Latens Backend - Contexto Técnico

**Propósito**: Backend API en Go para Latens, herramienta que "despierta" proyectos dormidos mediante análisis inteligente del estado de desarrollo.

## Arquitectura Core

### Principios arquitectónicos
- **Clean Architecture**: Separación estricta entre capas de dominio, aplicación e infraestructura
- **Dependency Inversion**: Dependencias apuntan hacia abstracciones, no implementaciones
- **Repository Pattern**: Abstracción completa de la capa de persistencia
- **Context Propagation**: Manejo consistente de timeouts y cancelaciones

### Estructura de capas
```
cmd/api/              # Application layer - Bootstrap y configuración
internal/controller/  # Interface adapters - HTTP handlers 
internal/services/    # Application layer - Lógica de negocio
internal/database/    # Infrastructure - Persistencia y repositorios
internal/models/      # Domain layer - Entidades de dominio
pkg/                  # Shared kernel - Utilidades reutilizables
```

## Tech Stack Fundamental

- **Lenguaje**: Go (latest stable)
- **Framework web**: Fiber (rápido, Express-like)
- **Base de datos**: MongoDB (document-oriented, flexible schema)
- **Autenticación**: OAuth GitHub (integración nativa con plataforma)
- **IA**: OpenAI API (análisis semántico de código)

## Patrones de diseño establecidos

### Dependency Injection
- Constructor functions `NewXxx()` que reciben dependencias
- Configuración centralizada en `main.go`
- Interfaces para todas las abstracciones

### Error Handling
- Context-aware error propagation
- HTTP status codes semánticamente correctos
- Response format estandarizado con `error` y `details`

### Configuration Management
- Environment-based configuration
- Validation mediante struct tags
- Configuración inmutable post-bootstrap

## Información de contexto adicional

Para detalles específicos que evolucionan con el proyecto, el sistema Serena mantiene memorias actualizadas con:

- **project-overview** - Estado actual del proyecto, decisiones técnicas recientes, deuda técnica identificada
- **suggested_commands** - Comandos de desarrollo, scripts de build, herramientas específicas del entorno
- **code-style-conventions** - Convenciones de naming observadas, patrones de código aplicados
- **task-completion-checklist** - Flujo de trabajo establecido, validaciones pre-commit, proceso de QA

Estas memorias contienen información complementaria que se actualiza conforme evoluciona el proyecto.

## Principios de desarrollo

### Code Quality
- Adherencia estricta a Go idioms y convenciones estándar
- Zero tolerance para warnings de `go vet`
- Formatted code con `go fmt` antes de cada commit

### API Design
- RESTful endpoints con versionado (`/api/v1/`)
- Consistent JSON responses
- Proper HTTP semantics y status codes

### Database Design
- Document-oriented modeling apropiado para MongoDB
- Índices optimizados para patrones de consulta
- Schema evolution-friendly structure

## Contexto de negocio

Latens analiza proyectos "dormidos" y genera insights sobre:
- Estado del desarrollo cuando se pausó
- Decisiones técnicas pendientes
- Contexto necesario para reactivar el proyecto
- Score de "profundidad del sueño" basado en métricas

Este contexto informa las decisiones arquitectónicas hacia flexibilidad y capacidad de análisis retrospectivo.

## Fórmula de Sleep Score

### Lógica de Cálculo

**Decisión simplificada:**
```
if repo.isArchived || repo.isDisabled || repo.isFork:
    return 100

if !hasMultipleBranches || días_desde_último_commit > 180:
    return InactivityScore
else:
    return (0.5 × InactivityScore) + (0.3 × FragmentationScore) + (0.2 × StalenessScore)
```

**Resultado final:**
```
SleepScore = max(0, min(100, score))
```

### Componentes

**InactivityScore = min(100, (días_desde_último_commit / 180) × 100)**

Donde:
- días_desde_último_commit = hoy - defaultBranchRef.target.committedDate
- 180 días = Umbral de "profundamente dormido" (6 meses)

---

**FragmentationScore = min(100, (branches_adelantadas / 3) × 100)**

Donde:
- branches_adelantadas = COUNT(refs WHERE target.committedDate > defaultBranch.target.committedDate)
- 3 ramas = Umbral de "dispersión alta"

---

**StalenessScore = min(100, (avg_días_último_commit_branches / 90) × 100)**

Donde:
- avg_días_último_commit_branches = promedio de (hoy - target.committedDate)
  para todas las ramas con target.committedDate > defaultBranch.target.committedDate
- 90 días = Umbral de "ramas obsoletas"

---

### Pesos y Umbrales

**Pesos de componentes:**
- Inactividad: 50% (factor dominante)
- Fragmentación: 30% (indicador secundario)
- Obsolescencia: 20% (complementaria)

**Umbrales configurados:**
- Inactividad: 180 días (6 meses)
- Fragmentación: 3 branches
- Obsolescencia: 90 días (3 meses)

**Complejidad computacional:**
- Mejor caso (repo simple/viejo): O(1)
- Peor caso (repo complejo < 6 meses): O(n) donde n ≤ 50 branches
- Totalmente viable para cálculo síncrono durante sync