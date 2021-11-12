# [Débogueur visuel SIMD](http://jpmiceli.com.ar/)

Langues prises en charge :

[🇪🇸](./../README.md)[🇬🇧](./README.en.md)[🇵🇹](./README.pt.md)[🇫🇷](./README.fr.md)[🇮🇳](./README.hi.md)

## Introduction:

Il s'agit d'un environnement de développement spécialement conçu pour travailler avec des instructions SIMD.

Il a un format similaire à celui du cahier jupyter, séparé en cellules dans lesquelles le code en cours de développement sera placé, par exemple le noyau d'un cycle dans un filtre d'image, et entre eux, il peut être affiché comme Dans un convivial manière, la valeur des enregistrements XMM, dans le format souhaité, dont vous souhaitez garder une trace.

## Mode d'utilisation:

### Rubrique de données :

La première cellule correspond toujours à la cellule de données et est signalée par l'étiquette`section .data`, les données sont définies exactement comme nous le ferions en assembleur.

Il est très important de noter que seule cette cellule peut être utilisée comme cellule de données, et de plus, cela doit en être la seule fonction. C'est-à-dire que vous ne pouvez pas y ajouter de code.

Cette cellule ne peut pas être supprimée.

### Section de texte :

Le début de cette section est marqué par la balise`section .text`, c'est-à-dire que tout le texte qui se trouve dans les cellules sous cette étiquette sera considéré comme du code.

La suppression de toutes les cellules de texte n'est pas autorisée, il y aura donc toujours au moins une cellule dans la section de texte.

### Coder en C :

Pour des raisons de sécurité, cet environnement n'autorise aucun appel aux fonctions libc. Pour aujourd'hui vous n'aurez qu'à programmer en assembleur.

### Appels système :

Cet environnement n'accepte aucun appel système. Si vous essayez d'exécuter l'un d'entre eux, le programme se terminera.

### Ajouter/Supprimer des cellules :

Toutes les cellules vous permettent d'ajouter des cellules, à la fois au-dessus et en dessous d'elles. Pour y parvenir, vous devez appuyer sur le bouton`+ Code`trouvé entre les cellules. Il est invisible jusqu'à ce que la souris se trouve à cette hauteur.
Pour supprimer une cellule, il suffit d'appuyer sur le bouton en forme de poubelle que chaque cellule possède en haut à droite.

### Imprimer les enregistrements :

Les enregistrements XMM acceptent les formats d'impression suivants :

-   16 registres entiers de 8 bits`.v16_int8`
-   8 registres entiers 16 bits`.v8_int16`
-   4 registres entiers 32 bits`.v4_int32`
-   2 registres d'entiers 64 bits`.v2_int64`
-   4 registres de nombres flottants simple précision (32 bits)`.v4_float`
-   2 registres de nombres flottants double précision (64 bits)`.v2_double`

À leur tour, les formats d'entiers peuvent être imprimés dans les bases de nombres suivantes :

-   Signé Base 10`/d`
-   Base 10 non signée`/u`
-   Base 16`/x`
-   Base 2 en complément A2`/t`

Pour demander qu'un enregistrement soit imprimé, le même format utilisé dans GDB doit être utilisé :

`;print<NB> xmm<PR><BF>`

Où`NB`est le format de la base numérique et est facultatif, par défaut c'est`/d`. Mais cette valeur par défaut est mise à jour lors de l'utilisation de l'instruction`;print`.

`PR`est l'enregistrement à imprimer.

`BF`est le format de bit et est également facultatif et sa valeur par défaut est`.v16_int8`. Il est mis à jour lors de l'utilisation de l'instruction`;print`. Vous pouvez définir ce format pour tous les enregistrements en même temps, en omettant la valeur de l'enregistrement à imprimer.

`;print`est analogue à`;p`

Quelques exemples de ceci :

`;p xmm.v4_int16`: J'ai défini la valeur par défaut des bits pour tous les registres à imprimer en entiers de 16 bits.

`;p/x xmm1.v2_int64`: j'imprime le registre xmm1 sous forme d'entiers 64 bits en hexadécimal

`;print xmm2.v2_double`: J'imprime le registre xmm2 sous forme de nombres flottants de 64 bits.

Les enregistrements sont imprimés dans 2 cas.

1) Si la valeur du registre a été modifiée dans l'exécution de la cellule.
2) Si l'utilisateur demande que l'enregistrement soit imprimé.

Si l'impression se fait par le premier cas, les valeurs par défaut que le registre a à ce moment-là seront utilisées.

Une façon de « commenter »`print`est en écrivant n'importe quel caractère au milieu. Par exemple:

`;=p/x xmm0.v2_int64`

### Copier le code :

Le bouton "copier dans le presse-papiers" fait exactement ce que vous pensez. Copiez tout le code qui se trouve dans les cellules dans le presse-papiers afin de pouvoir le coller facilement dans l'assembleur. Un détail de ceci est que s'il y a une cellule que vous ne voulez pas copier, vous pouvez ajouter un :`;nope`dans cette cellule et il n'apparaîtra pas dans la copie.

### Code propre :

Comme effacer le code de toutes les cellules à la fois est un peu gênant, vous pouvez appuyer sur le bouton`Clean Code`pour effectuer cette action.**Cela ne peut pas être annulé**, pour laquelle l'utilisation de ce bouton doit être confirmée au moyen d'une pop-up.

### Masquer les enregistrements :

Pour masquer les enregistrements que vous ne souhaitez pas imprimer il existe 2 alternatives : la première consiste à utiliser la commande`hide`qui est expliqué dans la section des commandes, cette commande nous permet de masquer les enregistrements individuellement. Si ce que vous voulez faire est de masquer tous les enregistrements provenant d'une cellule spécifique, vous pouvez appuyer sur le bouton en forme d'œil qui se trouve en haut à droite de la cellule.

### Raccourcis clavier:

-   `Ctrl+Enter`: Exécuter le code.
-   `Ctrl+ArrowDown`: insère une cellule sous la cellule actuelle. Si vous ne vous tenez sur aucune cellule, elle sera insérée à la fin.
-   `Ctrl+ArrowUp`: insère une cellule au-dessus de la cellule actuelle. Si vous ne vous tenez sur aucune cellule, elle sera insérée au début.
-   `Ctrl+Alt+D`: Supprimer la cellule actuelle.**Cette action ne peut pas être annulée.**
-   `Alt+ArrowDown`: déplace le curseur d'une cellule en dessous de la cellule actuelle.
-   `Alt+ArrowUp`: Déplace le curseur d'une cellule au-dessus de la cellule actuelle.
