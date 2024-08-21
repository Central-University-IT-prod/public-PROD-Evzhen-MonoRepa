import { SubmitHandler, useForm } from 'react-hook-form'
import Select from 'react-select'

import { YellowButton } from '../../../shared/YellowButton/YellowButton'

import styles from './CreateContest.module.scss'
import { useContext, useState } from 'react'
import { Context } from '../../../../main'
import { Loading } from '../../../shared/Loading/Loading'
import { useNavigate } from 'react-router-dom'

export interface IContestInput {
  name: string
  description: string
  start_date: Date | string
  end_date: Date | string
}

export default function CreateContest() {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IContestInput>()
  const [selected, setSelected] = useState('Программироание')
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const { store } = useContext(Context)
  const navigateTo = useNavigate()

  const onSubmit: SubmitHandler<IContestInput> = async (data): Promise<void> => {
    const createContest = async (): Promise<void> => {
      setIsLoading(true)
      await store.createContest({
        ...data,
        field: selected,
        start_date: Date.parse(new Date(data.start_date).toISOString()) / 1000,
        end_date: Date.parse(new Date(data.end_date).toISOString()) / 1000
      })
      setIsLoading(false)
      setTimeout(() => {
        navigateTo('/')
      }, 2500)
    }
    createContest()
  }

  const options = [
    { value: 'programming', label: 'Программирование' },
    { value: 'math', label: 'Математика' },
    { value: 'biology', label: 'Биология' },
    { value: 'russian', label: 'Русский язык' },
    { value: 'checmicks', label: 'Химия' },
    { value: 'english', label: 'Английский язык' },
    { value: 'physics', label: 'Физика' },
    { value: 'geographic', label: 'География' }
  ]

  return (
    <div className={`container ${styles.container}`}>
      {isLoading ? (
        <Loading top={10} />
      ) : (
        <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
          <label className={styles.label} htmlFor="title">
            Название
          </label>
          <input
            className={styles.input}
            id="title"
            {...register('name', { required: 'Это поле обязательное!' })}
            type="text"
            placeholder="PROD"
          />
          {errors.name && (
            <span className={styles.error} role="alert">
              {errors.name.message}
            </span>
          )}
          <label className={styles.label} htmlFor="description">
            Описание
          </label>
          <textarea
            className={styles.multi}
            id="description"
            {...register('description', { required: 'Это поле обязательное!' })}
            placeholder="Олимпиада по промышленной разработке"
          />
          {errors.description && (
            <span className={styles.error} role="alert">
              {errors.description.message}
            </span>
          )}
          <label className={styles.label} htmlFor="title">
            Даты проведения
          </label>
          <div className={styles['date-container']}>
            <input
              className={`${styles['date-input']} ${styles.input}`}
              id="startDate"
              {...register('start_date', { required: 'Это поле обязательное!' })}
              type="date"
            />
            <input
              className={`${styles['date-input']} ${styles.input}`}
              id="endDate"
              {...register('end_date', { required: 'Это поле обязательное!' })}
              type="date"
            />
          </div>
          {(errors.start_date || errors.end_date) && (
            <span className={styles.error} role="alert">
              Это поле обязательное!
            </span>
          )}
          <label className={styles.label} htmlFor="title">
            Направление
          </label>
          <Select
            className={styles.select}
            options={options}
            defaultValue={options[0]}
            placeholder="Выберите направление"
            value={options.find((item) => item.label === selected)}
            onChange={(event) => setSelected(event?.label ?? '')}
            theme={(theme) => ({
              ...theme,
              borderRadius: 5,
              colors: {
                ...theme.colors,
                primary: 'rgba(0, 0, 0, 0.7)',
                primary25: 'rgba(0, 0, 0, 0.5)',
                primary50: 'rgba(0, 0, 0, 0.6)'
              }
            })}
            styles={{
              option: (baseStyles, state) =>
                state.isFocused || state.isSelected
                  ? {
                      ...baseStyles,
                      color: 'white',
                      fontWeight: 400
                    }
                  : {
                      ...baseStyles,
                      color: 'black',
                      fontWeight: 400
                    }
            }}
          />
          <YellowButton type="submit">Создать</YellowButton>
        </form>
      )}
    </div>
  )
}
