package bigtable

//
//import (
//	proto "final/proto/external-api"
//	"github.com/syndtr/goleveldb/leveldb"
//	"github.com/syndtr/goleveldb/leveldb/util"
//	"log"
//)
//
//type GFS struct {
//	db *leveldb.DB
//}
//
//func NewGFS() *GFS {
//	return &GFS{}
//}
//
//func (this *GFS) Read(tableName string, prefix string, maxVersion int64) ([]byte, error) {
//	tableFile, err := leveldb.OpenFile(tableName, nil)
//
//	defer func(tableFile *leveldb.DB) {
//		closeErr := tableFile.Close()
//		if err != nil {
//			log.Fatalf(closeErr.Error())
//		}
//	}(tableFile)
//
//	iter := this.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
//	defer iter.Release()
//
//	values := make([]*proto.ValueWithTimestamps, 0)
//	var count int64 = 0
//
//	for iter.Next() {
//		if count >= maxVersion {
//			break
//		}
//
//		key := string(iter.Key())
//		value := string(iter.Value())
//		timestamp, err := extractTimestampFromKey(key)
//		if err != nil {
//			return nil, err
//		}
//
//		values = append(values, &proto.ValueWithTimestamps{
//			Timestamp: timestamp,
//			Value:     value,
//		})
//		count++
//	}
//
//	return values
//}
//
//func (this *GFS) Write(tableName string, key []byte, value []byte) bool {
//	tableFile, err := leveldb.OpenFile(tableName, nil)
//
//	defer func(tableFile *leveldb.DB) {
//		closeErr := tableFile.Close()
//		if err != nil {
//			log.Fatalf(closeErr.Error())
//		}
//	}(tableFile)
//
//	err = this.db.Put(key, value, nil)
//	if err != nil {
//		log.Fatalf("Failed to write to LevelDB: %v", err)
//		return false
//	}
//
//	return true
//}
//
//func (this *GFS) Delete(tableName string, key []byte) bool {
//	// delete all keys
//	if err := this.db.Delete(key, nil); err != nil {
//		log.Fatalf("Delete error: %v", err)
//	}
//
//	return true
//}
//
////func (this *GFS) List(tableName string) map[string][]byte {
////
////}
