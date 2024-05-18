"use client"
import React, { useContext, useEffect, useRef, useState } from "react";
import {
  Typography,
  MenuItem,
  Avatar,
  MenuHandler,
  Menu,
  MenuList
} from "@material-tailwind/react";
import { Select} from "flowbite-react";

import {
  UserIcon,
  ClipboardDocumentListIcon,
  UserGroupIcon,
  Cog6ToothIcon
} from "@heroicons/react/24/solid";
import { fetchData } from "@/app/utils/fetch";
import { formatDate } from "@/app/utils/formatDate";
import PostsPage from "../allPosts";
import Image from "next/image";
import { GroupIdContext } from "@/module/groupId_provide";
import { useHomeData } from "@/app/api/context/HomeDataContext";

// nav list component
const navListItems = [
    {
      label: "posts",
      icon: ClipboardDocumentListIcon,
      active: "bg-indigo-600",
    },
    {
      label: "followers",
      icon: UserGroupIcon,
      active: "bg-white",
    },
    {
      label: "followee",
      icon: UserIcon,
      active: "bg-white",
    },
  ];
export default function Profil({addPost}){
    const [field, setField] = useState("posts")
    const [activ, setActiv] = useState(0)
    const refs = useRef([]);
    const [profil, setProfil] = useState(null)
    const [posts, setPosts] = useState([])
    const [followers, setFollowers] = useState([])
    const [followees, setFollowees] = useState([])
    const [isMenuOpen, setIsMenuOpen] = React.useState(false);
    const [selectedOption, setSelectedOption] = useState("");
    const handleSelectChange = (event) => {
      setSelectedOption(event.target.value);
    };
    const {idUser} = useContext(GroupIdContext)
    const userLogin = useHomeData().Profil.Id

    console.log("Id => ",idUser)
    useEffect(()=>{
      fetchData(`/profil?id=${idUser}`)
      .then(data => {
        if(data?.Profil.IsPublic === false){
          setSelectedOption("private")
        }else{
          setSelectedOption("public")
        }
        setProfil(data?.Profil)
        setPosts(data?.Posts?.reverse())
      })
      .catch(error => {
        console.log("Erreur Profil => ",error)
      })
    },[field, addPost, idUser])
    useEffect(()=>{ 
      if (field==="followers"){
        fetchData(`/followers?id=${idUser}`)
        .then(data => {
            setFollowers(data)
          })
          .catch(erreur => console.log(erreur))
        }else if (field==="followee"){
          fetchData(`/followee?id=${idUser}`)
          .then(data => {
            setFollowees(data)
          })
          .catch(erreur => console.log(erreur))
        }
      },[field, idUser])
      
    return (
        <div className="flex flex-col gap-2 w-full h-[75vh] overflow-y-auto hide-scrollbar">
          <div className="flex bg-white w-full p-5 justify-between max-h-70">
            <div className="flex gap-5 bg-white">
              <Avatar
               src={profil?.Avatar && (profil?.Avatar.endsWith('.jpg') || profil?.Avatar.endsWith('.png') || profil?.Avatar.endsWith('.jpeg')) ? `/uploads/${profil?.Avatar}` : '/uploads/default.png'}
               alt="avatar"
               className="p-0.5 h-40"
               size="xxl"
              />
              <div className="mt-4">
                  <h4 className="text-base font-semibold mb-2">FirstName: <span className="text-sm font-medium ml-2">{profil?.FirstName}</span></h4>
                  <h4 className="text-base font-semibold mb-2">LastName: <span className="text-sm font-medium ml-2">{profil?.LastName}</span></h4>
                  <h4 className="text-base font-semibold mb-2">Birthday: <span className="text-sm font-medium ml-2">{formatDate(profil?.DateOfBirth)}</span></h4>
                  <h4 className="text-base font-semibold mb-2">Nickname: <span className="text-sm font-medium ml-2">{profil?.Nickname}</span></h4>
                  <h4 className="text-base font-semibold mb-2">Email: <span className="text-sm font-medium ml-2">{profil?.Email}</span></h4>
                  <h4 className="text-base font-semibold mb-2">Profil: <span className="text-sm font-medium ml-2">{selectedOption}</span></h4>
                  {profil?.About && <h4 className="text-base font-semibold mb-2">About: <span className="text-sm font-medium ml-2">{profil?.About}</span></h4>}
              </div>
            </div>
            {(idUser===0 || idUser===userLogin)? (<div className="">
              <Menu allowHover open={isMenuOpen} handler={setIsMenuOpen}>
                <MenuHandler>
                <Typography
                  variant="small"
                  color="gray"
                  className="font-medium text-blue-gray-500"
                  >
                  <MenuItem className="relative flex justify-center items-center max-h-7 font-semibold gap-1 text-indigo-600">
                    <span className="">Setting</span>
                    {React.createElement(Cog6ToothIcon, { className: "h-[20px] w-[20px] cursor-pointer" })}{" "}
                  </MenuItem>
                </Typography>
                </MenuHandler>
                <MenuList className="w-16 gap-3 overflow-visible">
                  <ul className="col-span-4 flex w-full flex-col gap-1 font-semibold">
                    <Select id="status" onChange={handleSelectChange} required>
                      <option>public</option>
                      <option>private</option>
                    </Select>
                  </ul>
                </MenuList>
              </Menu>
            </div>):''}
          </div>
          <div className="flex flex-col">
            <ul className="flex gap-0.5 bg-white ">
              {navListItems.map(({ label, icon, active }, key) => (
                <Typography 
                  onClick={()=>{
                      setField(label)
                      refs.current[activ].style.background="white"
                      setActiv(key)
                      refs.current[key].style.background="#5850ec"
                    }
                  }
                  key={label}
                  as="div"
                  variant="small"
                  color="gray"
                  className="font-medium text-blue-gray-500 w-1/3"
                >
                  <MenuItem ref={el =>refs.current[key]=el} className={`grid items-center gap-px  py-0.5 ${active} hover:bg-inherit text-black`}>
                    <div className="flex justify-center">
                    {React.createElement(icon, { className: "h-[25px] w-[25px]" })}{" "}
                    </div>
                    <span className="text-gray-900 text-center text-sm font-semibold"> {label}</span>
                  </MenuItem>
                </Typography>
              ))}
            </ul>
            <div className="border-y-2 flex ">
                {field ==="posts" ?(
                  <PostsPage propsPosts={posts} className="mx-0"/>
                ):field==="followers" ? (<Follow users={followers} />):(<Follow users={followees} />)}
            </div>
          </div>
        </div>
    )
}

function Follow({users}){
  const {setIdUser} = useContext(GroupIdContext)
  return (
    <ul className="posts grid gap-1 grid-cols-1 w-full ">
      {users?.map(({Id, Avatar, FirstName, LastName})=>(
        <li className="cursor-pointer hover:bg-blue-gray-50 flex items-center bg-white shadow-md border-solid rounded-lg p-2 pl-4 justify-between max-h-11" key={Id} onClick={()=>setIdUser(Id)}>
          <Image
            src={Avatar && (Avatar.endsWith('.jpg') || Avatar.endsWith('.png') || Avatar.endsWith('.jpeg')) ? `/uploads/${Avatar}` : '/uploads/default.png'}
            alt="avatar"
            className="rounded-full"
            width = {35}
            height = {35}
            priority={true}
            style={{ border: '2px solid black' }}
          />
          <span className="flex-auto ml-4">{FirstName} {LastName}</span>
        </li>
      ))}
    </ul>
  )
}

