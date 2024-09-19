package mw

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	raw, _ := protojson.Marshal((req).(proto.Message)) // для превращения protbuf структур в json используем google.golang.org/protobuf/encoding/protojson пакет а не encoding/json
	log.Printf("request: method: %v, req: %v\n", info.FullMethod, string(raw))

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("response: method: %v, err: %v\n", info.FullMethod, err)
		return
	}

	rawResp, _ := protojson.Marshal((resp).(proto.Message))
	log.Printf("response: method: %v, resp: %v\n", info.FullMethod, string(rawResp))

	return
}

func WithHTTPLoggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hello_World")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", SwaggerUrlForCors)

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")

			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.WriteHeader(200)
			return
		}
		log.Println(r.Method)
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
