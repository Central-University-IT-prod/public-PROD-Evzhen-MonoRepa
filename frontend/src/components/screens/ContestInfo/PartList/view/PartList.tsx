import styles from '../css/PartList.module.scss'
import readXlsxFile from 'read-excel-file'
import PartItem from '../PartItem/PartItem'
import IUser from '../../../../../models/User/IUser'
import { useState, useEffect } from 'react'
import { useOutside } from '../../../../../hooks/useOutside'
import { ErrorsModal } from '../ErrorsModal/ErrorsModal'

import Button from '@mui/material/Button'

import AuthService from '../../../../../services/AuthService'
import { IContestTeam } from '../../../../../models/User/IContestTeam'
import { APP_URL } from '../../../../../constants/constants'

interface IError {
  err: string
  login: string
}

interface IUserItem {
  surname: string
  name: string
  patronymic: string
  tg: string
  team_id: number
  role: string
  contest_id: number
  prev_points: number
  track: string
  user_login: string
  id: number
}

interface IUsersReq {
  error: false
  errors: IError[]
  usersList: IUserItem[]
}

export type { IUsersReq }

export default function PartList({
  contestId,
  usersReq,
  setUsersReq
}: {
  contestId: number
  usersReq: IUsersReq
  setUsersReq: (usersReq: IUsersReq) => void
}) {
  const [onlyNoTeam, setOnlyNoTeam] = useState(false)
  const [teams, setTeams] = useState<IContestTeam[]>([])

  const { ref, isShow, setIsShow } = useOutside(false)

  useEffect(() => {
    AuthService.getCommands(contestId).then((data) => {
      setTeams(data.data)
    })
  }, [])

  const removeById = (id: number) => {
    setUsersReq({ ...usersReq, usersList: usersReq.usersList.filter((item) => item.id !== id) })
  }

  const handleChangeFile = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      readXlsxFile(event.target.files[0]).then((data) => {
        const nameIndex = data[0].findIndex((item) => item.toString() === 'ФИО')
        const loginIndex = data[0].findIndex((item) => item.toString() === 'Логин на платформе')
        const tgIndex = data[0].findIndex((item) => item.toString() === 'Телеграм')
        const prevPointsIndex = data[0].findIndex((item) => item.toString() === 'Балл за предыдущий этап')
        const trackIndex = data[0].findIndex((item) => item.toString() === 'Трек')
        const maxPrevPoints = data[0].findIndex((item) => item.toString() === 'Максимальный балл в предыдущем этапе')
        if (nameIndex === -1 || loginIndex === -1 || tgIndex === -1) {
          console.log('Err ')
          return
        }
        data.shift()
        const parts: IUser[] = []
        data.forEach((item) => {
          const curItem: IUser = {
            name: item[nameIndex] + '',
            login: item[loginIndex] + '',
            tg: item[tgIndex] + '',
            prev_points: +item[prevPointsIndex],
            max_points: +item[maxPrevPoints],
            track: item[trackIndex] + ''
          }
          parts.push(curItem)
        })
        AuthService.loadContestUsers(contestId, parts).then((data) => {
          setUsersReq({
            error: false,
            errors: data.data.errors,
            usersList: [...usersReq.usersList, ...data.data.created_profiles]
          })
          if (data.data.errors.length > 0) {
            setIsShow(true)
          }
        })
      })
    }
  }

  return (
    <div ref={ref} className={styles['load__wrapper']}>
      {isShow && <ErrorsModal errors={usersReq.errors} setIsModalShow={setIsShow} />}
      <div className={styles.load}>
        <p className={styles.load__title}>Пример входного excel файла</p>
        <a href={`${APP_URL}/excel_example.png`} target="_blank" className={styles.load__example}></a>
        <Button component="label" role={undefined} variant="contained" tabIndex={-1} className={styles['load__button']}>
          Загрузить список участников
          <input type="file" className={styles.load__input} onChange={handleChangeFile} />
        </Button>
      </div>
      <div className={styles['part__filter__container']}>
        <p className={styles['part__title']}>Список участников {onlyNoTeam ? 'без команды' : ''}</p>
        <button
          className={`${onlyNoTeam ? styles.active : ''} ${styles.part__filter}`}
          onClick={() => setOnlyNoTeam(!onlyNoTeam)}
        >
          Участники без команды
        </button>
      </div>
      <div className={styles.part__list}>
        {usersReq.usersList
          .filter((item) => (onlyNoTeam ? !item.team_id : true))
          .map((user) => (
            <PartItem
              user={user}
              key={user.id}
              removeById={removeById}
              team={teams.find((team) => team.id === user.team_id)?.name ?? 'Нет команды'}
            />
          ))}
      </div>
    </div>
  )
}
