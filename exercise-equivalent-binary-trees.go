
package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	innerWalk(t, ch)
	close(ch)
}

func innerWalk(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		innerWalk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		innerWalk(t.Right, ch)
	}
}
	

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	for val1 := range ch1 {
		val2, ok := <-ch2
		if !ok || val1 != val2 {
			return false
		}
	}
	if _, ok := <- ch2; ok {
		return false
	}
	return true
}

func main() {
	fmt.Println(true == Same(tree.New(1), tree.New(1)))
	fmt.Println(false == Same(tree.New(1), tree.New(2)))
	fmt.Println(true == Same( &tree.Tree{nil,10,nil}, &tree.Tree{nil,10,nil}))
	fmt.Println(false == Same( &tree.Tree{nil,9,nil}, &tree.Tree{nil,10,nil}))
	fmt.Println(false == Same(tree.New(1), &tree.Tree{nil,10,nil}))
	fmt.Println(false == Same( &tree.Tree{nil,10,nil}, tree.New(1)))
}	


