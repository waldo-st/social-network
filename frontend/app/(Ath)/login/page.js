"use client"
import Link from "next/link"
import { useRef, useState } from "react"
import { useRouter } from 'next/navigation';
import { pushData } from "@/app/utils/fetch";

export default function Login() {
    const email = useRef()
    const password = useRef()
    const [error, setError] = useState("")

    const route = useRouter()
    const handleSubmit = (e) => {
        e.preventDefault()
        const newForm = {
            login: email.current.value,
            password: password.current.value,
        }
        pushData(newForm, '/login')
        .then(tokenData => {
            const daysValid = 1;
            const expiryDate = new Date(Date.now() + daysValid * 24 * 60 * 60 * 1000).toUTCString();
            document.cookie = `token=${tokenData}; expires=${expiryDate}; path=/; samesite=strict`;
            route.push("/")
        })
        .catch(error => {
            console.log(error.status)
            setError(error.data.message);
        });
    }
    return (
        <div className="grid place-items-center h-screen">
            <div className="shadow-lg p-5 rounded-lg border-t-4 border-blue-400">
                <h1 className=" text-xl font-bold my-4">Enter the details</h1>
                <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                    <input ref={email} type="text" placeholder="Email or Username" />
                    <input ref={password} type="password" placeholder="password" />
                    <button className=" bg-blue-600 text-white font-bold cursor-pointer px-6 py-2">Login</button>
                    {error && <div className=" bg-red-500 text-white w-fit text-sm py-1 px-3 rounded-md mt-2">{error}</div>}
                    <Link className=" text-sm mt-3 text-right " href={`/signup`}>Do not have an account? <span className="underline">Register</span></Link>
                </form>
            </div>
        </div>
    )
}
