export const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:3000"

export async function apiFetch(path: string, options?: RequestInit) {
    const token = typeof window !== "undefined" ? localStorage.getItem("token") : null

    const res = await fetch(`${API_URL}${path}`, {
        ...options,
        headers: {
            "Content-Type": "application/json",
            ...(token ? { Authorization: `Bearer ${token}` } : {}),
            ...options?.headers,
        },
    })

    if (res.status === 401) {
        if (typeof window !== "undefined") {
            localStorage.removeItem("token")
            localStorage.removeItem("user")
            window.location.href = "/login"
        }
        throw new Error("No autorizado")
    }

    return res
}

export async function apiGet<T = unknown>(path: string): Promise<T> {
    const res = await apiFetch(path)
    if (!res.ok) {
        const data = await res.json().catch(() => null)
        throw new Error(data?.error || `Error ${res.status}`)
    }
    return res.json()
}

export async function apiPost<T = unknown>(path: string, body: unknown): Promise<T> {
    const res = await apiFetch(path, {
        method: "POST",
        body: JSON.stringify(body),
    })
    if (!res.ok) {
        const data = await res.json().catch(() => null)
        throw new Error(data?.error || `Error ${res.status}`)
    }
    return res.json()
}

export async function apiPut<T = unknown>(path: string, body: unknown): Promise<T> {
    const res = await apiFetch(path, {
        method: "PUT",
        body: JSON.stringify(body),
    })
    if (!res.ok) {
        const data = await res.json().catch(() => null)
        throw new Error(data?.error || `Error ${res.status}`)
    }
    return res.json()
}

export async function apiPatch<T = unknown>(path: string, body?: unknown): Promise<T> {
    const res = await apiFetch(path, {
        method: "PATCH",
        body: body ? JSON.stringify(body) : undefined,
    })
    if (!res.ok) {
        const data = await res.json().catch(() => null)
        throw new Error(data?.error || `Error ${res.status}`)
    }
    return res.json()
}
