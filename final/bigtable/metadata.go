package bigtable

// TODO: HAVEN'T TEST
import (
	"context"
	"errors"
	proto "final/proto/external-api"
	"sync"
)

type TabletMetadata struct {
	ServerAddress  string
	TabletStartRow string
	TabletEndRow   string
}

type MetadataService struct {
	proto.UnimplementedMetadataServiceServer
	mu      sync.RWMutex
	tablets map[string][]TabletMetadata
}

func NewMetadataService() *MetadataService {
	return &MetadataService{
		tablets: make(map[string][]TabletMetadata),
	}
}

// RegisterTablet AddTablet when register a new tablet server
func (m *MetadataService) RegisterTablet(tableName string, metadata TabletMetadata) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tablets[tableName] = append(m.tablets[tableName], metadata)
}

// GetTabletLocation Get tablet server address according to the request
func (m *MetadataService) GetTabletLocation(ctx context.Context, req *proto.TabletLocationRequest) (*proto.TabletLocationResponse, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tablets, ok := m.tablets[req.TableName]
	if !ok {
		return nil, errors.New("table not found")
	}

	for _, tablet := range tablets {
		if req.RowKey >= tablet.TabletStartRow && req.RowKey < tablet.TabletEndRow {
			return &proto.TabletLocationResponse{
				ServerAddress:  tablet.ServerAddress,
				TabletStartRow: tablet.TabletStartRow,
				TabletEndRow:   tablet.TabletEndRow,
			}, nil
		}
	}
	return nil, errors.New("no tablet found for the given row key")
}
