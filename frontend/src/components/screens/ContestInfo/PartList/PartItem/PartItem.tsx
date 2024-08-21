import styles from './PartItem.module.scss'
import { HiOutlineTrash } from 'react-icons/hi'

import AuthService from '../../../../../services/AuthService'

interface IUser {
  name: string
  surname: string
  patronymic: string
  team?: string
  prevPoints?: number
  tg: string
  user_login: string
  id: number
}

export default function PartItem({
  user,
  removeById,
  team
}: {
  user: IUser
  removeById: (id: number) => void
  team: string
}) {
  return (
    <div className={styles['part__item']}>
      <div>
        <div className={styles.part__name}>
          {user.surname} {user.name} {user.patronymic}
        </div>
        {user.prevPoints && <div className={styles.part__prev}>Балл за предыдущий этап: {user.prevPoints}</div>}
        <a href={`https://t.me/${user.tg}`} target="blank" className={styles.part__tg}>
          Телеграм @{user.tg}
        </a>
        <p className={styles.part__tg}>Логин {user.user_login}</p>
        <p className={styles['part__team']}>{team ?? 'Нет команды'}</p>
      </div>
      <div className={styles.part__right}>
        <HiOutlineTrash
          className={styles.part__delete}
          onClick={() => {
            AuthService.removeUser(user.id)
            removeById(user.id)
          }}
        />
      </div>
    </div>
  )
}
