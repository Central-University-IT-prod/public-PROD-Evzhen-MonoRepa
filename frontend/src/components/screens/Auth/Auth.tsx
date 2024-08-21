import { FC, useContext } from 'react'
import { SubmitHandler, useForm } from 'react-hook-form'
import { YellowButton } from '../../shared/YellowButton/YellowButton'
import styles from './Auth.module.scss'
import { Context } from '../../../main'

interface IAuthInputs {
  login: string
  password: string
}

export const Auth: FC = () => {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IAuthInputs>()
  const { store } = useContext(Context)

  const onSubmit: SubmitHandler<IAuthInputs> = async (data): Promise<void> => {
    await store.login(data.login, data.password)
    if (!store.isAuth) {
      window.alert('указан неверный логин и/или пароль')
    }
  }

  return (
    <div className={styles.auth}>
      <div className="container">
        <div className={styles.auth__inner}>
          <h2 className={styles['auth__title']}>Вход</h2>
          <form className={styles['auth__form']} onSubmit={handleSubmit(onSubmit)}>
            <label className={styles['auth__form-label']} htmlFor="login">
              Логин
            </label>
            <input
              className={styles['auth__form-input']}
              id="login"
              {...register('login', { required: 'Это поле обязательное!' })}
              type="text"
              placeholder="admin"
            />
            {errors.login && (
              <span className={styles['auth__form-err']} role="alert">
                {errors.login.message}
              </span>
            )}
            <label className={styles['auth__form-label']} htmlFor="password">
              Пароль
            </label>
            <input
              className={styles['auth__form-input']}
              id="password"
              {...register('password', { required: 'Это поле обязательное!' })}
              type="password"
              placeholder="admin"
            />
            {errors.password && (
              <span className={styles['auth__form-err']} role="alert">
                {errors.password.message}
              </span>
            )}
            <YellowButton type="submit">Готово</YellowButton>
          </form>
        </div>
      </div>
    </div>
  )
}
