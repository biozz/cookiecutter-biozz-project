import { itemsClient } from "../clients/client";
import { Item } from "../types";
import { setItems } from "./store";

export const getItems = async () => {
  let [result] = await itemsClient.getItems();
  let resultData = result || [];
  let newItems: Item[] = resultData.map((i: any) => {
    return {
      id: i.id,
      name: i.name,
    };
  });
  setItems(newItems);
};
