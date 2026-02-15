// login
function loginForm() {
    return {
        email: "",
        password: "",
        loading: false,
        errorMessage: "",

        async submit() {
            this.errorMessage = "";
            this.loading = true;

            try {
                const res = await fetch("/api/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({
                        email: this.email,
                        password: this.password,
                    }),
                });

                if (!res.ok) {
                    const text = await res.text();
                    throw new Error(text.trim() || "Login failed");
                }

                const data = await res.json();

                localStorage.setItem("access_token", data.access_token);
                localStorage.setItem("refresh_token", data.refresh_token);

                document.cookie =
                    "access_token=" +
                    encodeURIComponent(data.access_token) +
                    "; path=/; SameSite=Strict";

                window.location.href = "/home";
            } catch (err) {
                this.errorMessage = err.message;
            } finally {
                this.loading = false;
            }
        },
    };
}

// Categories
function categoriesPage() {
    return {
        // Table state
        categories: [],
        page: 1,
        pageSize: 10,
        totalPages: 1,
        totalItems: 0,
        loadingTable: false,

        // Form state
        name: "",
        loading: false,
        errorMessage: "",

        openModal() {
            const modal = document.getElementById('categoryModal');
            M.Modal.getInstance(modal).open();
        },

        async fetchCategories() {
            this.loadingTable = true;
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/categories?page=${this.page}&page_size=${this.pageSize}`, {
                    headers: { "Authorization": "Bearer " + token },
                });

                if (!res.ok) {
                    throw new Error("Falha ao carregar categorias");
                }

                const data = await res.json();
                this.categories = data.data || [];
                this.page = data.page;
                this.pageSize = data.page_size;
                this.totalPages = data.total_pages;
                this.totalItems = data.total_items;
            } catch (err) {
                console.error(err);
            } finally {
                this.loadingTable = false;
            }
        },

        goToPage(p) {
            if (p < 1 || p > this.totalPages) return;
            this.page = p;
            this.fetchCategories();
        },

        paginationRange() {
            const range = [];
            const maxVisible = 5;
            let start = Math.max(1, this.page - Math.floor(maxVisible / 2));
            let end = start + maxVisible - 1;

            if (end > this.totalPages) {
                end = this.totalPages;
                start = Math.max(1, end - maxVisible + 1);
            }

            for (let i = start; i <= end; i++) {
                range.push(i);
            }
            return range;
        },

        async submit() {
            this.errorMessage = "";
            this.loading = true;

            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch("/api/categories", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + token,
                    },
                    body: JSON.stringify({ name: this.name }),
                });

                if (!res.ok) {
                    const text = await res.text();
                    throw new Error(text.trim() || "Falha ao criar categoria");
                }

                const modal = document.getElementById('categoryModal');
                M.Modal.getInstance(modal).close();
                this.name = "";
                this.page = 1;
                await this.fetchCategories();
                M.toast({ html: 'Categoria criada com sucesso!', classes: 'green' });
            } catch (err) {
                this.errorMessage = err.message;
            } finally {
                this.loading = false;
            }
        }
    };
}

// Products
function productsPage() {
    return {
        // Categories for dropdown
        categories: [],
        filteredCategories: [],
        categorySearch: "",
        categoryId: "",
        showDropdown: false,

        // Form state
        flavor: "",
        productionPrice: "",
        sellingPrice: "",
        loading: false,
        errorMessage: "",

        openModal() {
            const modal = document.getElementById('productModal');
            const instance = M.Modal.getInstance(modal);
            instance.open();
        },

        filterCategories() {
            const search = this.categorySearch.toLowerCase();
            if (search === "") {
                this.filteredCategories = this.categories;
            } else {
                this.filteredCategories = this.categories.filter(cat =>
                    cat.name.toLowerCase().startsWith(search)
                );
            }
            this.showDropdown = true;
        },

        selectCategory(category) {
            this.categoryId = category.id;
            this.categorySearch = category.name;
            this.showDropdown = false;
        },

        async fetchCategories() {
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/categories?page=1&page_size=100`, {
                    headers: { "Authorization": "Bearer " + token },
                });

                if (!res.ok) {
                    throw new Error("Falha ao carregar categorias");
                }

                const data = await res.json();
                this.categories = data.data || [];
                this.filteredCategories = this.categories;
            } catch (err) {
                console.error(err);
            }
        },

        async submit() {
            this.errorMessage = "";
            this.loading = true;

            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch("/api/products", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + token,
                    },
                    body: JSON.stringify({
                        category_id: this.categoryId,
                        flavor: this.flavor,
                        production_price: parseFloat(this.productionPrice),
                        selling_price: parseFloat(this.sellingPrice),
                    }),
                });

                if (!res.ok) {
                    const text = await res.text();
                    throw new Error(text.trim() || "Falha ao criar produto");
                }

                const modal = document.getElementById('productModal');
                M.Modal.getInstance(modal).close();

                // Reset form
                this.categoryId = "";
                this.categorySearch = "";
                this.flavor = "";
                this.productionPrice = "";
                this.sellingPrice = "";
                this.filteredCategories = this.categories;

                M.toast({ html: 'Produto criado com sucesso!', classes: 'green' });
            } catch (err) {
                this.errorMessage = err.message;
            } finally {
                this.loading = false;
            }
        }
    };
}

document.addEventListener('DOMContentLoaded', function () {
    initModal('categoryModal');
    initModal('productModal');
});

// Common
function initModal(id, options = {}) {
    let el = document.getElementById(id);
    if (el) {
        M.Modal.init(el, {
            inDuration: 0,
            outDuration: 0,
            onCloseEnd: function (modal) {
                modal.style.top = '';
                modal.style.opacity = '';
                modal.style.display = 'none';
            },
            ...options,
        });
    }
}