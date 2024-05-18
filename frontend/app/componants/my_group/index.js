import { ChatRoom } from "../chat";
import { useState, useRef, useEffect } from 'react'
import { Card, Modal, Label, TextInput, Datepicker, Button, Checkbox } from 'flowbite-react'
import { MdEventRepeat } from "react-icons/md";
import { HiOutlineViewGridAdd } from "react-icons/hi";
import { SelectMembers } from "../inputSelect";
import { fetchData, pushData } from "../../utils/fetch";
// import { Event } from "../rightSide/index";
import { formatDate } from "../allPosts";
import { useSearchParams } from "next/navigation";
import { useHomeData } from "@/app/api/context/HomeDataContext";

export function GroupeItem({ groupeId }) {
    const [user, setUser] = useState([])
    const [followers, setFollowers] = useState([])
    const searchParams = useSearchParams()
    const search = searchParams.get('group')
    const userLogin = useHomeData().Profil.Id
    
    useEffect(()=>{
        if (user.length!==0){
            pushData(user, `/joinGroup/request?gid=${groupeId}&type=invite`)
            .then(data =>setUser([]))
            .catch(error=>console.log("Erreur invite GROUP => ",error))
        }
    },[user])

    useEffect(()=>{
        fetchData(`/followers?id=${userLogin}`)
        .then(data => {
          setFollowers(data)
        })
        .catch(erreur => console.log(erreur))
    },[])

    return (
        <div className={`flex w-full`}>
            <SelectMembers heigth={"max-h-14"} forAddFriends={setUser} domaine={"invite friends"} listUsers={followers}/>
            <ChatRoom title={"Chat Room"} room={search} type={"chatGroup"} style={"w-1/3 max-h-14 bg-white"} id={`${groupeId}`} idLogin={userLogin}/>
            <ListEvents idGroup={`${groupeId}`} />
        </div>
    )
}
export const ListEvents = ({ idGroup }) => {
    const [openModal, setOpenModal] = useState(false)
    const [listEvents, setListEvents] = useState([])
    const onCloseModal = () => {
        setOpenModal(false)
        setListEvents([])
    }
    // useEffect(() => {
    //     fetchData(`/event?groupId=${idGroup}`)
    //         .then((response) => {
    //             setListEvents(response)
    //         });
    // }, [listEvents])
    return (
        <>
            <Card className="w-1/3 text-indigo-600 bg-white max-h-14" onClick={() => {
                setOpenModal(true)
            }}>
                <div className="flex flex-col items-center cursor-pointer">
                    <MdEventRepeat className="h-6 w-52" />
                    <span>Events</span>
                </div>
            </Card>
            <Modal show={openModal} size="md" onClose={onCloseModal} popup>
                <Modal.Header className="flex justify-center">List Events</Modal.Header>
                <Modal.Body>
                    <div className='flex justify-end'>
                        <EventPopup idGroup={idGroup} list={setListEvents} />
                    </div>
                    <div>
                        {/* {
                            listEvents?.map((event, i) => (
                                <Event key={i} {...event} />
                            ))
                        } */}
                    </div>
                </Modal.Body>
            </Modal>
        </>
    )
}
export const EventPopup = ({ idGroup, list }) => {
    const [openModal, setOpenModal] = useState(false)
    const [title, setTitle] = useState("")
    const [description, setDescription] = useState("")
    const [date, setDate] = useState(new Date())
    const [time, setTime] = useState('')
    const [goingCheckbox, setGoingCheckbox] = useState(false);
    const [notGoingCheckbox, setNotGoingCheckbox] = useState(false);
    let dateAreSelected = useRef(false)
    const [error, setError] = useState('')
    const onCloseModal = () => {
        setOpenModal(false)
        dateAreSelected.current = false
        setTitle('')
        setDescription('')
        setTime('')
        setDate(new Date())
        setGoingCheckbox(false)
        setNotGoingCheckbox(false)
        setError('')
    }
    function SelectedDate(date) {
        setDate(date)
        dateAreSelected.current = true
    }
    const managerData = () => {
        let option = [false, false]
        if (goingCheckbox) {
            option[0] = true
        }
        if (notGoingCheckbox) {
            option[1] = true
        }
        const [heures, minutes] = time.split(':')
        const data = {
            groupId: idGroup,
            title: title,
            description: description,
            date: date.setHours(parseInt(heures, 10), parseInt(minutes, 10)),
            option: option
        }
        if (title === '' || description === '' || data == NaN) {
            setError(`error un champs n'est pas remplit`)
        } else {
            // pushData(data, `/event?groupId=${idGroup}`)
            //     .then((resp) => {
            //         list(resp)
            //         onCloseModal()
            //     })
            //     .catch((error) => {
            //     })
        }
        // console.log(data);
    }
    return (
        <>
            <Button className="w-1/3 bg-indigo-600 text-white max-h-14" onClick={() => {
                setOpenModal(true)
            }}>
                <div className="flex flex-col items-center cursor-pointer">
                    <HiOutlineViewGridAdd className="h-6 w-52" />
                    <span>Create Event</span>
                </div>
            </Button>
            <Modal show={openModal} size="md" onClose={onCloseModal} popup>
                <Modal.Header className="flex justify-center">New Event</Modal.Header>
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
                    <div className="flex flex-col gap-2">
                        <div className='flex justify-between'>
                            <div className='w-2/3'>
                                <Label htmlFor="Data" value="Select Date:" />
                                <Datepicker
                                    minDate={new Date(2024, 0, 1)}
                                    maxDate={new Date(2050, 3, 30)}
                                    autoHide={true}
                                    value={date}
                                    onSelectedDateChanged={SelectedDate}
                                    required />
                            </div>
                            {dateAreSelected.current ? (
                                <div>
                                    <Label htmlFor="Time" value="Select Hours:" />
                                    <input type="time" id="time" className="bg-gray-50 border leading-none border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500" min="09:00" max="18:00" value={time} onChange={(e) => { setTime(e.target.value) }} required />
                                </div>
                            ) : ""}
                        </div>
                        <div className='flex gap-10'>
                            <span>Select options</span>
                            <div className='flex gap-2 items-center'>
                                <Label htmlFor="going" value="Going:" />
                                <Checkbox id='going'
                                    name="going"
                                    checked={goingCheckbox}
                                    onChange={(e) => {
                                        setGoingCheckbox(e.target.checked)
                                    }}
                                />
                            </div>
                            <div className='flex gap-2 items-center'>
                                <Label htmlFor="not_going" value="Not Going:" />
                                <Checkbox id='not_going'
                                    name="not_going"
                                    checked={notGoingCheckbox}
                                    onChange={(e) => {
                                        setNotGoingCheckbox(e.target.checked)
                                    }}
                                />
                            </div>
                        </div>
                        <span className="text-red-500 self-center">{error}</span>
                        <Button className='w-1/2 ml-20' onClick={managerData}>Create</Button>
                    </div>
                </Modal.Body>
            </Modal>
        </>
    )
}