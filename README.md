# Proyecto: Compver ***Rama: Pre-release***

En esta rama se propondran cambios en la estructura del funcionamiento de proyecto a la vez que se buscara el mantenimiento de este añadiendo nuevas funciones, cambios, solucion de bugs o errores, generacion de cambios y modificaciones varias.

Dicha rama contiene contenido en constante cambio, es libre del usuario descargar y usar las funciones tanto de prueba como normales del proyecto, pero igualmente puede contener errores y/o cambios futuros inesperados.

Dicho esto, esta bajo tu propio riesgo.


## Compver

Esta aplicacion ha sido desarrollada para implementar una solucion alternativa a la hora de usar innumerables herramientas que cuentan con versiones - pre-releases propias. Todas ellas contienen su propio codigo y es normal a la larga darse cuenta del caos de mantener versiones estables directamente a buscar versiones mas actuales sin necesidad de ir a lo ultimo.

Por ello se ha planteado esta solucion: ***Compver*** es una solucion de codigo abierto, es un comparador de versiones capaz de analizar la ruta que el usuario necesite o el filesystem del repositorio concreto donde sea necesario.

## Arquitectura

La arquitectura de Compver está diseñada para ser modular. El motor trabaja con detectores y proveedores mediante interfaces independientes.

Para conocer los detalles, consulta [ARCHITECTURE.md](ARCHITECTURE.md)


### Su funcion

***Compver*** compara las versiones usadas en una infraestructura concreta con las versiones mas estables y actuales de la misma herramienta, tambien da una lista sobre las pre-release que tienen las herramientas que este ha analizado. 

EL proyecto utiliza un motor generico, esto permite implementar adaptadores cuando se quiera añadir soporte a una nueva tecnologia, permitiendonos no necesitar ni modificar el motor.

En el futuro se buscara implementar y mostrar adaptadores de python para que cada usuario tenga la libertad de añadir soporte a su stack tecnologico concreto.

Su funcion ira evolucionando en funcion a como avance el proyecto.

### Comandos y funciones

`compver` analiza un directorio, detecta las herramientas y dependencias
presentes, consulta versiones disponibles y compara la versión actual con las
versiones publicadas y compatibles.

#### Analizar un proyecto

```bash
compver scan [ruta]
```

Si no se indica una ruta, se analiza el directorio actual:

```bash
compver scan
```

También existe un atajo para realizar un análisis:

```bash
compver
compver .
compver /ruta/al/proyecto
```

Estas formas son equivalentes:

```bash
compver scan .
compver scan /ruta/al/proyecto
```

#### Seleccionar una herramienta

Actualmente se puede limitar el análisis a Go o Terraform:

```bash
compver scan . --tool go
compver scan . --tool terraform
```

Si no se indica `--tool`, Compver analiza todas las herramientas compatibles
detectadas en la ruta.

#### Controlar las versiones mostradas

Por defecto se muestran las tres versiones estables más recientes:

```bash
compver scan . --limit 3
```

Para mostrar un número diferente:

```bash
compver scan . --limit 6
```

Para incluir versiones alpha, beta, rc y otras pre-releases:

```bash
compver scan . --include-prereleases
```

Para mostrar todas las versiones disponibles:

```bash
compver scan . --all-versions
```

Estas opciones pueden combinarse:

```bash
compver scan . \
  --tool terraform \
  --limit 10 \
  --include-prereleases
```

#### Mostrar información adicional

La salida por defecto es resumida. Se pueden solicitar secciones adicionales
mediante `--show`.

Mostrar releases seleccionadas:

```bash
compver scan . --show releases
```

Mostrar información del ciclo de vida y soporte:

```bash
compver scan . --show lifecycle
```

Mostrar vulnerabilidades conocidas:

```bash
compver scan . --show vulnerabilities
```

Mostrar fuentes y estado de la caché:

```bash
compver scan . --show sources
```

También pueden combinarse varias secciones separándolas por comas:

```bash
compver scan . \
  --show lifecycle,vulnerabilities,sources,releases
```

Las secciones disponibles son:

- `releases`: versiones seleccionadas y sus metadatos.
- `lifecycle`: estado de soporte y fecha de fin de soporte.
- `vulnerabilities`: vulnerabilidades conocidas y versiones que las corrigen.
- `sources`: fuentes consultadas, observaciones detectadas y estado de la
  caché.

> Breve recordatorio

Las opciones de `scan` pueden colocarse antes o después de la ruta:

Ejemplos:
```bash
compver scan --format json /ruta/proyecto
compver scan /ruta/proyecto --format json
compver scan /ruta/proyecto --tool terraform --limit 6
```

Compver reorganiza internamente los argumentos para procesarlos correctamente.  
La ruta continúa siendo única y puede combinarse con las opciones disponibles.  

#### Formato 

Añadido argumento `--format`: Nos permite modificar la salida de la CLI, ya sea por tabla o directamente en formato JSON, a su vez nos permite redireccionar dicho JSON a un archivo concreto que queramos con la redireccion de la propia shell.

#### Caché y conectividad

Compver utiliza una caché común para reducir consultas externas y permitir
análisis en proyectos con conectividad limitada.

Usar automáticamente la caché cuando sea válida:

```bash
compver scan .
```

Analizar únicamente con información almacenada en caché:

```bash
compver scan . --offline
```

Ignorar la caché y consultar de nuevo las fuentes disponibles:

```bash
compver scan . --refresh
```

`--offline` y `--refresh` no pueden utilizarse al mismo tiempo.


#### Opciones de personalizacion de la salida CLI

Por defecto la cli presenta colores y iconos ANSI cuando la salida es una terminal. Los colores no estaran cuando se use el parametro de `--format json`.

Se utiliza la ruta de `tmp/compver-fixture` como ejemplo didactico, al igual que usar `.`.

Si queremos quitar los colores:
```bash
compver scan /tmp/compver-fixture --no-color
```
Para mostrar todos los detalles disponibles:
```bash
compver scan . --verbose
```

Para mostrar únicamente los resúmenes:
```bash
compver scan . --summary
```

Para ocultar la salida normal:
```bash
compver scan . --quiet
```

Las opciones `--verbose`, `--summary` y `--quiet` no pueden utilizarse juntas.

#### Automatizacion y códigos de salida

Compver puede utilizarse en pipelines de CI/CD mediante `--fail-on`.

Fallar si se detectan vulnerabilidades:

```bash
compver scan . --fail-on vulnerable
```

Fallar si se detectan dependencias fuera de soporte:

```bash
compver scan . --fail-on end-of-life
```

También pueden combinarse varios criterios:

```bash
compver scan . --fail-on vulnerable,end-of-life
```

El alias `eol` equivale a `end-of-life`:

```bash
compver scan . --fail-on vulnerable,eol
```

| Código | Significado |
|---------|-------------|
| `0`| El análisi finalizó correctamente y no se activó ninguna politica. |
| `1`| Se produjo un error durante la ejecución |
| `2`| Se incumplió una condición definida mediante `--fail-on` |


Ejemplo para CI/CD con salida JSON:
```bash
compver scan . \
  --format json \
  --quiet \
  --fail-on vulnerable,end-of-life
```


#### Mostrar la versión
```bash
compver version
```

#### Mostrar la ayuda
```bash
compver help
compver -h
compver --help
```

### Ejemplo
```text
$ compver scan . --show lifecycle,vulnerabilities,sources

Ruta analizada: .
Dependencias encontradas: 2

=== GO ===
Nombre                                        Version         Tipo         Estado       Disponibles  Mostradas
------                                        -------         ----         ------       -----------  ---------
golang.org/x/text                             v0.14.0         directa      update-available 12           3
  Fuente del detector: go.mod
  Fuente de versiones: https://proxy.golang.org (red)
  Fuente de lifecycle: no disponible
  Fuente de seguridad: https://api.osv.dev/v1/query (red)
  Lifecycle: unknown
  Vulnerabilidades: ninguna conocida
  Estado: update-available
  Última publicada: v0.16.0
  Última compatible: v0.16.0

=== TERRAFORM ===
Nombre                                        Version         Tipo         Estado       Disponibles  Mostradas
------                                        -------         ----         ------       -----------  ---------
hashicorp/aws                                 v5.0.0          directa      current      20           3
  Fuente del detector: main.tf
  Fuente de versiones: https://registry.terraform.io (caché)
  Fuente de lifecycle: https://endoflife.date/api/terraform.json (red)
  Fuente de seguridad: https://api.osv.dev/v1/query (red)
  Lifecycle: active
  Vulnerabilidades: ninguna conocida
  Estado: current
  Última publicada: v5.0.0
  Última compatible: v5.0.0
```

#### Opciones disponibles

| Opción | Función |
|---|---|
| `--tool <nombre>` | Analiza únicamente una herramienta compatible, como `go` o `terraform`. |
| `--limit <número>` | Limita el número de versiones mostradas. |
| `--all-versions` | Muestra todas las versiones disponibles. |
| `--include-prereleases` | Incluye versiones `alpha`, `beta`, `rc` y otras pre-releases. |
| `--show <secciones>` | Muestra información adicional: `lifecycle`, `vulnerabilities`, `sources` y `releases`. |
| `--offline` | Utiliza únicamente datos almacenados en caché. |
| `--refresh` | Ignora la caché y consulta de nuevo las fuentes externas. |
| `--format <formato>` | Define el formato de salida: `table` o `json`. |
| `--no-color` | Desactiva los colores de la terminal. |
| `--verbose` | Activa todas las secciones detalladas del análisis. |
| `--summary` | Muestra únicamente el resumen general y los resúmenes por herramienta. |
| `--quiet` | Silencia la salida normal, manteniendo los errores y el código de salida. |
| `--fail-on <criterios>` | Devuelve un error si se cumplen determinadas condiciones. |

---

### Control sobre exclusiones de ***Compver***

A la hora de verificar versiones, es importante tambien contar con un archivo que nos permita colocar las extensiones, rutas o archivos que queramos que ***compver*** ignore a la hora de verificar sus versiones.

Para ello se ha añadido al motor la integracion de un `.compverignore` en el recorrido del motor.

Las exclusiones predeterminadas actuales ronda:

- `.git`
- `node_modules`
- `vendor`
- `.cache`
- `dist`
- `build`

Tambien contiene opciones para:

- Soporte para excluir directorios y archivos personalizados.
- Soporte para comentarios y lineas vacias.
- Validación de patrones inválidos.
- `.terraform/` no se excluya automáticamente para conservar la detección de providers instalados.

Para asignar estos cambios, se debe crear el archivo en la raiz del proyecto que se quiera analizar, manteniendo la misma metodologia que `.git`.

> Recordatorio: `.compverignore` afecta únicamente a Compver, no modifica el comportamiento de Git ni sustituye al `.gitignore`.

### Lenguaje usado (actual)

Se usara el lenguaje ***Go*** debido a que es facil de mantener, rapido y lo suficientemente estable.

Tambien debido a la amplia magnitud de dispositivos que usan dicho Lenguaje, todo ello suma.



---


## Metas especificas principales actuales:

 ✅ Creacion de un motor general.    
 ✅ Visualizacion de versiones.  
    Probar funcionamiento con adaptadores de prueba.  
    Compatibilidad con herramientas DevOps.  
    Visualizar su funcionamiento.  
    Solucion de errores o Bugs.  
    Crear binarios GO en pre-releases.  
    Verificacion.  
    Lanzar 1.0 en rama main.  

