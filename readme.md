# terminal device (go sdk)

## interface
### PubPropertySync
publish device property message to e-hub or i-hub (sync interface)
### PubPropertyAsync
publish device property message to e-hub or i-hub (async interface) 
### PubEventSync
publish device event message to e-hub or i-hub (sync interface)
### PubEventAsync
publish device event message to e-hub or i-hub (async interface)

### set property(SetProperty)
recv property map(need reply)
func(meta Metadata) (Metadata, error)

### call service(ServiceHandle)
recv service name and param map(need reply)
func(name string, meta Metadata) (Metadata, error)


## example
```go
	//init
	options := &index.Options{
		Token:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY3IiOiIxIiwiYXVkIjoiaWFtIiwiYXpwIjoiaWFtIiwiY3VpZCI6ImlhbXItb3g1Z3VwNzciLCJlaXNrIjoiYzdXZzdaR0RsenZXU1NDS1lfalEyRFZkM0VVN21DZUxjWEtyUF90U3RObz0iLCJleHAiOjE2MDExMDc4ODMsImlhdCI6MTU2OTU3MTg4MywiaXNzIjoic3RzIiwianRpIjoiQTVMVjAzT1Bsc1ZuYndZa1R4Z2dXWiIsIm5iZiI6MCwib3JnaSI6ImlvdGQtZjllZGUzOGYtM2YwMC00OTkyLThhMTMtMGVjNWEwMTk0MzI3Iiwib3d1ciI6InVzci1rTFZWQkRxZCIsInByZWYiOiJxcm46cWluZ2Nsb3VkOmlhbToiLCJydHlwIjoicm9sZSIsInN1YiI6InN0cyIsInRoaWQiOiJpb3R0LTZlYjdhNjUwLWExYTItNGQxMS04OGIwLWZhNmUxNmZjYWQwOCIsInR5cCI6IklEIn0.I0pbK6n--aGES9DNzfQcIDF2YJKfB2ANQdi1lQ5cBTgb1YMRTT7dCF6bQUGHm35uULbrZdZizFVl0s2isaq2TodW9nMVE1wd1BxZguEh5o7I5iMBuX33ws7lZPAdEVgMG6rslZXEsFuS7UVSi4iwJvNFqJKSHOsgJicnYSFrg3wt7uh2bgNJD6SKUmojmzS6dkcI8_4V9kkoZb_89PH5Y-zaw6CcrWDbF_KeDj7Jl6dmP4BmXouHFvv2vNdHTuk6M6tWDsPlVEMz5nD4flE5-i0rWmsm8qv2htX_-caRgWcVuhGwpqGG81aNMORnOqxe8bBz1UGtHT3k_u18e1pfZA",
		Server:   "tcp://192.168.14.120:8055",
		SetProperty: func(meta index.Metadata) (index.Metadata, error) {
			fmt.Println("SetProperty", meta)
			return nil, nil
		},
		ServiceHandle: func(name string, meta index.Metadata) (index.Metadata, error) {
			fmt.Println("ServiceHandle", name, meta)
			return nil, nil
		},
	}
	_, err := mqtt.NewMqtt(options)
	if err != nil {
		panic(err)
	}
	//publish property message
    data := index.Metadata{
        "int32":  10,
        "float":  rand.Float32(),
        "double": rand.Float64(),
        "string": "xxxxxxxxxxxxxxxxx",
    }
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    reply, err := m.PubPropertySync(ctx, data)
    cancel()
    fmt.Println(reply)
    data = index.Metadata{
        "int32":  10,
        "string": "hexing-string",
        "float":  rand.Float32(),
        "double": rand.Float64(),
    }
    //publish event message
    reply, err = m.PubEventSync(context.Background(), "he-event1", data)
    fmt.Println(reply)
```