"use client"
import styles1 from '../../styles/siderbar_right.module.css'
import React, { useEffect, useState } from 'react'
import styles2 from '../../styles/listeChat.module.css'
import styles3 from '../../styles/event.module.css'
import {Card,} from "@material-tailwind/react";
import { UserSuggest } from '../leftSide'

export default function RightSide({propsRight}) {

  return (
    <Card className="h-[calc(90vh-2rem)] w-[20rem] p-4 shadow-xl shadow-blue-gray-900/5 ">
        <div className={styles1.infoEvent}>
            <span className={styles1.title}>Events</span>             
            <div className={styles1.listEvents}>
                <Event/>
                <Event/>
            </div>
        </div>
        <div className={styles1.message}>
            <span className={styles1.title}>Chat Messages</span>
            <div className={styles1.listeChats}>
                <ul>
                    {propsRight.Connections?.map(({Id, Avatar, FirstName},index)=>(
                        <li key={Id}>
                            <UserSuggest avatar={Avatar} name={FirstName} id={Id} isFollow={"chat"} Login={propsRight.Profil.Id}/>
                        </li>
                    ))}
                </ul>
            </div>
        </div>
    </Card>
  )
}

function Event(){
    return (
        <div className={styles3.event}>
            <span className={styles3.span1}>
                <span className={styles3.span2}>18</span>
                March
            </span>
            <div className={styles3.description}>
                <span>Codding Talk</span>
                <span className={styles3.by}>By Nasser Nourdine</span>
                <span className={styles3.hours}>11:30</span>
            </div>
        </div>
    )
}
