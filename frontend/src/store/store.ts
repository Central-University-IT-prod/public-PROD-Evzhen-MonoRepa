import { makeAutoObservable } from 'mobx'
import AuthService from '../services/AuthService'
import { IAdmin } from '../models/User/IAdmin'
import IUser from '../models/User/IUser'
import { IContest, IContestRequest } from '../models/User/IContest'

export default class Store {
  admin = {} as IAdmin
  isAuth = false
  isLoading = false

  constructor() {
    makeAutoObservable(this)
  }

  setAuth(bool: boolean) {
    this.isAuth = bool
  }

  setAdmin(admin: IAdmin) {
    this.admin = admin
  }

  setLoading(bool: boolean) {
    this.isLoading = bool
  }

  async login(login: string, password: string) {
    this.setLoading(true)
    try {
      const response = await AuthService.login(login, password)
      localStorage.setItem('token', response.data.token)
      await this.getAdmin()
      this.setAuth(true)
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async getAdmin() {
    this.setLoading(true)
    try {
      const userResponse = await AuthService.getAdmin()
      const contests = await this.getContests()
      this.setAdmin({ ...userResponse.data, contests: contests || ([] as IContest[]) })
      return this.admin
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async getContests() {
    this.setLoading(true)
    try {
      const contests = await AuthService.getContests()
      return contests.data
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async loadContestFile(contestId: number, contestUsers: IUser[]) {
    this.setLoading(true)
    let res
    try {
      res = await AuthService.loadContestUsers(contestId, contestUsers)
      await this.getAdmin()
    } catch (e) {
      res = console.log(e)
    } finally {
      this.setLoading(false)
    }
    return res
  }

  async createContest(contest: IContestRequest) {
    this.setLoading(true)
    try {
      await AuthService.createContest(contest)
      await this.getAdmin()
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async getCommands(contestId: number) {
    this.setLoading(true)
    try {
      await AuthService.getCommands(contestId)
      await this.getAdmin()
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async setTeamLimit(contestId: number, min_teammates: number, max_teammates: number) {
    this.setLoading(true)
    try {
      await AuthService.setLimit(contestId, min_teammates, max_teammates)
      await this.getAdmin()
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async getCommandMembers(commandId: number) {
    this.setLoading(true)
    try {
      const commandMembers = await AuthService.getCommandMembers(commandId)
      return commandMembers
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }

  async changeContest(contestId: number, name: string, description: string, start_date: number, end_date: number) {
    this.setLoading(true)
    try {
      await AuthService.changeContest(contestId, name, description, start_date, end_date)
    } catch (e) {
      console.log(e)
    } finally {
      this.setLoading(false)
    }
  }
}
