import { Centrifuge } from "centrifuge";
import { clientID } from "../types";

export function createWebSocket() {
  if (import.meta.env.VITE_MOCK === "1") {
    return null;
  }
  let wsProto = "wss:";
  let proto = "http:";
  let host = window.location.host;
  let backendURL = import.meta.env.VITE_BACKEND;
  if (backendURL.length > 0) {
    [proto, host] = backendURL.split("//");
    if (proto === "http:") {
      wsProto = "ws:";
    }
  }
  let websokcetURL = `${wsProto}//${host}/api/ws?id=${clientID}`;
  let socket: WebSocket | null = null;
  if (!import.meta.env.VITE_MOCK) {
    socket = new WebSocket(websokcetURL);
  }
  return socket;
}

export async function getRealtimeToken(): Promise<string> {
  const response = await fetch(
    `${import.meta.env.VITE_BACKEND}/auth/realtime/token`,
  );
  const result = await response.json();
  return result["token"];
}

export async function goCentrifugo() {
  const wsURL = import.meta.env.VITE_CENTRIFUGO_URL.replace("http:", "ws:");
  const centrifuge = new Centrifuge(`${wsURL}/connection/websocket`, {
    token: await getRealtimeToken(),
  });
  try {
    await centrifuge.ready();
  } catch (e) {
    // console.log(e);
    return;
  }

  const sub = centrifuge.newSubscription(`items#${clientID}`);

  sub.on("publication", function (ctx) {
    console.log(ctx);
  });

  sub.subscribe();

  centrifuge.connect();
}
