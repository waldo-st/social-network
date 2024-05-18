"use client"
import CreatePostPopup from "../componants/createPostPopup";  
import PostsPage from "../componants/allPosts"; 
import { useContext, useEffect, useState } from "react";
import { fetchData } from "../utils/fetch";
import { getCookie } from "../utils/fetch";
import {WebsocketContext, URL_Websocket} from "../../module/websocket_provide";

export default function Home() {
    const [posts, setPosts] = useState([])
    const [childData, setChildData] = useState(null);
    const {setConn} = useContext(WebsocketContext)

    const handleCallback = (data) => {
        setChildData(data);
    }
    
    useEffect(() => {
        const token = getCookie('token');
        const ws = new WebSocket(`${URL_Websocket}/handshake?token=${token}`);
        ws.onopen = () => {
            setConn(ws);
            console.log("Connection established...");
        };
        ws.onerror = (error) => {
            console.log("WebSocket error: ", error);
        };
    },[]);
    
    useEffect(()=>{
        fetchData("/post")
        .then(dataPost =>{
            setPosts(dataPost)
        })
        .catch(error => console.log("ERROR fetcht posts => ",error))
    },[childData])

    return (
        <div className="flex flex-col w-full mx-8 items-center gap-2">
            <div className="flex justify-end border-t-4 border-indigo-500 p-2 w-full bg-white">
                <CreatePostPopup groupeId={`0`} parentCallback={handleCallback}/>
            </div>
            <div className="flex flex-col gap-2 w-full h-[75vh] overflow-y-auto hide-scrollbar">
                <PostsPage propsPosts={posts} />
            </div>
        </div>
    );
}
