package handler

import (
	"fmt"
	"net/http"

	"mengzhao/view/home"
)

func HandleHomeIndex(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	//account, err := db.GetAccountByID(r.Context(), user.ID)
	//if err != nil {
	//	return err
	//}

	fmt.Printf("%+v\n", user.Account)

	return home.Index().Render(r.Context(), w)
}
