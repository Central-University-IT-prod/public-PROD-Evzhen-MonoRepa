import { AxiosResponse } from 'axios'
import { AuthResponse } from '../models/Auth/AuthResponse'
import { IAdmin } from '../models/User/IAdmin'
import { $api, $authApi } from '../http'
import IUser from '../models/User/IUser'
import { IContest, IContestRequest } from '../models/User/IContest'
import { IContestTeam } from '../models/User/IContestTeam'

export interface ILoadedProfile {
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
interface ILoadContestUsers {
  errors: {
    err: string
    login: string
  }[]
  created_profiles: ILoadedProfile[]
}

interface IResProfile {
  login: string
  final_points: number
  max_final_points: number
}

interface IErrEndContest {
  id: number
  login: string
  err: string
}

export default class AuthService {
  static async login(login: string, password: string): Promise<AxiosResponse<AuthResponse>> {
    return $api.post<AuthResponse>('/login', { login, password })
  }
  static async getAdmin(): Promise<AxiosResponse<IAdmin>> {
    return $authApi.get<IAdmin>('/getAdmin')
  }
  static async loadContestUsers(contestId: number, contestUsers: IUser[]): Promise<AxiosResponse<ILoadContestUsers>> {
    return $authApi.post<ILoadContestUsers>('/contests/loadProfiles', { contest_id: contestId, profiles: contestUsers })
  }
  static async getContestUsers(contestId: number): Promise<AxiosResponse<ILoadedProfile[]>> {
    return $authApi.get<ILoadedProfile[]>(`/contest/${contestId}/profiles`)
  }
  static async createContest(contest: IContestRequest): Promise<AxiosResponse<void>> {
    return $authApi.post<void>('/contests/create', contest)
  }
  static async getContests(): Promise<AxiosResponse<IContest[]>> {
    return $authApi.get<IContest[]>('/getContests')
  }
  static async getCommands(contestId: number): Promise<AxiosResponse<IContestTeam[]>> {
    return $authApi.get<IContestTeam[]>(`/getCommands/${contestId}`)
  }
  static async setLimit(contestId: number, min_teammates: number, max_teammates: number): Promise<AxiosResponse<void>> {
    return $authApi.patch<void>('/contests/setTeamLimit', { contest_id: contestId, min_teammates, max_teammates })
  }
  static async removeUser(id: number): Promise<AxiosResponse<void>> {
    return $authApi.delete<void>(`/profile/${id}`)
  }
  static async getCommandMembers(commandId: number): Promise<AxiosResponse<ILoadedProfile[]>> {
    return $authApi.get<ILoadedProfile[]>(`/getTeammates/${commandId}`)
  }
  static async changeContest(
    contestId: number,
    name: string,
    description: string,
    start_date: number,
    end_date: number
  ): Promise<AxiosResponse<void>> {
    return $authApi.patch<void>(`/contest/change`, { contest_id: contestId, name, description, start_date, end_date })
  }
  static async loadResults(contestId: number, contestUsers: IResProfile[]): Promise<AxiosResponse<IErrEndContest>> {
    return $authApi.post<IErrEndContest>('/contest/end', { contest_id: contestId, profiles: contestUsers })
  }
  static async removeCommand(commandId: number): Promise<AxiosResponse<void>> {
    return $authApi.delete<void>('/contest/command', { data: { command_id: commandId } })
  }
  static async debatCommand(commandId: number, approved: number): Promise<AxiosResponse<void>> {
    return $authApi.patch<void>(`/contest/command/${commandId}/approve/${approved}`)
  }
}
