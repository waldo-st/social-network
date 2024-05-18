"use client"
import { Button, Label, Modal, TextInput, Select, Textarea, FileInput} from "flowbite-react";
import { useState } from "react";
import "../../globals.css"
import { SelectMembers } from "../inputSelect";
import { pushData } from "@/app/utils/fetch";

export default function CreatePostPopup({groupeId, parentCallback}) {
    const [openModal, setOpenModal] = useState(false)
    const [title, setTitle] = useState('')
    const [img, setImg] = useState('')
    const [content, setContent] = useState('')
    const [friends, setFriends] = useState([])
    const [selectedOption, setSelectedOption] = useState('public');
    function onCloseModal() {
        setOpenModal(false);
        setTitle('');
        setContent('');
        setSelectedOption('public');
        setImg('')
        setFriends([])
    }

    const handleChildData = (friendSelect) => {
        setFriends(friendSelect)
    }
    
    const handleSelectChange = (event) => {
        setSelectedOption(event.target.value);
    };

    function managerData() {
        if (title != "" && content != "") {
            const postData = {
                title: `${title}`,
                privacy: `${selectedOption}`,
                content: `${content}`,
                image: `${img}`,
                friends:`${friends}`||""
            }
            pushData(postData,`/post?groupId=${groupeId}`)
            .then(data =>{
                parentCallback(data)
                onCloseModal()
            })
            .catch(error => console.log("ERROR create post", error))
        }
    }
    return (
        <>
            <Button onClick={() => { setOpenModal(true) }} className="bg-indigo-600">Create Post</Button>
            <Modal show={openModal} size="md" onClose={onCloseModal} popup>
                <Modal.Header className="flex justify-center">New Post</Modal.Header>
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
                        {groupeId ==='0'? (<div className="mb-2 block">
                            <Label htmlFor="status" value="Selct Status" />
                            <Select id="countries" required onChange={handleSelectChange} value={selectedOption}>
                                <option value="public">public</option>
                                <option value="private">private</option>
                                <option value="almostPrivate">almost private</option>
                            </Select>
                        </div>):""}
                        {selectedOption == "almostPrivate" ? <SelectMembers forAddFriends={handleChildData} domaine={"identify friends"}/> : ""}
                    </div>
                    <div className="mb-2 block">
                        <Label htmlFor="text" value="Your Content" />
                        <Textarea
                            id="content"
                            placeholder="ex:c'est fun le codage"
                            value={content}
                            onChange={(e) => { setContent(e.target.value) }}
                            required />
                    </div>
                    <div className="flex justify-between">
                        <div className="mb-2 block">
                            <Label htmlFor="file-opload" value="Join File" />
                            <FileInput
                                id="file-opload"
                                value={img}
                                onChange={(e) => { setImg(e.target.value) }}
                                required />
                        </div>
                        <Button onClick={managerData}>Create</Button>
                    </div>
                    <span>error</span>
                </Modal.Body>
            </Modal>
        </>
    )
}
