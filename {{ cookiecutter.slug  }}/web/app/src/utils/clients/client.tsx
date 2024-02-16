import { ItemsClient } from "../types";
import { httpClient } from "./http";
import { mockClient } from "./mock";

function getItemsClient(): ItemsClient {
  if (import.meta.env.VITE_MOCK) return mockClient;
  return httpClient;
}

export const itemsClient = getItemsClient();
