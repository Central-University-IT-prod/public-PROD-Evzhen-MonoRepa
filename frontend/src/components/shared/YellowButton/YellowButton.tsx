import { FC } from 'react'
import styles from './YellowButton.module.scss'

interface IYellowButton {
  children: React.ReactNode
  type?: 'button' | 'submit' | 'reset' | undefined
  className?: string
  onClick?: () => void
}

export const YellowButton: FC<IYellowButton> = (props) => {
  return (
    <button className={`${styles.button} ${props.className}`} type={props.type} onClick={props.onClick}>
      {props.children}
    </button>
  )
}
