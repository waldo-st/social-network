"use client"
import { GroupeItem } from "@/app/componants/my_group";
import CreatePostPopup from "@/app/componants/createPostPopup";
import PostsPage from "@/app/componants/allPosts";
import { useContext, useEffect, useState } from "react";
import { fetchData } from "@/app/utils/fetch";

export default function GroupeItemPage(context) {
    const Id = context.params.slug
    
    const [posts,setPosts] = useState([])
    useEffect(()=>{
        fetchData(`/group/posts?groupId=${Id}`)
        .then(data=>{
            setPosts(data)
        })
        .catch(error=>console.log("Error Fetch posts group => ",error))
    },[])
    
    return (
        <div className="flex flex-col gap-2 w-full mx-8 items-center ">
            <div className="flex justify-end border-t-4 border-indigo-500 pt-2 w-full">
                <CreatePostPopup groupeId={Id}/>
            </div>
            <GroupeItem groupeId={Id}/>
            <PostsPage propsPosts={posts}/>
        </div>
    )
}

export async function getStaticPaths(){
    const data = await fetchData(`/joinedGroups`)
    const paths = data.map(item => ({
        params: {slug: item.Id.toString()}
    }))
    console.log("Param => ",paths)
    return {
        paths,
        fallback: false
    }
}
