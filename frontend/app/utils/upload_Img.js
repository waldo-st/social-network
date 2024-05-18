import axios from "axios";

export const handleImageUpload = async (event) => {
        const file = event.target.files[0];
        const formData = new FormData();
        formData.append('file', file);
    
        try {
          const response = await axios.post('/api/upload_img', formData);
          // Handle the response here
          console.log('Uploaded image:', response.data.Message);
          return response.data
          // Update state or display the uploaded image
        } catch (error) {
          console.error('Error uploading image:', error);
        }
    };