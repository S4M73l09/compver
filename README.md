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

***compver*** usa una serie de comandos para controlar de mejor manera el funcionamiento de este:

Podemos invocar ***compver*** con el comando para analizar el directorio actual:

#### Analizar un proyecto
```bash
compver scan [ruta]
```

Si no se indica una ruta, se analiza el directorio actual:

```bash
compver scan
```

#### Atajo de análisis

Estas formas son equivalente:
```bash
compver
compver .
compver /ruta/al/proyecto
compver scan /ruta/al/proyecto
compver scan --tool terraform --limit 6 /ruta/proyecto
compver scan --tool terraform --limit 3 --include-pre-releases
```

El argumento `--tool` se puede juntar con cualquier comando mostrado junto con los demas argumentos. En este caso se usa de ejemplo `Terraform`.

Todas ejecutan un análisis básico de la ruta indicada.

##### Argumentos añadidos al scan

Se añadieron estos argumentos para el comando anterior.

```bash
compver scan --limit <Numero> .
compver scan --include-pre-releases .
compver scan --all-versions .
```
Esto nos permite mejorar la auditoria de versiones de alguna herramienta o carpeta.

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
$ compver scan .

Ruta analizada: .
Dependencias encontradas: 2

Nombre                                        Versión         Tipo         Estado
------                                        -------         ----         ------
github.com/example/library                    v1.2.3          directa      stable
golang.org/x/text                             v0.14.0         indirecta    beta
```

#### Opciones previstas

Está previsto añadir opciones para personalizar el análisis:

```bash
compver scan /ruta --include-prereleases
compver scan /ruta --only-outdated
compver scan /ruta --format-json
```

Estas opciones todavía están en desarrollo.

---


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


