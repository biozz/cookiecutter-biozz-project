import { Signal, createEffect, createSignal } from "solid-js";
import { SetStoreFunction, Store, createStore } from "solid-js/store";

export function createLocalStore<T extends object>(
  key: string,
  initState: T
): [Store<T>, SetStoreFunction<T>] {
  const [state, setState] = createStore(initState);
  if (localStorage.getItem(key)) setState(JSON.parse(localStorage.getItem(key) || ''));
  createEffect(() => (localStorage.todos = JSON.stringify(state)));
  return [state, setState];
}

export function createStoredSignal<T>(
  key: string,
  defaultValue: T,
  storage = localStorage
): Signal<T> {

  const initialValue = storage.getItem(key)
    ? JSON.parse(storage.getItem(key) || '') as T
    : defaultValue;

  const [value, setValue] = createSignal<T>(initialValue);

  const setValueAndStore = ((arg: any) => {
    const v = setValue(arg);
    storage.setItem(key, JSON.stringify(v));
    return v;
  }) as typeof setValue;

  return [value, setValueAndStore];
}

