import { Component, For, onMount } from "solid-js";
import { items } from "./utils/context/store";
import { getItems } from "./utils/context/mutations";

const App: Component = () => {
  onMount(async () => {
    await getItems();
  });
  return (
    <div>
      <ul>
        <For each={items}>
          {(item) => (
            <li>
              {item.id} {item.name}
            </li>
          )}
        </For>
      </ul>
    </div>
  );
};

export default App;
