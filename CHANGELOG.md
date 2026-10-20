# Changelog

## Indice mensual

- [Octubre 2026](#octubre-2026)
- [Septiembre 2026](#septiembre-2026)

---


## Octubre 2026

### Resumen

Cambios realizados durante octubre de 2026.

---

## [20-10-26]

### Added

- Lista de prioridades añadida:
    - Añadido `--format json`.
    - Mantener `table` como formato predeterminado añadido.
    - El JSON incluye herramienta, dependencia, versión actual, estado, releases, lifecycle, vulnerabilidades y fuentes.
- Mejorar después el parser para permitir opciones tanto antes como después de la ruta.
    - Ejemplo de antes:
    ```bash
    compver scan --format json /tmp/compver-fixture
    ```
    - Ejemplo con la estructura actualizada:
    ```bash
    compver scan /tmp/compver-fixture --format json
    ```

- Añadido mejoras a la CLI de manera visual:
    - Colores por estado. Solo para salida terminal, JSON nunca tendra colores.
    - Iconos y bloques visuales.
    - Resumen inicial.
    - Argumento para quitar los colores con `--no-color`.
    - Resumen por herramienta añadido.
    - Tabla dinámica implementada.
      - Las columnas calcular su ancho según el contenido.
      - Los nombres largos se limitan y terminan en `…`.
      - Los estados extensos tienen un ancho máximo controlado.
      - Las columnas siguen alineadas.
      - El resumen por herramientas y los detalles inferiores se mantienen.
      - La salida JSON no se ve afectada por esta estructura.
    - Añadido argumentos para mejorar la visibilidad y maniobrabilidad:
      - `--verbose`: muestra todas las secciones detalladas de una herramienta o ruta.
      - `--summary`: muestra solo el resumen general y el resumen por herramienta.
      - `--quiet`: oculta la salida normal, manteniendo los errores.
      - Los tres modos son incompatibles entre si.
      - Se conserva el formato JSON mediante `--format-json`.
      - El resumen general vuelve a mostrar la ruta analizada.
      - Se reutiliza la agrupación de dependencias para evitar lógica duplicada.
      - Ejemplos:
        ```bash
        compver scan . --summary
        compver scan . --verbose
        compver scan . --quiet
        ```
    - Añadido argumentos para el uso de Automatizacion:
      - `--fail-on vulnerable`.
      - `--fail-on end-of-life`.
      - `--fail-on vulnerable,end-of-life`
    Tambien acepta:
      - `--fail-on eol`


### Changed

- Modificado archivo `scan_json_output` para arreglar fallo de reestructuracion a la hora de mostrar todas las versiones.


## [14-10-26]


### Added

- Añadido archivo para aplicar exclusiones al motor, añadiendo el archivo `ignore.go`
- Añadido exclusion predeterminadas al motor de `Compver`.
- Dichas exclusiones predeterminadas son:
    - `.git`
    - `node_modules`
    - `vendor`
    - `.cache`
    - `dist`
    - `build`
- Añadido la creacion de un archivo que ignore dependencias o tecnologias concretas, `.compverignore`.


### Changed



## [13-10-26]

### Added

- Añadido archivo `installed_provider_detector.go` en adaptador Terraform.  
- Añadido archivo `installed_provider_detector_test.go` en adaptador Terraform.  
- Añadido el detector de providers instalados de Terraform mediante la carpeta
  `.terraform/providers`.  
- Añadida la correlación de dependencias detectadas por distintas fuentes,
  agrupándolas por herramienta y nombre.  
- Añadido soporte para observaciones de versión:  
  - Restricción declarada.  
  - Versión bloqueada.  
  - Versión instalada.  
- Añadida la comparación explícita entre la versión instalada, bloqueada o
  declarada y las versiones disponibles.  
- Añadido el modelo `ReleaseInfo` para representar información detallada de
  cada versión:  
  - Ciclo de versión.  
  - Fecha de publicación.  
  - Estado de soporte.  
  - Fecha de fin de soporte.  
  - Vulnerabilidades asociadas.  
- Añadido el evaluador central de versiones con estados:  
  - `current`.  
  - `update-available`.  
  - `compatible-update-available`.  
  - `deprecated`.  
  - `end-of-life`.  
  - `vulnerable`.  
  - `vulnerable-update-available`.  
- Añadido soporte para proveedores de lifecycle.  
- Añadido proveedor de lifecycle para Terraform mediante `endoflife.date`.  
- Añadido proveedor de información de publicación para providers de Terraform
  mediante GitHub Releases.  
- Añadido soporte para providers de vulnerabilidades.  
- Añadido proveedor de vulnerabilidades para Go Modules mediante OSV.  
- Añadido proveedor de vulnerabilidades para Terraform mediante OSV.  
- Añadida información de versiones corregidas para las vulnerabilidades
  detectadas.  
- Ampliado el sistema de caché para almacenar releases, lifecycle y
  vulnerabilidades.  
- Añadido soporte de caché para los modos:  
  - Automático.  
  - `--offline`.  
  - `--refresh`.  
- Añadida la opción `--show` para controlar el nivel de detalle de la salida:  
  - `lifecycle`.  
  - `vulnerabilities`.  
  - `sources`.  
  - `releases`.  
- Añadida la posibilidad de combinar secciones mediante comas:  
  `--show lifecycle,vulnerabilities,sources`.  
- Añadida la visualización de metadatos de releases seleccionadas.  
- Añadida la visualización de fuentes de detección, versiones, lifecycle y
  seguridad.  
- Añadida la visualización del estado de caché de cada proveedor.  
- Añadida la visualización de vulnerabilidades y versiones en las que fueron
  corregidas.  

### Changed

- Modificado motor para separar las fases de:  
    - Detección.  
    - Correlación.  
    - Consulta de versiones.  
    - Enriquecimiento de lifecycle.  
    - Consulta de vulnerabilidades.  
    - Evaluación final.  
- Modificando el analizador central para registrar detectores y proveedores de forma modular.  
- Modificado el adaptador de Terraform para analizar:  
    - Archivos `.tf`.  
    - `.terraform.lock.hcl`.  
    - `.terraform-version`.  
    - Providers instalados dentro de `.terraform/providers`.  
- Modificado el adaptador de Go para conservar la información de dependencias directas o indirectas.  
- Modificada la salida de la CLI para agrupar los resultados por herramienta.  
- La salida por defecto de la CLI ahora es mas compacta.  
- Los detalles de lifecycle, vulnerabilidades, fuentes y releases se muestran bajo demanda mediante `--show`.  
- Modificado el sistema de evaluación para distinguir entre la ultima versión publicada y la última versión compatible con las restricciones del proyecto.  
- Modificada la comparación para priorizar la versión instalada o bloqueada cuando está disponible.  
- Modificados los proveedores para reutilizar resultados almacenados en caché cuando corresponde.  

### Fixed

- Corregidos errores de sintaxis y orientación en la ejecución del análisis.  
- Corregida la selección de versiones compatibles con las restricciones declaradas.  
- Corregida la agrupación de dependencias procedentes de diferentes fuentes.  
- Corregida la representación de versiones estables y pre-release.  
- Corregida la conservación de información de releases al utilizar caché.  
- Corregida la salida de errores de proveedores y consultas externas.  



## [12-10-26]

### Added

- Añadido el detector de restricciones `required_providers` en archivos `.tf`.
- Añadido el campo `Constraint` al modelo de dependencias.
- Modificado y añadido diferentes condiciones normalizadas. Agregando soporte para versiones: 
    * `=`  
    * `!=`  
    * `>`  
    * `>=`  
    * `<`  
    * `<=`  
    * `~>`  
- Añadido filtrado común mediante `FilterByConstraint`.  
- Añadido `Constraints` en el adaptador de Terraform para traducir expresiones de versiones.  
- Las versiones `Constraints` ahora son genericas, y solo es necesario agregarlas a la herramienta que necesites.
- Añadido archivo para contribucion `CONTRIBUTING.md`.
- Añadido archivo de seguridad `SECURITY.md`.
- Añadido explicacion de como crear PRs `PULL_REQUEST_TEMPLATE.md`.
- Añadido workflow de validacion `validation-ci.yaml`.

### Changed

- Modificado `config_detector.go` para convertir las restricciones de `*tf`.  
- Modificado `scan_output.go` para mostrar únicamente las versiones compatibles.  
- Modificado los diferentes Tests del motor de restricciones y del parser de Terraform.  




## [09-10-26]

### Added

- Añadido el adaptador de Terraform mediante `.terraform-version`.
- Añadido el proveedor de releases oficial de HashiCorp.
- Integrada la caché común en el proveedor de Terraform.
- Añadido soporte para `compver scan --tool terraform`.
- Añadidas pruebas para el detector y proveedor de Terraform.
- Añadido el detector de `.terraform.lock.hcl`.
- Añadido el proveedor del Terraform Registry para consultar providers.
- Añadidas pruebas para providers bloqueados y consultas al Registry.
- Creado archivo `scan_output.go`.
- `spinner.go` añadido como metodo de visualizacion durante el analisis.

### Changed

- El proveedor de Go Modules ahora solo procesa dependencias detectadas en
  archivos `go.mod`.

- Cambiado la logica del archivo `scan.go` para que solo ejecute analisis de versiones.
- Añadido la logica de agrupar y mostrar los resultados de `scan.go` al archivo `scan_output.go`.
- Cambios `model.Dependency.Tool` para identificar el grupo que pertenece cada resultado.




## [08-10-26]


### Added

- Mejora del selector de versiones.
- Mejorada la comparación con la version actual.
- Estados de comparación.
- Separación entre versiones disponibles y seleccionadas.
- Mejora del registro de proveedores.
- Mejora de consulta al proxy de Go.
- El proveedor Go ahora utiliza cache.cache.
- Añadido nuevos argumentos:
    - [`auto`] usa caché válida y consulta internet si caduca.
    - [`refresh`] ignora la caché.
    - [`offline`] no realiza conexiones y usa solo datos almacenados.


### Changed

---


## [05-10-26]


### Added

- Creado y añadido el contrato de proveedor para usar conexion ***solo*** para verificacion de versiones.
- Creado archivo `provider.go` en el modulo de GO para probar conectividad.
- Ahora el motor consulta versiones disponibles durante el analisis.
- Modo `Dependency` agregado:
    * Versiones disponibles.
    * Fuente consultada.
    * Errores del proveedor.
    * Si el resultado procede de caché.
- `internal/app/analyzer.go` centraliza detector y proveedor Go.
- `scan.go` utiliza ahora `app.NewAnalyzer()`.
- La CLI muestra el número de versiones remotas encontradas.
- Añadidas pruebas del proveedor y de la integración con el motor.
- Añadido un registro centralizado mediante `app.NewAnalyzer()`.
- Añadido extension de argumentos para el comando de `compver scan`.
    [`compver scan --limit <Numero>`]
    [`compver scan --include-prereleases`]
    [`compver scan --all-versions`]  
- Añadido un argumento nuevo para el comando `compver scan`.
    [`compver scan --tool <nombre>`] De momento solo existe compatibilidad con herramienta Go.
- Creacion de `ARCHITECTURE.md` para mostrar la arquitectura completa.


### Changed



## [01-10-26]


### Added

- Añadido en ***cli/*** varios archivos indicando el comando:  
    [`cli.go`]  
    [`help.go`]  
    [`scan.go`]  
    [`version.go`]  

- Mejora añadida en el archivo de `main.go`.
- Añadida explicacion sobre los comandos de la aplicacion.

### Changed  

Reformulacion del CLI

---


## Septiembre 2026

### Resumen

Cambios realizados durante septiembre de 2026.

---


## [30-09-26]


### Added

- Estructura inicial basada en worktrees.
- Diseño del motor genérico de análisis de versiones.
- Planificacion de adaptadores independientes.

- Creado el módulo inicial de Go.
- Añadido el motor genérico de análisis.
- Añadido el sistema de detectores extensibles.
- Añadido los modelos iniciales de dependencias y resultados.
- Añadida una prueba básica del motor.
- Añadido el primer adaptador gomod para Go.
- Añadido deteccion de // indirect en el adaptador de Go.
- Nuevo paquete ***internal/version***.
- Parseo de versiones como `v1.2.3` y `1.2.3-rc.1`.
- Clasificacion entre estables y pre-release.
- Comparación básica entre versiones.
- Integrado el análisis de versiones en el motor.
- Añadida la clasificación de versiones estables, alpha, beta, rc y pre-release.
- Añadida la clasificación `unknown` para versiones no reconocidas.


### Changed


Añadida nueva meta al README.md
Cambio en `main.go` para que la salida de CLI sea mas humana y consistente.
La CLI muestra ahora el estado de la versión detectada.

### Fixed

---
