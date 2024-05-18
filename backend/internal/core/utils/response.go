package util

// func SendResponse(res http.ResponseWriter, data interface{}, code int) {
// 	res.Header().Set("Content-Type", "application/json")
// 	res.WriteHeader(code)
// 	if err := json.NewEncoder(res).Encode(data); err != nil {
// 		log.Println("❌ Error encoding json response")
// 	}
// }
// func HandleError(res http.ResponseWriter, message string, code int) {
// 	log.Println("❌ ERROR: ", message)
// 	response := map[string]string{"message": message}
// 	SendResponse(res, response, code)
// }

// func BindJson(r *http.Request, data interface{}) error {
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(data); err != nil {
// 		return err
// 	}
// 	return nil
// }
