package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Product struct {
	ID          int
	Name        string
	Price       float64
	OldPrice    float64
	Description string
	Image       string
	HasDiscount bool
	Sizes       []string
}

var products = []Product{
	{1, "PALACE PULL A CAPUCHE UNISEXE CHASSEUR", 145, 0, "Un pull chaud et stylé pour affronter l'hiver.", "/img/products/16A.webp", false, []string{"S", "M", "L"}},
	{2, "PALACE PULL A CAPUCHON MARINE", 138, 0, "Confectionné dans un tissu doux et confortable pour l’hiver.", "/img/products/18A.webp", false, []string{"S", "XS", "M"}},
	{3, "PALACE PULL CREW PASSEPOSE NOIR", 128, 145, "Un sweat classique et sobre, parfait pour un look streetwear.", "/img/products/19A.webp", true, []string{"S", "M", "L"}},
	{4, "PALACE WASHED TERRY 1/4 PLACKET HOOD MOJITO", 165, 0, "Un hoodie léger et confortable.", "/img/products/21A.webp", false, []string{"M", "L"}},
	{5, "PALACE PANTALON BOSSY JEAN STONE", 125, 0, "Un jean au look vintage et décontracté.", "/img/products/22A.webp", false, []string{"S", "M", "L"}},
	{6, "PALACE PANTALON CARGO GORE‑TEX R‑TEX NOIR", 110, 0, "Un pantalon technique et résistant.", "/img/products/33B.webp", false, []string{"S", "M"}},
}

var tmpl *template.Template

func init() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Répertoire de travail actuel :", cwd)


	tmpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))


	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))

	// mes Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/product", productHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/addProduct", addProductHandler)

	log.Println("Serveur démarré sur :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", products)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	for _, p := range products {
		if p.ID == id {
			err := tmpl.ExecuteTemplate(w, "product.html", p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.NotFound(w, r)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "add.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	desc := r.FormValue("description")
	image := r.FormValue("image")

	newProduct := Product{
		ID:          len(products) + 1,
		Name:        name,
		Price:       price,
		Description: desc,
		Image:       image,
		HasDiscount: false,
		Sizes:       []string{"S", "M", "L"},
	}
	products = append(products, newProduct)
	http.Redirect(w, r, "/product?id="+strconv.Itoa(newProduct.ID), http.StatusSeeOther)
}
