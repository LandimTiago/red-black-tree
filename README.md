# Red-Black Tree

Estudos direcionados ao funcionamento de order books em uma aplicação golang

## 📌 Passo 1: Estruturando os Nós da Red-Black Tree

Cada nó precisa armazenar:

- Valor (price, por exemplo)
- Cor (vermelho ou preto)
- Filhos (esquerdo e direito)
- Pai (para facilitar rebalanceamento)
- Vamos definir essas propriedades no struct Node:

### 🌱 Criando o Arquivo red_black_tree.go

Aqui está a primeira parte do código:

```go
    package main

    import "fmt"

    // Definição das cores dos nós
    const (
        RED   = true  // Nó vermelho 🔴
        BLACK = false // Nó preto ⚫
    )

    // Estrutura do nó da Red-Black Tree
    type Node struct {
        price  int     // Valor do nó (pode ser um preço no order book)
        color  bool    // Cor do nó (RED ou BLACK)
        parent *Node   // Ponteiro para o nó pai
        left   *Node   // Ponteiro para o filho esquerdo
        right  *Node   // Ponteiro para o filho direito
    }

    // Definição da árvore Red-Black
    type RedBlackTree struct {
        root *Node // Nó raiz da árvore
    }
```

### 📝 Explicação

1. Criamos constantes RED e BLACK para representar as cores.
2. Criamos o struct Node com:

- price: O valor armazenado (no caso de um order book, poderia ser um preço de uma ordem).
- color: Cor do nó (vermelho 🔴 ou preto ⚫).
- parent, left, right: Ponteiros para estruturar a árvore.

3. Criamos o struct RedBlackTree que contém a root (raiz).

## 📌 Passo 2: Criando um Nó na Árvore

Agora, precisamos de uma função que cria um novo nó:

```go
// Função para criar um novo nó
    func newNode(price int) *Node {
        return &Node{
            price: price,
            color: RED, // Todo novo nó começa como vermelho 🔴
            left:  nil,
            right: nil,
            parent: nil,
        }
    }

```

### 📝 Explicação

1. A função newNode(price int) cria um nó vermelho 🔴.
2. O nó começa sem filhos (left e right são nil).
3. O pai (parent) também começa como nil.

## 📌 Passo 3: Inserindo um Nó na Red-Black Tree

A inserção segue as regras:

1. A dicionamos como em uma árvore binária de busca normal.
2. Se o pai for vermelho 🔴, fazemos ajustes (recolorimento e rotação).

Vamos criar a função de inserção:

```go
// Função para inserir um nó na árvore
    func (tree *RedBlackTree) Insert(price int) {
        newNode := newNode(price)

        if tree.root == nil {
            // Se a árvore estiver vazia, o novo nó se torna a raiz e deve ser preto ⚫
            newNode.color = BLACK
            tree.root = newNode
            return
        }

        var parent *Node
        current := tree.root

        // Encontrando a posição correta para inserir
        for current != nil {
            parent = current
            if newNode.price < current.price {
                current = current.left
            } else {
                current = current.right
            }
        }

        // Definindo o pai do novo nó
        newNode.parent = parent

        // Inserindo o nó no local correto
        if newNode.price < parent.price {
            parent.left = newNode
        } else {
            parent.right = newNode
        }

        // Corrigir a árvore para manter as regras Red-Black
        tree.fixInsert(newNode)
    }

```

### 📝 Explicação

1. Criamos um novo nó vermelho 🔴.
2. Se a árvore estiver vazia, o novo nó vira a raiz e precisa ser preto ⚫.
3. Buscamos o local correto na árvore para inserção.
4. Inserimos o novo nó como filho do nó encontrado.
5. Chamamos fixInsert(newNode) para corrigir possíveis violações da árvore Red-Black.

## 📌 Passo 4: Corrigindo a Árvore Após Inserção

Aqui começa a magia da Red-Black Tree! Precisamos de:

1. Recolorir nós se houver dois vermelhos 🔴 seguidos.
2. Rotacionar se necessário para manter o balanceamento.

```go
    // Função para corrigir a arvore após a inserção
    func (tree *RedBlackTree) fixInsert(node *Node) {
        for node.parent != nil && node.parent.color == RED {
            grandparent := node.parent.parent

            // O pai está a esquerda do avô
            if node.parent == grandparent.left {
                uncle := grandparent.right

                // Caso 1: O tio também é vermelho 🔴 → Recolorimos
                if uncle != nil && uncle.color == RED {
                    node.parent.color = BLACK
                    uncle.color = BLACK
                    grandparent.color = RED

                    node = grandparent // Continuamos verificando para cima
                } else {
                    // Caso 2: O nó é um filho à direita → Rotação para a esquerda
                    if node == node.parent.right {
                        node = node.parent
                        tree.leftRotate(node)
                    }

                    // Caso 3: Rotação para a direita e recoloração
                    node.parent.color = BLACK
                    grandparent.color = RED
                    tree.rightRotate(grandparent)
                }
            } else {
                // O pai está à direita do avô (espelho do caso anterior)
                uncle := grandparent.left

                if uncle != nil && uncle.color == RED {
                    node.parent.color = BLACK
                    uncle.color = BLACK
                    grandparent.color = RED
                    node = grandparent
                } else {
                    if node == node.parent.left {
                        node = node.parent
                        tree.rightRotate(node)
                    }

                    node.parent.color = BLACK
                    grandparent.color = RED
                    tree.leftRotate(grandparent)
                }
            }

        }

        // A raiz sempre deve ser preta ⚫
        tree.root.color = BLACK
    }
```

### 📝 Explicação

1. Se o pai do novo nó for vermelho 🔴, verificamos o tio (irmão do pai).
2. Caso 1: Se o tio for vermelho 🔴 → Recolorimos os nós para manter o balanceamento.
3. Caso 2: Se o tio for preto ⚫ e o novo nó for um filho direito → Fazemos rotação esquerda.
4. Caso 3: Se o tio for preto ⚫ e o novo nó for um filho esquerdo → Rotação direita e recoloração.
5. Garantimos que a raiz sempre seja preta ⚫.

## 📌 Passo 5: Implementando as Rotações

AS rotações são usadas para manter a árvore balanceada

```go
   // Rotação para a esquerda
    func (tree *RedBlackTree) leftRotate(node *Node) {
        rightChild := node.right
        if rightChild == nil {
            return
        }

        node.right = rightChild.left
        if rightChild.left != nil {
            rightChild.left.parent = node
        }

        rightChild.parent = node.parent

        if node.parent == nil {
            tree.root = rightChild
        } else if node == node.parent.left {
            node.parent.left = rightChild
        } else {
            node.parent.right = rightChild
        }

        rightChild.left = node
        node.parent = rightChild
    }

    // Rotação para a direita
    func (tree *RedBlackTree) rightRotate(node *Node) {
        leftChild := node.left
        node.left = leftChild.right

        if leftChild.right != nil {
            leftChild.right.parent = node
        }

        leftChild.parent = node.parent

        if node.parent == nil {
            tree.root = leftChild
        } else if node == node.parent.right {
            node.parent.right = leftChild
        } else {
            node.parent.left = leftChild
        }

        leftChild.right = node
        node.parent = leftChild
    }
```

## 📌 Criando um Teste Simples

Crie um novo arquivo main.go e adicione o seguinte código

```go
package main

import "fmt"

// Função auxiliar para imprimir a árvore em ordem
func inOrderTraversal(node *Node) {
    if node != nil {
        inOrderTraversal(node.left)
        fmt.Printf("%d (%s) ", node.price, colorToString(node.color))
        inOrderTraversal(node.right)
    }
}

// Converte a cor do nó para string (apenas para debug)
func colorToString(color bool) string {
    // if color == RED {
    //     return "🔴"
    // }

    // versão simplificada
    if color {
		return "🔴"
	}
    return "⚫"
}

func main() {
    tree := &RedBlackTree{}

    // Inserindo valores na árvore
    values := []int{10, 20, 30, 15, 25, 5, 1}
    for _, v := range values {
        tree.Insert(v)
    }

    // Exibindo a árvore em ordem
    fmt.Println("Árvore em ordem (in-order traversal):")
    inOrderTraversal(tree.root)
    fmt.Println("\n")

    // Exibindo a raiz da árvore
    fmt.Printf("Raiz: %d (%s)\n", tree.root.price, colorToString(tree.root.color))
}

```

## 📌 Como esse código testa a Red-Black Tree?

1. Criamos uma árvore vazia.
2. Inserimos os valores {10, 20, 30, 15, 25, 5, 1}.
3. Exibimos os valores ordenados com a função inOrderTraversal.
4. Exibimos a raiz da árvore e sua cor.

## 📌 Executando o Teste

Para rodar o código, basta usar:

```sh
    go run main.go
```

Se tudo estiver correto, o programa exibirá os valores em ordem crescente, com as cores dos nós indicando um balanceamento correto da árvore.
