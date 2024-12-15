package tablet

//import (
//	"context"
//	"fmt"
//)
//
//func (s *TabletServiceServer) recoverDataFromPersistent(targetStoragePath string) error {

//	deadKeys, err := deadClient.GetAllKeys()
//	if err != nil {
//		return fmt.Errorf("failed to get keys from dead LevelDB: %v", err)
//	}
//
//	ctx := context.Background()
//	for _, key := range deadKeys {
//		if err != nil {
//			return fmt.Errorf("Failed to get value for key %s from dead LevelDB: %v", key, err)
//		}
//		err = takeOverClient.Put(key, value)
//		s.Write(ctx)
//		if err != nil {
//			return fmt.Errorf("Failed to put value for key %s to take over LevelDB: %v", key, err)
//		}
//	}
//
//	return nil
//}
