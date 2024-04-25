package api
import(
	"net/http"
	"github.com/go-chi/render"
)


func renderStatusError(r *http.Request, w http.ResponseWriter, status int, errMsg string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]interface{}{"error": errMsg})
}