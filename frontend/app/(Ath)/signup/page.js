"use client"
import Link from "next/link"
import { useRouter } from 'next/navigation';
import { useRef, useState } from "react"
import checkChampForm from "@/app/utils/checkField"
import { pushData } from "@/app/utils/fetch";
import { handleImageUpload } from "@/app/utils/upload_Img";

// import { POST } from "../../api/upload";

export const ENDPOINT ="http://localhost:8080"


export default function Register(){
    const firstName = useRef()
    const lastName = useRef()
    const nickname = useRef()
    const dateOfBirth = useRef()
    const email = useRef()
    const password = useRef()
    const avatar = useRef()
    const  about = useRef()
    const [checkValid,setCheckValid] = useState("")
    const [selectedImage, setSelectedImage] = useState("");
    const [error, setError] = useState("")
    
    const route = useRouter()

   
    
    const handleSubmit = e =>{
        e.preventDefault()
        
        const newForm = {
            firstName: firstName.current.value,
            lastName: lastName.current.value,
            nickname: nickname.current.value,
            dateOfBirth: new Date(dateOfBirth.current.value),
            email: email.current.value,
            password: password.current.value,
            avatar: selectedImage !== "" ? selectedImage : "default.png",
            about: about.current.value,
        }

        setCheckValid(checkChampForm(newForm))
        
        if (checkValid===""){
            pushData(newForm, '/register')
            .then(data => {
                if (data.status === 201){
                    route.push('/login');
                }
            })
            .catch(error => {
                setError(error.data.message);
            });
            
        }else{
            setError(checkValid)
        }
    }
    return (
    <div className="grid place-items-center h-screen">
        <div className="shadow-lg p-5 rounded-lg border-t-4 border-blue-400">
            <h1 className=" text-xl font-bold my-4">Enter the details</h1>
            <form onSubmit={handleSubmit} className="flex flex-col gap-2">
                <input className="w-[400px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={firstName} type="text" placeholder="First Name" name="firstName" required />
                <div className="flex gap-1">
                    <input className="w-[200px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={lastName} type="text" placeholder="Last Name" name="lastName" required />
                    <input className="w-[197px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={nickname} type="text" placeholder="Nickname" name="nickname" />
                </div>
                <input className="w-[400px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={dateOfBirth} type="date" placeholder="Date of Birth" name="dateOfBirth" required />
                <input className="w-[400px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={email} type="email" placeholder="Email" name="email" required />
                <input className="w-[400px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={password} type="password" placeholder="Password" name="password" required />
 
                <input className="w-[400px] border border-gray-200 py-2 px-6 bg-zinc-100/40" ref={avatar} type="file" onChange={(e)=>{
                    handleImageUpload(e)
                    .then(data=>{
                        if (data.status===201){
                            setSelectedImage(data.Message)
                        }
                    }) 
                    .catch(error=>console.log(error))                   
                }} placeholder="Avatar / Image" name="avatar" accept="image/*"/>
                <textarea ref={about} placeholder="About Me" name="aboutMe" />
                <button className="bg-blue-600 text-white font-bold cursor-pointer px-6 py-2">Register</button>
                {/* Afficher un message d'erreur ici si nécessaire */}
                {error && <div className="bg-red-500 text-white w-fit text-sm py-1 px-3 rounded-md mt-2">{error}</div>}
                 {/* <div className="bg-red-500 text-white w-fit text-sm py-1 px-3 rounded-md mt-2">Error message</div> */}
                <div className="text-sm mt-3 text-right">
                    Already have an account?{' '}
                    <Link href="/login">
                    <span className="underline">Login</span>
                    </Link>
                </div>
            </form>
        </div>
    </div>
    )
}