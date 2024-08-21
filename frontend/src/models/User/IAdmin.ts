import { IContest } from './IContest'

export interface IAdmin {
  id: number
  login: string
  password: string
  contests: IContest[]
}
