import { FC, useEffect, useState } from 'react'
import styles from './ContestTeamModal.module.scss'
import { IoIosCloseCircleOutline } from 'react-icons/io'
import { IContestTeam } from '../../../../../models/User/IContestTeam'
import AuthService, { ILoadedProfile } from '../../../../../services/AuthService'
import { YellowButton } from '../../../../shared/YellowButton/YellowButton'

interface IContestTeamModalProps {
  team: IContestTeam
  setIsModalShow: React.Dispatch<React.SetStateAction<boolean>>
  contestId: number
}

export const ContestTeamModal: FC<IContestTeamModalProps> = ({
  team,
  setIsModalShow,
  contestId
}: IContestTeamModalProps) => {
  const [participants, setParticipants] = useState<ILoadedProfile[]>([] as ILoadedProfile[])
  const [commands, setCommands] = useState<IContestTeam[]>([] as IContestTeam[])
  const [isApproved, setIsApproved] = useState<boolean>(false)
  useEffect(() => {
    const asd = async (): Promise<void> => {
      const newParticipants = await AuthService.getCommandMembers(team.id)
      const commandss = await AuthService.getCommands(contestId)
      setParticipants(newParticipants.data)
      setCommands(commandss.data)
      commandss.data.forEach((command) => {
        if (command.id == team.id) {
          setIsApproved(Boolean(command.approved) || false)
        }
      })
    }
    asd()
  }, [])

  const handleDebatCommand = async (val: number): Promise<void> => {
    await AuthService.debatCommand(team.id, val)
    setIsModalShow(false)
  }

  return (
    <div className={styles.modal}>
      <button className={styles['modal__close']} type="button" onClick={() => setIsModalShow(false)}>
        <IoIosCloseCircleOutline className={styles['modal__close-icon']} />
      </button>
      <div className={styles['modal__item']}>
        <h4 className={styles['modal__item-title']}>Название команды</h4>
        <h4 className={styles['modal__item-value']}>{team.name}</h4>
      </div>
      <div className={styles['modal__item']}>
        <h4 className={styles['modal__item-title']}>Участники команды</h4>
        <div className={styles['modal__participant-list']}>
          {!participants.length && 'Участников пока что нет'}
          {participants.map(({ id, name, tg, role, prev_points }) => (
            <div className={styles['modal__hui']} key={id}>
              <div className={styles['modal__item-participant']} key={id}>
                <div className={styles['modal__item-value-wrapper']}>
                  <h4 className={styles['modal__item-login']}>{name}</h4>
                  <a href={`https://t.me/${tg}`} target="_blank" className={styles['modal__item-tg']}>
                    @{tg}
                  </a>
                </div>
                <div className={styles['modal__item-value-wrapper']}>
                  <h4 className={styles['modal__item-gpa']}>Предыдущий балл: {prev_points.toFixed(2)}</h4>
                  <h4 className={styles['modal__item-role']}>Роль: {role}</h4>
                </div>
              </div>
            </div>
          ))}
          {isApproved ? (
            <div onClick={() => handleDebatCommand(0)}>
              <YellowButton className={styles['modal__item-btn']}>Опровергнуть команду</YellowButton>
            </div>
          ) : (
            <div onClick={() => handleDebatCommand(1)}>
              <YellowButton className={styles['modal__item-btn']}>Утвердить команду</YellowButton>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
