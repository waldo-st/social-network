const axios = require('axios');

const ENDPOINTE ="http://localhost:8080"

export async function fetchData(path){
    try{
  
      const token = getCookie('token');
      
      const res = await axios.get(ENDPOINTE + path, {
        headers: {
          "Authorization": "Bearer " + token
        }
      });
      return res.data;
    } catch(error) {
      if (error.response.status === 401){
        deleteCookie("token")
      }
      console.error("Error fetch data: "+error.response);
      throw error.response;
    }
}

// Function to get the cookie value by name
export function getCookie(name) {
  let cookieValue = null;
  if (document.cookie && document.cookie !== '') {
    const cookies = document.cookie.split(';');
    for (let i = 0; i < cookies.length; i++) {
      const cookie = cookies[i].trim();
      if (cookie.substring(0, name.length + 1) === (name + '=')) {
        cookieValue = decodeURIComponent(cookie.substring(name.length + 1));
        break;
      }
    }
  }
  return cookieValue;
}

// function delete cookie
function deleteCookie(name) {
  document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

export async function pushData(payload, path) {
  try {
    const token = getCookie('token');
    const formatData = JSON.stringify(payload);
    const response = await axios.post(ENDPOINTE + path, formatData, {
      headers: {
        'Content-Type': 'application/json',
        "Authorization": token ? `Bearer ${token}`:""
      }
    });
  
    return response.data;
  } catch (error) {
    // deleteCookie("token")
    if (error.response.status === 401){
      deleteCookie("token")
    }
    console.error('Error pushing data: ', error.response);
    throw error.response;
  }
}
