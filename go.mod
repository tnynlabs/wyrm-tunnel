module github.com/tnynlabs/wyrm-tunnel

go 1.15

require (
	github.com/golang/protobuf v1.5.2
	github.com/joho/godotenv v1.3.0
	github.com/tnynlabs/wyrm v0.0.0
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace "github.com/tnynlabs/wyrm" v0.0.0 => "./wyrm"
