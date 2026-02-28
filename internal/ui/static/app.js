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
        // Table state
        products: [],
        page: 1,
        pageSize: 10,
        totalPages: 1,
        totalItems: 0,
        loadingTable: false,

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
            console.log(modal)
            const instance = M.Modal.getInstance(modal);
            console.log(instance)
            instance.open();
        },

        async fetchProducts() {
            this.loadingTable = true;
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/products?page=${this.page}&page_size=${this.pageSize}`, {
                    headers: { "Authorization": "Bearer " + token },
                });

                if (!res.ok) {
                    throw new Error("Falha ao carregar produtos");
                }

                const data = await res.json();
                this.products = data.data || [];
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
            this.fetchProducts();
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

        formatMoney(value) {
            return new Intl.NumberFormat('pt-BR', {
                style: 'currency',
                currency: 'BRL'
            }).format(value);
        },

        formatPercentage(value) {
            return new Intl.NumberFormat('pt-BR', {
                style: 'percent',
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            }).format(value / 100);
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

                // Refresh products list
                this.page = 1;
                await this.fetchProducts();

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
    initModal('productionModal');
    initModal('saleModal');
});

// Productions
function productionsPage() {
    return {
        // Table state
        productions: [],
        page: 1,
        pageSize: 10,
        totalPages: 1,
        totalItems: 0,
        loadingTable: false,

        // Modal / form state
        products: [],
        filteredProducts: [],
        productSearch: "",
        productId: "",
        showDropdown: false,
        quantity: 1,
        loading: false,
        errorMessage: "",

        openModal() {
            const modal = document.getElementById('productionModal');
            const instance = M.Modal.getInstance(modal);
            instance.open();
        },

        async fetchProducts() {
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/products?page=1&page_size=1000`, {
                    headers: { "Authorization": "Bearer " + token },
                });

                if (!res.ok) {
                    throw new Error("Falha ao carregar produtos");
                }

                const data = await res.json();
                this.products = data.data || [];
                this.filteredProducts = this.products;
            } catch (err) {
                console.error(err);
            }
        },

        async fetchProductions() {
            this.loadingTable = true;
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/productions?page=${this.page}&page_size=${this.pageSize}`, {
                    headers: { "Authorization": "Bearer " + token },
                });

                if (!res.ok) {
                    throw new Error("Falha ao carregar produções");
                }

                const data = await res.json();
                this.productions = data.data || [];
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
            this.fetchProductions();
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

        formatDate(value) {
            return new Date(value).toLocaleString('pt-BR');
        },

        filterProducts() {
            const search = this.productSearch.toLowerCase();
            if (search === "") {
                this.filteredProducts = this.products;
            } else {
                this.filteredProducts = this.products.filter(prod =>
                    prod.flavor.toLowerCase().includes(search) ||
                    prod.category.name.toLowerCase().includes(search)
                );
            }
            this.showDropdown = true;
        },

        selectProduct(product) {
            this.productId = product.id;
            this.productSearch = product.category.name + ' ' + product.flavor;
            this.showDropdown = false;
        },

        async submit() {
            this.errorMessage = "";
            this.loading = true;

            if (!this.productId) {
                this.errorMessage = "Selecione um produto da lista";
                this.loading = false;
                return;
            }

            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/products/${this.productId}/productions`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + token,
                    },
                    body: JSON.stringify({
                        quantity: parseInt(this.quantity, 10),
                    }),
                });

                if (!res.ok) {
                    const text = await res.text();
                    throw new Error(text.trim() || "Falha ao registrar produção");
                }

                const modal = document.getElementById('productionModal');
                M.Modal.getInstance(modal).close();

                // Reset form
                this.productId = "";
                this.productSearch = "";
                this.quantity = 1;
                this.filteredProducts = this.products;

                // Refresh table
                this.page = 1;
                await this.fetchProductions();

                M.toast({ html: 'Produção registrada com sucesso!', classes: 'green' });
            } catch (err) {
                this.errorMessage = err.message;
            } finally {
                this.loading = false;
            }
        }
    };
}

// Sales
function salesPage() {
    return {
        // Table state
        sales: [],
        page: 1,
        pageSize: 10,
        totalPages: 1,
        totalItems: 0,
        loadingTable: false,

        // Products for item selection
        products: [],

        // Form state
        source: "",
        paymentMethod: "",
        selfConsumption: false,
        items: [newSaleItem()],
        total: 0,
        loading: false,
        errorMessage: "",

        openModal() {
            const modal = document.getElementById('saleModal');
            M.Modal.getInstance(modal).open();
        },

        async fetchProducts() {
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/products?page=1&page_size=1000`, {
                    headers: { "Authorization": "Bearer " + token },
                });
                if (!res.ok) throw new Error("Falha ao carregar produtos");
                const data = await res.json();
                this.products = data.data || [];
            } catch (err) {
                console.error(err);
            }
        },

        async fetchSales() {
            this.loadingTable = true;
            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch(`/api/sales?page=${this.page}&page_size=${this.pageSize}`, {
                    headers: { "Authorization": "Bearer " + token },
                });
                if (!res.ok) throw new Error("Falha ao carregar vendas");
                const data = await res.json();
                this.sales = data.data || [];
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
            this.fetchSales();
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

        formatMoney(value) {
            return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
        },

        formatDate(value) {
            return new Date(value).toLocaleString('pt-BR');
        },

        addItem() {
            this.items.push(newSaleItem());
        },

        removeItem(index) {
            this.items.splice(index, 1);
            this.calculateTotal();
        },

        filterProductsForItem(index) {
            const search = this.items[index].productSearch.toLowerCase();
            if (search === "") {
                this.items[index].filteredProducts = this.products;
            } else {
                this.items[index].filteredProducts = this.products.filter(p =>
                    (p.category.name + ' ' + p.flavor).toLowerCase().includes(search)
                );
            }
            this.items[index].showDropdown = true;
        },

        selectProductForItem(index, product) {
            this.items[index].productId = product.id;
            this.items[index].productSearch = product.category.name + ' ' + product.flavor;
            this.items[index].unitPrice = product.selling_price;
            this.items[index].showDropdown = false;
            this.calculateTotal();
        },

        onSelfConsumptionChange() {
            if (this.selfConsumption) {
                this.items.forEach(item => { item.isGift = false; });
            }
            this.calculateTotal();
        },

        calculateTotal() {
            if (this.selfConsumption) {
                this.total = 0;
                return;
            }
            let t = 0;
            for (const item of this.items) {
                if (!item.isGift && item.unitPrice && item.quantity) {
                    t += item.unitPrice * item.quantity;
                }
            }
            this.total = t;
        },

        async submit() {
            this.errorMessage = "";
            this.loading = true;

            if (!this.source) {
                this.errorMessage = "Selecione a origem da venda";
                this.loading = false;
                return;
            }
            if (!this.paymentMethod) {
                this.errorMessage = "Selecione o método de pagamento";
                this.loading = false;
                return;
            }

            const saleItems = [];
            for (const item of this.items) {
                if (!item.productId) {
                    this.errorMessage = "Selecione um produto para todos os itens";
                    this.loading = false;
                    return;
                }
                if (!item.quantity || item.quantity <= 0) {
                    this.errorMessage = "Quantidade deve ser maior que 0";
                    this.loading = false;
                    return;
                }
                saleItems.push({
                    product_id: item.productId,
                    quantity: parseInt(item.quantity, 10),
                    is_gift: item.isGift,
                });
            }

            try {
                const token = localStorage.getItem("access_token");
                const res = await fetch("/api/sales", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": "Bearer " + token,
                    },
                    body: JSON.stringify({
                        source: this.source,
                        payment_method: this.paymentMethod,
                        self_consumption: this.selfConsumption,
                        items: saleItems,
                    }),
                });

                if (!res.ok) {
                    const text = await res.text();
                    throw new Error(text.trim() || "Falha ao registrar venda");
                }

                const modal = document.getElementById('saleModal');
                M.Modal.getInstance(modal).close();

                // Reset form
                this.source = "";
                this.paymentMethod = "";
                this.selfConsumption = false;
                this.items = [newSaleItem()];
                this.total = 0;

                // Refresh table
                this.page = 1;
                await this.fetchSales();

                M.toast({ html: 'Venda registrada com sucesso!', classes: 'green' });
            } catch (err) {
                this.errorMessage = err.message;
            } finally {
                this.loading = false;
            }
        }
    };
}

function newSaleItem() {
    return {
        productId: "",
        productSearch: "",
        filteredProducts: [],
        showDropdown: false,
        quantity: 1,
        unitPrice: 0,
        isGift: false,
    };
}

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