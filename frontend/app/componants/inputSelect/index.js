"use client"
import { Button, Label, Modal,Card, Checkbox, ModalFooter } from "flowbite-react";
import { list } from "postcss";
import { useState } from "react";
import { IoPersonAdd } from "react-icons/io5";

export function SelectMembers({ forAddFriends, domaine, heigth="" , listUsers }) {
    const [openModal, setOpenModal] = useState(false)
    const [member, setMember] = useState([])
    function onCloseModal() {
        setOpenModal(false);
        setMember([])
    }
    const listes = listUsers?.map(objet =>{
        return {id:objet.Id, username: objet.FirstName}
    })
   
    function handleCheckboxChange(index, isChecked) {
        if (isChecked) {
            setMember(prevMember => [...prevMember, listes[index]]);
        } else {
            setMember(prevMember => prevMember.filter(idUser => idUser !== listes[index]));
        }
    }

    function managerData() {
        const memberSelected = member
        forAddFriends(memberSelected.filter(idUser => idUser != 0))
        onCloseModal()
    }
    return (
        <>
            <Card className={` text-indigo-600 bg-white ${heigth} w-1/3 cursor-pointer`} onClick={() => {
                setOpenModal(true)
            }}>
                <div className="flex flex-col items-center">
                    <IoPersonAdd className="w-6 h-6" />
                    <span>{domaine}</span>
                </div>
            </Card>
            <Modal show={openModal} size="md" onClose={onCloseModal} popup>
                <Modal.Header className="flex justify-center">Add Friends</Modal.Header>
                <Modal.Body className="flex items-center p-3 m-0 border-b border-t ">
                    <ul className="flex flex-col h-44 w-full gap-2">
                        {listes?.map(({id, username}, index)=>(
                            <div key={index+id} className="flex items-center justify-between w-full pl-3 pr-3 p-1 shadow-md rounded-md ">
                                <Label htmlFor={index} value={username} className="cursor-pointer"/>
                                <Checkbox id={id} className="justify-self-end w-3 h-3 cursor-pointer" onChange={(e) => handleCheckboxChange(index, e.target.checked)}/>
                            </div>
                        ))}
                    </ul>
                </Modal.Body>
                <ModalFooter className="flex justify-end m-0 p-3">
                    <Button onClick={managerData}>Terminer</Button>
                </ModalFooter>
            </Modal >
        </>
    )
}
