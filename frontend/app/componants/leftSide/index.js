"use client"
import { Avatar } from "@material-tailwind/react";
import React, { useContext, useEffect, useState } from "react";
import Image from 'next/image';
// import { FaUserPlus, ImUserPlus } from "react-icons/fa6";
import { ImUserMinus, ImUserPlus, ImUserCheck } from "react-icons/im";
import { RiMessage2Fill } from "react-icons/ri";
import {
  Card,
} from "@material-tailwind/react";
import { pushData } from "@/app/utils/fetch";
import { ChatRoom } from "../chat";
import { useHomeData } from "@/app/api/context/HomeDataContext";
import { WebsocketContext } from "@/module/websocket_provide";

function UserCard({info}){
  return(
    <div className="grid items-center justify-center gap-4 border-blue-gray-700 border rounded-xl p-4">
      <div className="w-full grid justify-center">
        <Image
          src={`/uploads/${info?.Avatar}`}
          alt="avatar"
          className="rounded-full"
          width = {70}
          height = {70}
          priority={true}
        />
      </div>
      <div className="grid justify-center">
        <h6 className="block font-sans text-base antialiased font-semibold leading-relaxed tracking-normal text-inherit text-center">
          {`${info?.FirstName} ${info?.LastName}`}
        </h6>
        <p className="text-blue-400">{`${info?.Email}`}</p>
      </div>
      <div className="flex gap-2 border-separate">
        <p><span className="font-bold">{`${info?.NbrOfFollowers}`}</span> Followers</p>|
        <p><span className="font-bold">{`${info?.NbrOfFollowee}`}</span> Following</p>
      </div>
    </div>
  );
}

let tabMsgPrivates = []
export function UserSuggest({avatar, name, id, isFollow, Login}){
  const [listUserAdd, setListUserAdd] = useState([])
  const [msgs, setMsgs] = useState([])
  const {chat} = useContext(WebsocketContext)

  function Submit() {
    pushData(null, `/follow/request?id=${id}&type=follow`)
    .then(data=>{
      console.log("Data userSuggeste => ",data)
      setListUserAdd((prev)=> [...prev, id])
    })
    .catch(error => {
      console.log(error)
    })
  }

  if (chat!==null){
    chat.onmessage= (event) => {
      const data = JSON.parse(event.data);
        console.log("Chat Private => ",tabMsgPrivates)
        tabMsgPrivates=[...tabMsgPrivates, ...data]
        console.log("Chat Private => ",tabMsgPrivates)
        setMsgs(tabMsgPrivates)
    };
  }
  useEffect(()=>{},[listUserAdd])

   // Check if the avatar has a valid extension (.jpg or .png)
   const isValidExtension = avatar && (avatar.endsWith('.jpg') || avatar.endsWith('.png'));
   const avat = isValidExtension ? `/uploads/${avatar}` : '/uploads/default.png';

  return (
      <div className="flex items-center mt-1 shadow-md border-solid rounded-lg p-1 justify-between max-h-10">
        <Avatar
        src={avat}
        alt="avatar"
        withBorder={true}
        className="p-0.5"
        size="sm"
        />
        <span className="flex-auto ml-4">{name}</span>
        <div className="flex flex-1 justify-end p-2 ">
          {isFollow === "chat" ? <ChatRoom title={""} room={name} type={"chatPrivate"} data={msgs} id={id} style={"h-7 w-7"} idLogin={Login} /> : (listUserAdd.includes(id) || isFollow === "pending")?<ImUserCheck className="w-5 h-5"/>:isFollow ==="false"?<ImUserPlus className=" cursor-pointer w-5 h-5" onClick={Submit}/>:<ImUserMinus className=" cursor-pointer w-5 h-5"/>}          
        </div>
      </div>
  )
}

export default function LeftSide({propsLeft}) {
  
  return (
    <Card className="h-[calc(90vh-2rem)] w-full max-w-[20rem] p-4 shadow-xl shadow-blue-gray-900/5">
      <UserCard info={propsLeft?.Profil}/>
      <div className="grid mt-5 ml-1">
        <h2 className="font-semibold text-sm text-black">Suggestions</h2>
        <div className="flex gap-1 flex-col w-full overflow-y-auto h-[45vh] hide-scrollbar">
          {propsLeft?.Users.map(({Id,Avatar,FirstName, IsFollowee}) => (
            <UserSuggest key={Id} avatar={Avatar} name={FirstName} isFollow={IsFollowee} id={Id} Login={propsLeft.Profil.Id}/>
          ))}
        </div>
      </div>
    </Card>
  );
}
