import {useEffect, useRef} from 'react';
import { useUserStore } from '@/global-store/user_store';
import { HOST, PORT, VERSION } from '@/api/axios.config';

type Props = {
  receivedMessage: (message: any) => void;
}

export const useWebSockets = ({ receivedMessage }: Props) => {
  const user = useUserStore((state: any) => state.user);
  const access_token = useUserStore((state: any) => state.access_token);
  var ws = useRef(new WebSocket(`ws://${HOST}:${PORT}/api/${VERSION}/ws/connect/${user._id}?Authorization=${access_token}`)).current;

  useEffect(() => {

    ws.onopen = () => {
      console.log('WebSocket connection opened');
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onmessage = (event) => {
      console.log('WebSocket message received');
      const data = JSON.parse(event.data);

      if (data.type === 'new_message') {
        receivedMessage(data.payload);
      } else {
        console.log('Unknown message type:', data.type);
      }
    };

    return () => {
      ws.close();
    }
  }, [receivedMessage]);

  const sendMessage = (payload: any, direct_message_id: string ) => {
    if (ws.readyState === WebSocket.OPEN) {
      var msg = JSON.stringify({
        type: 'send_message',
        payload: payload,
        user_id: user._id,
        direct_message_id: direct_message_id,
        group_id: ''
      })

      ws.send(msg);
    } else {
      console.error('WebSocket is not open. Ready state:', ws.readyState);
    }
  }

  return {
    sendMessage
  };
};