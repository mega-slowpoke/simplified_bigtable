## Demo Video
- The demo video is included in the .tar.gz file, or you can access it here https://drive.google.com/file/d/1WTscIXuerrIedZLm_w8yzACyfGTpIs_C/view?usp=sharing

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
        
## Instruction to reproduce the demo
All LevelDB file will be stored in cmd/server/tablet

1. Just follow the video to start master server and tablet server with compiled file in the cmd/server
2. You can run test/client_test.go to do the client createTable/deleteTable/read/write/delete. You can create you own tests
3. if you want to reproduce sharding: 
   - Start a server and a tablet server (start only one tablet in the beginning, otherwise you might run into problems caused by LevelDB SDK version incompatibility) 
   - Uncomment "TestClientShard" in the client_test.go, this is used to create many tables, you can change the data
   - Now You can start more tablets one by one (as many as you want), you will see how extra tables be moved to available tablet severs
4. if you want to reproduce the fault tolerance (recovery), you can ctrl+c any of the tablet server, you might need to wait for a few seconds to see another tablet server recover its data


## Group Work
- Chaoxu Wu: Architecture Design, Tablet server, Writeup, Testing
- Yuzhe Ruan: Master server, Client library, Writeup, Testing
