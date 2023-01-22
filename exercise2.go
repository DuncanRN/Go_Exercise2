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


// Currently not being used - DELETE !!!! ???
type DeviceUnMarshalled struct {
    Name   string 
    Type   string 
    Info   string 
    Value  string 
    Timestamp string 
}

// Here we setup a Devices struct containing
// an array of Devices
type Devices struct {
    Devices []Device `json:"Devices"`
}

func removeDevice(s []Device, i int) []Device {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
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


    // initialize the Devices array
    var devices Devices

    // now unmarshal the byte Array - byteValue which contains the data from the file.
    // it is put in 'devices' which we just defined
    json.Unmarshal(byteValue, &devices)

    // now loop around the devices array, and for each device
    // print out the device's Name, Type, Info, Value and Timestamp
    
    // for i := 0; i < len(devices.Devices); i++ {
    //     fmt.Println("Device Name: " + devices.Devices[i].Name)
    // }


    // ====
    // TASK2:  Discard the devices where the timestamp value is before the current time. The timestamps are in UNIX format
    // ====

    for i := 0; i < len(devices.Devices); i++ {
        
        deviceTimestampInt64, err := strconv.ParseInt(devices.Devices[i].Timestamp, 10, 64)
        if err != nil {
            panic(err)
        }

        if (deviceTimestampInt64 < currentTime){

            // fmt.Println("YES! Device timestamp is earlier than current time. We should discard ")
            devices.Devices = removeDevice(devices.Devices, i)
            i = i -1 

            // this  previous line +i = i - 1+  fixes the problem where we remove, for example the 1st item 
            // in the list. Every other element moves up the list, the 2nd item becomes the 1st,
            // the 3rd item becomes the 2nd.....
            // our next iteration around the loop after removing the first item is to look at the 
            // 2nd element. But the object in the 2nd element is the original 3rd element.


        }else{
            // previously we created a whole new array with the correct items. 
            // this initial approach has been rejected as it would not scale well with 
            // larger data sets, and is wasteful of memory resources.
            // create a new DeviceUnMarshalled

            /*
            var newDevice DeviceUnMarshalled   
            newDevice.Name = devices.Devices[i].Name
            newDevice.Type = devices.Devices[i].Type
            newDevice.Info = devices.Devices[i].Info 
            newDevice.Value = devices.Devices[i].Value
            newDevice.Timestamp = devices.Devices[i].Timestamp
            devicesStilInDate = append(devicesStilInDate, newDevice)
            */
        }
    }

    fmt.Println("now we output the items, with the earlier timestamped ones removed...")
    fmt.Println("_______")

    for i := 0; i < len(devices.Devices); i++ {
        fmt.Println("Device Name: " + devices.Devices[i].Name)
    }
    fmt.Println("_______")





    
    // TASK3:  Get the total of all value entries, values are base64 encoded integers

    fmt.Println("Task 3")
    var total = 0
    for i := 0; i < len(devices.Devices); i++ {
        // fmt.Println("Device Name: " + devices.Devices[i].Name)

        decodedValue, err := base64.StdEncoding.DecodeString(devices.Devices[i].Value)
        if err != nil {
            fmt.Println("Error Found:", err)
            return
        }
        fmt.Println("Device Value: " + devices.Devices[i].Value + " Decoded value: " + string(decodedValue))

        decodedValueToInt, err := strconv.Atoi(string(decodedValue))

        if err != nil {
            fmt.Println("Error during conversion")
            return
        }

        total = total + decodedValueToInt
    }

    fmt.Println("total - ", total)

    // ====
    // TASK4:  Parse the uuid from the info field of each entry
    // ====
    fmt.Println("Task 4")

    // create array of strings uuids
    var uuids []string

    // loop round our devices, get the devices.Devices[i].Info
    for i := 0; i < len(devices.Devices); i++ {
        
        var Info = devices.Devices[i].Info
        fmt.Println("We want uuid from Device Info: " + Info)
        // get the position of the substring "uuid" withing the Info string
        // an example of this larger Info string is 
        // "A Bacnet device uuid:29446300-e583-11ec-8fea-0242ac120002, used to read the light-level"
        posOfUuid := strings.Index(string(Info), "uuid") 

        fmt.Println("uuid occurs at: " , posOfUuid)

        // now find next comma
        // Here we are making some assumptions, that the string will alwayszs be formated 
        // with a comma immediately after the uuid. 
        // An improvement might be to check for the first whitespace or comma after the uuid
        
        posOfComma := strings.Index(string(Info), ",") 

        fmt.Println("comma occurs at: " , posOfComma)

        substring := Info[(posOfUuid+5):posOfComma]

        // ideally here we would check that the posOfUuid is less than posOfComma,
        // but we are working on the assumption that there is no comma occuring in the string
        // before the uuid
        
        // add that grabbed uuid to our uuids array
        uuids = append(uuids, substring)
    }

    fmt.Println(uuids)

    /*
    var formattedUuids = "";

    for i := 0; i < len(uuids); i++ {
        formattedUuids +=  uuids[i] + ", ";
    }

    // this clears off the final ", "
    formattedUuids = formattedUuids[0:(len(formattedUuids)-2)]
    fmt.Println("formattedUuids ", formattedUuids)
    */

    // TASK5:  Output the values total and the list of uuids in the format described 
    // by the JSON schema. Write this data to a file

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

    // Currently this solution only works for an array of uuids with 3 elements.

    jsondat := &myJSON{Array: []string{string(uuids[0]), string(uuids[1]), string(uuids[2])}}
    encjson, _ := json.Marshal(jsondat)
    fmt.Println(string(encjson))

    // write this data to a file

    

    // AS a final thing do we close our JSON file? We Defered that earlier...
    // !!!!!

}