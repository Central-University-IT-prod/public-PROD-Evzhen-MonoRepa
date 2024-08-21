import { FC, useEffect, useState } from 'react'
import { SubmitHandler, useForm } from 'react-hook-form'
import { IContestInput } from '../../CreateContest/view/CreateContest'
import { IContest } from '../../../../models/User/IContest'
import styles from './ContestEdit.module.scss'
import { YellowButton } from '../../../shared/YellowButton/YellowButton'
import { convertDateToInput } from './ConvertDateToInput'
import AuthService from '../../../../services/AuthService'
import { Loading } from '../../../shared/Loading/Loading'
import Button from '@mui/material/Button'
import readXlsxFile from 'read-excel-file'

interface IResult {
  login: string
  final_points: number
  max_final_points: number
}

export const ContestEdit: FC<{ contestId: number }> = ({ contestId }: { contestId: number }) => {
  const [isEnd, setIsEnd] = useState(false)
  const [contest, setContest] = useState<IContest>()
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue
  } = useForm<IContestInput>()

  useEffect(() => {
    const getContests = async (): Promise<void> => {
      setIsLoading(true)
      try {
        const response = await AuthService.getContests()
        response.data?.forEach((contest) => {
          if (contest.id === contestId) {
            setContest(contest)
            setValue('name', contest?.name)
            setValue('description', contest?.description)
            setValue('start_date', convertDateToInput(contest?.start_date || -1) || '')
            setValue('end_date', convertDateToInput(contest?.end_date || -1) || '')
          }
        })
      } catch (e) {
        console.log(e)
      } finally {
        setIsLoading(false)
      }
    }
    getContests()
  }, [])

  const onSubmit: SubmitHandler<IContestInput> = async (data): Promise<void> => {
    AuthService.changeContest(
      contestId,
      data.name,
      data.description,
      Date.parse(new Date(data.start_date).toISOString()) / 1000,
      Date.parse(new Date(data.end_date).toISOString()) / 1000
    )
  }
  const handleChangeFile = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      readXlsxFile(event.target.files[0]).then((data) => {
        const loginIndex = data[0].findIndex((item) => item.toString() === 'Логин на платформе')
        const curPointsIndex = data[0].findIndex((item) => item.toString() === 'Балл за этот этап')
        const maxCurPoints = data[0].findIndex((item) => item.toString() === 'Максимальный балл в этапе')
        if (loginIndex === -1 || curPointsIndex === -1 || maxCurPoints === -1) {
          return
        }
        data.shift()
        const results: IResult[] = []
        data.forEach((item) => {
          const curItem: IResult = {
            login: item[loginIndex] + '',
            final_points: +item[curPointsIndex],
            max_final_points: +item[maxCurPoints]
          }
          results.push(curItem)
        })
        console.log(results)
        AuthService.loadResults(contestId, results)
      })
    }
  }

  let listBlock
  if (isEnd) {
    listBlock = (
      <>
        <div className={styles.load}>
          <Button
            component="label"
            role={undefined}
            variant="contained"
            tabIndex={-1}
            className={styles['load__button']}
          >
            Загрузить результаты
            <input type="file" className={styles.load__input} onChange={handleChangeFile} />
          </Button>
          <p className={styles.load__title}>Пример входного excel файла</p>
          <div className={styles.load__example}></div>
        </div>
      </>
    )
  } else {
    listBlock = (
      <>
        <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
          <label className={styles.label} htmlFor="title">
            Название
          </label>
          <input
            className={styles.input}
            id="title"
            defaultValue={contest?.name}
            {...register('name', { validate: (value) => value !== '' || 'Это поле обязательное!' })}
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
              defaultValue={convertDateToInput(contest?.start_date || -1) || undefined}
              id="startDate"
              {...register('start_date', { required: 'Это поле обязательное!' })}
              type="date"
            />
            <input
              className={`${styles['date-input']} ${styles.input}`}
              defaultValue={convertDateToInput(contest?.end_date || -1) || undefined}
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
          <label className={styles.label} style={{ fontSize: 15 }} htmlFor="title">
            Направление: {contest?.field}
          </label>
          <YellowButton type="submit">Применить изменения</YellowButton>
        </form>
        <button className={styles['finish-contest__button']} onClick={() => setIsEnd(true)}>
          Завершить олимпиаду
        </button>
      </>
    )
  }

  return <>{isLoading ? <Loading top={10} /> : listBlock}</>
}
