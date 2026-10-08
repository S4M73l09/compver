# Arquitectura de Compver

## Objetivo

Compver es una herramienta CLI escrita en Go para analizar proyectos,
detectar dependencias y comparar sus versiones actuales con versiones
disponibles en registros externos.

La arquitectura está diseñada para permitir añadir nuevas herramientas sin
modificar el núcleo del motor.

## Flujo general

```text
Usuario
  ↓
CLI
  ↓
App / registro central
  ↓
Engine
  ├── Detectores
  ├── Proveedores
  ├── Analizador de versiones
  └── Selector de versiones
  ↓
Resultado
  ↓
CLI
```

## Estructura del proyecto

```text
compver/
├── cmd/
│   └── compver/
│       └── main.go
├── internal/
│   ├── adapters/
│   │   └── gomod/
│   │       ├── detector.go
│   │       ├── detector_test.go
│   │       ├── provider.go
│   │       └── provider_test.go
│   ├── app/
│   │   └── analyzer.go
│   ├── cli/
│   │   ├── cli.go
│   │   ├── help.go
│   │   ├── scan.go
│   │   ├── scan_options.go
│   │   └── version.go
│   ├── engine/
│   │   ├── engine.go
│   │   └── engine_test.go
│   ├── model/
│   │   └── model.go
│   ├── providers/
│   │   └── provider.go
│   └── version/
│       ├── version.go
│       ├── version_test.go
│       └── selector.go
├── README.md
├── ARCHITECTURE.md
├── CHANGELOG.md
└── go.mod
```

## `cmd/compver`

Contiene el punto de entrada del programa. `main.go` solo inicia la CLI y le
pasa los argumentos recibidos.

## `internal/cli`

Gestiona la interacción con el usuario:

- Interpreta comandos y argumentos.
- Procesa opciones como `--tool`, `--limit` y `--include-prereleases`.
- Muestra resultados, ayuda y errores.

Comandos actuales:

```bash
compver
compver .
compver scan .
compver scan --tool go .
compver scan --limit 10 .
compver scan --include-prereleases .
compver scan --all-versions .
compver version
compver help
```

La CLI no debe contener lógica específica de Go, Terraform, npm u otras
herramientas.

## `internal/app`

Contiene la composición principal de la aplicación. `analyzer.go` registra
los detectores y proveedores que utilizará el motor.

Actualmente registra el detector y el proveedor de Go Modules.

## `internal/engine`

Es el núcleo de Compver. Se encarga de:

- Recorrer la ruta indicada.
- Ejecutar los detectores registrados.
- Analizar las versiones actuales.
- Ejecutar los proveedores compatibles.
- Guardar las versiones disponibles.
- Conservar errores individuales sin detener todo el análisis.
- Devolver un resultado estructurado.

El motor trabaja mediante interfaces genéricas:

```go
type Detector interface {
    CanHandle(path string) bool
    Detect(path string) ([]model.Dependency, error)
}
```

```go
type Provider interface {
    CanHandle(dependency model.Dependency) bool
    AvailableVersions(
        ctx context.Context,
        dependency model.Dependency,
        options QueryOptions,
    ) (Result, error)
}
```

## Detectores

Los detectores identifican dependencias y versiones dentro del filesystem.

El detector actual analiza archivos `go.mod`, extrae dependencias y
distingue entre dependencias directas e indirectas.

El detector no consulta Internet. Su responsabilidad termina cuando devuelve
los datos encontrados en el proyecto.

## Proveedores

Los proveedores consultan registros externos para obtener versiones
disponibles.

El proveedor actual utiliza el proxy de módulos de Go:

```text
https://proxy.golang.org
```

El flujo es:

```text
Dependencia actual
  ↓
Proveedor Go
  ↓
Consulta HTTP
  ↓
Lista de versiones
  ↓
Resultado del proveedor
```

Los proveedores deben controlar timeouts, errores HTTP, cancelación mediante
`context.Context`, versiones no reconocidas y el modo offline.

## Versiones

`internal/version` contiene la lógica común para trabajar con versiones.

Actualmente permite:

- Analizar versiones SemVer.
- Reconocer versiones estables.
- Reconocer `alpha`, `beta` y `rc`.
- Reconocer otras pre-releases.
- Marcar versiones no reconocidas como `unknown`.
- Comparar versiones.
- Seleccionar las versiones que se mostrarán.

Ejemplos:

```text
v1.2.3          → stable
v1.2.3-alpha.1 → alpha
v1.2.3-beta.1  → beta
v1.2.3-rc.1    → rc
v1.2.3-nightly → pre-release
```

## Selección de versiones

El proveedor puede devolver muchas versiones. El selector decide cuáles
mostrar:

```text
Versiones disponibles
  ↓
Filtrar pre-releases
  ↓
Ordenar versiones
  ↓
Aplicar límite
  ↓
Versiones seleccionadas
```

Comportamiento previsto:

```bash
compver scan .
compver scan --limit 10 .
compver scan --include-prereleases .
compver scan --all-versions .
```

## Modelo de datos

Cada dependencia puede contener la versión actual, las versiones disponibles,
la fuente consultada y los errores asociados a la consulta.

La diferencia entre versiones disponibles y seleccionadas es importante:

```text
AvailableVersions → todas las versiones devueltas por el proveedor
SelectedVersions  → versiones filtradas para mostrar al usuario
```

## Añadir una nueva herramienta

Para añadir Terraform, por ejemplo, se crearían un detector y un proveedor
específicos:

```text
internal/adapters/terraform/
├── detector.go
├── detector_test.go
├── provider.go
└── provider_test.go
```

El detector identificaría archivos como `*.tf` o `.terraform.lock.hcl`, y el
proveedor consultaría el registro correspondiente.

El motor no debería modificarse. Solo sería necesario:

1. Crear el detector.
2. Crear el proveedor.
3. Registrarlos en `internal/app/analyzer.go`.
4. Añadir tests.
5. Actualizar la documentación.

## Gestión de errores

Los errores de una dependencia concreta no deberían detener todo el análisis.
Por ejemplo, una dependencia puede consultarse correctamente, otra puede
obtenerse desde caché y otra puede devolver un timeout.

El error se conserva en el resultado para que la CLI pueda mostrarlo sin
perder el resto del análisis.

## Caché y modo offline

La interfaz de proveedores contempla los modos `auto`, `offline` y
`refresh`.

Está previsto implementar:

```bash
compver scan --offline .
compver scan --refresh .
```

El modo offline no debe realizar consultas de red. Si no existe una caché
válida, la dependencia se marcará como no consultable.

## Principios de diseño

- El motor no depende de herramientas concretas.
- Los detectores solo analizan archivos.
- Los proveedores gestionan consultas externas.
- La CLI no contiene lógica de negocio.
- Las versiones se procesan mediante lógica común.
- Los errores individuales no detienen todo el análisis.
- Los componentes deben poder probarse de forma independiente.
- La red debe ser opcional y controlable.
- No se deben enviar rutas ni contenido innecesario a servicios externos.

## Estado actual

Compver puede actualmente:

- Analizar rutas.
- Detectar dependencias Go desde `go.mod`.
- Reconocer dependencias directas e indirectas.
- Clasificar versiones.
- Consultar versiones mediante el proxy de Go.
- Aplicar límites y filtros básicos.
- Usar el comando `--tool go`.
- Ejecutar tests y validaciones con `go vet`.
- Seleccionar versiones.
- Comparacion con la versión actual entre la version estable mas reciente.
- Estados de comparación.
- Separacion entre versiones disponibles y seleccionadas.

Está previsto añadir:

- Caché persistente.
- Modo offline completo.
- Comparación con la versión estable más reciente.
- Soporte para Terraform.
- Soporte para npm.
- Salida JSON.
- Interfaz CLI con colores y tablas.
- Distribución como binario y paquete Linux.
