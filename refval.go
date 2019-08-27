package main

import (
	"fmt"
)

type mupp struct{
	x int
	y int
}

var global = &mupp{2,3}


func ChangeStruct(a *mupp,b *mupp){
    //copier.Copy(a,global)

	a.x = 10
	a.y = 15


    *b = *global

	fmt.Println(b.x,b.y)
	fmt.Println(global.x,global.y)
}



func CallChange(){

	var first = mupp{100,100}

	var second mupp


	ChangeStruct(&first,&second)

	fmt.Println(first.x,first.y)
	fmt.Println(second.x,second.y)
}