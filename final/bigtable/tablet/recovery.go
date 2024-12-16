package tablet

//func (s *TabletServiceServer) recoverDataFromPersistent(targetStoragePath string, startRow string, endRow string) error {
//	targetTableFile, err := s.CreateTable(targetStoragePath)
//	if err != nil {
//		return fmt.Errorf("failed to create target LevelDB: %v", err)
//	}
//	defer targetTableFile.Close()
//
//	// 遍历所有表
//	for tableName, tableFile := range s.Tables {
//		defer tableFile.Close()
//
//		// 构建用于范围查询的前缀
//		startPrefix := fmt.Sprintf("%s:", startRow)
//		endPrefix := fmt.Sprintf("%s:", endRow)
//
//		iter := tableFile.NewIterator(util.BytesPrefix([]byte(startPrefix)), util.BytesPrefix([]byte(endPrefix)))
//		defer iter.Release()
//
//		// 遍历范围内的数据
//		for iter.Next() {
//			key := string(iter.Key())
//			value := string(iter.Value())
//
//			// 解析出原表中的各个部分（假设和之前写入逻辑中构造key的格式一致，可根据实际调整）
//			parts := ParseKey(key)
//			if len(parts) < 4 {
//				continue
//			}
//			rowKey := parts[0]
//			columnFamily := parts[1]
//			columnQualifier := parts[2]
//			timestamp := parts[3]
//
//			// 构造写入新LevelDB的key（格式需和原写入逻辑一致，可根据实际调整）
//			newKey := BuildKey(rowKey, columnFamily, columnQualifier, timestamp)
//			// 将数据写入新的LevelDB实例
//			err := targetTableFile.Put([]byte(newKey), []byte(value), nil)
//			if err != nil {
//				return fmt.Errorf("failed to write data to target LevelDB for table %s: %v", tableName, err)
//			}
//		}
//	}
//	return nil
//}
//
//func ParseKey(key string) []string {
//	// 这里简单以":"分割来解析key，需根据实际构造key的格式调整
//	return strings.Split(key, ":")
//}
