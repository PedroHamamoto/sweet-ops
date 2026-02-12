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