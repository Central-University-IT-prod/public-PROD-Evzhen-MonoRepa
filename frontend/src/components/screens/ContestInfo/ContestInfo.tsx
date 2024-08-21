import { FC, useState } from 'react'
import { useParams } from 'react-router-dom'
import styles from './ContestInfo.module.scss'
import { TabContext, TabList, TabPanel } from '@mui/lab'
import { Box, Tab } from '@mui/material'
import { ContestEdit } from './ContestEdit/ContestEdit'
import PartList, { IUsersReq } from './PartList/view/PartList'
import { ContestTeams } from './ContestTeams/ContestTeams'
import { useEffect } from 'react'
import AuthService from '../../../services/AuthService'

export const ContestInfo: FC = () => {
  const params = useParams()
  const [value, setValue] = useState('1')

  const handleChange = (event: React.SyntheticEvent, newValue: string) => {
    setValue(newValue)
  }

  const contestId = Number(params.id?.split('=')[1])

  const [usersReq, setUsersReq] = useState<IUsersReq>({
    errors: [],
    error: false,
    usersList: []
  })

  useEffect(() => {
    AuthService.getContestUsers(contestId).then((data) => {
      setUsersReq({ ...usersReq, usersList: data.data })
    })
  }, [])

  return (
    <div className={styles.info}>
      <div className="container">
        <TabContext value={value}>
          <Box sx={{ borderBottom: 2, borderColor: 'divider' }} color={'white'} borderColor={'white'}>
            <TabList
              variant="scrollable"
              scrollButtons
              allowScrollButtonsMobile
              className={styles['info__tabs-list']}
              onChange={handleChange}
              aria-label="tabs"
              textColor="secondary"
              indicatorColor="secondary"
            >
              <Tab className={styles['info__tabs-item']} color="white" label="Редактирование олимпиады" value="1" />
              <Tab className={styles['info__tabs-item']} label="Список созданных команд" value="2" />
              <Tab className={styles['info__tabs-item']} label="Список участников" value="3" />
            </TabList>
          </Box>
          <TabPanel className={styles['info__edit']} value="1">
            <ContestEdit contestId={contestId} />
          </TabPanel>
          <TabPanel className={styles['info__teams']} value="2">
            <ContestTeams contestId={contestId} />
          </TabPanel>
          <TabPanel className={styles['info__participants']} value="3">
            <PartList contestId={contestId} usersReq={usersReq} setUsersReq={setUsersReq} />
          </TabPanel>
        </TabContext>
      </div>
    </div>
  )
}
