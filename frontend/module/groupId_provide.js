import React, {useState,createContext} from "react"

const GroupIdContext = createContext(
    {
        idUser:0,
        setIdUser:()=>{},
    }
);

const GroupIdProvider =({children})=>{
    
    const [idUser, setIdUser] = useState(0)
    
    return (
        <GroupIdContext.Provider value={{idUser:idUser, setIdUser:setIdUser}}>
            {children}
        </GroupIdContext.Provider>
    )
}

export {GroupIdContext}
export default GroupIdProvider