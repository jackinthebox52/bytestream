curl -X POST http://localhost:8080/api/stream/start \
  -H "Content-Type: application/json" \
  -d '{"name":"NFL", "origin":"https://reliabletv.me", "url":"https://s-c3.aistrem.net/plyvivo/d0do5e3uxi7a708opof0/chunklist.m3u8"}'