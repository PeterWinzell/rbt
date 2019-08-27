package main

import "fmt"


func setupDataPoints()[] GPSLocation{

	var locations []GPSLocation

	gps := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387401,
			Longitude:-122.035179,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233332",
		Timestamp:1,
	}
	gps.Location.Zindex = GetZorderIndex(gps.Location.Latitude,gps.Location.Longitude)
	locations = append(locations, *gps)

	gps2 := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387403,
			Longitude:-122.035170,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233333",
		Timestamp:1,
	}

	gps2.Location.Zindex = GetZorderIndex(gps2.Location.Latitude,gps2.Location.Longitude)
	locations = append(locations, *gps2)

	gps3 := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387404,
			Longitude:-122.035177,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233334",
		Timestamp:1,
	}

	gps3.Location.Zindex = GetZorderIndex(gps3.Location.Latitude,gps3.Location.Longitude)
	locations = append(locations, *gps3)

	gps4 := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387441,
			Longitude:-122.035150,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233335",
		Timestamp:1,
	}

	gps4.Location.Zindex = GetZorderIndex(gps4.Location.Latitude,gps4.Location.Longitude)
    locations = append(locations, *gps4)

	gps5 := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387501,
			Longitude:-122.035122,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233336",
		Timestamp:1,
	}

	gps5.Location.Zindex = GetZorderIndex(gps5.Location.Latitude,gps5.Location.Longitude)
	locations = append(locations, *gps5)

	gps6 := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387201,
			Longitude:-122.035199,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233337",
		Timestamp:1,
	}

	gps6.Location.Zindex = GetZorderIndex(gps6.Location.Latitude,gps6.Location.Longitude)
	locations = append(locations, *gps6)

	gps7 := &GPSLocation{
		Location: Locationdata{
			Latitude:37.387445,
			Longitude:-122.035111,
			Accuracy:1,
		},
		Gpsobject:0,
		Uuid:"ahash-1234-23233-233338",
		Timestamp:1,
	}

	gps7.Location.Zindex = GetZorderIndex(gps7.Location.Latitude,gps7.Location.Longitude)
	locations = append(locations, *gps7)

	return locations

}


func main() {

	// CallChange()

	//var location1 locations
	location1 := Queue{GetQueue()}

	locations := setupDataPoints()

	location2 := TreeExtended{GetTree()}

	fmt.Printf("type of location1 is %T\n", location1)
    fmt.Printf("type of location2 is %T\n", location2)


	location1.AddGPSPosition(locations[0])
	location1.AddGPSPosition(locations[1])
	location1.AddGPSPosition(locations[2])
	location1.AddGPSPosition(locations[3])
	location1.AddGPSPosition(locations[4])
	location1.AddGPSPosition(locations[5])
	l2,err :=location1.AddGPSPosition(locations[6])

	if err == nil{
		fmt.Printf("number of adjacent elemnts in queue is  %d\n", len(l2))
	}

	location2.AddGPSPosition(locations[0])
	location2.AddGPSPosition(locations[1])
	location2.AddGPSPosition(locations[2])
	location2.AddGPSPosition(locations[3])
	location2.AddGPSPosition(locations[4])
	location2.AddGPSPosition(locations[5])

	l4, err := location2.AddGPSPosition(locations[6])



	if err == nil{
		fmt.Printf("number of adjacent elemnts in rbt tree is  %d\n", len(l4))
	}
	//location2.AddGPSPosition(*gps)

	// fmt.Printf
}