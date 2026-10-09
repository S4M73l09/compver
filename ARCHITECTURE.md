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
  │    └── Caché
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
├── cmd/                                   # Puntos de entrada de los ejecutables.
│   └── compver/                           # Ejecutable principal de Compver.
│       └── main.go                        # Inicia la CLI.
├── internal/                              # Código interno de la aplicación.
│   ├── adapters/                          # Integraciones específicas por herramienta.
│   │   ├── gomod/                         # Adaptador para dependencias Go.
│   │   │   ├── detector.go                # Detecta dependencias en go.mod.
│   │   │   ├── detector_test.go           # Tests del detector de Go.
│   │   │   ├── provider.go                # Consulta versiones en proxy.golang.org.
│   │   │   └── provider_test.go           # Tests del proveedor de Go.
│   │   └── terraform/                     # Adaptador para Terraform.
│   │       ├── detector.go                # Detecta la versión de Terraform.
│   │       ├── detector_test.go           # Tests del detector de Terraform.
│   │       ├── lock_detector.go           # Lee versiones desde .terraform.lock.hcl.
│   │       ├── lock_detector_test.go      # Tests del detector de bloqueo.
│   │       ├── provider.go                # Consulta releases de Terraform.
│   │       ├── provider_test.go           # Tests del proveedor de Terraform.
│   │       ├── registry_provider.go       # Consulta providers del Registry.
│   │       └── registry_provider_test.go  # Tests del Registry.
│   ├── app/                               # Composición y registro de componentes.
│   │   └── analyzer.go                    # Registra detectores y proveedores.
│   ├── cache/                             # Sistema de caché persistente.
│   │   ├── cache.go                       # Contrato general de caché.
│   │   ├── file_cache.go                  # Caché basada en archivos JSON.
│   │   └── file_cache_test.go             # Tests de la caché.
│   ├── cli/                               # Interfaz de línea de comandos.
│   │   ├── cli.go                         # Procesa los comandos principales.
│   │   ├── help.go                        # Muestra la ayuda.
│   │   ├── scan.go                        # Ejecuta el análisis.
│   │   ├── scan_options.go                # Procesa las opciones de scan.
│   │   ├── scan_output.go                 # Agrupa y muestra resultados.
│   │   ├── spinner.go                     # Muestra progreso durante el análisis.
│   │   └── version.go                     # Muestra la versión de Compver.
│   ├── engine/                            # Motor genérico de análisis.
│   │   ├── engine.go                      # Coordina detectores y proveedores.
│   │   └── engine_test.go                 # Tests del motor.
│   ├── model/                             # Modelos de datos compartidos.
│   │   └── model.go                       # Dependencias y resultados.
│   ├── providers/                         # Contratos de proveedores externos.
│   │   └── provider.go                    # Interfaz y opciones de consulta.
│   └── version/                           # Lógica común de versiones.
│       ├── version.go                     # Parseo y clasificación.
│       ├── version_test.go                # Tests de versiones.
│       ├── selector.go                    # Selección y filtrado.
│       ├── comparison.go                  # Comparación de versiones.
│       └── comparison_test.go             # Tests de comparación.
├── README.md                              # Documentación de uso.
├── ARCHITECTURE.md                        # Diseño interno del proyecto.
├── CHANGELOG.md                           # Historial de cambios.
├── LICENSE                                # Licencia Apache 2.0.
└── go.mod                                 # Definición del módulo Go.
```

## `cmd/compver`

Contiene el punto de entrada del programa. `main.go` solo inicia la CLI y le
pasa los argumentos recibidos.

## `internal/cli`

Gestiona la interacción con el usuario:

- Interpreta comandos y argumentos.
- Procesa opciones como `--tool`, `--limit`, `--include-prereleases`,
  `--offline` y `--refresh`.
- Muestra resultados, ayuda y errores.

Comandos actuales:

```bash
compver
compver .
compver scan .
compver scan --tool go .
compver scan --tool terraform .
compver scan --limit 10 .
compver scan --include-prereleases .
compver scan --all-versions .
compver scan --offline .
compver scan --refresh .
compver version
compver help
```

La CLI no debe contener lógica específica de Go, Terraform, npm u otras
herramientas.

## `internal/app`

Contiene la composición principal de la aplicación. `analyzer.go` registra
los detectores y proveedores que utilizará el motor.

Actualmente registra los adaptadores de Go Modules y Terraform. Cuando no se
indica `--tool`, ambos detectores pueden participar en el análisis.

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

El detector de Terraform analiza archivos `.terraform-version` y obtiene la
versión fijada del binario de Terraform.

El detector de bloqueo analiza `.terraform.lock.hcl` y obtiene las versiones
exactas seleccionadas para los providers del proyecto.

El detector no consulta Internet. Su responsabilidad termina cuando devuelve
los datos encontrados en el proyecto.

## Proveedores

Los proveedores consultan registros externos para obtener versiones
disponibles.

El proveedor actual utiliza el proxy de módulos de Go:

```text
https://proxy.golang.org
```

El proveedor de Terraform consulta el índice oficial de releases de HashiCorp:

```text
https://releases.hashicorp.com/terraform/
```

Para los providers, el adaptador consulta el Terraform Registry:

```text
https://registry.terraform.io/v1/providers/{namespace}/{name}/versions
```

El proveedor extrae versiones estables y pre-releases del listado, dejando
que el selector común decida cuáles se muestran.

El flujo es:

```text
Dependencia actual
  ↓
Proveedor compatible
  ↓
Consulta de caché
  ├── Entrada válida → Resultado
  └── Sin entrada válida → Consulta HTTP
  ↓
Lista de versiones
  ↓
Resultado del proveedor
```

Los proveedores deben controlar timeouts, errores HTTP, cancelación mediante
`context.Context`, versiones no reconocidas y el modo offline.

Cuando obtienen versiones desde Internet, los proveedores pueden guardarlas en
la caché común. De esta forma, el motor no necesita conocer los detalles de
cada registro externo.

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

Para añadir otra herramienta se crearían un detector y un proveedor
específicos. El adaptador de Terraform contiene detectores y proveedores
separados para el binario y sus providers:

```text
internal/adapters/terraform/
├── detector.go
├── detector_test.go
├── lock_detector.go
├── lock_detector_test.go
├── provider.go
├── provider_test.go
├── registry_provider.go
└── registry_provider_test.go
```

El siguiente paso será analizar restricciones declaradas en archivos `*.tf`
y verificar los plugins instalados en `.terraform/providers`.

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

`internal/cache` proporciona una interfaz de caché independiente del
proveedor y una implementación persistente basada en archivos JSON. La caché
se guarda en el directorio de usuario de la aplicación, normalmente:

```text
~/.cache/compver/
```

Las claves se convierten en nombres de archivo seguros mediante un hash, por
lo que cada herramienta puede reutilizar la misma infraestructura sin
colisionar con otras.

La interfaz de proveedores contempla los modos `auto`, `offline` y
`refresh`:

- `auto`: usa una entrada reciente y consulta Internet si falta o está
  caducada.
- `offline`: no realiza consultas de red y depende exclusivamente de la
  caché.
- `refresh`: ignora la entrada existente, consulta Internet y actualiza la
  caché.

Ejemplos:

```bash
compver scan --offline .
compver scan --refresh .
```

El modo offline no realiza consultas de red. Si no existe una entrada de
caché válida, la dependencia se marca como no consultable y el resto del
análisis continúa.

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
- Comparar la versión actual con las versiones disponibles.
- Estados de comparación.
- Separación entre versiones disponibles y seleccionadas.
- Usar caché persistente compartida entre proveedores.
- Ejecutar los modos `auto`, `offline` y `refresh`.
- Detectar Terraform desde `.terraform-version`.
- Consultar releases de Terraform mediante el proveedor oficial de HashiCorp.
- Detectar providers y versiones bloqueadas desde `.terraform.lock.hcl`.
- Consultar versiones de providers mediante el Terraform Registry.

Está previsto añadir:

- Analizar restricciones de providers en archivos `*.tf`.
- Inspeccionar providers instalados en `.terraform/providers`.
- Soporte para npm.
- Salida JSON.
- Interfaz CLI con colores y tablas.
- Distribución como binario y paquete Linux.
