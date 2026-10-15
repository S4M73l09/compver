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
| `--tool <nombre>` | Analiza únicamente una herramienta compatible. |
| `--limit <número>` | Limita el número de versiones mostradas. |
| `--all-versions` | Muestra todas las versiones disponibles. |
| `--include-prereleases` | Incluye alpha, beta, rc y otras pre-releases. |
| `--show <secciones>` | Muestra información adicional. |
| `--offline` | Utiliza únicamente datos almacenados en caché. |
| `--refresh` | Ignora la caché y consulta de nuevo las fuentes externas. |

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

