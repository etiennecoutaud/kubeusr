package printer

import (
	"fmt"
	"os"
	"text/tabwriter"

	v1 "k8s.io/api/core/v1"
)

// PrintUsers print array with all users
func PrintUsers(users []v1.ServiceAccount) {

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAMESPACE\tUSER\tSCOPE\tCREATION DATE\t")

	for _, usr := range users {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\t\n", usr.Namespace, usr.Name, usr.Labels["scope"], usr.CreationTimestamp)
	}

	fmt.Fprint(w)
	w.Flush()
}
