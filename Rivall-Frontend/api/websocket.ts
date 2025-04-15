import { useUserStore } from "@/global-store/user_store";
import { useRouter } from "expo-router";
import { HOST, PORT, VERSION } from "@/api/axios.config";

export function connectToWebSocket(user_id: string, otp: string) : any {
  console.log('Connecting to Websocket...');
  const endpoint = `ws://${HOST}:${PORT}/api/${VERSION}/ws/connect/${user_id}?otp=${otp}`;
  return new WebSocket(endpoint);
}
