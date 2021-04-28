import React, { useEffect, useState } from 'react';
import BulletScreen, { StyledBullet } from 'rc-bullets';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const client = new W3CWebSocket('ws://127.0.0.1:8000/ws/messages');

const App: React.FC = () => {
  const [screen, setScreen] = useState<JSX.Element[]>([]);

  useEffect(() => {
    setScreen(new BulletScreen('.screen', { duration: 20 }))
  }, [])

  useEffect(() => {
    client.onopen = () => console.log("Websocket connected");
    client.onclose = () => console.error("Websocket disconnected");
    client.onmessage = (message) => {
      const dataFromServer = JSON.parse(message.data as string);
      console.log('got reply! ', dataFromServer);
      screen.push(
        <StyledBullet
          head={dataFromServer.avatarURL}
          msg={dataFromServer.text}
          backgroundColor={'#fff'}
          size='large'
        />
      )
    }
  })

  return (
    <>
      <div className="screen" style={{ width: '100vw', height: '100vh', background: '#000000' }}></div>
    </>
  );
}

export default App;
