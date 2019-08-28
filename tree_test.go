package main

import (
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"strconv"
	"testing"

)


var inorderfacit = [9]int64{5,6,7,9,10,11,12,15,20}


func getTestData()[] GPSLocation{
	datapoints := [9]uint64{10,5,7,6,15,12,11,9,20}
	var list  [] GPSLocation

	for i, e := range datapoints{
		var gps = &GPSLocation{
			Location: Locationdata{
				Latitude:  37.387401,
				Longitude: -122.035179,
				Accuracy:  1,
			},
			Gpsobject: 0,
			Uuid:strconv.FormatInt(int64(i), 10),
			Timestamp: 1,
		}
		gps.Location.Zindex = e
		list = append(list, *gps)
	}

    return list
}

func getDataTree() *TreeExtended{
    var tree = TreeExtended{GetTree()}

	data := getTestData()
	for _,e := range data {
		tree.AddGPSPosition(e)
	}
	return &tree
}


func testTreeData()*TreeExtended{
	location_tree := TreeExtended{GetTree()}
	data := getTestData()

	for _,e := range data{
		location_tree.AddGPSPosition(e)
	}
	return &location_tree
}

func TestTreeInsertion(t *testing.T){
	location_tree := TreeExtended{GetTree()}
	data := getTestData()

	for i,e := range data{
		neighbours,err := location_tree.AddGPSPosition(e)
		if err == nil {
			fmt.Print("index  [", i, "] = ",e.Location.Zindex," has succ/prec = ")
			for _, gpso := range neighbours {
				fmt.Print( "{",gpso.Location.Zindex, "}")
			}
			fmt.Println(" ")
		}else{
			t.Error("could not insert GPS position")
		}

	}

}


func PrintPreOrder(t *rbt.Node){

	if (t  == nil){
		return
	}
	fmt.Println(t.Key.(GPSLocation).Location.Zindex)
	PrintPreOrder(t.Left)
	//fmt.Println(t.Key.(GPSLocation).Location.Zindex)
	PrintPreOrder(t.Right)

}

func TestPrintTree(t *testing.T){
	location_tree := TreeExtended{GetTree()}

	PrintPreOrder(location_tree.Root)
	//PrintSuccOrder(location_tree)
}

func (t *TreeExtended)isValidBST(n *rbt.Node) bool{

	if (n == nil) {
		return true
	}

	if n.Left != nil && t.Comparator(n.Left.Key,n.Key) > 0{
		return false
	}

	if n.Right != nil && t.Comparator(n.Right.Key,n.Key) < 0{
		return false
	}

	t.isValidBST(n.Left)
    t.isValidBST(n.Right)

	return true
}

func TestValidBTree(t *testing.T){
	location_tree := testTreeData()

	if location_tree.isValidBST(location_tree.Root) ==  false {
		t.Errorf(" tree is nota valid binary tree %t ", false)
	}

}

func TestPreSuc(t *testing.T){
	location_tree := testTreeData()
	testdata := getTestData()

	pre := rbt.Node{}
	suc := rbt.Node{}
	location_tree.FindPreSuc(location_tree.Root,testdata[0],&pre,&suc)

	if (pre.Key.(GPSLocation).Location.Zindex != 9){
		t.Errorf(" expected 9 got %d",pre.Key.(GPSLocation).Location.Zindex)
	}

	if (suc.Key.(GPSLocation).Location.Zindex != 11){
		t.Errorf(" expected 11 got %d",suc.Key.(GPSLocation).Location.Zindex)
	}

	location_tree.FindPreSuc(location_tree.Root,testdata[4],&pre,&suc)

	if (pre.Key.(GPSLocation).Location.Zindex != 12){
		t.Errorf(" expected 12 got %d",pre.Key.(GPSLocation).Location.Zindex)
	}

	if (suc.Key.(GPSLocation).Location.Zindex != 20){
		t.Errorf(" expected 20 got %d",suc.Key.(GPSLocation).Location.Zindex)
	}

}


func inOrder(t *rbt.Node,counter *int ,val *bool){

	if (t  == nil){
		return
	}
	inOrder(t.Left,counter,val)

	key := t.Key.(GPSLocation).Location.Zindex
	if (inorderfacit[*counter] != int64(key)){
		*val = false
	}

	*counter = *counter + 1
	inOrder(t.Right,counter,val)

	return
}

func TestInorder(t *testing.T){
	location_tree := testTreeData()
	counter := 0
	inorder := true
	inOrder(location_tree.Root,&counter,&inorder)

    if (inorder == false){
    	t.Errorf(" expected true got %t ", inorder)
	}

	inorder = true
	counter = 0
	inorderfacit[7] = 14
	inOrder(location_tree.Root,&counter,&inorder)
	if (inorder == true){
		t.Errorf(" expected false got %t ", inorder)
	}
}

func TestNearbyNeighbours(t *testing.T) {

}