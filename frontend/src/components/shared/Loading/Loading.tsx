import { FC } from 'react'
import { PulseLoader } from 'react-spinners'

interface ILoadingProps {
  size?: number
  top?: number
  className?: string
}

export const Loading: FC<ILoadingProps> = (props: ILoadingProps) => {
  return (
    <PulseLoader
      className={props.className}
      color="#ffdd2d"
      size={props.size || 10}
      style={{
        width: '100vw',
        height: `${props.top || 100}vh`,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center'
      }}
    />
  )
}
