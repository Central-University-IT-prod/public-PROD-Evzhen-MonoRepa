import { Link } from 'react-router-dom'
import { IContest } from '../../../../models/User/IContest'
import styles from './ContestItem.module.scss'
import { convertDate } from './ConvertDate'
import { useEffect, useState } from 'react'
import AuthService from '../../../../services/AuthService'

export default function ContestItem({ contest }: { contest: IContest }) {
  const [count, setCount] = useState<number>(0)
  useEffect(() => {
    const f = async (): Promise<void> => {
      const response = await AuthService.getContestUsers(contest.id)
      setCount(response.data.length)
    }
    f()
  })
  return (
    <Link to={{ pathname: `/contest-info/:id=${contest.id}` }}>
      <div className={styles.contest__item}>
        <div className={styles['contest__item-title']}> {contest.name}</div>
        <div className={styles.col}>
          <p className={styles['contest__item-date']}>
            {convertDate(contest.start_date)} - {convertDate(contest.end_date)}
          </p>
          <div className={styles['contest__item-text']}>Количество участников: {count}</div>
        </div>
      </div>
    </Link>
  )
}
