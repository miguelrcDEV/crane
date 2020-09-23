package types

type MediaServer struct {
	Instance string
	Url      string
	User     string
	Secret   string
}

type Connections struct {
	NumberOfElements int
	Content          []Connection
}

type Connection struct {
	ConnectionId string
	CreatedAt    float64
	Location     string
	Platform     string
	Token        string
	Role         string
	ServerData   string
	ClientData   string
}

type Session struct {
	SessionId              string
	CreatedAt              float64
	MediaMode              string
	RecordingMode          string
	DefaultOutputMode      string
	DefaultRecordingLayout string
	CustomSessionId        string
	Connections            Connections
	Recording              bool
}

type MediaServerResponse struct {
	NumberOfElements int
	Content          []Session
	Health           bool
}

type MediaServerActiveSessions struct {
	MediaServer      MediaServer
	Health           bool
	NumberOfElements int
	Content          []Session
}
