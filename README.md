# Proyecto: Compver ***Rama: Pre-release***

En esta rama se propondran cambios en la estructura del funcionamiento de proyecto a la vez que se buscara el mantenimiento de este añadiendo nuevas funciones, cambios, solucion de bugs o errores, generacion de cambios y modificaciones varias.

Dicha rama contiene contenido en constante cambio, es libre del usuario descargar y usar las funciones tanto de prueba como normales del proyecto, pero igualmente puede contener errores y/o cambios futuros inesperados.

Dicho esto, esta bajo tu propio riesgo.


## Compver

Esta aplicacion ha sido desarrollada para implementar una solucion alternativa a la hora de usar innumerables herramientas que cuentan con versiones - pre-releases propias. Todas ellas contienen su propio codigo y es normal a la larga darse cuenta del caos de mantener versiones estables directamente buscar versiones mas actuales sin necesidad de ir a lo ultimo.

Por ello se ha planteado esta solucion: ***Compver*** es una solucion de codigo abierto, es un comparador de versiones capaz de analizar la ruta que el usuario necesite o el filesystem del repositorio concreto donde sea necesario.

### Su funcion

***Compver*** compara las versiones usadas en una infraestructura concreta con las versiones mas estables y actuales de la misma herramienta, tambien da una lista sobre las pre-release que tienen las herramientas que este ha analizado. 

EL proyecto utiliza un motor generico, esto permite implementar adaptadores cuando se quiera añadir soporte a una nueva tecnologia, permitiendonos no necesitar ni modificar el motor.

En el futuro se buscara implementar y mostrar adaptadores de python para que cada usuario tenga la libertad de añadir soporte a su stack tecnologico concreto.

Su funcion ira evolucionando en funcion a como avance el proyecto.

### Lenguaje usado (actual)

Se usara el lenguaje ***Go*** debido a que es facil de mantener, rapido y lo suficientemente estable.

Tambien debido a la amplia magnitud de dispositivos que usan dicho Lenguaje, todo ello suma.



---


## Metas especificas principales actuales:

    Creacion de un motor general  
    Visualizacion de versiones  
    Compatibilidad con herramientas DevOps  
    Visualizar su funcionamiento  
    Solucion de errores o Bugs
    Verificacion
    Lanzar 1.0 en rama main


