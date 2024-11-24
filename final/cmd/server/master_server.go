package main

func main() {
	// In order to reduce the load on the master server, all clients only get the
	// location information of the tablet server through the master, it does not
	// request the master node every time it reads or writes, but connects to the
	// tablet server directly, and the client itself keeps a cache of the location
	// of the tablet server in order to minimize the number of times and frequency
	// of communication with the master.
}
