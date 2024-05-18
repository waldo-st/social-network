"use client"
import CreatePostPopup from "../../componants/createPostPopup";

import { Button, Label, Modal, TextInput, } from "flowbite-react";

import React, { useEffect, useRef, useState } from "react";
import {
  Typography,
  MenuItem,
} from "@material-tailwind/react";
import {
    MagnifyingGlassIcon,
    UserGroupIcon,
} from "@heroicons/react/24/solid";
import { useRouter } from "next/navigation";
import { fetchData, pushData } from "@/app/utils/fetch";

export default function GroupeLayoutCompenent() {
    const [childData, setChildData] = useState(null);

    const handleCallback = (data) => {
        setChildData(data);
    }

    return (
        <div className="groupeLayout flex flex-col gap-4 text-blue-700 w-full">
            <div className="flex justify-end border-t-4 border-indigo-500 p-2 bg-white">
                <CreatePostPopup />
            </div>
            <div className="flex  justify-end ">
                <CreateGroupPopup parentCallback={handleCallback}/>
            </div>
            <InfoGroups newGroup={childData}/>
        </div >
    )
}


const navListItems = [
    {
        label: "Discovert Groups",
        icon: MagnifyingGlassIcon,
        active: "bg-indigo-600",
    },
    {
        label: "your Groups",
        icon: UserGroupIcon,
        active: "bg-white",
    },
  ];

function InfoGroups({newGroup}){
    const [field, setField]= useState("Discovert Groups")
    const [activ, setActiv] = useState(0)
    const [groupsJoin, setGroupsJoin] = useState([])
    const [groupsNoJoin, setGroupsNoJoin] = useState([])
    const refs = useRef([]);
    useEffect(()=>{
        if (field ==="your Groups"){
            fetchData(`/joinedGroups`)
            .then(data=>{
                setGroupsJoin(data)
            })
            .catch(error=>{
                console.log("Fetch all Groups => ", error)
            })
        }else{
            fetchData(`/unjoinedGroups`)
            .then(data=>{
                setGroupsNoJoin(data)
            })
            .catch(error=>{
                console.log("Fetch all Groups => ", error)
            })
        }
    },[field, newGroup])
    return (
        <div className="flex flex-col w-full h-500 gap-4">
            <div className="flex flex-col gap-2">
                <ul className="flex gap-0.5 bg-white rounded-md">
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
                        className="font-medium text-blue-gray-500 w-1/2"
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
                <div className="bg-white border-y-2 flex flex-col gap-1 w-full h-[50vh] overflow-y-auto hide-scrollbar rounded-md">
                    {field ==="Discovert Groups" ?(<div className="flex flex-col">
                        {
                            groupsNoJoin?.map(({Id, Title, Description}) => (
                                <GroupsDiscovert key={Id} title={Title} description={Description} I={Id} join={"Join"}/>
                            ))
                        }</div>):(<div>{
                            groupsJoin?.map(({Id, Title, Description}) => (
                                <GroupsDiscovert key={Id} title={Title} description={Description} I={Id} join={"View"}/>
                            ))
                        }</div>)}
                </div>
            </div>
        </div>
    )
}

function GroupsDiscovert({title, description, join, I}) {
    const route = useRouter()
    const joinGroup=()=>{
        if (join === "View"){
            route.push(`/groupes/${I}?group=${title}`)
        }
        if (join === "Join"){
            pushData(null, `/joinGroup/request?gid=${I}&type=join`)
        }
    }
    return (
        <div className=" flex justify-between items-center p-1 pr-5 pl-5 border-b" >
            <div className="flex gap-2 items-center">
                <div className="flex flex-col ">
                    <span className="nameGroup font-semibold text-sm title">{title}</span>
                    <span className="text-xs font-light description">{description}</span>
                </div>
            </div>
            <button  className="text-white bg-indigo-600 rounded-lg w-10 h-6 text-xs cursor-pointer" onClick={joinGroup}>{join}</button>
        </div>
    )
}

function CreateGroupPopup({ parentCallback }) {
    const [openModal, setOpenModal] = useState(false)
    const [title, setTitle] = useState('')
    const [description, setDescription] = useState('')
    function onCloseModal() {
        setOpenModal(false);
        setTitle('');
        setDescription('');
    }

    function managerData() {
        if (title != "" && description != "") {
            const data = {
                title: `${title}`,
                description: `${description}`
            }
            pushData(data,"/group")
            .then(data => {
                parentCallback(data)
                onCloseModal()
            })
            .catch(error => console.log(error))
            setOpenModal(false)
        }
    }
    
    return (
        <>
            <div onClick={() => { setOpenModal(true) }} className=" flex items-center text-indigo-600 cursor-pointer">
                <svg className="fill-current" xmlns="http://www.w3.org/2000/svg" height="25" viewBox="0 -960 960 960" width="30"><path d="M450-450H200v-60h250v-250h60v250h250v60H510v250h-60v-250Z" /></svg>
                <span className="font-semibold text-sm">Create Groupe</span>
            </div>
            <Modal show={openModal} size="md" onClose={onCloseModal} popup>
                <Modal.Header className="flex justify-center">New Group</Modal.Header>
                <Modal.Body>
                    <div className="flex gap-4">
                        <div className="mb-2 block">
                            <Label htmlFor="text" value="Your Title" />
                            <TextInput
                                id="title"
                                placeholder="ex:social-network"
                                value={title}
                                onChange={(e) => { setTitle(e.target.value) }}
                                required />
                        </div>
                        <div className="mb-2 block">
                            <Label htmlFor="text" value="Your Description" />
                            <TextInput
                                id="description"
                                placeholder="ex:c'est fun le codage"
                                value={description}
                                onChange={(e) => { setDescription(e.target.value) }}
                                required />
                        </div>
                    </div>
                    <div className="flex justify-between">
                        <Button onClick={managerData}>Create</Button>
                    </div>
                </Modal.Body>
            </Modal>
        </>
    )
}