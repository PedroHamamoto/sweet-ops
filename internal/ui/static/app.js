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
function categoryForm() {
    return {
        name: "",
        loading: false,
        errorMessage: "",

        openModal() {
            const modal = document.getElementById('categoryModal');
            M.Modal.getInstance(modal).open();
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
                location.reload();
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