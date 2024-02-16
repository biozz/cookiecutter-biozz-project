import { nanoid } from "nanoid";
import { Item, ItemsClient } from "../types";

export class MockClient implements ItemsClient {
  clientID: string;
  mockedItems: any;
  mockedEntries: any;

  constructor() {
    this.clientID = nanoid();
    this.mockedItems = [...mockedItems];
  }

  public async getItems(itemType: string, namespace: string): Promise<Result> {
    let result: Result = {
      status: Status.OK,
      error: null,
      data: this.mockedItems.filter(
        (x: Item) => x.type === itemType && x.namespace === namespace,
      ),
    };
    return new Promise((resolve, reject) => resolve(result));
  }
}

export const mockedItems = [
  {
    id: nanoid(),
    name: "🥛Молоко",
  },
  {
    id: nanoid(),
    name: "🍨Сливки",
  },
  {
    id: nanoid(),
    name: "🧀Сыр",
  },
  {
    id: nanoid(),
    name: "🥖Белый хлеб",
  },
  {
    id: nanoid(),
    name: "🍪Печенье",
  },
];

export const mockClient = new MockClient();
