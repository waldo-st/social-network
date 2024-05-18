export default function checkChampForm(form) {
    const regexfirstname = /^[a-zA-Z ]{2,15}$/g;
    const regexlastname = /^[a-zA-Z]{2,10}$/g;
    const regexEmail = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/g;
    const regexPassword = /^[^\0]{4,}$/g;
  
    if (regexfirstname.test(form.firstName) === false) {
      return `FirstName is not correct...`;
    }
    if (regexlastname.test(form.lastName) === false) {
      return `LastName is not correct...`;
    }
    if (regexEmail.test(form.email) === false) {
      return `Email is not correct...`;
    }
    if (regexPassword.test(form.password) === false) {
      return `Password must contain more than 4 characters...`;
    }
    return "";
  }
  