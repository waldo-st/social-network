"use client"
import CreatePostPopup from "@/app/componants/createPostPopup";
import Profil from "@/app/componants/profil";
import { useState } from "react";

export default function Profile() {
    const [childData, setChildData] = useState(null);
    const handleCallback = (data) => {
        setChildData(data);
    }
    return (
        <div className="flex flex-col w-full mx-8 items-center gap-2">
            <div className="flex justify-end items-center border-t-4 border-indigo-500 p-2 w-full bg-white">
                <CreatePostPopup groupeId={`0`} parentCallback={handleCallback}/>
            </div>
            <Profil addPost={childData}/>
        </div>
    );
}