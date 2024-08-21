import { FC, useLayoutEffect, useState } from 'react'
import styles from './ContestTeams.module.scss'
import { useOutside } from '../../../../hooks/useOutside'
import { ContestTeamModal } from './ContestModal/ContestTeamModal'
import { IContestTeam } from '../../../../models/User/IContestTeam'
import AuthService, { ILoadedProfile } from '../../../../services/AuthService'
import { Loading } from '../../../shared/Loading/Loading'
import { HiOutlineTrash } from 'react-icons/hi'

export const ContestTeams: FC<{ contestId: number }> = ({ contestId }: { contestId: number }) => {
  const [currentTeam, setCurrentTeam] = useState<IContestTeam>({} as IContestTeam)
  const [isLoading, setIsLoading] = useState<boolean>(false)
  const { ref, isShow, setIsShow } = useOutside(false)
  const [onlyNotReady, setOnlyNotReady] = useState(false)
  const [minLimits, setMinLimits] = useState<string | undefined>(undefined)
  const [maxLimits, setMaxLimits] = useState<string | undefined>(undefined)
  const [isLimits, setIsLimits] = useState<boolean>(false)

  const [participants, setParticipants] = useState<ILoadedProfile[]>([] as ILoadedProfile[])
  const [contestTeamsList, setContestTeamsList] = useState<IContestTeam[]>([] as IContestTeam[])

  let asd = false
  useLayoutEffect(() => {
    const defineLimits = async (): Promise<void> => {
      const contests = await AuthService.getContests()
      contests.data?.forEach((contest) => {
        if (contest.id === contestId && contest.min_teammates && contest.max_teammates) {
          setIsLimits(true)
          setMinLimits(String(contest.min_teammates))
          setMaxLimits(String(contest.max_teammates))
        }
      })
    }
    defineLimits()
    if (!contestTeamsList.length && !asd) {
      AuthService.getCommands(contestId).then((res) => {
        setContestTeamsList(res.data)
        asd = true
      })
    }
  })

  const handleSetLimits = async (): Promise<void> => {
    setIsLoading(true)
    try {
      await AuthService.setLimit(contestId, Number(minLimits), Number(maxLimits))
      const commands = await AuthService.getCommands(contestId)
    } catch (e) {
      console.log(e)
    } finally {
      setIsLoading(false)
    }
  }

  const handleRemoveCommand = async (commandId: number) => {
    await AuthService.removeCommand(commandId)
    await AuthService.getCommands(contestId).then((res) => {
      setContestTeamsList(res.data)
      asd = true
    })
  }

  if (isLoading) return <Loading top={10} />

  let limitBlock
  let mainBlock
  if (!isLimits) {
    limitBlock = (
      <div className={styles['teams__limits__container']}>
        <div className={styles['teams__limits__container-cont']}>
          <input
            type="number"
            placeholder="Минимальное число участников"
            className={styles['teams__limits__limit']}
            value={minLimits || ''}
            onChange={(event) => setMinLimits(event.target.value)}
          />
          <input
            type="number"
            placeholder="Максимальное число участников"
            className={styles['teams__limits__limit']}
            value={maxLimits || ''}
            onChange={(event) => setMaxLimits(event.target.value)}
          />
        </div>
        <button className={styles['teams__limits__submit']} onClick={handleSetLimits}>
          Установить
        </button>
      </div>
    )
    mainBlock = (
      <div className={styles.teams__placeholder}>
        Задайте ограничения, чтобы дать участникам возможность создавать команды
      </div>
    )
  } else {
    limitBlock = (
      <div className={styles['teams__limits__not-editable']}>
        <div className={styles['teams__limits__not-editable__limit']}>
          Минимальное число участников в команде: {minLimits}
        </div>
        <div className={styles['teams__limits__not-editable__limit']}>
          Максимальное число участников в команде: {maxLimits}
        </div>
      </div>
    )
    mainBlock = (
      <>
        <div className={styles['teams__filter__container']}>
          <p className={styles['teams__title']}>Список команд</p>
          <button
            className={`${onlyNotReady ? styles.active : ''} ${styles.teams__filter}`}
            onClick={() => setOnlyNotReady(!onlyNotReady)}
          >
            Не готовые команды
          </button>
        </div>
        <div className={styles['teams__list']} ref={ref}>
          {isShow && <ContestTeamModal contestId={contestId} team={currentTeam} setIsModalShow={setIsShow} />}
          {contestTeamsList
            .filter(async (item) => {
              if (!onlyNotReady) {
                return true
              } else {
                AuthService.getCommandMembers(item.id).then((res) => setParticipants(res?.data))
                if (!(participants.length >= +(minLimits ?? 0) && participants.length <= +(maxLimits ?? 0))) {
                  return true
                }
                return false
              }
            })
            .map(({ id, name, description, owner_id, contest_id }) => (
              <div
                className={styles['teams__item']}
                onClick={() => {
                  setCurrentTeam({ id, name, description, owner_id, contest_id })
                  setIsShow(!isShow)
                }}
                key={id}
              >
                <h4 className={styles['teams__item-title']}>{name}</h4>
                <HiOutlineTrash style={{ fontSize: 30 }} onClick={() => handleRemoveCommand(id)} />
              </div>
            ))}
        </div>
      </>
    )
  }
  return (
    <div className={styles.wrapper}>
      {limitBlock}
      {mainBlock}
    </div>
  )
}
