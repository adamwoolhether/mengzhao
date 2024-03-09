package handler

import (
	"fmt"
	"net/http"
	"time"

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

func HandleLongProcess(w http.ResponseWriter, r *http.Request) error {
	time.Sleep(time.Second * 5)

	return home.UserLikes(1000).Render(r.Context(), w)
}
