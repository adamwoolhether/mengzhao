package supabase

import (
	"errors"
	"os"

	"github.com/nedpals/supabase-go"
)

var Client *supabase.Client

func Connect() error {
	supabaseHost := os.Getenv("SUPABASE_HOST")
	if len(supabaseHost) == 0 {
		return errors.New("SUPABASE_HOST is not set")
	}
	supabaseSecret := os.Getenv("SUPABASE_SECRET")
	if len(supabaseSecret) == 0 {
		return errors.New("SUPABASE_SECRET is not set")
	}

	Client = supabase.CreateClient(supabaseHost, supabaseSecret)

	return nil
}
