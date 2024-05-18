import { Card, Modal, Badge, Avatar, TextInput } from "flowbite-react";
import { PiWechatLogoFill } from "react-icons/pi";
import { FaTelegramPlane } from "react-icons/fa";
import { useState, useContext, useEffect, useRef } from "react";
import { URL_Websocket, WebsocketContext } from "@/module/websocket_provide";
import { getCookie } from "@/app/utils/fetch";
import { formatDate } from "../allPosts";

export const ChatRoom = ({ title, room, type, style, id, idLogin }) => {
    const [openModal, setOpenModal] = useState(false)
    const [messages, setMessages] = useState([])
    const {connChat, setConnChat} = useContext(WebsocketContext)
    const {chat,setChat} = useContext(WebsocketContext)
    const textMsg = useRef()

    const onCloseModal = () => {
        setOpenModal(false)
        textMsg.current.value=""
    }

    const initChat = () =>{
        setOpenModal(true)
        const token = getCookie("token")
        if (type === "chatGroup"){
            const Ws = new WebSocket(`${URL_Websocket}/groupChat?roomId=${id}&token=${token}`)
            Ws.onopen = () => {
                setConnChat(Ws);
                console.log("Connection Chat Room established...");
            };
            Ws.onerror = (error) => {
                console.log("WebSocket error: ", error);
            };
        }
        else{
            const Ws = new WebSocket(`${URL_Websocket}/chat?id=${id}&token=${token}`)
            Ws.onopen = () => {
                setChat(Ws);
                console.log("Connection Chat Private established...");
            };
            Ws.onerror = (error) => {
                console.log("WebSocket error: ", error);
            };
        }
    }

    useEffect(()=>{
        if (type === "chatGroup"){
            if (connChat!==null){
                console.log("Cht msg => ")
                connChat.onmessage = (msg)=>{
                    const m = JSON.parse(msg.data);
                    if (Array.isArray(m)){
                        setMessages([...messages, ...m])

                    }else{
                        setMessages([...messages, m])
                    }
                }
                connChat.onclose=()=>{
                    console.log("Connexion is close...")
                }
            }
        }else{
            if (chat!==null){
                chat.onmessage = (msg)=>{

                }
                chat.onclose=()=>{
                    console.log("Connexion is close...")
                }
            }
        }
    },[textMsg, connChat, chat, messages])
    const SenMessage = () => {
        if (!textMsg.current?.value) return
        if (type === "chatGroup"){
            const dataMsg = {
                type :"group",
                content :textMsg.current.value
            }
            if (connChat!==null){
                connChat.send(JSON.stringify(dataMsg))
                textMsg.current.value=""
            }
        }else{
            const dataMsg = {
                type :"private",
                content :textMsg.current.value
            }
            if (chat!==null){
                chat.send(JSON.stringify(dataMsg))
                textMsg.current.value=""
            }
        }
    }
    
    return (
        <>
            <Card className={`${style} flex justify-center items-center text-indigo-600 cursor-pointer`} onClick={initChat}>
                <div className="flex flex-col items-center">
                    <PiWechatLogoFill className="h-6 w-6" />
                    <span>{title}</span>
                </div>
            </Card>
            <Modal show={openModal} size="md" onClose={onCloseModal} popup className="border">
                <Modal.Header className="flex justify-center items-center bg-gray-100">
                <Avatar img={""} rounded />
                <span>{room}</span>
                </Modal.Header>
                <Modal.Body className="flex flex-col justify-center gap-2">
                    <div className="flex flex-col gap-2 w-full min-h-50 overflow-y-auto hide-scrollbar">
                        <ul className="flex flex-col">
                            {Array.isArray(messages) && messages?.map(({Content, SenderId, Username, CreatedAt},index)=>(
                                <li key={index}>
                                    {SenderId===idLogin ? (<Send message={Content} username={Username} creatAt={formatDate(CreatedAt)}/>):(<Receiver  message={Content} username={Username} creatAt={formatDate(CreatedAt)}/>)}
                                </li>
                            ))}
                        </ul>
                    </div>
                </Modal.Body>
                <Modal.Footer className="h-12">
                    <div className="flex justify-center items-center gap-2 w-full">
                        <TextInput className="w-full"
                            ref={textMsg}
                            id="message"
                            placeholder="Saisisser votre message ici"
                            required />
                        <FaTelegramPlane className="h-8 w-8 p-1 text-white cursor-pointer rounded-full bg-cyan-500" onClick={SenMessage} />
                    </div>
                </Modal.Footer>
            </Modal>
        </>
    )
}
const Send = ({ username, message, creatAt }) => {
    return (
        <>
            <div className="flex flex-col items-center self-end">
                <span className="self-end text-xs">{username}</span>
                <Badge className="flex items-center self-end gap-2 min-w-40 text-white rounded-l-3xl bg-teal-400 rounded-br-3xl ">
                    {message.img ? (<image />) : ("")}
                    {message}
                </Badge>
                <span className="self-end text-[10px] text-gray-500">{creatAt}</span>
            </div>
        </>
    )
}
const Receiver = ({ username, message, creatAt }) => {
    return (
        <>
            <div className="flex flex-col items-center self-start">
                <span className="self-start text-xs">{username}</span>
                <Badge className="flex items-center self-start gap-2 min-w-40 bg-blue-500 text-white rounded-r-3xl rounded-bl-3xl">
                    {message}
                </Badge>
                <span className="self-start ml-2 text-[10px] text-gray-500">{creatAt}</span>
            </div>
        </>
    )
}