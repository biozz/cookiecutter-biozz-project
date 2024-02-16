import { nanoid } from "nanoid";

export type Item = {
  id: string;
  name: string;
};

export type MaybeString = string | null;

export interface ItemsClient {
  getItems(): Promise<[Item[], MaybeString]>;
}

export enum Mode {
  VIEW = 0,
  EDIT,
  CREATE,
}

export type DebtEntry = {
  id: string;
  amount: number;
  created_at: string;
  comment: string;
};

export const clientID = nanoid();
