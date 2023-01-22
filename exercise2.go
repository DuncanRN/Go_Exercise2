// Some notes about this exercise:
// I started with no knowledge of Go (of C for that matter) before starting this exercise.
// but because the brief states "(C or GO prefered)" I thought I should accept the challenge!
// I'm fairly happy with it all, except the final part, it will only work for an array of length 3 - a big 
// assumption!
// Otherwise it was enjoyable getting started on a new language.


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
        "io/ioutil"
        )

var  currentTime = time.Now().Unix()

// Here we setup a Device struct containing 
// a name, a type and a list of social links
type Device struct {
    Name   string `json:"Name"`
    Type   string `json:"Type"`
    Info   string `json:"Info"`
    Value  string  `json:"value"` 
    Timestamp string `json:"timestamp"` 
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


    // ====
    // TASK2:  Discard the devices where the timestamp value is before the current time. The timestamps are in UNIX format
    // ====

    for i := 0; i < len(devices.Devices); i++ {
        
        deviceTimestampInt64, err := strconv.ParseInt(devices.Devices[i].Timestamp, 10, 64)
        if err != nil {
            panic(err)
        }

        if (deviceTimestampInt64 < currentTime){

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


    // ====
    // TASK3:  Get the total of all value entries, values are base64 encoded integers
    // ====

    var total = 0
    for i := 0; i < len(devices.Devices); i++ {

        decodedValue, err := base64.StdEncoding.DecodeString(devices.Devices[i].Value)
        if err != nil {
            fmt.Println("Error Found:", err)
            return
        }

        decodedValueToInt, err := strconv.Atoi(string(decodedValue))

        if err != nil {
            fmt.Println("Error during conversion")
            return
        }

        total = total + decodedValueToInt
    }


    // ====
    // TASK4:  Parse the uuid from the info field of each entry
    // ====

    // create array of strings uuids
    var uuids []string

    // loop round our devices, get the devices.Devices[i].Info
    for i := 0; i < len(devices.Devices); i++ {
        
        var Info = devices.Devices[i].Info

        // get the position of the substring "uuid" withing the Info string
        // an example of this larger Info string is 
        // "A Bacnet device uuid:29446300-e583-11ec-8fea-0242ac120002, used to read the light-level"
        posOfUuid := strings.Index(string(Info), "uuid") 

        // now find next comma
        // Here we are making some assumptions, that the string will alwayszs be formated 
        // with a comma immediately after the uuid. 
        // An improvement might be to check for the first whitespace or comma after the uuid
        
        posOfComma := strings.Index(string(Info), ",") 

        substring := Info[(posOfUuid+5):posOfComma]

        // ideally here we would check that the posOfUuid is less than posOfComma,
        // but we are working on the assumption that there is no comma occuring in the string
        // before the uuid
        
        // add that grabbed uuid to our uuids array
        uuids = append(uuids, substring)
    }

    fmt.Println(uuids)


    // ====
    // TASK5:  Output the values total and the list of uuids in the format described 
    // by the JSON schema. Write this data to a file
    // ====

    type myJSON struct {
        Total int
        Uuids []string
    }

    
    // Currently this solution only works for an array of uuids with 3 elements.
    jsondat := &myJSON{Total: total, Uuids: []string{string(uuids[0]), string(uuids[1]), string(uuids[2])}}
    file, _ := json.Marshal(jsondat)
    fmt.Println(string(file))

    // write this data to a file

	_ = ioutil.WriteFile("result.json", file, 0644)


}