package cmds


import (
   "fmt"
   "github.com/astropatty/gbtest/auth"
)



func Deploy() {
   err := auth.CheckCredentials()
   if err != nil {
      panic(fmt.Sprintf("Unable authenticate: %s", err))
   }
   fmt.Println("Permissions valid!")
}
