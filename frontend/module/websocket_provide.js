import React, {useState,createContext} from "react"

const URL_Websocket ="ws://localhost:8080"
const WebsocketContext = createContext(
    {
        conn:null,
        setConn:(c)=>{},
        connChat:null,
        setConnChat:(c)=>{},
        chat:null,
        setChat:(c)=>{}
    }
);

const WebsocketProvider =({children})=>{
    
    const [conn, setConn] = useState(null)
    const [connChat, setConnChat] = useState(null)
    const [chat, setChat] = useState(null)
    
    return (
        <WebsocketContext.Provider
        value={
                {
                conn:conn,
                setConn:setConn,
                connChat:connChat,
                setConnChat:setConnChat,
                chat:chat,
                setChat:setChat,
                }
            }
            >
            {children}
        </WebsocketContext.Provider>
    )
}

export {URL_Websocket, WebsocketContext}
export default WebsocketProvider