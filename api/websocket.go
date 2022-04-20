package api

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize: 1024,
// 	WriteBufferSize: 1024,
// }

// func reader(conn *websocket.Conn) {
// 	for {
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		log.Println(string(p))

// 		if err:=conn.WriteMessage(messageType, p); err != nil{
// 			log.Println(err)
// 			return
// 		}

// 	}
// }

// func wsEndpoint(c *gin.Context) {
// 	fmt.Println("connected =====================")
// 	token := c.Request.URL.Query().Get("Authorization")
// 	fmt.Println("token: ", token)
// 	if token == "" {
// 		c.Writer.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	fmt.Println("token: ", token)
// 	claims, err := jwt.ExtractClaims(token, []byte("7VFGY5ArECvjhRU6wuLq"))
// 	if err != nil {
// 		c.Writer.WriteHeader(http.StatusForbidden)
// 		c.Writer.Write([]byte(err.Error()))
// 		c.Writer.Header().Set("Content-Type", "application/json")
// 		return
// 	}
// 	fmt.Println("user id : ", claims["id"].(string))
// 	webSkt.serveWs(hub, claims["id"].(string), c.Writer, c.Request)
// 	upgrader.CheckOrigin = func(r *http.Request) bool {return true}

// 	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	log.Println("successfully connected")

// 	bytemsg, _ := json.Marshal("halo")
// 	ws.WriteMessage(1, bytemsg)
// 	reader(ws)
// }
