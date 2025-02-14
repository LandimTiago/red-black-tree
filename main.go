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
				if node != nil && node == node.parent.left {
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

// Remove um nó com o valor especificado
func (tree *RedBlackTree) Delete(price int) {
	// Localizar o nó a ser removido
	nodeToDelete := tree.search(tree.root, price)

	if nodeToDelete == nil {
		// Valor não encontrado
		return
	}

	tree.deleteNode(nodeToDelete)
}

// Busca por um valor na árvore ( auxiliar para a remoção )
func (tree *RedBlackTree) search(node *Node, price int) *Node {

	if node == nil || node.price == price {
		return node
	}

	if price < node.price {
		return tree.search(node.left, price)
	}

	return tree.search(node.right, price)
}

// Função de remoção de um nó da árvore
func (tree *RedBlackTree) deleteNode(node *Node) {
	var child, replacement *Node
	originalColor := node.color

	if node.left == nil {
		// Caso 1: Sem filho à esquerda (inclui nó folha)
		child = node.right
		tree.transplant(node, node.right)
	} else if node.right == nil {
		// Caso 2: Sem filho à direita
		child = node.left
		tree.transplant(node, node.left)
	} else {
		// Caso 3: Nó com dois filhos → encontrar sucessor
		replacement = tree.minimum(node.right)
		originalColor = replacement.color
		child = replacement.right

		if replacement.parent == node {
			if child != nil {
				child.parent = replacement
			}
		} else {
			tree.transplant(replacement, replacement.right)
			replacement.right = node.right
			replacement.right.parent = replacement
		}

		tree.transplant(node, replacement)
		replacement.left = node.left
		replacement.left.parent = replacement
		replacement.color = node.color
	}

	// Corrigir as propriedades Red-Black após remoção de um nó preto
	if originalColor == BLACK {
		tree.fixDelete(child)
	}
}

// Substitui um nó por outro (auxiliar para a remoção)
func (tree *RedBlackTree) transplant(u, v *Node) {
	if u.parent == nil {
		tree.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}

	if v != nil {
		v.parent = u.parent
	}
}

// Corrige a árvore após a remoção
func (tree *RedBlackTree) fixDelete(node *Node) {
	for node != tree.root && (node == nil || node.color == BLACK) {
		if node == node.parent.left {
			sibling := node.parent.right

			// Caso 1: Irmão é vermelho
			if sibling != nil && sibling.color == RED {
				sibling.color = BLACK
				node.parent.color = RED
				tree.leftRotate(node.parent)
				sibling = node.parent.right
			}

			// Caso 2: Ambos os filhos do irmão são pretos
			if (sibling.left == nil || sibling.left.color == BLACK) &&
				(sibling.right == nil || sibling.right.color == BLACK) {
				if sibling != nil {
					sibling.color = RED
				}
				node = node.parent
			} else {
				// Caso 3: Irmão tem filho vermelho à esquerda
				if sibling.right == nil || sibling.right.color == BLACK {
					if sibling.left != nil {
						sibling.left.color = BLACK
					}
					sibling.color = RED
					tree.rightRotate(sibling)
					sibling = node.parent.right
				}

				// Caso 4: Irmão tem filho vermelho à direita
				if sibling != nil {
					sibling.color = node.parent.color
				}
				node.parent.color = BLACK
				if sibling.right != nil {
					sibling.right.color = BLACK
				}
				tree.leftRotate(node.parent)
				node = tree.root
			}
		} else {
			// Espelho dos casos acima para o irmão à esquerda
			sibling := node.parent.left

			if sibling != nil && sibling.color == RED {
				sibling.color = BLACK
				node.parent.color = RED
				tree.rightRotate(node.parent)
				sibling = node.parent.left
			}

			if (sibling.right == nil || sibling.right.color == BLACK) &&
				(sibling.left == nil || sibling.left.color == BLACK) {
				if sibling != nil {
					sibling.color = RED
				}
				node = node.parent
			} else {
				if sibling.left == nil || sibling.left.color == BLACK {
					if sibling.right != nil {
						sibling.right.color = BLACK
					}
					sibling.color = RED
					tree.leftRotate(sibling)
					sibling = node.parent.left
				}

				if sibling != nil {
					sibling.color = node.parent.color
				}
				node.parent.color = BLACK
				if sibling.left != nil {
					sibling.left.color = BLACK
				}
				tree.rightRotate(node.parent)
				node = tree.root
			}
		}
	}

	if node != nil {
		node.color = BLACK
	}
}

// Encontra o menor nó de uma subárvore
func (tree *RedBlackTree) minimum(node *Node) *Node {
	current := node
	for current.left != nil {
		current = current.left
	}
	return current
}

func main() {
	tree := &RedBlackTree{}

	// Inserir alguns valores
	values := []int{20, 15, 25, 10, 18, 22, 30}
	for _, v := range values {
		tree.Insert(v)
	}

	fmt.Println("Árvore em ordem (antes da remoção):")
	inOrderTraversal(tree.root)
	fmt.Println()

	// Remover um valor
	tree.Delete(15)
	fmt.Println("Árvore em ordem (após remover 15):")
	inOrderTraversal(tree.root)
	fmt.Println()
}
