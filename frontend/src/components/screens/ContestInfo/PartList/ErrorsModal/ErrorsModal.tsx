import { FC } from 'react'
import styles from './ErrorsModal.module.scss'
import { IoIosCloseCircleOutline } from 'react-icons/io'

interface IContestTeamModalProps {
  errors: {
    err: string
    login: string
  }[]
  setIsModalShow: React.Dispatch<React.SetStateAction<boolean>>
}

export const ErrorsModal: FC<IContestTeamModalProps> = ({ errors, setIsModalShow }: IContestTeamModalProps) => {
  return (
    <div className={styles.modal}>
      <button className={styles['modal__close']} type="button" onClick={() => setIsModalShow(false)}>
        <IoIosCloseCircleOutline className={styles['modal__close-icon']} />
      </button>
      <div className={styles['modal__item']}>
        <h4 className={styles['modal__item-title']}>Ошибки</h4>
      </div>
      <div className={styles['modal__item']}>
        <div className={styles['modal__errors-list']}>
          {errors.map((item) => (
            <div className={styles['error__item']} key={item.login}>
              <div className={styles['error__title']}>{item.login}</div>
              <div className={styles['error__title']}>{item.err}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
