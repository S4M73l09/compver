# Changelog


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


---

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
