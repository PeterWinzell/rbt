package main

import (
	"container/list"
	lls "github.com/emirpasic/gods/stacks/linkedliststack"
	rbt "github.com/emirpasic/gods/trees/redblacktree"
	log "github.com/sirupsen/logrus"
	"sync"
)

const (

	MaxDetections = 2
	Expirationtime = 5 // throw away any entries older that this (seconds)
	Timedepth = 5 // only consider queue neighbours with in this (seconds)

	Criticaldistance = 200 // The distance to a bike/car where we issue a warning (meters)
	GARBAGESIZE = 5
)


/*********************** GPSLOCATION DATA TYPE(s) ****************************************************************/
/*****************************************************************************************************************/

var queueMutex = &sync.Mutex{}
var once sync.Once
var once_2 sync.Once

var(
	q_instance *list.List
	t_instance *rbt.Tree
)

type GPSLocation struct{
	Location Locationdata `json:location`
	Gpsobject int	  `json:gpsobject`
	Uuid string       `json:uuid`
	Timestamp int64    `json:timestamp`
}

type Locationdata struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Accuracy  float64 `json:"accuracy"`
	Zindex    uint64  `json:zindex`
}

func withinTime(driver_ts int64,detect_ts int64) bool{
	return ((driver_ts - detect_ts)/1e+9 < Timedepth)
}

// Singleton, one instance of Queue...
func GetQueue() *list.List{

	once.Do(func(){
		q_instance = list.New()
	})

	return q_instance
}


/*****************************************A generic location sharing interface ****************************************/
/**********************************************************************************************************************/



type locations interface {

	AddGPSPosition(location GPSLocation) ([] GPSLocation, error)
	GarbageCollect()(error)

	getNearbyObjects(location GPSLocation,a *rbt.Node )[] GPSLocation
	nearbyObject(driver GPSLocation,detect GPSLocation,vtype int) bool

}


/****************************************************QUEUE IMPLEMENTATION *********************************************/
/**********************************************************************************************************************/


type Queue struct{
	Q *list.List
}


func (q Queue) nearbyObject(driver GPSLocation,detect GPSLocation,vehicletype int) bool{
	// if it is the same don't add
	if (driver.Uuid == detect.Uuid) {
		return false
	}

	// don't need to check type its going to be mutual exclusive
	if (withinTime(driver.Timestamp,detect.Timestamp)){
		return true
	}

	return false
}

func (q Queue) getNearbyObjects(location GPSLocation, qu interface{})[] GPSLocation{

	var listofdectees []GPSLocation

	listan := q.Q


	for element := listan.Front();element != nil;element = element.Next(){
		var detectee = element.Value.(GPSLocation)
		if q.nearbyObject(location,detectee,location.Gpsobject){
			listofdectees = append(listofdectees, detectee)
			log.Info(" collisions " , len(listofdectees))
		}
		if len(listofdectees) >= MaxDetections {
			break
		}
	}

	return listofdectees
}

func (q Queue) AddGPSPosition(location GPSLocation)([] GPSLocation, error){

	queueMutex.Lock()

	var neighbours []GPSLocation
	q.Q.PushFront(location)
	neighbours = q.getNearbyObjects(location,nil)

	queueMutex.Unlock()

	return neighbours,nil
}

func (q Queue) GarbageCollect() error{
	return nil
}



/****************************************************RED AND BLACK TREE IMPLEMENTATION ********************************/
/**********************************************************************************************************************/

type TreeExtended struct{
	*rbt.Tree
}

func (t *TreeExtended) GetNodeFromKey(key interface{}) (foundNode *rbt.Node){
	node :=  t.Root

	for node != nil {
		compare := t.Comparator(key, node.Key)
		switch {
		case compare == 0:
			return node
		case compare < 0:
			node = node.Left
		case compare > 0:
			node = node.Right
		}
	}
	return nil
}

// Find predecessor and successor to a tree node. O(h) where the h is height of the tree. h = log n at worst case
func (t* TreeExtended) FindPreSuc(root *rbt.Node,key interface{},pre *rbt.Node,suc *rbt.Node){
	if (root != nil) {

		if t.Comparator(root.Key, key) == 0 {

			// max value in left subtree is predecessor

			if root.Left != nil {
				tmp := root.Left
				for tmp.Right != nil {
					tmp = tmp.Right
				}
				*pre = *tmp

			}

			if root.Right != nil {
				tmp := root.Right
				for tmp.Left != nil {
					tmp = tmp.Left
				}
				*suc = *tmp
			}
			//return pre,suc
		}else if t.Comparator(root.Key, key) == 1 {
            *suc = *root
			t.FindPreSuc(root.Left, key, pre, suc)
		} else {
			*pre = *root
			t.FindPreSuc(root.Right, key, pre, suc)
		}
	}
}



func GetTree() *rbt.Tree{
	once_2.Do(func(){
		t_instance = rbt.NewWith(byGPSIndexation)
	})

	return t_instance
}


// we need a custom comparator for the Tree implementation
func byGPSIndexation(a,b interface {}) int {

	c1  := a.(GPSLocation)
	c2  := b.(GPSLocation)

	zindex_1 := c1.Location.Zindex
	zindex_2 := c2.Location.Zindex

	switch {
	case zindex_1 > zindex_2:
		return 1
	case zindex_1 < zindex_2:
		return -1
	default:
		return 0

	}
}

func (t TreeExtended) nearbyObject(driver GPSLocation,detect GPSLocation,vehicletype int) bool{
	// if it is the same don't add
	if (driver.Uuid == detect.Uuid) {
		return false
	}

	// don't need to check type its going to be mutual exclusive
	if (withinTime(driver.Timestamp,detect.Timestamp)){
		return true
	}

	return false
}

func IsMemberOf(list [] GPSLocation, Key GPSLocation)(bool){

	found := false
	i := 0
	for found == false && i < len(list){
		found = (list[i].Location.Zindex == Key.Location.Zindex) && (list[i].Uuid == Key.Uuid)
		i++
	}
	return found
}

func (t TreeExtended) getNearbyObjects(location GPSLocation)[] GPSLocation{


	var listofdectees []GPSLocation

	found := false
	tmplocation := location

	stack := lls.New()
	stack.Push(tmplocation)

	pre := rbt.Node{}
    var suc rbt.Node

	for found == false  && (stack.Empty() ==  false){

		key,_ := stack.Pop()
		t.FindPreSuc(t.Root,key,&pre,&suc )
		if pre.Key == nil && suc.Key == nil {
			break;
		}

		if pre.Key != nil {
			if (t.nearbyObject(location,pre.Key.(GPSLocation),0)) {
				if (!IsMemberOf(listofdectees, pre.Key.(GPSLocation))) {
					listofdectees = append(listofdectees, pre.Key.(GPSLocation))
					stack.Push(pre.Key) // we don't know if it was time or proximity
				}
			}
		}
		if suc.Key != nil {
			if (t.nearbyObject(location,suc.Key.(GPSLocation),0)) {
				if (!IsMemberOf(listofdectees, suc.Key.(GPSLocation))) {
					listofdectees = append(listofdectees, suc.Key.(GPSLocation))
					stack.Push(suc.Key)
				}
			}
			 // we don't know if it was time or proximity
		}


		found = (len(listofdectees) >= MaxDetections)
	}

	return listofdectees
}

func (t TreeExtended) AddGPSPosition(location GPSLocation)([] GPSLocation, error){

	queueMutex.Lock()

	t.Put(location,location)
	neighbours := t.getNearbyObjects(location)

	queueMutex.Unlock()

	return neighbours, nil

}

func (t TreeExtended) GarbageCollect() error{
	return nil
}

