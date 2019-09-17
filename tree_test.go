package main

import (
	"fmt"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	"strconv"
	"testing"

)


var inorderfacit = [9]int64{5,6,7,9,10,11,12,15,20}


func getOrderTestData(){

}


func getAGPSObject(lat float64,long float64,name string) *GPSLocation{
	obj := &GPSLocation{
		Location:Locationdata{
			Latitude:lat,
			Longitude:long,
			Accuracy:0,
			Zindex:0,
		},
		Gpsobject:0,
		Uuid:name,
		Timestamp:0,
	}
	obj.Location.Zindex = GetZorderIndex(lat,long)
	return obj
}

func getDataZOrderTest()[] *GPSLocation{

	// below equator
	var locations [] *GPSLocation


	locations = append(locations, getAGPSObject(37.391547,-122.034613,"Pastoria Avenue 2"))
	locations = append(locations, getAGPSObject(-32.576457,-66.228494,"San Loius Argentina"))
	locations = append(locations, getAGPSObject(37.390227,-122.034238,"Pastoria Avenue 1"))
	locations = append(locations, getAGPSObject(37.378128,-122.038240,"Pastoria Avenue 3"))
	locations = append(locations, getAGPSObject(37.369688,-122.041008,"Hollen beck av"))
	locations = append(locations, getAGPSObject(37.368328,-122.039109,"Holthouse - terrace"))
	locations = append(locations, getAGPSObject(37.367189,-122.039141,"Sunnyvale - west"))
	locations = append(locations, getAGPSObject(37.389324,-122.029620,"KFC SV"))
	locations = append(locations, getAGPSObject(37.388348, -122.030375,"Starbucks SV"))
	locations = append(locations, getAGPSObject(37.426914, -122.097788,"BMW"))
	locations = append(locations, getAGPSObject(37.426506, -122.0977128,"BMW 2"))

	return locations
}

func swap(location1 *GPSLocation,location2 *GPSLocation){
	temp := *location1
	*location1 = *location2
	*location2 = temp
}

func sort(locations [] *GPSLocation){

	for i := 1; i < len(locations);i++{
		j := i
		for j>0  && (byGPSIndexation(*locations[j-1],*locations[j]) == 1) {
            swap(locations[j],locations[j-1])
			j--;
		}
	}

}

func printlocations(locations []* GPSLocation){
	for _,e := range locations{
		fmt.Println(e.Uuid)
	}
	fmt.Println("")
}

func getLatLong(loc * GPSLocation) (float64,float64){
	return loc.Location.Latitude,loc.Location.Longitude
}

func TestZorder(t *testing.T){

	locations := getDataZOrderTest()
	sort(locations)
	printlocations(locations)
	lat1,long1 := getLatLong(locations[0])
	var distances [] float64

	for i := 1; i < len(locations);i++{
		lat2,long2 := getLatLong(locations[i])
		distances = append(distances,GetApproxDistance2(lat1, long1,lat2,long2))
	}

	distanceok := true
	index := 1
	shortestdistance := distances[0]
	for distanceok && index < len(distances){
		distanceok = shortestdistance < distances[index]
		index++
	}
   if (distanceok == false) {
	   t.Errorf("zorder did not keep distance")
   }
}

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
	PrintPreOrder(t.Right)

}

func PrintInOrder(t *rbt.Node){

	if (t  == nil){
		return
	}

	PrintInOrder(t.Left)
	fmt.Println(t.Key.(GPSLocation).Location.Zindex)
	PrintInOrder(t.Right)
}

func PrintPostOrder(t *rbt.Node) {
	if (t == nil) {
		return
	}

	PrintPostOrder(t.Left)
	PrintPostOrder(t.Right)
	fmt.Println(t.Key.(GPSLocation).Location.Zindex)
}

func TestPrintTree(t *testing.T){
	location_tree := TreeExtended{GetTree()}

	fmt.Println("************** PRE ORDER ***************")
	PrintPreOrder(location_tree.Root)
	fmt.Println("************** POST ORDER ***************")
	PrintPostOrder(location_tree.Root)
	fmt.Println("************** IN ORDER ***************")
	PrintInOrder(location_tree.Root)
	fmt.Println("************** END ORDER ***************")
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