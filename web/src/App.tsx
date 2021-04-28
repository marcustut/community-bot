import React, { useEffect, useState } from 'react';
import BulletScreen, { StyledBullet } from 'rc-bullets';

const headUrl='https://zerosoul.github.io/rc-bullets/assets/img/heads/girl.jpg';

const App: React.FC = () => {
  const [screen, setScreen] = useState<JSX.Element[]>([]);
  const [bullet, setBullet] = useState<string>('');

  useEffect(() => {
    setScreen(new BulletScreen('.screen', { duration: 20 }))
  }, [])

  const handleChange: React.ChangeEventHandler<HTMLInputElement> = ({ target: { value } }) => {
    setBullet(value);
  };

  const handleSend = () => {
    if (bullet) {
      screen.push(
        <StyledBullet
          head={headUrl}
          msg={bullet}
          backgroundColor={'#fff'}
          size='large'
        />
      );
    }
  };
  
  return (
    <>
      <div className="screen" style={{ width: '100vw', height: '80vh', background: '#000000' }}></div>
      <input value={bullet} onChange={handleChange} />
      <button onClick={handleSend}>发送</button>
    </>
  );
}

export default App;
