package main

import (
	"fmt"
)

// Definição das cores dos nós
const (
	RED   = true
	BLACK = false
)

// Estrutura do nó da Red-Black Tree
type Node struct {
	price  int   // Valor do nó (preço no order book)
	color  bool  // Cor do nó
	parent *Node // Nó pai
	left   *Node // Nó filho esquerdo
	right  *Node // Nó filho direito
}

// Definição da Red-Black Tree
type RedBlackTree struct {
	root *Node // Raiz da árvore
}

// Função para criar um novo nó, começa sem filhos e o pai também é nil
func newNode(price int) *Node {
	return &Node{
		price:  price,
		color:  RED,
		parent: nil,
		left:   nil,
		right:  nil,
	}
}

// Função para inserir um novo nó na árvore
func (tree *RedBlackTree) Insert(price int) {
	newNode := newNode(price)

	if tree.root == nil {
		// Se a arvore estiver vazia, o novo nó se torna a raiz e deve ser preto ⚫
		newNode.color = BLACK
		tree.root = newNode
		return
	}

	var parent *Node
	current := tree.root

	// Buscando a posição correta para a inserção
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

	if newNode.price < parent.price {
		parent.left = newNode
	} else {
		parent.right = newNode
	}

	// Corrigir a arvore para manter as regras Red-Black
	tree.fixInsert(newNode)
}

// Função para corrigir a arvore após a inserção
func (tree *RedBlackTree) fixInsert(node *Node) {
	for node.parent != nil && node.parent.color {
		grandparent := node.parent.parent

		// O pai está a esquerda do avô
		if node.parent == grandparent.left {
			uncle := grandparent.right

			// Caso 1: O tio também é vermelho 🔴 → Recolorimos
			if uncle != nil && uncle.color {
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

			if uncle != nil && uncle.color {
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
