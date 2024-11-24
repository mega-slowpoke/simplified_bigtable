package bigtable

//// TabletMetadata represents the metadata of a tablet.
//type TabletMetadata struct {
//	RowStart   string
//	RowEnd     string
//	TabletAddr string
//}
//
//// MetadataServer is responsible for storing tablet metadata and querying tablet locations.
//type MetadataServer struct {
//	sync.RWMutex
//	// A map of row key ranges to tablet server addresses.
//	TabletMap map[string]TabletMetadata
//}
//
//// NewMetadataServer creates a new MetadataServer.
//func NewMetadataServer() *MetadataServer {
//	return &MetadataServer{
//		TabletMap: make(map[string]TabletMetadata),
//	}
//}
//
//// AddTablet adds a new tablet to the metadata server.
//func (server *MetadataServer) AddTablet(rowStart, rowEnd, tabletAddr string) {
//	server.Lock()
//	defer server.Unlock()
//	server.TabletMap[rowStart] = TabletMetadata{
//		RowStart:   rowStart,
//		RowEnd:     rowEnd,
//		TabletAddr: tabletAddr,
//	}
//}
//
//func (server *MetadataServer) GetTablet(rowKey string) (string, error) {
//	server.RLock()
//	defer server.RUnlock()
//
//	for _, tablet := range server.TabletMap {
//		// Check if the rowKey is within the tablet's range.
//		if rowKey >= tablet.RowStart && rowKey <= tablet.RowEnd {
//			return tablet.TabletAddr, nil
//		}
//	}
//
//	return "", errors.New("tablet not found")
//}
