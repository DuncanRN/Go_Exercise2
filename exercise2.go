package main

import (
        "fmt"
        "os"
        "time"
        "io"
        "encoding/json"
        "strconv"
        "encoding/base64"
        // "reflect"
        "strings"
        )

var  currentTime = time.Now().Unix()

// Here we setup a Device struct containing 
// a name, a type and a list of social links
type Device struct {
    Name   string `json:"Name"`
    Type   string `json:"Type"`
    Info   string `json:"Info"`
    Value  string  `json:"value"` // !!! Duncan is there a better way to do this?
                                    // just treat it as a string then DecodeString later?
    Timestamp string `json:"timestamp"`  // again, should Timestamp be an int ? 
    // or even better should we be creating a new struct called UnixTime ?
}

type DeviceUnMarshalled struct {
    Name   string 
    Type   string 
    Info   string 
    Value  string 
    Timestamp string 
}

// current problem - making Timestamp a UnixTime (the type we created) makes us only get the first field

// i, err := strconv.ParseInt("1405544146", 10, 64)
// if err != nil {
//     panic(err)
// }
// tm := time.Unix(i, 0)
// fmt.Println(tm)

// Here we setup a Devices struct containing
// an array of Devices
type Devices struct {
    Devices []Device `json:"Devices"`
}

func removeDevice(s []Device, i int) []Device {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}


type response1 struct {
    Total   int
    Uuids []string
}

func main() {

    // TASK1:  Parse the data from exercise-02/data/data.json

    // Open the JSON File
    dataJsonFile, err := os.Open("data/data.json")
    // if this creates an error then output it.
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Opened data/data.json Successfully ")
    // defer closing the file so that it can be parsed later
    defer dataJsonFile.Close()


    // now we can unmarshall our dataJsonFile

    // read our opened jsonFile as a byte array.
    byteValue, _ := io.ReadAll(dataJsonFile)


    // initialize the Devices  array
    var devices Devices

    // now unmarshal the byte Array - byteValue which contains the data from the file.
    // it is put in 'devices' which we just defined
    json.Unmarshal(byteValue, &devices)


    // now loop around the devices array, and for each device
    // print out the device's Name, Type, Info, Value and Timestamp
    
    for i := 0; i < len(devices.Devices); i++ {
        fmt.Println("Device Name: " + devices.Devices[i].Name)
        // fmt.Println("Device Type: " + devices.Devices[i].Type)
        // fmt.Println("Device Info: " + devices.Devices[i].Info)
        // fmt.Println("Device Value: " + devices.Devices[i].Value)

        // fmt.Println("Device Timestamp: " + devices.Devices[i].Timestamp)
    }


    // fmt.Println("Device Timestamp : " + strconv.ParseInt(devices.Devices[i].Timestamp, 10, 64))

    // TASK2:  Discard the devices where the timestamp value is before the current time. The timestamps are in UNIX format
    
// !!! BIG PROBLEM HERE !!!
// we loop round, delete from the devices but only one is deleted

    // we know how to loop through each device
    // how do we discard from an array

    var devicesStilInDate []DeviceUnMarshalled

    // var lengthOfDevicesBeforeRemovals = len(devices.Devices) 
    for i := 0; i < len(devices.Devices); i++ {
        
        // fmt.Println("Is device timestamp - " + devices.Devices[i].Timestamp + " earlier than current time - "  + strconv.FormatInt(currentTime, 10) + "?")
        
        deviceTimestampInt64, err := strconv.ParseInt(devices.Devices[i].Timestamp, 10, 64)
        if err != nil {
            panic(err)
        }

        if (deviceTimestampInt64 < currentTime){

            // fmt.Println("YES! Device timestamp is earlier than current time. We should discard ")
            devices.Devices = removeDevice(devices.Devices, i)
            i = i -1 

            // this i = i - 1 is to fix the problem where we remove for example the first item 
            // in the list. Every other element moves up the list, the second item becomes the first,
            // the 3rd item becomes the 2nd.
            // our next iteration around the loop after removing the first item is to look at the 
            // 2nd element. But the object in the 2nd element is the original 3rd element.


        }else{
            // previously we created a whole new array with the correct items. 
            // this initial approach has been rejected as it would not scale well with 
            // larger data sets, and is wasteful of memory resources.

            // create a new DeviceUnMarshalled

            
            var newDevice DeviceUnMarshalled   
            newDevice.Name = devices.Devices[i].Name
            newDevice.Type = devices.Devices[i].Type
            newDevice.Info = devices.Devices[i].Info 
            newDevice.Value = devices.Devices[i].Value
            newDevice.Timestamp = devices.Devices[i].Timestamp
            devicesStilInDate = append(devicesStilInDate, newDevice)
            
        }
    }

    fmt.Println("now we output them again, with the earlier timestamped ones removed...")
    fmt.Println("_______")

    for i := 0; i < len(devices.Devices); i++ {
        fmt.Println("Device Name: " + devices.Devices[i].Name)
    }
    fmt.Println("_______")



    for i := 0; i < len(devicesStilInDate); i++ {
    
        fmt.Println("Device Name: " + devicesStilInDate[i].Name)
    }

    
    // TASK3:  Get the total of all value entries, values are base64 encoded integers

    fmt.Println("Task 3")
    var total = 0
    for i := 0; i < len(devicesStilInDate); i++ {
        // fmt.Println("Device Name: " + devices.Devices[i].Name)
        // fmt.Println("Device Type: " + devices.Devices[i].Type)
        // fmt.Println("Device Info: " + devices.Devices[i].Info)
        // fmt.Println("Device Value: " + devicesStilInDate[i].Value)
        // fmt.Println("Decoded Value: " + base64.StdEncoding.DecodeString(""))


        // fmt.Println("var1 = ", reflect.TypeOf(devicesStilInDate[i].Value)) 


        // using the function
        decodedValue, err := base64.StdEncoding.DecodeString(devicesStilInDate[i].Value)
        if err != nil {
            fmt.Println("Error Found:", err)
            return
        }
        fmt.Println("Device Value: " + devicesStilInDate[i].Value + " Decoded value: " + string(decodedValue))


        decodedValueToInt, err := strconv.Atoi(string(decodedValue))

        if err != nil {
            fmt.Println("Error during conversion")
            return
        }

        // fmt.Println(marks)

        // fmt.Println("Decoded decodedValueToInt: " , decodedValueToInt)

        total = total + decodedValueToInt
        // fmt.Println("have just added to the total : " , total)


    }

    fmt.Println("total - ", total)

    // TASK4:  Parse the uuid from the info field of each entry

    // create array of strings uuids
    var uuids []string


    // loop round our remaining devices, get the devices.Devices[i].Info
    for i := 0; i < len(devicesStilInDate); i++ {
        // find nice way to grab uuid, - find where "uuid" appears and start 4 past that...?

        var Info = devices.Devices[i].Info
        fmt.Println("We want uuid from Device Info: " + Info)

        posOfUuid := strings.Index(string(Info), "uuid") 

        fmt.Println("uuid occurs at: " , posOfUuid)

        // now find next comma
        // !!! Duncan this is making some assumptions, should we find next whitespace too 
        // and see if that is first, or get whitespace OR comma?
        posOfComma := strings.Index(string(Info), ",") 

        fmt.Println("comma occurs at: " , posOfComma)

        substring := Info[(posOfUuid+5):posOfComma]

        fmt.Println("substring is " + substring);
        
        // add that grabbed uuid to our uuids array
        uuids = append(uuids, substring)
    }

    fmt.Println(uuids)

    var formattedUuids = "";

    for i := 0; i < len(uuids); i++ {
        formattedUuids +=  uuids[i] + ", ";
    }

    // this clears off the final ", "
    formattedUuids = formattedUuids[0:(len(formattedUuids)-2)]
    fmt.Println("formattedUuids ", formattedUuids)


    // TASK5:  Output the values total and the list of uuids in the format described 
    // by the JSON schema. Write this data to a file

    // const jsonOutput = `{
    //         "uuids" : [` + formattedUuids +  `],
    //         "total" : "` + total + `"
    //         }`

    // fmt.Println("jsonOutput - " + jsonOutput)
    // so something like
    // {
    //     "uuids" : ["134124124", "1231232134", "1233"],
    //     "total" : "100"
    //     }

    // how do we write that to a file?


    res1D := &response1{
        Total:   1,
        Uuids: []string{"apple", "peach", "pear"}}
    res1B, _ := json.Marshal(res1D)
    fmt.Println(string(res1B))


    type response2 struct {
        Total     int `json:"total"`
        Uuids string `json:"uuids"`
    }
    
    
    res2 := &response2{
        Total:     total,
        Uuids: "likes to eat seed",
    }


    // we can use the json.Marhal function to
    // encode the pigeon variable to a JSON string
    data, _ := json.Marshal(res2)
    // data is the JSON string represented as bytes
    // the second parameter here is the error, which we
    // are ignoring for now, but which you should ideally handle
    // in production grade code

    // to print the data, we can typecast it to a string
    fmt.Println(string(data))
    

    val := []interface{}{}
	// val = append(val, 3)
	// val = append(val, "123123")
	// val = append(val, struct {
	// 	Status            string
	// 	CurrentTime       time.Time
	// 	HeartbeatInterval int
	// }{
	// 	"Accepted",
	// 	time.Now(),
	// 	300,
	// })
    val = append(val, struct {
		Total            int
		CurrentTime       time.Time
		HeartbeatInterval int
	}{
		total,
		time.Now(),
		300,
	})

	js, _ := json.Marshal(val)
	fmt.Printf("%#v", string(js))
    fmt.Println("")
    fmt.Println("")

    type myJSON struct {
        Array []string
    }

    jsondat := &myJSON{Array: []string{string(uuids[0]), string(uuids[1]), string(uuids[2])}}
    encjson, _ := json.Marshal(jsondat)
    fmt.Println(string(encjson))


    

    // jsondat := &myJSON{Array: []string{formattedUuids}}
    // encjson, _ := json.Marshal(jsondat)
    // fmt.Println(string(encjson))
    


}