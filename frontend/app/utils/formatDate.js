export const formatDate=(data)=>{
    let date = new Date(data);
    let formattedDate = date.toLocaleDateString('fr-FR', {day: '2-digit', month: '2-digit', year: 'numeric'});
    return formattedDate
}