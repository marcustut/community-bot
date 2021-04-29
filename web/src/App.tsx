import React, { useEffect, useState } from 'react';
import BulletScreen, { StyledBullet } from 'rc-bullets';
import { w3cwebsocket as W3CWebSocket } from 'websocket';

const client = new W3CWebSocket('ws://127.0.0.1:8000/ws/messages');

const App: React.FC = () => {
  const [screen, setScreen] = useState<JSX.Element[]>([]);

  //generates a random size and returns that said size as a string
  //outputs 'small', 'medium' or 'large'
  const generateRandomSize = ()=>{

    const sizes = [
      'small',
      'medium',
      'large',
    ]

    const index = Math.floor(Math.random() * sizes.length);
    return sizes[index]

  }

   //generates a random colour and returns that said size as a string
  //outputs rgb colour value
  const generateRandomColour = ()=>{

    const colours = [
      'rgb(255, 129, 120)',
      'rgb(213, 255, 161)',
      'rgb(182, 255, 179)',
      'rgb(179, 255, 242)',
      'rgb(181, 214, 255)',
      'rgb(146, 151, 247)',
      'rgb(205, 167, 252)',
      'rgb(205, 167, 252)',
      'rgb(245, 171, 208)',
      'rgb(255, 219, 120)'
    ]

    const index = Math.floor(Math.random() * colours.length);
    return colours[index]

  }

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
          backgroundColor={generateRandomColour()}
          size={generateRandomSize()}
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
