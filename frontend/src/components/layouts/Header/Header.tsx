import { Link, useLocation } from 'react-router-dom'
import { YellowButton } from '../../shared/YellowButton/YellowButton'
import styles from './Header.module.scss'
import { APP_URL } from '../../../constants/constants'

export default function Header() {
  const location = useLocation()

  return (
    <div className={styles.header}>
      <div className={`container ${styles['header__container']}`}>
        <Link to={'/'} className={styles.header__wrap}>
          <img src={`${APP_URL}/logo.svg`} alt="logo" width={50} style={{ borderRadius: 5 }} />
          <p className={styles['header__title']}>teamder admin</p>
        </Link>
        {location.pathname !== '/' ? (
          <Link to="/">
            <YellowButton className={styles['header__btn']}>Вернуться к списку событий</YellowButton>
          </Link>
        ) : (
          <YellowButton className={styles.add}>
            <Link to="/new-contest">Добавить событие</Link>
          </YellowButton>
        )}
      </div>
      {location.pathname === '/' && <hr className={styles['header__divider']} />}
    </div>
  )
}
