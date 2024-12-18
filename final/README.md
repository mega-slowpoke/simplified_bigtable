## Demo Video
- TODO

## Instruction to run code
1. Environment Setup: You should run code in the docker container provided by the lecture, you can use `go mod tidy` to download all 
2. You should the three steps below to run and test the code 
   - First, run `master` in the cmd/server/master with command 
      ```
     ./master --master_address=localhost:9090
     ``` 
     - If it succeeds, you will see a prompt `Master server is running on localhost:9090...`
     - You can change `localhost:9090` to other available ports. 
     
   - Second, run `tablet` in the cmd/server/tablet with command
      ```
     ./tablet --tablet_address=localhost:10000 --master_address=localhost:9090 --max_table_cnt=5 --check_max_period=100` 
      ```
     - If it succeeds, you will see `Tablet is listening on localhost:10000...` on the tablet side, and `Tablet server 'localhost:10000' registered successfully` on the master side. This service discovery is to notify master that new tablet go online, and it can use it later
     - You can run as many as tablet as you like (preferably 3 or 5). The `master_address` parameter here must be identical to the one you entered when you start master
     max_table_cnt (default 5) and check_max_period (default 100) are optional, they are used to specify how many tables a tablet can take care of at most and every how many microsecond the tablet will check
     if its table number is greater than max_table_cnt
   
    - Now, you can run client test provided in the test/integration_test/client_test.go to test the six interfaces we provide
      - CreateTable, DeleteTable, Read, Write, Delete (You have to CreateTable first before you Read, Write and Delete things in this table)
      - GetTabletLocation
        
    

## Group Work
- Chaoxu Wu: Architecture Design, Tablet server, Writeup, Testing
- Yuzhe Ruan: Master server, Client library, Writeup, Testing
