export async function request<T>(path: string, options?: RequestInit): Promise<T> {
    const response = await fetch(path, options);
    if (!response.ok) {
        const body = await response.json().catch(() => null);
        throw new APIError(body?.error ?? `Request failed (${response.status})`, response.status);
    }
    if (response.status === 204) return undefined as T;
    return response.json();
}

export class APIError extends Error {
    constructor(message: string, public status: number) {
        super(message);
    }
}
