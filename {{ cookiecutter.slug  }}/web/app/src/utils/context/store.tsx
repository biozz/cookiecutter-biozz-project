import { createStore } from "solid-js/store";
import { Item } from "../types";

export const [items, setItems] = createStore([] as Item[]);

