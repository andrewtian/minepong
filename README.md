# minepong

Golang Minecraft server querier. Not compliant with de facto Minecraft standards.

# Usage

The `minepong.Ping` method takes care of writing to and reading from a
`net.Conn` interface.

    pong, err := Ping(svr.conn, svr.host)
    if err != nil {
        fmt.Println(err)
        return
    }

A pong is returned information from the query.

    type Pong struct {
        Version struct {
            Name     string
            Protocol int
        } `json:"version"`
        Players struct {
            Max    int `json:"max"`
            Online int `json:"online"`
            Sample []map[string]string
        } `json:"players"`
        Description interface{} `json:"description"`
        FavIcon     string      `json:"favicon"`
    }
