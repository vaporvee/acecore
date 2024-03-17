module github.com/vaporvee/acecore

go 1.21.6

require (
	github.com/bwmarrin/discordgo v0.27.1
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.3
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/crypto v0.0.0-20210421170649-83a5a9bb288b // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/vaporvee/acecore/custom => ./custom
