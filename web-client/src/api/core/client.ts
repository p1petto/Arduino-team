export class Client {
  baseUrl: string;
  token: string;

  constructor(baseUrl: string, token: string = "") {
    this.token = token;
    this.baseUrl = baseUrl;
  }
  url(uri: string) {
    return new URL(uri, this.baseUrl);
  }

  headers() {
    return {
      Authorization: `Bearer ${this.token}`,
      "Content-Type": "application/x-www-form-urlencoded",
    };
  }

  post = (uri: string, data?: Record<string, string>) =>
    fetch(this.url(uri), {
      method: "POST",
      headers: this.headers(),
      body: new URLSearchParams({ ...data }),
    });

  get = async (uri: string) =>
    fetch(this.url(uri), { method: "GET", headers: this.headers() });
}
