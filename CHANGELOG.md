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
