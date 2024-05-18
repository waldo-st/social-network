"use client"
import React, { useState, useEffect } from 'react';
import { Avatar } from "@material-tailwind/react";
import { Button, Label, Modal, TextInput, Select, Checkbox, Textarea, FileInput, ModalFooter } from "flowbite-react";
import { fetchData, pushData } from '@/app/utils/fetch';

export function formatDate(dateString) {
  let date = new Date(dateString);
  let options = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit', second: '2-digit' };
  let formattedDate = date.toLocaleDateString('fr-FR', options);
  return formattedDate
}

const Post = ({ Username, CreatedAt, Content, Image, NbrComments, Title , UserAvatar, Id}) => {
  let date = formatDate(CreatedAt)
  const [showPopup, setShowPopup] = useState(false);
  const togglePopup = () => {
    setShowPopup(!showPopup);
  };

  // Check if the avatar has a valid extension (.jpg or .png)
  const isValidExtension = UserAvatar && (UserAvatar.endsWith('.jpg') || UserAvatar.endsWith('.png') || UserAvatar.endsWith('.jpeg'));
  const avat = isValidExtension ? `/uploads/${UserAvatar}` : '/uploads/default.png';
  return (
    <div className="mx-auto w-full border rounded-lg shadow-lg bg-white">
      <div className='m-3'>
        <div className="flex items-center mb-4 gap-3">
          <Avatar src={avat} alt="avatar" />
          <div>
            <h2 className="text-lg font-semibold">{Username}</h2>
            <p className="text-gray-600">{date}</p>
          </div>
        </div>
        {/* <hr></hr> */}
        <div className='pt-2'>
          <h3 className='text-lg'>{Title}</h3>
          <div className="text-gray-700 mb-4 whitespace-pre-line">{Content}</div>
          {/* <div> <Image>{Image}</Image> </div> */}
        </div>
        {/* <hr></hr> */}
        <div className='flex justify-end pr-9 pt-2'>
          <div className="flex items-center justify-between mb-2">
            <div className="flex items-center">
              <button className="flex items-center text-blue-500" onClick={togglePopup}>
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-6 h-6">
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8.625 9.75a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H8.25m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0H12m4.125 0a.375.375 0 1 1-.75 0 .375.375 0 0 1 .75 0Zm0 0h-.375m-13.5 3.01c0 1.6 1.123 2.994 2.707 3.227 1.087.16 2.185.283 3.293.369V21l4.184-4.183a1.14 1.14 0 0 1 .778-.332 48.294 48.294 0 0 0 5.83-.498c1.585-.233 2.708-1.626 2.708-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
                </svg>
                Comment({NbrComments})
              </button>
            </div>
        </div>
        </div>
        {showPopup && <Comments postId={Id}/>}
      </div>
    </div>
  );
};

const PostsPage = ({propsPosts}) => {
  return (
    <div className='posts grid gap-2 grid-cols-1 w-full overflow-y-auto hide-scrollbar'>
        <ul>
          {
            propsPosts?.map((post, index) => (
              <li key={index} >
                <Post {...post} />
              </li>
            ))
          }
        </ul>
    </div>
  );
};

const Comments = ({postId}) => {
  const [newComment, setNewComment] = useState('');
  const [comments, setComments] = useState([]);
  const [image, setImage] = useState('');
  const [content, setContent] = useState('');
  useEffect(() => {
    
    fetchData("/comment")
    .then((data)=>{

    })
    .catch(error => console.log("Error fecth comment => ",error))
      
  },[]);
  const handleSubmitComment = async () => {
    const dataComment = {
      comment: newComment,
      image: `${image}`
    }
    pushData(dataComment,`/comment?postId=${postId}`)
    // setNewComment('');
    // setImage('');
  };
  
  return (
    <div className='p-3 bg-blue-100 rounded-xl '>
      <hr/>
      <h3 className='text-black mb-4 font-extrabold'>Commentaires :</h3>
      <ul>
        {comments.map((comment) => (
          <li key={comment.id} className="text-gray-600 mb-2">
            <div className="flex items-center mb-2 gap-1">
              <Avatar src={comment.avatar} alt="avatar" />
              <div>
                <h2 className="text-lg font-semibold ">{comment.user}</h2>
                <h4 className=' font-light'>{formatDate(comment.date)}</h4>
              </div>
            </div>
            <div className='ml-12'>{comment.text}</div>
          </li>
        ))}
      </ul>
      <div className="mb-4 flex items-center"> {/* Utilisation de flexbox */}
        <input
          type="content"
          value={newComment}
          onChange={(e) => { setNewComment(e.target.value) }}
          className="flex-1 border rounded-lg py-2 px-4 mr-2" // Utilisation de flex-1 pour que le champ de texte prenne tout l'espace disponible
          placeholder="Ajouter un commentaire..."
        />
        <div className="mb-4 block w-20"> {/* Utilisation de classes tailwindcss pour définir la largeur */}
          <Label className=' font-bold' htmlFor="file-opload" value="Join File" />
          <FileInput
            id="file-opload"
            value={image}
            onChange={(e) => { setImage(e.target.value) }}
            required
          />
        </div>
        <button
          onClick={handleSubmitComment}
          className="mt-2 bg-blue-500 text-white py-2 px-4 rounded-lg hover:bg-blue-600"
        >
          Ajouter
        </button>
      </div>
    </div>
  );
};


export default PostsPage;
