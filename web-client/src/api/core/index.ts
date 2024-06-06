const API_URL = "http://localhost:1090";

export const get = (uri: string, headers?: HeadersInit) =>
  fetch(`${API_URL}${uri}`, { method: "GET", headers });

export const post = (uri: string, body?: object, headers?: HeadersInit) =>
  fetch(`${API_URL}${uri}`, {
    method: "POST",
    body: JSON.stringify(body),
    headers,
  });
