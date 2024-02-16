import { Item, ItemsClient, MaybeString, clientID } from "../types";

export class HttpClient implements ItemsClient {
  baseUrl: string;
  clientID: string;

  constructor(baseUrl: string, clientID: string) {
    this.baseUrl = baseUrl;
    this.clientID = clientID;
  }

  private async request(
    method: string,
    path: string,
    body: any,
  ): Promise<[any, MaybeString]> {
    let opts: RequestInit = {
      method: method,
      body: null,
      headers: {
        "X-Client-ID": this.clientID,
      },
    };
    if (body !== null) {
      opts.body = JSON.stringify(body);
    }
    let url = `${import.meta.env.VITE_BACKEND}${path}`;
    let response = await fetch(url, opts);
    if (response.status === 401) {
      window.location.assign(response.headers.get("location") || "/");
      return [null, null];
    }
    let result = await response.json();
    if (result.status === "error") {
      return [null, result.error];
    }
    return [result.data, null];
  }

  public async getItems(): Promise<[Item[], MaybeString]> {
    let url = "/api/items";
    let [result, err] = await this.request("GET", url, null);
    if (err) return [[], err];
    return [result, null];
  }
}

export const httpClient = new HttpClient(
  import.meta.env.VITE_BACKEND,
  clientID,
);
