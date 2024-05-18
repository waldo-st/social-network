"use client";
import Link from "next/link";
import React, { useContext, useState, useEffect } from "react";
import Logout from "@/app/utils/logout";
import {
  Navbar,
  Collapse,
  Typography,
  MenuItem,
  IconButton,
  Menu,
  MenuHandler,
  MenuList,
} from "@material-tailwind/react";
import {
  UserIcon,
  HomeIcon,
  UserGroupIcon,
  Bars2Icon,
  ChatBubbleOvalLeftEllipsisIcon,
  BellIcon,
} from "@heroicons/react/24/solid";
import { WebsocketContext } from "@/module/websocket_provide";
import Image from "next/image";
import { pushData } from "@/app/utils/fetch";
import { GroupIdContext } from "@/module/groupId_provide";
// nav list component
const navListItems = [
  {
    label: "Home",
    link: "/",
    icon: HomeIcon,
  },
  {
    label: "Groupes",
    link: "/groupes",
    icon: UserGroupIcon,
  },
  {
    label: "Profile",
    link: "/profile",
    icon: UserIcon,
  },
];
function NavList() {
  const {setIdUser} = useContext(GroupIdContext)
  return (
    <ul className="mt-10 mb-4 flex flex-col gap-3 lg:mb-0 lg:mt-0 lg:flex-row lg:items-center">
      {navListItems.map(({ label, link, icon }) => (
        <li key={label} className="font-medium rounded-xl p-3 w-20 text-blue-gray-500 hover:bg-blue-gray-200" onClick={()=>setIdUser(0)}>
          <Link href={link}>
            <div className="grid justify-center items-center gap-px rounded-full cursor-pointer">
              <div className="flex justify-center">
                {React.createElement(icon, { className: "h-[25px] w-[25px]" })}
              </div>
              <span className="text-gray-900 text-xs">{label}</span>
            </div>
          </Link>
        </li>
      ))}
    </ul>
  );
}
// nav list menu
const defaultItems = [
  {
    username: "Sarah",
    message: "Hello bro, have you finished your tasks?",
  },
  {
    username: "John",
    message: "Yo man I'm waiting for your component",
  },
  {
    username: "Michael",
    message: "Boy N'kamou?",
  },
];
const notifItems = [];
function MessageNotify({ nbr, msgNotificationItems }) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const renderItems = msgNotificationItems.map(({ username }, index) => (
    <li className="flex gap-2 justify-between hover:bg-blue-gray-50 rounded-md p-1" key={`${username}-${index}`}>
      <div>Message from </div>
      <div className="font-semibold">{username}</div>
      <div className="flex justify-items-end">
        <button>
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" className="w-6 h-6">
          <path fillRule="evenodd" d="M4.848 2.771A49.144 49.144 0 0 1 12 2.25c2.43 0 4.817.178 7.152.52 1.978.292 3.348 2.024 3.348 3.97v6.02c0 1.946-1.37 3.678-3.348 3.97a48.901 48.901 0 0 1-3.476.383.39.39 0 0 0-.297.17l-2.755 4.133a.75.75 0 0 1-1.248 0l-2.755-4.133a.39.39 0 0 0-.297-.17 48.9 48.9 0 0 1-3.476-.384c-1.978-.29-3.348-2.024-3.348-3.97V6.741c0-1.946 1.37-3.68 3.348-3.97ZM6.75 8.25a.75.75 0 0 1 .75-.75h9a.75.75 0 0 1 0 1.5h-9a.75.75 0 0 1-.75-.75Zm.75 2.25a.75.75 0 0 0 0 1.5H12a.75.75 0 0 0 0-1.5H7.5Z" clipRule="evenodd" />
        </svg>
        </button>
      </div>
    </li>
  ));
  return (
    <>
      <Menu allowHover open={isMenuOpen} handler={setIsMenuOpen}>
        <MenuHandler>
          <Typography
            variant="small"
            color="gray"
            className="font-medium text-blue-gray-500"
          >
            <MenuItem className="relative">
              {React.createElement(ChatBubbleOvalLeftEllipsisIcon, {
                className: "h-[30px] w-[30px]",
              })}{" "}
              {nbr!==0 && <span className="bg-red-900 flex justify-center items-center text-white text-sm absolute top-2 right-2 h-4 w-4 rounded-full">
                {nbr}
              </span>}
            </MenuItem>
          </Typography>
        </MenuHandler>
        <MenuList className="w-[20rem] gap-3 overflow-visible">
          <ul className="col-span-4 flex w-full flex-col gap-1">{renderItems}</ul>
        </MenuList>
      </Menu>
    </>
  );
}
const Paths=(type)=>{
  return type === "follow" ? "/follow/reply?id=" : (type === "join"||type==="invite") ? `/joinGroup/reply?type=${type}&id=`:""
}
function DemandeFollower({ userName, id, message, type }) {
  const path = Paths(type)
  console.log("paths =>",type)
  const handleAccept = async () => {
    try {
      await pushData(null, `${path}${id}&status=accepted`);
      console.log("User accepted");
    } catch (error) {
      console.error("Error accepting request:", error);
    }
  };
  const handleDeny = async () => {
    try {
      await pushData(null, `${path}${id}&status=rejected`);
      console.log("User denied");
    } catch (error) {
      console.error("Error denying request:", error);
    }
  };
  return (
    <div className="flex gap-1 items-center justify-between">
      <div className="flex gap-1 justify-center items-center">
        <span className="font-semibold text-sm">{userName}</span>
        <span className="text-xs">{message}</span>
      </div>
      <div className="flex justify-left gap-1 items-center">
        <button
          onClick={handleAccept}
          className="text-white bg-indigo-600 rounded-lg max-w-10 p-1 font-semibold flex justify-center items-center m-1 h-5 text-xs cursor-pointer"
        >
          Yes
        </button>
        <button
          onClick={handleDeny}
          className="text-white bg-red-600 rounded-lg max-w-10 p-1 font-semibold flex justify-center items-center m-1 h-5 text-xs cursor-pointer"
        >
          No
        </button>
      </div>
    </div>
  );
}
function Notification({ nbr, msgNotificationItems, onRespond }) {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const renderItems = msgNotificationItems?.map(({ username, id, message, type }, index) => (
    <li key={`${username}-${index}`}>
      <DemandeFollower userName={username} id={id} message={message} type={type}
      // onRespond={onRespond}
       />
      <hr />
    </li>
  ));
  return (
    <>
      <Menu allowHover open={isMenuOpen} handler={setIsMenuOpen}>
        <MenuHandler>
          <Typography
            variant="small"
            color="gray"
            className="font-medium text-blue-gray-500"
          >
            <MenuItem className="relative">
              {React.createElement(BellIcon, {
                className: "h-[30px] w-[30px]",
              })}{" "}
              {nbr!==0&&nbr!==undefined&&<span className="bg-red-900 flex justify-center items-center text-white text-sm absolute top-2 right-2 h-4 w-4 rounded-full">
                {nbr}
              </span>}
            </MenuItem>
          </Typography>
        </MenuHandler>
        <MenuList className="w-[20rem] gap-3 overflow-visible">
          <ul className="col-span-4 flex w-full flex-col gap-1">{renderItems}</ul>
        </MenuList>
      </Menu>
    </>
  );
}
function ComplexNavbar() {
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [notifications, setNotifications] = useState(notifItems);
  const toggleIsNavOpen = () => setIsNavOpen((cur) => !cur);
  const { conn } = useContext(WebsocketContext);
    if (conn !== null) {
      conn.addEventListener("message", (event) => {
        const data = JSON.parse(event.data);
        console.log("Message from server ", data);
        setNotifications(data);
      });
    }

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth >= 960) {
        setIsNavOpen(false);
      }
    };
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, []);
  return (
    <Navbar className="mx-auto max-w-screen-xl p-2 lg:rounded-full lg:pl-6">
      <div className="relative mx-auto flex items-center justify-between text-blue-gray-900">
        <h1 className="mr-4 ml-2 py-1.5 font-bold">SOCIAL-NETWORK</h1>
        <div className="hidden lg:block">
          <NavList />
        </div>
        <div className="flex gap-4">
          <IconButton
            size="sm"
            color="blue-gray"
            variant="text"
            onClick={toggleIsNavOpen}
            className="ml-auto mr-2 lg:hidden"
          >
            <Bars2Icon className="h-6 w-6" />
          </IconButton>
          <div className="flex">
            <MessageNotify nbr={3} msgNotificationItems={defaultItems} />
            <Notification
              nbr={notifications?.length}
              msgNotificationItems={notifications}
              // onRespond={handleRespond}
            />
          </div>
          <div className="h-8 my-auto flex items-center justify-center">
            <button className="Btn" onClick={Logout}>
              <div className="sign">
                <svg viewBox="0 0 512 512">
                  <path
                    d="M377.9 105.9L500.7 228.7c7.2 7.2 11.3 17.1 11.3 27.3s-4.1 20.1-11.3 27.3L377.9 406.1c-6.4 6.4-15 9.9-24 9.9c-18.7 0-33.9-15.2-33.9-33.9l0-62.1-128 0c-17.7 0-32-14.3-32-32l0-64c0-17.7 14.3-32 32-32l128 0 0-62.1c0-18.7 15.2-33.9 33.9-33.9c9 0 17.6 3.6 24 9.9zM160 96L96 96c-17.7 0-32 14.3-32 32l0 256c0 17.7 14.3 32 32 32l64 0c17.7 0 32 14.3 32 32s-14.3 32-32 32l-64 0c-53 0-96-43-96-96L0 128C0 75 43 32 96 32l64 0c17.7 0 32 14.3 32 32s-14.3 32-32 32z"
                  ></path>
                </svg>
              </div>
              <div className="text">Logout</div>
            </button>
          </div>
        </div>
      </div>
      <Collapse open={isNavOpen} className="overflow-scroll">
        <NavList />
      </Collapse>
    </Navbar>
  );
}
export default ComplexNavbar;