import ContestItem from '../ContestItem/ContestItem'
import styles from './ContestList.module.scss'
import { IContest } from '../../../../models/User/IContest'
import { useEffect, useState } from 'react'
import AuthService from '../../../../services/AuthService'
import { Loading } from '../../../shared/Loading/Loading'

export const ContestList = () => {
  const [contests, setContests] = useState<IContest[]>([] as IContest[])
  const [isLoading, setIsLoading] = useState<boolean>(false)

  useEffect(() => {
    const getContests = async (): Promise<void> => {
      setIsLoading(true)
      try {
        const response = await AuthService.getContests()
        setContests(response.data)
      } catch (e) {
        console.log(e)
      } finally {
        setIsLoading(false)
      }
    }
    getContests()
  }, [])

  return (
    <div className="container">
      <div className={styles.contest__list}>
        <div className={styles.head}>
          <h3 className={styles.head__title}>Список событий</h3>
        </div>
        <div>
          {isLoading ? (
            <Loading top={10} className={styles.loading} />
          ) : (
            <>
              {contests?.length ? (
                contests.map((contest) => <ContestItem contest={contest} key={contest.id} />)
              ) : (
                <>
                  <h3 className={styles['no__contests']} style={{ textAlign: 'start', fontSize: 20 }}>
                    У вас пока что нет событий :(
                  </h3>
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
