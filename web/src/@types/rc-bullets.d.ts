declare module 'rc-bullets' {
  const BulletScreen: any;
  export default BulletScreen;

  export class StyledBullet extends React.Component<StyledBulletProps & any, any> {};

  interface StyledBulletProps {
    msg: string;
    head: string;
    size: 'small' | 'normal' | 'large' | 'huge';
    color: string;
    backgroundColor: string;
  }
}